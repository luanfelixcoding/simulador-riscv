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

## Formato e resultado do exemplo

O arquivo `.hex` contém uma palavra de 32 bits por linha, com oito dígitos hexadecimais. O prefixo `0x` é opcional, `#` inicia comentário. A primeira palavra entra no endereço zero da RAM, com o byte menos significativo primeiro. O PC inicia em zero e avança quatro bytes por instrução. Não há instrução de parada, a execução termina após a última palavra carregada.

[`programs/instrucoes-iniciais.hex`](programs/primeiros-passos.hex) executa `ADDI x1,x0,5`, `ADDI x2,x1,-2` e `ADD x3,x1,x2`. O resultado esperado é `x1=5`, `x2=3`, `x3=8` e PC final `0x0000000c`. A CLI mostra a palavra, o PC e os 32 registradores após cada passo. `x0` permanece zero. Codificações não suportadas causam erro sem avançar o PC.

## Feito por 👥
- [Luan](https://github.com/luanfelixcoding)
- [Felipe](https://github.com/Mafra013)