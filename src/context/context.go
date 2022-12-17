package context

import (
	"context"

	"lenslocked/models"
)

const (
	userKey privateKey = "user"
)

type privateKey string

func WithUser(ctx context.Context, user *models.User) context.Context {
	return context.WithValue(ctx, userKey, user)
}

func User(ctx context.Context) *models.User {
	temp := ctx.Value(userKey)
	if temp == nil {
		return nil
	}
	if user, ok := temp.(*models.User); ok {
		return user
	}
	return nil
}
