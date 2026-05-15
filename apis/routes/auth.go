package routes

import (
	"errors"
	"strings"
	"time"

	"github.com/chibx/vendor-pulse/internal/auth"
	"github.com/chibx/vendor-pulse/internal/constants"
	"github.com/chibx/vendor-pulse/internal/db"
	"github.com/chibx/vendor-pulse/internal/model"
	server "github.com/chibx/vendor-pulse/internal/server_errors"
	"github.com/chibx/vendor-pulse/internal/types/request"
	"github.com/chibx/vendor-pulse/internal/types/response"
	"github.com/chibx/vendor-pulse/internal/utils"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func UserRegister() fiber.Handler {
	err500 := fiber.NewError(fiber.StatusInternalServerError, "Error occurred while creating your account, please try again")
	return func(ctx fiber.Ctx) error {
		var err error
		// var errorBag = []serverErrors.ErrorDetail{}
		userForRegister := &request.RegisterUserRequest{}

		err = ctx.Bind().Body(userForRegister)
		if err != nil {
			errorBag := server.ValErrToBag(err)
			logger.Error().Err(err).Msg("Error occured while parsing login values")
			return response.FromFiberError(ctx, fiber.ErrBadRequest, errorBag)
		}

		// now := time.Now()

		user, err := utils.UserRegisterToDBBackendUser(ctx, userForRegister)
		if err != nil {
			serverErr := new(server.ServerErr)
			if errors.As(err, &serverErr) {
				logger.Error().Msg(serverErr.Message)
				return response.WriteResponse(ctx, serverErr.Code, serverErr.Message)
			}

			logger.Error().Err(err).Msg("Error converting user to db backend user")
			return response.FromFiberError(ctx, err500)
		}

		err = db.Users().CreateUser(ctx.Context(), user)
		if err != nil {
			logger.Error().Err(err).Msg("DB Error creating backend user")
			return response.FromFiberError(ctx, err500)
		}

		// TODO: Add an ip addition feature

		return response.WriteResponse(ctx, fiber.StatusOK, "User created successfully.")
	}
}

func UserLogin() fiber.Handler {
	err500 := fiber.NewError(fiber.StatusInternalServerError, "Error occurred while signing you in, please try again")

	return func(ctx fiber.Ctx) error {
		userLogin := &request.UserLoginRequest{}
		var err error

		err = ctx.Bind().Body(userLogin)
		if err != nil {
			errorBag := server.ValErrToBag(err)
			logger.Error().Err(err).Msg("Error occured while parsing login values")
			return response.FromFiberError(ctx, fiber.ErrBadRequest, errorBag)
		}

		userDB, err := db.Users().GetUserLoginByEmail(ctx.Context(), userLogin.Email)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return response.WriteResponse(ctx, fiber.StatusUnauthorized, "Invalid username and/or password")
			}

			logger.Error().Err(err).Msg("Database error during backend login")
			return response.FromFiberError(ctx, err500)
		}

		match, err := auth.CompareRawAndHash(userLogin.Password, userDB.Password)
		if err != nil {
			logger.Error().Err(err).Msg("Error verifying backend login password")
			return response.FromFiberError(ctx, err500)
		}

		if !match {
			return response.WriteResponse(ctx, fiber.StatusUnauthorized, "Invalid username and/or password")
		}

		refreshTokenExp := time.Now().Add(constants.UserRefreshTkDur)
		accessTokenExp := time.Now().Add(constants.UserAccessTkDur)
		deviceUUID := ctx.Locals(constants.DeviceIDKey).(uuid.UUID)
		ipAddr := ctx.IP()

		refreshToken, refreshTokenHash, err := auth.CompositeRefreshToken()
		if err != nil {
			logger.Error().Err(err).Msg("Error generating composite refresh token")
			return response.FromFiberError(ctx, err500)
		}

		userSession := &model.UserSession{
			UserId:           userDB.ID,
			RefreshTokenHash: refreshTokenHash,
			LastIP:           ipAddr,
			DeviceId:         deviceUUID,
			ExpiresAt:        refreshTokenExp,
		}
		err = auth.CreateBackendSession(ctx.Context(), userSession)
		if err != nil {
			return response.FromFiberError(ctx, err500)
		}

		accessToken, err := auth.GenerateBackendAccessToken(userDB.ID)
		if err != nil {
			return response.FromFiberError(ctx, err500)
		}

		ctx.Cookie(&fiber.Cookie{
			Name:     constants.UserRefreshTkKey,
			Value:    refreshToken,
			Expires:  refreshTokenExp,
			SameSite: "Strict",
			HTTPOnly: true,
			Secure:   true,
		})

		ctx.Cookie(&fiber.Cookie{
			Name:     constants.UserAccessTkKey,
			Value:    accessToken,
			Expires:  accessTokenExp,
			SameSite: "Strict",
			HTTPOnly: true,
			Secure:   true,
		})

		return response.WriteResponse(ctx, fiber.StatusOK, "Success")
	}
}

