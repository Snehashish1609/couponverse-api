package main

import (
	"fmt"

	"github.com/Snehashish1609/couponverse-api/config"
	"github.com/Snehashish1609/couponverse-api/db"
	"github.com/Snehashish1609/couponverse-api/db/coupon"
	v1 "github.com/Snehashish1609/couponverse-api/handlers/v1"
	"github.com/Snehashish1609/couponverse-api/router"
	"github.com/Snehashish1609/couponverse-api/server"
	"github.com/rs/zerolog/log"
)

const (
	AppName = "CouponVerse"
	Port    = "8080"
)

func main() {

	// Init global config
	config.InitConfig(AppName, fmt.Sprintf(":"+Port))
	appConfig := config.GetConfig()

	dbConnString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		appConfig.DBConfig.Host,
		appConfig.DBConfig.Port,
		appConfig.DBConfig.User,
		appConfig.DBConfig.Password,
		appConfig.DBConfig.Name,
	)

	gdb, err := db.Conn(dbConnString)
	if err != nil {
		log.Fatal().Err(err).Msg("Error connecting to DB")
	}

	CouponClient := coupon.NewClient(gdb)
	CouponHandler := v1.NewCouponsHandler(CouponClient)

	r, err := router.CreateRouter(CouponHandler)
	if err != nil {
		log.Fatal().Err(err).Msg("Error creating router")
	}

	fmt.Printf("Starting couponverse on port %s...\n", Port)
	err = server.ServeCouponVerse(r, appConfig.Port)
	if err != nil {
		log.Fatal().Err(err).Msg("Error initiating server")
	}
}
