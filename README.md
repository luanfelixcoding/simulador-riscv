# Simulador RISC-V em Golang

Projeto para matéria de `Organização de Computadores`.

## Organização e funcionamento

- `cmd/simriscv` contém a entrada do programa
- `internal/loader` lê o arquivo hexadecimal
- `internal/isa` extrai os campos e decodifica a instrução
- `internal/machine` implementa CPU, RAM, barramento, VRAM e periféricos de E/S programada

O arquivo aceita uma palavra de oito dígitos hexadecimais por linha, prefixo `0x` opcional, linhas vazias e comentários com `#`. As palavras são gravadas em little-endian na RAM de 512 KiB.

`Passo()` executa uma instrução. O PC avança quatro bytes, exceto nos desvios tomados. A execução termina quando o PC chega exatamente ao fim do programa carregado. Um PC fora desse trecho produz erro na próxima busca. `x0` permanece zero. Erros de instrução ou acesso são reportados sem avançar o PC. Leituras da porta de entrada consomem um caractere, o barramento registra a última transferência concluída.

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

Barramento, VRAM e E/S programada estão implementados. A interrupção básica por temporizador também está implementada. Teclado interativo, cache e pipeline continuam como extras não implementados.

## Programa de teste de todas as instruções

O arquivo `programs/teste-todas-instrucoes.hex` executa **as 27 instruções obrigatórias e MIN custom-0**. Cada linha contém a palavra de máquina, o PC, a instrução e seu resultado esperado. Os seis desvios são exercitados tanto tomados quanto não tomados, há também um laço com deslocamento negativo e uma tentativa de escrita em x0.

Na raiz do repositório, no Terminal macOS ou no PowerShell Windows:

```sh
go run ./cmd/simriscv -programa programs/teste-todas-instrucoes.hex
```

Resultado esperado: **49 passos**, **PC=0x000000d4**, `x30=6`, `x31=0`, `x24=0` e `x0=0`. MIN coloca -8 (`0xfffffff8`) em x22, SW grava esse valor em RAM[0x400] e LW o recupera em x23. RAM[0x404] recebe 6. Os dados ficam fora dos 212 bytes de código.

`x30` conta os seis pares de desvios verificados e `x31` deve permanecer zero. Esses contadores não validam sozinhos todas as operações: compare também os demais resultados com os comentários do arquivo. Há testes automatizados Go de decodificação/execução das 28 operações, barramento, memória, E/S, CLI e interrupções, o teste do arquivo também confere os 32 registradores finais e a cobertura dinâmica das instruções.

## Barramento e mapa de memória

A CPU acessa os dispositivos pelo barramento, que usa endereços e dados de 32 bits e controle de leitura/escrita. O carregamento do programa, a busca de instruções, LW e SW utilizam esse caminho. O barramento registra a última transferência concluída.

| Faixa | Uso |
|---|---|
| `0x00000 - 0x7ffff` | RAM principal: 512 KiB |
| `0x80000 - 0x8ffff` | VRAM: 64 KiB |
| `0x90000 - 0x9fbff` | Reservada para expansão futura e acessos rejeitados |
| `0x9fc00 - 0x9ffff` | Região de periféricos, somente as portas abaixo estão implementadas |

Os acessos LW/SW exigem alinhamento de quatro bytes. Endereços reservados, portas não implementadas e operações sem permissão geram erro. Para acessar o vetor da VRAM, o barramento subtrai `0x80000` do endereço global.

## VRAM e exibição ASCII

SW escreve quatro bytes em little-endian na VRAM e LW pode lê-los de volta. Cada byte representa um caractere. Por exemplo, `0x00004b4f` grava os bytes `4f 4b 00 00`, formando "OK".

O terminal apresenta a VRAM a cada `-video-a-cada N` instruções (padrão: 10) e ao sair da execução, inclusive por erro ou limite de passos. Conteúdo vazio não é impresso. Zeros finais são omitidos da exibição, zeros internos e controles não exibíveis viram espaços. São preservados ASCII imprimível, tabulação e quebra de linha. Não há tela gráfica.

## E/S programada

Os endereços específicos e bits de estado são convenções do simulador dentro da região de periféricos do enunciado, não instruções adicionais da ISA.

| Porta | Permissão | Comportamento |
|---|---|---|
| `0x9fc00` | Escrita com SW | Envia o byte inferior à fila de saída, os outros 24 bits são descartados |
| `0x9fc04` | Leitura com LW | Retira um byte da fila de entrada e o devolve em 32 bits, retorna zero se vazia |
| `0x9fc08` | Leitura com LW | Estado: bit 0 indica entrada disponível, bit 1 indica saída pronta |
| `0x9fc0c` | Escrita com SW durante interrupção | Retorna ao programa principal, restaurando PC e registradores |

O estado é 3 quando há entrada e 2 quando a fila está vazia. A saída está sempre pronta neste modelo. Consulte o estado para distinguir ausência de entrada de um caractere de valor zero. A saída é consumida e exibida com o prefixo `[E/S]` após cada instrução.

