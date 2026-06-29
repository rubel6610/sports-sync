package service

import (
	"spotsync-api/dto"
	"spotsync-api/models"
	"spotsync-api/repository"
)

type ZoneService interface {
	CreateZone(req dto.CreateZoneRequest) (*dto.ZoneResponse, error)
	GetAllZones() ([]dto.ZoneResponse, error)
	GetZoneByID(id uint) (*dto.ZoneResponse, error)
}

type zoneService struct {
	zoneRepo repository.ZoneRepository
}

func NewZoneService(repo repository.ZoneRepository) ZoneService {
	return &zoneService{zoneRepo: repo}
}

func (s *zoneService) CreateZone(req dto.CreateZoneRequest) (*dto.ZoneResponse, error) {
	zone := &models.ParkingZone{
		Name:          req.Name,
		Type:          req.Type,
		TotalCapacity: req.TotalCapacity,
		PricePerHour:  req.PricePerHour,
	}
	if err := s.zoneRepo.Create(zone); err != nil {
		return nil, err
	}
	return &dto.ZoneResponse{
		ID:            zone.ID,
		Name:          zone.Name,
		Type:          zone.Type,
		TotalCapacity: zone.TotalCapacity,
		PricePerHour:  zone.PricePerHour,
		CreatedAt:     zone.CreatedAt,
		UpdatedAt:     zone.UpdatedAt,
	}, nil
}

func (s *zoneService) GetAllZones() ([]dto.ZoneResponse, error) {
	zones, spots, err := s.zoneRepo.FindAllWithAvailableSpots()
	if err != nil {
		return nil, err
	}
	res := make([]dto.ZoneResponse, len(zones))
	for i, zone := range zones {
		res[i] = dto.ZoneResponse{
			ID:             zone.ID,
			Name:           zone.Name,
			Type:           zone.Type,
			TotalCapacity:  zone.TotalCapacity,
			AvailableSpots: spots[i],
			PricePerHour:   zone.PricePerHour,
			CreatedAt:      zone.CreatedAt,
		}
	}
	return res, nil
}

func (s *zoneService) GetZoneByID(id uint) (*dto.ZoneResponse, error) {
	zone, spots, err := s.zoneRepo.FindByIDWithAvailableSpots(id)
	if err != nil {
		return nil, err
	}
	return &dto.ZoneResponse{
		ID:             zone.ID,
		Name:           zone.Name,
		Type:           zone.Type,
		TotalCapacity:  zone.TotalCapacity,
		AvailableSpots: spots,
		PricePerHour:   zone.PricePerHour,
		CreatedAt:      zone.CreatedAt,
	}, nil
}