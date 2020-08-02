package main

import (
	"context"
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
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	conn, err := grpc.DialContext(ctx, conf.GRPCServer.Host+conf.GRPCServer.Port, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()

	portDomainClient := portdomain.NewClient(portdomainv1.NewPortDomainAPIClient(conn))

	s := server.NewServer(portDomainClient)
	log.Fatal(http.ListenAndServe(":"+conf.HTTPServer.Port, s.Router))
}
