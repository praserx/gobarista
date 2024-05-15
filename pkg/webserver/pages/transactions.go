package pages

import (
	"net/http"
	"sort"

	"github.com/gin-gonic/gin"
	"github.com/praserx/gobarista/pkg/database"
)

const TransactionsTitle = "Transactions | GoBarista"

func TransactionsGET(c *gin.Context) {
	transactions, _ := database.SelectAllTransactionsForUser(1)
	// if err != nil {
	// 	logger.Error(fmt.Sprintf("error: cannot get users: %v", err.Error()))
	// 	c.HTML(http.StatusInternalServerError, "admin_billing.go.tmpl", gin.H{
	// 		"title": "Uživatelské účty",
	// 		"users": "",
	// 	})
	// }

	sort.SliceStable(transactions, func(i, j int) bool {
		return transactions[i].ID > transactions[j].ID
	})

	c.HTML(http.StatusOK, "transactions.go.tmpl", gin.H{
		"title":        TransactionsTitle,
		"transactions": transactions,
	})
}
