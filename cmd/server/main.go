package main

import (
	"database/sql"
	"fmt"
	"net"
	"net/http"

	graphql_handler "github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gbgomes/GoExpert/CleanArch/configs"
	"github.com/gbgomes/GoExpert/CleanArch/internal/event/handler"
	"github.com/gbgomes/GoExpert/CleanArch/internal/infra/database"
	"github.com/gbgomes/GoExpert/CleanArch/internal/infra/graph"
	"github.com/gbgomes/GoExpert/CleanArch/internal/infra/grpc/pb"
	"github.com/gbgomes/GoExpert/CleanArch/internal/infra/grpc/service"
	"github.com/gbgomes/GoExpert/CleanArch/internal/infra/web/webserver"
	"github.com/gbgomes/GoExpert/CleanArch/internal/usecase"
	"github.com/gbgomes/GoExpert/CleanArch/pkg/events"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	_ "github.com/gbgomes/GoExpert/CleanArch/docs"

	"github.com/streadway/amqp"

	// mysql
	_ "github.com/go-sql-driver/mysql"
)

// @title           Clean Architeture project API Documentation
// @version         1.0
// @description     Clean Architeture API Documentation
// @termsOfService  http://swagger.io/terms/

// @contact.name   Guilherme Gomes
// @contact.url    http://github.com/gbgomes
// @contact.email  guilherme.gomes@gmail.com

// @license.name   Free License
// @license.url    http://github.com/gbgomes

// @host      localhost:8000
// @BasePath  /
func main() {
	configs, err := configs.LoadConfig(".")
	if err != nil {
		panic(err)
	}

	db, err := sql.Open(configs.DBDriver, fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", configs.DBUser, configs.DBPassword, configs.DBHost, configs.DBPort, configs.DBName))
	if err != nil {
		panic(err)
	}
	defer db.Close()

	rabbitMQChannel := getRabbitMQChannel(configs)

	eventDispatcher := events.NewEventDispatcher()
	eventDispatcher.Register("OrderCreated", &handler.OrderCreatedHandler{
		RabbitMQChannel: rabbitMQChannel,
	})

	webserver := webserver.NewWebServer(configs.WebServerPort)
	webOrderHandler := NewWebOrderHandler(db, eventDispatcher)
	webserver.AddHandler("POST /order", webOrderHandler.Create)
	webserver.AddHandler("GET /order", webOrderHandler.List)
	fmt.Println("Starting web server on port", configs.WebServerPort)
	go webserver.Start()

	createOrderUseCase := NewCreateOrderUseCase(db, eventDispatcher)
	OrderRepository := database.NewOrderRepository(db)
	listOrdersUseCase := usecase.NewListOrdersUseCase(OrderRepository)
	grpcServer := grpc.NewServer()
	createOrderService := service.NewOrderService(*createOrderUseCase, *listOrdersUseCase)
	pb.RegisterOrderServiceServer(grpcServer, createOrderService)
	reflection.Register(grpcServer)

	fmt.Println("Starting gRPC server on port", configs.GRPCServerPort)
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", configs.GRPCServerPort))
	if err != nil {
		panic(err)
	}
	go grpcServer.Serve(lis)

	srv := graphql_handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{
		CreateOrderUseCase: *createOrderUseCase,
		ListOrderUseCase:   *listOrdersUseCase,
	}}))
	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)
	fmt.Println("Starting GraphQL server on port", configs.GraphQLServerPort)
	http.ListenAndServe(":"+configs.GraphQLServerPort, nil)
}

func getRabbitMQChannel(configs *configs.Conf) *amqp.Channel {
	conn, err := amqp.Dial(fmt.Sprintf("%s://%s:%s@%s:%s/", configs.MQDriver, configs.MQUser, configs.MQPassword, configs.MQHost, configs.MQPort))
	if err != nil {
		panic(err)
	}
	ch, err := conn.Channel()
	if err != nil {
		panic(err)
	}
	return ch
}
