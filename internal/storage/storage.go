package storage

import (
	"context"
	"errors"
	"fmt"
	"os"

	pgx "github.com/jackc/pgx/v4"
	"github.com/usmartpro/banner-rotation/internal/app"
)

var (
	ErrObjectNotFound = errors.New("object not found")
	ErrGetRowsError   = errors.New("get rows error")
)

type Storage struct {
	ctx  context.Context
	conn *pgx.Conn
	dsn  string
}

func New(ctx context.Context, dsn string) *Storage {
	return &Storage{
		ctx: ctx,
		dsn: dsn,
	}
}

func (s *Storage) Connect(ctx context.Context) app.Storage {
	conn, err := pgx.Connect(ctx, s.dsn)
	if err != nil {
		if _, err := fmt.Fprintf(os.Stderr, "Error connect to database: %v\n", err); err != nil {
			return nil
		}
		os.Exit(1)
	}
	s.conn = conn
	return s
}

func (s *Storage) Close(ctx context.Context) error {
	return s.conn.Close(ctx)
}

func (s *Storage) AddBannerToSlot(bannerID, slotID int) error {
	sql := `INSERT INTO banner_slot (banner_id, slot_id) VALUES ($1, $2)`
	_, err := s.conn.Exec(
		s.ctx,
		sql,
		bannerID,
		slotID,
	)

	return err
}

func (s *Storage) DeleteBannerFromSlot(bannerID, slotID int) error {
	sql := `DELETE FROM banner_slot WHERE banner_id=$1 AND slot_id=$2`
	_, err := s.conn.Exec(
		s.ctx,
		sql,
		bannerID,
		slotID,
	)

	return err
}

func (s *Storage) ClickBanner(bannerID, slotID, socialGroupID int) error {
	sql := `INSERT INTO banner_clicks (banner_id, slot_id, social_group_id, date) VALUES ($1, $2, $3, current_timestamp)`
	_, err := s.conn.Exec(
		s.ctx,
		sql,
		bannerID,
		slotID,
		socialGroupID,
	)

	return err
}

func (s *Storage) GetBannersInfo(slotID, socialGroupID int) ([]app.BannerStats, error) {
	sql := `SELECT bs.banner_id, bs.slot_id, bv.social_group_id, 
				   count(distinct bv.id) view_count, count(distinct cl.id) click_count
			FROM banner_slot bs
			LEFT JOIN banner_views bv ON bv.slot_id = bs.slot_id AND bv.banner_id = bs.banner_id
			LEFT JOIN banner_clicks cl ON bv.slot_id = cl.slot_id AND bv.banner_id = cl.banner_id AND 
											   bv.social_group_id = cl.social_group_id
			WHERE bs.slot_id = $1 AND (bv.social_group_id = $2 OR bv.social_group_id is null)
			GROUP BY bs.banner_id, bs.slot_id, bv.social_group_id
			ORDER BY bv.social_group_id`

	rows, err := s.conn.Query(s.ctx, sql, slotID, socialGroupID)
	if err != nil {
		return nil, ErrGetRowsError
	}
	defer rows.Close()

	var banners []app.BannerStats
	for rows.Next() {
		var banner app.BannerStats

		err := rows.Scan(
			&banner.BannerID,
			&banner.SlotID,
			&banner.SocialGroupID,
			&banner.ViewCount,
			&banner.ClickCount)
		if err != nil {
			return nil, err
		}
		banners = append(banners, banner)
	}

	if len(banners) == 0 {
		return nil, ErrObjectNotFound
	}

	return banners, nil
}

func (s *Storage) IncrementBannerView(bannerID, slotID, socialGroupID int) error {
	sql := `INSERT INTO banner_views (banner_id, slot_id, social_group_id, date) VALUES ($1, $2, $3, current_timestamp)`
	_, err := s.conn.Exec(
		s.ctx,
		sql,
		bannerID,
		slotID,
		socialGroupID,
	)

	return err
}

func (s *Storage) GetRandomBanner(slotID int) (int, error) {
	sql := `SELECT banner_id FROM banner_slot WHERE slot_id = $1 ORDER BY random() LIMIT 1`

	var bannerID int
	err := s.conn.QueryRow(s.ctx, sql, slotID).Scan(
		&bannerID,
	)
	if err != nil {
		return 0, ErrObjectNotFound
	}

	return bannerID, nil
}
