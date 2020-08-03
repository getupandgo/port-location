package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"google.golang.org/grpc"

	portdomainv1 "port-location/api/proto/portdomain/v1"
	"port-location/internal/clientapi"
	"port-location/internal/clientapi/portdomain"
	"port-location/internal/clientapi/server"
)

// nolint: gochecknoglobals
var conf = clientapi.Config{
	HTTPServer: clientapi.HTTPServer{
		Host: "localhost",
		Port: "8000",
	},
	GRPCServer: clientapi.GRPCServer{
		Host: "localhost",
		Port: "9000",
	},
	PortFilePath: "./port_data/ports.json",
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	serverAddr := fmt.Sprintf(`%s:%s`, conf.GRPCServer.Host, conf.GRPCServer.Port)

	conn, err := grpc.DialContext(ctx, serverAddr, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()

	portDomainClient := portdomain.NewClient(portdomainv1.NewPortDomainAPIClient(conn))

	s := server.NewServer(portDomainClient)
	if err := s.ParsePortFile(context.Background(), conf.PortFilePath); err != nil {
		log.Fatalf("fail to dial: %v", err)
	}

	log.Fatal(http.ListenAndServe(":"+conf.HTTPServer.Port, s.Router))
}
