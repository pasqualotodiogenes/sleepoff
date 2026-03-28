package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/pasqualotodiogenes/sleepoff/internal/buildinfo"
	"github.com/pasqualotodiogenes/sleepoff/internal/model"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

var (
	version = buildinfo.VersionString()
	dryRun  bool
)

// Cores para output bonito
var (
	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#DA70D6")).
			Bold(true)

	subtitleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#888888"))

	successStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#00CED1"))

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF6B6B"))

	dimStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#666666"))

	accentStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#DA70D6"))
)

func printBanner() {
	banner := `
   ███████╗██╗     ███████╗███████╗██████╗  ██████╗ ███████╗███████╗
   ██╔════╝██║     ██╔════╝██╔════╝██╔══██╗██╔═══██╗██╔════╝██╔════╝
   ███████╗██║     █████╗  █████╗  ██████╔╝██║   ██║█████╗  █████╗  
   ╚════██║██║     ██╔══╝  ██╔══╝  ██╔═══╝ ██║   ██║██╔══╝  ██╔══╝  
   ███████║███████╗███████╗███████╗██║     ╚██████╔╝██║     ██║     
   ╚══════╝╚══════╝╚══════╝╚══════╝╚═╝      ╚═════╝ ╚═╝     ╚═╝     
`
	fmt.Println(titleStyle.Render(banner))
	fmt.Println(subtitleStyle.Render("   Timer de desligamento automático • v" + version))
	fmt.Println()
}

func printHelp() {
	printBanner()

	fmt.Println(titleStyle.Render("   COMANDOS"))
	fmt.Println()
	fmt.Printf("   %s              %s\n",
		accentStyle.Render("sleepoff"),
		dimStyle.Render("Abre o menu interativo"))
	fmt.Printf("   %s %s      %s\n",
		accentStyle.Render("sleepoff"),
		successStyle.Render("<tempo>"),
		dimStyle.Render("Inicia timer direto"))
	fmt.Println()

	fmt.Println(titleStyle.Render("   FORMATOS DE TEMPO"))
	fmt.Println()
	fmt.Printf("   %s  %s  %s  %s\n",
		successStyle.Render("30"),
		dimStyle.Render("30 minutos"),
		successStyle.Render("1h"),
		dimStyle.Render("1 hora"))
	fmt.Printf("   %s  %s  %s  %s\n",
		successStyle.Render("30m"),
		dimStyle.Render("30 minutos"),
		successStyle.Render("90s"),
		dimStyle.Render("90 segundos"))
	fmt.Printf("   %s  %s\n",
		successStyle.Render("1h30m"),
		dimStyle.Render("1 hora e 30 minutos"))
	fmt.Println()

	fmt.Println(titleStyle.Render("   OPÇÕES"))
	fmt.Println()
	fmt.Printf("   %s        %s\n",
		accentStyle.Render("--dry-run"),
		dimStyle.Render("Modo teste (não desliga de verdade)"))
	fmt.Printf("   %s %s   %s\n",
		accentStyle.Render("-h"),
		accentStyle.Render("--help"),
		dimStyle.Render("Mostra esta ajuda"))
	fmt.Printf("   %s %s  %s\n",
		accentStyle.Render("-v"),
		accentStyle.Render("--version"),
		dimStyle.Render("Mostra a versão"))
	fmt.Println()

	fmt.Println(titleStyle.Render("   EXEMPLOS"))
	fmt.Println()
	fmt.Printf("   %s\n", dimStyle.Render("# Timer de 30 minutos"))
	fmt.Printf("   %s\n", successStyle.Render("sleepoff 30m"))
	fmt.Println()
	fmt.Printf("   %s\n", dimStyle.Render("# Timer de 1 hora (modo teste)"))
	fmt.Printf("   %s\n", successStyle.Render("sleepoff 1h --dry-run"))
	fmt.Println()
	fmt.Printf("   %s\n", dimStyle.Render("# Menu interativo"))
	fmt.Printf("   %s\n", successStyle.Render("sleepoff"))
	fmt.Println()

	fmt.Println(titleStyle.Render("   EXECUÇÃO NO WINDOWS"))
	fmt.Println()
	fmt.Printf("   %s  %s\n",
		accentStyle.Render("Instalado no PATH:"),
		successStyle.Render("sleepoff 90s"))
	fmt.Printf("   %s  %s\n",
		accentStyle.Render("Binário local (PowerShell):"),
		successStyle.Render(".\\sleepoff.exe 90s"))
	fmt.Println(dimStyle.Render("   O PowerShell não executa arquivos da pasta atual sem .\\ por padrão."))
	fmt.Println()

	fmt.Println(titleStyle.Render("   CONTROLES (durante o timer)"))
	fmt.Println()
	fmt.Printf("   %s  %s\n", accentStyle.Render("[p] [espaço]"), dimStyle.Render("Pausar/Retomar"))
	fmt.Printf("   %s  %s\n", accentStyle.Render("[+] [-]"), dimStyle.Render("Ajustar ±5 minutos"))
	fmt.Printf("   %s  %s\n", accentStyle.Render("[c] [esc]"), dimStyle.Render("Cancelar"))
	fmt.Println()
}

