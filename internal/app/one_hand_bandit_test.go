package app

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPopularBanner(t *testing.T) {
	bannerStats := []BannerStats{
		{1, 1, 1, 10, 0},
		{2, 1, 1, 10, 10},
		{3, 1, 1, 10, 0},
	}

	expectedBannerID := 2
	var BannerID int
	require.NotPanics(t, func() {
		BannerID, _ = OneHandBandit(bannerStats)
	})
	require.Equal(t, expectedBannerID, BannerID)
}

func TestShowEveryBanners(t *testing.T) {
	bannerStats := []BannerStats{
		{1, 1, 1, 1000, 1},
		{2, 1, 1, 1000, 1},
		{3, 1, 1, 1000, 1},
	}

	var BannerID int
	var banners []int
	for i := 0; i < 20; i++ {
		BannerID, _ = OneHandBandit(bannerStats)
		banners = append(banners, BannerID)
	}

	require.Contains(t, banners, 1)
	require.Contains(t, banners, 2)
	require.Contains(t, banners, 3)
}
