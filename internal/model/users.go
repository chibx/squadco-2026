package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// User represents a customer/user in the system
type User struct {
	ID              int64     `json:"user_id" gorm:"column:id;primaryKey"`
	FullName        string    `json:"full_name" gorm:"column:full_name;not null"`
	Email           string    `json:"email" gorm:"column:email;not null;unique;index:idx_users_email"`
	PhoneNumber     string    `json:"phone" gorm:"column:phone_number;not null"`
	Password        string    `json:"password" gorm:"column:password;not null"`
	Status          string    `json:"status" gorm:"column:status;not null;default:'active'"` // "active", "suspended"
	IsEmailVerified bool      `json:"is_email_verified" gorm:"column:is_email_verified;not null;default:false"`
	IsBusiness      bool      `json:"is_business" gorm:"column:is_business;not null;default:false;index:idx_users_is_business"` // default: false
	CreatedAt       time.Time `json:"created_at" gorm:"column:created_at;not null"`
	UpdatedAt       time.Time `json:"updated_at" gorm:"column:updated_at;not null"`
}

// Vendor represents a business/vendor in the system
type Business struct {
	BusinessID   int64     `json:"business_id" gorm:"column:business_id;primaryKey"`
	Name         string    `json:"name" gorm:"column:name;not null"`
	Email        string    `json:"email" gorm:"column:email;not null"`
	Password     string    `json:"password" gorm:"column:password;not null"`
	TIN          string    `json:"tin" gorm:"column:tin;not null"`
	NIN          string    `json:"nin" gorm:"column:nin;not null;unique"`
	BusinessType string    `json:"business_type" gorm:"column:business_type;not null"` // "restaurant", "ecommerce", "food_delivery"
	RCNumber     string    `json:"rc_number" gorm:"column:rc_number;not null"`         // CAC/RC number
	Address      string    `json:"address" gorm:"column:address;not null"`
	Status       string    `json:"status" gorm:"column:status;not null"` // "pending_documents", "documents_received", "payment_pending", "processing", "approved", "restricted", "flagged", "suspended"
	VendorID     string     `json:"vendor_id" gorm:"column:vendor_id;not null;index"`
	CreatedAt    time.Time `json:"created_at" gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP"`
}

func (b Business) TableName() string {
	return "business"
}

type UserSession struct {
	UserId           int64     `json:"user_id"`
	RefreshTokenHash string    `json:"-"`
	LastIP           string    `json:"last_ip"`
	DeviceId         uuid.UUID `json:"device_id"`
	UserAgent        string    `json:"user_agent"`
	CreatedAt        time.Time `json:"created_at"`
	ExpiresAt        time.Time `json:"expires_at"`
}

type UserAddress struct {
	ID      int64  `json:"id"`
	UserID  int64  `json:"user_id"`
	Address string `json:"address"`
}

type UserWallet struct {
	UserID       int64           `json:"user_id"`
	Amount       decimal.Decimal `json:"amount"`
	LastFundedAt *time.Time      `json:"last_funded_at"`
}

// On insert do nothing (for conflict)
type UserPurchasedMeal struct {
	UserID int64 `json:"user_id"`
	MealID int64 `json:"meal_id"`
}
