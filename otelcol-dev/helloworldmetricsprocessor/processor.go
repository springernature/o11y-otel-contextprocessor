package helloworldmetricsprocessor

import (
	"context"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.uber.org/zap"
)

type helloWorldMetricsProcessor struct {
	logger               *zap.Logger
	ExampleParameterAttr string
	cancel               context.CancelFunc
}

// Constructor. This class implements component.Component
// but such implementations are not really needed because we are using the processorhelper
// package functions in the Factory.
func NewHelloWorldMetricsProcessor(logger *zap.Logger, cfg *Config) *helloWorldMetricsProcessor {
	return &helloWorldMetricsProcessor{
		logger:               logger,
		ExampleParameterAttr: cfg.ExampleParameterAttr,
	}
}

// implements https://pkg.go.dev/go.opentelemetry.io/collector/component#Component  Start
// Normally processors do not need this function at all
func (hwp *helloWorldMetricsProcessor) Start(ctx context.Context, host component.Host) error {
	hwp.logger.Info("Start!")
	ctx = context.Background()
	ctx, hwp.cancel = context.WithCancel(ctx)
	// This is an example of how to interact with host. This is not needed at all
	for k, _ := range host.GetExtensions() {
		hwp.logger.Info("Detected extension", zap.String("id", k.String()))
	}
	return nil
}

// implements https://pkg.go.dev/go.opentelemetry.io/collector/component#Component  Shutdown
// Normally processors do not need this function at all
func (hwp *helloWorldMetricsProcessor) Shutdown(ctx context.Context) error {
	hwp.logger.Info("Shutdown!")
	hwp.cancel()
	return nil
}

// https://pkg.go.dev/go.opentelemetry.io/collector/processor/processorhelper#ProcessMetricsFunc
// Main processor function
func (hwp *helloWorldMetricsProcessor) ProcessMetrics(ctx context.Context, md pmetric.Metrics) (pmetric.Metrics, error) {
	rms := md.ResourceMetrics()
	// iterate the ResourceMetricsSlice getting its attributes
	for i := 0; i < rms.Len(); i++ {
		attrs := rms.At(i).Resource().Attributes()
		if exampleParameterAttrValue, found := hwp.getAttrKey(attrs); found {
			hwp.logger.Info("Resource Attribute found!", zap.String("value", exampleParameterAttrValue.AsString()))
		}
	}
	return md, nil
}

// Finds the ExampleParameterAttr parameter in the metric attributes and returns its value
func (hwp *helloWorldMetricsProcessor) getAttrKey(attrs pcommon.Map) (result pcommon.Value, found bool) {
	found = false
	attrs.Range(func(k string, v pcommon.Value) bool {
		hwp.logger.Debug("Resource Attribute", zap.String("key", k), zap.String("type", v.Type().String()))
		if hwp.ExampleParameterAttr == k {
			result = v
			found = true
			// Returning false stops the iteration, but here we want to iterate to see all attributes
			//return false
		}
		return true
	})
	return
}
