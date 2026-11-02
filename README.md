# Simulador RISC-V em Golang

Projeto para matéria de `Organização de Computadores`.

## Organização e funcionamento

- `cmd/simriscv` contém a entrada do programa
- `internal/loader` lê o arquivo hexadecimal
- `internal/isa` extrai os campos e decodifica a instrução
- `internal/machine` guarda os registradores e busca a palavra na memória usando o PC
- `inicio/main.go` centralizador da execução do programa

O arquivo aceita uma palavra de oito dígitos hexadecimais por linha, prefixo `0x` opcional, linhas vazias e comentários com `#`. As palavras são gravadas em bytes, do menos ao mais significativo. Nesta etapa apenas o trecho do programa é alocado, limitado a 512 KiB.

`Passo()` executa uma instrução e avança o PC em quatro bytes. O programa termina ao atingir o fim das palavras carregadas, um limite configurável impede execução ilimitada. `x0` permanece zero. Codificações inválidas deixam PC e registradores inalterados e geram diagnóstico. O exemplo termina em três passos com `x1=5`, `x2=3`, `x3=8` e PC `0x0c` (12 decimal).

## Demonstração inicial

Para apresentar o começo do projeto, execute na raiz desta pasta:

```sh
go run ./inicio
```

O arquivo [`inicio/main.go`](inicio/main.go) é um exemplo pequeno e independente: executa duas instruções `ADDI` e uma `ADD`, mostra o PC e os resultados nos registradores. O resultado esperado é `x1=5`, `x2=3` e `x3=8`.

## Executar

É necessário Go 1.22 ou posterior. Os comandos abaixo devem ser executados **na raiz deste repositório**.

macOS, no Terminal:

```sh
go version
go vet ./...
go run ./cmd/simriscv -programa programs/instrucoes-iniciais.hex
```

Windows, no PowerShell:

```powershell
go version
go vet ./...
go run ./cmd/simriscv -programa programs/instrucoes-iniciais.hex
```

Também é possível usar `go run ./cmd/simriscv`, pois o arquivo acima é o padrão. `-max-passos N` limita o número de instruções, o valor padrão é 100. `go run ./cmd/simriscv -help` mostra as opções.

Para compilar, use `go build -o simriscv ./cmd/simriscv` no macOS ou `go build -o simriscv.exe ./cmd/simriscv` no Windows.

## Feito por 👥
- [Luan](https://github.com/luanfelixcoding)
- [Felipe](https://github.com/Mafra013)