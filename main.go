package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/lugingf/vestiarie-test-task/internal/domain"
	"github.com/lugingf/vestiarie-test-task/internal/handler"
	"github.com/lugingf/vestiarie-test-task/internal/storage"
	"github.com/lugingf/vestiarie-test-task/resources"
	"log"
	"net/http"
	"time"
)

func main() {
	di := resources.Init()

	payoutStorage := storage.NewPayoutStorageSQL(di.SQLShard)
	itemStorage := storage.NewItemStorageSQL(di.SQLShard)
	s := domain.NewPayoutService(&payoutStorage, &itemStorage)

	ph, err := handler.NewPayoutHandler(s)
	if err != nil {
		log.Fatal("failed to init handler")
	}

	h := router(ph)
	listenAddress := fmt.Sprintf(":%v", di.AppConfig.Server.Port)
	server := &http.Server{
		Addr:        listenAddress,
		Handler:     h,
		ReadTimeout: 1 * time.Second,
	}

	runHTTPServer(server)
}

func runHTTPServer(server *http.Server) {
	log.Println(fmt.Sprintf("starting HTTP server: listening on %s", server.Addr))

	err := server.ListenAndServe()
	if err != http.ErrServerClosed {
		log.Fatal("failed to listen server")
	}
}

func router(ph *handler.PayoutHandler) http.Handler {
	r := mux.NewRouter()
	payoutRoute := r.Path("/payouts").Subrouter()
	payoutRoute.Methods(http.MethodPost).HandlerFunc(ph.PostPayouts)
	return r
}
