package pages

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/praserx/gobarista/pkg/database"
	"github.com/praserx/gobarista/pkg/logger"
	"github.com/praserx/gobarista/pkg/models"
)

func AdminUsersAddGET(c *gin.Context) {
	c.HTML(http.StatusOK, "admin_users_add.go.tmpl", gin.H{
		"title": "Add user",
	})
}

func AdminUsersAddPOST(c *gin.Context) {
	var ok bool
	var err error
	var eid, email, firstname, lastname, location string

	fmt.Println(c.Request.Form)

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

	if _, err = database.SelectUserByEID(user.EID); err != nil {
		id, err := database.InsertUser(user)
		if err != nil {
			logger.Info(fmt.Sprintf("error: cannot create user: %v", err.Error()))
			c.Status(http.StatusInternalServerError)
			return
		}
		logger.Info(fmt.Sprintf("user successfully created: new user id: %d", id))
	} else {
		logger.Info("user already exists")
		c.Status(http.StatusBadRequest)
		return
	}

	c.HTML(http.StatusOK, "admin_users_add.go.tmpl", gin.H{
		"title":        "Add user",
		"message":      "User was successfuly added!",
		"message_type": "success",
	})
}
