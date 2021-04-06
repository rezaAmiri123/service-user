package tracing

import (
	"context"
	"fmt"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	zipkinot "github.com/openzipkin-contrib/zipkin-go-opentracing"
	"github.com/openzipkin/zipkin-go"
	zipkinhttp "github.com/openzipkin/zipkin-go/reporter/http"
	"github.com/sirupsen/logrus"
)

// Trace instance
var tracer opentracing.Tracer

// SetTracer can be used by unit tests to provide a NoopTracer instance. Real users should always
// use the InitTracing func.
func SetTracer(initiateTracer opentracing.Tracer){
	tracer = initiateTracer
}

// InitTracing connects the calling service to Zipkin and initializes the tracer.
func InitTracing(zipkinURL, serviceName, endPointAddress string){
	logrus.Infof("Connecting to zipkin server at %v", zipkinURL)
	reporter := zipkinhttp.NewReporter(fmt.Sprintf("%s/api/v1/spans", zipkinURL))
	//endPointAddress = "127.0.0.1:0"
	endpoint,err := zipkin.NewEndpoint(serviceName, endPointAddress)
	if err != nil {
		logrus.Fatalf("unable to create local endpoint: %+v\n", err)
	}
	nativeTracer,err := zipkin.NewTracer(reporter, zipkin.WithLocalEndpoint(endpoint))
	if err != nil {
		logrus.Fatalf("unable to create tracer: %+v\n", err)
	}

	// use zipkin-go-opentracing to wrap our tracer
	tracer = zipkinot.Wrap(nativeTracer)
	logrus.Infof("Successfully started zipkin tracer for service '%v'", serviceName)
}

// StartHTTPTrace loads tracing information from an INCOMING HTTP request.
//func StartHTTPTrace(r *http.Request, opName string) opentracing.Span {
//	carrier := opentracing.HTTPHeadersCarrier(r.Header)
//	clientContext, err := tracer.Extract(opentracing.HTTPHeaders, carrier)
//	if err == nil {
//		return tracer.StartSpan(opName, ext.RPCServerOption(clientContext))
//	} else {
//		return tracer.StartSpan(opName)
//	}
//}
//StartHTTPTrace loads tracing information from an INCOMING HTTP request.
func StartHTTPTrace(ctx context.Context, opName string) opentracing.Span {
	return tracer.StartSpan(opName,ext.RPCServerOption(ctx))
}
