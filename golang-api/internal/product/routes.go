package product

import (
	"github.com/fixia/golang-api/internal/auth"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(rg *gin.RouterGroup, h *Handler) {
	products := rg.Group("/products")
	{
		products.GET("", h.List)
		products.GET("/:id", h.GetByID)

		admin := products.Group("", auth.JWTMiddleware(), auth.RequireRole("admin"))
		{
			admin.POST("", h.Create)
			admin.PUT("/:id", h.Update)
			admin.DELETE("/:id", h.Delete)
		}
	}
}
