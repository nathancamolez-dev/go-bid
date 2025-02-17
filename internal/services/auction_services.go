package services

import (
	"context"
	"errors"
	"log/slog"
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type MessageKind int

const (
	//Requests
	PlaceBid MessageKind = iota

	// Ok/ Success
	SuccessfullyPlacedBid

	//Errors
	FailedToPlaceBid

	//Info
	AuctionFinished
	NewBidPlaced
)

type Message struct {
	Message string
	Kind    MessageKind
	UserID  uuid.UUID
	Amount  float64
}

type AuctionLobby struct {
	sync.Mutex
	Rooms map[uuid.UUID]*AuctionRoom
}

type AuctionRoom struct {
	Id         uuid.UUID
	Context    context.Context
	Broadcast  chan Message
	Unregister chan *Client
	Register   chan *Client
	Clients    map[uuid.UUID]*Client

	BidsServices BidsServices
}

func (r *AuctionRoom) registerClient(c *Client) {
	slog.Info("New user connected", "Client", c)

	r.Clients[c.UserId] = c

}

func (r *AuctionRoom) unregisterClient(c *Client) {
	slog.Info("User disconnected", "Client", c)

	delete(r.Clients, c.UserId)
}

func (r *AuctionRoom) broadcastMessage(m Message) {
	slog.Info("New messge recieved", "RoomID", r.Id, "message", m.Message, "user_id", m.UserID)
	switch m.Kind {
	case PlaceBid:
		bid, err := r.BidsServices.PlaceBid(r.Context, r.Id, m.UserID, m.Amount)
		if err != nil {
			if errors.Is(err, ErrBidIsToLow) {
				if client, ok := r.Clients[m.UserID]; ok {
					client.Send <- Message{Kind: FailedToPlaceBid, Message: ErrBidIsToLow.Error()}
				}
			}
			return
		}

		if client, ok := r.Clients[m.UserID]; ok {
			client.Send <- Message{Kind: SuccessfullyPlacedBid, Message: "Successfully placed bid"}
		}

		for id, client := range r.Clients {
			newBidMessage := Message{Kind: NewBidPlaced, Message: "A new bid has been placed", Amount: bid.BidAmount}
			if id == m.UserID {
				continue
			}
			client.Send <- newBidMessage
		}
	}
}

func (r *AuctionRoom) Run() {
	slog.Info("Room stareted", "AuctionID", r.Id)

	defer func() {
		close(r.Broadcast)
		close(r.Register)
		close(r.Unregister)
	}()

	for {
		select {
		case client := <-r.Register:
			r.registerClient(client)
		case client := <-r.Unregister:
			r.unregisterClient(client)
		case message := <-r.Broadcast:
			r.broadcastMessage(message)
		case <-r.Context.Done():
			slog.Info("Auction has ended", "auctionID", r.Id)
			for _, client := range r.Clients {
				client.Send <- Message{Kind: AuctionFinished, Message: "auction has been finished"}
			}
			return

		}
	}
}

func NewAuctionRoom(ctx context.Context, id uuid.UUID, BidsService BidsServices) *AuctionRoom {
	return &AuctionRoom{
		Id:           id,
		Broadcast:    make(chan Message),
		Register:     make(chan *Client),
		Unregister:   make(chan *Client),
		Clients:      make(map[uuid.UUID]*Client),
		Context:      ctx,
		BidsServices: BidsService,
	}
}

type Client struct {
	Room   *AuctionRoom
	Conn   *websocket.Conn
	Send   chan Message
	UserId uuid.UUID
}

func NewClient(room *AuctionRoom, conn *websocket.Conn, userId uuid.UUID) *Client {
	return &Client{
		Room:   room,
		Conn:   conn,
		Send:   make(chan Message, 512),
		UserId: userId,
	}
}
