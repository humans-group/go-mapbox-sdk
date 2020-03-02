package mapbox

// Client covers all Mabpox API
type Client interface {
	// Geocoder covers forward and reverse geocoding mapbox API
	Geocoder
}