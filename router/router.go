package router

import (
	v1 "github.com/Snehashish1609/couponverse-api/handlers/v1"
	"github.com/gorilla/mux"
)

func CreateRouter(ch *v1.CouponsHandler) (*mux.Router, error) {
	r := mux.NewRouter()

	// routes
	r.HandleFunc("/", ch.HomeHandler) // blocked for now
	r.HandleFunc("/coupons", ch.CreateCoupon).Methods("POST")
	r.HandleFunc("/coupons/{id}", ch.GetCoupon).Methods("GET")
	r.HandleFunc("/coupons", ch.GetAllCoupons).Methods("GET")
	r.HandleFunc("/coupons/{id}", ch.UpdateCoupon).Methods("PUT")
	r.HandleFunc("/coupons/{id}", ch.DeleteCoupon).Methods("DELETE")

	r.HandleFunc("/applicable-coupons", ch.GetApplicableCoupons).Methods("POST")
	r.HandleFunc("/apply-coupon/{id}", ch.ApplyCoupon).Methods("POST")

	return r, nil
}
