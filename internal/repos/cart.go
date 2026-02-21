package repos

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/richardbowden/degrees/internal/dbpg"
	"github.com/richardbowden/degrees/internal/services"
)

type Cart struct {
	store dbpg.Storer
}

func NewCartRepo(store dbpg.Storer) *Cart {
	return &Cart{
		store: store,
	}
}

func (r *Cart) CreateCartSession(ctx context.Context, params dbpg.CreateCartSessionParams) (dbpg.CartSession, error) {
	return r.store.CreateCartSession(ctx, params)
}

func (r *Cart) GetCartBySessionToken(ctx context.Context, token string) (dbpg.CartSession, error) {
	session, err := r.store.GetCartBySessionToken(ctx, dbpg.GetCartBySessionTokenParams{SessionToken: token})
	if err != nil {
		if dbpg.IsErrNoRows(err) {
			return dbpg.CartSession{}, services.ErrNoRecord
		}
		return dbpg.CartSession{}, err
	}
	return session, nil
}

func (r *Cart) GetCartByUserID(ctx context.Context, userID int64) (dbpg.CartSession, error) {
	session, err := r.store.GetCartByUserID(ctx, dbpg.GetCartByUserIDParams{
		UserID: pgtype.Int8{Int64: userID, Valid: true},
	})
	if err != nil {
		if dbpg.IsErrNoRows(err) {
			return dbpg.CartSession{}, services.ErrNoRecord
		}
		return dbpg.CartSession{}, err
	}
	return session, nil
}

func (r *Cart) ListCartItems(ctx context.Context, cartSessionID int64) ([]dbpg.ListCartItemsRow, error) {
	return r.store.ListCartItems(ctx, dbpg.ListCartItemsParams{CartSessionID: cartSessionID})
}

func (r *Cart) AddCartItem(ctx context.Context, params dbpg.AddCartItemParams) (dbpg.CartItem, error) {
	return r.store.AddCartItem(ctx, params)
}

func (r *Cart) AddCartItemOption(ctx context.Context, params dbpg.AddCartItemOptionParams) (dbpg.CartItemOption, error) {
	return r.store.AddCartItemOption(ctx, params)
}

func (r *Cart) UpdateCartItemQuantity(ctx context.Context, params dbpg.UpdateCartItemQuantityParams) (dbpg.CartItem, error) {
	item, err := r.store.UpdateCartItemQuantity(ctx, params)
	if err != nil {
		if dbpg.IsErrNoRows(err) {
			return dbpg.CartItem{}, services.ErrNoRecord
		}
		return dbpg.CartItem{}, err
	}
	return item, nil
}

func (r *Cart) RemoveCartItem(ctx context.Context, id int64) error {
	return r.store.RemoveCartItem(ctx, dbpg.RemoveCartItemParams{ID: id})
}

func (r *Cart) ClearCart(ctx context.Context, cartSessionID int64) error {
	return r.store.ClearCart(ctx, dbpg.ClearCartParams{CartSessionID: cartSessionID})
}
