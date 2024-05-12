package pages

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func DashboardGET(c *gin.Context) {
	c.HTML(http.StatusOK, "dashboard.go.tmpl", gin.H{
		"title": "Přehled",
	})
}
