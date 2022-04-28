package app

type BannerStats struct {
	BannerID      int64 `db:"banner_id"`
	SlotID        int64 `db:"slot_id"`
	SocialGroupID int64 `db:"social_group_id"`
	ViewCount     int64 `db:"view_count"`
	ClickCount    int64 `db:"click_count"`
}
