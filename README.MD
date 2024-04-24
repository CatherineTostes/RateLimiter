# Rate Limiter Middleware

Esssa aplicação tem como objetivo realizar uma consulta na rota `https://economia.awesomeapi.com.br/json/last/` passando `USD-BRL` e `EUR-BRL`, para retornar o valor do câmbio.

Dentro dos limites estabelecidos, quando é verificado por IP a quantidade máxima de requisições é de 10.
Já quando é enviado pelo header `API_TOKEN` a quantidade máxima de requisições é de 12 em um limite de 10 segundos.

Para subir a aplicação é necessário rodar o comando:
```
docker-compose up -d
```

Para rodar os testes é necessário rodar o comando:
```
go test ./...
```

## Verificação por IP
Para realizar a requisição é necessário bater na rota:
```
http://localhost:8080/exchange
```

exemplo de chamada com loop:
```
for i in {1..11}; do curl -X GET http://localhost:8080/exchange; done
```

## Verificação por API_TOKEN
Para realizar a requisição é necessário bater na rota:
```
- Header: API_TOKEN: 123456
http://localhost:8080/exchange
```

exemplo de chamada com loop:
```
for i in {1..13}; do curl --header "API_TOKEN:abc123" -X GET http://localhost:8080/exchange; done
```