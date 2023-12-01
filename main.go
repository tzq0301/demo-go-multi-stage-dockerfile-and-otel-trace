package main

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

var tracer trace.Tracer

func main() {
	tracer = sdktrace.NewTracerProvider().Tracer("DockerDemoService")

	ctx := parseSpanFromString(context.TODO())
	log(ctx)

	f(ctx)
}

func log(ctx context.Context) {
	spanContext := trace.SpanContextFromContext(ctx)
	fmt.Printf("TraceID = %s, SpanID = %s\n", spanContext.TraceID().String(), spanContext.SpanID().String())
}

func parseSpanFromString(ctx context.Context) context.Context {
	traceparent := "00-80e1afed08e019fc1110464cfa66635c-7a085853722dc6d2-01"

	splits := strings.Split(traceparent, "-")
	mustTrue(len(splits) == 4)

	var (
		version       = splits[0]
		rawTraceID    = splits[1]
		rawParentID   = splits[2]
		rawTraceFlags = splits[3]
	)

	mustTrue(version == "00")

	spanContext := trace.NewSpanContext(trace.SpanContextConfig{
		TraceID:    must(trace.TraceIDFromHex(rawTraceID)),
		SpanID:     must(trace.SpanIDFromHex(rawParentID)),
		TraceFlags: trace.TraceFlags(must(strconv.Atoi(rawTraceFlags))),
	})

	return trace.ContextWithSpanContext(ctx, spanContext)
}

func f(ctx context.Context) {
	ctx, span := tracer.Start(ctx, "f")
	defer span.End()

	log(ctx)
}

func must[T any](t T, err error) T {
	if err != nil {
		panic(err)
	}
	return t
}

func mustTrue(b bool) {
	if !b {
		panic(b)
	}
}
