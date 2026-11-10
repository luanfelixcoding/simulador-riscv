// O pacote isa decodifica as 27 instruções obrigatórias e MIN.
package isa

import "fmt"

type Operacao string

const (
	ADD   Operacao = "ADD"
	SUB   Operacao = "SUB"
	SLL   Operacao = "SLL"
	SLT   Operacao = "SLT"
	SLTU  Operacao = "SLTU"
	XOR   Operacao = "XOR"
	SRL   Operacao = "SRL"
	SRA   Operacao = "SRA"
	OR    Operacao = "OR"
	AND   Operacao = "AND"
	ADDI  Operacao = "ADDI"
	SLTI  Operacao = "SLTI"
	SLTIU Operacao = "SLTIU"
	XORI  Operacao = "XORI"
	ORI   Operacao = "ORI"
	ANDI  Operacao = "ANDI"
	SLLI  Operacao = "SLLI"
	SRLI  Operacao = "SRLI"
	SRAI  Operacao = "SRAI"
	LW    Operacao = "LW"
	SW    Operacao = "SW"
	BEQ   Operacao = "BEQ"
	BNE   Operacao = "BNE"
	BLT   Operacao = "BLT"
	BGE   Operacao = "BGE"
	BLTU  Operacao = "BLTU"
	BGEU  Operacao = "BGEU"
	MIN   Operacao = "MIN"
)

type Instrucao struct {
	Palavra                                                   uint32
	Operacao                                                  Operacao
	CodigoOperacao, Destino, Fonte1, Fonte2, Funcao3, Funcao7 uint32
	ImediatoI, ImediatoS, ImediatoB                           int32
}

func Decodificar(p uint32) (Instrucao, error) {
	i := Instrucao{Palavra: p, CodigoOperacao: p & 127, Destino: (p >> 7) & 31, Funcao3: (p >> 12) & 7, Fonte1: (p >> 15) & 31, Fonte2: (p >> 20) & 31, Funcao7: p >> 25, ImediatoI: int32(p) >> 20}
	s := (p>>25)<<5 | (p>>7)&31
	b := (p>>31)<<12 | ((p>>7)&1)<<11 | ((p>>25)&63)<<5 | ((p>>8)&15)<<1

	i.ImediatoS = int32(s<<20) >> 20
	i.ImediatoB = int32(b<<19) >> 19

	switch i.CodigoOperacao {
	case 0x33:
		if i.Funcao7 == 0 {
			i.Operacao = [8]Operacao{ADD, SLL, SLT, SLTU, XOR, SRL, OR, AND}[i.Funcao3]
		}
		if i.Funcao7 == 0x20 {
			i.Operacao = map[uint32]Operacao{0: SUB, 5: SRA}[i.Funcao3]
		}
	case 0x0b: // MIN customizada (grupo 6)
		if i.Funcao7 == 0 && i.Funcao3 == 0 {
			i.Operacao = MIN
		}
	case 0x13:
		i.Operacao = map[uint32]Operacao{0: ADDI, 2: SLTI, 3: SLTIU, 4: XORI, 6: ORI, 7: ANDI}[i.Funcao3]
		if i.Funcao3 == 1 {
			if i.Funcao7 == 0 {
				i.Operacao = SLLI
			}
		}
		if i.Funcao3 == 5 {
			if i.Funcao7 == 0 {
				i.Operacao = SRLI
			}
			if i.Funcao7 == 0x20 {
				i.Operacao = SRAI
			}
		}
	case 0x03:
		if i.Funcao3 == 2 {
			i.Operacao = LW
		}
	case 0x23:
		if i.Funcao3 == 2 {
			i.Operacao = SW
		}
	case 0x63:
		i.Operacao = map[uint32]Operacao{0: BEQ, 1: BNE, 4: BLT, 5: BGE, 6: BLTU, 7: BGEU}[i.Funcao3]
	}

	if i.Operacao == "" {
		return Instrucao{}, fmt.Errorf("instrução não suportada: palavra=0x%08x", p)
	}
	return i, nil
}
