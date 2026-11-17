package maquina

import "fmt"

// Convenções locais: tratador nos últimos 4 KiB da RAM e retorno via MMIO.
const (
	VetorInterrupcao uint32 = 0x7f000
	PortaRetorno     uint32 = 0x9fc0c
)

type controladorInterrupcao struct {
	periodo, contagem, fim uint32
	ativo, retorno         bool
	pcSalvo                uint32
	registradoresSalvos    [32]uint32
	atendidas              uint64
}

// ConfigurarTemporizador deve ser chamada antes da primeira instrução concluída.
// O temporizador não conta instruções do tratador e não permite aninhamento.
func (m *Maquina) ConfigurarTemporizador(periodo uint32, palavras []uint32) error {
	if m.iniciada {
		return fmt.Errorf("temporizador só pode ser configurado antes da execução")
	}
	if periodo == 0 || len(palavras) == 0 || len(palavras) > int((TamanhoRAM-VetorInterrupcao)/4) {
		return fmt.Errorf("período deve ser positivo e tratador deve ocupar entre 1 e 1024 palavras")
	}
	if m.bytesPrograma > VetorInterrupcao {
		return fmt.Errorf("programa sobrepõe região do tratador")
	}
	// Validação completa precede qualquer escrita. Todos os endereços são RAM alinhada.
	for indice, palavra := range palavras {
		if erro := m.EscreverPalavra(VetorInterrupcao+uint32(indice)*4, palavra); erro != nil {
			return erro
		}
	}
	m.barramento.interrupcao = controladorInterrupcao{periodo: periodo, fim: VetorInterrupcao + uint32(len(palavras))*4}
	return nil
}

func (m *Maquina) EmInterrupcao() bool           { return m.barramento.interrupcao.ativo }
func (m *Maquina) InterrupcoesAtendidas() uint64 { return m.barramento.interrupcao.atendidas }

func (m *Maquina) atualizarInterrupcao() {
	c := &m.barramento.interrupcao
	if c.ativo {
		if c.retorno {
			m.pc, m.registradores = c.pcSalvo, c.registradoresSalvos
			c.ativo, c.retorno = false, false
		}
		return
	}
	if c.periodo == 0 || m.Concluida() {
		return
	}
	c.contagem++
	if c.contagem == c.periodo {
		c.contagem = 0
		c.pcSalvo, c.registradoresSalvos = m.pc, m.registradores
		c.ativo = true
		c.atendidas++
		m.pc = VetorInterrupcao
	}
}
