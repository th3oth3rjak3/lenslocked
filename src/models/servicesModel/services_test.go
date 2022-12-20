package servicesModel

import (
	"strings"
	"testing"
	"time"

	"lenslocked/config"
	"lenslocked/models/errorsModel"
	"lenslocked/models/usersModel"
	"lenslocked/rand"
)

func mockServices(causeDbError bool) (*Services, error) {
	dbCfg := config.TestPostgresConfig()
	psqlInfo := dbCfg.ConnectionInfo()
	if causeDbError {
		psqlInfo = psqlInfo + "/////"
	}
	services, err := NewServices(
		WithGorm(dbCfg.Dialect(), psqlInfo),
		WithUser(config.DefaultHashKeyConfig()),
		WithGallery(),
		WithImages(),
		WithLogMode(false),
	)
	if err != nil {
		return nil, err
	}

	// Clear the users table between tests.
	services.DestructiveReset()
	return services, nil
}

func fakeUserService() usersModel.User {
	remember, err := rand.RememberToken()
	if err != nil {
		panic(err)
	}
	name := "Fake User"
	email := "fake.user@email.com"
	password := "some special password"

	return usersModel.User{
		Name:     name,
		Email:    email,
		Password: password,
		Remember: remember,
	}
}

func TestCreateUser(t *testing.T) {
	s, err := mockServices(false)
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()

	fakeUser := fakeUserService()
	user := fakeUserService()

	if err := s.User.Create(&user); err != nil {
		t.Fatal(err)
	}

	if user.ID < 1 {
		t.Errorf("User ID is less than 1. Have: %d, Want: %d", user.ID, 1)
	}
	if user.Name != fakeUser.Name {
		t.Errorf("User Name is incorrect. Have: %s, Want: %s", user.Name, fakeUser.Name)
	}
	if user.Email != fakeUser.Email {
		t.Errorf("User Email is incorrect. Have: %s, Want: %s", user.Email, fakeUser.Email)
	}
	if time.Since(user.CreatedAt) > time.Duration(10*time.Second) {
		expectedMin := time.Now()
		expectedMax := expectedMin.Add(10 * time.Second)
		// Because I guess the Go epoch is 01-02-2006 03:04:05 PM (15:04:05)
		fString := "15:04:05"
		minStr := expectedMin.Format(fString)
		maxStr := expectedMax.Format(fString)
		actual := user.CreatedAt.Format(fString)
		t.Errorf("Expected user to be created recently. Have: %v, Want: %v - %v", actual, minStr, maxStr)
	}
	if time.Since(user.UpdatedAt) > time.Duration(10*time.Second) {
		expectedMin := time.Now()
		expectedMax := expectedMin.Add(10 * time.Second)
		// Because I guess the Go epoch is 01-02-2006 03:04:05 PM (15:04:05)
		fString := "15:04:05"
		minStr := expectedMin.Format(fString)
		maxStr := expectedMax.Format(fString)
		actual := user.UpdatedAt.Format(fString)
		t.Errorf("Expected user to be updated recently. Have: %v, Want: %v - %v", actual, minStr, maxStr)
	}
}

func TestCreateDuplicateEmailUser(t *testing.T) {
	s, err := mockServices(false)
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()

	fakeUser := fakeUserService()
	user := fakeUserService()

	if err := s.User.Create(&user); err != nil {
		t.Fatal(err)
	}
	if err := s.User.Create(&fakeUser); err != errorsModel.ErrEmailTaken {
		t.Errorf("Expected ErrEmailTaken, Got: %s", err.Error())
	}
}

