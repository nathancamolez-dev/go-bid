package api

import (
	"net/http"

	"github.com/nathancamolez-dev/go-bid/internal/jsonutils"
)

func (api *Api) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !api.Sessions.Exists(r.Context(), "AuthenticateUserId") {
			jsonutils.EncodeJson(w, r, http.StatusUnauthorized, map[string]any{
				"message": "unauthorized",
			})
			return
		}
		next.ServeHTTP(w, r)

	})
}
