package models

import "gorm.io/gorm"

type Location struct {
	gorm.Model
	ID     uint   `gorm:"primaryKey"`
	UserId uint   `gorm:"not_null;index"`
	Lon    string `gorm:"not null;unique"`
	Lat    string `gorm:"not null;unique"`
	Name   string `gorm:"not null"`
}

type LocationDB interface {
	ByEmail(email string) (*Location, error)
	// CRUD
	Create(location *Location) error
	Delete(id uint) error
	GetAll() ([]Location, error)
}
