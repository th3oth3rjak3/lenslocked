package models

import (
	"os"
	"testing"
	"time"

	"github.com/joho/godotenv"
)

func mockUserService(causeError bool) (*UserService, error) {
	err := godotenv.Load("../.env")
	if err != nil {
		panic(err)
	}
	psqlInfo := os.Getenv("DB_CONNECTION_STRING_TEST")
	if causeError {
		psqlInfo = os.Getenv("DB_CONNECTION_STRING_ERROR")
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

func TestCreateUser(t *testing.T) {
	us, err := mockUserService(false)
	if err != nil {
		t.Fatal(err)
	}
	defer us.Close()

	name := "Fake User"
	email := "fake.user@email.com"

	usr := User{
		Name:  name,
		Email: email,
	}
	if err := us.Create(&usr); err != nil {
		t.Fatal(err)
	}

	if usr.ID < 1 {
		t.Errorf("User ID is less than 1. Have: %d, Want: %d", usr.ID, 1)
	}
	if usr.Name != name {
		t.Errorf("User Name is incorrect. Have: %s, Want: %s", usr.Name, name)
	}
	if usr.Email != email {
		t.Errorf("User Email is incorrect. Have: %s, Want: %s", usr.Email, email)
	}
	if time.Since(usr.CreatedAt) > time.Duration(10*time.Second) {
		expectedMin := time.Now()
		expectedMax := expectedMin.Add(10 * time.Second)
		// Because I guess the Go epoch is 01-02-2006 03:04:05 PM (15:04:05)
		fString := "15:04:05"
		minStr := expectedMin.Format(fString)
		maxStr := expectedMax.Format(fString)
		actual := usr.CreatedAt.Format(fString)
		t.Errorf("Expected user to be created recently. Have: %v, Want: %v - %v", actual, minStr, maxStr)
	}
	if time.Since(usr.UpdatedAt) > time.Duration(10*time.Second) {
		expectedMin := time.Now()
		expectedMax := expectedMin.Add(10 * time.Second)
		// Because I guess the Go epoch is 01-02-2006 03:04:05 PM (15:04:05)
		fString := "15:04:05"
		minStr := expectedMin.Format(fString)
		maxStr := expectedMax.Format(fString)
		actual := usr.UpdatedAt.Format(fString)
		t.Errorf("Expected user to be updated recently. Have: %v, Want: %v - %v", actual, minStr, maxStr)
	}
}

func TestUserById(t *testing.T) {
	us, err := mockUserService(false)
	if err != nil {
		t.Fatal(err)
	}
	defer us.Close()
	name := "Fake User"
	email := "fake.user@email.com"

	usr := User{
		Name:  name,
		Email: email,
	}
	if err := us.Create(&usr); err != nil {
		t.Fatal(err)
	}

	newUsr, err := us.ByID(usr.ID)
	if err != nil {
		t.Fatal(err)
	}
	if newUsr.ID < 1 {
		t.Errorf("User ID is less than 1. Have: %d, Want: %d", usr.ID, 1)
	}
	if newUsr.Name != name {
		t.Errorf("User Name is incorrect. Have: %s, Want: %s", usr.Name, name)
	}
	if newUsr.Email != email {
		t.Errorf("User Email is incorrect. Have: %s, Want: %s", usr.Email, email)
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
	us, err := mockUserService(false)
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
	us, err := mockUserService(false)
	if err != nil {
		t.Fatal(err)
	}
	if err := us.Close(); err != nil {
		t.Fatalf("User service db connection errored on close: %s", err)
	}
}

func TestBadDatabaseConnection(t *testing.T) {
	us, err := mockUserService(true)
	if us != nil {
		t.Errorf("Expected no user service, Have: %+v", us)
	}
	if err == nil {
		t.Errorf("Expected an error, Have: err=%s", err)
	}
}

func TestQueryWithClosedUserService(t *testing.T) {
	us, err := mockUserService(false)
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
	us, err := mockUserService(false)
	if err != nil {
		t.Fatal(err)
	}
	defer us.Close()
	name := "Fake User"
	email := "fake.user@email.com"

	usr := User{
		Name:  name,
		Email: email,
	}
	if err := us.Create(&usr); err != nil {
		t.Fatal(err)
	}
	newEmail := "fake.user.new@email.com"
	newName := "Fake User New"
	usr.Name = newName
	usr.Email = newEmail
	if err = us.Update(&usr); err != nil {
		t.Fatalf("Update user failed: %s", err)
	}
	newUser, err := us.ByID(usr.ID)
	if err != nil {
		t.Fatalf("Get user by ID failed: %s", err)
	}
	if newUser.ID != usr.ID {
		t.Errorf("User ID doesn't match. Have: %d, Want: %d", newUser.ID, usr.ID)
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
	us, err := mockUserService(false)
	if err != nil {
		t.Fatal(err)
	}
	defer us.Close()

	name := "Fake User"
	email := "fake.user@email.com"

	usr := User{
		Name:  name,
		Email: email,
	}
	if err := us.Create(&usr); err != nil {
		t.Fatal(err)
	}
	newUsr, err := us.ByEmail(email)
	if err != nil {
		t.Fatalf("Error getting user by email. %s", err)
	}
	if newUsr.Email != email {
		t.Errorf("Email invalid. Have: %s, Want: %s", newUsr.Email, email)
	}
}

func TestUserByInvalidEmail(t *testing.T) {
	us, err := mockUserService(false)
	if err != nil {
		t.Fatal(err)
	}
	defer us.Close()

	email := "fake.user@email.com"

	_, err = us.ByEmail(email)
	if err == nil {
		t.Fatalf("Expected an error. Got: %s", err)
	}
	if err != ErrNotFound {
		t.Fatalf("Expected ErrNotFound. Got: %s", err)
	}
}

func TestUserByEmailWithClosedConnection(t *testing.T) {
	us, err := mockUserService(false)
	if err != nil {
		t.Fatal(err)
	}
	err = us.Close()
	if err != nil {
		t.Fatal(err)
	}

	email := "fake.user@email.com"

	_, err = us.ByEmail(email)
	if err == nil {
		t.Fatalf("Expected an error. Got: %s", err)
	}
	if err == ErrNotFound {
		t.Fatalf("Expected Some Other Error. Got: %s", err)
	}
}

func TestDeleteUserById(t *testing.T) {
	us, err := mockUserService(false)
	if err != nil {
		t.Fatal(err)
	}
	defer us.Close()
	name := "Fake User"
	email := "fake.user@email.com"

	usr := User{
		Name:  name,
		Email: email,
	}
	if err := us.Create(&usr); err != nil {
		t.Fatal(err)
	}
	if err := us.Delete(usr.ID); err != nil {
		t.Fatalf("Expected no errors, Got: %s", err)
	}
}

func TestDeleteUserByInvalidId(t *testing.T) {
	us, err := mockUserService(false)
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
	us, err := mockUserService(false)
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
