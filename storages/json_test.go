// PromHouse
// Copyright (C) 2017 Percona LLC
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published
// by the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program. If not, see <http://www.gnu.org/licenses/>.

package storages

import (
	"encoding/json"
	"testing"

	"github.com/prometheus/common/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	prom2 "github.com/Percona-Lab/PromHouse/prompb/prom2"
)

func TestMarshalMetricsAndLabels(t *testing.T) {
	for _, labels := range [][]*prom2.Label{
		{
			&prom2.Label{Name: "__name__", Value: "normal"},
			&prom2.Label{Name: "instance", Value: "foo"},
			&prom2.Label{Name: "job", Value: "bar"},
		}, {
			&prom2.Label{Name: "__name__", Value: "funny_1"},
			&prom2.Label{Name: "label", Value: ""},
		}, {
			&prom2.Label{Name: "__name__", Value: "funny_2"},
			&prom2.Label{Name: "label", Value: "'`\"\\"},
		}, {
			&prom2.Label{Name: "__name__", Value: "funny_3"},
			&prom2.Label{Name: "label", Value: "''``\"\"\\\\"},
		}, {
			&prom2.Label{Name: "__name__", Value: "funny_4"},
			&prom2.Label{Name: "label", Value: "'''```\"\"\"\\\\\\"},
		}, {
			&prom2.Label{Name: "__name__", Value: "funny_5"},
			&prom2.Label{Name: "label", Value: `\ \\ \\\\ \\\\`},
		}, {
			&prom2.Label{Name: "__name__", Value: "funny_6"},
			&prom2.Label{Name: "label", Value: "🆗"},
		},
	} {
		b1 := marshalLabels(labels, nil)
		b2, err := json.Marshal(makeMetric(labels))
		require.NoError(t, err)

		m1 := make(model.Metric)
		require.NoError(t, json.Unmarshal(b1, &m1))
		m2 := make(model.Metric)
		require.NoError(t, json.Unmarshal(b2, &m2))
		assert.Equal(t, m2, m1)

		l1, err := unmarshalLabels(b1)
		require.NoError(t, err)
		l2, err := unmarshalLabels(b2)
		require.NoError(t, err)
		sortLabels(l1)
		sortLabels(l2)
		assert.Equal(t, labels, l1)
		assert.Equal(t, labels, l2)
	}
}

var sink []byte

func BenchmarkMarshalJSON(b *testing.B) {
	var err error
	metric := makeMetric(labelsB)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sink, err = json.Marshal(metric)
	}
	b.StopTimer()

	require.NoError(b, err)
}

func BenchmarkMarshalLabels(b *testing.B) {
	for i := 0; i < b.N; i++ {
		sink = marshalLabels(labelsB, sink[:0])
	}
}
