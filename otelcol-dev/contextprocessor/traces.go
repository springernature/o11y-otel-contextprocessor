package contextprocessor

import (
	"context"

	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/pdata/ptrace"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type contextTracesProcessor struct {
	contextProcessor
	nextConsumer consumer.Traces
}

func NewContextTracesProcessor(
	logger *zap.Logger,
	nextConsumer consumer.Traces,
	eventOptions trace.SpanStartEventOption,
	actions []ActionConfig) (*contextTracesProcessor, error) {
	exeRunner := NewExeActionsRunner()
	for _, action := range actions {
		if err := exeRunner.AddAction(action); err != nil {
			return nil, err
		}
	}
	return &contextTracesProcessor{
		contextProcessor: contextProcessor{
			logger:           logger,
			exeActionsRunner: exeRunner,
			eventOptions:     eventOptions,
		},
		nextConsumer: nextConsumer,
	}, nil
}

// implements https://pkg.go.dev/go.opentelemetry.io/collector/consumer#Traces
func (ctxt *contextTracesProcessor) ConsumeTraces(ctx context.Context, td ptrace.Traces) error {
	span := trace.SpanFromContext(ctx)
	span.AddEvent("Start processing.", ctxt.eventOptions)
	rss := td.ResourceSpans()
	newCtx := ctx
	if rss.Len() > 0 {
		// Only first batch
		attrs := rss.At(0).Resource().Attributes()
		newCtx = ctxt.exeActionsRunner.Apply(ctx, attrs)
	}
	span.AddEvent("End processing.", ctxt.eventOptions)
	return ctxt.nextConsumer.ConsumeTraces(newCtx, td)
}
