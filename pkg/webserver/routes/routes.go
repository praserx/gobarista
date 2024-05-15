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
	PageRoot             = "page_root"
	PageCodeVerification = "page_code_verification"
	PageDashboard        = "page_dashboard"
	PageBilling          = "page_billing"
	PageTransactions     = "page_transactions"
	PageAdminBilling     = "page_admin_billing"
	PageAdminUsers       = "page_admin_users"
	PageAdminUsersAdd    = "page_admin_users_add"
)

var Routes = map[string]Route{
	PageRoot:             {Title: pages.RootTitle, Path: "/"},
	PageCodeVerification: {Title: pages.CodeVerificationTitle, Path: "/code-verification"},
	PageDashboard:        {Title: pages.DashboardTitle, Path: "/dashboard"},
	PageBilling:          {Title: pages.BillingTitle, Path: "/billing"},
	PageTransactions:     {Title: pages.TransactionsTitle, Path: "/transactions"},
	PageAdminBilling:     {Title: pages.AdminBillingTitle, Path: "/admin/billing"},
	PageAdminUsers:       {Title: pages.AdminUsersTitle, Path: "/admin/users"},
	PageAdminUsersAdd:    {Title: pages.AdminUsersAddTitle, Path: "/admin/users/add"},
}

func SetupRoutes(router *gin.Engine) {
	router.GET(Routes[PageRoot].Path, pages.RootGET)
	router.POST(Routes[PageRoot].Path, pages.RootPOST)
	router.GET(Routes[PageCodeVerification].Path, pages.CodeVerificationGET)
	router.POST(Routes[PageCodeVerification].Path, pages.CodeVerificationPOST)
	router.GET(Routes[PageDashboard].Path, pages.DashboardGET)
	router.GET(Routes[PageBilling].Path, pages.BillingGET)
	router.GET(Routes[PageTransactions].Path, pages.TransactionsGET)
	router.GET(Routes[PageAdminBilling].Path, pages.AdminBillingGET)
	router.GET(Routes[PageAdminUsers].Path, pages.AdminUsersGET)
	router.GET(Routes[PageAdminUsersAdd].Path, pages.AdminUsersAddGET)
	router.POST(Routes[PageAdminUsersAdd].Path, pages.AdminUsersAddPOST)
}
