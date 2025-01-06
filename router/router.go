package router

import (
	v1 "github.com/Snehashish1609/couponverse-api/handlers/v1"
	"github.com/gorilla/mux"
)

func CreateRouter() (*mux.Router, error) {
	r := mux.NewRouter()

	// routes
	r.HandleFunc("/", v1.HomeHandler)
	r.HandleFunc("/api/v1/user/{user_id}", v1.GetUserHandler).Methods("GET")
	r.HandleFunc("/coupons", v1.CreateCoupon).Methods("POST")
	r.HandleFunc("/coupons/{id}", v1.GetCoupon).Methods("GET")
	r.HandleFunc("/coupons", v1.GetAllCoupons).Methods("GET")
	r.HandleFunc("/coupons/{id}", v1.UpdateCoupon).Methods("PUT")
	r.HandleFunc("/coupons/{id}", v1.DeleteCoupon).Methods("DELETE")

	r.HandleFunc("/applicable-coupons", v1.GetApplicableCoupons).Methods("POST")
	r.HandleFunc("/apply-coupon/{id}", v1.ApplyCoupon).Methods("POST")

	return r, nil
}
