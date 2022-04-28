package internalhttp

type BannerSlot struct {
	BannerID int `json:"bannerId"`
	SlotID   int `json:"slotId"`
}

type BannerSlotSocialGroup struct {
	BannerID      int `json:"bannerId"`
	SlotID        int `json:"slotId"`
	SocialGroupID int `json:"socialGroupId"`
}

type SlotSocialGroup struct {
	SlotID        int `json:"slotId"`
	SocialGroupID int `json:"socialGroupId"`
}

type BannerResponse struct {
	BannerID int `json:"bannerId"`
}

type Error struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
}
