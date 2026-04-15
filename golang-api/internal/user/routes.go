package user

import (
	"github.com/fixia/golang-api/internal/auth"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(rg *gin.RouterGroup, h *Handler) {
	authGroup := rg.Group("/auth")
	{
		authGroup.POST("/register", h.Register)
		authGroup.POST("/login", h.Login)
		authGroup.GET("/me", auth.JWTMiddleware(), h.Me)
	}

	users := rg.Group("/users", auth.JWTMiddleware(), auth.RequireRole("admin"))
	{
		users.GET("", h.List)
		users.GET("/:id", h.GetByID)
		users.PUT("/:id", h.Update)
		users.DELETE("/:id", h.Delete)
	}
}
