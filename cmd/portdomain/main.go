package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"

	portdomainv1 "port-location/api/proto/portdomain/v1"
	"port-location/internal/portdomain"
	"port-location/internal/portdomain/server"
	"port-location/internal/portdomain/storage"
)

func main() {
	confPath := os.Getenv("CONFIG_PATH")

	var conf portdomain.Config
	if err := conf.Read(confPath); err != nil {
		log.Fatal(err.Error())
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	psqlConf := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		conf.DB.Host, conf.DB.Port, conf.DB.Username, conf.DB.Password, conf.DB.Name,
	)

	db, err := sqlx.ConnectContext(ctx, "postgres", psqlConf)
	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}
	defer db.Close()

	driver, err := postgres.WithInstance(db.DB, &postgres.Config{})
	if err != nil {
		log.Fatalf("failed to apply migrations: %v", err)
	}

	m, err := migrate.NewWithDatabaseInstance(conf.DB.MigrationsDir, conf.DB.Name, driver)
	if err != nil {
		log.Fatalf("failed to apply migrations: %v", err)
	}

	if err := m.Up(); err != nil {
		if err != migrate.ErrNoChange {
			log.Fatalf("failed to apply migrations: %v", err)
		}
	}

	sc := storage.NewClient(db)
	portDomainServer := server.NewServer(sc)

	grpcServerConf := fmt.Sprintf("%s:%s", conf.GRPCServer.Host, conf.GRPCServer.Port)

	lis, err := net.Listen("tcp", grpcServerConf)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	portdomainv1.RegisterPortDomainAPIServer(grpcServer, portDomainServer)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to start grpc server: %v", err)
	}
}
