package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/adityawiryaa/api/domain/request"
	"github.com/adityawiryaa/api/pkg/response"
)

func (h *Handler) RegisterAgent(c *gin.Context) {
	var req request.RegisterAgentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "INVALID_REQUEST", err.Error())
		return
	}

	resp, err := h.commandUC.RegisterAgent(c.Request.Context(), &req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "REGISTRATION_FAILED", err.Error())
		return
	}

	response.Success(c, http.StatusCreated, resp)
}
