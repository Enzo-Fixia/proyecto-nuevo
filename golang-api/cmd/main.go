package main

import (
	"log"
	"net/http"

	"github.com/fixia/golang-api/config"
	"github.com/fixia/golang-api/database"
	"github.com/fixia/golang-api/internal/order"
	"github.com/fixia/golang-api/internal/product"
	"github.com/fixia/golang-api/internal/user"
	appmiddleware "github.com/fixia/golang-api/middleware"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()

	if cfg.AppEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	db := database.Connect(cfg)

	userRepo := user.NewRepository(db)
	userSvc := user.NewService(userRepo)
	userHandler := user.NewHandler(userSvc)

	productRepo := product.NewRepository(db)
	productSvc := product.NewService(productRepo)
	productHandler := product.NewHandler(productSvc)

	orderRepo := order.NewRepository(db)
	orderSvc := order.NewService(orderRepo)
	orderHandler := order.NewHandler(orderSvc)

	r := gin.New()
	r.Use(gin.Recovery(), appmiddleware.Logger())

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	v1 := r.Group("/api/v1")
	user.RegisterRoutes(v1, userHandler)
	product.RegisterRoutes(v1, productHandler)
	order.RegisterRoutes(v1, orderHandler)

	addr := ":" + cfg.AppPort
	log.Printf("🚀 Server running at %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
