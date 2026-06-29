package repository

import (
	"errors"
	"spotsync-api/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ReservationRepository interface {
	CreateAtomicReservation(userID uint, zoneID uint, licensePlate string) (*models.Reservation, error)
	FindByUserID(userID uint) ([]models.Reservation, error)
	FindAllWithPreload() ([]models.Reservation, error)
	FindByID(id uint) (*models.Reservation, error)
	UpdateStatus(id uint, status string) error
}

type reservationRepository struct {
	db *gorm.DB
}

func NewReservationRepository(db *gorm.DB) ReservationRepository {
	return &reservationRepository{db: db}
}

func (r *reservationRepository) CreateAtomicReservation(userID uint, zoneID uint, licensePlate string) (*models.Reservation, error) {
	var reservation models.Reservation

	err := r.db.Transaction(func(tx *gorm.DB) error {
		var zone models.ParkingZone
		// Concurrency Rule Lock: FOR UPDATE row level locking
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&zone, zoneID).Error; err != nil {
			return errors.New("parking zone not found")
		}

		var activeReservations int64
		tx.Model(&models.Reservation{}).Where("zone_id = ? AND status = 'active'", zoneID).Count(&activeReservations)

		if int(activeReservations) >= zone.TotalCapacity {
			return errors.New("zone_full")
		}

		reservation = models.Reservation{
			UserID:       userID,
			ZoneID:       zoneID,
			LicensePlate: licensePlate,
			Status:       "active",
		}

		if err := tx.Create(&reservation).Error; err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	return &reservation, nil
}

func (r *reservationRepository) FindByUserID(userID uint) ([]models.Reservation, error) {
	var res []models.Reservation
	err := r.db.Preload("Zone").Where("user_id = ?", userID).Find(&res).Error
	return res, err
}

func (r *reservationRepository) FindAllWithPreload() ([]models.Reservation, error) {
	var res []models.Reservation
	err := r.db.Preload("User").Preload("Zone").Find(&res).Error
	return res, err
}

func (r *reservationRepository) FindByID(id uint) (*models.Reservation, error) {
	var res models.Reservation
	if err := r.db.First(&res, id).Error; err != nil {
		return nil, err
	}
	return &res, nil
}

func (r *reservationRepository) UpdateStatus(id uint, status string) error {
	return r.db.Model(&models.Reservation{}).Where("id = ?", id).Update("status", status).Error
}