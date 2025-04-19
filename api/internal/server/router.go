package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mislu/market-api/internal/server/controllers"
	"github.com/mislu/market-api/internal/utils/app"
)

func (s *Server) serve() {

	userRouter := s.engine.Group("/api/user")
	mockRouter := s.engine.Group("/api/mock")

	// setup routers
	s.registerUserGroup(userRouter)
	s.registerMockGroup(mockRouter)

	// run
	srv := &http.Server{
		Addr:    app.GetConfig().Server.Addr,
		Handler: s.engine,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			panic(err)
		}
	}()
	log.Printf("Server started at %s", app.GetConfig().Server.Addr)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	log.Println("Server exiting")
}

func (s *Server) registerUserGroup(group *gin.RouterGroup) {
	// TODO add auth

	group.POST("/register", controllers.CreateUser())
	group.POST("/login", controllers.Login())
	group.GET("/:userID", controllers.GetUserInfo())
	group.PUT("/:userID/avatar", controllers.UploadAvatar())
	group.PUT("/:userID/basic", controllers.UpdateBasic())
	group.PUT("/:userID/password", controllers.UpdatePassword())
}

func (s *Server) registerMockGroup(group *gin.RouterGroup) {
	group.POST("", controllers.MockPost())
	group.GET("", controllers.MockGet())
	group.GET("/error", controllers.MockError())
}
