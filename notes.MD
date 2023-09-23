
Migrate
Instalar os pacotes do migrate
    https://github.com/golang-migrate/migrate

rodar o comando abaixo para inciar a estrutura de migration e em seguida prencher oa arquivos comos sqls de criação e drops
    migrate create -ext=sql -dir=sql/migrations -seq init

rodar o comando abaixo para iniciar o BD no docker (precisa das definições em docker-compose.yaml)
    docker-compose up -d

rodar o comando abaixo para entrar no docker do mysql
    docker-compose exec mysql bash
    mysql -uroot -p orders

rodar o comando abaixo para adicionar o driver do mysql
    go install -tags mysql github.com/golang-migrate/migrate/v4/cmd/migrate

rodar o comando abaixo para inciar o migrate
    migrate -path=sql/migrations -database="mysql://root:root@tcp(localhost:3306)/orders" -verbose up 


Wire
precisa instalar o goolgle wire
Ele resolve as dependencias sem usar reflexion. Ele gera o código e com isso não é necessário usr um framework durante a execução
    https://github.com/google/wire

Para instalar:
    go install github.com/google/wire/cmd/wire@latest

Em seguida, cria o arquivo wire.go na pasta do server com o mapeamento das dependencias
Nesta mesma pasta, executar o comando abaixo
    wire
O arquivo wire_gen.go será gerado contendo os atalhos de dependencias


gRPC
ir para a pasta internal/infra/grpc e passar a usar esta pasta como base
criar a pasta protofiles, e nela o arquivo order.proto (notar a configuração de pasta a geração dos arquivos)
rodar o comando
    protoc --go_out=. --go-grpc_out=. protofiles/order.proto
implementar a arquivo de serviço: order_service.go

Usando o gRPC com o evans
    evans -r repl
    package pb
    service <ServicoDesejado>
        service OrderService
    call <serviço>
        call CreateOrder
        call ListOrders

graphQL
criar o arquivo tools.go na pasta raiz do projeto, com o conteúdo necessário
Na pasta internal/infra, executar o comando:
    go run github.com/99designs/gqlgen init
que irá inicializar toda a estrutura do graphQL

Editar o arquivo schema.graphqls com as definições das queries e estruturas
executar (na pasta onde está o arquivo de squema)
    go run github.com/99designs/gqlgen generate
para atualizar toda a estrutura após alterações

Editar o arquivo resolver.go, incluindo o usecase na estrutura 
Editar o arquivo schema.resolvers.go para implementar os serviços


Swagger
instalar o swag do site https://github.com/swaggo/swag
copiar/criar a documentação do projeto no main.go
na pasta raiz do projeto, rodar o comando swag init -g cmd/server/main.go
no webserver.go incluir o httpSwagger "github.com/swaggo/http-swagger" nos imports do main.go
no main.go incluir o _ "github.com/gbgomes/GoExpert/APIS/docs"
no webserver.go adicionar o Handler para a doc
    s.Router.Get("/docs/*", httpSwagger.Handler(httpSwagger.URL("http://localhost:8000/docs/doc.json")))



Para executar o sistema
    go run main.go wire_gen.go

