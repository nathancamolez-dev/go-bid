package services

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/nathancamolez-dev/go-bid/internal/store/pgstore"
)

type BidsServices struct {
	pool    *pgxpool.Pool
	queries *pgstore.Queries
}

var ErrBidIsToLow = errors.New("the bid value is too low")

func NewBidsServices(pool *pgxpool.Pool) *BidsServices {
	return &BidsServices{
		pool:    pool,
		queries: pgstore.New(pool),
	}
}

func (bs *BidsServices) PlaceBid(
	ctx context.Context,
	product_id, bidder_id uuid.UUID,
	amount float64,
) (pgstore.Bid, error) {
	product, err := bs.queries.GetProductById(ctx, product_id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return pgstore.Bid{}, err
		}
	}

	highestBid, err := bs.queries.GetHighestBidByProductId(ctx, product_id)
	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			return pgstore.Bid{}, err
		}
	}

	if product.Baseprice > amount || highestBid.BidAmount >= amount {
		return pgstore.Bid{}, ErrBidIsToLow
	}

	highestBid, err = bs.queries.CreateBid(ctx, pgstore.CreateBidParams{
		ProductID: product_id,
		UserID:    bidder_id,
		BidAmount: amount,
	})
	if err != nil {
		return pgstore.Bid{}, err
	}

	return highestBid, nil
}
