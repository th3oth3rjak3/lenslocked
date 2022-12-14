package models

import (
	"os"
	"strings"
	"testing"
	"time"

	"github.com/joho/godotenv"
)

func mockUserService(causeDbError bool, causeEnvError bool) (*UserService, error) {
	if os.Getenv("GITHUB_ACTION_STATUS_INDICATOR") != "true" {
		err := godotenv.Load("../.env")
		if err != nil {
			panic(err)
		}
	}
	psqlInfo := os.Getenv("DB_CONNECTION_STRING_TEST")
	if causeDbError {
		psqlInfo = os.Getenv("DB_CONNECTION_STRING_ERROR")
	}
	if causeEnvError {
		os.Unsetenv("HASH_KEY")
	}
	us, err := NewUserService(psqlInfo)
	if err != nil {
		return nil, err
	}
	us.db.LogMode(false)

	// Clear the users table between tests.
	us.DestructiveReset()
	return us, nil
}

func fakeUserService() User {
	name := "Fake User"
	email := "fake.user@email.com"
	remember := "special_remember_token"
	password := "some special password"

	return User{
		Name:     name,
		Email:    email,
		Password: password,
		Remember: remember,
	}
}

func TestCreateUser(t *testing.T) {
	us, err := mockUserService(false, false)
	if err != nil {
		t.Fatal(err)
	}
	defer us.Close()

	fakeUser := fakeUserService()
	user := fakeUserService()

	if err := us.Create(&user); err != nil {
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

func TestCreateWithEmptyPassword(t *testing.T) {
	us, err := mockUserService(false, false)
	if err != nil {
		t.Fatal(err)
	}
	defer us.Close()
	user := fakeUserService()
	user.Password = ""
	err = us.Create(&user)
	if err != ErrInvalidPassword {
		t.Errorf("Expected an ErrInvalidPassword error, Got: %s", err)
	}
}

func TestUserById(t *testing.T) {
	us, err := mockUserService(false, false)
	if err != nil {
		t.Fatal(err)
	}
	defer us.Close()

	fakeUser := fakeUserService()
	user := fakeUserService()

	if err := us.Create(&user); err != nil {
		t.Fatal(err)
	}

	newUsr, err := us.ByID(user.ID)
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
	us, err := mockUserService(false, false)
	if err != nil {
		t.Fatal(err)
	}
	defer us.Close()
	var invalidId uint = 100
	_, err = us.ByID(invalidId)
	if err == nil {
		t.Errorf("Have: %s, Want: %s", err, ErrNotFound)
	}
}

func TestCloseUserServiceConnection(t *testing.T) {
	us, err := mockUserService(false, false)
	if err != nil {
		t.Fatal(err)
	}
	if err := us.Close(); err != nil {
		t.Fatalf("User service db connection errored on close: %s", err)
	}
}

func TestBadDatabaseConnection(t *testing.T) {
	us, err := mockUserService(true, false)
	if us != nil {
		t.Errorf("Expected no user service, Have: %+v", us)
	}
	if err == nil {
		t.Errorf("Expected an error, Have: err=%s", err)
	}
}

func TestQueryWithClosedUserService(t *testing.T) {
	us, err := mockUserService(false, false)
	if err != nil {
		t.Fatal(err)
	}
	err = us.Close()
	if err != nil {
		t.Fatal(err)
	}
	var id uint = 1
	_, err = us.ByID(id)
	if err == nil {
		t.Errorf("Have: %s, Want: %s", err, "Some Other Error")
	}
}

func TestUpdateUser(t *testing.T) {
	us, err := mockUserService(false, false)
	if err != nil {
		t.Fatal(err)
	}
	defer us.Close()

	user := fakeUserService()
	remember := user.Remember

	if err := us.Create(&user); err != nil {
		t.Fatal(err)
	}
	user.Remember = remember
	newEmail := "fake.user.new@email.com"
	newName := "Fake User New"
	user.Name = newName
	user.Email = newEmail
	if err = us.Update(&user); err != nil {
		t.Fatalf("Update user failed: %s", err)
	}
	newUser, err := us.ByID(user.ID)
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

func TestUserByEmail(t *testing.T) {
	us, err := mockUserService(false, false)
	if err != nil {
		t.Fatal(err)
	}
	defer us.Close()

	user := fakeUserService()
	// Try creating a user with an empty remember token
	user.Remember = ""

	if err := us.Create(&user); err != nil {
		t.Fatal(err)
	}
	newUsr, err := us.ByEmail(user.Email)
	if err != nil {
		t.Fatalf("Error getting user by email. %s", err)
	}
	if newUsr.Email != user.Email {
		t.Errorf("Email invalid. Have: %s, Want: %s", newUsr.Email, user.Email)
	}
}

func TestUserByInvalidEmail(t *testing.T) {
	us, err := mockUserService(false, false)
	if err != nil {
		t.Fatal(err)
	}
	defer us.Close()

	user := fakeUserService()

	_, err = us.ByEmail(user.Email)
	if err == nil {
		t.Fatalf("Expected an error. Got: %s", err)
	}
	if err != ErrNotFound {
		t.Fatalf("Expected ErrNotFound. Got: %s", err)
	}
}

func TestUserByEmailWithClosedConnection(t *testing.T) {
	us, err := mockUserService(false, false)
	if err != nil {
		t.Fatal(err)
	}
	err = us.Close()
	if err != nil {
		t.Fatal(err)
	}

	user := fakeUserService()

	_, err = us.ByEmail(user.Email)
	if err == nil {
		t.Fatalf("Expected an error. Got: %s", err)
	}
	if err == ErrNotFound {
		t.Fatalf("Expected Some Other Error. Got: %s", err)
	}
}

func TestDeleteUserById(t *testing.T) {
	us, err := mockUserService(false, false)
	if err != nil {
		t.Fatal(err)
	}
	defer us.Close()
	user := fakeUserService()
	if err := us.Create(&user); err != nil {
		t.Fatal(err)
	}
	if err := us.Delete(user.ID); err != nil {
		t.Fatalf("Expected no errors, Got: %s", err)
	}
	_, err = us.ByID(user.ID)
	if err != ErrNotFound {
		t.Errorf("Expected ErrNotFound, Got: %s", err)
	}
}

func TestDeleteUserByInvalidId(t *testing.T) {
	us, err := mockUserService(false, false)
	if err != nil {
		t.Fatal(err)
	}
	defer us.Close()

	err = us.Delete(0)
	if err == nil {
		t.Fatalf("Expected an error, Got: %s", err)
	}
	if err != ErrInvalidId {
		t.Errorf("Expected ErrInvalidId, Got: %s", err)
	}
}

func TestDestructiveReset(t *testing.T) {
	us, err := mockUserService(false, false)
	if err != nil {
		t.Fatal(err)
	}
	if err = us.Close(); err != nil {
		t.Fatalf("Error closing connection: %s", err)
	}
	if err := us.AutoMigrate(); err == nil {
		t.Errorf("Expected an automigrate error with closed database: %s", err)
	}
	if err := us.DestructiveReset(); err == nil {
		t.Errorf("Expected a desctructive reset error with closed database: %s", err)
	}
}

func TestAuthenticateValidUser(t *testing.T) {
	us, err := mockUserService(false, false)
	if err != nil {
		t.Fatal(err)
	}
	defer us.Close()
	user := fakeUserService()

	email := user.Email
	capsEmail := strings.ToUpper(user.Email)
	badEmail := user.Email + "123"
	password := user.Password
	badPassword := user.Password + user.Password

	err = us.Create(&user)
	if err != nil {
		t.Errorf("Expected a successful user creation. %s", err)
	}
	_, err = us.Authenticate(email, password)
	if err != nil {
		t.Errorf("Email and password should have been correct: %s", err)
	}
	_, err = us.Authenticate(capsEmail, password)
	if err != nil {
		t.Errorf("Should be case insensitive: %s", err)
	}
	_, err = us.Authenticate(email, badPassword)
	if err != ErrInvalidPassword {
		t.Errorf("Expected ErrInvalidPassword: %s", err)
	}
	_, err = us.Authenticate(badEmail, password)
	if err != ErrNotFound {
		t.Errorf("Expected ErrNotFound: %s", err)
	}
}

func TestByRemember(t *testing.T) {
	us, err := mockUserService(false, false)
	if err != nil {
		t.Fatal(err)
	}
	defer us.Close()
	user := fakeUserService()
	remember := user.Remember
	err = us.Create(&user)
	if err != nil {
		t.Errorf("Expected a successful user creation. %s", err)
	}
	newUser, err := us.ByRemember(remember)
	if err != nil {
		t.Errorf("Should have found the user, %s", err)
	}
	if newUser.ID != user.ID {
		t.Errorf("Got the wrong user. Have: %+v, Want: %+v", newUser, user)
	}
}

func TestMissingEnvironment(t *testing.T) {
	_, err := mockUserService(false, true)
	if err != ErrEnvironmentUnset {
		t.Errorf("Expected ErrEnvironmentUnset, Got: %s", err)
	}
}
