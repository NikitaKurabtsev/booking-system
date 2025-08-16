package handlers

import "github.com/gin-gonic/gin"

type errorResponse struct {
	Error string `json:"error"`
}

type statusResponse struct {
	Message string `json:"message"`
}

// TODO: add logger interface instead of slog injection
func (h *Handler) errorResponse(c *gin.Context, statusCode int, message string) {
	h.logger.Error("error", message)
	c.AbortWithStatusJSON(statusCode, errorResponse{Error: message})
}
