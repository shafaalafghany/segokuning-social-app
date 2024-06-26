package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5"
	"github.com/shafaalafghany/segokuning-social-app/internal/common/response"
	"github.com/shafaalafghany/segokuning-social-app/internal/common/utils/validation"
	dto "github.com/shafaalafghany/segokuning-social-app/internal/domain/dto/user"
	"go.uber.org/zap"
)

func (uh *UserHandler) Update(w http.ResponseWriter, r *http.Request) {
	var (
		userId string
		data   dto.UserUpdate
	)

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		uh.log.Info("required fields are missing or invalid", zap.Error(err))
		(&response.Response{
			HttpStatus: http.StatusBadRequest,
			Message:    "required fields are missing or invalid",
		}).GenerateResponse(w)
		return
	}

	if err := validation.UrlValidation(data.ImageUrl); err != nil {
		uh.log.Info("failed to validate url", zap.Error(err))
		(&response.Response{
			HttpStatus: http.StatusBadRequest,
			Message:    "URL malformed",
		}).GenerateResponse(w)

		return
	}

	if err := uh.val.Struct(data); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		for _, e := range validationErrors {
			uh.log.Info(validation.CustomError(e), zap.Error(err))
			(&response.Response{
				HttpStatus: http.StatusBadRequest,
				Message:    validation.CustomError(e),
			}).GenerateResponse(w)
			return
		}
	}

	ctx := r.Context()
	userId = ctx.Value("user_id").(string)

	result, err := uh.ur.FindById(ctx, userId)
	if err != nil {
		if err == pgx.ErrNoRows {
			uh.log.Info("user is not found", zap.Error(err))
			(&response.Response{
				HttpStatus: http.StatusNotFound,
				Message:    "user not found",
			}).GenerateResponse(w)
			return
		}

		uh.log.Info("failed to get user", zap.Error(err))
		(&response.Response{
			HttpStatus: http.StatusInternalServerError,
			Message:    err.Error(),
		}).GenerateResponse(w)
		return
	}

	result.Name = data.Name
	result.ImageUrl = data.ImageUrl

	if err := uh.ur.Update(ctx, *result); err != nil {
		uh.log.Info("failed to update user", zap.Error(err))
		(&response.Response{
			HttpStatus: http.StatusInternalServerError,
			Message:    err.Error(),
		}).GenerateResponse(w)
		return
	}

	(&response.Response{
		HttpStatus: http.StatusOK,
		Message:    "successfully update user profile",
	}).GenerateResponse(w)
}
