package kepler

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	prometheus "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
	"golang.org/x/net/context"

	v1 "github.com/re-cinq/aether/pkg/types/v1"
)

const query = "sum without (mode) (rate(kepler_container_%s_joules_total[%s]))"

// Stop is used to adhere to the interface and will be called when Aether is
// shut down
func (k *KeplerSource) Stop(ctx context.Context) error {
	return nil
}

// Fetch is what gets run by aether, it should return a list of instances with
// the metrics attached to them (CPU, Memory, Networking, Storage)
func (k *KeplerSource) Fetch(ctx context.Context) ([]*v1.Instance, error) {

	// Get API for the Prometheus client
	k.v1API = prometheus.NewAPI(*k.client)

	// Get CPU energy consumption for each container
	err := k.cpuMetrics(ctx, fmt.Sprintf(query, "core", k.cfg.Interval.String()))
	if err != nil {
		return nil, err
	}

	// Get Memory energy consumption for each container
	err = k.memMetrics(ctx, fmt.Sprintf(query, "dram", k.cfg.Interval.String()))
	if err != nil {
		return nil, err
	}

	// The Fetch signature requires a slice as a return
	instances := []*v1.Instance{}
	for _, i := range k.instancesMap {
		instances = append(instances, i)
	}

	k.logger.Debug("kepler source: fetched instances", "instance Count", len(instances))
	return instances, nil
}

func (k *KeplerSource) cpuMetrics(ctx context.Context, query string) error {
	//r := prometheus.Range{
	//	Start: time.Now().Add(-k.cfg.Interval),
	//	End:   time.Now(),
	//	Step:  time.Minute,
	//}
	result, warnings, err := k.v1API.Query(ctx, query, time.Now())
	if err != nil {
		k.logger.Error("kelper source: error querying prometheus", "error", err)
		return err
	}

	if len(warnings) > 0 {
		k.logger.Warn("kepler source: prometheus query warnings", "warnings", warnings)
	}

	// Check if the result is a model.Vector
	vector, ok := result.(model.Vector)
	if !ok {
		return fmt.Errorf("expected a model.Vector but got %T", result)
	}

	for _, sample := range vector {
		// store labels to be added to metric and construct the instance
		labels := make(map[string]string)
		for k, v := range sample.Metric {
			labels[string(k)] = string(v)
		}

		if sample.Value == 0 {
			k.logger.Debug("kepler source: CPU energy consumption is 0", "containerID", labels["container_id"], "container", labels["container_name"])
			continue
		}

		// Convert the Joules energy consumption collected from kepler
		// into KwH
		// 1 Joule = 1 Watt second, so divide Joules by the interval of time in seconds
		// and divide by 3600 to get watt hours, and divide by 1000 to get kilowatt hours.
		ekWh := float64(sample.Value) / k.cfg.Interval.Seconds() / 3600 / 1000

		// Create the CPU metric with the energy consumtion
		// collected from the query
		m := v1.NewMetric(v1.CPU.String())
		m.ResourceType = v1.CPU
		m.Energy = ekWh
		m.Labels = labels

		// since container_id is unique, we can use it as the instanceMap key
		// and ID for the instance
		id := labels["container_id"]
		name := labels["container_name"]

		k.logger.Debug("kepler source: CPU energy consumption kWh", "containerID", id, "container", name, "energy", m.Energy)

		// if the instance already exists in the instancesMap
		// upsert the metric and continue looping
		if instance, exists := k.instancesMap[id]; exists {
			instance.Metrics.Upsert(m)
			continue
		}

		// if the instance doesn't already exist, create one
		// from the labels provided
		region, exists := labels["region"]
		if !exists {
			region, err = getRegionFromInstance(labels["instance"])
			if err != nil {
				k.logger.Error("kepler source: error getting region from instance", "error", err)
				continue
			}
		}

		p, ok := v1.Providers[k.cfg.Provider]
		if !ok {
			return fmt.Errorf("provider %s not found", k.cfg.Provider)
		}

		instance := &v1.Instance{
			ID:       id,
			Provider: p,
			Service:  "kepler",
			Name:     name,
			Region:   region,
			Status:   v1.InstanceRunning,
			Labels:   labels,
		}

		instance.Metrics.Upsert(m)
		k.instancesMap[id] = instance
	}

	return nil
}

