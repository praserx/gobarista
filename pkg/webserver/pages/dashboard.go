package pages

import (
	"fmt"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/praserx/gobarista/pkg/database"
	"github.com/praserx/gobarista/pkg/logger"
	"github.com/praserx/gobarista/pkg/webserver/security"
)

const DashboardTitle = "Dashboard | GoBarista"

func DashboardGET(c *gin.Context) {
	sessionID := sessions.Default(c).Get(security.SessionKey).(string)
	sessionData, ok := security.SessionGet(sessionID)
	if !ok {
		logger.Warning("cannot get session data")
		c.HTML(http.StatusInternalServerError, "root.go.tmpl", gin.H{
			"title":        BillingTitle,
			"message":      "Ooops! Something goes wrong.",
			"message_type": "danger",
		})
	}

	bills, err := database.SelectAllBillsForUser(sessionData.UserID)
	if err != nil {
		logger.Error(fmt.Sprintf("error: cannot get users: %v", err.Error()))
		c.HTML(http.StatusInternalServerError, "billing.go.tmpl", gin.H{
			"title": BillingTitle,
			"bills": "",
		})
	}

	c.HTML(http.StatusOK, "dashboard.go.tmpl", gin.H{
		"title": DashboardTitle,
		"bills": bills,
	})
}
