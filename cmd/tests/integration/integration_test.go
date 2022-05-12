package integration

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"net/http"
	"strconv"
	"testing"

	"github.com/jackc/pgx/v4"
	"github.com/stretchr/testify/suite"
	internalhttp "github.com/usmartpro/banner-rotation/internal/server/http"
	sqlstorage "github.com/usmartpro/banner-rotation/internal/storage"
)

var apiURL = "http://banner:8000"

type GetBannerResponse struct {
	BannerID int `json:"bannerId"`
}

type IntegrationTestSuite struct {
	suite.Suite
	dsn           string
	ctx           context.Context
	conn          *pgx.Conn
	storage       *sqlstorage.Storage
	httpClient    *http.Client
	slotID        int
	banners       []int
	socialGroupID int
}

func (s *IntegrationTestSuite) SetupTest() {
	s.dsn = "postgres://postgres:postgres@postgres/banners?sslmode=disable"
	s.storage = sqlstorage.New(context.Background(), s.dsn)
	s.conn = s.initDB()
	s.httpClient = &http.Client{}
	s.slotID = s.getRandValue()
	s.banners = s.getRandSlice()
	s.socialGroupID = s.getRandValue()
	s.saveValuesInDB()
}

func (s *IntegrationTestSuite) TearDownTest() {
	s.removeValuesFromDB()
	_ = s.conn.Close(s.ctx)
	_ = s.storage.Close(s.ctx)
}

func (s *IntegrationTestSuite) initDB() *pgx.Conn {
	conn, err := pgx.Connect(context.Background(), s.dsn)
	if err != nil {
		log.Fatal(err.Error())
	}

	if err = conn.Ping(context.Background()); err != nil {
		log.Fatal(err.Error())
	}

	return conn
}

func (s *IntegrationTestSuite) TestShowAllBanners() {
	var banners []int
	var result GetBannerResponse
	var body []byte
	for i := 0; i < 20; i++ {
		req, err := http.NewRequestWithContext(context.Background(), "Get",
			apiURL+fmt.Sprintf("/banner?slotId=%d&socialGroupId=%d", s.slotID, s.socialGroupID), nil)
		resp, _ := s.httpClient.Do(req)

		s.Require().NoError(err)
		s.Require().Equal(http.StatusOK, resp.StatusCode)
		body, err = ioutil.ReadAll(resp.Body)
		s.Require().NoError(err)
		err = json.Unmarshal(body, &result)
		s.Require().NoError(err)
		banners = append(banners, result.BannerID)
		err = resp.Body.Close()
		s.Require().NoError(err)
	}

	for _, expectedBannerID := range s.banners {
		s.Contains(banners, expectedBannerID)
	}
}

func (s *IntegrationTestSuite) TestMoreShowsForPopularBanner() {
	for i := 0; i < 20; i++ {
		req, err := http.NewRequestWithContext(context.Background(), "Get",
			apiURL+fmt.Sprintf("/banner?slotId=%d&socialGroupId=%d", s.slotID, s.socialGroupID), nil)
		resp, _ := s.httpClient.Do(req)
		s.Require().NoError(err)
		s.Require().Equal(http.StatusOK, resp.StatusCode)
		_ = resp.Body.Close()
	}

	for i := 0; i < 10; i++ {
		bSS := internalhttp.BannerSlotSocialGroup{BannerID: s.banners[0], SlotID: s.slotID, SocialGroupID: s.socialGroupID}
		body, err := json.Marshal(&bSS)

		s.Require().NoError(err)

		req, err := http.NewRequestWithContext(context.Background(), "POST", apiURL+"/click", bytes.NewReader(body))
		resp, _ := s.httpClient.Do(req)
		s.Require().NoError(err)
		s.Require().Equal(http.StatusOK, resp.StatusCode)
		_ = resp.Body.Close()
	}

	var banners []int
	var result GetBannerResponse
	var body []byte
	for i := 0; i < 20; i++ {
		req, err := http.NewRequestWithContext(context.Background(), "Get",
			apiURL+fmt.Sprintf("/banner?slotId=%d&socialGroupId=%d", s.slotID, s.socialGroupID), nil)
		resp, _ := s.httpClient.Do(req)

		s.Require().NoError(err)
		s.Require().Equal(http.StatusOK, resp.StatusCode)

		body, err = ioutil.ReadAll(resp.Body)
		s.Require().NoError(err)
		_ = resp.Body.Close()
		err = json.Unmarshal(body, &result)
		s.Require().NoError(err)
		banners = append(banners, result.BannerID)
	}

	freq := make(map[int]int)
	for _, bannerID := range banners {
		freq[bannerID]++
	}

	var (
		maxFreqBannerID = 0
		maxFreqValue    = -1
	)

	for bannerID, freqValue := range freq {
		if freqValue > maxFreqValue {
			maxFreqValue = freqValue
			maxFreqBannerID = bannerID
		}
	}

	s.Require().Equal(s.banners[0], maxFreqBannerID)
}

