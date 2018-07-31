package newrelic

import (
	"fmt"
	"istio.io/istio/mixer/template/metric"
	"strconv"
)

func (h *handler) sendMetricsToNewRelic(eventData map[string]interface{}) (int, error) {

	fmt.Println("------->  sendMetricToNewRelic()")

	h.client.PostEvent(eventData)

	return 0, nil
}

func (h *handler) generateMetricData(insts []*metric.Instance) map[string]interface{} {

	fmt.Println("-------> generateMetricData")

	dataMap := map[string]interface{}{
		"eventType": "Istio",
	}

	for _, metricInstance := range insts {
		dataMap["metricInstance.name"] = metricInstance.Name

		value, err := getNumericValue(metricInstance.Value)
		if err != nil {
			//not failing, moving on
			_ = h.env.Logger().Errorf("could not parse value %v", err)
		}
		dataMap["metricInstance.value"] = value

		for key, val := range metricInstance.Dimensions {
			dimensionLabel := "metricInstance.dimension." + key
			dataMap[dimensionLabel] = val
			fmt.Printf("metricInstance.dimension: %s - %s \n", dimensionLabel, val)
		}

		dataMap["monitoredResourceType"] = metricInstance.MonitoredResourceType

		for key, val := range metricInstance.MonitoredResourceType {
			resourceTypeKey := "metricInstance.resourceType." + string(key)
			dataMap[resourceTypeKey] = val
			fmt.Printf("metricInstance.resourceType: %s - %s \n", resourceTypeKey, val)
		}
	}
	return dataMap
}

func getNumericValue(numValue interface{}) (float64, error) {

	switch val := numValue.(type) {

	case string:
		value, err := strconv.ParseFloat(val, 64)
		if err != nil {
			return 0, fmt.Errorf("can't parse string %s into float", val)
		}
		return value, nil

	case int:
		return float64(val), nil

	case int64:
		return float64(val), nil

	case float32:
		return float64(val), nil

	case float64:
		return float64(val), nil

	default:
		return 0, fmt.Errorf("unsupported value type %T. Only strings and numbers allowed", val)
	}
}
