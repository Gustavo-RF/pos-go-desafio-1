# Desafio pos Go 
Esse desafio consiste em ter duas apliações go:
- Server
- Client

O servidor é responsável por buscar as informações de cotação do dólar em https://economia.awesomeapi.com.br/json/last/USD-BRL, salvar em banco e retornar para o client através da rota /cotacao

O cliente é responsável por fazer a chamada ao server, mostrar no writer padrão o resultado e salvar em um arquivo txt

## Rodando o projeto
  go run server/server.go

## Fazendo a chamada pelo client
  go run client/client.go

Ao rodar o client, os resultados serão salvos em banco, arquivo e mostrado na tela para o usuário.
