package contextprocessor

import (
	"context"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type contextProcessor struct {
	logger        *zap.Logger
	actionsRunner *ActionsRunner
	cancel        context.CancelFunc
	eventOptions  trace.SpanStartEventOption
}

// implements https://pkg.go.dev/go.opentelemetry.io/collector/component#Component  Start
func (ctxt *contextProcessor) Start(ctx context.Context, host component.Host) error {
	ctx = context.Background()
	ctx, ctxt.cancel = context.WithCancel(ctx)
	for k, _ := range host.GetExtensions() {
		ctxt.logger.Info("Extension", zap.String("id", k.String()))
	}
	return nil
}

// implements https://pkg.go.dev/go.opentelemetry.io/collector/component#Component  Shutdown
func (ctxt *contextProcessor) Shutdown(ctx context.Context) error {
	ctxt.cancel()
	return nil
}
