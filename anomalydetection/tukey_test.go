//
// +build small

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
	"math/rand"
	"testing"
	"time"

	"github.com/intelsdi-x/snap/control/plugin"
	"github.com/intelsdi-x/snap/core"
	. "github.com/smartystreets/goconvey/convey"
)

func TestTukeyMethod(t *testing.T) {

	Convey("So all values should be a floats", t, func() {
		metrics := make([]plugin.MetricType, 10)
		someNumbers := []float64{100, 1, 2, 10, 10, 2, 3, 60, 80, 100}
		for i := range metrics {
			time.Sleep(3)
			rand.Seed(time.Now().UTC().UnixNano())
			data := someNumbers[i]
			metrics[i] = *plugin.NewMetricType(core.NewNamespace("foo", "bar"), time.Now(), nil, "", data)
		}
		values, err := unpackData(metrics)
		So(err, ShouldBeNil)
		So(len(values), ShouldResemble, 10)
		var expectedType float64
		for _, v := range values {
			So(v, ShouldHaveSameTypeAs, expectedType)
		}

	})
	Convey("So should return proper slice for round size of slice", t, func() {
		metrics := make([]plugin.MetricType, 6)
		someNumbers := []float64{99, 2, 25, 20, 2, 20}
		for i := range metrics {
			time.Sleep(3)
			rand.Seed(time.Now().UTC().UnixNano())
			data := someNumbers[i]
			metrics[i] = *plugin.NewMetricType(core.NewNamespace("foo", "bar"), time.Now(), nil, "", data)
		}
		values, _ := unpackData(metrics)
		factor := 0.8
		value, outlayerIndex := getOutliers(values, factor)
		expectedVal := 41.0
		expencedSize := 4
		So(value, ShouldResemble, expectedVal)
		So(len(outlayerIndex), ShouldResemble, expencedSize)
		So(outlayerIndex, ShouldNotBeEmpty)
	})
	Convey("So should return proper slice for even size of slice", t, func() {
		metrics := make([]plugin.MetricType, 7)
		someNumbers := []float64{99, 2, 25, 20, 2, 20, 60}
		for i := range metrics {
			time.Sleep(3)
			rand.Seed(time.Now().UTC().UnixNano())
			data := someNumbers[i]
			metrics[i] = *plugin.NewMetricType(core.NewNamespace("foo", "bar"), time.Now(), nil, "", data)
		}

		values, _ := unpackData(metrics)

		factor := 0.8
		value, outlayerIndex := getOutliers(values, factor)

		expectedVal := 61.333333333333336
		expencedSize := 3

		So(value, ShouldResemble, expectedVal)
		So(outlayerIndex, ShouldNotBeEmpty)
		So(len(outlayerIndex), ShouldResemble, expencedSize)

	})
	Convey("So interfacetoString should return always float64", t, func() {

		var (
			integer          int
			integer16        int16
			integer32        int32
			integer64        int64
			uinteger         uint
			uinteger16       uint16
			uinteger32       uint32
			uinteger64       uint64
			float_32         float32
			float_64         float64
			valid_string     string
			non_valid_string string
			expectedValue    float64
		)

		value, err := interfaceToFloat(integer)
		So(value, ShouldHaveSameTypeAs, expectedValue)
		So(err, ShouldBeNil)
		value, err = interfaceToFloat(integer16)
		So(value, ShouldHaveSameTypeAs, expectedValue)
		So(err, ShouldBeNil)
		value, err = interfaceToFloat(integer32)
		So(value, ShouldHaveSameTypeAs, expectedValue)
		So(err, ShouldBeNil)
		value, err = interfaceToFloat(integer64)
		So(value, ShouldHaveSameTypeAs, expectedValue)
		So(err, ShouldBeNil)
		value, err = interfaceToFloat(uinteger)
		So(value, ShouldHaveSameTypeAs, expectedValue)
		So(err, ShouldBeNil)
		value, err = interfaceToFloat(uinteger16)
		So(value, ShouldHaveSameTypeAs, expectedValue)
		So(err, ShouldBeNil)
		value, err = interfaceToFloat(uinteger32)
		So(value, ShouldHaveSameTypeAs, expectedValue)
		So(err, ShouldBeNil)
		value, err = interfaceToFloat(uinteger64)
		So(value, ShouldHaveSameTypeAs, expectedValue)
		So(err, ShouldBeNil)
		value, err = interfaceToFloat(float_32)
		So(value, ShouldHaveSameTypeAs, expectedValue)
		So(err, ShouldBeNil)
		value, err = interfaceToFloat(float_64)
		So(value, ShouldHaveSameTypeAs, expectedValue)
		So(err, ShouldBeNil)
		So(err, ShouldBeNil)
		valid_string = "10"
		value, err = interfaceToFloat(valid_string)
		So(value, ShouldHaveSameTypeAs, expectedValue)
		So(err, ShouldBeNil)
		non_valid_string = "error"
		value, err = interfaceToFloat(non_valid_string)
		So(value, ShouldHaveSameTypeAs, expectedValue)
		So(err, ShouldNotBeEmpty)
	})

}
