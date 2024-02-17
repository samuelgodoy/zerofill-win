# ZeroFiller

**Nota:** Este programa foi desenvolvido para funcionar apenas no sistema operacional Windows.

O ZeroFiller é uma ferramenta de linha de comando (CLI) dedicada a realizar preenchimento zero em discos, oferecendo um ajuste flexível de throughput. Seu objetivo é permitir o preenchimento de discos com grande capacidades mas de baixa qualidade, ao mesmo tempo em que permite o ajuste da taxa de escrita para reduzir o estresse térmico do dispositivo evitando corromper, travamentos ou ejeções.

## Requisitos

- Sistema operacional Windows

## Instalação

Baixe o arquivo executável `zerofill-win.exe` e execute como administrador.

## Uso

### Listar todos os discos físicos disponíveis

```bash
zerofill-win.exe list
```

Este comando lista todos os discos físicos disponíveis no sistema.

### Preencher um disco com zeros

```bash
zerofill-win.exe zero --device [caminho_do_dispositivo]
```

- \[caminho_do_dispositivo\]: O caminho do dispositivo que você deseja preencher com zeros. Você pode obter o caminho do dispositivo usando o comando `list`.

### Ajustar o throughput durante o preenchimento zero

Durante o processo de preenchimento zero, você pode ajustar o throughput pressionando as teclas `+` e `-`:

- Pressione `+` para aumentar a taxa de preenchimento em 100KB/s.
- Pressione `-` para diminuir a taxa de preenchimento em 100KB/s.

## Observações

- Tenha cuidado ao selecionar o dispositivo alvo.
- Pode não ser totalmente seguro contra ferramentas de recuperação de arquivos. 
- Certifique-se de ter permissões adequadas para acessar e escrever no dispositivo selecionado.
- Ajustar o throughput pode ajudar a reduzir o estresse térmico do disco, mas leve em consideração a capacidade de refrigeração do dispositivo.

## Exemplo

```bash
zerofill-win.exe zero --device '\\.\PHYSICALDRIVE2'
```

Este comando iniciará o processo de preenchimento zero no disco físico referenciado pelo caminho `\\.\PHYSICALDRIVE2`.
