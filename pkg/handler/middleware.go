package handler

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"math/rand"
	"net/http"
	"net/smtp"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	authorizationHeader = "Authorization"
	userCtx             = "userId"
)

func generatePasswordHash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func comparePasswordHash(hash, inputPassword string) error {
	hashedPassword := []byte(hash)
	inputPasswordBytes := []byte(inputPassword)

	return bcrypt.CompareHashAndPassword(hashedPassword, inputPasswordBytes)
}

func validateEmail(email string) bool {
	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(emailRegex)

	return re.MatchString(email)
}

func generateRecoveryCode() (string, string) {
	rand.New(rand.NewSource(time.Now().UnixNano()))

	code := rand.Intn(900000) + 100000
	codeString := strconv.Itoa(code)

	currentTime := time.Now().UTC()
	expiredAt := currentTime.Add(5 * time.Minute).UTC()
	formattedExpiresAt := expiredAt.Format("2006-01-02 15:04:05")

	return codeString, formattedExpiresAt
}

func sendCode(toEmail string) (string, string, error) {
	code, expiredAt := generateRecoveryCode()

	message := fmt.Sprintf("Subject: Recovery code\r\n\r\nYour recovery code: %s. This code expired at: %s", code, expiredAt)

	to := []string{toEmail}
	from := os.Getenv("SMTP_EMAIL_FROM")
	password := os.Getenv("SMTP_PASSWORD")
	host := os.Getenv("SMTP_HOST")
	url := os.Getenv("SMTP_HOST") + ":" + os.Getenv("SMTP_PORT")
	auth := smtp.PlainAuth("", from, password, host)

	err := smtp.SendMail(
		url,
		auth,
		from,
		to,
		[]byte(message),
	)
	if err != nil {
		return "", "", errors.New("error with sending recovery mail")
	}
	return code, expiredAt, nil
}

func (h *Handler) userIdentity(c *gin.Context) {
	header := c.GetHeader(authorizationHeader)
	if header == "" {
		newErrorResponse(c, http.StatusUnauthorized, "empty auth header")
		return
	}

	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 {
		newErrorResponse(c, http.StatusUnauthorized, "invalid auth header")
		return
	}

	//parse token
	userId, err := h.services.Authorization.ParseToken(headerParts[1])
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}
	c.Set(userCtx, userId)
}
