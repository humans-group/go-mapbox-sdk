package mapbox

import (
	"context"
	"testing"

	"github.com/valyala/fasthttp"
)

type fastHttpClient struct {
}

func (_ *fastHttpClient) Do(req *fasthttp.Request, resp *fasthttp.Response) error {
	resp.SetBodyRaw(testRespBody)
	return nil
}

var resp1 *ReverseGeocodeResponse

func Benchmark_Geocoder(b *testing.B) {
	g := NewFastHttpGeocoder(HttpClient(&fastHttpClient{}))
	for i := 0; i <= b.N; i++ {
		resp1, _ = g.ReverseGeocode(context.Background(), &ReverseGeocodeRequest{})
	}
}

var testRespBody = []byte(`{"type":"FeatureCollection","query":[-77.05,38.889],"features":[{"id":"address.6707678235122794","type":"Feature","place_type":["address"],"relevance":1,"properties":{"accuracy":"rooftop"},"text":"Lincoln Memorial Circle SW","place_name":"2 Lincoln Memorial Circle SW, Washington, District of Columbia 20024, United States","center":[-77.0501629,38.8892227],"geometry":{"type":"Point","coordinates":[-77.0501629,38.8892227]},"address":"2","context":[{"id":"neighborhood.295198","text":"National Mall"},{"id":"postcode.4419139247733840","text":"20024"},{"id":"place.7673410831246050","wikidata":"Q61","text":"Washington"},{"id":"region.1753213251667470","short_code":"US-DC","wikidata":"Q3551781","text":"District of Columbia"},{"id":"country.9053006287256050","short_code":"us","wikidata":"Q30","text":"United States"}]},{"id":"neighborhood.295198","type":"Feature","place_type":["neighborhood"],"relevance":1,"properties":{},"text":"National Mall","place_name":"National Mall, Washington, District of Columbia 20024, United States","bbox":[-77.056852,38.8788473,-77.0140495,38.893034],"center":[-77.02,38.89],"geometry":{"type":"Point","coordinates":[-77.02,38.89]},"context":[{"id":"postcode.4419139247733840","text":"20024"},{"id":"place.7673410831246050","wikidata":"Q61","text":"Washington"},{"id":"region.1753213251667470","short_code":"US-DC","wikidata":"Q3551781","text":"District of Columbia"},{"id":"country.9053006287256050","short_code":"us","wikidata":"Q30","text":"United States"}]},{"id":"postcode.4419139247733840","type":"Feature","place_type":["postcode"],"relevance":1,"properties":{},"text":"20024","place_name":"Washington, District of Columbia 20024, United States","bbox":[-77.0644108917888,38.8501751868964,-77.0036921626302,38.8928826270284],"center":[-77.03,38.89],"geometry":{"type":"Point","coordinates":[-77.03,38.89]},"context":[{"id":"place.7673410831246050","wikidata":"Q61","text":"Washington"},{"id":"region.1753213251667470","short_code":"US-DC","wikidata":"Q3551781","text":"District of Columbia"},{"id":"country.9053006287256050","short_code":"us","wikidata":"Q30","text":"United States"}]},{"id":"place.7673410831246050","type":"Feature","place_type":["place"],"relevance":1,"properties":{"wikidata":"Q61"},"text":"Washington","place_name":"Washington, District of Columbia, United States","bbox":[-77.1197609567342,38.79155738,-76.909391,38.99555093],"center":[-77.0366,38.895],"geometry":{"type":"Point","coordinates":[-77.0366,38.895]},"context":[{"id":"region.1753213251667470","short_code":"US-DC","wikidata":"Q3551781","text":"District of Columbia"},{"id":"country.9053006287256050","short_code":"us","wikidata":"Q30","text":"United States"}]},{"id":"region.1753213251667470","type":"Feature","place_type":["region"],"relevance":1,"properties":{"short_code":"US-DC","wikidata":"Q3551781"},"text":"District of Columbia","place_name":"District of Columbia, United States","bbox":[-77.208138,38.717703,-76.909393,38.995548],"center":[-77.03667,38.895],"geometry":{"type":"Point","coordinates":[-77.03667,38.895]},"context":[{"id":"country.9053006287256050","short_code":"us","wikidata":"Q30","text":"United States"}]},{"id":"country.9053006287256050","type":"Feature","place_type":["country"],"relevance":1,"properties":{"short_code":"us","wikidata":"Q30"},"text":"United States","place_name":"United States","bbox":[-179.9,18.765563,-66.885444,71.540724],"center":[-100,40],"geometry":{"type":"Point","coordinates":[-100,40]}}],"attribution":"NOTICE: © 2020 Mapbox and its suppliers. All rights reserved. Use of this data is subject to the Mapbox Terms of Service (https://www.mapbox.com/about/maps/). This response and the information it contains may not be retained. POI(s) provided by Foursquare."}`)
