package pages

import (
	"net/http"
	"sort"

	"github.com/gin-gonic/gin"
	"github.com/praserx/gobarista/pkg/database"
)

const BillingTitle = "Billing | GoBarista"

func BillingGET(c *gin.Context) {
	bills, _ := database.SelectAllBillsForUser(1)
	// if err != nil {
	// 	logger.Error(fmt.Sprintf("error: cannot get users: %v", err.Error()))
	// 	c.HTML(http.StatusInternalServerError, "admin_billing.go.tmpl", gin.H{
	// 		"title": "Uživatelské účty",
	// 		"users": "",
	// 	})
	// }

	sort.SliceStable(bills, func(i, j int) bool {
		return bills[i].ID > bills[j].ID
	})

	c.HTML(http.StatusOK, "billing.go.tmpl", gin.H{
		"title": "Billing",
		"bills": bills,
	})
}
