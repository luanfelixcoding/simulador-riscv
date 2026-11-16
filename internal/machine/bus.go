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
	Endereco    uint32
	Dados       uint32
	Controle    Controle
	ram         *RAM
	vram        *VRAM
	perifericos *Perifericos
}

func NovoBarramento(ram *RAM) *Barramento {
	return &Barramento{ram: ram, vram: &VRAM{}, perifericos: &Perifericos{}}
}

func (b *Barramento) transferir(
	endereco, dados uint32,
	controle Controle,
) (uint32, error) {
	if endereco%4 != 0 {
		return 0, fmt.Errorf("endereço desalinhado: 0x%08x", endereco)
	}
	var dispositivo interface {
		LerPalavra(uint32) (uint32, error)
		EscreverPalavra(uint32, uint32) error
	}
	local := endereco
	switch {
	case endereco < TamanhoRAM:
		if b.ram == nil {
			return 0, fmt.Errorf("RAM não conectada")
		}
		dispositivo = b.ram
	case endereco >= InicioVRAM && endereco < InicioVRAM+TamanhoVRAM:
		dispositivo = b.vram
		local -= InicioVRAM
	case endereco >= PortaSaida && endereco <= 0x9ffff:
		dispositivo = b.perifericos
	default:
		return 0, fmt.Errorf("endereço reservado ou sem dispositivo: 0x%08x", endereco)
	}

	switch controle {
	case Leitura:
		valor, erro := dispositivo.LerPalavra(local)
		if erro != nil {
			return 0, erro
		}
		dados = valor

	case Escrita:
		if erro := dispositivo.EscreverPalavra(local, dados); erro != nil {
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
