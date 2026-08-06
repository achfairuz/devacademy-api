package models

type Role string

const (
	RoleAdmin   Role = "admin"
	RoleMentor  Role = "mentor"
	RoleStudent Role = "student"
)

func (r Role) IsValid() bool {
	switch r {
	case RoleAdmin, RoleMentor, RoleStudent:
		return true
	}
	return false
}

func (r Role) String() string {
	return string(r)
}
