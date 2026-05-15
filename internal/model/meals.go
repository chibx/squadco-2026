package model

import (
	"time"

	"github.com/shopspring/decimal"
)

// ==================== MEALS & MENU ====================

// Meal represents a meal offered by a vendor
type Meal struct {
	ID          int64     `json:"meal_id" gorm:"column:id;primaryKey"`
	VendorID    int64     `json:"vendor_id" gorm:"column:vendor_id;not null;index:idx_meals_vendor_id"`
	Name        string    `json:"name" gorm:"column:name;not null"`
	Description string    `json:"description" gorm:"column:description"`
	Price       int64     `json:"price" gorm:"column:price;not null"` // in kobo
	Enabled     bool      `json:"enabled" gorm:"column:enabled;not null;default:false"`
	Category    string    `json:"category" gorm:"column:category;not null;index:idx_meals_category"`
	Status      string    `json:"status" gorm:"column:status;not null;default:'active';index:idx_meals_status"`
	Score       float64   `json:"score" gorm:"column:score;not null;default:0"`
	CreatedAt   time.Time `json:"created_at" gorm:"column:created_at;not null"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"column:updated_at;not null"`
}

// MealPicture represents pictures for a meal (multiple per meal)
type MealPicture struct {
	ID        int64     `json:"picture_id" gorm:"column:id;primaryKey"`
	MealID    int64     `json:"meal_id" gorm:"column:meal_id;not null;index:idx_meal_pictures_meal_id"`
	ImageURL  string    `json:"image_url" gorm:"column:image_url;not null"`
	PublicID  string    `json:"-" gorm:"column:public_id;not null"`
	IsPrimary bool      `json:"is_primary" gorm:"column:is_primary;not null;default:false"`
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at;not null"`
}

// ==================== ORDERS ====================

type CartItem struct {
	ID         int64 `gorm:"primaryKey"`
	CustomerID int64
	MealID     int64
	Quantity   int16
}

// Order represents a customer order
type Order struct {
	ID              string          `json:"order_id" gorm:"column:id;primaryKey"`
	CustomerID      int64           `json:"customer_id" gorm:"column:customer_id;not null;index:idx_orders_customer_id"`
	VendorID        int64           `json:"vendor_id" gorm:"column:vendor_id;not null;index:idx_orders_vendor_id"`
	Status          string          `json:"status" gorm:"column:status;not null;default:'pending';index:idx_orders_status"`
	TotalAmount     decimal.Decimal `json:"total_amount" gorm:"column:total_amount;not null"` // in kobo
	DeliveryFee     decimal.Decimal `json:"delivery_fee" gorm:"column:delivery_fee;not null"` // in kobo
	DeliveryAddress string          `json:"delivery_address" gorm:"column:delivery_address;not null"`
	CreatedAt       time.Time       `json:"created_at" gorm:"column:created_at;not null"`
	UpdatedAt       time.Time       `json:"updated_at" gorm:"column:updated_at;not null"`
}

// OrderItem represents items in an order
type OrderItem struct {
	ID        int64           `json:"id" gorm:"column:id;primaryKey"`
	OrderID   string          `json:"order_id" gorm:"column:order_id;not null;index:idx_order_items_order_id"`
	MealID    int64           `json:"meal_id" gorm:"column:meal_id;not null;index:idx_order_items_meal_id"`
	Quantity  int             `json:"quantity" gorm:"column:quantity;not null"`
	Price     decimal.Decimal `json:"price" gorm:"column:price;not null"` // per item in kobo
	CreatedAt time.Time       `json:"created_at" gorm:"column:created_at;not null"`
}

// ==================== REVIEWS & RATINGS ====================

// Review represents customer review for a meal/order (only buyers can review)
type Review struct {
	ID             int64     `json:"id" gorm:"column:id;primaryKey"`
	CustomerID     int64     `json:"customer_id" gorm:"column:customer_id;not null;index:idx_reviews_customer_id"`
	VendorID       int64     `json:"vendor_id" gorm:"column:vendor_id;not null;index:idx_reviews_vendor_id"`
	MealID         int64     `json:"meal_id" gorm:"column:meal_id;not null;index:idx_reviews_meal_id"`
	Rating         int16     `json:"rating" gorm:"column:rating;not null;check:rating >= 1 AND rating <= 5"` // 1-5
	Comment        string    `json:"comment" gorm:"column:comment"`
	Edits          int16     `json:"edits" gorm:"column:edits;default:0"`
	Sentiment      *string   `json:"sentiment" gorm:"column:sentiment"`
	SentimentScore *float64  `json:"sentiment_score" gorm:"column:sentiment_score"`
	CreatedAt      time.Time `json:"created_at" gorm:"column:created_at;not null"`
	UpdatedAt      time.Time `json:"updated_at" gorm:"column:updated_at;not null"`
}
