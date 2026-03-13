package httpserver

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/cihub/seelog"
)

type HttpServer struct {
	Ctx context.Context
	srv *http.Server

	debug   bool
	running bool
}

func NewHttpServer(port int) *HttpServer {
	serv := new(HttpServer)
	serv.srv = new(http.Server)
	serv.srv.Addr = fmt.Sprintf(":%d", port)
	serv.srv.Handler = InitRoute()
	serv.srv.ReadTimeout = 120 * time.Second
	serv.srv.WriteTimeout = 120 * time.Second
	serv.Ctx = context.Background()
	serv.running = true
	return serv
}

func (h *HttpServer) Start() {
	seelog.Info(context.Background(), "http server info ", h.srv.Addr)
	go func(h *HttpServer) {
		err := h.srv.ListenAndServe()
		if err != nil {
			seelog.Warnf("listenAndServe error : %s\n", err.Error())
			return
		}
	}(h)

}

func (h *HttpServer) Stop() {

	err := h.srv.Shutdown(h.Ctx)
	if err != nil {
		seelog.Error("http server shutdown error : ", err.Error())
	}
}

func (h *HttpServer) GetProcessName() string {
	return "httpserver"
}
