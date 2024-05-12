package pages

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/praserx/gobarista/pkg/database"
	"github.com/praserx/gobarista/pkg/logger"
)

func AdminBillingGET(c *gin.Context) {
	periods, err := database.SelectAllPeriods()
	if err != nil {
		logger.Error(fmt.Sprintf("error: cannot get users: %v", err.Error()))
		c.HTML(http.StatusInternalServerError, "admin_billing.go.tmpl", gin.H{
			"title": "Uživatelské účty",
			"users": "",
		})
	}

	c.HTML(http.StatusOK, "admin_billing.go.tmpl", gin.H{
		"title":   "Billing periods",
		"periods": periods,
	})
}
