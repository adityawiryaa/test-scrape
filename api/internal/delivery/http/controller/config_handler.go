package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/adityawiryaa/api/domain/request"
	"github.com/adityawiryaa/api/pkg/response"
)

func (h *Handler) UpdateConfig(c *gin.Context) {
	var req request.UpdateConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "INVALID_REQUEST", err.Error())
		return
	}

	cfg, err := h.commandUC.UpdateConfig(c.Request.Context(), &req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "UPDATE_FAILED", err.Error())
		return
	}

	c.Header("ETag", strconv.FormatInt(cfg.Version, 10))
	response.Success(c, http.StatusCreated, cfg)
}

func (h *Handler) GetConfig(c *gin.Context) {
	ifNoneMatch := c.GetHeader("If-None-Match")

	cfg, err := h.queryUC.GetLatestConfig(c.Request.Context())
	if err != nil {
		response.Error(c, http.StatusNotFound, "NOT_FOUND", "no config available")
		return
	}

	etag := strconv.FormatInt(cfg.Version, 10)
	if ifNoneMatch == etag {
		c.Status(http.StatusNotModified)
		return
	}

	c.Header("ETag", etag)
	response.Success(c, http.StatusOK, cfg)
}

func (h *Handler) GetConfigByVersion(c *gin.Context) {
	versionStr := c.Param("version")
	version, err := strconv.ParseInt(versionStr, 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "INVALID_VERSION", "version must be a number")
		return
	}

	cfg, err := h.queryUC.GetConfigByVersion(c.Request.Context(), version)
	if err != nil {
		response.Error(c, http.StatusNotFound, "NOT_FOUND", "config version not found")
		return
	}

	response.Success(c, http.StatusOK, cfg)
}
