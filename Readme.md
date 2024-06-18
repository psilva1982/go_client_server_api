# Client-Server API

Neste repositório possuímos duas aplicações:
- **server.go** 
- **client.go**

## server.go

Esta aplicação é responsável por consumir a API contendo o câmbio de Dólar e Real no endereço: https://economia.awesomeapi.com.br/json/last/USD-BRL e em seguida retorna no formato JSON o resultado para o **client.go**.

Este serviço roda no endereço http://localhost:8080/cotacao

Além disso, o **server.go**  registra no banco de dados SQLite cada cotação recebida, sendo que o **timeout** máximo para chamar a **API** de cotação do dólar deverá ser de **200ms** e o timeout máximo para conseguir persistir os dados no **banco** deverá ser de **10ms**.

## client.go

A aplicação recebe do **server.go** apenas o valor atual do câmbio (campo "bid" do JSON), e através o package "context", tem um timeout máximo de 300ms para receber o resultado do **server.go**.


### Executando

Efetuando o clone do repositório

```
git clone https://github.com/psilva1982/go_client_server_api.git
```

Baixando os módulos necessários

```
git mod tidy
```

Executando o **server/server.go**

```
go run server/server.go
```

Em um outro terminal execute o **client/client.go**

```
go run client/client.go
```

Visualize os resultados no SQLite e no arquivo **cotacao.txt**