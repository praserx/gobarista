package pages

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func MainGET(c *gin.Context) {
	fmt.Println("index")
	c.HTML(http.StatusOK, "index.go.tmpl", gin.H{
		"title": "home",
	})
}
