package routes

import (
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/chibx/vendor-pulse/internal/constants"
	"github.com/chibx/vendor-pulse/internal/db"
	"github.com/chibx/vendor-pulse/internal/global"
	"github.com/chibx/vendor-pulse/internal/model"
	server "github.com/chibx/vendor-pulse/internal/server_errors"
	"github.com/chibx/vendor-pulse/internal/types"
	"github.com/chibx/vendor-pulse/internal/types/request"
	"github.com/chibx/vendor-pulse/internal/types/response"
	"github.com/chibx/vendor-pulse/internal/utils"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"golang.org/x/sync/errgroup"
)

var _ = strconv.ParseInt

func sanitizeFileName(filename string) string {
	clean := filepath.Base(filename)
	clean = strings.ReplaceAll(clean, " ", "_")
	clean = strings.ReplaceAll(clean, "\\", "_")
	clean = strings.ReplaceAll(clean, "/", "_")
	clean = strings.Trim(clean, ".")
	if clean == "" {
		return "media"
	}
	return clean
}

func CreateMeal() fiber.Handler {
	err500 := fiber.NewError(fiber.StatusInternalServerError, "Couldn't create meal, please try again")
	return func(ctx fiber.Ctx) error {
		reqBody := &request.CreateMealRequest{}

		err := ctx.Bind().Body(reqBody)
		if err != nil {
			errorBags := server.ValErrToBag(err)
			return response.FromFiberError(ctx, fiber.ErrBadRequest, errorBags)
		}

		form, err := ctx.MultipartForm()
		if err != nil {
			return response.FromFiberError(ctx, err500)
		}
		files := form.File["files"]
		if len(files) == 0 {
			return response.FromFiberError(ctx, fiber.NewError(fiber.StatusBadRequest, "At least one image of the meal needs to be uploaded"))
		}

		if len(files) > 7 {
			return response.WriteResponse(ctx, fiber.StatusBadRequest, "Maximum image upload limit of 7 exceeded.")
		}

		vendorID := int64(0)
		if userCtx, ok := ctx.Locals(constants.UserCtxKey).(request.UserCtx); ok {
			vendorID = userCtx.ID
		}

		if vendorID == 0 {
			return response.FromFiberError(ctx, fiber.ErrUnauthorized)
		}

		failedIndexes := make([]int32, 0)
		uploadedPictures := make([]*model.MealPicture, 0, len(files))

		for index, fileHeader := range files {
			err := validateFileType(fileHeader, constants.AllowedImages)
			if err != nil {
				return response.WriteResponse(ctx, fiber.StatusBadRequest, "Invalid file type, only supports jpeg, png and webp.")
			}
			fileReader, err := fileHeader.Open()
			if err != nil {
				failedIndexes = append(failedIndexes, int32(index+1))
				continue
			}

			defer func() {
				_ = fileReader.Close()
			}()

			sanitized := sanitizeFileName(fileHeader.Filename)
			fileName, err := uuid.NewRandom()
			if err != nil {
				return response.FromFiberError(ctx, err500)
			}
			publicID := strings.TrimSuffix(fileName.String(), filepath.Ext(sanitized))
			uploadResult, err := global.Cloudinary.Upload.Upload(ctx.Context(), fileReader, uploader.UploadParams{
				Folder:         fmt.Sprintf("meals/%d", vendorID),
				PublicID:       publicID,
				UniqueFilename: utils.AsPointer(true),
				Overwrite:      utils.AsPointer(false),
			})
			if err != nil {
				failedIndexes = append(failedIndexes, int32(index+1))
				continue
			}

			uploadedPictures = append(uploadedPictures, &model.MealPicture{
				ImageURL:  uploadResult.SecureURL,
				PublicID:  uploadResult.PublicID,
				IsPrimary: index == 0,
			})
		}

		if len(failedIndexes) > 0 {
			for _, picture := range uploadedPictures {
				_, _ = global.Cloudinary.Upload.Destroy(ctx.Context(), uploader.DestroyParams{PublicID: picture.PublicID, ResourceType: "image"})
			}

			return response.FromFiberError(ctx, err500, &response.MediaUploadsResponse{FailedMedias: failedIndexes})
		}

		meal := &model.Meal{
			VendorID:    vendorID,
			Name:        reqBody.Name,
			Description: reqBody.Description,
			Price:       reqBody.Price.IntPart(),
			Category:    reqBody.Category,
			Status:      "active",
			Enabled:     reqBody.Enabled,
		}

		mealID, err := db.Meals().CreateMeal(ctx.Context(), meal, uploadedPictures)
		if err != nil {
			return response.FromFiberError(ctx, err500)
		}

		return response.WriteResponse(ctx, fiber.StatusCreated, "Meal created successfully", fiber.Map{"meal_id": mealID})
	}
}

