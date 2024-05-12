package pages

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func BillingGET(c *gin.Context) {
	c.HTML(http.StatusOK, "billing.go.tmpl", gin.H{
		"title":   "Vyúčtování",
		"billing": "",
	})
}
