package main

import (
	"database/sql"
	"fmt"
	"log"
	"net"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"

	portdomainv1 "port-location/api/proto/portdomain/v1"
	"port-location/internal/portdomain"
	"port-location/internal/portdomain/server"
	"port-location/internal/portdomain/storage"
)

// nolint: gochecknoglobals
var conf = portdomain.Config{
	GRPCServer: portdomain.GRPCServer{
		Host: "localhost",
		Port: "9000",
	},
	DB: portdomain.DB{
		Host:          "localhost",
		Port:          "5432",
		Name:          "postgres",
		Username:      "postgres",
		Password:      "example",
		MigrationsDir: "file:///Users/getupandgo/projects/port-location/migrations",
	},
}

func main() {
	psqlConf := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		conf.DB.Host, conf.DB.Port, conf.DB.Username, conf.DB.Password, conf.DB.Name,
	)

	db, err := sql.Open("postgres", psqlConf)
	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}
	defer db.Close()

	if err = db.Ping(); err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})
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