func (k *KeplerSource) memMetrics(ctx context.Context, query string) error {
	//r := prometheus.Range{
	//	Start: time.Now().Add(-k.cfg.Interval),
	//	End:   time.Now(),
	//	Step:  time.Minute,
	//}
	result, warnings, err := k.v1API.Query(ctx, query, time.Now())
	if err != nil {
		k.logger.Error("kelper source: error querying prometheus", "error", err)
		return err
	}

	if len(warnings) > 0 {
		k.logger.Warn("kepler source: prometheus query warnings", "warnings", warnings)
	}

	// Check if the result is a model.Vector
	vector, ok := result.(model.Vector)
	if !ok {
		return fmt.Errorf("expected a model.Vector but got %T", result)
	}

	for _, sample := range vector {
		labels := make(map[string]string)
		for k, v := range sample.Metric {
			labels[string(k)] = string(v)
		}

		if sample.Value == 0 {
			k.logger.Debug("kepler source: memory energy consumption is 0", "containerID", labels["container_id"], "container", labels["container_name"])
			continue
		}

		// Convert the Joules energy consumption collected from kepler
		// into KwH
		// 1 Joule = 1 Watt second, so divide Joules by the interval of time in seconds
		// and divide by 3600 to get watt hours, and divide by 1000 to get kilowatt hours.
		ekWh := float64(sample.Value) / k.cfg.Interval.Seconds() / 3600 / 1000

		// Create the memory metric with the energy consumtion
		// collected from the query
		m := v1.NewMetric(v1.Memory.String())
		m.ResourceType = v1.Memory
		m.Energy = ekWh
		m.Labels = labels

		// since container_id is unique, we can use it as the instanceMap key
		// and ID for the instance
		id := labels["container_id"]
		name := labels["container_name"]

		k.logger.Debug("kepler source: memory energy consumption kWh", "containerID", id, "container", name, "energy", m.Energy)

		// if the instance already exists in the instancesMap
		// upsert the metric and continue looping
		if instance, exists := k.instancesMap[id]; exists {
			instance.Metrics.Upsert(m)
			continue
		}

		// if the instance doesn't yet exist, create it from
		// the labels
		region, exists := labels["region"]
		if !exists {
			region, err = getRegionFromInstance(labels["instance"])
			if err != nil {
				k.logger.Error("kepler source: error getting region from instance", "error", err)
				continue
			}
		}

		p, ok := v1.Providers[k.cfg.Provider]
		if !ok {
			return fmt.Errorf("provider %s not found", k.cfg.Provider)
		}

		instance := &v1.Instance{
			ID:       id,
			Provider: p,
			Service:  "kepler",
			Name:     name,
			Status:   v1.InstanceRunning,
			Region:   region,
			Labels:   labels,
		}

		instance.Metrics.Upsert(m)
		k.instancesMap[id] = instance
	}

	return nil
}

// getRegionFromInstance takes a prometheus metrics label of
// instance and extracts the region the cluster is running in
// from the string. For example:
// gke-gc0-apps-europe-west1-default-ca15cfa4-hrb7 => europe-west1
// ip-192-168-29-157.eu-north-1.compute.internal => eu-north-1
func getRegionFromInstance(s string) (string, error) {
	if strings.HasPrefix(s, "gke-") {
		// The regex is looking for the region in the instance name
		re := regexp.MustCompile("(europe|asia|australia|southamerica|me|africa|us)-[a-z]+[0-9]+")
		match := re.FindString(s)
		if match == "" {
			return "", fmt.Errorf("unable to get region from instance: %s", s)
		}
		return match, nil
	}

	// Handle AWS instance names
	if strings.Contains(s, ".") {
		parts := strings.Split(s, ".")
		// The region is typically the second part
		if len(parts) >= 2 {
			return parts[1], nil
		}
	}
	return "", fmt.Errorf("unable to get region from instance: %s", s)
}