var rootCmd = &cobra.Command{
	Use:     "sleepoff [tempo]",
	Short:   "Timer de desligamento automático",
	Version: version,
	Args:    cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// Se pediu ajuda explícita
		if h, _ := cmd.Flags().GetBool("help"); h {
			printHelp()
			return
		}

		var m model.Model

		if len(args) == 1 {
			// Modo CLI direto
			duration, err := parseDuration(args[0])
			if err != nil {
				fmt.Println()
				fmt.Println(errorStyle.Render("   ✗ Duração inválida: " + args[0]))
				fmt.Println()
				fmt.Println(dimStyle.Render("   Formatos aceitos: 30, 30m, 90s, 1h, 1h30m"))
				fmt.Println(dimStyle.Render("   Exemplo: sleepoff 30m"))
				fmt.Println()
				os.Exit(1)
			}

			// Feedback visual
			fmt.Println()
			fmt.Printf("   %s Timer de %s iniciado\n",
				successStyle.Render("✓"),
				successStyle.Render(formatDuration(duration)))
			if dryRun {
				fmt.Printf("   %s Modo teste ativado\n", dimStyle.Render("ℹ"))
			}
			fmt.Println()

			m = model.NewWithDuration(duration, dryRun)
		} else {
			// Modo interativo
			m = model.New()
			m.DryRun = dryRun
		}

		p := tea.NewProgram(m, tea.WithAltScreen())
		finalModel, err := p.Run()
		if err != nil {
			fmt.Println(errorStyle.Render("   ✗ Erro: " + err.Error()))
			os.Exit(1)
		}

		fm := finalModel.(model.Model)

		// Mensagem de saída
		if fm.State == model.StateRunning {
			// Cancelou durante o timer
			showCancelSummary(fm)
		} else if fm.ShutdownError != "" {
			showShutdownError(fm.ShutdownError)
		} else if fm.State == model.StateConfirmation || fm.State == model.StateDone {
			// Timer terminou
			showFinalSuccess(dryRun)
		} else if fm.Quitting {
			// Saiu do menu
			showGoodbye()
		}
	},
}

func init() {
	// Permite executar por duplo clique sem o aviso padrão do Cobra no Windows.
	cobra.MousetrapHelpText = ""

	rootCmd.Flags().BoolVar(&dryRun, "dry-run", false, "Modo teste (não desliga)")
	rootCmd.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		printHelp()
	})
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func parseDuration(s string) (time.Duration, error) {
	if n, err := strconv.Atoi(s); err == nil {
		if n <= 0 {
			return 0, errors.New("duração deve ser maior que zero")
		}
		return time.Duration(n) * time.Minute, nil
	}
	d, err := time.ParseDuration(s)
	if err != nil {
		return 0, err
	}
	if d <= 0 {
		return 0, errors.New("duração deve ser maior que zero")
	}
	return d, nil
}

