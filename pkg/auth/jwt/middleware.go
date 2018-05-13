package jwt

import (
	"context"

	"github.com/dgrijalva/jwt-go"
	kitjwt "github.com/go-kit/kit/auth/jwt"
	"github.com/go-kit/kit/endpoint"
)

func MakeAuthenticatorMiddleware(signingSecret string) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		tokenValidationKey := []byte(signingSecret)
		key := func(token *jwt.Token) (interface{}, error) {
			return tokenValidationKey, nil
		}
		authenticator := kitjwt.NewParser(key, jwt.SigningMethodHS256, kitjwt.MapClaimsFactory)

		endpoint := func(ctx context.Context, request interface{}) (interface{}, error) {
			claims := ctx.Value(kitjwt.JWTClaimsContextKey).(jwt.MapClaims)
			subject := claims["sub"].(string)

			ctxWithIdentity := context.WithValue(ctx, "Identity", subject)
			return next(ctxWithIdentity, request)
		}

		return authenticator(endpoint)
	}
}
