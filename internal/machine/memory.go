package maquina

import (
	"encoding/binary"
	"fmt"
)

const TamanhoRAM = 0x80000

// RAM armazena os bytes; o barramento encaminha os acessos.
type RAM struct {
	dados []byte
}

func NovaRAM() *RAM {
	return &RAM{dados: make([]byte, TamanhoRAM)}
}

func validarEndereco(endereco uint32) error {
	if endereco%4 != 0 {
		return fmt.Errorf("endereço desalinhado: 0x%08x", endereco)
	}
	if endereco > TamanhoRAM-4 {
		return fmt.Errorf("endereço fora da RAM: 0x%08x", endereco)
	}
	return nil
}

func (ram *RAM) LerPalavra(endereco uint32) (uint32, error) {
	if erro := validarEndereco(endereco); erro != nil {
		return 0, erro
	}

	return binary.LittleEndian.Uint32(
		ram.dados[endereco : endereco+4],
	), nil
}

func (ram *RAM) EscreverPalavra(endereco, valor uint32) error {
	if erro := validarEndereco(endereco); erro != nil {
		return erro
	}

	binary.LittleEndian.PutUint32(
		ram.dados[endereco:endereco+4],
		valor,
	)
	return nil
}

// Mantêm a interface da máquina; agora todo acesso passa pelo barramento.
func (maquina *Maquina) LerPalavra(endereco uint32) (uint32, error) {
	return maquina.barramento.LerPalavra(endereco)
}

func (maquina *Maquina) EscreverPalavra(endereco, valor uint32) error {
	return maquina.barramento.EscreverPalavra(endereco, valor)
}
