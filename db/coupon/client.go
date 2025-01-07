package coupon

import (
	"github.com/Snehashish1609/couponverse-api/models"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type CouponClient struct {
	DB *gorm.DB
}

type Client interface {
	GetAllCoupons() ([]models.Coupon, error)
	CreateCoupon(coupon *models.Coupon) error
	GetCouponById(id int) (*models.Coupon, error)
	UpdateCoupon(id int, newCoupon *models.Coupon) (*models.Coupon, error)
	DeleteCoupon(id int) error
	MigrateDB() error
}

func NewClient(gdb *gorm.DB) Client {
	return &CouponClient{
		DB: gdb,
	}
}

func (c *CouponClient) MigrateDB() error {
	logger := log.With().Logger()
	logger.Info().Msg("Starting auto migration")
	err := c.DB.
		AutoMigrate(
			models.Coupon{},
		)
	if err != nil {
		logger.Err(err).Msg("Auto migration failed")
		return err
	}
	logger.Info().Msg("Auto migration complete")
	return nil
}

func (c *CouponClient) GetAllCoupons() ([]models.Coupon, error) {
	log.Debug().
		Msg("GetAllCoupons called")

	var coupons []models.Coupon
	result := c.DB.Find(&coupons)
	return coupons, result.Error
}

func (c *CouponClient) CreateCoupon(coupon *models.Coupon) error {
	log.Debug().
		Msg("CreateCoupon called")

	result := c.DB.Create(&coupon)
	return result.Error
}

func (c *CouponClient) GetCouponById(id int) (*models.Coupon, error) {
	log.Debug().
		Msg("GetCouponById called")

	var coupon *models.Coupon
	result := c.DB.First(&coupon, id)
	return coupon, result.Error
}

func (c *CouponClient) UpdateCoupon(id int, newCoupon *models.Coupon) (*models.Coupon, error) {
	log.Debug().
		Msg("UpdateCoupon called")

	var existingCoupon *models.Coupon

	result := c.DB.First(&existingCoupon, id)
	if result.Error != nil {
		return nil, result.Error
	}

	existingCoupon.Type = newCoupon.Type
	existingCoupon.Details = newCoupon.Details
	result = c.DB.Save(&existingCoupon)
	return newCoupon, result.Error
}

func (c *CouponClient) DeleteCoupon(id int) error {
	log.Debug().
		Msg("DeleteCoupon called")

	result := c.DB.Delete(&models.Coupon{}, id)
	return result.Error
}
