package v1

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Snehashish1609/couponverse-api/common"
	"github.com/Snehashish1609/couponverse-api/db"
	"github.com/Snehashish1609/couponverse-api/models"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
)

func GetAllCoupons(w http.ResponseWriter, r *http.Request) {
	log.Info().
		Msg("GetAllCoupons called")

	var coupons []models.Coupon
	result := db.DB.Find(&coupons)
	if result.Error != nil {
		log.Error().Msgf("%s", result.Error.Error())
		response := common.DataResponse{Message: result.Error.Error()}

		common.WriteResponse(response, w, http.StatusInternalServerError)
		return
	}

	common.WriteResponse(coupons, w, http.StatusOK)
}

func CreateCoupon(w http.ResponseWriter, r *http.Request) {
	log.Info().
		Msg("CreateCoupon called")

	var coupon models.Coupon
	if err := json.NewDecoder(r.Body).Decode(&coupon); err != nil {
		log.Error().Msgf("Invalid coupon payload: %s", err.Error())
		response := common.DataResponse{
			Message: "Invalid coupon payload",
		}

		common.WriteResponse(response, w, http.StatusBadRequest)
		return
	}

	result := db.DB.Create(&coupon)
	if result.Error != nil {
		log.Error().Msgf("%s", result.Error.Error())
		response := common.DataResponse{Message: result.Error.Error()}

		common.WriteResponse(response, w, http.StatusInternalServerError)
		return
	}

	response := common.DataResponse{
		Message: "Created coupon successfully",
		Data:    coupon,
	}

	common.WriteResponse(response, w, http.StatusCreated)
}

func GetCoupon(w http.ResponseWriter, r *http.Request) {
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

	var coupon models.Coupon

	result := db.DB.First(&coupon, couponId)
	if result.Error != nil {
		log.Error().Msgf("%s", result.Error.Error())
		response := common.DataResponse{Message: result.Error.Error()}

		common.WriteResponse(response, w, http.StatusNotFound)
		return
	}

	response := common.DataResponse{
		Message: "Found coupon",
		Data:    coupon,
	}

	common.WriteResponse(response, w, http.StatusOK)
}

func UpdateCoupon(w http.ResponseWriter, r *http.Request) {
	log.Info().
		Msg("UpdateCoupon called")

	var coupon, updatedCoupon models.Coupon
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

	result := db.DB.First(&coupon, couponID)
	if result.Error != nil {
		log.Error().Msgf("%s", result.Error.Error())
		response := common.DataResponse{Message: result.Error.Error()}

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
	result = db.DB.Save(&coupon)
	if result.Error != nil {
		log.Error().Msgf("%s", result.Error.Error())
		response := common.DataResponse{Message: result.Error.Error()}

		common.WriteResponse(response, w, http.StatusInternalServerError)
		return
	}

	response := common.DataResponse{
		Message: "Updated coupon successfully",
		Data:    coupon,
	}

	common.WriteResponse(response, w, http.StatusOK)
}

func DeleteCoupon(w http.ResponseWriter, r *http.Request) {
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

	result := db.DB.Delete(&models.Coupon{}, couponId)
	if result.Error != nil {
		log.Error().Msgf("%s", result.Error.Error())
		response := common.DataResponse{Message: result.Error.Error()}

		common.WriteResponse(response, w, http.StatusNotFound)
		return
	}

	response := common.DataResponse{
		Message: "Deleted coupon",
		Data:    models.Coupon{},
	}

	common.WriteResponse(response, w, http.StatusOK)
}
