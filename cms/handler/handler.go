package handler

import (
	"io/fs"
	"net/http"

	"practice/hack-and-solve/cms/conn"
	"practice/hack-and-solve/cms/handler/common"

	"github.com/benbjohnson/hashfs"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Svc struct {
	*common.Server
}

func NewServer(
	config *viper.Viper,
	logger *logrus.Entry,
	assets fs.FS,
	decoder *schema.Decoder,
	conn *conn.Conn,
) (*mux.Router, error) {
	cs := &common.Server{
		Config:     config,
		Logger:     logger,
		Assets:     hashfs.NewFS(assets),
		Decoder:    decoder,
	}

	if err := cs.ParseTemplates(); err != nil {
		return nil, err
	}

	r := mux.NewRouter()
	r.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", common.CacheStaticFiles(http.FileServer(http.FS(cs.Assets)))))

	r.NotFoundHandler = cs.GetErrorHandler()
	return r, nil
}
