package context

import (
	"context"

	"lenslocked/models/usersModel"
)

const (
	userKey privateKey = "user"
)

type privateKey string

func WithUser(ctx context.Context, user *usersModel.User) context.Context {
	return context.WithValue(ctx, userKey, user)
}

func User(ctx context.Context) *usersModel.User {
	temp := ctx.Value(userKey)
	if temp == nil {
		return nil
	}
	if user, ok := temp.(*usersModel.User); ok {
		return user
	}
	return nil
}
