package middleware

import (
	"net/http"

	"lenslocked/context"
	"lenslocked/models/galleriesModel"
	"lenslocked/models/usersModel"
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

//  func (mw *RequireUser) ImageSafety(next echo.HandlerFunc) echo.HandlerFunc {
// 	return (func(c echo.Context) error {
// 		r := c.Request()
// 		cookie, err := r.Cookie("remember_token")
// 		if err != nil {
// 			log.Println(err)
// 			return nil
// 		}
// 		usr, err := mw.UserService.ByRemember(cookie.Value)
// 		if err != nil {
// 			log.Println(err)
// 			return nil
// 		}

// 		galleryId := c.Param("galleryId")
// 		gid, err := strconv.Atoi(galleryId)
// 		if err != nil {
// 			log.Println(err)
// 			return nil
// 		}
// 		gallery, err := mw.GalleryService.ByID(uint(gid))
// 		if err != nil {
// 			log.Println(err)
// 			return nil
// 		}
// 		if gallery.UserID != usr.ID {
// 			log.Println(errors.New("gallery user_id doesn't match user id"))
// 			return nil
// 		}
// 		return next(c)
// 	})
// }
