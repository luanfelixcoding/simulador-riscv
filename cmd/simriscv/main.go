package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	carregador "simriscv/internal/loader"
	maquina "simriscv/internal/machine"
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Uso: simriscv [opções]\nOpções:\n")
		flag.PrintDefaults()
	}
	programa := flag.String("programa", "programs/instrucoes-iniciais.hex", "arquivo .hex com instruções de 32 bits por linha")
	maximoPassos := flag.Int("max-passos", 100, "limite de passos (maior que zero)")
	intervalo := flag.Int("video-a-cada", 10, "exibir VRAM a cada N instruções (maior que zero)")
	entrada := flag.String("entrada", "", "texto ASCII disponível na porta de entrada")
	flag.Parse()
	if *intervalo <= 0 {
		falhar("-video-a-cada deve ser maior que zero")
	}
	arquivo, erro := os.Open(*programa)
	if erro != nil {
		falhar("abrir programa: %v", erro)
	}
	defer arquivo.Close()
	palavras, erro := carregador.LerHex(arquivo)
	if erro != nil {
		falhar("carregar programa: %v", erro)
	}
	computador, erro := maquina.Nova(palavras)
	if erro != nil {
		falhar("inicializar máquina: %v", erro)
	}
	if erro := computador.FornecerEntrada(*entrada); erro != nil {
		falhar("%v", erro)
	}
	if erro := executarPrograma(os.Stdout, computador, *maximoPassos, *intervalo); erro != nil {
		falhar("%v", erro)
	}
}

func executarPrograma(saida io.Writer, computador *maquina.Maquina, maximoPassos, intervalo int) error {
	if intervalo <= 0 {
		return fmt.Errorf("intervalo de vídeo deve ser maior que zero")
	}
	defer exibirVideo(saida, computador)
	if maximoPassos <= 0 {
		return fmt.Errorf("--max-passos deve ser maior que zero")
	}
	for executados := 0; executados < maximoPassos; executados++ {
		if computador.Concluida() {
			fmt.Fprintf(saida, "PC=0x%08x\nFim da execução após %d passo(s)\n", computador.PC(), executados)
			return nil
		}
		resultado, erro := computador.Passo()
		if erro != nil {
			return fmt.Errorf("passo %d: %w", executados+1, erro)
		}
		fmt.Fprintf(saida, "Passo %d: PC=0x%08x palavra=0x%08x %s próximo PC=0x%08x\n", executados+1, resultado.PC, resultado.Palavra, resultado.Operacao, resultado.ProximoPC)
		exibirRegistradores(saida, computador.Registradores())
		if texto := computador.ConsumirSaida(); texto != "" {
			fmt.Fprintf(saida, "[E/S] %s\n", texto)
		}
		if (executados+1)%intervalo == 0 && !computador.Concluida() {
			exibirVideo(saida, computador)
		}
	}
	if computador.Concluida() {
		fmt.Fprintf(saida, "PC=0x%08x\nFim da execução após %d passo(s)\n", computador.PC(), maximoPassos)
		return nil
	}
	return fmt.Errorf("limite de %d passos atingido", maximoPassos)
}

func exibirRegistradores(saida io.Writer, registradores [32]uint32) {
	for indice, valor := range registradores {
		fmt.Fprintf(saida, "x%-2d=0x%08x", indice, valor)
		if indice%4 == 3 {
			fmt.Fprintln(saida)
		} else {
			fmt.Fprint(saida, "  ")
		}
	}
	fmt.Fprintln(saida)
}

func falhar(formato string, argumentos ...any) {
	fmt.Fprintf(os.Stderr, "Erro: "+formato+"\n", argumentos...)
	os.Exit(1)
}

func exibirVideo(saida io.Writer, computador *maquina.Maquina) {
	if texto := computador.TextoVRAM(); texto != "" {
		fmt.Fprintf(saida, "[VRAM]\n%s\n[/VRAM]\n", texto)
	}
}