func DeleteMealMedia() fiber.Handler {
	err500 := fiber.NewError(fiber.StatusInternalServerError, "Couldn't delete meal media, please try again")
	return func(ctx fiber.Ctx) error {
		reqBody := &request.DeleteMealPicturesRequest{}

		err := ctx.Bind().Body(reqBody)
		if err != nil {
			errorBags := server.ValErrToBag(err)
			return response.FromFiberError(ctx, fiber.ErrBadRequest, errorBags)
		}

		vendorID := int64(0)
		if userCtx, ok := ctx.Locals(constants.UserCtxKey).(request.UserCtx); ok {
			vendorID = userCtx.ID
		}

		if vendorID == 0 {
			return response.FromFiberError(ctx, fiber.ErrUnauthorized)
		}

		pictures, err := db.Meals().GetMealPicturesForVendor(ctx.Context(), vendorID, reqBody.MealID, reqBody.MediaIDs)
		if err != nil {
			return response.FromFiberError(ctx, err500)
		}
		if len(pictures) == 0 {
			return response.FromFiberError(ctx, fiber.NewError(fiber.StatusBadRequest, "No matching media found for deletion"))
		}

		var errG errgroup.Group

		for _, picture := range pictures {
			errG.Go(func() error {
				_, err := global.Cloudinary.Upload.Destroy(ctx.Context(), uploader.DestroyParams{PublicID: picture.PublicID, ResourceType: "image"})
				return err
			})
			// if err != nil {
			// 	return response.FromFiberError(ctx, err500)
			// }
		}

		err = errG.Wait()
		if err != nil {
			return response.FromFiberError(ctx, err500)
		}

		err = db.Meals().DeleteMealPictures(ctx.Context(), vendorID, reqBody.MealID, reqBody.MediaIDs)
		if err != nil {
			return response.FromFiberError(ctx, err500)
		}

		return response.WriteResponse(ctx, fiber.StatusOK, "Meal media deleted successfully")
	}
}

func UpdateMeal() fiber.Handler {
	err500 := fiber.NewError(fiber.StatusInternalServerError, "Couldn't update meal, please try again")
	return func(ctx fiber.Ctx) error {
		mealID := ctx.Params("id")
		if mealID == "" {
			return response.FromFiberError(ctx, fiber.ErrBadRequest)
		}

		reqBody := &request.UpdateMealRequest{}

		err := ctx.Bind().Body(reqBody)
		if err != nil {
			errorBags := server.ValErrToBag(err)
			return response.FromFiberError(ctx, fiber.ErrBadRequest, errorBags)
		}

		vendorID := int64(0)
		if userCtx, ok := ctx.Locals(constants.UserCtxKey).(request.UserCtx); ok {
			vendorID = userCtx.ID
		}
		if vendorID == 0 {
			return response.FromFiberError(ctx, fiber.ErrUnauthorized)
		}

		mealIDInt, err := strconv.ParseInt(mealID, 10, 64)
		if err != nil {
			return response.FromFiberError(ctx, fiber.ErrBadRequest)
		}

		err = db.Meals().UpdateMeal(ctx.Context(), vendorID, mealIDInt, reqBody)
		if err != nil {
			return response.FromFiberError(ctx, err500)
		}

		return response.WriteResponse(ctx, fiber.StatusOK, "Meal updated successfully")
	}
}

func DeleteMeal() fiber.Handler {
	err500 := fiber.NewError(fiber.StatusInternalServerError, "Couldn't delete meal, please try again")
	return func(ctx fiber.Ctx) error {
		mealID := ctx.Params("id")
		if mealID == "" {
			return response.FromFiberError(ctx, fiber.ErrBadRequest)
		}

		vendorID := int64(0)
		if userCtx, ok := ctx.Locals(constants.UserCtxKey).(request.UserCtx); ok {
			vendorID = userCtx.ID
		}
		if vendorID == 0 {
			return response.FromFiberError(ctx, fiber.ErrUnauthorized)
		}

		mealIDInt, err := strconv.ParseInt(mealID, 10, 64)
		if err != nil {
			return response.FromFiberError(ctx, fiber.ErrBadRequest)
		}

		pictures, err := db.Meals().GetMealPicturesByMealID(ctx.Context(), vendorID, mealIDInt)
		if err != nil {
			return response.FromFiberError(ctx, err500)
		}

		var wg sync.WaitGroup

		for _, picture := range pictures {
			wg.Go(func() {
				_, _ = global.Cloudinary.Upload.Destroy(ctx.Context(), uploader.DestroyParams{PublicID: picture.PublicID, ResourceType: "image"})
			})
		}

		wg.Wait()

		err = db.Meals().DeleteMeal(ctx.Context(), vendorID, mealIDInt)
		if err != nil {
			return response.FromFiberError(ctx, err500)
		}

		return response.WriteResponse(ctx, fiber.StatusOK, "Meal deleted successfully")
	}
}

