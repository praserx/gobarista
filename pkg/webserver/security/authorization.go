package authorization

import (
	"github.com/praserx/gobarista/pkg/webserver/routes"
)

const RoleAdmin = "admin"
const RoleUser = "user"
const RoleGuest = "guest"

// C - R - U - D

var Permissions = map[string]map[string][]rune{
	RoleAdmin: {
		routes.APIv1Users:      {'C', 'R', 'U', 'D'},
		routes.APIv1UsersID:    {'C', 'R', 'U', 'D'},
		routes.APIv1Bills:      {'C', 'R', 'U', 'D'},
		routes.APIv1BillsID:    {'C', 'R', 'U', 'D'},
		routes.APIv1Accounts:   {'C', 'R', 'U', 'D'},
		routes.APIv1AccountsID: {'C', 'R', 'U', 'D'},
	},
	RoleUser: {
		routes.APIv1Users:      {},
		routes.APIv1UsersID:    {'R'},
		routes.APIv1Bills:      {},
		routes.APIv1BillsID:    {'R'},
		routes.APIv1Accounts:   {'R'},
		routes.APIv1AccountsID: {'R'},
	},
	RoleGuest: {
		routes.APIv1Users:      {},
		routes.APIv1UsersID:    {},
		routes.APIv1Bills:      {},
		routes.APIv1BillsID:    {},
		routes.APIv1Accounts:   {},
		routes.APIv1AccountsID: {},
	},
}

var ExtendedPermission = []string{RoleAdmin}

func IsAuthorized(route, role string, perm rune, reqExt bool) bool {
	if p, ok := Permissions[role]; ok {
		if r, ok := p[route]; ok {
			for _, op := range r {
				if perm == op && (!reqExt || HasExtendedPermission(role)) {
					return true
				}
			}
		}
	}

	return false
}

func HasExtendedPermission(role string) bool {
	for _, r := range ExtendedPermission {
		if r == role {
			return true
		}
	}
	return false
}

func GetPermissionFromRequestMethod(method string) rune {
	switch method {
	case "POST":
		return 'C'
	case "GET":
		return 'R'
	case "PUT":
		return 'U'
	case "DELETE":
		return 'D'
	default:
		return '-'
	}
}
