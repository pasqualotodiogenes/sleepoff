package windowsux

import (
	"testing"
)

func TestSingleInstanceLockExclusive(t *testing.T) {
	// Primeiro lock: deve adquirir com sucesso
	lock1, alreadyRunning1, err1 := Lock()
	if err1 != nil {
		t.Fatalf("Lock 1 falhou com erro: %v", err1)
	}
	if alreadyRunning1 {
		t.Fatalf("Lock 1 retornou alreadyRunning = true, esperado false")
	}
	if lock1 == nil {
		t.Fatalf("Lock 1 retornou instância nil")
	}

	// Segundo lock concorrente: deve detectar conflito (alreadyRunning = true)
	lock2, alreadyRunning2, err2 := Lock()
	if err2 != nil {
		lock1.Release()
		t.Fatalf("Lock 2 falhou com erro inesperado: %v", err2)
	}
	if !alreadyRunning2 {
		if lock2 != nil {
			lock2.Release()
		}
		lock1.Release()
		t.Fatalf("Lock 2 concorrente esperado alreadyRunning = true, obteve false")
	}
	if lock2 != nil {
		lock2.Release()
		lock1.Release()
		t.Fatalf("Lock 2 concorrente deveria retornar lock nil quando alreadyRunning = true")
	}

	// Libera lock 1
	lock1.Release()

	// Terceiro lock após release: deve adquirir com sucesso novamente
	lock3, alreadyRunning3, err3 := Lock()
	if err3 != nil {
		t.Fatalf("Lock 3 pós-release falhou com erro: %v", err3)
	}
	if alreadyRunning3 {
		t.Fatalf("Lock 3 pós-release retornou alreadyRunning = true, esperado false")
	}
	if lock3 == nil {
		t.Fatalf("Lock 3 pós-release retornou instância nil")
	}
	lock3.Release()
}

func TestAcquireSingleInstanceAlias(t *testing.T) {
	lock, alreadyRunning, err := AcquireSingleInstance()
	if err != nil {
		t.Fatalf("AcquireSingleInstance() falhou: %v", err)
	}
	if alreadyRunning {
		t.Fatalf("AcquireSingleInstance() retornou alreadyRunning = true")
	}
	if lock != nil {
		lock.Release()
	}
}
