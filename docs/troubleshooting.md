# 🔧 Guia de Solução de Problemas (Troubleshooting)

Este guia documenta o comportamento em baixo nível do `sleepoff` nos sistemas operacionais suportados (Linux, macOS e Windows), detalhando mecanismos de IPC, permissões de segurança de kernel, subsistemas de sessão e passos exatos para resolução de falhas.

---

## 1. macOS: Permissões de Automação e Apple Events

### Mecânica Interna
No macOS, o desligamento ordenado da máquina é delegado ao subsistema de eventos do sistema operacional através do AppleScript / OSA (Open Scripting Architecture):

```bash
osascript -e 'tell application "System Events" to shut down'
```

Essa chamada instrui o daemon `System Events` (`/System/Library/CoreServices/System Events.app`) a iniciar o encerramento seguro dos processos em espaço de usuário, salvando o estado das aplicações e executando o shutdown através da interface gráfica do macOS.

### Sintoma da Falha
A partir do macOS Mojave (10.14) e reforçado no macOS Catalina, Ventura, Sonoma e Sequoia, a Apple implementou o **TCC (Transparency, Consent, and Control)**. Se o emulador de terminal onde o `sleepoff` está rodando não possuir autorização prévia para enviar Apple Events ao `System Events`, o comando falha imediatamente:

```text
execution error: Not authorized to send Apple events to System Events. (-1743)
```
ou
```text
36:58: execution error: "System Events" got an error: User canceled. (-128)
```

### Resolução Passo a Passo

#### Via Interface Gráfica (Recomendado)
1. Abra as **Configurações do Sistema** (*System Settings*).
2. Navegue até **Privacidade e Segurança** (*Privacy & Security*) > **Automação** (*Automation*).
3. Localize o seu emulador de terminal na lista (ex.: `Terminal`, `iTerm2`, `Alacritty`, `kitty`, `WezTerm` ou `Visual Studio Code`).
4. Expanda a entrada do terminal e ative a chave ao lado de **System Events**.
5. Reinicie a sessão do terminal e teste com:
   ```bash
   sleepoff 10s --dry-run
   ```

#### Diagnóstico e Reset via CLI
Caso a permissão tenha sido negada inadvertidamente e o prompt de autorização não reapareça, você pode resetar o banco de dados do TCC para a categoria de Apple Events:

```bash
# Para o Terminal padrão da Apple:
tccutil reset AppleEvents com.apple.Terminal

# Para o iTerm2:
tccutil reset AppleEvents com.googlecode.iterm2

# Para o Alacritty:
tccutil reset AppleEvents org.alacritty

# Para o Kitty:
tccutil reset AppleEvents net.kovidgoyal.kitty
```

Após executar o comando, inicie o `sleepoff`. No momento do shutdown, o macOS exibirá uma caixa de diálogo nativa solicitando autorização para enviar eventos ao `System Events`. Clique em **Permitir** (*OK*).

> [!NOTE]
> Sessões remotas SSH puras ou daemons rodando desanexados de uma sessão de janela Aqua (login gráfico) não têm acesso ao contexto TCC do usuário logado. Nesses casos, o envio de Apple Events será rejeitado pelo kernel Darwin.

---

## 2. Linux: systemd, D-Bus e Permissões sem Sudo

### Mecânica Interna
No Linux moderno gerenciado por `systemd`, o `sleepoff` despacha a solicitação de desligamento através do comando:

```bash
systemctl poweroff
```

Internamente, o `systemctl` não executa diretamente uma syscall `sys_reboot(LINUX_REBOOT_CMD_POWER_OFF)`. Em vez disso, ele se conecta ao barramento de sistema **D-Bus** (`/run/dbus/system_bus_socket`) e invoca o método `PowerOff` na interface do serviço `systemd-logind`:

- **Destino D-Bus:** `org.freedesktop.login1`
- **Objeto D-Bus:** `/org/freedesktop/login1`
- **Método:** `org.freedesktop.login1.Manager.PowerOff(true)`

### Por que funciona sem sudo em Desktops?
O `systemd-logind` integra-se ao subsistema de privilégios **Polkit** (anteriormente PolicyKit). O Polkit avalia se a ação `org.freedesktop.login1.power-off` é permitida com base na política da sessão. Por padrão, se o usuário estiver em uma **sessão local ativa** conectada ao assento `seat0` (display local X11 ou Wayland), a ação recebe o status `yes` (permitido sem requisição de senha de root).

