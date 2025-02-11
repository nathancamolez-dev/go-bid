package user

import (
	"context"

	"github.com/nathancamolez-dev/go-bid/internal/validator"
)

type CreateUserReq struct {
	UserName string `json:"user_name"`
	Password string `json:"password"`
	Email    string `json:"email"`
	Bio      string `json:"bio"`
}

func (req CreateUserReq) Valid(ctx context.Context) validator.Evaluator {
	var eval validator.Evaluator

	eval.CheckField(validator.NotBlank(req.UserName), "user_name", "this field cannot be empty")
	eval.CheckField(
		validator.Matches(req.Email, validator.EmailRX),
		"email",
		"must be a valid email address",
	)
	eval.CheckField(validator.NotBlank(req.Email), "email", "this field cannot be empty")
	eval.CheckField(validator.NotBlank(req.Bio), "bio", "this field cannot be empty")
	eval.CheckField(
		validator.MinChars(req.Bio, 10) && validator.MaxChars(req.Bio, 100),
		"bio",
		"this field must have at least 10 characters and at most 100 characters",
	)
	eval.CheckField(
		validator.MinChars(req.Password, 8),
		"password",
		"this field must have at least 8 characters",
	)

	return eval
}
