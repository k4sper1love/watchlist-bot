package roles

// Role represents a user role in the application.
// Roles are ordered hierarchically, with higher values indicating greater privileges.
type Role int

// Predefined roles in ascending order of privilege.
const (
	User       Role = iota // Basic user with minimal permissions.
	Helper                 // User with helper-level permissions.
	Admin                  // Administrator with elevated permissions.
	SuperAdmin             // Super administrator with extended permissions.
	Root                   // Root user with full system control.
)

// roleNames maps each Role to its corresponding human-readable name.
var roleNames = map[Role]string{
	User:       "user",
	Helper:     "helper",
	Admin:      "admin",
	SuperAdmin: "superAdmin",
	Root:       "root",
}

// String returns the human-readable name of the role.
// If the role is unknown, it returns "unknown".
func (r Role) String() string {
	if name, exists := roleNames[r]; exists {
		return name
	}
	return "unknown"
}

// HasAccess checks if the current role has sufficient access for the required role.
// Returns true if the current role is equal to or higher than the required role.
func (r Role) HasAccess(required Role) bool {
	return r >= required && r <= Root
}

// Compare compares the current role with another role.
// Returns:
// - 1 if the current role is higher than the other role.
// - -1 if the current role is lower than the other role.
// - 0 if the roles are equal.
func (r Role) Compare(other Role) int {
	switch {
	case r > other:
		return 1
	case r < other:
		return -1
	default:
		return 0
	}
}

// NextRole returns the next role in the hierarchy.
// If the current role is already the highest (Root), it returns Root.
func (r Role) NextRole() Role {
	if r < Root {
		return r + 1
	}
	return Root
}

// PrevRole returns the previous role in the hierarchy.
// If the current role is already the lowest (User), it returns User.
func (r Role) PrevRole() Role {
	if r > User {
		return r - 1
	}
	return User
}
