package clientapi

import (
	"context"
	"log"
	"net/http"
	"time"

	"google.golang.org/grpc"

	portdomainV1 "port-location/api/proto/portdomain/v1"
	"port-location/internal/clientapi"
	"port-location/internal/clientapi/portdomain"
	"port-location/internal/clientapi/server"
)

func main() {
	conf := clientapi.Config{
		clientapi.HTTPServer{},
		clientapi.GRPCServer{},
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	conn, err := grpc.DialContext(ctx, conf.GRPCServer.Host+conf.GRPCServer.Port, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()

	portDomainClient := portdomain.NewClient(portdomainV1.NewPortDomainAPIClient(conn))

	s := server.NewServer(portDomainClient)
	log.Fatal(http.ListenAndServe(":"+conf.HTTPServer.Port, s.Router))
}
