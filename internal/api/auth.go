package api

import (
	"net/http"

	"github.com/gorilla/csrf"

	"github.com/nathancamolez-dev/go-bid/internal/jsonutils"
)

func (api *Api) HandleGetCSRFtoken(w http.ResponseWriter, r *http.Request) {
	token := csrf.Token(r)
	_ = jsonutils.EncodeJson(w, r, http.StatusOK, map[string]any{
		"csrf_token": token,
	})
}

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
