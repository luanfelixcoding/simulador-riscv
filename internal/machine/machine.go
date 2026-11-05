// O pacote maquina implementa o estado da CPU e a execução de instruções
package maquina

import (
	"encoding/binary" // Leitura e escrita de dados binários
	"errors"          // Criação e manipulação de erros simples
	"fmt"             // Formatação e exibição de texto no terminal

	decodificador "simriscv/internal/isa"
)

const TamanhoRAM = 0x80000

var ErroProgramaConcluido = errors.New("fim do programa carregado")

type Maquina struct {
	registradores [32]uint32
	pc            uint32
	memoria       []byte
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
		return nil, fmt.Errorf("programa excede a RAM: %d palavras", len(palavras))
	}
	maquina := &Maquina{memoria: make([]byte, TamanhoRAM), bytesPrograma: uint32(len(palavras) * 4)}
	for indice, palavra := range palavras {
		binary.LittleEndian.PutUint32(maquina.memoria[indice*4:indice*4+4], palavra)
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
	palavra := binary.LittleEndian.Uint32(maquina.memoria[pc : pc+4])
	instrucao, erro := decodificador.Decodificar(palavra)
	if erro != nil {
		return ResultadoPasso{}, fmt.Errorf("PC=0x%08x: %w", pc, erro)
	}
	var valor uint32
	escreverRegistrador := true
	switch instrucao.Operacao {
	case decodificador.LW:
		endereco := maquina.registradores[instrucao.Fonte1] + uint32(instrucao.ImediatoI)
		valor, erro = maquina.LerPalavra(endereco)
		if erro != nil {
			return ResultadoPasso{}, fmt.Errorf("PC=0x%08x LW: %w", pc, erro)
		}
	case decodificador.SW:
		endereco := maquina.registradores[instrucao.Fonte1] + uint32(instrucao.ImediatoS)
		if erro := maquina.EscreverPalavra(endereco, maquina.registradores[instrucao.Fonte2]); erro != nil {
			return ResultadoPasso{}, fmt.Errorf("PC=0x%08x SW: %w", pc, erro)
		}
		escreverRegistrador = false
	case decodificador.ADDI:
		valor = maquina.registradores[instrucao.Fonte1] + uint32(instrucao.ImediatoI)
	case decodificador.ADD:
		valor = maquina.registradores[instrucao.Fonte1] + maquina.registradores[instrucao.Fonte2]
	default:
		return ResultadoPasso{}, fmt.Errorf("operação não implementada: %s", instrucao.Operacao)
	}
	if escreverRegistrador && instrucao.Destino != 0 {
		maquina.registradores[instrucao.Destino] = valor
	}
	maquina.pc += 4
	return ResultadoPasso{PC: pc, ProximoPC: maquina.pc, Palavra: palavra, Operacao: instrucao.Operacao}, nil
}
