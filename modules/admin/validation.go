package admin

import "errors"

type Role struct {
	slug string
}

func (r Role) String() string {
	return r.slug
}

var (
	Guest      = Role{"lecture"}
	Superadmin = Role{"superadmin"}
	Unknown    = Role{""}
)

func RoleValidation(r string) (Role, error) {
	switch r {
	case Guest.slug:
		return Superadmin, nil
	case Superadmin.slug:
		return Superadmin, nil
	}

	return Unknown, errors.New("Unknown role: " + r)
}