func formatDuration(d time.Duration) string {
	h := int(d.Hours())
	m := int(d.Minutes()) % 60
	s := int(d.Seconds()) % 60

	if h > 0 && m > 0 {
		return fmt.Sprintf("%dh%dm", h, m)
	} else if h > 0 {
		return fmt.Sprintf("%dh", h)
	} else if m > 0 {
		return fmt.Sprintf("%dm", m)
	}
	return fmt.Sprintf("%ds", s)
}

func showCancelSummary(m model.Model) {
	// Limpa terminal
	fmt.Print("\033[H\033[2J")

	elapsed := m.TotalDuration - m.Remaining
	totalMins := int(m.TotalDuration.Minutes())
	elapsedMins := int(elapsed.Minutes())
	elapsedSecs := int(elapsed.Seconds()) % 60
	percentage := 0
	if m.TotalDuration.Seconds() > 0 {
		percentage = int((elapsed.Seconds() / m.TotalDuration.Seconds()) * 100)
	}

	// ASCII art grande
	cancelArt := `
   ╔═══════════════════════════════════════════════════════════════╗
   ║                                                               ║
   ║      ████████╗██╗███╗   ███╗███████╗██████╗                   ║
   ║      ╚══██╔══╝██║████╗ ████║██╔════╝██╔══██╗                  ║
   ║         ██║   ██║██╔████╔██║█████╗  ██████╔╝                  ║
   ║         ██║   ██║██║╚██╔╝██║██╔══╝  ██╔══██╗                  ║
   ║         ██║   ██║██║ ╚═╝ ██║███████╗██║  ██║                  ║
   ║         ╚═╝   ╚═╝╚═╝     ╚═╝╚══════╝╚═╝  ╚═╝                  ║
   ║                                                               ║
   ║                     ⏹  CANCELADO                              ║
   ║                                                               ║
   ╚═══════════════════════════════════════════════════════════════╝
`
	fmt.Print(errorStyle.Render(cancelArt))
	fmt.Println()

	// Barra de progresso grande
	barLen := 50
	filled := (percentage * barLen) / 100
	if filled > barLen {
		filled = barLen
	}
	bar := ""
	for i := 0; i < barLen; i++ {
		if i < filled {
			bar += "█"
		} else {
			bar += "░"
		}
	}

	fmt.Printf("   Progresso: %s %d%%\n", successStyle.Render(bar), percentage)
	fmt.Println()

	// Stats em box
	fmt.Println(dimStyle.Render("   ┌─────────────────────────────────────┐"))
	fmt.Printf("   │  %s  %s%s│\n",
		dimStyle.Render("Tempo decorrido:"),
		successStyle.Render(fmt.Sprintf("%02dm %02ds", elapsedMins, elapsedSecs)),
		"              ")
	fmt.Printf("   │  %s  %s%s│\n",
		dimStyle.Render("Tempo total:    "),
		accentStyle.Render(fmt.Sprintf("%02dm", totalMins)),
		"                   ")
	fmt.Println(dimStyle.Render("   └─────────────────────────────────────┘"))
	fmt.Println()

	// Mensagem motivacional grande
	var msg, emoji string
	switch {
	case percentage >= 80:
		emoji = "🏆"
		msg = "QUASE LÁ! Tão perto do final!"
	case percentage >= 50:
		emoji = "💪"
		msg = "PASSOU DA METADE! Respeitável."
	case percentage >= 20:
		emoji = "🤔"
		msg = "Bem... pelo menos tentou."
	default:
		emoji = "😅"
		msg = "Desistiu rápido hein?"
	}
	fmt.Printf("   %s  %s\n", emoji, titleStyle.Render(msg))
	fmt.Println()
	fmt.Println()
	fmt.Println(dimStyle.Render("   ─────────────────────────────────────────────────────────────────"))
	fmt.Println(dimStyle.Render("   sleepoff v" + version + " • " + buildinfo.RepoRef))
	fmt.Println()
}