func TestCreateUserWithInvalidEmail(t *testing.T) {
	s, err := mockServices(false)
	if err != nil {
		t.Fatalf("Expected mockServices, %s", err.Error())
	}
	defer s.Close()
	user := fakeUserService()
	user.Email = ""
	err = s.User.Create(&user)
	if err != errorsModel.ErrEmailMissing {
		t.Errorf("Expected ErrEmailMissing, Got: %s", err.Error())
	}
	invalidEmails := []string{
		"fake.com",
		"fake@123",
		"fake@.com",
		"////@/./",
		"fake",
		"fake@something.z",
		"@something.com",
	}

	for _, email := range invalidEmails {
		user = fakeUserService()
		user.Email = email
		err = s.User.Create(&user)
		if err != errorsModel.ErrEmailInvalid {
			t.Log(user.Email)
			t.Errorf("Expected ErrEmailInvalid, Got: %s", err.Error())
		}
	}
}

func TestCreateWithInvalidPassword(t *testing.T) {
	s, err := mockServices(false)
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()
	user := fakeUserService()
	user.Password = ""
	err = s.User.Create(&user)
	if err != errorsModel.ErrPasswordRequired {
		t.Errorf("Expected an ErrPasswordRequired error, Got: %s", err.Error())
	}
	user = fakeUserService()
	user.Password = "short"
	err = s.User.Create(&user)
	if err != errorsModel.ErrPasswordTooShort {
		t.Errorf("Expected an ErrPasswordTooShort error, Got: %s", err.Error())
	}
}

func TestUserById(t *testing.T) {
	s, err := mockServices(false)
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()

	fakeUser := fakeUserService()
	user := fakeUserService()

	if err := s.User.Create(&user); err != nil {
		t.Fatal(err)
	}

	newUsr, err := s.User.ByID(user.ID)
	if err != nil {
		t.Fatal(err)
	}
	if newUsr.ID < 1 {
		t.Errorf("User ID is less than 1. Have: %d, Want: %d", newUsr.ID, 1)
	}
	if newUsr.Name != user.Name || newUsr.Name != fakeUser.Name {
		t.Errorf("User Name is incorrect. Have: %s, Want: %s", newUsr.Name, fakeUser.Name)
	}
	if newUsr.Email != user.Email || newUsr.Email != fakeUser.Email {
		t.Errorf("User Email is incorrect. Have: %s, Want: %s", newUsr.Email, fakeUser.Email)
	}
	if time.Since(newUsr.CreatedAt) > time.Duration(10*time.Second) {
		expectedMin := time.Now()
		expectedMax := expectedMin.Add(10 * time.Second)
		// Because I guess the Go epoch is 01-02-2006 03:04:05 PM (15:04:05)
		fString := "15:04:05"
		minStr := expectedMin.Format(fString)
		maxStr := expectedMax.Format(fString)
		actual := newUsr.CreatedAt.Format(fString)
		t.Errorf("Expected user to be created recently. Have: %v, Want: %v - %v", actual, minStr, maxStr)
	}
	if time.Since(newUsr.UpdatedAt) > time.Duration(10*time.Second) {
		expectedMin := time.Now()
		expectedMax := expectedMin.Add(10 * time.Second)
		// Because I guess the Go epoch is 01-02-2006 03:04:05 PM (15:04:05)
		fString := "15:04:05"
		minStr := expectedMin.Format(fString)
		maxStr := expectedMax.Format(fString)
		actual := newUsr.UpdatedAt.Format(fString)
		t.Errorf("Expected user to be updated recently. Have: %v, Want: %v - %v", actual, minStr, maxStr)
	}
}

func TestUserByInvalidId(t *testing.T) {
	s, err := mockServices(false)
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()
	var invalidId uint = 100
	_, err = s.User.ByID(invalidId)
	if err == nil {
		t.Errorf("Have: %s, Want: %s", err.Error(), errorsModel.ErrUserNotFound.Error())
	}
}

func TestClosedUserServiceConnection(t *testing.T) {
	s, err := mockServices(false)
	if err != nil {
		t.Fatal(err)
	}
	if err := s.Close(); err != nil {
		t.Fatalf("User service db connection errored on close: %s", err.Error())
	}
}

