package middleware

import (
	"net/http"
	"strings"

	"github.com/google/uuid"

	v1 "github.com/anemptyemptiness/Go-Rocket-App/order/internal/client/grpc/iam/v1"
	"github.com/anemptyemptiness/Go-Rocket-App/platform/pkg/auth"
)

const SessionHeaderKey = "X-Session-UUID"

func AuthMiddleware(iamClient v1.Client) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			bearer := r.Header.Get("Authorization")
			sessionUUID := strings.TrimPrefix(bearer, "Bearer ")
			if sessionUUID == "" {
				sessionUUID = r.Header.Get(SessionHeaderKey)
				if sessionUUID == "" {
					w.WriteHeader(http.StatusUnauthorized)
					return
				}
			}

			session, user, err := iamClient.Whoiam(r.Context(), sessionUUID)
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			sessionUuid := session.UUID
			userUuid := user.UUID

			if sessionUuid == "" || userUuid == "" {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			userUUID, err := uuid.Parse(userUuid)
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			ctx := auth.WithSessionUUID(r.Context(), sessionUuid)
			ctx = auth.WithUserUUID(ctx, userUUID)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
