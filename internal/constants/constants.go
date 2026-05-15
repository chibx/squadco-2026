package constants

import (
	"os"
	"time"

	"github.com/go-redis/redis_rate/v10"
)

const (
	APPNAME       = "CloveDelight"
	FAKE_DELIVERY = 700
)

var AllowedImages = []string{"image/jpeg", "image/png", "image/webp"}

var PublicFolder = func() string {
	if folder := os.Getenv("PUBLIC_FOLDER"); folder != "" {
		return folder
	}
	return "dist"
}()

var (
	// Token bucket configurations
	GlobalLimit = redis_rate.Limit{
		Rate:   50000, // Total requests across the entire app
		Period: time.Minute,
		Burst:  100000, // Allow large bursts (e.g., traffic spikes)
	}

	UserLimit = redis_rate.Limit{
		Rate:   100, // Per customer
		Period: time.Minute,
		Burst:  200,
	}
)

const (
	USER_ACTIVE    = "active"
	USER_SUSPENDED = "suspended"
)

// Max allowed image size in bytes i.e 5MB
const (
	MaxPasswordLimit = 30
	MaxUsernameLimit = 30
	MaxImageUpload   = 5 * 1024 * 1024
	GlobalLimitKey   = "rl_global:app"
	UserLimitKey     = "rl_user:"
)

// const BackendSessionTimeout = 30 * time.Minute
const (
	// UserAccessTkDur  = 15 * time.Minute
	UserAccessTkDur  = 2 * 24 * time.Hour
	UserRefreshTkDur = 7 * 24 * time.Hour
	DeviceIDDur      = 365 * 24 * time.Hour
	ApiKeyCtxKey     = "api_key"
	DeviceIDKey      = "device_id"
	UserCtxKey       = "user"
	UserRefreshTkKey = "refresh_token"
	UserAccessTkKey  = "access_token"
)
