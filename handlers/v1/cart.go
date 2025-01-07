package v1

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/Snehashish1609/couponverse-api/common"
	"github.com/Snehashish1609/couponverse-api/db"
	"github.com/Snehashish1609/couponverse-api/models"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
)

func GetApplicableCoupons(w http.ResponseWriter, r *http.Request) {
	log.Info().
		Msg("GetApplicableCoupons called")

	// body, _ := ioutil.ReadAll(r.Body)
	// fmt.Println("Body:", string(body))

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

	fmt.Println("Cart:", request.Cart)

	var coupons []models.Coupon
	if err := db.DB.Find(&coupons).Error; err != nil {
		http.Error(w, "Could not retrieve coupons", http.StatusInternalServerError)
		return
	}

	fmt.Println("got coupons:", coupons)

	applicableCoupons := []common.ApplicableCouponResponse{}

	for _, coupon := range coupons {
		fmt.Println(coupon.Type)
		var discount float64
		switch coupon.Type {
		case models.CartWise:
			fmt.Println("calculateCartWiseDiscount")
			discount = calculateCartWiseDiscount(coupon, request.Cart)
		case models.ProductWise:
			fmt.Println("calculateProductWiseDiscount")
			discount = calculateProductWiseDiscount(coupon, request.Cart)
		case models.BxGy:
			fmt.Println("calculateBxGyDiscount")
			discount = calculateBxGyDiscount(coupon, request.Cart)
		}

		fmt.Println("discount:", discount)
		if discount > 0 {
			applicableCoupons = append(applicableCoupons, common.ApplicableCouponResponse{
				CouponID: coupon.ID,
				Type:     string(coupon.Type),
				Discount: discount,
			})
		}
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"applicable_coupons": applicableCoupons,
	})
}

func calculateCartWiseDiscount(coupon models.Coupon, cart models.Cart) float64 {
	var details struct { // TODO: better way to handle details per coupon type
		Threshold float64 `json:"threshold"`
		Discount  float64 `json:"discount"`
	}
	json.Unmarshal([]byte(coupon.Details), &details)

	fmt.Println(details)

	totalCartValue := 0.0
	for _, item := range cart.Items {
		totalCartValue += float64(item.Quantity) * item.Price
	}

	fmt.Println("total cart value:", totalCartValue)

	if totalCartValue > details.Threshold {
		return totalCartValue * (details.Discount / 100)
	}

	return 0
}

func calculateProductWiseDiscount(coupon models.Coupon, cart models.Cart) float64 {
	var details struct {
		ProductID int     `json:"product_id"`
		Discount  float64 `json:"discount"`
	}
	json.Unmarshal([]byte(coupon.Details), &details)
	fmt.Println(details)

	totalDiscount := 0.0
	for _, item := range cart.Items {
		if item.ProductID == details.ProductID {
			productDiscount := float64(item.Quantity) * item.Price * (details.Discount / 100)
			totalDiscount += productDiscount
		}
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

func ApplyCoupon(w http.ResponseWriter, r *http.Request) {
	log.Info().
		Msg("ApplyCoupon called")

	couponID, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "Invalid coupon ID", http.StatusBadRequest)
		return
	}

	var coupon models.Coupon
	if err := db.DB.First(&coupon, couponID).Error; err != nil {
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
		discount = calculateCartWiseDiscount(coupon, applyCouponsRequest.Cart)
	case models.ProductWise:
		discount = calculateProductWiseDiscount(coupon, applyCouponsRequest.Cart)
	case models.BxGy:
		discount = calculateBxGyDiscount(coupon, applyCouponsRequest.Cart)
	}

	totalCartValue := 0.0
	for _, item := range applyCouponsRequest.Cart.Items {
		totalCartValue += float64(item.Quantity) * item.Price
	}

	response := map[string]interface{}{
		"updated_cart":   applyCouponsRequest.Cart,
		"total_discount": discount,
		"final_price":    totalCartValue - discount,
	}
	json.NewEncoder(w).Encode(response)
}
