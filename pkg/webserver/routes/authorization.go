package routes

const RoleAdmin = "admin"
const RoleUser = "user"
const RoleGuest = "guest"

// C - R - U - D

var Permissions = map[string]map[string][]rune{
	RoleAdmin: {
		PageRoot:             {'C', 'R', 'U', 'D'},
		PageCodeVerification: {'C', 'R', 'U', 'D'},
		PageDashboard:        {'C', 'R', 'U', 'D'},
		PageBilling:          {'C', 'R', 'U', 'D'},
		PageAdminBilling:     {'C', 'R', 'U', 'D'},
		PageAdminUsers:       {'C', 'R', 'U', 'D'},
		PageAdminUsersAdd:    {'C', 'R', 'U', 'D'},
	},
	RoleUser: {
		PageRoot:             {'C', 'R'},
		PageCodeVerification: {'C', 'R'},
		PageDashboard:        {'R'},
		PageBilling:          {'R'},
		PageAdminBilling:     {},
		PageAdminUsers:       {},
		PageAdminUsersAdd:    {},
	},
	RoleGuest: {
		PageRoot:             {'C', 'R'},
		PageCodeVerification: {'C', 'R'},
		PageDashboard:        {},
		PageBilling:          {},
		PageAdminBilling:     {},
		PageAdminUsers:       {},
		PageAdminUsersAdd:    {},
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
