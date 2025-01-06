package models

import "gorm.io/gorm"

type CouponType string

var AllCouponType = [...]string{"cart-wise", "product-wise", "bxgy"}

var (
	CartWise    CouponType = CouponType(AllCouponType[0])
	ProductWise CouponType = CouponType(AllCouponType[1])
	BxGy        CouponType = CouponType(AllCouponType[2])
)

type Coupon struct {
	gorm.Model
	Type    CouponType `gorm:"type:varchar(100)" json:"type"`
	Details string     `gorm:"type:jsonb" json:"details"`
}

func IsValidCouponType(couponType string) bool {
	for _, item := range AllCouponType {
		if item == couponType {
			return true
		}
	}
	return false
}

// type Product struct {
// 	gorm.Model
// 	Name  string  `gorm:"not null" json:"name"`
// 	Price float64 `gorm:"not null" json:"price"`
// }

type CouponProductRelation struct {
	gorm.Model
	CouponID  string     `gorm:"not null" json:"coupon_id"`
	ProductID string     `gorm:"not null" json:"product_id"`
	Type      CouponType `gorm:"type:varchar(100)" json:"type"`
	Quantity  string     `gorm:"not null" json:"quantity"`
}

type CartItem struct {
	ProductID int     `json:"product_id"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
}

type Cart struct {
	Items []CartItem `json:"items"`
}
