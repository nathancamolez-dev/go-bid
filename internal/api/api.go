package api

import (
	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi/v5"

	"github.com/nathancamolez-dev/go-bid/internal/services"
)

type Api struct {
	Router         *chi.Mux
	UserService    services.UserServices
	ProductService services.ProductService
	Sessions       *scs.SessionManager
}
