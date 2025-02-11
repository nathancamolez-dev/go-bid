package api

import (
	"errors"
	"net/http"

	"github.com/nathancamolez-dev/go-bid/internal/jsonutils"
	"github.com/nathancamolez-dev/go-bid/internal/services"
	"github.com/nathancamolez-dev/go-bid/internal/usecase/user"
)

func (api *Api) handleSignupUser(w http.ResponseWriter, r *http.Request) {
	data, problems, err := jsonutils.DecodeValidJson[user.CreateUserReq](r)
	if err != nil {
		_ = jsonutils.EncodeJson(w, r, http.StatusUnprocessableEntity, problems)
		return
	}
	id, err := api.UserService.CreateUser(r.Context(),
		data.UserName, data.Email, data.Password, data.Bio)
	if err != nil {
		if errors.Is(err, services.ErrDuplicatedEmailOrUsername) {
			_ = jsonutils.EncodeJson(w, r, http.StatusUnprocessableEntity, map[string]any{
				"error": "Already exists a user with this email or username",
			})
			return
		}
	}

	_ = jsonutils.EncodeJson(w, r, http.StatusUnprocessableEntity, map[string]any{
		"user_id": id,
	})
}

func (api *Api) handleLoginUser(w http.ResponseWriter, r *http.Request) {

	panic("TODO")
}

func (api *Api) handleLogout(w http.ResponseWriter, r *http.Request) {

	panic("TODO")
}
