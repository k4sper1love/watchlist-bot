package roles

type Role int

const (
	User Role = iota
	Helper
	Admin
	SuperAdmin
	Root
)

var roleNames = map[Role]string{
	User:       "user",
	Helper:     "helper",
	Admin:      "admin",
	SuperAdmin: "superAdmin",
	Root:       "root",
}

func (r Role) String() string {
	if name, exists := roleNames[r]; exists {
		return name
	}
	return "unknown"
}

func (r Role) HasAccess(required Role) bool {
	return r >= required && r <= Root
}

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

func (r Role) NextRole() Role {
	if r < Root {
		return r + 1
	}
	return Root
}

func (r Role) PrevRole() Role {
	if r > User {
		return r - 1
	}
	return User
}
