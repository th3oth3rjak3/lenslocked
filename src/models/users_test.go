package models

import (
	"os"
	"strings"
	"testing"
	"time"

	"lenslocked/rand"

	"github.com/joho/godotenv"
)

func mockUserService(causeDbError bool, causeEnvError bool) (UserService, error) {
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
	// Log mode set to false...
	us.LogMode(false)

	// Clear the users table between tests.
	us.DestructiveReset()
	return us, nil
}

func fakeUserService() User {
	remember, err := rand.RememberToken()
	if err != nil {
		panic(err)
	}
	name := "Fake User"
	email := "fake.user@email.com"
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

func TestCreateDuplicateEmailUser(t *testing.T) {
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
	if err := us.Create(&fakeUser); err != ErrEmailTaken {
		t.Errorf("Expected ErrEmailTaken, Got: %s", err.Error())
	}
}

func TestCreateUserWithInvalidEmail(t *testing.T) {
	us, err := mockUserService(false, false)
	if err != nil {
		t.Fatalf("Expected mockUserService, %s", err.Error())
	}
	defer us.Close()
	user := fakeUserService()
	user.Email = ""
	err = us.Create(&user)
	if err != ErrEmailMissing {
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
		err = us.Create(&user)
		if err != ErrEmailInvalid {
			t.Log(user.Email)
			t.Errorf("Expected ErrEmailInvalid, Got: %s", err.Error())
		}
	}
}

func TestCreateWithInvalidPassword(t *testing.T) {
	us, err := mockUserService(false, false)
	if err != nil {
		t.Fatal(err)
	}
	defer us.Close()
	user := fakeUserService()
	user.Password = ""
	err = us.Create(&user)
	if err != ErrPasswordRequired {
		t.Errorf("Expected an ErrPasswordRequired error, Got: %s", err.Error())
	}
	user = fakeUserService()
	user.Password = "short"
	err = us.Create(&user)
	if err != ErrPasswordTooShort {
		t.Errorf("Expected an ErrPasswordTooShort error, Got: %s", err.Error())
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
		t.Errorf("Have: %s, Want: %s", err.Error(), ErrNotFound.Error())
	}
}

func TestClosedUserServiceConnection(t *testing.T) {
	us, err := mockUserService(false, false)
	if err != nil {
		t.Fatal(err)
	}
	if err := us.Close(); err != nil {
		t.Fatalf("User service db connection errored on close: %s", err.Error())
	}
}

func TestBadDatabaseConnection(t *testing.T) {
	us, err := mockUserService(true, false)
	if us != nil {
		t.Errorf("Expected no user service, Have: %+v", us)
	}
	if err == nil {
		t.Errorf("Expected an error, Have: err=%s", err.Error())
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
		t.Errorf("Have: %s, Want: %s", err.Error(), "Some Other Error")
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
	password := user.Password
	if err := us.Create(&user); err != nil {
		t.Fatal(err)
	}
	user.Remember = remember
	newEmail := "fake.user.new@email.com"
	newName := "Fake User New"
	user.Name = newName
	user.Email = newEmail
	user.Password = password
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

func TestUpdateUserWithNoChanges(t *testing.T) {
	us, err := mockUserService(false, false)
	if err != nil {
		t.Fatal(err)
	}
	defer us.Close()

	user := fakeUserService()
	if err := us.Create(&user); err != nil {
		t.Fatalf("Expected to create user: %s", err.Error())
	}
	userCopy := fakeUserService()
	userCopy.ID = user.ID
	if err := us.Update(&userCopy); err != nil {
		t.Fatalf("Expected a successful update with no changes: %s", err.Error())
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
	// We should expect the "ByEmail" method to handle normalization to lowercase.
	user.Email = strings.ToUpper(user.Email)
	newUsr, err := us.ByEmail(user.Email)
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
	us, err := mockUserService(false, false)
	if err != nil {
		t.Fatal(err)
	}
	defer us.Close()

	user := fakeUserService()

	_, err = us.ByEmail(user.Email)
	if err == nil {
		t.Fatalf("Expected an error. Got: %s", err.Error())
	}
	if err != ErrNotFound {
		t.Fatalf("Expected ErrNotFound. Got: %s", err.Error())
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
		t.Fatalf("Expected an error. Got: %s", err.Error())
	}
	if err == ErrNotFound {
		t.Fatalf("Expected Some Other Error. Got: %s", err.Error())
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
		t.Fatalf("Expected no errors, Got: %s", err.Error())
	}
	_, err = us.ByID(user.ID)
	if err != ErrNotFound {
		t.Errorf("Expected ErrNotFound, Got: %s", err.Error())
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
		t.Fatalf("Expected an error, Got: %s", err.Error())
	}
	if err != ErrIdInvalid {
		t.Errorf("Expected ErrIdInvalid, Got: %s", err.Error())
	}
}

func TestDestructiveReset(t *testing.T) {
	us, err := mockUserService(false, false)
	if err != nil {
		t.Fatal(err)
	}
	if err = us.Close(); err != nil {
		t.Fatalf("Error closing connection: %s", err.Error())
	}
	if err := us.AutoMigrate(); err == nil {
		t.Errorf("Expected an automigrate error with closed database: %s", err.Error())
	}
	if err := us.DestructiveReset(); err == nil {
		t.Errorf("Expected a desctructive reset error with closed database: %s", err.Error())
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
	badEmail := user.Email + "123"
	password := user.Password
	badPassword := user.Password + user.Password

	err = us.Create(&user)
	if err != nil {
		t.Errorf("Expected a successful user creation. %s", err.Error())
	}
	_, err = us.Authenticate(email, password)
	if err != nil {
		t.Errorf("Email and password should have been correct: %s", err.Error())
	}
	_, err = us.Authenticate(email, badPassword)
	if err != ErrPasswordIncorrect {
		t.Errorf("Expected ErrPasswordIncorrect: %s", err.Error())
	}
	_, err = us.Authenticate(badEmail, password)
	if err != ErrNotFound {
		t.Errorf("Expected ErrNotFound: %s", err.Error())
	}
}

func TestInvalidRememberToken(t *testing.T) {
	us, err := mockUserService(false, false)
	if err != nil {
		t.Fatal(err)
	}
	defer us.Close()
	user := fakeUserService()
	remember, err := rand.String(10)
	if err != nil {
		t.Errorf("Expected to be able to generate a random base64 string, Got: %s", err.Error())
	}
	user.Remember = remember
	err = us.Create(&user)
	if err != ErrRememberTokenTooShort {
		t.Errorf("Expected ErrRememberTokenTooShort, Got: %s", err.Error())
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
		t.Errorf("Expected a successful user creation. %s", err.Error())
	}
	newUser, err := us.ByRemember(remember)
	if err != nil {
		t.Errorf("Should have found the user, %s", err.Error())
	}
	if newUser.ID != user.ID {
		t.Errorf("Got the wrong user. Have: %+v, Want: %+v", newUser, user)
	}
}

func TestMissingEnvironment(t *testing.T) {
	_, err := mockUserService(false, true)
	if err != ErrEnvironmentUnset {
		t.Errorf("Expected ErrEnvironmentUnset, Got: %s", err.Error())
	}
}
