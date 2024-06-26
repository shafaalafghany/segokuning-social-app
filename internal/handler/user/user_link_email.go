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

func (uh *UserHandler) LinkEmail(w http.ResponseWriter, r *http.Request) {
	var (
		userId string
		data   dto.UserEmail
	)

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		uh.log.Info("required fields are missing or invalid", zap.Error(err))
		(&response.Response{
			HttpStatus: http.StatusBadRequest,
			Message:    "required fields are missing or invalid",
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

	if err := validation.EmailValidation(data.Email); err != nil {
		uh.log.Info("failed to validate email credential", zap.Error(err))
		(&response.Response{
			HttpStatus: http.StatusBadRequest,
			Message:    err.Error(),
		}).GenerateResponse(w)
		return
	}

	ctx := r.Context()
	userId = ctx.Value("user_id").(string)

	resultEmail, err := uh.ur.FindByEmail(ctx, data.Email)
	if err != nil && err != pgx.ErrNoRows {
		uh.log.Info("failed to get email", zap.Error(err))
		(&response.Response{
			HttpStatus: http.StatusInternalServerError,
			Message:    err.Error(),
		}).GenerateResponse(w)
		return
	}

	if resultEmail != nil {
		if resultEmail.ID == userId && resultEmail.Email != "" {
			uh.log.Info("you already have an email")
			(&response.Response{
				HttpStatus: http.StatusBadRequest,
				Message:    "You already have an email",
			}).GenerateResponse(w)
			return
		}

		uh.log.Info("email already existed")
		(&response.Response{
			HttpStatus: http.StatusConflict,
			Message:    "email already existed",
		}).GenerateResponse(w)
		return
	}

	resUser, err := uh.ur.FindById(ctx, userId)
	if err != nil {
		if err == pgx.ErrNoRows {
			uh.log.Info("user is not found", zap.Error(err))
			(&response.Response{
				HttpStatus: http.StatusNotFound,
				Message:    "User not found",
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

	if resUser.Email != "" {
		uh.log.Info("cannot change email if you already have one")
		(&response.Response{
			HttpStatus: http.StatusBadRequest,
			Message:    "cannot change email if you already have one",
		}).GenerateResponse(w)
		return
	}

	resUser.Email = data.Email

	if err := uh.ur.Update(ctx, *resUser); err != nil {
		uh.log.Info("failed to update user", zap.Error(err))
		(&response.Response{
			HttpStatus: http.StatusInternalServerError,
			Message:    err.Error(),
		}).GenerateResponse(w)
		return
	}

	(&response.Response{
		HttpStatus: http.StatusOK,
		Message:    "successfully link your email to email",
	}).GenerateResponse(w)
}
