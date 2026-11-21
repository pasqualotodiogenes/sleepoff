# 📊 Arquitetura do sleepoff

Documento técnico de engenharia de software detalhando a arquitetura interna, o modelo de estados, a camada de abstração multiplataforma via Go build tags, a estratégia de compilação sem CGO e os fundamentos de sistemas operacionais (POSIX vs. Win32).

---

## 1. Princípios de Engenharia e Design

O `sleepoff` segue uma filosofia estrita de fundamentos de sistemas:
- **Zero-CGO (`CGO_ENABLED=0`):** Binários 100% estáticos, sem acoplamento a bibliotecas dinâmicas do sistema C (`glibc`, `musl`, `libSystem`). Isso viabiliza cross-compilação trivial no pipeline CI e portabilidade universal em qualquer distribuição Linux ou versão do macOS/Windows.
- **Isolamento por Build Tags:** Ausência de condicionais dinâmicas de sistema operacional (`if runtime.GOOS == "windows"`) no core da aplicação (`internal/model`). Cada subsistema de baixo nível é isolado em arquivos com tags de compilação (`_windows.go`, `_linux.go`, `_darwin.go`, `_posix.go`, `_other.go`).
- **Arquitetura TEA (The Elm Architecture):** Unidirecionalidade estrita no fluxo de dados através do framework [Bubble Tea](https://github.com/charmbracelet/bubbletea), garantindo previsibilidade de estado, renderização sem flickering e teste unitário determinístico.

---

## 2. Camada de Abstração Multiplataforma

A tabela abaixo sintetiza como cada primitiva de sistema operacional é implementada nativamente nas plataformas suportadas:

| Subsistema / Recurso | Windows (`_windows.go`) | Linux (`_linux.go`) | macOS (`_darwin.go`) |
| :--- | :--- | :--- | :--- |
| **Execução de Desligamento** | Invocação Win32 do executável nativo `shutdown.exe /s /t 0` | IPC D-Bus com `systemd-logind` via `systemctl poweroff` (fallbacks: `shutdown -h now`, `poweroff`) | Apple Events via motor OSA: `osascript -e 'tell app "System Events" to shut down'` |
| **Instância Única (Lock)** | Win32 Named Mutex no namespace de sessão: `Local\sleepoff_single_instance` | POSIX advisory lock via `syscall.Flock(fd, LOCK_EX\|LOCK_NB)` em `$XDG_RUNTIME_DIR/sleepoff.lock` | POSIX advisory lock via `syscall.Flock(fd, LOCK_EX\|LOCK_NB)` em `/tmp/sleepoff.lock` |
| **Feedback Sonoro (Beep)** | Chamada síncrona nativa `kernel32.dll` `Beep(dwFreq, dwDuration)` | Reprodução encadeada de pipelines de áudio do usuário: `paplay` (PulseAudio), `pw-play` (PipeWire), `aplay` (ALSA), fallback: `\a` | Utilitário nativo Darwin `afplay /System/Library/Sounds/Tink.aiff`, fallback: `\a` |
| **Integração Desktop / Tray** | Notificações ricas em System Tray e menus de contexto via `lxn/walk` | Headless graceful no-op (evita dependências pesadas de bibliotecas gráficas C) | Headless graceful no-op (evita dependências pesadas de bibliotecas gráficas C) |
| **Auto-Atualização** | Verificação de release no GitHub e bootstrap do instalador Inno Setup | Delegado a gerenciadores de pacote do sistema ou atualização do binário | Delegado a gerenciadores de pacote do sistema ou atualização do binário |

---

## 3. Mecânica em Baixo Nível: POSIX vs. Win32

### A. Trava de Instância Única (Single Instance Lock)
Para prevenir condições de corrida em que múltiplos timers disputam o desligamento da máquina ou sobrecarregam o terminal:

- **No Windows (Win32):**
  Usa-se o subsistema de sincronização de objetos do kernel NT via `CreateMutexW`:
  ```go
  handle, err := syscall.CreateMutex(nil, false, syscall.StringToUTF16Ptr("Local\\sleepoff_single_instance"))
  ```
  Se o erro retornado for `ERROR_ALREADY_EXISTS` (`183`), o app sabe que outro processo já segura o mutex e aborta graciosamente.

- **No Linux e macOS (POSIX):**
  Utiliza-se descritores de arquivo combinados com a syscall `flock(2)` em modo exclusivo e não-bloqueante (`LOCK_EX | LOCK_NB`):
  ```go
  fd, err := syscall.Open(lockPath, syscall.O_CREAT|syscall.O_RDWR, 0600)
  err = syscall.Flock(fd, syscall.LOCK_EX|syscall.LOCK_NB)
  ```
  O kernel do Linux/BSD mantém a trava vinculada à tabela de arquivos abertos da tabela de processos. Se o processo morrer (por `SIGKILL`, crash ou encerramento normal), o kernel libera o lock automaticamente ao fechar o descritor de arquivo, eliminando o risco de "stale locks" persistentes em disco.

### B. Desligamento do Sistema (Poweroff Subsystems)
- **Windows:** Disparo direto do utilitário `shutdown.exe /s /t 0`, solicitando encerramento imediato ao subsistema `csrss.exe` / `winlogon.exe`.
- **Linux:** O sleepoff prioriza `systemctl poweroff`. Isso ativa o canal D-Bus com `systemd-logind` (`org.freedesktop.login1`). O Polkit autoriza sessões de desktop locais no assento `seat0` sem necessidade de privilégios de superusuário (`sudo`). Caso o host não execute `systemd`, o sleepoff tenta em cascata `shutdown -h now` e `poweroff`.
- **macOS:** Emite uma instrução OSA para o processo `System Events.app`. Este evento é intermediado pelo subsistema TCC da Apple, garantindo que aplicações com janelas abertas possam encerrar graciosamente ou salvar documentos não gravados.

---

## 4. Máquina de Estados Finita (FSM) e Fluxo de Execução

```
┌─────────────────────────────────────────────────────────────────────┐
│                           INÍCIO                                    │
└────────────────────────────┬────────────────────────────────────────┘
                             │
                             ▼
                    ┌────────────────┐
                    │  Argumentos?   │
                    └───────┬────────┘
                            │
              ┌─────────────┴─────────────┐
              │                           │
              ▼                           ▼
    ┌─────────────────┐         ┌─────────────────┐
    │  sleepoff 30m   │         │    sleepoff     │
    │  (CLI direto)   │         │  (interativo)   │
    └────────┬────────┘         └────────┬────────┘
             │                           │
             │                           ▼
             │                  ┌─────────────────┐
             │                  │  SPLASH SCREEN  │
             │                  │   (2 segundos)  │
             │                  └────────┬────────┘
             │                           │
             │                           ▼
             │                  ┌─────────────────┐
             │                  │      MENU       │
             │                  │ ↑↓ navegação    │
             │                  │ Enter seleciona │
             │                  │ Q sai           │
             │                  └───────┬─────────┘
             │                          │
             │         ┌────────────────┼────────────────┐
             │         │                │                │
             │         ▼                ▼                ▼
             │  ┌──────────┐    ┌──────────┐    ┌──────────────┐
             │  │ 15/30/45 │    │  1h/2h   │    │ Personalizar │
             │  │   min    │    │          │    │  (input)     │
             │  └────┬─────┘    └────┬─────┘    └──────┬───────┘
             │       │               │                 │
             └───────┴───────────────┴─────────────────┘
                                     │
                                     ▼
                          ┌─────────────────────┐
                          │    TIMER RODANDO    │
                          │                     │
                          │  [p] pausar         │
                          │  [+] +5 min         │
                          │  [-] -5 min         │
                          │  [c] cancelar (2x)  │
                          └─────────┬───────────┘
                                    │
                    ┌───────────────┼───────────────┐
                    │               │               │
                    ▼               ▼               ▼
           ┌──────────────┐ ┌──────────────┐ ┌──────────────┐
           │   PAUSADO    │ │ TEMPO ACABA  │ │  CANCELADO   │
           │  (amarelo)   │ │              │ │  (stats box) │
           └──────┬───────┘ └──────┬───────┘ └──────┬───────┘
                  │                │                │
                  │                ▼                │
                  │     ┌─────────────────┐         │
                  │     │  TELA DE PÂNICO │         │
                  │     │   10 segundos   │         │
                  │     │   [c] cancelar  │         │
                  │     └────────┬────────┘         │
                  │              │                  │
                  │    ┌─────────┴─────────┐        │
                  │    ▼                   ▼        │
                  │ ┌────────┐      ┌──────────┐    │
                  │ │ TIMER  │      │ SHUTDOWN │    │
                  │ │ FINAL  │      │ (real)   │    │
                  │ └───┬────┘      └──────────┘    │
                  │     │                           │
                  └─────┴───────────────────────────┘
                                    │
                                    ▼
                          ┌─────────────────┐
                          │    RESUMO       │
                          │ "Aguentou X%"   │
                          └─────────────────┘
```

### Detalhamento dos Estados (`model.State`)

| Constante de Estado | Descrição Funcional | Transições e Eventos de Teclado |
| :--- | :--- | :--- |
| `StateSplash` | Splash screen com animação de loading progressivo (2000ms). | Transita automaticamente para `StateMenu` via `TickMsg`. |
| `StateMenu` | Menu interativo de seleção de intervalos predefinidos. | `↑`/`↓`/`k`/`j` navegam; `Enter` escolhe; `q` encerra o app. |
| `StateCustomInput` | Prompt de digitação numérica de duração customizada em minutos. | `Enter` valida valor positivo e inicia timer; `Esc` retorna ao menu. |
| `StateRunning` | Timer ativo com barra de gradiente, logs de eventos e tempo final. | `p`/`Espaço` pausa/retoma; `+`/`-` ajusta ±5m; `c`/`q`/`Esc` (2x em 2s) cancela. |
| `StateConfirmation` | Tela de pânico com contagem regressiva de 10s e beeps sonoros progressivos. | `c` ou `Esc` (1x) aborta o desligamento; timeout executa o shutdown. |
| `StateDone` | Encerramento do ciclo e apresentação da tela final ASCII de resumo. | Dispara encerramento do programa (`tea.Quit`). |

---

## 5. Fluxo de Dados (TEA Pattern - The Elm Architecture)

```
┌─────────────────────────────────────────────────────────────┐
│                      Bubble Tea Loop                        │
└─────────────────────────────────────────────────────────────┘
                               │
         ┌─────────────────────┼─────────────────────┐
         │                     │                     │
         ▼                     ▼                     ▼
    ┌─────────┐          ┌──────────┐          ┌──────────┐
    │  Init   │          │  Update  │          │   View   │
    │         │◄─────────│          │◄─────────│          │
    │ Retorna │   Msg    │ Processa │   Model  │ Renderiza│
    │   Cmd   │          │   Msgs   │          │   Tela   │
    └─────────┘          └──────────┘          └──────────┘
         │                     │                     │
         │                     │                     │
         ▼                     ▼                     ▼
    TickMsg              KeyMsg              String (UI)
    (50ms)               WindowSizeMsg       
```

- **Ciclo de Tick:** Um ticker de 50ms (`time.Millisecond * 50`) mantém a barra de progresso suave e avalia o deadline do timer contra `time.Now()`.
- **Prevenção de Drifting:** A contagem regressiva não subtrai frações de segundo iterativamente. Ela calcula `m.Remaining = time.Until(m.FinishTime)`. Isso garante exatidão cronométrica mesmo se a CPU estiver sob carga pesada ou o processo sofrer atraso de escalonamento no sistema operacional.
- **Cancelamento Anti-Acidente:** O cancelamento no `StateRunning` exige confirmação dupla em 2 segundos (`CancelPending`), impedindo que toques acidentais no teclado interrompam um timer longo.

---

## 6. Estrutura de Código e Organização de Pacotes

```
sleepoff/
├── main.go                     # Ponto de entrada CLI, Cobra, tratamento de flags
├── main_test.go                # Testes de parsing e integração CLI
├── .goreleaser.yml             # Pipeline de empacotamento multiplataforma (zip/tar.gz)
├── .github/workflows/          # CI e automação de release
├── docs/                       # Documentação e ativos
│   ├── assets/                 # SVGs de logotipos, badges e previews
│   └── troubleshooting.md      # Guia de solução de permissões TCC e D-Bus
├── packaging/windows/          # Scripts Inno Setup, manifesto e ícone Win32
├── internal/
│   ├── buildinfo/              # Variáveis injetadas via ldflags (-X)
│   ├── config/                 # Constantes de paleta de cores e presets de tempo
│   ├── model/                  # State machine pura, Bubble Tea model, view e update
│   ├── shutdown/               # Abstração de desligamento e som por build tags
│   │   ├── shutdown.go         # Interface pública e fallback
│   │   ├── shutdown_windows.go # Win32 API e kernel32 Beep
│   │   ├── shutdown_linux.go   # D-Bus / systemd logind e reprodutores de áudio
│   │   └── shutdown_darwin.go  # Apple Events (osascript) e afplay
│   ├── ui/                     # Definições de estilos Lipgloss
│   ├── update/                 # Verificação de atualizações no GitHub
│   └── windowsux/              # Integração de desktop (Tray via lxn/walk e locks)
├── README.md                   # Apresentação do projeto e guia do usuário
└── ARCHITECTURE.md             # Esta especificação técnica
```
