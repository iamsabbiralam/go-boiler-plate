package main

import (
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"

	"practice/hack-and-solve/cms/conn"
	"practice/hack-and-solve/cms/handler"
	"practice/hack-and-solve/utility"
	"practice/hack-and-solve/utility/logging"

	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/spf13/viper"
	"github.com/urfave/negroni"
	"github.com/yookoala/realpath"
)

const (
	svcName = "CMS"
	version = "1.0.0"
)

func main() {
	cfg, err := utility.NewConfig("env/config")
	if err != nil {
		log.Fatal(err)
	}
	env := cfg.GetString("runtime.environment")
	logger := logging.NewLogger(cfg).WithFields(logrus.Fields{
		"environment": env,
		"service":     svcName,
		"version":     version,
	})

	conns := conn.NewConns(logger, cfg)
	defer conns.Close()

	server, err := setupServer(logger, cfg, conns)
	if err != nil {
		logger.WithError(err).Error("failed to run setup server")
	}

	server.Use(func(h http.Handler) http.Handler {
		recov := negroni.NewRecovery()
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			recov.ServeHTTP(w, r, h.ServeHTTP)
		})
	})

	l, err := net.Listen("tcp", ":"+cfg.GetString("server.port"))
	if err != nil {
		logger.WithError(err).Error("failed to listen server")
	}

	log.Printf("starting %s service on %s", svcName, l.Addr())
	if err := http.Serve(l, server); err != nil {
		logger.WithError(err).Error("failed to serve server")
	}
}

func setupServer(logger *logrus.Entry, cfg *viper.Viper, conn *conn.Conn) (*mux.Router, error) {
	decoder := schema.NewDecoder()
	decoder.IgnoreUnknownKeys(true)

	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	assetPath, err := realpath.Realpath(filepath.Join(wd, "assets"))
	if err != nil {
		return nil, err
	}
	asst := afero.NewIOFS(afero.NewBasePathFs(afero.NewOsFs(), assetPath))

	srv, err := handler.NewServer(cfg, logger, asst, decoder, conn)
	if err != nil {
		return nil, err
	}
	return srv, nil
}

type rt func(*http.Request) (*http.Response, error)

func (t rt) RoundTrip(req *http.Request) (*http.Response, error) { return t(req) }

func fakeTLS(r http.RoundTripper) http.RoundTripper {
	return rt(func(req *http.Request) (*http.Response, error) {
		rq := req.Clone(req.Context())
		rq.Header.Set("X-Forwarded-Proto", "https")
		return r.RoundTrip(rq)
	})
}
