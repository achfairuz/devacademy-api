package models

type Role string

const (
	RoleAdmin      Role = "admin"
	RoleInstructor Role = "instructor"
	RoleStudent    Role = "student"
)

func (r Role) IsValid() bool {
	switch r {
	case RoleAdmin, RoleInstructor, RoleStudent:
		return true
	}
	return false
}

func (r Role) String() string {
	return string(r)
}
