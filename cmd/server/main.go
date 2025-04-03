package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	// "github.com/go-kit/kit/transport/grpc"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/golang-migrate/migrate/source/file"

	employeside "github.com/Employes-Side/employee-side"
	"github.com/Employes-Side/employee-side/internal/endpoints"
	"github.com/Employes-Side/employee-side/internal/handlers"
	"github.com/Employes-Side/employee-side/internal/repositories"
	"github.com/gorilla/mux"
	"k8s.io/klog"
)

func main() {

	cfg := employeside.LoadConfiguration()

	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?parseTime=true",
		cfg.DB.Username,
		cfg.DB.Password,
		cfg.DB.Hostname,
		cfg.DB.Port,
		cfg.DB.Database,
	)

	migrateSQL(dsn, cfg.DB.Database)

	dbConn, err := sql.Open("mysql", cfg.DB.CreateDSN())
	if err != nil {
		panic(err)
	}

	defer dbConn.Close()

	router := mux.NewRouter()

	var userManager repositories.UserRepository
	{
		userManager = *repositories.NewManager(dbConn)
	}

	var writerManager repositories.WriterRepository
	{
		writerManager = *repositories.NewWriterManager(dbConn)
	}

	var blogManager repositories.BlogRepository
	{
		blogManager = *repositories.NewBlogManager(dbConn)
	}

	var modulesManager repositories.ModulesRepository
	{
		modulesManager = *repositories.NewModulesManger(dbConn)
	}

	protectedRoutes := router.PathPrefix("/").Subrouter()
	protectedRoutes.Use(jwtMiddlewareHandler)

	userEndpoint := endpoints.NewUserEndpoint(userManager)
	{
		handlers.NewHandler(protectedRoutes, userEndpoint)
	}

	writerEndpoint := endpoints.NewWriterEndpoint(writerManager)
	{
		handlers.NewWriterHandler(protectedRoutes, writerEndpoint)
	}

	blogEndpoint := endpoints.NewBlogEndpoint(blogManager)
	{
		handlers.NewBlogHandler(protectedRoutes, blogEndpoint)
	}

	modulesEndpoint := endpoints.NewModuleEndpoint(modulesManager)
	{
		handlers.NewModuleHandler(protectedRoutes, modulesEndpoint)
	}

	err = router.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		path, err := route.GetPathTemplate()
		if err != nil {
			return nil
		}

		methods, err := route.GetMethods()
		if err != nil {
			return nil
		}

		klog.Infof("\t%v %s\n", methods, path)

		return nil
	})
	if err != nil {
		klog.Errorf("cannot print routes: %v", err)
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
