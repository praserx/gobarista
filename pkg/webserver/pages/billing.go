package pages

import (
	"fmt"
	"net/http"
	"sort"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/praserx/gobarista/pkg/database"
	"github.com/praserx/gobarista/pkg/logger"
	"github.com/praserx/gobarista/pkg/webserver/security"
)

const BillingTitle = "Billing | GoBarista"

func BillingGET(c *gin.Context) {
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

	sort.SliceStable(bills, func(i, j int) bool {
		return bills[i].ID > bills[j].ID
	})

	c.HTML(http.StatusOK, "billing.go.tmpl", gin.H{
		"title": BillingTitle,
		"bills": bills,
	})
}
