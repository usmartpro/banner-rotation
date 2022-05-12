package app

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPopularBanner(t *testing.T) {
	bannerStats := []BannerStats{
		{BannerID: 1, SlotID: 1, SocialGroupID: 1, ViewCount: 10, ClickCount: 0},
		{BannerID: 2, SlotID: 1, SocialGroupID: 1, ViewCount: 10, ClickCount: 10},
		{BannerID: 3, SlotID: 1, SocialGroupID: 1, ViewCount: 10, ClickCount: 0},
	}

	expectedBannerID := 2
	BannerID, err := MultiArmedBandit(bannerStats)
	require.NoError(t, err)
	require.Equal(t, expectedBannerID, BannerID)
}

func TestShowEveryBanners(t *testing.T) {
	bannerStats := []BannerStats{
		{BannerID: 1, SlotID: 1, SocialGroupID: 1, ViewCount: 1000, ClickCount: 1},
		{BannerID: 2, SlotID: 1, SocialGroupID: 1, ViewCount: 1000, ClickCount: 1},
		{BannerID: 3, SlotID: 1, SocialGroupID: 1, ViewCount: 1000, ClickCount: 1},
	}

	var BannerID int
	var banners []int
	var err error
	for i := 0; i < 20; i++ {
		BannerID, err = MultiArmedBandit(bannerStats)
		require.NoError(t, err)
		banners = append(banners, BannerID)
	}

	require.Contains(t, banners, 1)
	require.Contains(t, banners, 2)
	require.Contains(t, banners, 3)
}