func showFinalSuccess(dryRun bool) {
	// Limpa terminal
	fmt.Print("\033[H\033[2J")

	if dryRun {
		successArt := `
   ╔═══════════════════════════════════════════════════════════════╗
   ║                                                               ║
   ║      ████████╗██╗███╗   ███╗███████╗██████╗                   ║
   ║      ╚══██╔══╝██║████╗ ████║██╔════╝██╔══██╗                  ║
   ║         ██║   ██║██╔████╔██║█████╗  ██████╔╝                  ║
   ║         ██║   ██║██║╚██╔╝██║██╔══╝  ██╔══██╗                  ║
   ║         ██║   ██║██║ ╚═╝ ██║███████╗██║  ██║                  ║
   ║         ╚═╝   ╚═╝╚═╝     ╚═╝╚══════╝╚═╝  ╚═╝                  ║
   ║                                                               ║
   ║                  ✓  FINALIZADO (TESTE)                        ║
   ║                                                               ║
   ╚═══════════════════════════════════════════════════════════════╝
`
		fmt.Print(successStyle.Render(successArt))
		fmt.Println()
		fmt.Println(dimStyle.Render("   ℹ  O computador NÃO foi desligado."))
		fmt.Println(dimStyle.Render("   ℹ  Use sem --dry-run para desligar de verdade."))
	} else {
		shutdownArt := `
   ╔═══════════════════════════════════════════════════════════════╗
   ║                                                               ║
   ║               ⚡ DESLIGANDO AGORA... ⚡                        ║
   ║                                                               ║
   ╚═══════════════════════════════════════════════════════════════╝
`
		fmt.Print(errorStyle.Render(shutdownArt))
	}
	fmt.Println()
	fmt.Println()
	fmt.Println(dimStyle.Render("   ─────────────────────────────────────────────────────────────────"))
	fmt.Println(dimStyle.Render("   sleepoff v" + version + " • " + buildinfo.RepoRef))
	fmt.Println()
}

func showGoodbye() {
	// Limpa terminal
	fmt.Print("\033[H\033[2J")

	goodbyeArt := `
   ╔═══════════════════════════════════════════════════════════════╗
   ║                                                               ║
   ║       ██████╗ ██╗   ██╗███████╗██╗                            ║
   ║       ██╔══██╗╚██╗ ██╔╝██╔════╝██║                            ║
   ║       ██████╔╝ ╚████╔╝ █████╗  ██║                            ║
   ║       ██╔══██╗  ╚██╔╝  ██╔══╝  ╚═╝                            ║
   ║       ██████╔╝   ██║   ███████╗██╗                            ║
   ║       ╚═════╝    ╚═╝   ╚══════╝╚═╝                            ║
   ║                                                               ║
   ║                     👋 Até mais!                              ║
   ║                                                               ║
   ╚═══════════════════════════════════════════════════════════════╝
`
	fmt.Print(titleStyle.Render(goodbyeArt))
	fmt.Println()
	fmt.Println()
	fmt.Println(dimStyle.Render("   ─────────────────────────────────────────────────────────────────"))
	fmt.Println(dimStyle.Render("   sleepoff v" + version + " • " + buildinfo.RepoRef))
	fmt.Println()
}

func showShutdownError(errMsg string) {
	fmt.Print("\033[H\033[2J")

	fmt.Println()
	fmt.Println(errorStyle.Render("   ✗ Falha ao executar shutdown"))
	fmt.Println()
	fmt.Println(dimStyle.Render("   Detalhes: " + errMsg))
	fmt.Println()
	fmt.Println(dimStyle.Render("   Verifique permissões, políticas do Windows ou use --dry-run para testar o fluxo."))
	fmt.Println()
}
