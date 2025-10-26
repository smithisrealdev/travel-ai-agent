package orchestrator

import (
	"github.com/smithisrealdev/travel-ai-agent/backend/internal/models"
)

// SocialServiceAdapter adapts the services.SocialService to the orchestrator interface
type SocialServiceAdapter struct {
	service interface {
		GetTopRatedPlaces(keyword, location string, limit int) ([]models.SocialPlace, error)
	}
}

// NewSocialServiceAdapter creates a new adapter
func NewSocialServiceAdapter(service interface {
	GetTopRatedPlaces(keyword, location string, limit int) ([]models.SocialPlace, error)
}) *SocialServiceAdapter {
	return &SocialServiceAdapter{service: service}
}

// GetTopRatedPlaces adapts the service call and converts models
func (a *SocialServiceAdapter) GetTopRatedPlaces(keyword, location string, limit int) ([]SocialPlace, error) {
	if a.service == nil {
		return nil, nil
	}
	
	places, err := a.service.GetTopRatedPlaces(keyword, location, limit)
	if err != nil {
		return nil, err
	}
	
	// Convert models.SocialPlace to orchestrator.SocialPlace
	result := make([]SocialPlace, len(places))
	for i, p := range places {
		result[i] = SocialPlace{
			PlaceID:     p.PlaceID,
			Name:        p.Name,
			Address:     p.Address,
			Rating:      p.Rating,
			ReviewCount: p.ReviewCount,
			Types:       p.Types,
		}
	}
	
	return result, nil
}
