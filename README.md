# Simulador RISC-V em Golang

Projeto para matéria de `Organização de Computadores`.

## Organização e funcionamento

- `cmd/simriscv` contém a entrada do programa
- `internal/loader` lê o arquivo hexadecimal
- `internal/isa` extrai os campos e decodifica a instrução
- `internal/machine` guarda os registradores e busca a palavra na memória usando o PC

O arquivo aceita uma palavra de oito dígitos hexadecimais por linha, prefixo `0x` opcional, linhas vazias e comentários com `#`. As palavras são gravadas em little-endian na RAM de 512 KiB.

`Passo()` executa uma instrução. O PC avança quatro bytes, exceto nos desvios tomados. A execução termina quando o PC chega exatamente ao fim do programa carregado. Um PC fora desse trecho produz erro na próxima busca. `x0` permanece zero. Erros preservam PC, registradores e memória.

## Instruções disponíveis

O simulador implementa **28 instruções**: as 27 obrigatórias para todos os grupos e `MIN` do grupo 6. É um subconjunto do RV32I, com a operação MIN.

- Aritmética/lógica (10): `ADD SUB SLL SLT SLTU XOR SRL SRA OR AND`.
- Imediatos (9): `ADDI SLTI SLTIU XORI ORI ANDI SLLI SRLI SRAI`.
- Memória (2): `LW SW`.
- Desvios condicionais (6): `BEQ BNE BLT BGE BLTU BGEU`.
- Especial do grupo (1): `MIN`.

`LW` e `SW` acessam palavras de 32 bits e exigem endereços alinhados em quatro bytes. Desvios tomados também exigem destinos alinhados em quatro bytes. Overflow aritmético mantém os 32 bits inferiores. Instruções fora desse conjunto são rejeitadas.

### MIN - instrução especial do grupo 6

`MIN rd, rs1, rs2` grava em rd o menor dos dois valores, interpretados como inteiros de 32 bits **com sinal**. Por exemplo, entre 15 e -3, o resultado é -3. Escritas em x0 são descartadas.

Palavra de máquina: `(rs2 << 20) | (rs1 << 15) | (rd << 7) | 0x0B`.

Exemplo: `MIN x6, x4, x5` corresponde a `0x0052030b`.

## Demonstração

```sh
go run ./cmd/simriscv -programa programs/instrucoes-iniciais.hex
go run ./cmd/simriscv -programa programs/instrucoes-completas.hex
```

O exemplo inicial resulta em x1=5, x2=3 e x3=8, em três passos. O arquivo `instrucoes-completas.hex` demonstra um laço de soma (1+2+3+4+5), armazenamento e leitura de 15 na memória e MIN entre 15 e -3. Ao final, `x3=x4=15` e `x5=x6=0xfffffffd (-3)`, em 22 passos.

Falta Barramento, VRAM, periféricos, interrupções e cache.

## Programa de teste de todas as instruções

O arquivo `programs/teste-todas-instrucoes.hex` executa **as 27 instruções obrigatórias e MIN custom-0**. Cada linha contém a palavra de máquina, o PC, a instrução e seu resultado esperado. Os seis desvios são exercitados tanto tomados quanto não tomados, há também um laço com deslocamento negativo e uma tentativa de escrita em x0.

Na raiz do repositório, no Terminal macOS ou no PowerShell Windows:

```sh
go run ./cmd/simriscv -programa programs/teste-todas-instrucoes.hex
```

Resultado esperado: **49 passos**, **PC=0x000000d4**, `x30=6`, `x31=0`, `x24=0` e `x0=0`. MIN coloca -8 (`0xfffffff8`) em x22, SW grava esse valor em RAM[0x400] e LW o recupera em x23. RAM[0x404] recebe 6. Os dados ficam fora dos 212 bytes de código.

`x30` conta os seis pares de desvios verificados e `x31` deve permanecer zero, esses contadores não validam sozinhos todas as operações. O teste Go confere os 32 registradores finais, as duas palavras de dados, o término e a execução efetiva das 28 operações, além dos dois caminhos de cada desvio. Os testes unitários continuam cobrindo outros casos-limite. Este exemplo não inclui os recursos ainda pendentes da Etapa 3.

## Executar

É necessário Go 1.22 ou posterior. Os comandos abaixo devem ser executados **na raiz deste repositório**.

macOS, no Terminal:

```sh
go version
go run ./cmd/simriscv -programa programs/instrucoes-iniciais.hex
```

Windows, no PowerShell:

```powershell
go version
go run ./cmd/simriscv -programa programs/instrucoes-iniciais.hex
```

Também é possível usar `go run ./cmd/simriscv`, pois o arquivo acima é o padrão. `-max-passos N` limita o número de instruções, o valor padrão é 100. `go run ./cmd/simriscv -help` mostra as opções.

Para compilar, use `go build -o simriscv ./cmd/simriscv` no macOS ou `go build -o simriscv.exe ./cmd/simriscv` no Windows.

## Feito por 👥
- [Luan](https://github.com/luanfelixcoding)
- [Felipe](https://github.com/Mafra013)
