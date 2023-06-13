package models

import (
	"gorm.io/gorm"
)

type Location struct {
	gorm.Model
	ID      uint   `gorm:"primaryKey"`
	UserId  uint   `gorm:"not_null;index"`
	Lon     string `gorm:"not null; unique"`
	Lat     string `gorm:"not null; unique"`
	Name    string `gorm:"not null"`
	IsSaved bool   `gorm:"not null"`
}

type LocationDB interface {
	Create(location *Location) error
	Delete(id uint) error
	FindByLonAndLat(lon, lat float64) (Location, error)
	GetByUserId(userId uint) ([]Location, error)
}

type LocationService interface {
	LocationDB
}

func NewLocationService(db *gorm.DB) LocationService {
	return &locationService{
		LocationDB: &locationValidator{
			&locationGorm{db}},
	}
}

type locationService struct {
	LocationDB
}

type locationValidator struct {
	LocationDB
}

var _ LocationDB = &locationValidator{}

func newLocationValidator(ldb LocationDB) *locationValidator {
	return &locationValidator{
		LocationDB: ldb,
	}
}

func (gv *locationValidator) userIDRequired(l *Location) error {
	if l.UserId <= 0 {
		return ErrUserIDRequired
	}
	return nil
}

func (lv *locationValidator) Delete(id uint) error {
	var location Location
	location.ID = id
	if id <= 0 {
		return ErrIDInvalid
	}

	return lv.LocationDB.Delete(id)
}

type locationGorm struct {
	db *gorm.DB
}

func (lg *locationGorm) Create(location *Location) error {
	return lg.db.Create(location).Error
}

func (lv *locationValidator) Create(location *Location) error {
	// Order of functions passed in to validator is important!
	err := runModelValFuncs(
		location, lv.userIDRequired,
	)
	if err != nil {
		return err
	}

	return lv.LocationDB.Create(location)
}

func (lv *locationValidator) GetByLonAndLat(lon, lat float64) (Location, error) {
	// TODO: does this need to be validated?
	return lv.LocationDB.FindByLonAndLat(lon, lat)
}

func (lv *locationValidator) GetByUserId(userId uint) ([]Location, error) {
	// Order of functions passed in to validator is important!
	if userId <= 0 {
		return nil, ErrUserIDRequired
	}
	return lv.LocationDB.GetByUserId(userId)
}

func (lg *locationGorm) Delete(id uint) error {
	location := Location{Model: gorm.Model{ID: id}}
	return lg.db.Unscoped().Delete(&location, id).Error
}

func (lg *locationGorm) GetByUserId(userID uint) ([]Location, error) {
	var locations []Location
	err := lg.db.Where("user_id = ?", userID).Find(&locations).Error
	if err != nil {
		return nil, err
	}
	return locations, nil
}

func (lg *locationGorm) FindByLonAndLat(lon, lat float64) (Location, error) {
	var location Location

	err := lg.db.Where("lon = ? AND lat = ?", lon, lat).Find(&location).Error
	if err != nil {
		return Location{}, err
	}
	return location, nil
}
