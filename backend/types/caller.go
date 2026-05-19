package types

// CallerContext holds the resolved identity of the authenticated caller.
// It is populated by the auth middleware after verifying the Cognito token
// and looking up the corresponding user in the application database.
type CallerContext struct {
	UserID       int    // Primary key in database
	CognitoSub   string // Cognito's stable user identifier (the "sub" claim)
	Role         string
	DepartmentID *int
}

func (c CallerContext) IsAdmin() bool {
	return c.Role == RoleAdmin
}

func (c CallerContext) IsDepartmentMember() bool {
	return c.Role == RoleMember && c.UserID != 0
}

// ContextKey is used to store and retrieve CallerContext from request contexts
// without colliding with keys from other packages.
type ContextKey string

const (
	CallerContextKey ContextKey = "caller_context"
)
