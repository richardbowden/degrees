package services

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/richardbowden/degrees/internal/dbpg"
	"github.com/richardbowden/degrees/internal/problems"
)

type CartRepository interface {
	CreateCartSession(ctx context.Context, params dbpg.CreateCartSessionParams) (dbpg.CartSession, error)
	GetCartBySessionToken(ctx context.Context, token string) (dbpg.CartSession, error)
	GetCartByUserID(ctx context.Context, userID int64) (dbpg.CartSession, error)
	ListCartItems(ctx context.Context, cartSessionID int64) ([]dbpg.ListCartItemsRow, error)
	AddCartItem(ctx context.Context, params dbpg.AddCartItemParams) (dbpg.CartItem, error)
	AddCartItemOption(ctx context.Context, params dbpg.AddCartItemOptionParams) (dbpg.CartItemOption, error)
	UpdateCartItemQuantity(ctx context.Context, params dbpg.UpdateCartItemQuantityParams) (dbpg.CartItem, error)
	RemoveCartItem(ctx context.Context, id int64) error
	ClearCart(ctx context.Context, cartSessionID int64) error
}

type CartService struct {
	repo CartRepository
}

func NewCartService(repo CartRepository) *CartService {
	return &CartService{
		repo: repo,
	}
}

// CartResult holds the cart session and its items for building the response.
type CartResult struct {
	Session  dbpg.CartSession
	Items    []dbpg.ListCartItemsRow
	Subtotal int64
}

// GetOrCreateCart retrieves an existing cart or creates a new one.
// For authenticated users, it looks up by userID.
// For guests, it looks up by sessionToken. If neither exists, creates a new guest session.
func (s *CartService) GetOrCreateCart(ctx context.Context, userID int64, sessionToken string) (*CartResult, error) {
	var session dbpg.CartSession
	var err error

	if userID > 0 {
		session, err = s.repo.GetCartByUserID(ctx, userID)
	} else if sessionToken != "" {
		session, err = s.repo.GetCartBySessionToken(ctx, sessionToken)
	}

	if err != nil && !errors.Is(err, ErrNoRecord) {
		return nil, problems.New(problems.Database, "failed to get cart", err)
	}

	// Create a new session if none found
	if errors.Is(err, ErrNoRecord) || session.ID == 0 {
		token, genErr := generateSessionToken()
		if genErr != nil {
			return nil, problems.New(problems.Internal, "failed to generate session token", genErr)
		}

		params := dbpg.CreateCartSessionParams{
			SessionToken: token,
			ExpiresAt:    pgtype.Timestamptz{Time: time.Now().Add(7 * 24 * time.Hour), Valid: true},
		}
		if userID > 0 {
			params.UserID = pgtype.Int8{Int64: userID, Valid: true}
		}

		session, err = s.repo.CreateCartSession(ctx, params)
		if err != nil {
			return nil, problems.New(problems.Database, "failed to create cart session", err)
		}
	}

	items, err := s.repo.ListCartItems(ctx, session.ID)
	if err != nil {
		return nil, problems.New(problems.Database, "failed to list cart items", err)
	}

	subtotal := calculateSubtotal(items)

	return &CartResult{
		Session:  session,
		Items:    items,
		Subtotal: subtotal,
	}, nil
}

// AddItem adds a service to the cart with optional service options.
func (s *CartService) AddItem(ctx context.Context, userID int64, sessionToken string, serviceID int64, vehicleID int64, quantity int32, optionIDs []int64) (*CartResult, error) {
	cart, err := s.GetOrCreateCart(ctx, userID, sessionToken)
	if err != nil {
		return nil, err
	}

	addParams := dbpg.AddCartItemParams{
		CartSessionID: cart.Session.ID,
		ServiceID:     serviceID,
		Quantity:      quantity,
	}
	if vehicleID > 0 {
		addParams.VehicleID = pgtype.Int8{Int64: vehicleID, Valid: true}
	}

	item, err := s.repo.AddCartItem(ctx, addParams)
	if err != nil {
		return nil, problems.New(problems.Database, "failed to add cart item", err)
	}

	for _, optID := range optionIDs {
		_, err := s.repo.AddCartItemOption(ctx, dbpg.AddCartItemOptionParams{
			CartItemID:      item.ID,
			ServiceOptionID: optID,
		})
		if err != nil {
			return nil, problems.New(problems.Database, "failed to add cart item option", err)
		}
	}

	// Re-fetch items to get updated totals
	items, err := s.repo.ListCartItems(ctx, cart.Session.ID)
	if err != nil {
		return nil, problems.New(problems.Database, "failed to list cart items", err)
	}

	return &CartResult{
		Session:  cart.Session,
		Items:    items,
		Subtotal: calculateSubtotal(items),
	}, nil
}

// UpdateItemQuantity updates the quantity of a cart item.
func (s *CartService) UpdateItemQuantity(ctx context.Context, userID int64, sessionToken string, itemID int64, quantity int32) (*CartResult, error) {
	cart, err := s.GetOrCreateCart(ctx, userID, sessionToken)
	if err != nil {
		return nil, err
	}

	_, err = s.repo.UpdateCartItemQuantity(ctx, dbpg.UpdateCartItemQuantityParams{
		ID:       itemID,
		Quantity: quantity,
	})
	if err != nil {
		if errors.Is(err, ErrNoRecord) {
			return nil, problems.New(problems.NotExist, "cart item not found")
		}
		return nil, problems.New(problems.Database, "failed to update cart item", err)
	}

	items, err := s.repo.ListCartItems(ctx, cart.Session.ID)
	if err != nil {
		return nil, problems.New(problems.Database, "failed to list cart items", err)
	}

	return &CartResult{
		Session:  cart.Session,
		Items:    items,
		Subtotal: calculateSubtotal(items),
	}, nil
}

// RemoveItem removes a cart item.
func (s *CartService) RemoveItem(ctx context.Context, userID int64, sessionToken string, itemID int64) (*CartResult, error) {
	cart, err := s.GetOrCreateCart(ctx, userID, sessionToken)
	if err != nil {
		return nil, err
	}

	err = s.repo.RemoveCartItem(ctx, itemID)
	if err != nil {
		return nil, problems.New(problems.Database, "failed to remove cart item", err)
	}

	items, err := s.repo.ListCartItems(ctx, cart.Session.ID)
	if err != nil {
		return nil, problems.New(problems.Database, "failed to list cart items", err)
	}

	return &CartResult{
		Session:  cart.Session,
		Items:    items,
		Subtotal: calculateSubtotal(items),
	}, nil
}

// ClearCart removes all items from the cart.
func (s *CartService) ClearCart(ctx context.Context, userID int64, sessionToken string) error {
	cart, err := s.GetOrCreateCart(ctx, userID, sessionToken)
	if err != nil {
		return err
	}

	err = s.repo.ClearCart(ctx, cart.Session.ID)
	if err != nil {
		return problems.New(problems.Database, "failed to clear cart", err)
	}
	return nil
}

func calculateSubtotal(items []dbpg.ListCartItemsRow) int64 {
	var subtotal int64
	for _, item := range items {
		subtotal += item.ServicePrice * int64(item.Quantity)
	}
	return subtotal
}

func generateSessionToken() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
