//go:build linux

package windowsux

import (
	"fmt"
	"os"
	"os/exec"
)

// ApplyProcessIdentity é no-op no Linux.
func ApplyProcessIdentity() error {
	return nil
}

// ShowAlreadyRunningMessage avisa no terminal/stderr quando o sleepoff já estiver rodando.
func ShowAlreadyRunningMessage() {
	fmt.Fprintln(os.Stderr, "   ✗ sleepoff já está em execução.")
	fmt.Fprintln(os.Stderr, "   Use a instância atual pelo terminal ou encerre-a antes de iniciar outra.")
}

// ShowUpdateInstallingMessage imprime aviso de atualização no Linux.
func ShowUpdateInstallingMessage(version string) {
	fmt.Fprintf(os.Stderr, "   ℹ Atualização encontrada: v%s\n", version)
}

// ShowInfo exibe mensagem informativa no terminal/stderr.
func ShowInfo(title, message string) {
	fmt.Fprintf(os.Stderr, "   ℹ %s: %s\n", title, message)
}

// OpenURL tenta abrir uma URL no navegador padrão usando xdg-open.
func OpenURL(url string) error {
	return exec.Command("xdg-open", url).Start()
}
