package models

// userValidator is a chained type that performs validation and
// normalization of data before being passed to the final UserDB implementation
type userValidator struct {
	UserDB
}

// Query methods

func (uv *userValidator) ByID(id uint) (*User, error) {
	// TODO: validate ByID data
	return uv.UserDB.ByID(id)
}

func (uv *userValidator) ByEmail(email string) (*User, error) {
	// TODO: validate by email
	return uv.UserDB.ByEmail(email)
}

func (uv *userValidator) ByRemember(token string) (*User, error) {
	// TODO: validate by remember token
	return uv.UserDB.ByRemember(token)
}

// Data Alteration Methods

func (uv *userValidator) Create(user *User) error {
	// TODO: validation
	return uv.UserDB.Create(user)
}

func (uv *userValidator) Update(user *User) error {
	// TODO: validation
	return uv.UserDB.Update(user)
}

func (uv *userValidator) Delete(id uint) error {
	// TODO: validation
	return uv.UserDB.Delete(id)
}
