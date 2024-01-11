package handler

import (
	pasteShare "PasteShare"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type signUpInput struct {
	NickName        string `json:"nickname" binding:"required"`
	Email           string `json:"email" binding:"required"`
	Password        string `json:"password" binding:"required"`
	PasswordConfirm string `json:"passwordConfirm" binding:"required"`
}

func (h *Handler) signUp(c *gin.Context) {
	var input signUpInput

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if !validateEmail(input.Email) {
		newErrorResponse(c, http.StatusBadRequest, "enter a different email")
		return
	}

	if input.Password != input.PasswordConfirm {
		newErrorResponse(c, http.StatusBadRequest, "passwords does not match")
		return
	}

	_, err := h.services.Authorization.CheckNickNameAndEmail(input.NickName, input.Email)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	user := pasteShare.User{
		NickName: input.NickName,
		Email:    input.Email,
		Password: input.Password,
	}

	user.Password, err = generatePasswordHash(user.Password)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	id, err := h.services.Authorization.CreateUser(user)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

type signInInput struct {
	NickName string `json:"nickname" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (h *Handler) signIn(c *gin.Context) {
	var input signInInput

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	passwordHash, err := h.services.GetPasswordHash(input.NickName)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err = comparePasswordHash(passwordHash, input.Password); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	token, err := h.services.Authorization.GenerateToken(input.NickName, passwordHash)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"token": token,
	})
}

type forgotPasswordInput struct {
	Email string `json:"email" binding:"required"`
}

func (h *Handler) forgotPassword(c *gin.Context) {
	var input forgotPasswordInput

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	_, err := h.services.Authorization.CheckEmail(input.Email)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	code, expiredAt, err := sendCode(input.Email)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	err = h.services.Authorization.SendRecoveryMail(input.Email, code, expiredAt)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"status": "ok",
	})
}

type resetPasswordInput struct {
	RecoveryCode       string `json:"recovery_code" binding:"required"`
	NewPassword        string `json:"new_password" binding:"required"`
	NewPasswordConfirm string `json:"new_password_confirm" binding:"required"`
}

func (h *Handler) resetPassword(c *gin.Context) {
	var input resetPasswordInput

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	id, expiredAt, err := h.services.Authorization.CheckCode(input.RecoveryCode)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	currentTime := time.Now().UTC()
	formattedExpiredAt, err := time.Parse("2006-01-02T15:04:05Z07:00", expiredAt)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "fail to parse expired at code time")
		return
	}

	if currentTime.After(formattedExpiredAt) {
		newErrorResponse(c, http.StatusBadRequest, "code expired")
		return
	}

	if input.NewPassword != input.NewPasswordConfirm {
		newErrorResponse(c, http.StatusBadRequest, "passwords does not match")
		return
	}

	passwordHash, err := generatePasswordHash(input.NewPassword)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.services.SetNewPassword(id, passwordHash); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"status": "ok",
	})

}