func TestBadDatabaseConnection(t *testing.T) {
	s, err := mockServices(true)
	if s != nil {
		t.Errorf("Expected no user service, Have: %+v", s)
	}
	if err == nil {
		t.Errorf("Expected an error, Have: err=%s", err.Error())
	}
}

func TestQueryWithClosedUserService(t *testing.T) {
	s, err := mockServices(false)
	if err != nil {
		t.Fatal(err)
	}
	err = s.Close()
	if err != nil {
		t.Fatal(err)
	}
	var id uint = 1
	_, err = s.User.ByID(id)
	if err == nil {
		t.Errorf("Have: %s, Want: %s", err.Error(), "Some Other Error")
	}
}

func TestUpdateUser(t *testing.T) {
	s, err := mockServices(false)
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()

	user := fakeUserService()
	remember := user.Remember
	password := user.Password
	if err := s.User.Create(&user); err != nil {
		t.Fatal(err)
	}
	user.Remember = remember
	newEmail := "fake.user.new@email.com"
	newName := "Fake User New"
	user.Name = newName
	user.Email = newEmail
	user.Password = password
	if err = s.User.Update(&user); err != nil {
		t.Fatalf("Update user failed: %s", err)
	}
	newUser, err := s.User.ByID(user.ID)
	if err != nil {
		t.Fatalf("Get user by ID failed: %s", err)
	}
	if newUser.ID != user.ID {
		t.Errorf("User ID doesn't match. Have: %d, Want: %d", newUser.ID, user.ID)
	}
	if newUser.Name != newName {
		t.Errorf("User Name wasn't updated. Have: %s, Want: %s", newUser.Name, newName)
	}
	if newUser.Email != newEmail {
		t.Errorf("User Email wasn't updated. Have: %s, Want: %s", newUser.Email, newEmail)
	}
	if time.Since(newUser.UpdatedAt) > time.Since(newUser.CreatedAt) {
		t.Errorf("Expected update time > created time. Updated: %+v, Created: %+v", newUser.UpdatedAt, newUser.CreatedAt)
	}
}

func TestUpdateUserWithNoChanges(t *testing.T) {
	s, err := mockServices(false)
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()

	user := fakeUserService()
	if err := s.User.Create(&user); err != nil {
		t.Fatalf("Expected to create user: %s", err.Error())
	}
	userCopy := fakeUserService()
	userCopy.ID = user.ID
	if err := s.User.Update(&userCopy); err != nil {
		t.Fatalf("Expected a successful update with no changes: %s", err.Error())
	}
}

func TestUserByEmail(t *testing.T) {
	s, err := mockServices(false)
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()

	user := fakeUserService()
	// Try creating a user with an empty remember token
	user.Remember = ""

	if err := s.User.Create(&user); err != nil {
		t.Fatal(err)
	}
	// We should expect the "ByEmail" method to handle normalization to lowercase.
	user.Email = strings.ToUpper(user.Email)
	newUsr, err := s.User.ByEmail(user.Email)
	if err != nil {
		t.Fatalf("Error getting user by email. %s", err.Error())
	}
	// convert original email back to lowercase. Expect them both to be lowercase.
	user.Email = strings.ToLower(user.Email)
	if newUsr.Email != user.Email {
		t.Errorf("Email invalid. Have: %s, Want: %s", newUsr.Email, user.Email)
	}
}

func TestUserByInvalidEmail(t *testing.T) {
	s, err := mockServices(false)
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()

	user := fakeUserService()

	_, err = s.User.ByEmail(user.Email)
	if err == nil {
		t.Fatalf("Expected an error. Got: %s", err.Error())
	}
	if err != errorsModel.ErrUserNotFound {
		t.Fatalf("Expected ErrNotFound. Got: %s", err.Error())
	}
}

