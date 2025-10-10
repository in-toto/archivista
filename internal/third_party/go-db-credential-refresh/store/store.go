package store

// Credential implements the Credentials interface.
type Credential struct {
	Username string
	Password string
}

// GetUsername implements the Credentials interface.
func (c Credential) GetUsername() string {
	return c.Username
}

// GetPassword implements the Credentials interface.
func (c Credential) GetPassword() string {
	return c.Password
}
