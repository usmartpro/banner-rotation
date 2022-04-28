package app

import (
	"crypto/rand"
	"math"
	"math/big"
)

func OneHandBandit(banners []BannerStats) (int, error) {
	var bannerID int
	var bannerIds []int
	var totalViewCount int64
	var maxIncome float64 = -1
	for _, banner := range banners {
		totalViewCount += banner.ViewCount
	}

	for _, banner := range banners {
		bannerIncome := (float64(banner.ClickCount) / float64(banner.ViewCount)) +
			math.Sqrt((2.0*math.Log(float64(totalViewCount)))/float64(banner.ViewCount))
		if bannerIncome > maxIncome {
			maxIncome = bannerIncome
			// очищаем слайс
			bannerIds = bannerIds[:0]
			// добавляем первый баннер
			bannerIds = append(bannerIds, int(banner.BannerID))
		} else if bannerIncome == maxIncome {
			// при одинаковых доходах складываем в слайс
			bannerIds = append(bannerIds, int(banner.BannerID))
		}
	}

	if len(bannerIds) == 0 {
		return 0, ErrGetBanner
	}
	if len(bannerIds) == 1 {
		return bannerIds[0], nil
	}

	// рандомный баннер при одинаковых доходах
	// стоило ли так делать хз, реализовал для требования:
	// Перебор всех: после большого количества показов, каждый баннер должен быть показан хотя один раз
	//
	// момент номер 2: линтер заставил использовать crypto/rand вместо math/rand, пришлось городить конвертации )
	index, _ := rand.Int(rand.Reader, big.NewInt(int64(len(bannerIds))))
	bannerID = bannerIds[index.Int64()]

	return bannerID, nil
}
