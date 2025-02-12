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
		jsonutils.EncodeJson(w, r, http.StatusUnprocessableEntity, map[string]any{
			"error": err.Error(),
		})
		if problems != nil {
			jsonutils.EncodeJson(w, r, http.StatusUnprocessableEntity, map[string]any{
				"error":    err.Error(),
				"problems": problems,
			})

		}
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

	_ = jsonutils.EncodeJson(w, r, http.StatusOK, map[string]any{
		"user_id": id,
	})
}

func (api *Api) handleLoginUser(w http.ResponseWriter, r *http.Request) {
	data, problems, err := jsonutils.DecodeValidJson[user.LoginUserReq](r)
	if err != nil {
		jsonutils.EncodeJson(w, r, http.StatusUnprocessableEntity, map[string]any{
			"error": err.Error(),
		})
		if problems != nil {
			jsonutils.EncodeJson(w, r, http.StatusUnprocessableEntity, map[string]any{
				"error":    err.Error(),
				"problems": problems,
			})

		}
		return
	}

	id, err := api.UserService.AuthenticateUser(r.Context(), data.Email, data.Password)
	if err != nil {
		if errors.Is(err, services.ErrInvalidCredentials) {
			_ = jsonutils.EncodeJson(w, r, http.StatusBadRequest, map[string]any{
				"error": "Invalid email or password",
			})
			return
		}
		jsonutils.EncodeJson(w, r, http.StatusInternalServerError, map[string]any{
			"error": "unexpected internal server error",
		})
		return
	}
	err = api.Sessions.RenewToken(r.Context())
	if err != nil {
		jsonutils.EncodeJson(w, r, http.StatusInternalServerError, map[string]any{
			"error": "unexpected internal server error",
		})
		return
	}

	api.Sessions.Put(r.Context(), "AuthenticateUserId", id)

	_ = jsonutils.EncodeJson(w, r, http.StatusOK, map[string]any{
		"message": "!!successfully logged in",
	})

}

func (api *Api) handleLogout(w http.ResponseWriter, r *http.Request) {
	err := api.Sessions.RenewToken(r.Context())
	if err != nil {
		jsonutils.EncodeJson(w, r, http.StatusInternalServerError, map[string]any{
			"error": "unexpected internal server error",
		})
		return
	}

	api.Sessions.Remove(r.Context(), "AuthenticateUserId")

	_ = jsonutils.EncodeJson(w, r, http.StatusOK, map[string]any{
		"message": "successfully logged out",
	})

}
