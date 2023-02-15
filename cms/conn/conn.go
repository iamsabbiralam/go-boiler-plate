package conn

import (
	"log"
	"practice/hack-and-solve/utility/middleware"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

type Conn struct {
	Hrm *grpc.ClientConn
}

func NewConns(logger *logrus.Entry, cfg *viper.Viper) *Conn {
	log.Printf("starting to dialing hrm service port: %s", cfg.GetString("services.hrmURL"))
	opts := getGRPCOpts(cfg)
	hrm, err := grpc.Dial(cfg.GetString("services.hrmURL"), opts...)
	if err != nil {
		logger.WithError(err).Fatal("failed to dial hrm service")
	}

	return &Conn{
		Hrm: hrm,
	}
}

func (co *Conn) Close() {
	co.Hrm.Close()
}

func getGRPCOpts(cnf *viper.Viper) []grpc.DialOption {
	opts := []grpc.DialOption{grpc.WithInsecure()}
	opts = append(opts,
		grpc.WithBlock(), grpc.WithChainUnaryInterceptor(
			middleware.AuthForwarder(),
		))
	return opts
}
