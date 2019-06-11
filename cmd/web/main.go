package main

import (
	"fmt"
	"net/http"
	// "log"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	log "github.com/sirupsen/logrus"

	service "gitlab.com/standard-go/project/internal/app/project/service"
	"gitlab.com/standard-go/project/internal/app/config/env"
	api "gitlab.com/standard-go/project/internal/app/project/web-api.v1"
	"gitlab.com/standard-go/project/internal/pkg/server"
)

/*
====================================
    MAIN
====================================
*/

func main() {
	dbUsername 	:= env.Get("DB_USERNAME")
	dbPassword 	:= env.Get("DB_PASSWORD")
	dbHost 		:= env.Get("DB_HOST")
	dbPort 		:= env.Get("DB_PORT")
	dbName 		:= env.Get("DB_DATABASE")
	dbURL 		:= fmt.Sprintf("mongodb://%s:%s@%s:%s", dbUsername, dbPassword, dbHost, dbPort)

	srv, err := service.New(dbURL, dbName)
	if err != nil {
		log.WithFields(log.Fields{
			"tag"	: "create_db_connection_failed",
			"error"	: err,
		}).Fatal("Failed to connect to MongoDB")
	}

	defer srv.Close()
	api.SetService(srv)
	r := chi.NewRouter()

	// Basic CORS
	cors := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
	})
	r.Use(cors.Handler)
	// r.Use(api.CommonMiddleware)
	// r.Use(api.JWTMiddleware)

	r.Route("/api/v1", api.InitAPIRoutes)

	appHost := env.Get("APP_HOST")
	fmt.Println("Running at : " + appHost)

	/*
	   Gracefully stop application
	   reference : https://medium.com/@kpbird/golang-gracefully-stop-application-23c2390bb212
	*/
	s := server.New(&http.Server{
		Addr	: appHost,
		Handler : r,
	})
	s.Start()
	defer s.Stop()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	fmt.Println("Server is shutting down")
}
