package webserver

import (
	"fmt"
	"html/template"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/praserx/gobarista/pkg/webserver/routes"
	"github.com/praserx/gobarista/resources"
)

// This feature is unimplemented.
type Server struct {
	host string
	port int
}

func NewServer(host string, port int) *Server {
	var srv Server

	srv.host = host
	srv.port = port

	return &srv
}

func (s *Server) Run() {
	gin.SetMode(gin.DebugMode)
	router := gin.Default()

	// router.Use(authorization.AuthorizationMiddleware())

	htmlTemplates := template.New("default").Funcs(template.FuncMap{
		"formatAsDate":  formatAsDate,
		"formatAsMoney": formatAsMoney,
	})
	htmlTemplates = template.Must(htmlTemplates.ParseGlob("resources/templates/web/*.go.tmpl"))
	htmlTemplates = template.Must(htmlTemplates.ParseGlob("resources/templates/web/partials/*.go.tmpl"))

	router.SetHTMLTemplate(htmlTemplates)

	router.StaticFS("/static", http.FS(resources.DirAssets))

	routes.Initialize(router)

	router.Run(fmt.Sprintf("%s:%d", s.host, s.port))
}

func formatAsDate(t time.Time) string {
	year, month, day := t.Date()
	return fmt.Sprintf("%d-%02d-%02d", year, month, day)
}

func formatAsMoney(val float32) string {
	return fmt.Sprintf("%.2f", val)
}
