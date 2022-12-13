package models

import (
	"fmt"
	"testing"
	"time"
)

func mockUserService(causeError bool) (*UserService, error) {
	var (
		host = "localhost"
		port = 5432
		user = "Jake"
		// Add password after user in the connection string if you have one.
		// password = ""
		dbname = "lenslocked_test"
	)
	if causeError {
		user = "SomeUserThatDoesntExist"
	}
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=disable", host, port, user, dbname)
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

	name := "Jake"
	email := "jake.d.hathaway@icloud.com"

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
	name := "Jake"
	email := "jake.d.hathaway@icloud.com"

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

func TestInvalidUser(t *testing.T) {
	us, err := mockUserService(false)
	if err != nil {
		t.Fatal(err)
	}
	defer us.Close()
	var invalidId uint = 100
	newUsr, err := us.ByID(invalidId)
	if err == nil {
		t.Errorf("Have: %s, Want: %s", err, ErrNotFound)
	}
	if newUsr != nil {
		t.Errorf("Have a user: %+v, Want: %+v", newUsr, nil)
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
	newUsr, err := us.ByID(id)
	if err == nil {
		t.Errorf("Have: %s, Want: %s", err, "Some Other Error")
	}
	if newUsr != nil {
		t.Errorf("Have a user: %+v, Want: %+v", newUsr, nil)
	}
}
