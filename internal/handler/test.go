package handler

import (
	"github.com/gin-gonic/gin"
	"prosto-delaj-api/internal/exceptions"
	"prosto-delaj-api/models"
)

func (h *Handler) create(c *gin.Context) {
	var input models.CreateInput

	if err := c.ShouldBindJSON(&input); err != nil {
		h.sendResponseSuccess(c, nil, exceptions.BadRequest)
		return
	}

	processStatus := h.usecase.Create(input)

	h.sendResponseSuccess(c, nil, processStatus)
}
