package worker

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/adityawiryaa/api/pkg/response"
)

func (h *Handler) ExecuteHit(c *gin.Context) {
	resp, err := h.commandUC.EnqueueHit(c.Request.Context())
	if err != nil {
		log.Printf("[handler] enqueue hit failed: %v", err)
		response.Error(c, http.StatusInternalServerError, "ENQUEUE_FAILED", err.Error())
		return
	}

	response.Success(c, http.StatusAccepted, resp)
}

func (h *Handler) GetHitResult(c *gin.Context) {
	taskID := c.Param("taskId")
	if taskID == "" {
		response.Error(c, http.StatusBadRequest, "INVALID_TASK_ID", "task ID is required")
		return
	}

	log.Printf("[handler] getting result for task: %s", taskID)

	resp, err := h.queryUC.GetHitResult(c.Request.Context(), taskID)
	if err != nil {
		log.Printf("[handler] get hit result failed: task=%s error=%v", taskID, err)
		response.Error(c, http.StatusInternalServerError, "RESULT_FETCH_FAILED", err.Error())
		return
	}

	log.Printf("[handler] result found: task=%s status=%s", taskID, resp.Status)
	response.Success(c, http.StatusOK, resp)
}
