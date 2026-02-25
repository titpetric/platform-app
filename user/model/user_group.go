package model

// NewUserGroup creates a new empty UserGroup.
func NewUserGroup() *UserGroup {
	return &UserGroup{}
}

// String returns the title of the UserGroup.
func (u *UserGroup) String() string {
	return u.Title
}
