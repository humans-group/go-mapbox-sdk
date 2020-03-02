gen:
	easyjson --all mapbox/entities.go
	easyjson mapbox/geocode.go
	minimock -g -i ./mapbox.Geocoder -o ./mapbox -s _mock.go
	minimock -g -i ./mapbox.Logger -o ./mapbox -s _mock.go

test:
	go test -race -coverprofile=coverage.txt -covermode=atomic -v ./...