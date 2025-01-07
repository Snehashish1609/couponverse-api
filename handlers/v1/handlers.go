package v1

import (
	"net/http"
	"strconv"

	"github.com/Snehashish1609/couponverse-api/common"
	"github.com/Snehashish1609/couponverse-api/db/coupon"
	"github.com/Snehashish1609/couponverse-api/models"
	"github.com/gorilla/mux"
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

func (ch *CouponsHandler) GetUserHandler(w http.ResponseWriter, r *http.Request) {
	log.Info().
		Msg("GetUserHandler called")

	vars := mux.Vars(r)
	userId, _ := strconv.Atoi(vars["user_id"])

	response := getDummyUser(userId)
	common.WriteResponse(response, w, http.StatusOK)
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

func getDummyUser(id int) models.User {
	return models.User{
		UserID:    id,
		FirstName: "foo",
		LastName:  "bar",
		Age:       100,
	}
}
