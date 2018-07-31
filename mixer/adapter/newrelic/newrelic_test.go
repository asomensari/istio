// Copyright 2018 Istio Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package newrelic

import (
	"fmt"
	"testing"

	"io/ioutil"
	"istio.io/istio/mixer/pkg/adapter"
	adapter_integration "istio.io/istio/mixer/pkg/adapter/test"
	"os"
	"strings"
)

func TestReport(t *testing.T) {
	operatorCfgBytes, err := ioutil.ReadFile("testdata/newrelic.yaml")
	if err != nil {
		t.Fatalf("could not read file: %v", err)
	}
	operatorCfg := string(operatorCfgBytes)

	defer func() {
		if removeErr := os.Remove("out.txt"); removeErr != nil {
			t.Logf("Could not remove temporary file %s: %v", "out.txt", removeErr)
		}
	}()

	adapter_integration.RunTest(
		t,
		func() adapter.Info {
			return GetInfo()
		},
		adapter_integration.Scenario{
			ParallelCalls: []adapter_integration.Call{
				{
					CallKind: adapter_integration.REPORT,
				},
			},

			GetState: func(ctx interface{}) (interface{}, error) {
				// validate if the content of "out.txt" is as expected
				bytes, err := ioutil.ReadFile("out.txt")
				if err != nil {
					return nil, err
				}
				s := string(bytes)
				wantStr := `
HandleMetric invoke for :
				Instance Name  :'requestcount.metric.istio-system'
				Instance Value : {requestcount.metric.istio-system 1 map[target:svc.cluster.local]  map[]},
				Type           : {INT64 map[target:STRING] map[]}
`
				if normalize(s) != normalize(wantStr) {
					return nil, fmt.Errorf("got adapters state as : '%s'; want '%s'", s, wantStr)
				}
				return nil, nil
			},

			GetConfig: func(ctx interface{}) ([]string, error) {
				return []string{
					operatorCfg,
				}, nil
			},

			Want: `
            		{
		 "AdapterState": null,
		 "Returns": [
		  {
		   "Check": {
		    "Status": {},
		    "ValidDuration": 0,
		    "ValidUseCount": 0
		   },
		   "Quota": null,
		   "Error": null
		  }
		 ]
		}`,
		},
	)
}

func normalize(s string) string {
	s = strings.TrimSpace(s)
	s = strings.Replace(s, "\t", "", -1)
	s = strings.Replace(s, "\n", "", -1)
	s = strings.Replace(s, " ", "", -1)
	return s
}
