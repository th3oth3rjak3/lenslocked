package usersController

import (
	"net/http"
	"time"

	"lenslocked/context"
	"lenslocked/models/usersModel"
	"lenslocked/rand"
	"lenslocked/views"
)

// The Users controller object.
type UsersController struct {
	SignupView  *views.View
	LoginView   *views.View
	userService usersModel.UserService
}

// Instantiates a new Users controller.
// This will panic if templates are not parsed correctly.
// Only used during initial startup.
func NewUsersController(us usersModel.UserService) *UsersController {
	return &UsersController{
		SignupView:  views.NewView("bootstrap", "users/new"),
		LoginView:   views.NewView("bootstrap", "users/login"),
		userService: us,
	}
}

// Used to show the signup page to the user who wishes to signup.
//
// GET /signup
func (u *UsersController) Signup(w http.ResponseWriter, r *http.Request) {
	var form SignupForm
	form.BindURLParams(r)
	u.SignupView.Render(w, r, &form)
}

// Used to process the signup request for a new user.
//
// POST /signup
func (u *UsersController) Create(w http.ResponseWriter, r *http.Request) {
	formData := &SignupForm{}
	var vd views.Data
	vd.Payload = &formData
	if err := formData.Bind(r); err != nil {
		vd.SetAlert(err)
		u.SignupView.Render(w, r, vd)
		return
	}
	user := &usersModel.User{
		Name:     formData.Name,
		Email:    formData.Email,
		Password: formData.Password,
	}
	if err := u.userService.Create(user); err != nil {
		vd.SetAlert(err)
		u.SignupView.Render(w, r, vd)
		return
	}
	err := u.signIn(w, user)
	if err != nil {
		vd.SetAlert(err)
		u.LoginView.Render(w, r, vd)
		return
	}

	alert := views.Alert{
		Level:   views.AlertLevelSuccess,
		Message: "Welcome to LensLocked.com!",
	}
	views.RedirectAlert(w, r, "/galleries", http.StatusFound, alert)
}

// Login is used to verify the provided email address and password and log
// the user in if they have an account and the credentials are correct.
//
// POST /login
func (u *UsersController) Login(w http.ResponseWriter, r *http.Request) {
	var vd views.Data
	formData := &LoginForm{}
	if err := formData.Bind(r); err != nil {
		vd.SetAlert(err)
		u.LoginView.Render(w, r, vd)
		return
	}
	usr, err := u.userService.Authenticate(formData.Email, formData.Password)
	if err != nil {
		vd.SetAlert(err)
		u.LoginView.Render(w, r, vd)
		return
	}

	err = u.signIn(w, usr)
	if err != nil {
		vd.SetAlert(err)
		u.LoginView.Render(w, r, vd)
		return
	}
	http.Redirect(w, r, "/galleries", http.StatusFound)
}

// signIn is used to attach the signed-in cookie to the http response.
func (u *UsersController) signIn(w http.ResponseWriter, usr *usersModel.User) error {
	// If the remember token is empty, generate one.
	if usr.Remember == "" {
		token, err := rand.RememberToken()
		if err != nil {
			return err
		}
		usr.Remember = token
		err = u.userService.Update(usr)
		if err != nil {
			return err
		}
	}

	cookie := http.Cookie{
		Name:     "remember_token",
		Value:    usr.Remember,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	}
	http.SetCookie(w, &cookie)
	return nil
}

// Logout is used to delete a user's session cookie, remember token, and then will update
// the user resource with a new remember token.
//
// POST /logout
func (u *UsersController) Logout(w http.ResponseWriter, r *http.Request) {
	exp := time.Now().Add(time.Hour * -24)
	cookie := http.Cookie{
		Name:     "remember_token",
		Value:    "",
		Expires:  exp,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	}
	http.SetCookie(w, &cookie)
	user := context.User(r.Context())
	token, _ := rand.RememberToken()
	user.Remember = token
	u.userService.Update(user)
	http.Redirect(w, r, "/", http.StatusFound)
}
