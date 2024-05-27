package pages

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/praserx/gobarista/pkg/database"
	"github.com/praserx/gobarista/pkg/logger"
	"github.com/praserx/gobarista/pkg/webserver/security"
)

const RootTitle = "Main | GoBarista"

func RootGET(c *gin.Context) {
	c.HTML(http.StatusOK, "root.go.tmpl", gin.H{
		"title": "home",
	})
}

func RootPOST(c *gin.Context) {
	var ok bool
	var email string

	if email, ok = c.GetPostForm("email"); !ok {
		logger.Warning("request with malformed data input")
		c.HTML(http.StatusBadRequest, "root.go.tmpl", gin.H{
			"title":        AdminUsersAddTitle,
			"message":      "The data you sent are not correct!",
			"message_type": "warning",
		})
		return
	}

	if user, err := database.SelectUserByEmail(email); err == nil {
		sessionID := sessions.Default(c).Get(security.SessionKey).(string)
		userSession := security.Session{
			UserID:      user.ID,
			UserRole:    "",
			Code:        security.Code(),
			CodeUsed:    false,
			CodeValidTo: time.Now().Add(5 * time.Minute),
			Logged:      false,
		}

		security.SessionSet(sessionID, userSession)
		logger.Info(fmt.Sprintf("code: %s", userSession.Code))
		// TODO: Send verification code here via email
	} else {
		logger.Warning("login attempt for nonexisting user")
		c.HTML(http.StatusUnauthorized, "root.go.tmpl", gin.H{
			"title":        AdminUsersAddTitle,
			"message":      "We do not recognize you! Is this your account?",
			"message_type": "warning",
		})
		return
	}

	c.Redirect(http.StatusSeeOther, "/code-verification")
}
