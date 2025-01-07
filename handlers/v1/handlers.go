package v1

import (
	"net/http"

	"github.com/Snehashish1609/couponverse-api/common"
	"github.com/Snehashish1609/couponverse-api/db/coupon"
	"github.com/rs/zerolog/log"
)

type CouponsHandler struct {
	DBClient coupon.Client
}

func NewCouponsHandler(client coupon.Client) *CouponsHandler {
	return &CouponsHandler{
		DBClient: client,
	}
}

func (ch *CouponsHandler) HomeHandler(w http.ResponseWriter, r *http.Request) {
	log.Info().
		Msg("handling home")

	response := http.Response{
		Status:     "api request not allowed on home",
		StatusCode: http.StatusBadRequest,
	}

	common.WriteResponse(response, w, http.StatusBadRequest)
}
