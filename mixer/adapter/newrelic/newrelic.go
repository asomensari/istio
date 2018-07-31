//go:generate $GOPATH/src/istio.io/istio/bin/mixer_codegen.sh -f mixer/adapter/newrelic/config/config.proto

package newrelic

import (
	// "github.com/gogo/protobuf/types"
	"context"
	//"os"
	"fmt"
	"istio.io/istio/mixer/adapter/newrelic/config"
	"istio.io/istio/mixer/pkg/adapter"
	"istio.io/istio/mixer/template/metric"

	"go-insights/client"
)

type (
	builder struct {
		adpCfg      *config.Params
		metricTypes map[string]*metric.Type
	}

	handler struct {
		client      *client.InsertClient
		cfg         *config.Params
		metricTypes map[string]*metric.Type
		env         adapter.Env
	}
)

// ensure types implement the requisite interfaces
var _ metric.HandlerBuilder = &builder{}
var _ metric.Handler = &handler{}

///////////////// Configuration-time Methods ///////////////

// adapter.HandlerBuilder#Build
func (b *builder) Build(ctx context.Context, env adapter.Env) (adapter.Handler, error) {

	var err error

	fmt.Println("-------> Build ")
	cli := client.NewInsertClient(b.adpCfg.InsightsKey, b.adpCfg.NewrelicId)

	fmt.Printf("Insights Key:%s, New Relic Account: %s \n", b.adpCfg.InsightsKey, b.adpCfg.NewrelicId)

	return &handler{client: cli, metricTypes: b.metricTypes, env: env}, err

}

// adapter.HandlerBuilder#SetAdapterConfig
func (b *builder) SetAdapterConfig(cfg adapter.Config) {

	fmt.Println("------->  SetAdapterConfig")
	b.adpCfg = cfg.(*config.Params)

}

// adapter.HandlerBuilder#Validate
func (b *builder) Validate() (ce *adapter.ConfigErrors) {
	fmt.Println("-------> Validate")

	if b.adpCfg.InsightsKey == "" {
		ce = ce.Append("InsightsKey", fmt.Errorf("insights_key cannot be empty"))
	}
	if b.adpCfg.NewrelicId == "" {
		ce = ce.Append("NewrelicId", fmt.Errorf("newrelic_id cannot be empty"))
	}

	return ce

}

// metric.HandlerBuilder#SetMetricTypes
func (b *builder) SetMetricTypes(types map[string]*metric.Type) {
	fmt.Println("------->  SetMetricTypes")
	b.metricTypes = types

}

////////////////// Request-time Methods //////////////////////////
// metric.Handler#HandleMetric
func (h *handler) HandleMetric(ctx context.Context, insts []*metric.Instance) error {

	fmt.Println("------->  At Handle Metric")

	metricData := h.generateMetricData(insts)
	_, err := h.sendMetricsToNewRelic(metricData)

	if err != nil {
		return fmt.Errorf("Error posting to Insights %s", metricData)
	}
	return nil
}

// adapter.Handler#Close
func (h *handler) Close() error {
	fmt.Println("------->  Close ")
	//return  h.client.Flush()
	//client.InsertClient{}//h.f.Close()
	return nil
}

////////////////// Bootstrap //////////////////////////
// GetInfo returns the adapter.Info specific to this adapter.
func GetInfo() adapter.Info {
	return adapter.Info{
		Name:        "newrelic",
		Description: "Logs the metric calls into New Relic",
		SupportedTemplates: []string{
			metric.TemplateName,
		},
		NewBuilder:    func() adapter.HandlerBuilder { return &builder{} },
		DefaultConfig: &config.Params{},
	}
}