func (s *IntegrationTestSuite) saveValuesInDB() {
	var sql string
	var err error
	for _, bannerID := range s.banners {
		sql := `INSERT INTO banners (id, description) VALUES ($1, $2)`
		_, err = s.conn.Exec(s.ctx, sql, bannerID, "description of banner "+strconv.Itoa(bannerID))
		if err != nil {
			log.Fatalln("Error save banners", err.Error())
		}
	}

	sql = `INSERT INTO slots (id, description) VALUES ($1, $2)`
	_, err = s.conn.Exec(s.ctx, sql, s.slotID, "description of slot "+strconv.Itoa(s.slotID))
	if err != nil {
		log.Fatalln("Error save slot", err.Error())
	}

	sql = `INSERT INTO social_groups (id, description) VALUES ($1, $2)`
	_, err = s.conn.Exec(s.ctx, sql, s.slotID, "description of social_group "+strconv.Itoa(s.socialGroupID))
	if err != nil {
		log.Fatalln("Error save social_group", err.Error())
	}

	for _, bannerID := range s.banners {
		err := s.storage.AddBannerToSlot(bannerID, s.slotID)
		if err != nil {
			log.Fatalln("Error add banner to slot", err.Error())
		}
	}
}

func (s *IntegrationTestSuite) getRandValue() int {
	value, _ := rand.Int(rand.Reader, big.NewInt(int64(1000)))
	return int(value.Int64()) + 1000
}

func (s *IntegrationTestSuite) getRandSlice() []int {
	return []int{s.getRandValue(), s.getRandValue(), s.getRandValue()}
}

func (s *IntegrationTestSuite) removeValuesFromDB() {
	var sql string
	var err error
	sql = `DELETE FROM banner_views WHERE slot_id=$1`
	_, err = s.conn.Exec(s.ctx, sql, s.slotID)
	if err != nil {
		log.Fatalln("Error delete banner_views", err.Error())
	}

	sql = `DELETE FROM banner_clicks WHERE slot_id=$1`
	_, err = s.conn.Exec(s.ctx, sql, s.slotID)
	if err != nil {
		log.Fatalln("Error delete banner_clicks", err.Error())
	}

	sql = `DELETE FROM banner_slot WHERE slot_id=$1`
	_, err = s.conn.Exec(s.ctx, sql, s.slotID)
	if err != nil {
		log.Fatalln("Error delete banner_slot", err.Error())
	}

	for _, bannerID := range s.banners {
		sql = `DELETE FROM banners WHERE id=$1`
		_, err = s.conn.Exec(s.ctx, sql, bannerID)
		if err != nil {
			log.Fatalln("Error delete banner", err.Error())
		}
	}

	sql = `DELETE FROM slots WHERE id=$1`
	_, err = s.conn.Exec(s.ctx, sql, s.slotID)
	if err != nil {
		log.Fatalln("Error delete slot", err.Error())
	}

	sql = `DELETE FROM social_groups WHERE id=$1`
	_, err = s.conn.Exec(s.ctx, sql, s.socialGroupID)
	if err != nil {
		log.Fatalln("Error delete social_group", err.Error())
	}
}

func TestIntegrationSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}
