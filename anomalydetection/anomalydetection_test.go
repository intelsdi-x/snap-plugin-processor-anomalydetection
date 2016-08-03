//
// +build unit

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
	"math/rand"
	"testing"
	"time"

	"github.com/intelsdi-x/snap/control/plugin"
	"github.com/intelsdi-x/snap/control/plugin/cpolicy"
	"github.com/intelsdi-x/snap/core"
	"github.com/intelsdi-x/snap/core/ctypes"
	. "github.com/smartystreets/goconvey/convey"
)

//Random number generator
func randInt(min int, max int) int {
	return min + rand.Intn(max-min)
}

func TestAnomalyProcessor(t *testing.T) {
	meta := Meta()
	Convey("Meta should return metadata for the plugin", t, func() {
		Convey("So meta.Name should equal tag", func() {
			So(meta.Name, ShouldEqual, "anomalydetection")
		})
		Convey("So meta.Version should equal version", func() {
			So(meta.Version, ShouldEqual, version)
		})
		Convey("So meta.Type should be of type plugin.ProcessorPluginType", func() {
			So(meta.Type, ShouldResemble, plugin.ProcessorPluginType)
		})
	})

	proc := NewAnomalydetectionProcessor()
	Convey("Create tag processor", t, func() {
		Convey("So proc should not be nil", func() {
			So(proc, ShouldNotBeNil)
		})
		Convey("So proc should be of type tagProcessor", func() {
			So(proc, ShouldHaveSameTypeAs, &anomalyDetectionProcessor{})
		})
		Convey("proc.GetConfigPolicy should return a config policy", func() {
			configPolicy, _ := proc.GetConfigPolicy()
			Convey("So config policy should be a cpolicy.ConfigPolicy", func() {
				So(configPolicy, ShouldHaveSameTypeAs, &cpolicy.ConfigPolicy{})
			})

		})
	})
}

func TestAnomalyProcessorMetrics(t *testing.T) {
	Convey("Anomaly Processor tests", t, func() {
		metrics := make([]plugin.MetricType, 10)
		config := make(map[string]ctypes.ConfigValue)
		config["BufLength"] = ctypes.ConfigValueInt{Value: 10}
		config["Factor"] = ctypes.ConfigValueFloat{Value: 0.8}

		Convey("Check if data is transfered properly", func() {
			for i := range metrics {
				time.Sleep(3)
				rand.Seed(time.Now().UTC().UnixNano())
				data := randInt(65, 90)
				metrics[i] = *plugin.NewMetricType(core.NewNamespace("foo", "bar"), time.Now(), nil, "", data)
			}
			So(metrics[0].Tags_, ShouldBeNil)
			var buf bytes.Buffer
			enc := gob.NewEncoder(&buf)
			enc.Encode(metrics)
			So(metrics[0].Tags_, ShouldBeNil)
			anomalyObj := NewAnomalydetectionProcessor()

			_, receivedData, _ := anomalyObj.Process("snap.gob", buf.Bytes(), config)

			var metricsNew []plugin.MetricType

			//Decodes the content into pluginMetricType
			dec := gob.NewDecoder(bytes.NewBuffer(receivedData))
			dec.Decode(&metricsNew)
			So(metricsNew[0].Tags_, ShouldBeNil)
			So(metrics, ShouldNotResemble, metricsNew)

		})

	})
	Convey("Check Data integrity tests", t, func() {
		metrics := make([]plugin.MetricType, 9)
		config := make(map[string]ctypes.ConfigValue)
		config["BufLength"] = ctypes.ConfigValueInt{Value: 10}
		config["Factor"] = ctypes.ConfigValueFloat{Value: 0.8}

		Convey("Check if data have a proper format", func() {
			for i := range metrics {
				time.Sleep(3)
				rand.Seed(time.Now().UTC().UnixNano())
				data := randInt(65, 90)
				metrics[i] = *plugin.NewMetricType(core.NewNamespace("foo", "bar"), time.Now(), nil, "", data)
			}
			So(metrics[0].Tags_, ShouldBeNil)
			var buf bytes.Buffer
			enc := gob.NewEncoder(&buf)
			enc.Encode(metrics)
			So(metrics[0].Tags_, ShouldBeNil)
			anomalyObj := NewAnomalydetectionProcessor()

			_, receivedData, _ := anomalyObj.Process("snap.gob", buf.Bytes(), config)

			So(receivedData, ShouldNotBeNil)
			var expectedType float64
			for _, v := range anomalyObj.BufferMetric.Buffer["/foo/bar"].Metrics {
				So(v.Data_, ShouldHaveSameTypeAs, expectedType)
			}
			So(len(anomalyObj.BufferMetric.Buffer["/foo/bar"].Metrics), ShouldEqual, 9)

		})

	})
}
