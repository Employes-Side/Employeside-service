package main

import (
	"context"
	"database/sql"
	"errors"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/go-sql-driver/mysql"

	employeside "github.com/Employes-Side/employee-side"
	"github.com/Employes-Side/employee-side/internal/endpoints"
	"github.com/Employes-Side/employee-side/internal/handlers"
	"github.com/Employes-Side/employee-side/internal/repositories"
	"github.com/gorilla/mux"
	"k8s.io/klog"
)

func main() {

	cfg := employeside.LoadConfiguration()
	dbConn, err := sql.Open("mysql", cfg.DB.CreateDSN())
	if err != nil {
		panic(err)
	}

	defer dbConn.Close()

	var userManager repositories.UserRepository
	{
		userManager = *repositories.NewManager(dbConn)
	}

	UserEndpoint := endpoints.NewUserEndpoint(userManager)

	router := mux.NewRouter()
	{
		appHandler := handlers.NewHandler(UserEndpoint)

		router.PathPrefix("/").Handler(appHandler)

	}

	httpServer := &http.Server{Handler: router}
	{
		lis, err := net.Listen("tcp", cfg.Bind.HTTP)
		if err != nil {
			panic(err)
		}

		klog.Infof("starting http server on %q", lis.Addr())
		go func() {
			if err := httpServer.Serve(lis); err != nil && !errors.Is(err, http.ErrServerClosed) {
				panic(err)
			}
		}()
	}

	defer func() {
		klog.Info("closing http server...")
		klog.Infof("http server close error: %v",
			httpServer.Shutdown(context.Background()))
	}()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	klog.Infof("received shutdown signal %q", <-sig)
	klog.Infof("bye")
}