A opção `-entrada` fornece texto ASCII antes da execução. Não há captura de teclado durante a execução. A interrupção básica é gerada pelo temporizador descrito abaixo. Letras acentuadas são rejeitadas nessa opção.

### Exemplo de VRAM e E/S

```sh
go run ./cmd/simriscv -programa programs/vram-es.hex -entrada A -video-a-cada 10
```

O exemplo termina em 13 passos, com PC `0x34`. Mostra "OK" na VRAM e `[E/S] A`. Resultados: `x1=0x80000`, `x2=x3=0x4b4f`, `x4=0x9fc00`, `x5=3`, `x6=0x41` e `x7=2`.

O programa lê o estado, mas não aguarda entrada em um laço: tenta ler um caractere uma vez. Se nenhuma entrada for fornecida, retorna zero e a saída é apresentada como espaço. Para exercitar as 28 instruções, use também `teste-todas-instrucoes.hex`.

## Demonstração integrada: mínimo de uma lista

`programs/minimo-lista.hex` inicializa [7,3,9,5] em RAM[0x100–0x10f], percorre os elementos com LW e BNE e calcula o menor com MIN custom-0. Grava o resultado 3 em RAM[0x110] e os bytes ASCII `MIN=3` seguidos de quebra de linha na VRAM por SW. O texto exibido vem da memória, não de uma mensagem fixa do simulador.

```sh
go run ./cmd/simriscv -programa programs/minimo-lista.hex
```

Resultado: **43 passos**, PC final **0x84**, x3=3 e VRAM `MIN=3`. A conversão ASCII desse exemplo atende ao mínimo de um dígito da lista fixa, não é um formatador genérico de inteiros negativos ou de vários dígitos. A instrução MIN continua assinada de 32 bits.

## Interrupção básica por temporizador

O PDF pede aceitação de interrupções, mas não fixa o protocolo. A convenção técnica deste simulador usa:

- `-interromper-a-cada N` e `-tratador ARQUIVO.hex`, fornecidos juntos. Sem essas opções, a interrupção permanece desativada.
- Rotina na RAM a partir de **0x7f000**, no máximo 1024 palavras (4 KiB). O programa principal não pode sobrepor essa área quando o temporizador está habilitado. Configuração permitida apenas antes da primeira instrução concluída.
- Após N instruções principais concluídas, o controlador ligado ao barramento salva o PC da próxima instrução e os 32 registradores em estado interno e desvia para o tratador. Não há interrupção depois do término do programa.
- O tratador usa as mesmas instruções e o mesmo barramento. O temporizador fica suspenso durante seu atendimento, não há aninhamento.
- **SW em 0x9fc0c**, com qualquer valor, solicita retorno: ao concluir essa instrução, PC e registradores são restaurados. Leituras nessa porta e escritas fora do atendimento geram erro. As alterações na RAM/VRAM e os efeitos de E/S do tratador permanecem.
- Sair do trecho carregado do tratador sem retornar gera erro na próxima busca. Um tratador em laço é limitado por `-max-passos`, que conta instruções principais e do tratador. A cadência de vídeo também conta ambos.

Este é um protocolo local via MMIO, não implementa CSRs, MRET ou o modo privilegiado RISC-V. Os detalhes podem ser ajustados se o professor trouxer outro contrato.

O exemplo `programs/tratador-temporizador.hex` incrementa RAM[0x200] e retorna, usando seis instruções:

```sh
go run ./cmd/simriscv -programa programs/minimo-lista.hex -interromper-a-cada 8 -tratador programs/tratador-temporizador.hex
```

Resultado: **73 passos** (43 principais + 30 do tratador), **cinco interrupções**, RAM[0x200]=5, PC final 0x84 e a mesma saída `MIN=3`. A CLI mostra mensagens de entrada e retorno, nos passos de entrada/retorno, o próximo PC mostrado já inclui a transição. Os testes comparam o banco final e o contexto restaurado com a execução sem interrupções.

## Verificações automatizadas

```sh
go test -count=1 ./...
go vet ./...
```

Os comandos servem para macOS e PowerShell. Foram executados no macOS com Go 1.27.2 em 17/11/2026. A execução/validação no Windows, o relatório e o ensaio ficaram fora desta sessão por solicitação da dupla. O [checklist](planejamento.md) registra o estado atual. A suíte cobre as operações, casos-limite, os programas reais, o mapa, as portas, vídeo por intervalo/término e interrupções repetidas com restauração do contexto.

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

Também é possível usar `go run ./cmd/simriscv`, pois o arquivo acima é o padrão. `-max-passos N` limita o número de instruções, o valor padrão é 100. `-video-a-cada N` define o intervalo de exibição da VRAM e deve ser positivo. `-entrada TEXTO` prepara a fila de entrada ASCII. `go run ./cmd/simriscv -help` mostra as opções.

Para compilar, use `go build -o simriscv ./cmd/simriscv` no macOS ou `go build -o simriscv.exe ./cmd/simriscv` no Windows.

## Feito por 👥
- [Luan](https://github.com/luanfelixcoding)
- [Felipe](https://github.com/Mafra013)
