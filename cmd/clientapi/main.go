package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"google.golang.org/grpc"

	portdomainv1 "port-location/api/proto/portdomain/v1"
	"port-location/internal/clientapi"
	"port-location/internal/clientapi/portdomain"
	"port-location/internal/clientapi/server"
)

func main() {
	confPath := os.Getenv("CONFIG_PATH")

	var conf clientapi.Config
	if err := conf.Read(confPath); err != nil {
		log.Fatal(err.Error())
		return
	}

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
		log.Fatalf("failed to parse file: %v", err)
	}

	log.Fatal(http.ListenAndServe(":"+conf.HTTPServer.Port, s.Router))
}
