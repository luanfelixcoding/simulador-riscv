package maquina

import "fmt"

// Controle informa qual operação deve ser realizada.
type Controle uint8

const (
	Leitura Controle = iota + 1
	Escrita
)

// Barramento conecta a CPU aos dispositivos de memória.
// Os campos públicos registram a última transferência concluída.
type Barramento struct {
	Endereco uint32
	Dados    uint32
	Controle Controle
	ram      *RAM
}

func NovoBarramento(ram *RAM) *Barramento {
	return &Barramento{ram: ram}
}

func (b *Barramento) transferir(
	endereco, dados uint32,
	controle Controle,
) (uint32, error) {
	// Nesta etapa, somente a faixa da RAM está conectada.
	if endereco >= TamanhoRAM {
		return 0, fmt.Errorf(
			"endereço sem dispositivo conectado: 0x%08x",
			endereco,
		)
	}

	if b.ram == nil {
		return 0, fmt.Errorf("RAM não conectada ao barramento")
	}

	switch controle {
	case Leitura:
		valor, erro := b.ram.LerPalavra(endereco)
		if erro != nil {
			return 0, erro
		}
		dados = valor

	case Escrita:
		if erro := b.ram.EscreverPalavra(endereco, dados); erro != nil {
			return 0, erro
		}

	default:
		return 0, fmt.Errorf("controle inválido: %d", controle)
	}

	b.Endereco = endereco
	b.Dados = dados
	b.Controle = controle

	return dados, nil
}

func (b *Barramento) LerPalavra(endereco uint32) (uint32, error) {
	return b.transferir(endereco, 0, Leitura)
}

func (b *Barramento) EscreverPalavra(endereco, valor uint32) error {
	_, erro := b.transferir(endereco, valor, Escrita)
	return erro
}
