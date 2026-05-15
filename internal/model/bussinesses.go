package model

import (
	"time"

	"github.com/shopspring/decimal"
)

// ==================== VENDOR CORE ====================

type VendorReview struct {
	ID          int64
	VendorID    string
	BusinessID int64
	Comment     string
	RatingRatio float64
	Score       float64
}

// ==================== DOCUMENTS & KYC ====================

// VendorDocument represents documents submitted for verification
type VendorDocument struct {
	DocumentID         string     `json:"document_id" gorm:"column:document_id;primaryKey"`
	VendorID           int64      `json:"vendor_id" gorm:"column:vendor_id;not null;index"`
	DocumentType       string     `json:"document_type" gorm:"column:document_type;not null"` // "cac_certificate", "valid_id", "bank_details", "business_photo", "product_photo", "exterior_photo", "selfie_photo"
	DocumentName       string     `json:"document_name" gorm:"column:document_name;not null"`
	FilePath           string     `json:"file_path" gorm:"column:file_path;not null"`
	VerificationStatus string     `json:"verification_status" gorm:"column:verification_status;not null;default:'pending'"` // "pending", "verified", "rejected"
	VerifiedBy         *string    `json:"verified_by,omitempty" gorm:"column:verified_by"`
	VerifiedAt         *time.Time `json:"verified_at,omitempty" gorm:"column:verified_at"`
	ExpiryDate         *time.Time `json:"expiry_date,omitempty" gorm:"column:expiry_date"`
	CreatedAt          time.Time  `json:"created_at" gorm:"column:created_at;not null"`
	UpdatedAt          time.Time  `json:"updated_at" gorm:"column:updated_at;not null"`
}

// VendorKYC stores personal KYC data for vendor
type VendorKYC struct {
	KYCID            string    `json:"kyc_id" gorm:"column:kyc_id;primaryKey"`
	VendorID         int64     `json:"vendor_id" gorm:"column:vendor_id;not null"`
	NINNumber        string    `json:"nin_number" gorm:"column:nin_number;not null"`
	BankName         string    `json:"bank_name" gorm:"column:bank_name;not null"`
	AccountNumber    string    `json:"account_number" gorm:"column:account_number;not null"`
	AccountName      string    `json:"account_name" gorm:"column:account_name;not null"`
	Latitude         *float64  `json:"latitude,omitempty" gorm:"column:latitude"`
	Longitude        *float64  `json:"longitude,omitempty" gorm:"column:longitude"`
	LocationAccuracy *float64  `json:"location_accuracy,omitempty" gorm:"column:location_accuracy"`
	CreatedAt        time.Time `json:"created_at" gorm:"column:created_at;not null"`
	UpdatedAt        time.Time `json:"updated_at" gorm:"column:updated_at;not null"`
}

// ==================== VERIFICATION & SCORING ====================

// VendorVerification represents the main verification record with AI scoring
type VendorVerification struct {
	VerificationID   string     `json:"verification_id" gorm:"column:verification_id;primaryKey"`
	VendorID         int64      `json:"vendor_id" gorm:"column:vendor_id;not null;index"`
	TrustScore       float64    `json:"trust_score" gorm:"column:trust_score;not null;default:0"`     // 0-100, updated based on interactions
	InitialScore     float64    `json:"initial_score" gorm:"column:initial_score;not null;default:0"` // Initial score granted after registration
	Verdict          string     `json:"verdict" gorm:"column:verdict;not null"`                       // "approved", "restricted", "flagged"
	BreakdownJSON    string     `json:"breakdown_json" gorm:"column:breakdown_json;type:jsonb"`       // JSON-encoded breakdown
	Flags            string     `json:"flags" gorm:"column:flags;type:jsonb"`                         // JSON-encoded array
	VerdictReason    string     `json:"verdict_reason" gorm:"column:verdict_reason"`
	ValidationStatus string     `json:"validation_status" gorm:"column:validation_status;not null;default:'pending'"` // "pending", "processing", "complete", "failed"
	JobID            string     `json:"job_id,omitempty" gorm:"column:job_id"`
	ValidatedBy      *string    `json:"validated_by,omitempty" gorm:"column:validated_by"`
	ValidatedAt      *time.Time `json:"validated_at,omitempty" gorm:"column:validated_at"`
	AdminNote        *string    `json:"admin_note,omitempty" gorm:"column:admin_note"`
	AdminDecision    *string    `json:"admin_decision,omitempty" gorm:"column:admin_decision"`
	CreatedAt        time.Time  `json:"created_at" gorm:"column:created_at;not null"`
	UpdatedAt        time.Time  `json:"updated_at" gorm:"column:updated_at;not null"`
}

// VerificationJob tracks async verification processing
type VerificationJob struct {
	JobID            string    `json:"job_id" gorm:"column:job_id;primaryKey"`
	VendorID         int64     `json:"vendor_id" gorm:"column:vendor_id;not null;index"`
	Status           string    `json:"status" gorm:"column:status;not null;default:'pending'"` // "pending", "processing", "complete", "failed"
	CurrentStep      string    `json:"current_step" gorm:"column:current_step"`                // "cac_check", "nin_check", "image_analysis", "score_fusion"
	StepsJSON        string    `json:"steps_json" gorm:"column:steps_json;type:jsonb"`         // JSON-encoded step statuses
	ErrorMessage     *string   `json:"error_message,omitempty" gorm:"column:error_message"`
	Result           *string   `json:"result,omitempty" gorm:"column:result"`
	EstimatedSeconds int       `json:"estimated_seconds" gorm:"column:estimated_seconds;not null;default:0"`
	CreatedAt        time.Time `json:"created_at" gorm:"column:created_at;not null"`
	UpdatedAt        time.Time `json:"updated_at" gorm:"column:updated_at;not null"`
}

