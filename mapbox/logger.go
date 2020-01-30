package mapbox

import (
	"context"
)

type Logger interface {
	Debugf(msg string, params ...interface{})
	Errorf(msg string, params ...interface{})
}

func (c *config) withLogger(ctx context.Context, do func(Logger)) {
	if c.requestLogger != nil  {
		do(c.requestLogger(ctx))
		return
	}

	if c.logger != nil {
		do(c.logger)
	}
}