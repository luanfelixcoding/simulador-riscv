// O pacote isa decodifica o subconjunto de instruções
package isa

import "fmt" // Formatação e exibição de texto no terminal

type Operacao string

const (
	ADDI Operacao = "ADDI"
	ADD  Operacao = "ADD"
)

// Struct para mapeamento dos campos de uma instrução.
type Instrucao struct {
	Palavra        uint32   // Instrução completa
	Operacao       Operacao // Nome Operação
	CodigoOperacao uint32   // campo Opcode, identificar classe de instrução
	Destino        uint32   // Registrador que recebe o resultado
	Fonte1         uint32   // Registrador de entrada rs1
	Fonte2         uint32   // Registrador de entrada rs2
	Funcao3        uint32   // Campo de 3 bits
	Funcao7        uint32   // Campo de 7 bits
	ImediatoI      int32    // Valor constante com sinal
}

// Decodificador ADDI e ADD
func Decodificar(palavra uint32) (Instrucao, error) {
	instrucao := Instrucao{
		Palavra:        palavra,
		CodigoOperacao: palavra & 0x7f,
		Destino:        (palavra >> 7) & 0x1f,
		Funcao3:        (palavra >> 12) & 0x7,
		Fonte1:         (palavra >> 15) & 0x1f,
		Fonte2:         (palavra >> 20) & 0x1f,
		Funcao7:        (palavra >> 25) & 0x7f,
		ImediatoI:      int32(palavra) >> 20,
	}
	switch {
	case instrucao.CodigoOperacao == 0x13 && instrucao.Funcao3 == 0:
		instrucao.Operacao = ADDI
	case instrucao.CodigoOperacao == 0x33 && instrucao.Funcao3 == 0 && instrucao.Funcao7 == 0:
		instrucao.Operacao = ADD
	default:
		return Instrucao{}, fmt.Errorf("Instrução não suportada: palavra=0x%08x", palavra)
	}
	return instrucao, nil
}
