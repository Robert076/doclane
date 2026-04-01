package types

const (
	RoleAdmin  = "admin"
	RoleMember = "member"
)

func IsValidRole(role string) bool {
	return role == RoleAdmin || role == RoleMember
}
