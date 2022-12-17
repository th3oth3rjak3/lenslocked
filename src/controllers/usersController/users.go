package usersController

import (
	"fmt"
	"net/http"

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
	u.SignupView.Render(w, nil)
}

// Used to process the signup request for a new user.
//
// POST /signup
func (u *UsersController) Create(w http.ResponseWriter, r *http.Request) {
	formData := &SignupForm{}
	var vd views.Data
	if err := formData.Bind(r); err != nil {
		vd.SetAlert(err, true)
		u.SignupView.Render(w, vd)
		return
	}
	user := &usersModel.User{
		Name:     formData.Name,
		Email:    formData.Email,
		Password: formData.Password,
	}
	if err := u.userService.Create(user); err != nil {
		vd.SetAlert(err, true)
		u.SignupView.Render(w, vd)
		return
	}
	err := u.signIn(w, user)
	if err != nil {
		// In this case the user was created, but couldn't login for some weird reason.
		// We're going to handle this by letting the user attempt to login after
		// a redirect to the login page. This should be an edge case if it ever happens.
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	http.Redirect(w, r, "/cookietest", http.StatusFound)
}

// Login is used to verify the provided email address and password and log
// the user in if they have an account and the credentials are correct.
//
// POST /login
func (u *UsersController) Login(w http.ResponseWriter, r *http.Request) {
	var vd views.Data
	formData := &LoginForm{}
	if err := formData.Bind(r); err != nil {
		vd.SetAlert(err, true)
		u.LoginView.Render(w, vd)
		return
	}
	usr, err := u.userService.Authenticate(formData.Email, formData.Password)
	if err != nil {
		vd.SetAlert(err, true)
		u.LoginView.Render(w, vd)
		return
	}

	err = u.signIn(w, usr)
	if err != nil {
		vd.SetAlert(err, true)
		u.LoginView.Render(w, vd)
		return
	}
	http.Redirect(w, r, "/cookietest", http.StatusFound)
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

// CookieTest is a route handler to display cookie information for testing purposes only.
func (u *UsersController) CookieTest(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("remember_token")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	usr, err := u.userService.ByRemember(cookie.Value)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	fmt.Fprintf(w, "Should be a user: %+v", usr)
}
