package order

import (
	"github.com/fixia/golang-api/internal/auth"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(rg *gin.RouterGroup, h *Handler) {
	orders := rg.Group("/orders", auth.JWTMiddleware())
	{
		orders.POST("", h.Create)
		orders.GET("", h.ListMine)
		orders.GET("/:id", h.GetByID)
		orders.PUT("/:id/status", auth.RequireRole("admin"), h.UpdateStatus)
	}
}