func AddReview() fiber.Handler {
	err500 := fiber.NewError(fiber.StatusInternalServerError, "Couldn't add your review, please try again")
	return func(ctx fiber.Ctx) error {
		mealID := ctx.Params("id")
		if mealID == "" {
			return response.FromFiberError(ctx, fiber.ErrBadRequest)
		}

		reqBody := &request.AddReviewRequest{}

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

		mealIDInt, err := strconv.ParseInt(mealID, 10, 64)
		if err != nil {
			return response.FromFiberError(ctx, fiber.ErrBadRequest)
		}

		now := time.Now()

		review := &model.Review{
			CustomerID: customerID,
			VendorID:   reqBody.VendorID,
			MealID:     mealIDInt,
			Rating:     reqBody.Rating,
			Comment:    reqBody.Comment,
			Edits:      0,
			CreatedAt:  now,
			UpdatedAt:  now,
		}

		err = db.Meals().AddReview(ctx.Context(), mealIDInt, review)
		if err != nil {
			logger.Err(err).Msg("Failed to add review to database")
			return response.FromFiberError(ctx, err500)
		}

		// TODO: Integrate AI/ML SDK for sentiment analysis on review.Comment in a background job or service
		// Calculate Sentiment and SentimentScore based on comment

		return response.WriteResponse(ctx, fiber.StatusCreated, "Review added successfully")
	}
}

func EditReview() fiber.Handler {
	err500 := fiber.NewError(fiber.StatusInternalServerError, "Couldn't edit this review, please try again")
	return func(ctx fiber.Ctx) error {
		reqBody := &request.EditReviewRequest{}

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

		err = db.Meals().EditReview(ctx.Context(), reqBody.ReviewID, customerID, reqBody)
		if err != nil {
			return response.FromFiberError(ctx, err500)
		}

		return response.WriteResponse(ctx, fiber.StatusOK, "Review edited successfully")
	}
}

func ListReviews() fiber.Handler {
	err500 := fiber.NewError(fiber.StatusInternalServerError, "Couldn't load meal reviews")
	return func(ctx fiber.Ctx) error {
		mealID := ctx.Params("meal_id")
		if mealID == "" {
			return response.FromFiberError(ctx, fiber.ErrBadRequest)
		}

		mealIDInt, err := strconv.ParseInt(mealID, 10, 64)
		if err != nil {
			return response.FromFiberError(ctx, fiber.ErrBadRequest)
		}

		page := 1
		limit := 10

		if q := ctx.Query("page", ""); q != "" {
			if p, err := strconv.Atoi(q); err == nil && p > 0 {
				page = p
			}
		}
		if q := ctx.Query("limit", ""); q != "" {
			if l, err := strconv.Atoi(q); err == nil && l > 0 {
				limit = l
			}
		}

		reviews, err := db.Meals().ListReviews(ctx.Context(), mealIDInt, types.Pagination{Page: page, PageSize: limit})
		if err != nil {
			return response.FromFiberError(ctx, err500)
		}

		customerID := int64(0)
		if userCtx, ok := ctx.Locals(constants.UserCtxKey).(request.UserCtx); ok {
			customerID = userCtx.ID
		}

		responseReviews := make([]*response.MealReviewResponse, 0, len(reviews))
		for _, review := range reviews {
			isOwner := review.CustomerID != 0 && review.CustomerID == customerID
			responseReviews = append(responseReviews, &response.MealReviewResponse{
				ReviewID:  review.ID,
				IsOwner:   isOwner,
				CanEdit:   isOwner && review.Edits <= 5,
				Rating:    int(review.Rating),
				Comment:   review.Comment,
				CreatedAt: review.CreatedAt,
				UpdatedAt: review.UpdatedAt,
			})
		}

		return response.WriteResponse(ctx, fiber.StatusOK, "Meal reviews loaded successfully", &response.ListMealReviewResponse{
			Reviews: responseReviews,
			Total:   len(responseReviews),
			Page:    page,
		})
	}
}

func SearchMeals() fiber.Handler {
	err500 := fiber.NewError(fiber.StatusInternalServerError, "Failed to load meals")

	return func(ctx fiber.Ctx) error {
		search := ctx.Query("search")
		if search == "" {
			return response.FromFiberError(ctx, fiber.ErrBadRequest)
		}

		page := 1
		limit := 10

		if q := ctx.Query("page", ""); q != "" {
			if p, err := strconv.Atoi(q); err == nil && p > 0 {
				page = p
			}
		}
		if q := ctx.Query("limit", ""); q != "" {
			if l, err := strconv.Atoi(q); err == nil && l > 0 {
				limit = l
			}
		}

		meals, err := db.Meals().SearchMeals(ctx.Context(), search, types.Pagination{Page: page, PageSize: limit})
		if err != nil {
			return response.FromFiberError(ctx, err500)
		}

		return response.WriteResponse(ctx, fiber.StatusOK, "Meal reviews loaded successfully", &response.SearchMealsResponse{
			Meals: meals,
			Total: len(meals),
			Page:  page,
		})
	}
}
