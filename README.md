Para executar o sistema:
    iniciar o doker com o mysql e o rabbitMQ:
        Na pasta raiz do projeto, executar 
            docker-compose up -d

    na pasta cmd/server
        go run main.go wire_gen.go

Serviços disponíveis no localhost nas seguintes portas:
    web server Rest: 8000
    gRPC server    : 50051
    GraphQL server : 8080

Uso da API Rest consultar swagger em:
    http://localhost:8000/docs/index.html


Serviços gRPC:
    evans -r repl
    package pb
    service <ServicoDesejado>
        service OrderService
    call <serviço>
        call CreateOrder
        call ListOrders

Serviços GraphQL
    mutation createOrder {
      createOrder(
        input: {id: "sql1", Price: 100, Tax:10}
      ) {
        id
        Price
        Tax
        FinalPrice
      }
    }

    query queryOrders {
      orders {
        id
        Price
        Tax
        FinalPrice
      }
    }


