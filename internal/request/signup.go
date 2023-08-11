package request

import (
	"github.com/arsalanaa44/rate_limiter/internal/model"
)

type SignUp struct {
	model.Quota
}

func (req SignUp) Validate() error {
	return nil
	// TODO - fix this
	//validation.ValidateStruct(&req,
	//	validation.Field(&req.MinuteRateLimit, is.Int, validation.Required),
	//)
}
