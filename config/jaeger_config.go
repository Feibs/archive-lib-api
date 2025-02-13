package config

import (
	"io"
	"net/http"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
)

var Tracer opentracing.Tracer

func InitJaeger() (opentracing.Tracer, io.Closer, error) {
	cfg := config.Configuration{
		Sampler: &config.SamplerConfig{
			Type:  jaeger.SamplerTypeRateLimiting,
			Param: 100,
		},
		Reporter: &config.ReporterConfig{
			// Using docker backend
			CollectorEndpoint:  "http://jaeger:14268/api/traces",
			
			// // Using local backend
			// CollectorEndpoint:  "http://localhost:14268/api/traces",
		},
		ServiceName: "archive-lib-service",
	}

	var err error
	var closer io.Closer
	Tracer, closer, err = cfg.NewTracer()
	return Tracer, closer, err
}

func StartSpanFromRequest(tracer opentracing.Tracer, r *http.Request, funcDesc string) opentracing.Span {
	spanCtx, _ := Extract(tracer, r)
	return tracer.StartSpan(funcDesc, ext.RPCServerOption(spanCtx))
}

func Inject(span opentracing.Span, request *http.Request) error {
	return span.Tracer().Inject(
		span.Context(),
		opentracing.HTTPHeaders,
		opentracing.HTTPHeadersCarrier(request.Header))
}

func Extract(tracer opentracing.Tracer, r *http.Request) (opentracing.SpanContext, error) {
	return tracer.Extract(
		opentracing.HTTPHeaders,
		opentracing.HTTPHeadersCarrier(r.Header))
}
