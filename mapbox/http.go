package mapbox

import (
	"github.com/valyala/fasthttp"
)

type FastHttpClient interface {
	Do(req *fasthttp.Request, resp *fasthttp.Response) error
}
