package metrics

import (
	"time"

	"github.com/prometheus/common/model"
)

type Sample struct {
	Metric map[string]string
	Value  float64
	Time   time.Time
}

func metricToMap(m model.Metric) map[string]string {
	result := make(map[string]string)
	for k, v := range m {
		result[string(k)] = string(v)
	}
	return result
}

func FromPromSample(s *model.Sample) Sample {
	return Sample{
		Metric: metricToMap(s.Metric),
		Value: float64(s.Value),
		Time: s.Timestamp.Time(),
	}
}

func FromPromSampleStream(ss *model.SampleStream) []Sample {
	result := make([]Sample, len(ss.Values))

	metric := metricToMap(ss.Metric)

	for i, v := range ss.Values {
		result[i] = Sample{
			Metric: metric,
			Value: float64(v.Value),
			Time: v.Timestamp.Time(),
		}
	}

	return result
}
