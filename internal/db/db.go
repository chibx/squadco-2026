package db

import (
	"github.com/chibx/vendor-pulse/internal/global"
	"gorm.io/gorm"
)

type (
	usersRepo    struct{ db *gorm.DB }
	mealsRepo    struct{ db *gorm.DB }
	ordersRepo   struct{ db *gorm.DB }
	paymentsRepo struct{ db *gorm.DB }
	cloveRepo    struct{ db *gorm.DB }
)

var (
	usersR    *usersRepo
	mealsR    *mealsRepo
	ordersR   *ordersRepo
	paymentsR *paymentsRepo
	cloveR    *cloveRepo
)

func Users() *usersRepo {
	if usersR == nil {
		usersR = &usersRepo{global.DB}
	}
	return usersR
}

func Meals() *mealsRepo {
	if mealsR == nil {
		mealsR = &mealsRepo{global.DB}
	}
	return mealsR
}

func Orders() *ordersRepo {
	if ordersR == nil {
		ordersR = &ordersRepo{global.DB}
	}
	return ordersR
}

func Payments() *paymentsRepo {
	if paymentsR == nil {
		paymentsR = &paymentsRepo{global.DB}
	}
	return paymentsR
}

func Clove() *cloveRepo {
	if cloveR == nil {
		cloveR = &cloveRepo{global.DB}
	}
	return cloveR
}
