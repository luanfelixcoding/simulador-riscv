# Simulador RISC-V em Golang

Projeto para matéria de `Organização de Computadores`.

## Organização e funcionamento

- `cmd/simriscv` contém a entrada do programa
- `internal/loader` lê o arquivo hexadecimal
- `internal/isa` extrai os campos e decodifica a instrução
- `internal/machine` guarda os registradores e busca a palavra na memória usando o PC

O arquivo aceita uma palavra de oito dígitos hexadecimais por linha, prefixo `0x` opcional, linhas vazias e comentários com `#`. As palavras são gravadas em little-endian na RAM de 512 KiB.

`Passo()` executa uma instrução. O PC avança quatro bytes, exceto nos saltos e desvios tomados. A execução termina quando o PC chega exatamente ao fim do programa carregado. Um PC fora desse trecho produz erro na próxima busca. `x0` permanece zero. Erros preservam PC, registradores e memória.

## Instruções disponíveis

O simulador implementa **28 instruções**: as 27 obrigatórias para todos os grupos e `MIN`, a especial do grupo 6. É um subconjunto do RV32I, com a operação MIN.

- Aritmética/lógica (10): `ADD SUB SLL SLT SLTU XOR SRL SRA OR AND`.
- Imediatos (9): `ADDI SLTI SLTIU XORI ORI ANDI SLLI SRLI SRAI`.
- Memória (2): `LW SW`.
- Desvios condicionais (6): `BEQ BNE BLT BGE BLTU BGEU`.
- Especial do grupo (1): `MIN`.

`LW` e `SW` acessam palavras de 32 bits e exigem endereços alinhados em quatro bytes. Desvios tomados também exigem destinos alinhados em quatro bytes. Overflow aritmético mantém os 32 bits inferiores. Instruções fora desse conjunto são rejeitadas.

### MIN — instrução especial do grupo 6

`MIN rd, rs1, rs2` grava em rd o menor dos dois valores, interpretados como inteiros de 32 bits **com sinal**. Por exemplo, entre 15 e -3, o resultado é -3. Escritas em x0 são descartadas.

Foi mantida a codificação de MIN da extensão Zbb: formato R, opcode `0x33`, funct3 `0x4` e funct7 `0x05`. Somente MIN dessa extensão é implementada.

Palavra de máquina: `(0x05 << 25) | (rs2 << 20) | (rs1 << 15) | (0x4 << 12) | (rd << 7) | 0x33`.

Exemplo: `MIN x6, x4, x5` corresponde a `0x0a524333`.

## Demonstração

```sh
go vet ./...
go run ./cmd/simriscv -programa programs/instrucoes-iniciais.hex
go run ./cmd/simriscv -programa programs/instrucoes-completas.hex
```

O exemplo inicial resulta em x1=5, x2=3 e x3=8, em três passos. O arquivo `instrucoes-completas.hex` demonstra um laço de soma (1+2+3+4+5), armazenamento e leitura de 15 na memória e MIN entre 15 e -3. Ao final, x3=x4=15 e x5=x6=0xfffffffd (-3), em 22 passos. É uma demonstração de parte das instruções, não de todas as 28.

Barramento, VRAM, periféricos, interrupções e cache permanecem para outras etapas.

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