func UserLogout() fiber.Handler {
	return func(ctx fiber.Ctx) error {
		userID := int64(0)
		if userCtx, ok := ctx.Locals(constants.UserCtxKey).(request.UserCtx); ok {
			userID = userCtx.ID
		}
		refreshToken := strings.TrimSpace(ctx.Cookies(constants.UserRefreshTkKey))
		if refreshToken != "" {
			tokenHash, err := auth.GenerateHashFromString(refreshToken, auth.DefaultHashParams)
			if err != nil {
				logger.Error().Err(err).Msg("Failed to generate refresh token hash for logout")
				// utils.ClearCookies(ctx, constants.UserAccessTkKey, constants.UserRefreshTkKey)
				// return nil
			}

			if err == nil {
				// TODO: Implement retry logic or use some package
				err = db.Users().DeleteSession(ctx, &model.UserSession{
					RefreshTokenHash: tokenHash,
					UserId:           userID,
				})
				if err != nil {
					logger.Error().Err(err).Msg("Failed to generate refresh token hash for logout")
					// utils.ClearCookies(ctx, constants.UserAccessTkKey, constants.UserRefreshTkKey)
					// return nil
				}
			}
		}
		utils.ClearCookies(ctx, constants.UserAccessTkKey, constants.UserRefreshTkKey)
		return nil
	}
}

func GetUserIdentity() fiber.Handler {
	return func(ctx fiber.Ctx) error {
		userID := int64(0)
		if userCtx, ok := ctx.Locals(constants.UserCtxKey).(request.UserCtx); ok {
			userID = userCtx.ID
		}
		if userID == 0 {
			return response.WriteResponse(ctx, fiber.StatusUnauthorized, "You aren't signed in")
		}
		return ctx.JSON(fiber.Map{
			"user_id": userID,
		})
	}
}

func UpgradeToBusinessAccount() fiber.Handler {
	err500 := fiber.NewError(fiber.StatusInternalServerError, "Error occurred while upgrading you to business mode, please try again")

	return func(ctx fiber.Ctx) error {
		if !ctx.IsMultipart() {
			return response.WriteResponse(ctx, fiber.StatusBadRequest, "Expect Multipart FormData")
		}
		reqBody := &request.UpgradeToBusinessAccountRequest{}

		err := ctx.Bind().Body(reqBody)
		if err != nil {
			errorBags := server.ValErrToBag(err)
			return response.FromFiberError(ctx, fiber.ErrBadRequest, errorBags)
		}

		customerID := int64(0)
		if userCtx, ok := ctx.Locals(constants.UserCtxKey).(request.UserCtx); ok {
			customerID = userCtx.ID
		}
		if customerID == 0 {
			return response.FromFiberError(ctx, fiber.ErrUnauthorized)
		}

		form, err := ctx.MultipartForm()
		if err != nil {
			return response.FromFiberError(ctx, err500)
		}
		fileMap, errorBag, err := allFilesExist(form, businessDocs)
		if err != nil {
			return response.FromFiberError(ctx, err500)
		}

		if len(errorBag) > 0 {
			return response.FromFiberError(ctx, fiber.ErrBadRequest, errorBag)
		}

		for _, k := range fileMap {
			if k.Size > constants.MaxImageUpload {
				return response.WriteResponse(ctx, fiber.StatusBadRequest, "Maximum Upload of 5MB exceeded")
			}
		}

		// logo := fileMap["logo"]
		// face := fileMap["face"]
		// cacDocument := fileMap["cac_document"]
		// ninDocument := fileMap["nin_document"]

		db.Clove()

		return nil
	}
}
