package response

import (
	"time"

	"github.com/shopspring/decimal"
)

// ==================== CUSTOMER AUTH RESPONSES ====================

// CustomerAuthResponse represents customer authentication response
type CustomerAuthResponse struct {
	CustomerID string `json:"customer_id"`
	FullName   string `json:"full_name"`
	Email      string `json:"email"`
	Token      string `json:"token"`
}

// ==================== MEAL RESPONSES ====================

// MealResponse represents a meal in responses
type MealResponse struct {
	MealID      string                 `json:"meal_id"`
	VendorID    string                 `json:"vendor_id"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Price       int64                  `json:"price"`
	Category    string                 `json:"category"`
	Status      string                 `json:"status"`
	Pictures    []*MealPictureResponse `json:"pictures"`
	CreatedAt   time.Time              `json:"created_at"`
}

// MealPictureResponse represents a meal picture
type MealPictureResponse struct {
	PictureID string `json:"picture_id"`
	ImageURL  string `json:"image_url"`
	IsPrimary bool   `json:"is_primary"`
}

type MediaUploadsResponse struct {
	FailedMedias []int32 `json:"failed_uploads,omitempty"`
}

// MealListResponse represents paginated list of meals
type MealListResponse struct {
	Meals []*MealResponse `json:"meals"`
	Total int             `json:"total"`
	Page  int             `json:"page"`
}

// ==================== CUSTOMER RESPONSES ====================

// CustomerProfileResponse represents customer profile
type CustomerProfileResponse struct {
	CustomerID string    `json:"customer_id"`
	FullName   string    `json:"full_name"`
	Email      string    `json:"email"`
	Phone      string    `json:"phone"`
	Status     string    `json:"status"`
	CreatedAt  time.Time `json:"created_at"`
}

// ==================== ORDER RESPONSES ====================

// OrderResponse represents an order
type OrderResponse struct {
	OrderID         string               `json:"order_id"`
	CustomerID      string               `json:"customer_id"`
	VendorID        string               `json:"vendor_id"`
	Status          string               `json:"status"`
	TotalAmount     int64                `json:"total_amount"`
	DeliveryFee     int64                `json:"delivery_fee"`
	DeliveryAddress string               `json:"delivery_address"`
	Items           []*OrderItemResponse `json:"items"`
	CreatedAt       time.Time            `json:"created_at"`
	UpdatedAt       time.Time            `json:"updated_at"`
}

// OrderItemResponse represents an item in an order
type OrderItemResponse struct {
	OrderItemID string `json:"order_item_id"`
	MealID      string `json:"meal_id"`
	MealName    string `json:"meal_name"`
	Quantity    int    `json:"quantity"`
	Price       int64  `json:"price"`
}

// OrderListResponse represents list of orders
type OrderListResponse struct {
	Orders []*OrderResponse `json:"orders"`
	Total  int              `json:"total"`
	Page   int              `json:"page"`
}

// ==================== REVIEW RESPONSES ====================

// MealReviewResponse represents a review
type MealReviewResponse struct {
	ReviewID  int64     `json:"review_id"`
	IsOwner   bool      `json:"is_owner"` // if the review belongs to the current logged in customer
	CanEdit   bool      `json:"can_edit"` // if the the comment belongs to the logged in customer and the model.Review.Edits is less than or equal to 5
	Rating    int       `json:"rating"`
	Comment   string    `json:"comment"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ListMealReviewResponse represents list of reviews for a meal/vendor
type ListMealReviewResponse struct {
	Reviews []*MealReviewResponse `json:"reviews"`
	Total   int                   `json:"total"`
	Page    int                   `json:"page"`
}

// ==================== VENDOR MENU RESPONSES ====================

// VendorMenuResponse represents vendor's menu with meals
type VendorMenuResponse struct {
	VendorID     string          `json:"vendor_id"`
	BusinessName string          `json:"business_name"`
	Meals        []*MealResponse `json:"meals"`
	Rating       float64         `json:"rating"`
	TotalOrders  int             `json:"total_orders"`
}

// ==================== SEARCH RESPONSES ====================

type SearchMealBody struct {
	Name         string          `json:"name"`
	Price        decimal.Decimal `json:"price"`
	BusinessName string          `json:"vendor_name"`
	AvgRating    float64         `json:"avg_rating"`
	TotalReview  int32           `json:"total_reviews"`
	TotalScore   float64         `json:"-"`
}

// SearchMealsResponse represents search results for meals
type SearchMealsResponse struct {
	Meals []*SearchMealBody `json:"meals"`
	Total int             `json:"total"`
	Page  int             `json:"page"`
}

// SearchVendorsResponse represents search results for vendors
type SearchVendorsResponse struct {
	Vendors []*VendorProfileResponse `json:"vendors"`
	Total   int                      `json:"total"`
	Page    int                      `json:"page"`
}
