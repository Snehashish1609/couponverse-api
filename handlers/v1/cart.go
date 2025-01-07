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

func (ch *CouponsHandler) GetApplicableCoupons(w http.ResponseWriter, r *http.Request) {
	log.Info().
		Msg("GetApplicableCoupons called")

	var request common.ApplicableCouponsRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		log.Error().Msgf("Invalid cart payload: %s", err.Error())
		response := common.DataResponse{
			Message: "Invalid cart payload",
		}

		common.WriteResponse(response, w, http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	allCoupons, err := ch.DBClient.GetAllCoupons()
	if err != nil {
		http.Error(w, "Could not retrieve coupons", http.StatusInternalServerError)
		return
	}

	applicableCoupons := []common.ApplicableCouponResponse{}

	for _, coupon := range allCoupons {
		var discount float64
		switch coupon.Type {
		case models.CartWise:
			discount = calculateCartWiseDiscount(coupon, request.Cart)
		case models.ProductWise:
			discount = calculateProductWiseDiscount(coupon, request.Cart)
		case models.BxGy:
			discount = calculateBxGyDiscount(coupon, request.Cart)
		}

		if discount > 0 {
			applicableCoupons = append(applicableCoupons, common.ApplicableCouponResponse{
				CouponID: coupon.ID,
				Type:     string(coupon.Type),
				Discount: discount,
			})
		}
	}

	// can be separated out
	response := map[string]interface{}{
		"applicable_coupons": applicableCoupons,
	}
	common.WriteResponse(response, w, http.StatusOK)
}

func calculateCartWiseDiscount(coupon models.Coupon, cart models.Cart) float64 {
	var details struct { // TODO: better way to handle details per coupon type
		Threshold float64  `json:"threshold"`
		Discount  float64  `json:"discount"`
		Cap       *float64 `json:"cap"`
	}
	json.Unmarshal([]byte(coupon.Details), &details)

	totalCartValue := 0.0
	for _, item := range cart.Items {
		totalCartValue += float64(item.Quantity) * item.Price
	}

	if totalCartValue > details.Threshold {
		if details.Cap != nil {
			return min(totalCartValue*(details.Discount/100), *details.Cap)
		}
		return totalCartValue * (details.Discount / 100)
	}

	return 0
}

func calculateProductWiseDiscount(coupon models.Coupon, cart models.Cart) float64 {
	var details struct {
		ProductID int      `json:"product_id"`
		Discount  float64  `json:"discount"`
		Cap       *float64 `json:"cap"`
	}
	json.Unmarshal([]byte(coupon.Details), &details)

	totalDiscount := 0.0
	for _, item := range cart.Items {
		if item.ProductID == details.ProductID {
			productDiscount := float64(item.Quantity) * item.Price * (details.Discount / 100)
			totalDiscount += productDiscount
		}
	}
	if details.Cap != nil {
		return min(totalDiscount, *details.Cap)
	}
	return totalDiscount
}

func calculateBxGyDiscount(coupon models.Coupon, cart models.Cart) float64 {
	var details struct {
		BuyProducts []struct {
			ProductID int `json:"product_id"`
			Quantity  int `json:"quantity"`
		} `json:"buy_products"`
		GetProducts []struct {
			ProductID int `json:"product_id"`
			Quantity  int `json:"quantity"`
		} `json:"get_products"`
		RepetitionLimit int `json:"repition_limit"`
	}
	json.Unmarshal([]byte(coupon.Details), &details)

	// map product and quantity eligible for coupon
	buyCounts := make(map[int]int)
	for _, item := range cart.Items {
		for _, buy := range details.BuyProducts {
			if item.ProductID == buy.ProductID {
				buyCounts[buy.ProductID] += item.Quantity
			}
		}
	}

	appliedCouponCount := 0
	productExhausted := false
	// find eligible free items from cart
	freeItemMap := make(map[int]int)
	for appliedCouponCount < details.RepetitionLimit && !productExhausted {
		eligibleForFreeItem := false // set base condition
		for _, buyProd := range details.BuyProducts {
			if buyCounts[buyProd.ProductID] >= buyProd.Quantity {
				eligibleForFreeItem = true
				continue // keep checking for match
			} else {
				productExhausted = true     // not enough product quantities to apply coupon
				eligibleForFreeItem = false // reset if match not found
			}
		}
		if eligibleForFreeItem {
			for _, buyProd := range details.BuyProducts {
				buyCounts[buyProd.ProductID] -= buyProd.Quantity
			}
			for _, getProd := range details.GetProducts {
				freeItemMap[getProd.ProductID] += getProd.Quantity
			}
			appliedCouponCount++
		}
	}

	// calculate free items cost
	discount := 0.0
	for _, item := range cart.Items {
		if freeItemMap[item.ProductID] > 0 {
			discount = item.Price * float64(freeItemMap[item.ProductID])
		}
	}

	return discount
}

func (ch *CouponsHandler) ApplyCoupon(w http.ResponseWriter, r *http.Request) {
	log.Info().
		Msg("ApplyCoupon called")

	couponID, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "Invalid coupon ID", http.StatusBadRequest)
		return
	}

	coupon, err := ch.DBClient.GetCouponById(couponID)
	if err != nil {
		http.Error(w, "Coupon not found", http.StatusNotFound)
		return
	}

	var applyCouponsRequest common.ApplyCouponRequest
	if err := json.NewDecoder(r.Body).Decode(&applyCouponsRequest); err != nil {
		http.Error(w, "Invalid cart input", http.StatusBadRequest)
		return
	}

	// applying the coupon
	var discount float64
	switch coupon.Type {
	case models.CartWise:
		discount = calculateCartWiseDiscount(*coupon, applyCouponsRequest.Cart)
	case models.ProductWise:
		discount = calculateProductWiseDiscount(*coupon, applyCouponsRequest.Cart)
	case models.BxGy:
		discount = calculateBxGyDiscount(*coupon, applyCouponsRequest.Cart)
	}

	totalCartValue := 0.0
	for _, item := range applyCouponsRequest.Cart.Items {
		totalCartValue += float64(item.Quantity) * item.Price
	}

	// can be separated out
	response := map[string]interface{}{
		"updated_cart":   applyCouponsRequest.Cart,
		"total_discount": discount,
		"final_price":    totalCartValue - discount,
	}
	common.WriteResponse(response, w, http.StatusAccepted)
}
