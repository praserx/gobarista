package pages

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const CodeVerificationTitle = "Code Verification | GoBarista"

func CodeVerificationGET(c *gin.Context) {
	c.HTML(http.StatusOK, "code_verification.go.tmpl", gin.H{
		"title": CodeVerificationTitle,
	})
}
