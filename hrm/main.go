package main

import (
	"fmt"
	"log"
	"net"
	
	"practice/hack-and-solve/hrm/storage/postgres"
	"practice/hack-and-solve/utility"
	"practice/hack-and-solve/utility/logging"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

var (
	svcName = "hrm"
	version = "development"
)

func main() {
	log.Printf("starting %s service", svcName)
	cfg, err := utility.NewConfig("env/config")
	if err != nil {
		log.Fatal(err)
	}

	logger := logging.NewLogger(cfg).WithFields(logrus.Fields{
		"service": svcName,
		"version": version,
	})

	dbString := utility.NewDBString(cfg)
	db, err := postgres.NewStorage(dbString, logger)
	if err != nil {
		logger.WithError(err).Error("unable to connect DB")
		return
	}

	if err := db.RunMigration(cfg.GetString("database.migrationDir")); err != nil {
		logger.WithError(err).Error("migration failed")
	}

	if err := runGRPCServer(db, cfg); err != nil {
		logger.WithError(err).Error("unable to setup grpc service")
	}
}

func runGRPCServer(store *postgres.Storage, config *viper.Viper) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", config.GetInt("server.port")))
	if err != nil {
		return err
	}

	grpcServer := grpc.NewServer()
	log.Printf("Server %s management listening at: %+v", svcName, lis.Addr())
	if err := grpcServer.Serve(lis); err != nil {
		return err
	}
	return nil
}
