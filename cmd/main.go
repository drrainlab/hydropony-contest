// Hydropony 🦄 contest

package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/drrainlab/hydropony-contest/cacher"
	"github.com/drrainlab/hydropony-contest/controllers"
	"github.com/drrainlab/hydropony-contest/models"
	"github.com/gin-gonic/gin"
)

func main() {

	InitConfig()

	_, srv := runServer()

	// graceful shutdown

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// give some time to finish processing request

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")

}

func runServer() (router *gin.Engine, srv *http.Server) {
	// init database connection
	err := models.InitDBConnection(GetConfig().DBDSN)
	if err != nil {
		panic("DB connection error")
	}

	cacher.InitCache(GetConfig().CacherDSN)

	// init Gin router
	router = gin.Default()
	// init routes
	router.GET("/users", controllers.GetUserInfoHandler)
	router.POST("/users", controllers.CreateUserHandler)

	// http server settings
	srv = &http.Server{
		Addr:    fmt.Sprintf(":%s", GetConfig().Port),
		Handler: router,
	}

	// starting server in goroutine
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	return router, srv
}
