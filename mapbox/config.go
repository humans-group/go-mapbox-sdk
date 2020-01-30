package mapbox

import (
	"context"
	"os"

	"github.com/valyala/fasthttp"
)

const (
	defaultAPI = "https://api.mapbox.com"
)

type Option func(c config) config

type config struct {
	accessToken   string
	rootAPI       string
	client        fasthttp.Client
	logger        Logger
	requestLogger func(ctx context.Context) Logger

	accessTokenGetValue string
	geocodeEndpoint string
}

func (c config) withEnv() config {
	at := os.Getenv("MAPBOX_ACCESS_TOKEN")
	if at != "" {
		c.accessToken = at
	}

	return c
}

func (c config) prepare() config {
	c.accessTokenGetValue = questionMark + access_token + equalMark + c.accessToken

	return c
}

func newConfig() config {
	return config{
		rootAPI:         defaultAPI,
		client:          fasthttp.Client{},
		geocodeEndpoint: "mapbox.places",
	}
}

func Log(l Logger) Option {
	return func(c config) config {
		c.logger = l
		return c
	}
}

func RequestLogger(extract func(ctx context.Context) Logger) Option {
	return func(c config) config {
		c.requestLogger = extract
		return c
	}
}

func AccessToken(at string) Option {
	return func(c config) config {
		c.accessToken = at
		return c
	}
}

// default to https://api.mapbox.com
func RootAPI(rootAPI string) Option {
	return func(c config) config {
		c.rootAPI = rootAPI
		return c
	}
}

func HttpClient(c fasthttp.Client) Option {
	return func(fhc config) config {
		fhc.client = c
		return fhc
	}
}

// could be set to mapbox.places-permanent, defualt to mapbox.places
func GeocodeEndpoint(endpoint string) Option {
	return func(c config) config {
		c.geocodeEndpoint = endpoint
		return c
	}
}
