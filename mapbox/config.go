package mapbox

import (
	"context"
	"os"

	"github.com/valyala/fasthttp"
)

const (
	defaultAPI = "https://api.mapbox.com"
)

// Option allows gradually modify config
type Option func(c config) config

type config struct {
	accessToken   string
	rootAPI       string
	client        FastHttpClient
	logger        Logger
	// requestLogger will be called instead of testLogger if set.
	requestLogger func(ctx context.Context) Logger

	accessTokenGetValue []byte
	geocodeEndpoint string
}

// withEnv overwrites config values with env is not empty
func (c config) withEnv() config {
	at := os.Getenv("MAPBOX_ACCESS_TOKEN")
	if at != "" {
		c.accessToken = at
	}

	return c
}

// prepare prebuilds some reused api parts like access token http get value
func (c config) prepare() config {
	c.accessTokenGetValue = []byte(questionMark + access_token + string(equalMark) + c.accessToken)

	return c
}

func newConfig() config {
	return config{
		rootAPI:         defaultAPI,
		client:          &fasthttp.Client{},
		geocodeEndpoint: "mapbox.places",
	}
}

// Log used to debug traces and to log errors.
func Log(l Logger) Option {
	return func(c config) config {
		c.logger = l
		return c
	}
}

// RequestLogger sets the way testLogger could be extracted from request context.
// If set will be used instead of Log.
func RequestLogger(extract func(ctx context.Context) Logger) Option {
	return func(c config) config {
		c.requestLogger = extract
		return c
	}
}
// AccessToken sets access_token get param.
// Could be set with MAPBOX_ACCESS_TOKEN too.
func AccessToken(at string) Option {
	return func(c config) config {
		c.accessToken = at
		return c
	}
}

// RootAPI allows to change root api address.
// default to https://api.mapbox.com
func RootAPI(rootAPI string) Option {
	return func(c config) config {
		c.rootAPI = rootAPI
		return c
	}
}

// HttpClient allows to change default fast http client
func HttpClient(c FastHttpClient) Option {
	return func(fhc config) config {
		fhc.client = c
		return fhc
	}
}

// GeocodeEndpoint sets geocode endpoint.
// could be set to mapbox.places-permanent, defualt to mapbox.places
func GeocodeEndpoint(endpoint string) Option {
	return func(c config) config {
		c.geocodeEndpoint = endpoint
		return c
	}
}
