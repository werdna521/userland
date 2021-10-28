package request

import (
	"context"

	e "github.com/werdna521/userland/api/error"
	"github.com/werdna521/userland/api/middleware"
	"github.com/werdna521/userland/security/jwt"
)

func GetAccessTokenFromCtx(ctx context.Context) (*jwt.AccessToken, e.Error) {
	at, ok := ctx.Value(middleware.AccessTokenCtxKey).(*jwt.AccessToken)
	if !ok {
		return nil, e.NewBadRequestError("cannot parse access token")
	}

	return at, nil
}

func GetRefreshTokenFromCtx(ctx context.Context) (*jwt.RefreshToken, e.Error) {
	rt, ok := ctx.Value(middleware.RefreshTokenCtxKey).(*jwt.RefreshToken)
	if !ok {
		return nil, e.NewBadRequestError("cannot parse refresh token")
	}

	return rt, nil
}
