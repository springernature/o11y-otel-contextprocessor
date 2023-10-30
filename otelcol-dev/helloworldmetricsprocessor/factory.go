package helloworldmetricsprocessor

import (
	"context"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/processor"
	"go.opentelemetry.io/collector/processor/processorhelper"
)

const (
	typeStr = "helloworld"
)

var (
	processorCapabilities                  = consumer.Capabilities{MutatesData: true}
	_                     component.Config = (*Config)(nil)
)

// Note: This isn't a valid configuration because the processor would do no work.
func createDefaultConfig() component.Config {
	return &Config{}
}

// NewFactory returns a new factory for the HelloWorld processor.
// This is a metrics only processor, alpha
func NewFactory() processor.Factory {
	return processor.NewFactory(
		typeStr,
		createDefaultConfig,
		processor.WithMetrics(createMetricsProcessor, component.StabilityLevelAlpha),
	//  This is a processor only for metrics!
	// 	processor.WithTraces(createTracesProcessor, metadata.TracesStability),
	// 	processor.WithLogs(createLogsProcessor, metadata.LogsStability))
	)
}

func createMetricsProcessor(
	ctx context.Context,
	set processor.CreateSettings,
	cfg component.Config,
	nextConsumer consumer.Metrics) (processor.Metrics, error) {

	// Get configuration
	helloworldProcessorCfg := cfg.(*Config)
	// Create instance
	helloworldProcessor := NewHelloWorldMetricsProcessor(set.Logger, helloworldProcessorCfg)
	// Use processorhelper package functions. No interface implementations are needed in this way.
	// Otherwise this interface will need to be implementedhttps://pkg.go.dev/go.opentelemetry.io/collector/consumer#Metrics
	// Set proper traces and call nextConsumer.
	return processorhelper.NewMetricsProcessor(
		ctx,
		set,
		cfg,
		nextConsumer,
		// Main function
		helloworldProcessor.ProcessMetrics,
		// Optional start function
		processorhelper.WithStart(helloworldProcessor.Start),
		// Optional Shutdown function
		processorhelper.WithShutdown(helloworldProcessor.Shutdown),
		processorhelper.WithCapabilities(processorCapabilities),
	)
}
