package pages

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/praserx/gobarista/pkg/database"
	"github.com/praserx/gobarista/pkg/logger"
	"github.com/praserx/gobarista/pkg/models"
)

const AdminUsersAddTitle = "Add user | GoBarista"

func AdminUsersAddGET(c *gin.Context) {
	c.HTML(http.StatusOK, "admin_users_add.go.tmpl", gin.H{
		"title": "Add user",
	})
}

func AdminUsersAddPOST(c *gin.Context) {
	var ok bool
	var err error
	var eid, email, firstname, lastname, location string

	if eid, ok = c.GetPostForm("eid"); !ok {
		c.Status(http.StatusBadRequest)
		return
	}

	if email, ok = c.GetPostForm("email"); !ok {
		c.Status(http.StatusBadRequest)
		return
	}

	if firstname, ok = c.GetPostForm("firstname"); !ok {
		c.Status(http.StatusBadRequest)
		return
	}

	if lastname, ok = c.GetPostForm("lastname"); !ok {
		c.Status(http.StatusBadRequest)
		return
	}

	if location, ok = c.GetPostForm("location"); !ok {
		c.Status(http.StatusBadRequest)
		return
	}

	user := models.User{
		EID:       eid,
		Firstname: firstname,
		Lastname:  lastname,
		Email:     email,
		Location:  location,
	}

	if _, err = database.SelectUserByEID(user.EID); err == nil {
		logger.Warning("user already exists")
		c.HTML(http.StatusBadRequest, "admin_users_add.go.tmpl", gin.H{
			"title":        AdminUsersAddTitle,
			"message":      "User already exists!",
			"message_type": "danger",
		})
		return
	}

	id, err := database.InsertUser(user)
	if err != nil {
		logger.Info(fmt.Sprintf("error: cannot create user: %v", err.Error()))
		c.HTML(http.StatusInternalServerError, "admin_users_add.go.tmpl", gin.H{
			"title":        AdminUsersAddTitle,
			"message":      "Ooops! Something goes wrong.",
			"message_type": "danger",
		})
		return
	}

	logger.Info(fmt.Sprintf("new user successfully created: new user id: %d", id))
	c.HTML(http.StatusOK, "admin_users_add.go.tmpl", gin.H{
		"title":        AdminUsersAddTitle,
		"message":      "User was successfuly added!",
		"message_type": "success",
	})
}
