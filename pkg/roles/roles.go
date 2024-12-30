package roles

type Role int

const (
	User Role = iota
	Helper
	Admin
	SuperAdmin
	Root
)

func (r Role) String() string {
	switch r {
	case User:
		return "user"
	case Helper:
		return "helper"
	case Admin:
		return "admin"
	case SuperAdmin:
		return "superAdmin"
	case Root:
		return "root"
	default:
		return "unknown"
	}
}

func (r Role) HasAccess(required Role) bool {
	return r >= required
}

func (r Role) Higher(other Role) bool {
	return r > other
}

func (r Role) Equal(other Role) bool {
	return r == other
}

func (r Role) NextRole() Role {
	next := r + 1

	if next <= Root {
		return next
	}

	return r
}

func (r Role) PrevRole() Role {
	prev := r - 1

	if prev >= User {
		return prev
	}

	return r
}
