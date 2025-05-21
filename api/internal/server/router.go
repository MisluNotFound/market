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
	productRouter := s.engine.Group("/api/product")
	assertRouter := s.engine.Group("/api/assert")
	orderRouter := s.engine.Group("/api/order")
	searchRouter := s.engine.Group("/api/search")
	conversationRouter := s.engine.Group("/api/conversation")

	// setup routers
	s.registerUserGroup(userRouter)
	s.registerMockGroup(mockRouter)
	s.registerProductGroup(productRouter)
	s.registerAssertGroup(assertRouter)
	s.registerOrderGroup(orderRouter)
	s.registerSearchGroup(searchRouter)
	s.registerConversationGroup(conversationRouter)

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

func (s *Server) registerProductGroup(group *gin.RouterGroup) {
	// TODO add auth
	group.POST("/:userID", controllers.CreateProduct())
	group.GET("/:userID/:productID", controllers.GetProduct())
	group.PUT("/:userID/:productID", controllers.UpdateProduct())
	group.PUT("/:userID/:productID/off-shelves", controllers.OffShelves())
	group.PUT("/:userID/:productID/on-shelves", controllers.OnShelves())
	group.PUT("/:userID/:productID/sold", controllers.SoldOut())
	group.PUT("/:userID/:productID/selling", controllers.NotSold())
	group.GET("/:userID", controllers.GetUserProducts())
	group.GET("/products", controllers.GetProductList())
	group.GET("/category", controllers.GetAllCategory())
	group.PUT("/:userID/:productID/price", controllers.UpdateProductPrice())
}

func (s *Server) registerAssertGroup(group *gin.RouterGroup) {
	// TODO add auth
	group.GET("/:type/:owner/:key", controllers.GetAssert())
}

func (s *Server) registerOrderGroup(group *gin.RouterGroup) {
	group.POST("/:userID/:productID", controllers.PurchaseProduct())
	group.GET("/:userID/list", controllers.GetOrderList())
	group.GET("/:userID/:orderID", controllers.GetOrder())
	group.PUT("/shipped/:userID/:orderID", controllers.ConfirmOrderShipped())
	group.PUT("/signed/:userID/:orderID", controllers.ConfirmOrderSigned())
	group.PUT("/pay/:userID/:orderID", controllers.PayOrder())
	group.GET("/:userID/status", controllers.GetAllOrderStatus())
	group.POST("/refund/:userID/:orderID", controllers.RefoundOrder())
	group.PUT("/cancel/:userID/:orderID", controllers.CancelOrder())
}

func (s *Server) registerSearchGroup(group *gin.RouterGroup) {
	group.POST("/products", controllers.SearchProduct())
}

func (s *Server) registerConversationGroup(group *gin.RouterGroup) {
	group.POST("/create", controllers.CreateConversation())
	group.GET("/:userID", controllers.GetConversationList())
	group.GET("/messages", controllers.GetMessages())
	group.GET("/:userID/list", controllers.GetConversationList())
}
