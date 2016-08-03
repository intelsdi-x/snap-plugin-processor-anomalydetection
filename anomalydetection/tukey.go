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
	"errors"
	"strconv"

	"github.com/intelsdi-x/snap/control/plugin"
)

func parseValues(values []float64, q1 float64, q3 float64, factor float64) (float64, []int) {
	var (
		outliers []int
		value    float64
	)
	fence1 := q1 - factor*(q3-q1)
	fence2 := q3 + factor*(q3-q1)
	for i, v := range values {
		if v < fence1 || v > fence2 {
			value = value + v
			outliers = append(outliers, i)

		}
	}

	if len(outliers) != 0 {
		return value / float64(len(outliers)), outliers
	}
	return 0.0, outliers
}

func getOutliers(values []float64, factor float64) (float64, []int) {

	l := len(values)

	if l%2 == 0 {
		q1 := values[l/4]
		q3 := values[3*l/4]

		return parseValues(values, q1, q3, factor)

	}
	i := values[l/4]
	j := values[l/4+1]
	q1 := i + (j-i)*0.25
	i = values[3*l/4-1]
	j = values[3*l/4]
	q3 := i + (j-i)*0.75
	return parseValues(values, q1, q3, factor)

}

func interfaceToFloat(face interface{}) (float64, error) {
	var (
		ret float64
		err error
	)
	switch val := face.(type) {
	case string:
		ret, err = strconv.ParseFloat(val, 64)
	case int:
		ret = float64(val)
	case int16:
		ret = float64(val)
	case int32:
		ret = float64(val)
	case int64:
		ret = float64(val)
	case uint:
		ret = float64(val)
	case uint16:
		ret = float64(val)
	case uint32:
		ret = float64(val)
	case uint64:
		ret = float64(val)
	case float32:
		ret = float64(val)
	case float64:
		ret = val

	default:
		err = errors.New("unsupported type")
	}
	return ret, err
}

func unpackData(values []plugin.MetricType) ([]float64, error) {
	metrics := []float64{}
	for _, v := range values {
		metrics = append(metrics, v.Data_.(float64))
	}
	return metrics, nil
}

func dataToFloat(m plugin.MetricType) (plugin.MetricType, error) {
	var err error
	m.Data_, err = interfaceToFloat(m.Data_)
	return m, err
}
