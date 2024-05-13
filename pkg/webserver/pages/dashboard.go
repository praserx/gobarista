package pages

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const DashboardTitle = "Dashboard | GoBarista"

func DashboardGET(c *gin.Context) {
	c.HTML(http.StatusOK, "dashboard.go.tmpl", gin.H{
		"title": DashboardTitle,
	})
}
