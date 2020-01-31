gen:
	easyjson --all mapbox/entities.go
	easyjson mapbox/geocode.go
	minimock -g -i ./mapbox.Geocoder -o ./mapbox -s _mock.go