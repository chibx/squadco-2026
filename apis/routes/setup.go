package routes

import (
	"fmt"
	"io/fs"
	"mime/multipart"
	"net/http"
	"slices"

	"github.com/chibx/vendor-pulse/apis/middlewares"
	"github.com/chibx/vendor-pulse/internal/global"
	"github.com/chibx/vendor-pulse/output"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/static"
)

var (
	logger              = global.Logger
	indexPageSize int64 = 0
)

func validateFileType(file *multipart.FileHeader, allowed []string) error {
	f, err := file.Open()
	if err != nil {
		return err
	}
	defer func() {
		_ = f.Close()
	}()

	// Read the first 512 bytes for MIME detection
	buf := make([]byte, 512)
	n, err := f.Read(buf)
	if err != nil {
		return err
	}

	mimeType := http.DetectContentType(buf[:n])

	if slices.Contains(allowed, mimeType) {
		return nil
	}

	return fmt.Errorf("file type %s is not allowed", mimeType)
}

func AddRoutes(app *fiber.App) {
	api := app.Group("/api", middlewares.EnsureDeviceIDIsSet, middlewares.AuthMiddleware)
	auth := api.Group("/auth")
	auth.Get("/get-identity", GetUserIdentity())
	auth.Post("/signup", UserRegister())
	auth.Post("/login", UserLogin())
	auth.Post("/logout", UserLogout())
	api.Post("/promote-account", UpgradeToBusinessAccount())

	// -------------------MEAL------------------------------
	api.Post("/meal", CreateMeal())
	api.Post("/meal/delete-media", DeleteMealMedia())
	api.Put("/meal/:id", UpdateMeal())
	api.Delete("/meal/:id", DeleteMeal())
	api.Get("/meal/search", SearchMeals())

	// -------------------ORDER-------------------------------
	api.Post("/place-order", PlaceOrder())
	api.Post("/cancel-order", CancelOrder())

	// --------------------REVIEWS----------------------------
	api.Post("/meal/:id/review", AddReview())
	api.Post("/review/edit", EditReview())
	api.Get("/meal/:id/review", ListReviews())

	sub, _ := fs.Sub(output.EmbedFS, "dist")
	app.Get("/*", static.New("", static.Config{
		FS: sub,
	}))

	app.Get("*", func(ctx fiber.Ctx) error {
		ctx.Set(fiber.HeaderContentType, "text/html")

		fileBytes, err := sub.Open("index.html")
		if err != nil {
			logger.Err(err).Any("fs", output.EmbedFS).Msg("")
			return fiber.ErrNotFound
		}

		if indexPageSize == 0 {
			stats, err := fileBytes.Stat()
			if err != nil {
				logger.Err(err).Msg("Could not read file stats")
				return fiber.ErrNotFound
			}
			indexPageSize = stats.Size()
		}
		b := make([]byte, indexPageSize)
		_, err = fileBytes.Read(b)
		if err != nil {
			logger.Err(err).Msg("Failed to read index.html bytes")
			return fiber.ErrNotFound
		}

		return ctx.Send(b)
	})
}
