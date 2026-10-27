package main

import "fmt"

// Três instruções RISC-V de 32 bits
func main() {
	programa := []uint32{
		0x00500093, // ADDI x1, x0, 5
		0xffe08113, // ADDI x2, x1, -2
		0x002081b3, // ADD  x3, x1, x2
	}

	var registradores [32]uint32
	var pc uint32

	for _, palavra := range programa {
		opcode := palavra & 0x7f
		destino := (palavra >> 7) & 0x1f
		funcao3 := (palavra >> 12) & 0x7
		fonte1 := (palavra >> 15) & 0x1f

		var resultado uint32
		var nome string

		switch {
		case opcode == 0x13 && funcao3 == 0: // ADDI
			imediato := int32(palavra) >> 20
			resultado = registradores[fonte1] + uint32(imediato)
			nome = "ADDI"
		case opcode == 0x33 && funcao3 == 0 && palavra>>25 == 0: // ADD
			fonte2 := (palavra >> 20) & 0x1f
			resultado = registradores[fonte1] + registradores[fonte2]
			nome = "ADD"
		default:
			fmt.Printf("Instrução não suportada: 0x%08x\n", palavra)
			return
		}

		if destino != 0 { // x0 é sempre zero
			registradores[destino] = resultado
		}
		fmt.Printf("PC=%d  %s  x%d=%d\n", pc, nome, destino, registradores[destino])
		pc += 4
	}
}
