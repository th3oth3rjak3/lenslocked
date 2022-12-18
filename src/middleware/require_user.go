package middleware

import (
	"fmt"
	"net/http"

	"lenslocked/context"
	"lenslocked/models/usersModel"
)

type User struct {
	UserService usersModel.UserService
}

type RequireUser struct {
	User
}

type UserCtx string

func (mw *User) Invoke(next http.Handler) http.Handler {
	return mw.InvokeFn(next.ServeHTTP)
}

func (mw *User) InvokeFn(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("remember_token")
		if err != nil {
			next(w, r)
			return
		}
		usr, err := mw.UserService.ByRemember(cookie.Value)
		if err != nil {
			next(w, r)
			return
		}

		ctx := context.WithUser(r.Context(), usr)
		r = r.WithContext(ctx)
		fmt.Println("User Found: ", usr)
		next(w, r)
	})
}

func (mw *RequireUser) Invoke(next http.Handler) http.Handler {
	return mw.InvokeFn(next.ServeHTTP)
}

func (mw *RequireUser) InvokeFn(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := context.User(r.Context())
		if user == nil {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
		next(w, r)
	})
}
