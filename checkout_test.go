package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCheckout(t *testing.T) {
	ProductList = initProductList()
	PromoList = initPromoList()

	t.Run("product not found", func(t *testing.T) {
		got, err := Checkout([]CheckoutRequest{
			{
				SKU: "not found",
				Qty: 1,
			},
		})
		assert.Error(t, err)
		assert.Equal(t, CheckoutResponse{
			Products: []Product{},
		}, got)
	})

	t.Run("qty insufficient", func(t *testing.T) {
		got, err := Checkout([]CheckoutRequest{
			{
				SKU: "43N23P",
				Qty: 100,
			},
		})
		assert.Error(t, err)
		assert.Equal(t, CheckoutResponse{
			Products: []Product{},
		}, got)
	})

	t.Run("success without promo", func(t *testing.T) {
		got, err := Checkout([]CheckoutRequest{
			{
				SKU: "43N23P",
				Qty: 1,
			},
		})
		assert.NoError(t, err)
		assert.Equal(t, CheckoutResponse{
			Products: []Product{
				{
					SKU:      "43N23P",
					Name:     "MacBook Pro",
					Price:    5399.99,
					Qty:      1,
					Discount: 0,
					Subtotal: 5399.99,
				},
				{
					SKU:      "234234",
					Name:     "Raspberry Pi B",
					Price:    0,
					Qty:      1,
					Discount: 0,
					Subtotal: 0,
				},
			},
			Total: float64(5399.99),
		}, got)
	})

	t.Run("success with promo", func(t *testing.T) {
		got, err := Checkout([]CheckoutRequest{
			{
				SKU: "43N23P",
				Qty: 1,
			},
			{
				SKU: "234234",
				Qty: 1,
			},
		})
		assert.NoError(t, err)
		assert.Equal(t, CheckoutResponse{
			Products: []Product{
				{
					SKU:      "43N23P",
					Name:     "MacBook Pro",
					Price:    5399.99,
					Qty:      1,
					Discount: 0,
					Subtotal: 5399.99,
				},
				{
					SKU:      "234234",
					Name:     "Raspberry Pi B",
					Price:    0,
					Qty:      1,
					Discount: 0,
					Subtotal: 0,
				},
			},
			Total: float64(5399.99),
		}, got)
	})

	t.Run("success with promo - alexa", func(t *testing.T) {
		got, err := Checkout([]CheckoutRequest{
			{
				SKU: "A304SD",
				Qty: 4,
			},
		})
		assert.NoError(t, err)
		assert.Equal(t, CheckoutResponse{
			Products: []Product{
				{
					SKU:      "A304SD",
					Name:     "Alexa Speaker",
					Price:    109.50,
					Qty:      4,
					Discount: 43.800000000000004,
					Subtotal: 394.2,
				},
			},
			Total: float64(394.2),
		}, got)
	})

	t.Run("success with promo - google home", func(t *testing.T) {
		got, err := Checkout([]CheckoutRequest{
			{
				SKU: "120P90",
				Qty: 3,
			},
		})
		assert.NoError(t, err)
		assert.Equal(t, CheckoutResponse{
			Products: []Product{
				{
					SKU:      "120P90",
					Name:     "Google Home",
					Price:    49.99,
					Qty:      3,
					Discount: 49.99,
					Subtotal: 99.97999999999999,
				},
			},
			Total: float64(99.97999999999999),
		}, got)
	})
}
