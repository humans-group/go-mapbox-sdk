package mapbox

import (
	"context"
	"testing"

	"github.com/gojuno/minimock/v3"
)

func Test_config_withLogger(t *testing.T) {
	tests := []struct {
		name   string
		logger        func(mc *minimock.Controller) Logger
		requestLogger func(mc *minimock.Controller) func(context.Context) Logger
	}{
		{
			name:"testLogger set",
			logger : func(mc *minimock.Controller) Logger {
				mock := NewLoggerMock(mc)
				mock.DebugfMock.Return()
				return mock
			},
		},
		{
			name:"request testLogger set",
			requestLogger: func(mc *minimock.Controller) func(context.Context) Logger {
				mock := NewLoggerMock(mc)
				mock.DebugfMock.Return()
				return func(context.Context) Logger{
					return mock
				}
			},
		},
		{
			name:"both loggers set",
			logger : func(mc *minimock.Controller) Logger {
				mock := NewLoggerMock(mc)
				return mock
			},
			requestLogger: func(mc *minimock.Controller) func(context.Context) Logger {
				mock := NewLoggerMock(mc)
				mock.DebugfMock.Return()
				return func(context.Context) Logger{
					return mock
				}
			},
		},

	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mc := minimock.NewController(t)
			c := config{}
			if tt.logger != nil {
				c.logger = tt.logger(mc)
			}
			if tt.requestLogger != nil {
				c.requestLogger = tt.requestLogger(mc)
			}
			c.withLogger(context.Background(), func(l Logger) {
				l.Debugf("")
			})
			mc.Finish()
		})
	}
}
