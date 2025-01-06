package common

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/Snehashish1609/couponverse-api/models"
	"github.com/rs/zerolog/log"
)

type DataResponse struct {
	Data    interface{} `json:"data,omitempty"` // data
	Message string      `json:"message"`        // error message
}

type ApplicableCouponsRequest struct {
	Cart models.Cart `json:"cart"`
}

type ApplyCouponRequest struct {
	Cart models.Cart `json:"cart"`
}

type ApplicableCouponResponse struct {
	CouponID uint    `json:"coupon_id"`
	Type     string  `json:"type"`
	Discount float64 `json:"discount"`
}

func WriteResponse(data interface{}, respWriter http.ResponseWriter, status int) {
	respWriter.Header().Set("content-type", "application/json; charset=utf-8")
	respWriter.WriteHeader(status)
	json.NewEncoder(respWriter).Encode(data)
}

func GetEnvOrDie(key string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	log.Fatal().
		Str("key", key).
		Msg("failed to get environment variable")
	return ""
}
