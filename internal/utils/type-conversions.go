package utils

import (
	"context"

	"github.com/chibx/vendor-pulse/internal/auth"
	"github.com/chibx/vendor-pulse/internal/global"
	"github.com/chibx/vendor-pulse/internal/model"
	"github.com/chibx/vendor-pulse/internal/types/request"
)

func UserRegisterToDBBackendUser(c context.Context, req *request.RegisterUserRequest) (*model.User, error) {
	logger := global.Logger

	passwordHash, err := auth.GenerateHashFromString(req.Password, auth.DefaultHashParams)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to hash password for new backend user")
		return nil, err
	}

	encryptedFullname, err := auth.Encrypt(req.FullName, global.SecretKey)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to encrypt fullname for new backend user")
		return nil, err
	}
	// if req.UserName != nil {
	// 	username = *req.UserName // There is no need to encrypt the username (App specific)
	// }

	encryptedEmail, err := auth.Encrypt(req.Email, global.SecretKey)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to encrypt email for new backend user")
		return nil, err
	}

	// encryptedPhoneNumber, err := auth.Encrypt(req.Phone, global.SecretKey)
	// if err != nil {
	// 	logger.Error().Err(err).Msg("Failed to encrypt phone number for new backend user")
	// 	return nil, err
	// }

	// var countryId uint
	// if req.Country != nil {
	// 	// err = api.Deps.DB.Model(&dbModels.Country{}).Where(dbModels.Country{Code: *req.Country}).Row().Scan(&countryId)
	// 	countryId, err = db.BackendUsers().GetCountryIdByCode(ctx, *req.Country)
	// 	if err != nil {
	// 		logger.Error(fmt.Sprintf("Failed to get country ID for `%s` for new backend user", *req.Country), zap.Error(err))
	// 		if errors.Is(err, pgx.ErrNoRows) {
	// 			return nil, server.NewServerErr(fiber.StatusBadRequest, fmt.Sprintf("Record for Country %s does not exist", *req.Country))
	// 		}
	// 		return nil, err
	// 	}
	// }

	// TODO: Implement image upload

	customer := &model.User{
		FullName: encryptedFullname,
		Email:    encryptedEmail,
		// PhoneNumber:     encryptedPhoneNumber,
		Password:        passwordHash,
		Status:          "active",
		IsEmailVerified: true,
		// Image:           nil,

	}

	return customer, nil
}
