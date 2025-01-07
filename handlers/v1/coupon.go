package v1

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Snehashish1609/couponverse-api/common"
	"github.com/Snehashish1609/couponverse-api/models"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
)

func (ch *CouponsHandler) GetAllCoupons(w http.ResponseWriter, r *http.Request) {
	log.Info().
		Msg("GetAllCoupons called")

	allCoupons, err := ch.DBClient.GetAllCoupons()
	if err != nil {
		log.Error().Msgf("Could not get all coupons: %s", err.Error())
		response := common.DataResponse{
			Message: err.Error(),
		}

		common.WriteResponse(response, w, http.StatusInternalServerError)
		return
	}

	common.WriteResponse(allCoupons, w, http.StatusOK)
}

func (ch *CouponsHandler) CreateCoupon(w http.ResponseWriter, r *http.Request) {
	log.Info().
		Msg("CreateCoupon called")

	var coupon *models.Coupon
	if err := json.NewDecoder(r.Body).Decode(&coupon); err != nil {
		log.Error().Msgf("Invalid coupon payload: %s", err.Error())
		response := common.DataResponse{
			Message: "Invalid coupon payload",
		}

		common.WriteResponse(response, w, http.StatusBadRequest)
		return
	}

	err := ch.DBClient.CreateCoupon(coupon)
	if err != nil {
		log.Error().Msgf("Could not create coupon: %s", err.Error())
		response := common.DataResponse{Message: err.Error()}

		common.WriteResponse(response, w, http.StatusInternalServerError)
		return
	}

	response := common.DataResponse{
		Message: "Created coupon successfully",
		Data:    coupon,
	}

	common.WriteResponse(response, w, http.StatusCreated)
}

func (ch *CouponsHandler) GetCoupon(w http.ResponseWriter, r *http.Request) {
	log.Info().
		Msg("GetCoupon called")

	vars := mux.Vars(r)
	couponId, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Error().Msgf("Invalid coupon ID")
		response := common.DataResponse{
			Message: "Invalid coupon ID",
		}

		common.WriteResponse(response, w, http.StatusBadRequest)
		return
	}

	coupon, err := ch.DBClient.GetCouponById(couponId)
	if err != nil {
		log.Error().Msgf("%s", err.Error())
		response := common.DataResponse{Message: err.Error()}

		common.WriteResponse(response, w, http.StatusNotFound)
		return
	}

	response := common.DataResponse{
		Message: "Found coupon",
		Data:    coupon,
	}

	common.WriteResponse(response, w, http.StatusOK)
}

func (ch *CouponsHandler) UpdateCoupon(w http.ResponseWriter, r *http.Request) {
	log.Info().
		Msg("UpdateCoupon called")

	var updatedCoupon *models.Coupon
	vars := mux.Vars(r)
	couponID, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Error().Msgf("Invalid coupon ID")
		response := common.DataResponse{
			Message: "Invalid coupon ID",
		}

		common.WriteResponse(response, w, http.StatusBadRequest)
		return
	}

	coupon, err := ch.DBClient.GetCouponById(couponID)
	if err != nil {
		log.Error().Msgf("%s", err.Error())
		response := common.DataResponse{Message: err.Error()}

		common.WriteResponse(response, w, http.StatusNotFound)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&updatedCoupon); err != nil {
		log.Error().Msgf("Invalid coupon payload")
		response := common.DataResponse{
			Message: "Invalid coupon payload",
		}

		common.WriteResponse(response, w, http.StatusBadRequest)
		return
	}

	coupon.Type = updatedCoupon.Type
	coupon.Details = updatedCoupon.Details
	updatedCoupon, err = ch.DBClient.UpdateCoupon(couponID, coupon)
	if err != nil {
		log.Error().Msgf("%s", err.Error())
		response := common.DataResponse{Message: err.Error()}

		common.WriteResponse(response, w, http.StatusInternalServerError)
		return
	}

	response := common.DataResponse{
		Message: "Updated coupon successfully",
		Data:    *updatedCoupon,
	}

	common.WriteResponse(response, w, http.StatusOK)
}

func (ch *CouponsHandler) DeleteCoupon(w http.ResponseWriter, r *http.Request) {
	log.Info().
		Msg("DeleteCoupon called")

	vars := mux.Vars(r)
	couponId, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Error().Msgf("Invalid coupon ID")
		response := common.DataResponse{
			Message: "Invalid coupon ID",
		}

		common.WriteResponse(response, w, http.StatusBadRequest)
		return
	}

	err = ch.DBClient.DeleteCoupon(couponId)
	if err != nil {
		log.Error().Msgf("%s", err.Error())
		response := common.DataResponse{Message: err.Error()}

		common.WriteResponse(response, w, http.StatusNotFound)
		return
	}

	response := common.DataResponse{
		Message: "Deleted coupon",
	}

	common.WriteResponse(response, w, http.StatusOK)
}
