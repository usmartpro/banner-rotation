package app

import (
	"errors"
	"fmt"
	"time"
)

type App struct {
	Logger               Logger
	Storage              Storage
	NotificationReceiver NotificationReceiver
}

type NotificationReceiver interface {
	Add(StatEvent) error
}

type Logger interface {
	Error(format string, params ...interface{})
	Info(format string, params ...interface{})
}

type Storage interface {
	AddBannerToSlot(bannerID, slotID int) error
	DeleteBannerFromSlot(bannerID, slotID int) error
	ClickBanner(bannerID, slotID, socialGroupID int) error
	GetBannersInfo(slotID, socialGroupID int) ([]BannerStats, error)
	IncrementBannerView(bannerID, slotID, socialGroupID int) error
	GetRandomBanner(slotID int) (int, error)
}

var ErrGetBanner = errors.New("error get banner")

func New(logger Logger, storage Storage, rcv NotificationReceiver) *App {
	return &App{
		Logger:               logger,
		Storage:              storage,
		NotificationReceiver: rcv,
	}
}

func (a *App) AddBannerToSlot(bannerID, slotID int) error {
	if err := a.Storage.AddBannerToSlot(bannerID, slotID); err != nil {
		a.Logger.Error("Error add banner to slot: %s", err)
		return err
	}

	return nil
}

func (a *App) DeleteBannerFromSlot(bannerID, slotID int) error {
	if err := a.Storage.DeleteBannerFromSlot(bannerID, slotID); err != nil {
		a.Logger.Error("Error delete banner from slot: %s", err)
		return err
	}

	return nil
}

func (a *App) ClickBanner(bannerID, slotID, socialGroupID int) error {
	if err := a.Storage.ClickBanner(bannerID, slotID, socialGroupID); err != nil {
		a.Logger.Error("Error click banner from slot: %s", err)
		return err
	}

	statEvent := StatEvent{
		EventType:     "click",
		SlotID:        slotID,
		BannerID:      bannerID,
		SocialGroupID: socialGroupID,
		DateTime:      time.Now(),
	}

	if err := a.NotificationReceiver.Add(statEvent); err != nil {
		return fmt.Errorf("error add statEvent for event 'click':  %w", err)
	}

	a.Logger.Info("Event 'click' on banner id=%s sent", bannerID)

	return nil
}

func (a *App) GetBanner(slotID, socialGroupID int) (int, error) {
	banners, err := a.Storage.GetBannersInfo(slotID, socialGroupID)

	var bannerID int
	if err == nil {
		if bannerID, err = OneHandBandit(banners); err != nil {
			return 0, err
		}
	} else {
		// получение случайного баннера
		// по требованиям нет такого пункта, но если есть проблемы с получением статистических данных из БД
		// для бандита, но система по хорошему все-равно должна выдать какой-нибудь баннер, доступный для слота
		bannerID, err = a.Storage.GetRandomBanner(slotID)
		if err != nil {
			a.Logger.Error("Error get banner for slot: %s", err)
			return 0, err
		}
	}

	if bannerID == 0 {
		return 0, ErrGetBanner
	}

	err = a.Storage.IncrementBannerView(bannerID, slotID, socialGroupID)
	if err != nil {
		a.Logger.Error("error add banner view: %s", err.Error())
	}

	// event
	statEvent := StatEvent{
		EventType:     "view",
		SlotID:        slotID,
		BannerID:      bannerID,
		SocialGroupID: socialGroupID,
		DateTime:      time.Now(),
	}

	if err := a.NotificationReceiver.Add(statEvent); err != nil {
		a.Logger.Info("error add statEvent for event 'view':  %w", err)
	} else {
		a.Logger.Info("Event 'view' on banner id=%s sent", bannerID)
	}

	return bannerID, nil
}
