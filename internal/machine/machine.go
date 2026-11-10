// O pacote maquina implementa o estado da CPU e a execução de instruções
package maquina

import (
	"errors" // Criação e manipulação de erros simples
	"fmt"

	decodificador "simriscv/internal/isa"
)

var ErroProgramaConcluido = errors.New("fim do programa carregado")

type Maquina struct {
	registradores [32]uint32
	pc            uint32
	barramento    *Barramento
	bytesPrograma uint32
}

type ResultadoPasso struct {
	PC        uint32
	ProximoPC uint32
	Palavra   uint32
	Operacao  decodificador.Operacao
}

// As palavras de máquina são colocadas a partir do endereço zero da RAM
// Os bytes de cada palavra são armazenados do menos ao mais significativo
func Nova(palavras []uint32) (*Maquina, error) {
	if len(palavras) > TamanhoRAM/4 {
		return nil, fmt.Errorf(
			"programa excede a RAM: %d palavras",
			len(palavras),
		)
	}

	ram := NovaRAM()

	maquina := &Maquina{
		barramento:    NovoBarramento(ram),
		bytesPrograma: uint32(len(palavras) * 4),
	}

	for indice, palavra := range palavras {
		endereco := uint32(indice * 4)

		if erro := maquina.EscreverPalavra(endereco, palavra); erro != nil {
			return nil, fmt.Errorf("carregar programa: %w", erro)
		}
	}

	return maquina, nil
}

// Retorna o contador de programa (Program Counter) atual da máquina
func (maquina *Maquina) PC() uint32 { return maquina.pc }

// Verifica se a execução da máquina foi finalizada (quando o PC atinge o fim do programa carregado)
func (maquina *Maquina) Concluida() bool { return maquina.pc == maquina.bytesPrograma }

// Registradores devolve uma cópia para impedir alterações diretas no estado.
func (maquina *Maquina) Registradores() [32]uint32 { return maquina.registradores }

// Passo busca, decodifica e executa
// Um erro deixa o PC e os registradores inalterados
func (maquina *Maquina) Passo() (ResultadoPasso, error) {
	pc := maquina.pc
	if maquina.Concluida() {
		return ResultadoPasso{}, ErroProgramaConcluido
	}
	if pc%4 != 0 || pc > maquina.bytesPrograma || maquina.bytesPrograma-pc < 4 {
		return ResultadoPasso{}, fmt.Errorf("PC fora do programa ou desalinhado: 0x%08x", pc)
	}
	palavra, erro := maquina.LerPalavra(pc)
	if erro != nil {
		return ResultadoPasso{}, fmt.Errorf(
			"busca no PC=0x%08x: %w",
			pc,
			erro,
		)
	}

	instrucao, erro := decodificador.Decodificar(palavra)
	if erro != nil {
		return ResultadoPasso{}, fmt.Errorf("PC=0x%08x: %w", pc, erro)
	}

	// Calcula em uma cópia: erros não deixam escritas parciais nos registradores.
	regs := maquina.registradores
	a, b := regs[instrucao.Fonte1], regs[instrucao.Fonte2]
	imm := uint32(instrucao.ImediatoI)
	proximoPC := pc + 4
	var valor uint32
	escreverRegistrador := true
	switch instrucao.Operacao {
	case decodificador.ADD:
		valor = a + b
	case decodificador.SUB:
		valor = a - b
	case decodificador.ADDI:
		valor = a + imm
	case decodificador.SLL:
		valor = a << (b & 31)
	case decodificador.SRL:
		valor = a >> (b & 31)
	case decodificador.SRA:
		valor = uint32(int32(a) >> (b & 31))
	case decodificador.SLLI:
		valor = a << (imm & 31)
	case decodificador.SRLI:
		valor = a >> (imm & 31)
	case decodificador.SRAI:
		valor = uint32(int32(a) >> (imm & 31))
	case decodificador.SLT:
		valor = bit(int32(a) < int32(b))
	case decodificador.SLTU:
		valor = bit(a < b)
	case decodificador.SLTI:
		valor = bit(int32(a) < instrucao.ImediatoI)
	case decodificador.SLTIU:
		valor = bit(a < imm)
	case decodificador.XOR:
		valor = a ^ b
	case decodificador.OR:
		valor = a | b
	case decodificador.AND:
		valor = a & b
	case decodificador.XORI:
		valor = a ^ imm
	case decodificador.ORI:
		valor = a | imm
	case decodificador.ANDI:
		valor = a & imm
	case decodificador.BEQ, decodificador.BNE, decodificador.BLT, decodificador.BGE, decodificador.BLTU, decodificador.BGEU:
		escreverRegistrador = false
		tomar := false
		switch instrucao.Operacao {
		case decodificador.BEQ:
			tomar = a == b
		case decodificador.BNE:
			tomar = a != b
		case decodificador.BLT:
			tomar = int32(a) < int32(b)
		case decodificador.BGE:
			tomar = int32(a) >= int32(b)
		case decodificador.BLTU:
			tomar = a < b
		case decodificador.BGEU:
			tomar = a >= b
		}
		if tomar {
			proximoPC = pc + uint32(instrucao.ImediatoB)
		}
	case decodificador.LW:
		valor, erro = maquina.LerPalavra(a + imm)
		if erro != nil {
			return ResultadoPasso{}, fmt.Errorf("PC=0x%08x LW: %w", pc, erro)
		}
	case decodificador.SW:
		escreverRegistrador = false
		if erro = maquina.EscreverPalavra(a+uint32(instrucao.ImediatoS), b); erro != nil {
			return ResultadoPasso{}, fmt.Errorf("PC=0x%08x SW: %w", pc, erro)
		}
	case decodificador.MIN:
		valor = a
		if int32(b) < int32(a) {
			valor = b
		}
	default:
		return ResultadoPasso{}, fmt.Errorf("operação não implementada: %s", instrucao.Operacao)
	}
	if proximoPC%4 != 0 {
		return ResultadoPasso{}, fmt.Errorf("PC de destino desalinhado: 0x%08x", proximoPC)
	}
	if escreverRegistrador && instrucao.Destino != 0 {
		regs[instrucao.Destino] = valor
	}
	regs[0] = 0
	maquina.registradores = regs
	maquina.pc = proximoPC
	return ResultadoPasso{PC: pc, ProximoPC: proximoPC, Palavra: palavra, Operacao: instrucao.Operacao}, nil
}

func bit(condicao bool) uint32 {
	if condicao {
		return 1
	}
	return 0
}
