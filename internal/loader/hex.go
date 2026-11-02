// O pacote carregador lê programas com uma palavra hexadecimal de 32 bits por linha
package carregador

import (
	"bufio"   // Leitura para entrada e saida
	"fmt"     // Formatação e exibição de texto no terminal
	"io"      // Interfaces básicas de leitura e escrita
	"strconv" // Conversão de tipos de dados
	"strings" // Manipulação de textos/strings
)

// LerHex ignora linhas vazias e comentários iniciados por #
func LerHex(leitor io.Reader) ([]uint32, error) {
	var palavras []uint32
	varredor := bufio.NewScanner(leitor)
	for linha := 1; varredor.Scan(); linha++ {
		parte, _, _ := strings.Cut(varredor.Text(), "#")
		parte = strings.TrimSpace(parte)
		if parte == "" {
			continue
		}
		parte = strings.TrimPrefix(strings.ToLower(parte), "0x")
		if len(parte) != 8 {
			return nil, fmt.Errorf("linha %d: esperados 8 dígitos hexadecimais", linha)
		}
		palavra, erro := strconv.ParseUint(parte, 16, 32)
		if erro != nil {
			return nil, fmt.Errorf("linha %d: palavra inválida %q: %w", linha, parte, erro)
		}
		palavras = append(palavras, uint32(palavra))
	}
	if erro := varredor.Err(); erro != nil {
		return nil, fmt.Errorf("ler programa: %w", erro)
	}
	return palavras, nil
}
