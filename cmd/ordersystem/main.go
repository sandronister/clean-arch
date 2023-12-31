package main

import (
	"database/sql"
	"fmt"
	"net"
	"net/http"

	graphql_handler "github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/sandronister/clean-arch/configs"
	"github.com/sandronister/clean-arch/internal/infra/event/handler"
	"github.com/sandronister/clean-arch/internal/infra/graphql/graph"
	"github.com/sandronister/clean-arch/internal/infra/grpc/pb"
	"github.com/sandronister/clean-arch/internal/infra/grpc/services"
	"github.com/sandronister/clean-arch/internal/infra/web/webserver"
	"github.com/sandronister/clean-arch/pkg/events"
	"github.com/streadway/amqp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	//mysql
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	config, err := configs.LoadConfig(".")
	if err != nil {
		panic(err)
	}

	db, err := sql.Open(config.DBDriver, fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", config.DBUser, config.DBPassword, config.DBHost, config.DBPort, config.DBName))
	if err != nil {
		panic(err)
	}

	defer db.Close()

	rabbitMQChannel := getRabbitMQChannel(config)
	eventDispatcher := events.NewEventDispatcher()
	eventDispatcher.Register("OrderCreated", &handler.OrderCreatedHandler{
		RabbitMQChannel: rabbitMQChannel,
	})

	createOrderUseCase := NewCreateOrderUseCase(db, eventDispatcher)
	listOrderUserCase := NewListOrderUseCase(db)

	webServer := webserver.NewWebServer(config.WebServerPort)
	webOrderHandler := NewWebOrderHandler(db, eventDispatcher)
	webServer.AddHandlers("/orders", http.MethodPost, webOrderHandler.Create)
	webServer.AddHandlers("/orders", http.MethodGet, webOrderHandler.List)
	fmt.Println("Starting web server on port", config.WebServerPort)
	go webServer.Start()

	grpcServer := grpc.NewServer()
	createOrderService := services.NewOrderService(*createOrderUseCase, *listOrderUserCase)
	pb.RegisterOrderServiceServer(grpcServer, createOrderService)
	reflection.Register(grpcServer)

	fmt.Println("Starting gRPC on port ", config.GRPCServerPort)
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", config.GRPCServerPort))

	if err != nil {
		panic(err)
	}

	go grpcServer.Serve(lis)

	srv := graphql_handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{
		CreateOrderUseCase: *createOrderUseCase,
		ListOrderUseCase:   *listOrderUserCase,
	}}))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	fmt.Println("Starting GraphQL server on port", config.GraphQLServerPort)
	http.ListenAndServe(":"+config.GraphQLServerPort, nil)

}

func getRabbitMQChannel(config *configs.Conf) *amqp.Channel {
	conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%s/", config.RabbitMQUser, config.RabbitMQPassword, config.RabbitMQHost, config.RabbitMQPort))
	if err != nil {
		panic(err)
	}
	ch, err := conn.Channel()
	if err != nil {
		panic(err)
	}

	return ch
}
