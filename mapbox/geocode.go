package mapbox

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/pkg/errors"

	"github.com/valyala/fasthttp"
)

const (
	limit        = "limit"
	types        = "types"
	country      = "country"
	language     = "language"
	reverseMode  = "reverseMode"
	autocomplete = "autocomplete"
	fuzzymatch   = "fuzzymatch"
	bbox         = "bbox"
	proximity    = "proximity"
	routing      = "routing"
	trueStr      = "true"
	oneStr       = "1"

	access_token = "access_token"

	floatFormatNoExponent = 'f'

	respHeaderRateLimitInterval = "X-Rate-Limit-Interval"
	respHeaderRateLimitLimit    = "X-Rate-Limit-Limit"
	respHeaderRateLimitReset    = "X-Rate-Limit-Reset"
)

var (
	responseFormatJSON = []byte(".json")
	getMethod          = []byte("GET")
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

// RateLimit wraps mapbox API rate limit resp headers
type RateLimit struct {
	Interval []byte
	Limit    []byte
	Reset    []byte
}

// easyjson:json
type rawReverseGeoResp struct {
	Features []Feature `json:"features"`
	Query    []float64 `json:"query"`
}

// easyjson:json
type rawForwardGeoResp struct {
	Features []Feature `json:"features"`
	Query    []string  `json:"query"`
}

// GeocodeResponse
type GeocodeResponse struct {
	RateLimit RateLimit
	// Raw mapbox API response
	RawResp []byte
	// passed query to mapbox
	ReverseQuery GeoPoint
	ForwardQuery []string
	// response result type
	Type string
	// response data
	Features []Feature
}

type ForwardGeocodeRequest struct {
	//The feature you’re trying to look up.
	//This could be an address, a point of interest name, a city name, etc.
	//When searching for points of interest, it can also be a category name (for example, “coffee shop”).
	//For information on categories, see the Point of interest category coverage section.
	//The search text should be expressed as a URL-encoded UTF-8 string,
	//and must not contain the semicolon character (either raw or URL-encoded).
	//Your search text, once decoded, must consist of at most 20 words and numbers separated by spacing and punctuation,
	//and at most 256 characters.
	//
	//The accuracy of coordinates returned by a forward geocoding request can be impacted
	//by how the addresses in the query are formatted. Learn more about address formatting
	//best practices in the https://docs.mapbox.com/help/troubleshooting/address-geocoding-format-guide.
	SearchText string

	//Specify whether to return autocomplete results (true, default) or not (false).
	//When autocomplete is enabled, results will be included that start with the requested string,
	//rather than just responses that match it exactly.
	//For example, a query for India might return both India and Indiana with autocomplete enabled,
	//but only India if it’s disabled.
	//
	//When autocomplete is enabled, each user keystroke counts as one request to the Geocoding API.
	//For example, a search for "coff" would be reflected as four separate Geocoding API requests.
	//To reduce the total requests sent, you can configure your application
	//to only call the Geocoding API after a specific number of characters are typed.
	Autocomplete *bool // default true

	//Limit results to only those contained within the supplied bounding box
	//Bounding boxes should be supplied as four numbers separated by commas,
	//in  minLon,minLat,maxLon,maxLat order.
	//The bounding box cannot cross the 180th meridian.
	Bbox []float64

	//Limit results to one or more countries.
	//Permitted values are ISO 3166 alpha 2 country codes separated by commas.
	Country string

	//Specify whether the Geocoding API should attempt approximate,
	//as well as exact, matching when performing searches (true, default),
	//or whether it should opt out of this behavior and only attempt exact matching (false).
	//For example, the default setting might return Washington, DC for a query of wahsington,
	//even though the query was misspelled.
	FuzzyMatch *bool // default true

	//Specify the user’s language.
	//This parameter controls the language of the text supplied in responses, and also affects result scoring,
	//with results matching the user’s query in the requested language being preferred over results
	//that match in another language. For example, an autocomplete query for things
	//that start with Frank might return Frankfurt as the first result with an English (en) language parameter,
	//but Frankreich (“France”) with a German (de) language parameter.
	//
	//Options are IETF language tags comprised of a mandatory ISO 639-1 language code and, optionally,
	//one or more IETF subtags for country or script.
	//
	//More than one value can also be specified, separated by commas,
	//for applications that need to display labels in multiple languages.
	//
	//For more information on which specific languages are supported, see the https://docs.mapbox.com/api/search/#language-coverage.
	Language string

	//Specify the maximum number of results to return. The default is 5 and the maximum supported is 10.
	Limit int // default 5

	//Bias the response to favor results that are closer to this location
	Proximity *GeoPoint

	//Specify whether to request additional metadata about the recommended navigation destination
	//corresponding to the feature (true) or not (false, default). Only applicable for address features.
	//
	//For example, if routing=true the response could include data about a point on the road the feature fronts.
	//Response features may include an array containing one or more routable points.
	//Routable points cannot always be determined.
	//Consuming applications should fall back to using the feature’s normal geometry for routing
	//if a separate routable point is not returned.
	Routing bool //default false

	//Filter results to include only a subset (one or more) of the available feature types.
	//Options are country, region, postcode, district, place, locality, neighborhood, address, and poi.
	//Multiple options can be comma-separated. Note that poi.landmark is a deprecated type that,
	//while still supported, returns the same data as is returned using the poi type.
	//
	//For more information on the available types, see the https://docs.mapbox.com/api/search/#data-types.
	Types []string
}

// Geocoder encapsulates forward and reverse geocode calls.
type Geocoder interface {
	// ReverseGeocode calls geocode/v5 reverse mapbox API
	ReverseGeocode(ctx context.Context, req *ReverseGeocodeRequest) (*GeocodeResponse, error)
	// ReverseGeocode calls geocode/v5 reverse mapbox API
	ForwardGeocode(ctx context.Context, req *ForwardGeocodeRequest) (*GeocodeResponse, error)
}

// FastHttpGeocoder is a fasthttp Geocoder implementation
type FastHttpGeocoder struct {
	config

	geocodeAPIURL []byte

	stringBufPull *stringsBufferPool
}

// ReverseGeocode calls geocode/v5 reverse mapbox API thought fasthttp client.
func (c *FastHttpGeocoder) ReverseGeocode(ctx context.Context, req *ReverseGeocodeRequest) (*GeocodeResponse, error) {
	freq := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(freq)

	fresp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(fresp)

	// split multivalues to limit memory consumption
	values := make(map[string]string, 5)

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
	if len(req.Types) > 0 {
		values[types] = strings.Join(req.Types, ",")
	}

	buf := c.stringBufPull.acquireStringsBuilder()
	defer c.stringBufPull.releaseStringsBuilder(buf)

	buf.Write(c.geocodeAPIURL)
	buf.WriteString(strconv.FormatFloat(req.GeoPoint.Lon, floatFormatNoExponent, 6, 64))
	buf.WriteByte(comma)
	buf.WriteString(strconv.FormatFloat(req.GeoPoint.Lat, floatFormatNoExponent, 6, 64))
	buf.Write(responseFormatJSON)
	buf.Write(c.accessTokenGetValue)

	encodeValues(buf, values)

	reqURI := buf.Bytes()

	c.withLogger(ctx, func(logger Logger) {
		logger.Debugf("mapbox_sdk: reverse geocode request %s", buf.String())
	})

	freq.Header.SetMethodBytes(getMethod)
	freq.SetRequestURIBytes(reqURI)

	if err := c.client.Do(freq, fresp); err != nil {
		return nil, err
	}

	respBytes := make([]byte, len(fresp.Body()))
	copy(respBytes, fresp.Body())

	c.withLogger(ctx, func(logger Logger) {
		logger.Debugf("mapbox_sdk: reverse geocode response %s", string(respBytes))
	})

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

	return &GeocodeResponse{
		RateLimit: readRespRateLimit(fresp),
		RawResp:   respBytes,
		ReverseQuery: GeoPoint{
			Lon: respRaw.Query[0],
			Lat: respRaw.Query[1],
		},
		Features: respRaw.Features,
	}, nil
}

// ReverseGeocode calls geocode/v5 reverse mapbox API thought fasthttp client.
func (c *FastHttpGeocoder) ForwardGeocode(ctx context.Context, req *ForwardGeocodeRequest) (*GeocodeResponse, error) {
	freq := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(freq)

	fresp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(fresp)

	// split multivalues to limit memory consumption
	values := make(map[string]string, 9)

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
	if req.Autocomplete != nil {
		values[autocomplete] = fmt.Sprint(*req.Autocomplete)
	} else {
		values[autocomplete] = trueStr
	}
	if req.FuzzyMatch != nil {
		values[fuzzymatch] = fmt.Sprint(*req.FuzzyMatch)
	} else {
		values[fuzzymatch] = trueStr
	}
	if len(req.Bbox) == 4 {
		values[bbox] = fmt.Sprintf("%f,%f,%f,%f", req.Bbox[0], req.Bbox[1], req.Bbox[2], req.Bbox[3])
	}
	if req.Proximity != nil {
		values[proximity] = fmt.Sprintf("%f,%f", req.Proximity.Lon, req.Proximity.Lat)
	}
	values[routing] = fmt.Sprint(req.Routing)
	if len(req.Types) > 0 {
		values[types] = strings.Join(req.Types, ",")
	}

	buf := c.stringBufPull.acquireStringsBuilder()
	defer c.stringBufPull.releaseStringsBuilder(buf)

	buf.Write(c.geocodeAPIURL)
	buf.WriteString(req.SearchText)
	buf.Write(responseFormatJSON)
	buf.Write(c.accessTokenGetValue)

	encodeValues(buf, values)

	reqURI := buf.Bytes()

	c.withLogger(ctx, func(logger Logger) {
		logger.Debugf("mapbox_sdk: forward geocode request %s", buf.String())
	})

	freq.Header.SetMethodBytes(getMethod)
	freq.SetRequestURIBytes(reqURI)

	if err := c.client.Do(freq, fresp); err != nil {
		return nil, err
	}

	respBytes := make([]byte, len(fresp.Body()))
	copy(respBytes, fresp.Body())

	c.withLogger(ctx, func(logger Logger) {
		logger.Debugf("mapbox_sdk: forward geocode response %s", string(respBytes))
	})

	if fresp.Header.StatusCode() != http.StatusOK {
		return nil, errors.Errorf("failed to reverse geocode URI %s statusCode %d resp %s",
			reqURI, fresp.Header.StatusCode(), string(respBytes))
	}

	respRaw := rawForwardGeoResp{}
	if err := respRaw.UnmarshalJSON(respBytes); err != nil {
		return nil, errors.Wrapf(err, "failed to unmarshall raw reverse geocode resp %s", string(respBytes))
	}

	return &GeocodeResponse{
		RateLimit:    readRespRateLimit(fresp),
		RawResp:      respBytes,
		Features:     respRaw.Features,
		ForwardQuery: respRaw.Query,
	}, nil
}

func NewFastHttpGeocoder(opts ...Option) *FastHttpGeocoder {
	c := FastHttpGeocoder{
		config:        newConfig(),
		stringBufPull: newStringsBufferPool(),
		geocodeAPIURL: []byte("/geocoding/v5/"),
	}

	for _, o := range opts {
		c.config = o(c.config)
	}

	c.config = c.config.withEnv()
	c.config = c.config.prepare()

	c.geocodeAPIURL = []byte(c.rootAPI + string(c.geocodeAPIURL) + c.geocodeEndpoint + slash)

	return &c
}

func readRespRateLimit(resp *fasthttp.Response) RateLimit {
	return RateLimit{
		Interval: resp.Header.Peek(respHeaderRateLimitInterval),
		Limit:    resp.Header.Peek(respHeaderRateLimitLimit),
		Reset:    resp.Header.Peek(respHeaderRateLimitReset),
	}
}
