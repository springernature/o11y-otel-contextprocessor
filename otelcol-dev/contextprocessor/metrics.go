package contextprocessor

import (
	"context"

	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type contextMetricsProcessor struct {
	contextProcessor
	nextConsumer consumer.Metrics
}

func NewContextMetricsProcessor(
	logger *zap.Logger,
	nextConsumer consumer.Metrics,
	eventOptions trace.SpanStartEventOption,
	actions []ActionConfig) (*contextMetricsProcessor, error) {
	exeRunner := NewExeActionsRunner()
	for _, action := range actions {
		if err := exeRunner.AddAction(action); err != nil {
			return nil, err
		}
	}
	return &contextMetricsProcessor{
		contextProcessor: contextProcessor{
			logger:           logger,
			exeActionsRunner: exeRunner,
			eventOptions:     eventOptions,
		},
		nextConsumer: nextConsumer,
	}, nil
}

// implements https://pkg.go.dev/go.opentelemetry.io/collector/consumer#Metrics
func (ctxt *contextMetricsProcessor) ConsumeMetrics(ctx context.Context, md pmetric.Metrics) error {
	span := trace.SpanFromContext(ctx)
	span.AddEvent("Start processing.", ctxt.eventOptions)
	rms := md.ResourceMetrics()
	newCtx := ctx
	if rms.Len() > 0 {
		// Only first batch
		attrs := rms.At(0).Resource().Attributes()
		newCtx = ctxt.exeActionsRunner.Apply(ctx, attrs)
	}
	span.AddEvent("End processing.", ctxt.eventOptions)
	return ctxt.nextConsumer.ConsumeMetrics(newCtx, md)
}