func TestUserByEmailWithClosedConnection(t *testing.T) {
	s, err := mockServices(false)
	if err != nil {
		t.Fatal(err)
	}
	err = s.Close()
	if err != nil {
		t.Fatal(err)
	}

	user := fakeUserService()

	_, err = s.User.ByEmail(user.Email)
	if err == nil {
		t.Fatalf("Expected an error. Got: %s", err.Error())
	}
	if err == errorsModel.ErrUserNotFound {
		t.Fatalf("Expected Some Other Error. Got: %s", err.Error())
	}
}

func TestDeleteUserById(t *testing.T) {
	s, err := mockServices(false)
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()
	user := fakeUserService()
	if err := s.User.Create(&user); err != nil {
		t.Fatal(err)
	}
	if err := s.User.Delete(user.ID); err != nil {
		t.Fatalf("Expected no errors, Got: %s", err.Error())
	}
	_, err = s.User.ByID(user.ID)
	if err != errorsModel.ErrUserNotFound {
		t.Errorf("Expected ErrNotFound, Got: %s", err.Error())
	}
}

func TestDeleteUserByInvalidId(t *testing.T) {
	s, err := mockServices(false)
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()

	err = s.User.Delete(0)
	if err == nil {
		t.Fatalf("Expected an error, Got: %s", err.Error())
	}
	if err != errorsModel.ErrIdInvalid {
		t.Errorf("Expected ErrIdInvalid, Got: %s", err.Error())
	}
}

func TestDestructiveReset(t *testing.T) {
	s, err := mockServices(false)
	if err != nil {
		t.Fatal(err)
	}
	if err = s.Close(); err != nil {
		t.Fatalf("Error closing connection: %s", err.Error())
	}
	if err := s.AutoMigrate(); err == nil {
		t.Errorf("Expected an automigrate error with closed database: %s", err.Error())
	}
	if err := s.DestructiveReset(); err == nil {
		t.Errorf("Expected a desctructive reset error with closed database: %s", err.Error())
	}
}

func TestAuthenticateValidUser(t *testing.T) {
	s, err := mockServices(false)
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()
	user := fakeUserService()

	email := user.Email
	badEmail := user.Email + "123"
	password := user.Password
	badPassword := user.Password + user.Password

	err = s.User.Create(&user)
	if err != nil {
		t.Errorf("Expected a successful user creation. %s", err.Error())
	}
	_, err = s.User.Authenticate(email, password)
	if err != nil {
		t.Errorf("Email and password should have been correct: %s", err.Error())
	}
	_, err = s.User.Authenticate(email, badPassword)
	if err != errorsModel.ErrPasswordIncorrect {
		t.Errorf("Expected ErrPasswordIncorrect: %s", err.Error())
	}
	_, err = s.User.Authenticate(badEmail, password)
	if err != errorsModel.ErrUserNotFound {
		t.Errorf("Expected ErrNotFound: %s", err.Error())
	}
}

func TestInvalidRememberToken(t *testing.T) {
	s, err := mockServices(false)
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()
	user := fakeUserService()
	remember, err := rand.String(10)
	if err != nil {
		t.Errorf("Expected to be able to generate a random base64 string, Got: %s", err.Error())
	}
	user.Remember = remember
	err = s.User.Create(&user)
	if err != errorsModel.ErrRememberTokenTooShort {
		t.Errorf("Expected ErrRememberTokenTooShort, Got: %s", err.Error())
	}
}

func TestByRemember(t *testing.T) {
	s, err := mockServices(false)
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()
	user := fakeUserService()
	remember := user.Remember
	err = s.User.Create(&user)
	if err != nil {
		t.Errorf("Expected a successful user creation. %s", err.Error())
	}
	newUser, err := s.User.ByRemember(remember)
	if err != nil {
		t.Errorf("Should have found the user, %s", err.Error())
	}
	if newUser.ID != user.ID {
		t.Errorf("Got the wrong user. Have: %+v, Want: %+v", newUser, user)
	}
}
