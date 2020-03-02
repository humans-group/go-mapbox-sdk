package mapbox

import (
	"bytes"
)

const (
	slash         = "/"
	comma         = ','
	questionMark  = "?"
	equalMark     = '='
	ampersandMark = '&'
)

// encodeValues do almost the same as url.Values.Encode() but faster and reuses *strings.Builder
func encodeValues(buf *bytes.Buffer, values map[string]string, valuesMulti map[string][]string) {
	for k, v := range values {
		encodeHttpGetKeyValue(buf, k, v)
	}

	for k, vs := range valuesMulti {
		for _, v := range vs {
			encodeHttpGetKeyValue(buf, k, v)
		}
	}
}

func encodeHttpGetKeyValue(buf *bytes.Buffer, k string, v string) {
	buf.WriteByte(ampersandMark)
	buf.WriteString(k)
	buf.WriteByte(equalMark)
	buf.WriteString(v)
}
