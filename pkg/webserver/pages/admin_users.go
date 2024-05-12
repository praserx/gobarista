package pages

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/praserx/gobarista/pkg/database"
	"github.com/praserx/gobarista/pkg/logger"
)

func AdminUsersGET(c *gin.Context) {
	users, err := database.SelectAllUsers()
	if err != nil {
		logger.Error(fmt.Sprintf("error: cannot get users: %v", err.Error()))
		c.HTML(http.StatusInternalServerError, "admin_users.go.tmpl", gin.H{
			"title": "Uživatelské účty",
			"users": "",
		})
	}

	c.HTML(http.StatusOK, "admin_users.go.tmpl", gin.H{
		"title": "Uživatelské účty",
		"users": users,
	})
}
