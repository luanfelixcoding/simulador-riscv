package maquina

import (
	"encoding/binary"
	"fmt"
	"strings"
)

const (
	InicioVRAM   uint32 = 0x80000
	TamanhoVRAM  uint32 = 0x10000
	PortaSaida   uint32 = 0x9fc00
	PortaEntrada uint32 = 0x9fc04
	PortaEstado  uint32 = 0x9fc08
)

// VRAM: um byte por caractere, palavras little-endian via LW/SW.
type VRAM struct{ dados [TamanhoVRAM]byte }

func (v *VRAM) LerPalavra(e uint32) (uint32, error) {
	if e%4 != 0 || e > TamanhoVRAM-4 {
		return 0, fmt.Errorf("acesso inválido à VRAM: 0x%x", e)
	}
	return binary.LittleEndian.Uint32(v.dados[e : e+4]), nil
}
func (v *VRAM) EscreverPalavra(e, d uint32) error {
	if e%4 != 0 || e > TamanhoVRAM-4 {
		return fmt.Errorf("acesso inválido à VRAM: 0x%x", e)
	}
	binary.LittleEndian.PutUint32(v.dados[e:e+4], d)
	return nil
}

func textoASCII(dados []byte) string {
	var s strings.Builder
	for _, c := range dados {
		switch {
		case c == 10 || c == 9 || c >= 32 && c <= 126:
			s.WriteByte(c)
		default:
			s.WriteByte(' ')
		}
	}
	return s.String()
}
func (m *Maquina) TextoVRAM() string {
	d := m.barramento.vram.dados[:]
	n := len(d)
	for n > 0 && d[n-1] == 0 {
		n--
	}
	return textoASCII(d[:n])
}

type Perifericos struct{ entrada, saida []byte }

func (p *Perifericos) LerPalavra(e uint32) (uint32, error) {
	switch e {
	case PortaEntrada:
		if len(p.entrada) == 0 {
			return 0, nil
		}
		c := p.entrada[0]
		p.entrada = p.entrada[1:]
		return uint32(c), nil
	case PortaEstado:
		estado := uint32(2)
		if len(p.entrada) > 0 {
			estado |= 1
		}
		return estado, nil
	default:
		return 0, fmt.Errorf("porta sem permissão de leitura: 0x%x", e)
	}
}
func (p *Perifericos) EscreverPalavra(e, d uint32) error {
	if e != PortaSaida {
		return fmt.Errorf("porta sem permissão de escrita: 0x%x", e)
	}
	p.saida = append(p.saida, byte(d))
	return nil
}

// A entrada antes da execução
func (m *Maquina) FornecerEntrada(texto string) error {
	for _, c := range []byte(texto) {
		if c > 127 {
			return fmt.Errorf("entrada deve conter apenas caracteres ASCII")
		}
	}
	m.barramento.perifericos.entrada = append(m.barramento.perifericos.entrada, []byte(texto)...)
	return nil
}
func (m *Maquina) ConsumirSaida() string {
	p := m.barramento.perifericos
	s := textoASCII(p.saida)
	p.saida = nil
	return s
}
