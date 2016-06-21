/*
http://www.apache.org/licenses/LICENSE-2.0.txt

Copyright 2016 Intel Corporation

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package anomalydetection

import (
	"bytes"
	"encoding/gob"

	log "github.com/Sirupsen/logrus"

	"github.com/intelsdi-x/snap/control/plugin"
	"github.com/intelsdi-x/snap/control/plugin/cpolicy"
	"github.com/intelsdi-x/snap/core/ctypes"
)

const (
	name                = "anomalydetection"
	version             = 1
	pluginType          = plugin.ProcessorPluginType
	defaultBufferLength = 30
	defaultFactor       = 3.0
)

// Buffer struct stores []plugin.MetricType for specific namespace
type Buffer struct {
	Metrics []plugin.MetricType
}

// BufferMetric struct, stores all Buffers
type BufferMetric struct {
	Buffer map[string]*Buffer
}

// Meta returns a plugin meta data
func Meta() *plugin.PluginMeta {
	return plugin.NewPluginMeta(name, version, pluginType, []string{plugin.SnapGOBContentType}, []string{plugin.SnapGOBContentType})
}

// NewAnomalydetectionProcessor creates new processor
func NewAnomalydetectionProcessor() *anomalyDetectionProcessor {
	buffer := make(map[string]*Buffer, 30)
	return &anomalyDetectionProcessor{
		BufferMetric: BufferMetric{
			Buffer: buffer,
		},
	}
}

type anomalyDetectionProcessor struct {
	BufferMetric BufferMetric
}

func (p *anomalyDetectionProcessor) addToBuffer(m plugin.MetricType, logger *log.Logger) error {
	ns := m.Namespace().String()
	m, err := dataToFloat(m)
	if err != nil {
		return err
	}
	if _, ok := p.BufferMetric.Buffer[ns]; ok {
		p.BufferMetric.Buffer[ns].Metrics = append(p.BufferMetric.Buffer[ns].Metrics, m)

	} else {
		vMet := []plugin.MetricType{m}
		p.BufferMetric.Buffer[ns] = &Buffer{
			Metrics: vMet,
		}
	}
	logger.Debug("Buffer lenght: ", len(p.BufferMetric.Buffer[ns].Metrics))
	return nil

}

func (p *anomalyDetectionProcessor) clearBuffer(ns string) {

	vMet := []plugin.MetricType{}
	p.BufferMetric.Buffer[ns] = &Buffer{
		Metrics: vMet,
	}
}

func (p *anomalyDetectionProcessor) getBuffer(ns string) []plugin.MetricType {

	return p.BufferMetric.Buffer[ns].Metrics
}

func (p *anomalyDetectionProcessor) getBufferLength(ns string) int {

	return len(p.BufferMetric.Buffer[ns].Metrics)
}

func (p *anomalyDetectionProcessor) calculateTukeyMethod(m plugin.MetricType, factor float64, logger *log.Logger) ([]plugin.MetricType, error) {

	ns := m.Namespace().String()
	metrics := p.getBuffer(ns)
	values, err := unpackData(metrics)
	if err != nil {
		return nil, err
	}
	_, outliersIndex := getOutliers(values, factor)

	ret := []plugin.MetricType{}
	for _, v := range outliersIndex {
		ret = append(ret, metrics[v])
	}
	return ret, nil

}

func handleErr(e error) {
	if e != nil {
		panic(e)
	}
}

func (p *anomalyDetectionProcessor) GetConfigPolicy() (*cpolicy.ConfigPolicy, error) {
	cp := cpolicy.New()
	config := cpolicy.NewPolicyNode()
	r1, err := cpolicy.NewIntegerRule("BufLength", true, defaultBufferLength)
	handleErr(err)
	r1.Description = "Buffer Length for tukey method "
	config.Add(r1)
	r2, err := cpolicy.NewFloatRule("Factor", false, defaultFactor)
	handleErr(err)
	r2.Description = "Buffer Length for tukey method "
	config.Add(r2)
	cp.Add([]string{""}, config)
	return cp, nil
}

func (p *anomalyDetectionProcessor) Process(contentType string, content []byte, config map[string]ctypes.ConfigValue) (string, []byte, error) {
	var (
		metrics, metricsTemp []plugin.MetricType
		bufferLength         int
		factor               float64
	)

	logger := log.New()
	logger.Debug("anomalyDetection Processor started")

	if config["BufLength"].(ctypes.ConfigValueInt).Value > 0 {
		bufferLength = config["BufLength"].(ctypes.ConfigValueInt).Value

	}
	if config["Factor"].(ctypes.ConfigValueFloat).Value > 0 {
		factor = config["Factor"].(ctypes.ConfigValueFloat).Value

	}
	//Decodes the content into MetricType
	dec := gob.NewDecoder(bytes.NewBuffer(content))
	if err := dec.Decode(&metrics); err != nil {
		logger.Printf("Error decoding: error=%v content=%v", err, content)
		return "", nil, err
	}

	for _, m := range metrics {

		ns := m.Namespace().String()
		if _, ok := p.BufferMetric.Buffer[ns]; ok {
			if p.getBufferLength(ns) == bufferLength-1 {
				mVal, err := p.calculateTukeyMethod(m, factor, logger)
				if err != nil {
					return "", nil, err
				}
				metricsTemp = append(metricsTemp, mVal...)
				p.clearBuffer(ns)

			} else {
				p.addToBuffer(m, logger)
			}
		} else {
			p.addToBuffer(m, logger)
		}

	}
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	enc.Encode(metricsTemp)
	return contentType, buf.Bytes(), nil
}
