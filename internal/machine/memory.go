package maquina

import (
	"encoding/binary"
	"fmt"
)

func validarEndereco(endereco uint32) error {
	if endereco%4 != 0 {
		return fmt.Errorf("endereço desalinhado: 0x%08x", endereco)
	}
	if endereco > TamanhoRAM-4 {
		return fmt.Errorf("endereço fora da RAM: 0x%08x", endereco)
	}
	return nil
}

// LerPalavra lê quatro bytes da RAM em ordem little-endian
func (maquina *Maquina) LerPalavra(endereco uint32) (uint32, error) {
	if erro := validarEndereco(endereco); erro != nil {
		return 0, erro
	}
	return binary.LittleEndian.Uint32(maquina.memoria[endereco : endereco+4]), nil
}

// EscreverPalavra valida antes de gravar, preservando a RAM em caso de erro
// Código e dados compartilham a RAM, onde as escritas podem modificar instruções.
func (maquina *Maquina) EscreverPalavra(endereco, valor uint32) error {
	if erro := validarEndereco(endereco); erro != nil {
		return erro
	}
	binary.LittleEndian.PutUint32(maquina.memoria[endereco:endereco+4], valor)
	return nil
}
