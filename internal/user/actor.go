package user

// Actor is the logged-in user performing an action. The service uses it for
// rules like "you can only update yourself unless you're an admin".
type Actor struct {
	ID   uint
	Role Role
}

func (a Actor) IsAdmin() bool { return a.Role == RoleAdmin }
