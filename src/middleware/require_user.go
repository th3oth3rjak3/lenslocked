package middleware

import (
	"net/http"
	"strconv"

	"lenslocked/context"
	"lenslocked/models/galleriesModel"
	"lenslocked/models/usersModel"

	"github.com/go-chi/chi/v5"
)

type User struct {
	UserService    usersModel.UserService
	GalleryService galleriesModel.GalleryService
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

func (mw *RequireUser) ImageSafety(next http.Handler) http.Handler {
	return mw.ImageSafetyFn(next.ServeHTTP)
}

func (mw *RequireUser) ImageSafetyFn(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := context.User(r.Context())
		galleryIDStr := chi.URLParam(r, "galleryID")
		galleryID, err := strconv.Atoi(galleryIDStr)
		if err != nil || galleryIDStr == "" {
			http.Error(w, "404 page not found", http.StatusNotFound)
			return
		}
		gallery, err := mw.User.GalleryService.ByID(uint(galleryID))
		if err != nil {
			http.Error(w, "404 page not found", http.StatusNotFound)
			return
		}
		if gallery.UserID != user.ID {
			http.Error(w, "404 page not found", http.StatusNotFound)
			return
		}
		next(w, r)
	})
}
