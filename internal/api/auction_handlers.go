package api

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/nathancamolez-dev/go-bid/internal/jsonutils"
	"github.com/nathancamolez-dev/go-bid/internal/services"
)

func (api *Api) handleSubscribeToAuction(w http.ResponseWriter, r *http.Request) {
	rawProductId := chi.URLParam(r, "product_id")

	productId, err := uuid.Parse(rawProductId)
	if err != nil {

		jsonutils.EncodeJson(w, r, http.StatusBadRequest, map[string]any{
			"message": "invalid uuid",
		})
		return

	}

	_, err = api.ProductService.GetProductById(r.Context(), productId)
	if err != nil {
		if errors.Is(err, services.ErrProductNotFound) {

			jsonutils.EncodeJson(w, r, http.StatusNotFound, map[string]any{
				"message": "product not found",
			})
			return
		}
		jsonutils.EncodeJson(w, r, http.StatusInternalServerError, map[string]any{
			"message": "unexpected internal server error",
		})
		return

	}

	userId, ok := api.Sessions.Get(r.Context(), "AuthenticateUserId").(uuid.UUID)
	if !ok {
		jsonutils.EncodeJson(w, r, http.StatusInternalServerError, map[string]any{
			"message": "unexpected internal server error",
		})

	}

	conn, err := api.WsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		jsonutils.EncodeJson(w, r, http.StatusInternalServerError, map[string]any{
			"message": "unexpected internal server error",
		})
	}

	fmt.Println("here")

	api.AuctionLobby.Lock()
	room, ok := api.AuctionLobby.Rooms[productId]
	api.AuctionLobby.Unlock()

	if !ok {
		jsonutils.EncodeJson(w, r, http.StatusBadRequest, map[string]any{
			"message": "the auction has ended",
		})
	}

	client := services.NewClient(room, conn, userId)
	fmt.Println(client)

	room.Register <- client
	go client.ReadEventLoop()
	go client.WriteEventLoop()

}
