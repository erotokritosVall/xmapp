package util

import "context"

type authTokenKey int

func SetAuthToken(ctx context.Context, value string) context.Context {
	c := context.WithValue(ctx, authTokenKey(1), value)

	return c
}

func GetAuthToken(ctx context.Context) any {
	return ctx.Value(authTokenKey(1))
}
