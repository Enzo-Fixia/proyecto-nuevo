package user

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

func (h *Handler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	u, err := h.svc.Register(req)
	if err != nil {
		if errors.Is(err, utils.ErrDuplicateEmail) {
			utils.ErrorResponse(c, http.StatusConflict, err.Error())
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.SuccessResponse(c, http.StatusCreated, "user registered", u)
}

func (h *Handler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	resp, err := h.svc.Login(req)
	if err != nil {
		if errors.Is(err, utils.ErrInvalidCredentials) || errors.Is(err, utils.ErrUnauthorized) {
			utils.ErrorResponse(c, http.StatusUnauthorized, err.Error())
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "login successful", resp)
}

func (h *Handler) Me(c *gin.Context) {
	uid, _ := c.Get(auth.CtxUserID)
	u, err := h.svc.GetByID(uid.(uint))
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "user not found")
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "", u)
}

func (h *Handler) List(c *gin.Context) {
	users, err := h.svc.ListAll()
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "", users)
}

func (h *Handler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "invalid id")
		return
	}
	u, err := h.svc.GetByID(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.ErrorResponse(c, http.StatusNotFound, "user not found")
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "", u)
}

func (h *Handler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "invalid id")
		return
	}
	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	u, err := h.svc.Update(uint(id), req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "user updated", u)
}

func (h *Handler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "invalid id")
		return
	}
	if err := h.svc.Delete(uint(id)); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "user deleted", nil)
}
