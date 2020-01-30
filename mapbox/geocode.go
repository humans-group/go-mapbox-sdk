package mapbox

import (
	"context"
	"net/http"
	"strconv"

	"github.com/pkg/errors"

	"github.com/valyala/fasthttp"
)

const (
	limit       = "limit"
	types       = "types"
	country     = "country"
	language    = "language"
	reverseMode = "reverseMode"
	routing     = "routing"
	trueStr     = "true"
	oneStr      = "1"

	responseFormatJSON = ".json"
	access_token       = "access_token"

	getMethod = "GET"

	floatFormatNoExponent = 'f'

	respHeaderRateLimitInterval = "X-Rate-Limit-Interval"
	respHeaderRateLimitLimit    = "X-Rate-Limit-Limit"
	respHeaderRateLimitReset    = "X-Rate-Limit-Reset"
)

type GeoPoint struct {
	Lon float64
	Lat float64
}

type ReverseGeocodeRequest struct {
	GeoPoint GeoPoint
	// Limit results to one or more countries.
	Limit int
	// Filter results to include only a subset (one or more) of the available feature types.
	// Options are country, region, postcode, district, place, locality, neighborhood, address, and poi.
	// Multiple options can be comma-separated. Note that poi.landmark is a deprecated type that, while still supported,
	// returns the same data as is returned using the poi type.
	Types []string
	// Permitted values are ISO 3166 alpha 2(https://en.wikipedia.org/wiki/ISO_3166-1_alpha-2) country codes separated by commas.
	Country string
	// Specify the user’s language. This parameter controls the language of the text supplied in responses.
	// Options are IETF language tags comprised of a mandatory ISO 639-1 language code and, optionally,
	// one or more IETF subtags for country or script.
	// More than one value can also be specified, separated by commas,
	// for applications that need to display labels in multiple languages.
	// For more information on which specific languages are supported, see https://docs.mapbox.com/api/search/#language-coverage
	Language string
	// Decides how results are sorted in a reverse geocoding query
	// if multiple results are requested using a limit other than 1.
	// Options are distance (default), which causes the closest feature
	// to always be returned first, and score, which allows high-prominence features
	// to be sorted higher than nearer, lower-prominence features.
	ReverseMode int
	// Specify whether to request additional metadata about the recommended navigation destination corresponding
	// to the feature (true) or not (false, default). Only applicable for address features.
	// For example, if routing=true the response could include data about a point on the road the feature fronts.
	// Response features may include an array containing one or more routable points.
	// Routable points cannot always be determined.
	// Consuming applications should fall back to using the feature’s normal geometry for routing
	// if a separate routable point is not returned.
	Routing bool
}

type RateLimit struct {
	Interval string
	Limit    string
	Reset    string
}

//easyjson:json
type rawReverseGeoResp struct {
	Features []Feature `json:"features"`
	Query    []float64 `json:"query"`
}

type ReverseGeocodeResponse struct {
	RateLimit RateLimit
	RawResp   []byte
	Query     GeoPoint
	Type      string
	Features  []Feature
}

type Geocoder interface {
	ReverseGeocode(ctx context.Context, req *ReverseGeocodeRequest) (*ReverseGeocodeResponse, error)
}

type FastHttpGeocoder struct {
	config

	geocodeAPIURL string

	stringBufPull *stringsBufferPool
}

func (c *FastHttpGeocoder) ReverseGeocode(ctx context.Context, req *ReverseGeocodeRequest) (*ReverseGeocodeResponse, error) {
	freq := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(freq)

	fresp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(fresp)

	values := map[string]string{}
	var valuesMulti map[string][]string

	if req.Country != "" {
		values[country] = req.Country
	}
	if req.Limit != 0 {
		values[limit] = strconv.Itoa(req.Limit)
	}
	if req.Language != "" {
		values[language] = req.Language
	}
	if req.Routing {
		values[routing] = trueStr
	}
	if req.ReverseMode == 1 {
		values[reverseMode] = oneStr
	}

	switch len(req.Types) {
	case 0:
	case 1:
		values[types] = req.Types[0]
	default:
		valuesMulti = map[string][]string{}
		for _, t := range req.Types {
			valuesMulti[types] = append(valuesMulti[types], t)
		}
	}

	buf := c.stringBufPull.acquireStringsBuilder()
	defer c.stringBufPull.releaseStringsBuilder(buf)

	buf.WriteString(c.geocodeAPIURL)
	buf.WriteString(strconv.FormatFloat(req.GeoPoint.Lon, floatFormatNoExponent, 6, 64))
	buf.WriteString(comma)
	buf.WriteString(strconv.FormatFloat(req.GeoPoint.Lat, floatFormatNoExponent, 6, 64))
	buf.WriteString(responseFormatJSON)
	buf.WriteString(c.accessTokenGetValue)

	encodeValues(buf, values, valuesMulti)

	reqURI := buf.String()

	c.withLogger(ctx, func(logger Logger) {
		logger.Debugf("mapbox_sdk: reverse geocode request %s", reqURI)
	})

	freq.Header.SetMethod(getMethod)
	freq.SetRequestURI(reqURI)

	if err := c.client.Do(freq, fresp); err != nil {
		return nil, err
	}

	respBytes := make([]byte, len(fresp.Body()))
	copy(respBytes, fresp.Body())

	if fresp.Header.StatusCode() != http.StatusOK {
		return nil, errors.Errorf("failed to reverse geocode URI %s statusCode %d resp %s",
			reqURI, fresp.Header.StatusCode(), string(respBytes))
	}

	respRaw := rawReverseGeoResp{}
	if err := respRaw.UnmarshalJSON(respBytes); err != nil {
		return nil, errors.Wrapf(err, "failed to unmarshall raw reverse geocode resp %s", string(respBytes))
	}

	if len(respRaw.Query) != 2 {
		return nil, errors.Errorf("unexpected len of query coordinates in resp %s", string(respBytes))
	}

	return &ReverseGeocodeResponse{
		RateLimit: readRespRateLimit(fresp),
		RawResp:   respBytes,
		Query: GeoPoint{
			Lon: respRaw.Query[0],
			Lat: respRaw.Query[1],
		},
		Features: respRaw.Features,
	}, nil
}

func NewFastHttpGeocoder(opts ...Option) *FastHttpGeocoder {
	c := FastHttpGeocoder{
		config:        newConfig(),
		stringBufPull: newStringsBufferPool(),
		geocodeAPIURL: "/geocoding/v5/",
	}

	for _, o := range opts {
		c.config = o(c.config)
	}

	c.config = c.config.withEnv()
	c.config = c.config.prepare()

	c.geocodeAPIURL = c.rootAPI + c.geocodeAPIURL + c.geocodeEndpoint + slash

	return &c
}

func readRespRateLimit(resp *fasthttp.Response) RateLimit {
	return RateLimit{
		Interval: string(resp.Header.Peek(respHeaderRateLimitInterval)),
		Limit:    string(resp.Header.Peek(respHeaderRateLimitLimit)),
		Reset:    string(resp.Header.Peek(respHeaderRateLimitReset)),
	}
}
