package logger

import (
	"finance-api/config/env"
	"io"
	"time"
	"tracing"

	"github.com/opentracing/opentracing-go"
	jlog "github.com/opentracing/opentracing-go/log"
)

func JaegerStart(config env.Config, serviceName string, handlerName string, message string) (io.Closer, opentracing.Span, time.Time) {
	if config.GetBool(`logger.jaeger.jaeger_enabled`) {

		startTime := time.Now()

		tracer, closer := tracing.Init(serviceName, config.GetString(`logger.jaeger.jaeger_url`))
		opentracing.SetGlobalTracer(tracer)

		/* starting span */
		span := tracer.StartSpan(handlerName)
		span.SetTag("Method", handlerName)
		span.LogFields(
			jlog.String("Message", startTime.String()+"  ==> "+message),
		)

		return closer, span, startTime
	}
	return nil, nil, time.Now()
}

func JaegerEnd(config env.Config, startTime time.Time, span opentracing.Span, status string, message string) {
	if config.GetBool(`logger.jaeger.jaeger_enabled`) {
		/* End of Tracing */
		endtime := time.Now()
		elapsed := endtime.Sub(startTime)
		span.LogFields(
			jlog.String("Status", status),
			jlog.String("Message", endtime.String()+"  ==> "+message),
			jlog.String("Duration", elapsed.String()),
		)
	}
}
