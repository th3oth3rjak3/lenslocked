package middleware

import (
	"fmt"
	"net/http"

	"lenslocked/context"
	"lenslocked/models/usersModel"
)

type RequireUser struct {
	UserService usersModel.UserService
}

type UserCtx string

func (mw *RequireUser) Invoke(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("remember_token")
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
		usr, err := mw.UserService.ByRemember(cookie.Value)
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		ctx := context.WithUser(r.Context(), usr)
		r = r.WithContext(ctx)
		fmt.Println("User Found: ", usr)
		next.ServeHTTP(w, r)
	})
}
