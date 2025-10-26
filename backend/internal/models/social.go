package models

// SocialPlaceRequest represents a request to fetch socially popular places
type SocialPlaceRequest struct {
	Keyword  string  `json:"keyword"`
	Location string  `json:"location"`
	Radius   int     `json:"radius,omitempty"` // in meters
	Limit    int     `json:"limit,omitempty"`
}

// SocialPlace represents a socially popular place
type SocialPlace struct {
	PlaceID      string   `json:"placeId"`
	Name         string   `json:"name"`
	Address      string   `json:"address"`
	Rating       float64  `json:"rating"`
	ReviewCount  int      `json:"reviewCount"`
	PriceLevel   int      `json:"priceLevel,omitempty"` // 0-4 scale
	Types        []string `json:"types,omitempty"`
	PhotoURL     string   `json:"photoUrl,omitempty"`
	Latitude     float64  `json:"latitude"`
	Longitude    float64  `json:"longitude"`
	OpenNow      bool     `json:"openNow,omitempty"`
}

// SocialPlacesResponse represents the response for social places
type SocialPlacesResponse struct {
	Places []SocialPlace `json:"places"`
	Query  string        `json:"query"`
	Count  int           `json:"count"`
}
