# bodocongo-torrent

## GRUPO

- Samuel Cabral de Luna - 121210376
- Victor Vinicius Freire de Araújo - 1211110361
- Paola Kathrein Moura Marques - 121111034
- Marcos Antônio Cardoso Pereira - 121210927

## COMO EXECUTAR

## Setup inicial

Edite o arquivo `peers.txt` e adicione os IPs e portas das máquinas que participarão da execução do programa, exceto o IP da própria máquina. Cada linha deve conter um IP e uma porta no formato `IP:PORTA`.

- Exemplo de `peers.txt`:

```
192.168.1.2:8080
192.168.1.3:8080
192.168.1.4:8080
```
## Executando o programa

Execute o programa em máquinas diferentes utilizando o comando abaixo:

```
go run *.go
```

## Buscando um Hash de Arquivo

Após a execução do programa, você pode buscar um hash de arquivo digitando o `filehash` no terminal.

O programa irá procurar o hash informado entre os peers configurados.

