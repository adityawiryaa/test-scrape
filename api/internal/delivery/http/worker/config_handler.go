package worker

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/adityawiryaa/api/domain/entity"
	"github.com/adityawiryaa/api/pkg/response"
)

func (h *Handler) ReceiveConfig(c *gin.Context) {
	var cfg entity.Config
	if err := c.ShouldBindJSON(&cfg); err != nil {
		response.Error(c, http.StatusBadRequest, "INVALID_REQUEST", err.Error())
		return
	}

	h.commandUC.ReceiveConfig(&cfg)
	response.Success(c, http.StatusOK, map[string]any{"version": cfg.Version})
}

func (h *Handler) GetCurrentConfig(c *gin.Context) {
	cfg := h.queryUC.CurrentConfig()
	if cfg == nil {
		response.Error(c, http.StatusNotFound, "NOT_FOUND", "no config loaded")
		return
	}
	response.Success(c, http.StatusOK, cfg)
}
