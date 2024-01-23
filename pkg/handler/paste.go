package handler

import (
	pasteShare "PasteShare"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) createPaste(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		return
	}

	var input pasteShare.Paste
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	input.Date = getCurrentTime()
	id, err := h.services.Paste.CreatePaste(userId, input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

type getAllPasteResponse struct {
	Data []pasteShare.Paste `json:"user pastes"`
}

func (h *Handler) getAllPastes(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		return
	}

	pastes, err := h.services.Paste.GetAll(userId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, getAllPasteResponse{
		Data: pastes,
	})
}
