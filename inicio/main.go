package main

import (
	"fmt" // Formatação e exibição de texto no terminal
	"os"  // Interação com o sistema operacional

	decodificador "simriscv/internal/isa"
	carregador "simriscv/internal/loader"
	maquina "simriscv/internal/machine"
)

func main() {
	if erro := executar(); erro != nil {
		fmt.Fprintln(os.Stderr, "Erro:", erro)
		os.Exit(1)
	}
}

// Demonstra os componentes do simulador usando o programa inicial em arquivo.
func executar() error {
	arquivo, erro := os.Open("programs/instrucoes-iniciais.hex")
	if erro != nil {
		return fmt.Errorf("abrir programa: %w", erro)
	}
	defer arquivo.Close()

	palavras, erro := carregador.LerHex(arquivo)
	if erro != nil {
		return erro
	}
	computador, erro := maquina.Nova(palavras)
	if erro != nil {
		return erro
	}

	for !computador.Concluida() {
		antes := computador.Registradores()
		passo, erro := computador.Passo()
		if erro != nil {
			return erro
		}
		// A instrução já executada é decodificada aqui apenas para exibição.
		instrucao, erro := decodificador.Decodificar(passo.Palavra)
		if erro != nil {
			return erro
		}
		fmt.Printf("PC=%d | palavra=0x%08x\n", passo.PC, passo.Palavra)
		switch instrucao.Operacao {
		case decodificador.ADDI:
			fmt.Printf("  ADDI x%d, x%d, %d\n", instrucao.Destino, instrucao.Fonte1, instrucao.ImediatoI)
			fmt.Printf("  Cálculo: %d + (%d)\n", int32(antes[instrucao.Fonte1]), instrucao.ImediatoI)
		case decodificador.ADD:
			fmt.Printf("  ADD x%d, x%d, x%d\n", instrucao.Destino, instrucao.Fonte1, instrucao.Fonte2)
			fmt.Printf("  Cálculo: %d + %d\n", int32(antes[instrucao.Fonte1]), int32(antes[instrucao.Fonte2]))
		}
		depois := computador.Registradores()
		fmt.Printf("  Registrador x%d: %d --> %d\n", instrucao.Destino, int32(antes[instrucao.Destino]), int32(depois[instrucao.Destino]))
		if instrucao.Destino == 0 {
			fmt.Println("  Escrita descartada: x0 permanece zero.")
		}
		fmt.Printf("  PC: %d --> %d\n\n", passo.PC, passo.ProximoPC)
	}
	fmt.Printf("PC final = %d\n", computador.PC())
	return nil
}
