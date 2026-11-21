package shutdown

import "testing"

func TestExecuteDryRun(t *testing.T) {
	t.Parallel()

	if err := Execute(true); err != nil {
		t.Fatalf("Execute(true) deve retornar nil em dry-run, obteve: %v", err)
	}

	if err := Shutdown(true); err != nil {
		t.Fatalf("Shutdown(true) deve retornar nil em dry-run, obteve: %v", err)
	}
}

func TestBeepDoesNotPanic(t *testing.T) {
	t.Parallel()

	// Beep roda em goroutine; valida que não ocorre panic com diferentes frequências.
	Beep(800, 20)
	Beep(600, 10)
	Beep(0, 0)
}
