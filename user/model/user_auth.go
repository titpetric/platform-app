package model

// Valid reports whether the UserAuth has an email and password.
func (u *UserAuth) Valid() bool {
	if u == nil || u.Email == "" || u.Password == "" {
		return false
	}
	return true
}
