package product

import (
	"context"
	"strconv"
	"time"

	"github.com/google/uuid"

	"github.com/nathancamolez-dev/go-bid/internal/validator"
)

type CreateProductReq struct {
	SellerID    uuid.UUID `json:"seller_id"`
	ProductName string    `json:"product_name"`
	Description string    `json:"description"`
	Baseprice   float64   `json:"baseprice"`
	AuctionEnd  time.Time `json:"auction_end"`
}

const minAuctionDuration = 2 * time.Hour

func isFloat64(value string) bool {
	_, err := strconv.ParseFloat(value, 64)
	return err == nil
}

func (req CreateProductReq) Valid(ctx context.Context) validator.Evaluator {
	var eval validator.Evaluator
	eval.CheckField(
		validator.NotBlank(req.ProductName),
		"product_name",
		"this field cannot be empty",
	)

	eval.CheckField(
		validator.NotBlank(req.Description),
		"description",
		"this field cannot be empty",
	)

	eval.CheckField(
		validator.MinChars(req.Description, 10) && validator.MaxChars(req.Description, 100),
		"description",
		"this field must have at least 10 characters and at most 100 characters",
	)

	eval.CheckField(
		validator.NonNegativeValue(req.Baseprice, 0),
		"baseprice",
		"must be greater than 0",
	)

	eval.CheckField(
		req.AuctionEnd.Sub(time.Now()) >= minAuctionDuration,
		"auction_end",
		"must be at least 2 hours from now",
	)

	return eval
}
