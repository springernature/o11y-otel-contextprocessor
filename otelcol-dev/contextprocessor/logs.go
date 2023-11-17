package contextprocessor

import (
	"context"

	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/pdata/plog"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type contextLogsProcessor struct {
	contextProcessor
	nextConsumer consumer.Logs
}

func NewContextLogsProcessor(
	logger *zap.Logger,
	nextConsumer consumer.Logs,
	eventOptions trace.SpanStartEventOption,
	actions []ActionConfig) (*contextLogsProcessor, error) {
	exeRunner := NewExeActionsRunner()
	for _, action := range actions {
		if err := exeRunner.AddAction(action); err != nil {
			return nil, err
		}
	}
	return &contextLogsProcessor{
		contextProcessor: contextProcessor{
			logger:           logger,
			exeActionsRunner: exeRunner,
			eventOptions:     eventOptions,
		},
		nextConsumer: nextConsumer,
	}, nil
}

// implements https://pkg.go.dev/go.opentelemetry.io/collector/consumer#Logs
func (ctxt *contextLogsProcessor) ConsumeLogs(ctx context.Context, ld plog.Logs) error {
	span := trace.SpanFromContext(ctx)
	span.AddEvent("Start processing.", ctxt.eventOptions)
	rsl := ld.ResourceLogs()
	newCtx := ctx
	if rsl.Len() > 0 {
		// Only first batch
		attrs := rsl.At(0).Resource().Attributes()
		newCtx = ctxt.exeActionsRunner.Apply(ctx, attrs)
	}
	span.AddEvent("End processing.", ctxt.eventOptions)
	return ctxt.nextConsumer.ConsumeLogs(newCtx, ld)
}
