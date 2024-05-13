package pages

import (
	"net/http"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/praserx/gobarista/pkg/logger"
	"github.com/praserx/gobarista/pkg/webserver/security"
)

const CodeVerificationTitle = "Code Verification | GoBarista"

func CodeVerificationGET(c *gin.Context) {
	c.HTML(http.StatusOK, "code_verification.go.tmpl", gin.H{
		"title": CodeVerificationTitle,
	})
}

func CodeVerificationPOST(c *gin.Context) {
	var ok bool
	var code string

	if code, ok = c.GetPostForm("code"); !ok {
		logger.Warning("request with malformed data input")
		c.HTML(http.StatusBadRequest, "root.go.tmpl", gin.H{
			"title":        CodeVerificationTitle,
			"message":      "The data you sent are not correct!",
			"message_type": "warning",
		})
		return
	}

	sessionID := sessions.Default(c).Get(security.SessionKey).(string)
	sessionData, ok := security.SessionGet(sessionID)
	if !ok {
		logger.Warning("cannot get session data")
		c.HTML(http.StatusInternalServerError, "root.go.tmpl", gin.H{
			"title":        CodeVerificationTitle,
			"message":      "Ooops! Something goes wrong.",
			"message_type": "danger",
		})
	}

	if time.Now().After(sessionData.CodeValidTo) {
		logger.Warning("security code is stale")
		c.HTML(http.StatusUnauthorized, "code_verification.go.tmpl", gin.H{
			"title":        CodeVerificationTitle,
			"message":      "You shall not pass! Your security code is stale.",
			"message_type": "warning",
		})
		return
	}

	if sessionData.CodeUsed {
		logger.Warning("security code was already used by someone")
		c.HTML(http.StatusUnauthorized, "code_verification.go.tmpl", gin.H{
			"title":        CodeVerificationTitle,
			"message":      "You shall not pass! Your security one-time code has been used already.",
			"message_type": "warning",
		})
		return
	}

	if sessionData.CodeUsed || sessionData.Code != code {
		logger.Warning("security code not match")
		c.HTML(http.StatusUnauthorized, "code_verification.go.tmpl", gin.H{
			"title":        CodeVerificationTitle,
			"message":      "You shall not pass! Your security code is invalid.",
			"message_type": "warning",
		})
		return
	}

	sessionData.Code = ""
	sessionData.CodeUsed = true
	security.SessionSet(sessionID, sessionData)

	c.Redirect(http.StatusSeeOther, "/dashboard")
}
