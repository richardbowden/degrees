package accesscontrol

import (
	"context"

	"time"

	"github.com/openfga/openfga/pkg/logger"
	"github.com/rs/zerolog"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type zerologAdapter struct {
	logger zerolog.Logger
}

func NewZerologAdapter(logger zerolog.Logger) logger.Logger {
	return &zerologAdapter{logger: logger}
}

type zerologEncoder struct {
	event *zerolog.Event
	ctx   *zerolog.Context
}

func (e *zerologEncoder) AddArray(key string, marshaler zapcore.ArrayMarshaler) error {
	if e.event != nil {
		e.event.Interface(key, marshaler)
	} else if e.ctx != nil {
		*e.ctx = e.ctx.Interface(key, marshaler)
	}
	return nil
}

func (e *zerologEncoder) AddObject(key string, marshaler zapcore.ObjectMarshaler) error {
	if e.event != nil {
		e.event.Interface(key, marshaler)
	} else if e.ctx != nil {
		*e.ctx = e.ctx.Interface(key, marshaler)
	}
	return nil
}

func (e *zerologEncoder) AddBinary(key string, value []byte) {
	if e.event != nil {
		e.event.Bytes(key, value)
	} else if e.ctx != nil {
		*e.ctx = e.ctx.Bytes(key, value)
	}
}

func (e *zerologEncoder) AddByteString(key string, value []byte) {
	if e.event != nil {
		e.event.Str(key, string(value))
	} else if e.ctx != nil {
		*e.ctx = e.ctx.Str(key, string(value))
	}
}

func (e *zerologEncoder) AddBool(key string, value bool) {
	if e.event != nil {
		e.event.Bool(key, value)
	} else if e.ctx != nil {
		*e.ctx = e.ctx.Bool(key, value)
	}
}

func (e *zerologEncoder) AddComplex128(key string, value complex128) {
	if e.event != nil {
		e.event.Interface(key, value)
	} else if e.ctx != nil {
		*e.ctx = e.ctx.Interface(key, value)
	}
}

func (e *zerologEncoder) AddComplex64(key string, value complex64) {
	if e.event != nil {
		e.event.Interface(key, value)
	} else if e.ctx != nil {
		*e.ctx = e.ctx.Interface(key, value)
	}
}

func (e *zerologEncoder) AddDuration(key string, value time.Duration) {
	if e.event != nil {
		e.event.Dur(key, value)
	} else if e.ctx != nil {
		*e.ctx = e.ctx.Dur(key, value)
	}
}

func (e *zerologEncoder) AddFloat64(key string, value float64) {
	if e.event != nil {
		e.event.Float64(key, value)
	} else if e.ctx != nil {
		*e.ctx = e.ctx.Float64(key, value)
	}
}

func (e *zerologEncoder) AddFloat32(key string, value float32) {
	if e.event != nil {
		e.event.Float32(key, value)
	} else if e.ctx != nil {
		*e.ctx = e.ctx.Float32(key, value)
	}
}

func (e *zerologEncoder) AddInt(key string, value int) {
	if e.event != nil {
		e.event.Int(key, value)
	} else if e.ctx != nil {
		*e.ctx = e.ctx.Int(key, value)
	}
}

func (e *zerologEncoder) AddInt64(key string, value int64) {
	if e.event != nil {
		e.event.Int64(key, value)
	} else if e.ctx != nil {
		*e.ctx = e.ctx.Int64(key, value)
	}
}

func (e *zerologEncoder) AddInt32(key string, value int32) {
	if e.event != nil {
		e.event.Int32(key, value)
	} else if e.ctx != nil {
		*e.ctx = e.ctx.Int32(key, value)
	}
}

func (e *zerologEncoder) AddInt16(key string, value int16) {
	if e.event != nil {
		e.event.Int16(key, value)
	} else if e.ctx != nil {
		*e.ctx = e.ctx.Int16(key, value)
	}
}

func (e *zerologEncoder) AddInt8(key string, value int8) {
	if e.event != nil {
		e.event.Int8(key, value)
	} else if e.ctx != nil {
		*e.ctx = e.ctx.Int8(key, value)
	}
}

func (e *zerologEncoder) AddString(key, value string) {
	if e.event != nil {
		e.event.Str(key, value)
	} else if e.ctx != nil {
		*e.ctx = e.ctx.Str(key, value)
	}
}

func (e *zerologEncoder) AddTime(key string, value time.Time) {
	if e.event != nil {
		e.event.Time(key, value)
	} else if e.ctx != nil {
		*e.ctx = e.ctx.Time(key, value)
	}
}

func (e *zerologEncoder) AddUint(key string, value uint) {
	if e.event != nil {
		e.event.Uint(key, value)
	} else if e.ctx != nil {
		*e.ctx = e.ctx.Uint(key, value)
	}
}

func (e *zerologEncoder) AddUint64(key string, value uint64) {
	if e.event != nil {
		e.event.Uint64(key, value)
	} else if e.ctx != nil {
		*e.ctx = e.ctx.Uint64(key, value)
	}
}

func (e *zerologEncoder) AddUint32(key string, value uint32) {
	if e.event != nil {
		e.event.Uint32(key, value)
	} else if e.ctx != nil {
		*e.ctx = e.ctx.Uint32(key, value)
	}
}

func (e *zerologEncoder) AddUint16(key string, value uint16) {
	if e.event != nil {
		e.event.Uint16(key, value)
	} else if e.ctx != nil {
		*e.ctx = e.ctx.Uint16(key, value)
	}
}

func (e *zerologEncoder) AddUint8(key string, value uint8) {
	if e.event != nil {
		e.event.Uint8(key, value)
	} else if e.ctx != nil {
		*e.ctx = e.ctx.Uint8(key, value)
	}
}

func (e *zerologEncoder) AddUintptr(key string, value uintptr) {
	if e.event != nil {
		e.event.Interface(key, value)
	} else if e.ctx != nil {
		*e.ctx = e.ctx.Interface(key, value)
	}
}

func (e *zerologEncoder) AddReflected(key string, value interface{}) error {
	if e.event != nil {
		e.event.Interface(key, value)
	} else if e.ctx != nil {
		*e.ctx = e.ctx.Interface(key, value)
	}
	return nil
}

func (e *zerologEncoder) OpenNamespace(key string) {
	// zerolog doesn't have namespaces, we can prefix future keys or ignore
}

func (z *zerologAdapter) Debug(msg string, fields ...zap.Field) {
	z.log(z.logger.Debug(), msg, fields)
}

func (z *zerologAdapter) Info(msg string, fields ...zap.Field) {
	z.log(z.logger.Info(), msg, fields)
}

func (z *zerologAdapter) Warn(msg string, fields ...zap.Field) {
	z.log(z.logger.Warn(), msg, fields)
}

func (z *zerologAdapter) Error(msg string, fields ...zap.Field) {
	z.log(z.logger.Error(), msg, fields)
}

func (z *zerologAdapter) Panic(msg string, fields ...zap.Field) {
	z.log(z.logger.Panic(), msg, fields)
}

func (z *zerologAdapter) Fatal(msg string, fields ...zap.Field) {
	z.log(z.logger.Fatal(), msg, fields)
}

func (z *zerologAdapter) With(fields ...zap.Field) logger.Logger {
	ctx := z.logger.With()
	enc := &zerologEncoder{ctx: &ctx}
	for _, field := range fields {
		field.AddTo(enc)
	}
	return &zerologAdapter{logger: ctx.Logger()}
}

func (z *zerologAdapter) DebugWithContext(ctx context.Context, msg string, fields ...zap.Field) {
	z.log(z.logger.Debug().Ctx(ctx), msg, fields)
}

func (z *zerologAdapter) InfoWithContext(ctx context.Context, msg string, fields ...zap.Field) {
	z.log(z.logger.Info().Ctx(ctx), msg, fields)
}

func (z *zerologAdapter) WarnWithContext(ctx context.Context, msg string, fields ...zap.Field) {
	z.log(z.logger.Warn().Ctx(ctx), msg, fields)
}

func (z *zerologAdapter) ErrorWithContext(ctx context.Context, msg string, fields ...zap.Field) {
	z.log(z.logger.Error().Ctx(ctx), msg, fields)
}

func (z *zerologAdapter) PanicWithContext(ctx context.Context, msg string, fields ...zap.Field) {
	z.log(z.logger.Panic().Ctx(ctx), msg, fields)
}

func (z *zerologAdapter) FatalWithContext(ctx context.Context, msg string, fields ...zap.Field) {
	z.log(z.logger.Fatal().Ctx(ctx), msg, fields)
}

func (z *zerologAdapter) log(event *zerolog.Event, msg string, fields []zap.Field) {
	enc := &zerologEncoder{event: event}
	for _, field := range fields {
		field.AddTo(enc)
	}
	event.Msg(msg)
}
