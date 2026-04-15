package order

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/fixia/golang-api/internal/auth"
	"github.com/fixia/golang-api/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Handler struct {
	svc Service
}

func NewHandler(svc Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) Create(c *gin.Context) {
	var req CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	uid, _ := c.Get(auth.CtxUserID)
	o, err := h.svc.Create(req, uid.(uint))
	if err != nil {
		if errors.Is(err, utils.ErrInsufficientStock) {
			utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
			return
		}
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.ErrorResponse(c, http.StatusNotFound, "product not found")
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.SuccessResponse(c, http.StatusCreated, "order created", o)
}

func (h *Handler) ListMine(c *gin.Context) {
	uid, _ := c.Get(auth.CtxUserID)
	orders, err := h.svc.ListByUser(uid.(uint))
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "", orders)
}

func (h *Handler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "invalid id")
		return
	}
	o, err := h.svc.GetByID(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.ErrorResponse(c, http.StatusNotFound, "order not found")
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	uid, _ := c.Get(auth.CtxUserID)
	role, _ := c.Get(auth.CtxRole)
	if o.UserID != uid.(uint) && role != "admin" {
		utils.ErrorResponse(c, http.StatusForbidden, "forbidden")
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "", o)
}

func (h *Handler) UpdateStatus(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "invalid id")
		return
	}
	var req UpdateStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	o, err := h.svc.UpdateStatus(uint(id), req.Status)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "order status updated", o)
}
