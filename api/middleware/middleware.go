package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/golang-jwt/jwt/v5"

	"github.com/stanislavCasciuc/atom-fit/api/response"
	"github.com/stanislavCasciuc/atom-fit/internal/auth"
	"github.com/stanislavCasciuc/atom-fit/internal/store"
)

type userKey string

const UserCtx userKey = "user"

type Middleware struct {
	store         store.Storage
	resp          response.Responser
	authenticator auth.Authenticator
}

func New(
	store store.Storage,
	resp response.Responser,
	authenticator auth.Authenticator,
) Middleware {
	return Middleware{store, resp, authenticator}
}

func (m *Middleware) AuthTokenMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			m.resp.UnauthorizedError(w, r, fmt.Errorf("authorization header is missing"))
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			m.resp.UnauthorizedError(w, r, fmt.Errorf("authorization header is malformed"))
			return
		}

		token := parts[1]
		jwtToken, err := m.authenticator.ValidateToken(token)
		if err != nil {
			m.resp.UnauthorizedError(w, r, err)
			return
		}

		claims, _ := jwtToken.Claims.(jwt.MapClaims)

		userID, err := strconv.ParseInt(fmt.Sprintf("%.f", claims["sub"]), 10, 64)
		if err != nil {
			m.resp.UnauthorizedError(w, r, err)
			return
		}

		ctx := r.Context()

		user, err := m.store.Users.GetByID(ctx, userID)
		if err != nil {
			m.resp.UnauthorizedError(w, r, err)
			return
		}

		ctx = context.WithValue(ctx, UserCtx, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
