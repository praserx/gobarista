package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/praserx/gobarista/pkg/webserver/pages"
)

type Route struct {
	Title string
	Path  string
}

var (
	PageIndex       = "page_index"
	PageLoginCode   = "page_login_code"
	PageDashboard   = "page_dashboard"
	APIv1Users      = "api_v1_users"
	APIv1UsersID    = "api_v1_users_id"
	APIv1Bills      = "api_v1_bills"
	APIv1BillsID    = "api_v1_bills_id"
	APIv1Accounts   = "api_v1_accounts"
	APIv1AccountsID = "api_v1_accounts_id"
)

func Initialize(router *gin.Engine) {
	router.GET("/", pages.MainGET)
	router.GET("/dashboard", pages.DashboardGET)
	router.GET("/billing", pages.BillingGET)
	router.GET("/admin/billing", pages.AdminBillingGET)
	router.GET("/admin/users", pages.AdminUsersGET)
	router.GET("/admin/users/add", pages.AdminUsersAddGET)
	router.POST("/admin/users/add", pages.AdminUsersAddPOST)
}

var Routes = map[string]Route{
	PageIndex:       {Title: "Vítejte | GoBarista", Path: "/"},
	PageLoginCode:   {Title: "Přihlásit se | GoBarista", Path: "/login"},
	PageDashboard:   {Title: "Přehled | GoBarista", Path: "/dashboard"},
	APIv1Users:      {Title: "Users", Path: "/api/v1/users"},
	APIv1UsersID:    {Title: "UsersID", Path: "/api/v1/users/:id"},
	APIv1Bills:      {Title: "Bills", Path: "/api/v1/bills"},
	APIv1BillsID:    {Title: "BillsID", Path: "/api/v1/bills/:id"},
	APIv1Accounts:   {Title: "Accounts", Path: "/api/v1/accounts"},
	APIv1AccountsID: {Title: "AccountsID", Path: "/api/v1/accounts/:id"},
}
