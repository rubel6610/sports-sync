package repository

import (
	"spotsync-api/models"
	"gorm.io/gorm"
)

type ZoneRepository interface {
	Create(zone *models.ParkingZone) error
	FindAllWithAvailableSpots() ([]models.ParkingZone, []int, error)
	FindByIDWithAvailableSpots(id uint) (*models.ParkingZone, int, error)
}

type zoneRepository struct {
	db *gorm.DB
}

func NewZoneRepository(db *gorm.DB) ZoneRepository {
	return &zoneRepository{db: db}
}

func (r *zoneRepository) Create(zone *models.ParkingZone) error {
	return r.db.Create(zone).Error
}

func (r *zoneRepository) FindAllWithAvailableSpots() ([]models.ParkingZone, []int, error) {
	var zones []models.ParkingZone
	if err := r.db.Find(&zones).Error; err != nil {
		return nil, nil, err
	}

	spots := make([]int, len(zones))
	for i, zone := range zones {
		var count int64
		r.db.Model(&models.Reservation{}).Where("zone_id = ? AND status = 'active'", zone.ID).Count(&count)
		spots[i] = zone.TotalCapacity - int(count)
	}
	return zones, spots, nil
}

func (r *zoneRepository) FindByIDWithAvailableSpots(id uint) (*models.ParkingZone, int, error) {
	var zone models.ParkingZone
	if err := r.db.First(&zone, id).Error; err != nil {
		return nil, 0, err
	}
	var count int64
	r.db.Model(&models.Reservation{}).Where("zone_id = ? AND status = 'active'", zone.ID).Count(&count)
	return &zone, zone.TotalCapacity - int(count), nil
}