package routes

import (
	"strings"

	"github.com/chibx/vendor-pulse/internal/constants"
	server "github.com/chibx/vendor-pulse/internal/server_errors"
	"github.com/chibx/vendor-pulse/internal/types/request"
	"github.com/chibx/vendor-pulse/internal/types/response"
	"github.com/gofiber/fiber/v3"
)

func MakePayment() fiber.Handler {
	err500 := fiber.NewError(fiber.StatusInternalServerError, "Couldn't place order, please try again")
	return func(ctx fiber.Ctx) error {
		reqBody := &request.MakePaymentRequest{}

		err := ctx.Bind().Body(reqBody)
		if err != nil {
			errorBags := server.ValErrToBag(err)
			return response.FromFiberError(ctx, fiber.ErrBadRequest, errorBags)
		}

		if !strings.HasPrefix(reqBody.OrderID, "ord_") {
			return response.FromFiberError(ctx, fiber.ErrBadRequest)
		}

		customerID := int64(0)
		if userCtx, ok := ctx.Locals(constants.UserCtxKey).(request.UserCtx); ok {
			customerID = userCtx.ID
		}
		if customerID == 0 {
			return response.FromFiberError(ctx, fiber.ErrUnauthorized)
		}

		// orderID, err := db.Orders().CreateOrder(ctx.Context(), customerID, reqBody.VendorID, reqBody)
		if err != nil {
			return response.FromFiberError(ctx, err500)
		}

		

		return response.WriteResponse(ctx, fiber.StatusCreated, "Order placed successfully")
	}
}
