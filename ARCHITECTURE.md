# 📊 Arquitetura do sleepoff

## Fluxograma de Estados

```
┌─────────────────────────────────────────────────────────────────────┐
│                           INÍCIO                                     │
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
                          │  [c] cancelar       │
                          └─────────┬───────────┘
                                    │
                    ┌───────────────┼───────────────┐
                    │               │               │
                    ▼               ▼               ▼
           ┌──────────────┐ ┌──────────────┐ ┌──────────────┐
           │   PAUSADO    │ │ TEMPO ACABA  │ │  CANCELADO   │
           │  (amarelo)   │ │              │ │              │
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

## Estados do App

| Estado | Descrição | Teclas |
|--------|-----------|--------|
| `StateSplash` | Animação inicial (2s) | Automático |
| `StateMenu` | Seleção de tempo | ↑↓ Enter Q |
| `StateCustomInput` | Input de minutos | Enter Esc |
| `StateRunning` | Timer ativo | P + - C Esc |
| `StateConfirmation` | Pânico (10s) | C Esc |
| `StateDone` | Finalizado | - |

## Estrutura de Arquivos

```
sleepoff/
├── main.go              # CLI, argumentos, Cobra
├── .goreleaser.yml      # Release zip e checksums
├── RELEASING.md         # Fluxo de corte de release
├── packaging/
│   └── windows/
│       ├── sleepoff.iss # Instalador Inno Setup
│       └── ...          # Icone e metadata de release
├── scripts/
│   ├── generate-resource.ps1
│   ├── install-local.ps1
│   └── package-installer.ps1
├── internal/
│   ├── config/
│   │   └── config.go    # Cores, opções de tempo, constantes
│   ├── buildinfo/
│   │   └── buildinfo.go # Versão, repo e metadata
│   ├── model/
│   │   ├── model.go     # State machine, tipos, constructors
│   │   ├── update.go    # Lógica de Update (Bubble Tea)
│   │   └── view.go      # Renderização de cada tela
│   ├── shutdown/
│   │   └── shutdown.go  # Beep, Execute shutdown
│   └── ui/
│       └── styles.go    # Estilos Lipgloss
├── README.md
├── TODO.md
└── sleepoff.bat         # Launcher para duplo clique
```

## Fluxo de Dados (TEA Pattern)

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

## Comandos de Teste

```bash
# Testar help (binário local)
.\sleepoff.exe --help

# Testar versão
.\sleepoff.exe --version

# Testar menu interativo (dry-run)
.\sleepoff.exe --dry-run

# Testar CLI direto
.\sleepoff.exe 1m --dry-run

# Testar erro de input
.\sleepoff.exe abc

# Instalar a build local no PATH do usuário
.\scripts\install-local.ps1 -Build
sleepoff 90s --dry-run
```
