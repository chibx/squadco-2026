package request

// RegisterUserRequest represents customer registration
type RegisterUserRequest struct {
	FullName string `form:"full_name" json:"full_name" name:"Full Name" validate:"required,lte=40"`
	Email    string `form:"email" json:"email" name:"Email Address" validate:"required,email"`
	Password string `form:"password" json:"password" name:"Password" validate:"required,lte=30"`
}

// UserLoginRequest represents customer login
type UserLoginRequest struct {
	Email    string `form:"email" json:"email" name:"Email Address" validate:"required,email"`
	Password string `form:"password" json:"password" name:"Password" validate:"required,lte=30"`
}

// UpdateCustomerProfileRequest represents customer profile update
type UpdateUserProfileRequest struct {
	FullName string `form:"full_name" json:"full_name" name:"Full Name"`
	Phone    string `form:"phone" json:"phone" name:"Phone Number"`
}

type UpgradeToBusinessAccountRequest struct {
	BusinessName string `form:"business_name" json:"business_name"`
	Phone        string `form:"phone" json:"phone" name:"Phone Number" validate:"required,gte=10"`
	Address      string `form:"address" json:"address" validate:"required,gte=6"`
	CAC          string `form:"cac" json:"cac" validate:"required,gte=7"`
	NIN          string `form:"nin" json:"nin" validate:"required,gte=10"`
}
