package types

import (
	"github.com/golang-jwt/jwt"
)

type JWTClaims struct {
	UserID       int    `json:"user_id"`
	Role         string `json:"role"`
	DepartmentID *int   `json:"department_id,omitempty"`
	jwt.StandardClaims
}

func (c JWTClaims) IsAdmin() bool {
	return c.Role == RoleAdmin
}

func (c JWTClaims) IsDepartmentMember() bool {
	return c.DepartmentID != nil && c.Role == RoleMember
}
