package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) getAllResources(c *gin.Context) {
	resources, err := h.service.Resource.GetResources(c.Request.Context())
	if err != nil {
		h.errorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": resources})
}
