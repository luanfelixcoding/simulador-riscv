# Simulador RISC-V em Golang

Projeto para matéria de `Organização de Computadores`.

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
