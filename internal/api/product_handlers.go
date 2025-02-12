package api

import (
	"net/http"

	"github.com/google/uuid"

	"github.com/nathancamolez-dev/go-bid/internal/jsonutils"
	"github.com/nathancamolez-dev/go-bid/internal/usecase/product"
)

func (api *Api) handleCreateProduct(w http.ResponseWriter, r *http.Request) {
	data, problems, err := jsonutils.DecodeValidJson[product.CreateProductReq](r)
	if err != nil {
		jsonutils.EncodeJson(w, r, http.StatusUnprocessableEntity, problems)
		return
	}

	userID, ok := api.Sessions.Get(r.Context(), "AuthenticateUserId").(uuid.UUID)
	if !ok {
		jsonutils.EncodeJson(w, r, http.StatusInternalServerError, map[string]any{
			"error": "Unexpected internal server error",
		})
	}

	id, err := api.ProductService.CreateProduct(r.Context(),
		userID,
		data.ProductName,
		data.Description,
		data.Baseprice,
		data.AuctionEnd,
	)
	if err != nil {
		jsonutils.EncodeJson(w, r, http.StatusInternalServerError, map[string]any{
			"error": "failed to create product",
		})
		return
	}

	_ = jsonutils.EncodeJson(w, r, http.StatusOK, map[string]any{
		"message":    "Sucessfully created product",
		"product_id": id,
	})

}