// ==================== PAYMENT & TRANSACTIONS ====================

// PaymentTransaction tracks verification fee payments
type PaymentTransaction struct {
	TransactionID       string          `json:"transaction_id" gorm:"column:transaction_id;primaryKey"`
	VendorID            int64           `json:"vendor_id" gorm:"column:vendor_id;not null;index"`
	TransactionRef      string          `json:"transaction_ref" gorm:"column:transaction_ref;not null;index"`
	Amount              decimal.Decimal `json:"amount" gorm:"column:amount;not null"`                     // in kobo
	Currency            string          `json:"currency" gorm:"column:currency;not null;default:'NGN'"`   // "NGN"
	Status              string          `json:"status" gorm:"column:status;not null;default:'pending'"`   // "pending", "success", "failed"
	TransactionType     string          `json:"transaction_type" gorm:"column:transaction_type;not null"` // "verification_fee", "payout"
	SquadCheckoutURL    *string         `json:"squad_checkout_url,omitempty" gorm:"column:squad_checkout_url"`
	SquadTransactionRef *string         `json:"squad_transaction_ref,omitempty" gorm:"column:squad_transaction_ref"`
	CreatedAt           time.Time       `json:"created_at" gorm:"column:created_at;not null"`
	UpdatedAt           time.Time       `json:"updated_at" gorm:"column:updated_at;not null"`
}

// VendorVirtualAccount stores Squad virtual account details
type VendorVirtualAccount struct {
	VirtualAccountID     string    `json:"virtual_account_id" gorm:"column:virtual_account_id;primaryKey"`
	VendorID             int64     `json:"vendor_id" gorm:"column:vendor_id;not null;index:idx_vendor_virtual_accounts_vendor_id"`
	VirtualAccountNumber string    `json:"virtual_account_number" gorm:"column:virtual_account_number;not null"`
	BankName             string    `json:"bank_name" gorm:"column:bank_name;not null"`
	CustomerIdentifier   string    `json:"customer_identifier" gorm:"column:customer_identifier;not null"`
	Status               string    `json:"status" gorm:"column:status;not null"`
	CreatedAt            time.Time `json:"created_at" gorm:"column:created_at;not null"`
	UpdatedAt            time.Time `json:"updated_at" gorm:"column:updated_at;not null"`
}

// ==================== ADMIN & USERS ====================

// AdminUser represents admin staff
type AdminUser struct {
	AdminID   int64     `json:"admin_id" gorm:"column:admin_id;primaryKey"`
	Email     string    `json:"email" gorm:"column:email;not null;unique;index:idx_admin_users_email"`
	Password  string    `json:"password" gorm:"column:password;not null"`
	Name      string    `json:"name" gorm:"column:name;not null"`
	Role      string    `json:"role" gorm:"column:role;not null"`
	Status    string    `json:"status" gorm:"column:status;not null"`
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at;not null"`
	UpdatedAt time.Time `json:"updated_at" gorm:"column:updated_at;not null"`
}

// ==================== METRICS & ANALYTICS ====================

// VendorMetrics stores aggregated metrics for a vendor
type VendorMetrics struct {
	MetricsID            string    `json:"metrics_id"`
	VendorID             int64     `json:"vendor_id"`
	TotalRatings         int       `json:"total_ratings"`
	AverageRating        float64   `json:"average_rating"`
	TotalComments        int       `json:"total_comments"`
	PositiveComments     int       `json:"positive_comments"`
	NegativeComments     int       `json:"negative_comments"`
	PositiveCommentRatio float64   `json:"positive_comment_ratio"`
	ResponseRate         float64   `json:"response_rate"`
	AverageResponseTime  float64   `json:"average_response_time"`
	TotalOrders          int       `json:"total_orders"`
	OrderFulfillmentRate float64   `json:"order_fulfillment_rate"`
	LastUpdated          time.Time `json:"last_updated"`
}

// PlatformMetrics tracks overall platform analytics
type PlatformMetrics struct {
	MetricsID           string    `json:"metrics_id" gorm:"column:metrics_id;primaryKey"`
	TotalVendors        int       `json:"total_vendors" gorm:"column:total_vendors;not null;default:0"`
	ApprovedVendors     int       `json:"approved_vendors" gorm:"column:approved_vendors;not null;default:0"`
	RestrictedVendors   int       `json:"restricted_vendors" gorm:"column:restricted_vendors;not null;default:0"`
	FlaggedVendors      int       `json:"flagged_vendors" gorm:"column:flagged_vendors;not null;default:0"`
	PendingManualReview int       `json:"pending_manual_review" gorm:"column:pending_manual_review;not null;default:0"`
	AvgTrustScore       float64   `json:"avg_trust_score" gorm:"column:avg_trust_score;not null;default:0"`
	FraudAttempts       int       `json:"fraud_attempts" gorm:"column:fraud_attempts;not null;default:0"`
	TotalOrders         int       `json:"total_orders" gorm:"column:total_orders;not null;default:0"`
	TotalRevenue        int64     `json:"total_revenue" gorm:"column:total_revenue;not null;default:0"`
	LastUpdated         time.Time `json:"last_updated" gorm:"column:last_updated;not null"`
}
