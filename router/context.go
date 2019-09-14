package router

import (
	"context"

	"bounzr/iam/oauth2"
	"bounzr/iam/repository"
)

// key is an unexported type for keys defined in this package.
// This prevents collisions with keys defined in other packages.
type key int

// userCtxKey is the key for user.User values in Contexts. It is
// unexported; clients use user.NewContext and user.FromContext
// instead of using this key directly.
const (
	userCtxKey   key = 0
	clientCtxKey key = 1
)

func fromContextGetClient(ctx context.Context) (*oauth2.ClientCtx, bool) {
	c, ok := ctx.Value(clientCtxKey).(*oauth2.ClientCtx)
	return c, ok
}

// FromContext returns the User value stored in ctx, if any.
func fromContextGetUser(ctx context.Context) (*repository.UserCtx, bool) {
	u, ok := ctx.Value(userCtxKey).(*repository.UserCtx)
	return u, ok
}

func newContextWithClient(ctx context.Context, c *oauth2.ClientCtx) context.Context {
	return context.WithValue(ctx, clientCtxKey, c)
}

// NewContext returns a new Context that carries value u.
func newContextWithUser(ctx context.Context, u *repository.UserCtx) context.Context {
	return context.WithValue(ctx, userCtxKey, u)
}
