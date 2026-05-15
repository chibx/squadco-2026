package request

import "github.com/shopspring/decimal"

// ==================== MEALS ====================

// CreateMealRequest represents creating a new meal
type CreateMealRequest struct {
	Name        string          `form:"name" json:"name" name:"Meal Name"`
	Description string          `form:"description" json:"description" name:"Description"`
	Price       decimal.Decimal `form:"price" json:"price" name:"Price"` // in kobo
	Category    string          `form:"category" json:"category" name:"Category"`
	Enabled     bool            `form:"enabled" json:"enabled" name:"Enabled"`
	// files is the form field for the image files
}

// UpdateMealRequest represents updating an existing meal
type UpdateMealRequest struct {
	Name        string          `form:"name" json:"name" name:"Meal Name" validate:"required"`
	Description string          `form:"description" json:"description" name:"Description" validate:"required,gte=20,lte=1000"`
	Price       decimal.Decimal `form:"price" json:"price" name:"Price" validate:"required"`
	Category    string          `form:"category" json:"category" name:"Category"`
	Status      string          `form:"status" json:"status" name:"Status"`
	Enabled     bool            `form:"enabled" json:"enabled" name:"Enabled"`
}

// UploadMealPicturesRequest represents uploading pictures for a meal
type UploadMealPicturesRequest struct {
	MealID    int64 `form:"meal_id" json:"meal_id" name:"Meal ID" validate:"required,gte=1"`
	IsPrimary bool  `form:"is_primary" json:"is_primary" name:"Is Primary"`
	// Pictures handled via multipart/form-data
}

// DeleteMealPicturesRequest represents a request to delete meal media
type DeleteMealPicturesRequest struct {
	MealID   int64   `json:"meal_id" name:"Meal ID" validate:"required,gte=1"`
	MediaIDs []int64 `json:"media_ids" name:"Media IDs" validate:"required,dive,gte=1"`
}

// ==================== ORDERS ====================

type AddToCartRequest struct {
	MealID   int64 `json:"meal_id" validate:"required,gte=1"`
	Quantity int16 `json:"quantity" validate:"required,gte=1"`
}

type UpdateCartRequest struct {
	CartItemID int64 `json:"item_id" validate:"required,gte=1"`
	Quantity   int16 `json:"quantity" validate:"required,gte=0"`
}

// PlaceOrderRequest represents placing a new order
type PlaceOrderRequest struct {
	VendorID        int64        `json:"vendor_id" name:"Vendor ID" validate:"required,gte=1"`
	Items           []*OrderItem `json:"items" name:"Order Items" validate:"required,gte=1"`
	DeliveryAddress string       `json:"delivery_address" name:"Delivery Address" validate:"required,gte=10"`
}

type MakePaymentRequest struct {
	OrderID string `json:"order_id" validate:"required"`
}

// OrderItem represents an item in the order
type OrderItem struct {
	MealID   int64 `json:"meal_id" name:"Meal ID" validate:"required,gte=1"`
	Quantity int   `json:"quantity" name:"Quantity" validate:"required,gte=1"`
}

type CancelOrder struct {
	OrderID string `json:"order_id" validate:"required"`
}

// UpdateOrderStatusRequest represents updating order status (vendor/admin)
type UpdateOrderStatusRequest struct {
	OrderID string `json:"order_id" validate:"required,gte=1"`
	Status  string `form:"status" json:"status" name:"Status" validate:"required"`
}

// ==================== REVIEWS ====================

// AddReviewRequest represents adding a review for an order/meal
type AddReviewRequest struct {
	VendorID int64  `json:"vendor_id" name:"Vendor ID" validate:"required,gte=1"`
	Rating   int16  `json:"rating" name:"Rating" validate:"required,gte=0,lte=5"` // 1-5
	Comment  string `json:"comment" name:"Comment"`
}

type EditReviewRequest struct {
	ReviewID int64  `json:"review_id" name:"Review ID" validate:"required,gte=1"`
	Rating   int16  `json:"rating" name:"Rating" validate:"required,gte=0,lte=5"` // 1-5
	Comment  string `json:"comment" name:"Comment"`
}

// ==================== SEARCH & DISCOVERY ====================

// SearchMealsRequest represents searching for meals
type SearchMealsRequest struct {
	Query    string `form:"query" json:"query" name:"Search Query" validate:"required"`
	Category string `form:"category" json:"category" name:"Category" validate:"required"`
	VendorID int64  `form:"vendor_id" json:"vendor_id" name:"Vendor ID" validate:"required,gte=1"`
	Page     int    `form:"page" json:"page" name:"Page" validate:"required,gte=1"`
	Limit    int    `form:"limit" json:"limit" name:"Limit" validate:"required,gte=5"`
}

// SearchVendorsRequest represents searching for vendors/restaurants
type SearchVendorsRequest struct {
	Query string `form:"query" json:"query" name:"Search Query"`
	State string `form:"state" json:"state" name:"State"`
	Page  int    `form:"page" json:"page" name:"Page" validate:"required,gte=1"`
	Limit int    `form:"limit" json:"limit" name:"Limit" validate:"required,gte=5"`
}