### Quando e por que a permissão é negada?
O erro ocorre normalmente em três cenários:

```text
Failed to power off system via logind: Access denied
Failed to call RPC: Access denied
```

1. **Sessão SSH Remota:** O usuário logado via SSH não reside no assento físico local (`seat0`), sendo classificado pelo logind como sessão inativa/remota.
2. **Múltiplos Usuários Conectados:** Se outro usuário estiver com sessão aberta simultaneamente, a política exige confirmação de privilégio administrativo para evitar perda de dados de terceiros (`org.freedesktop.login1.power-off-multiple-sessions`).
3. **Inibidores Ativos (Inhibitors):** Aplicações segurando travas de inibição de shutdown (`power-off-ignore-inhibit`).

### Resolução

#### Opção A: Regra Polkit (Recomendado, Seguro e sem Sudo)
Crie uma regra declarativa no diretório de políticas do Polkit para permitir que usuários do grupo `users` ou `wheel` executem o poweroff mesmo em sessões remotas ou concorrentes.

Crie o arquivo `/etc/polkit-1/rules.d/50-sleepoff-poweroff.rules`:

```javascript
/* Permite desligamento sem sudo via systemd-logind / Polkit */
polkit.addRule(function(action, subject) {
    if ((action.id == "org.freedesktop.login1.power-off" ||
         action.id == "org.freedesktop.login1.power-off-multiple-sessions" ||
         action.id == "org.freedesktop.login1.power-off-ignore-inhibit") &&
        (subject.isInGroup("users") || subject.isInGroup("wheel"))) {
        return polkit.Result.YES;
    }
});
```

Recarregue o daemon do polkit ou reinicie o serviço:
```bash
sudo systemctl restart polkit
```

#### Opção B: Ambientes sem systemd (SysVinit, OpenRC, runit)
Em distribuições como Alpine Linux, Devuan ou Void Linux sem `systemd`, o `sleepoff` recorre em cascata aos seguintes binários:

```bash
shutdown -h now
poweroff
```

Esses comandos exigem capacidade `CAP_SYS_BOOT` ou permissão de superusuário. Para permitir a execução sem senha para o seu usuário, adicione uma regra cirúrgica no `sudoers`:

Crie o arquivo `/etc/sudoers.d/sleepoff`:
```text
%wheel ALL=(ALL) NOPASSWD: /sbin/shutdown, /sbin/poweroff
```
Garanta que as permissões do arquivo estejam restritas a `0440`:
```bash
sudo chmod 0440 /etc/sudoers.d/sleepoff
```

---

## 3. Windows: Privilégios e Travas de Processos

### Mecânica Interna
No Windows, o `sleepoff` invoca o binário nativo do subsistema Win32:

```cmd
shutdown.exe /s /t 0
```

Este comando solicita o encerramento gracioso dos serviços e aplicações do Windows e desliga o computador imediatamente (`/t 0`).

### Possíveis Problemas e Resoluções

1. **`SeShutdownPrivilege` em Ambientes Corporativos (GPO):**
   - Se o usuário estiver sob uma conta corporativa com políticas restritivas do Active Directory onde o privilégio `SeShutdownPrivilege` ("Desligar o sistema") foi removido do grupo `Users`, o comando retornará:
     ```text
     Acesso negado. (1314) / Access is denied.
     ```
   - **Solução:** É necessário conceder o direito `SeShutdownPrivilege` localmente via `secpol.msc` (*Políticas Locais -> Atribuição de Direitos de Usuário -> Desligar o sistema*).

2. **Instância Única Já em Execução:**
   - O `sleepoff` impede duas instâncias simultâneas usando um Win32 Named Mutex (`Local\sleepoff_single_instance`). Se você tentar abrir uma nova instância com uma anterior ativa, ela emitirá um alerta e encerrará com código limpo.
   - Para encerrar uma instância órfã em segundo plano:
     ```powershell
     Stop-Process -Name sleepoff -Force
     ```

3. **Validação em Modo Seguro:**
   - Em caso de dúvida ou testes em scripts automatizados, sempre utilize a flag `--dry-run`:
     ```powershell
     sleepoff 10s --dry-run
     ```
   - No modo teste, todo o pipeline visual, sonoro e a tela de confirmação de 10s são executados, mas a chamada para o comando de desligamento é ignorada com segurança.
