package mapbox

import (
	"strings"
)

const (
	slash         = "/"
	comma         = ","
	questionMark  = "?"
	equalMark     = "="
	ampersandMark = "&"
)

func encodeValues(buf *strings.Builder, values map[string]string, valuesMulti map[string][]string) {
	for k, v := range values {
		buf.WriteString(ampersandMark)

		buf.WriteString(k)
		buf.WriteString(equalMark)
		buf.WriteString(v)
	}

	for k, vs := range valuesMulti {
		for _, v := range vs {
			buf.WriteString(ampersandMark)

			buf.WriteString(k)
			buf.WriteString(equalMark)
			buf.WriteString(v)
		}
	}
}
