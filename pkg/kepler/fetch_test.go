package kepler

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConvertJoulesTokWh(t *testing.T) {
	type testcase struct {
		name            string
		joules          float64
		intervalSeconds float64
		kWh             float64
		err             string
	}

	for _, test := range []*testcase{
		func() *testcase {
			return &testcase{
				name:            "valid conversion",
				joules:          1000,
				intervalSeconds: 300, // 5 minutes
				kWh:             9.259259259259259e-07,
			}
		}(),
		func() *testcase {
			return &testcase{
				name:            "valid conversion 1 minute",
				joules:          1000,
				intervalSeconds: 60, // 1 minute
				kWh:             4.62962962962963e-06,
			}
		}(),
		func() *testcase {
			return &testcase{
				name:            "valid conversion 372 J",
				joules:          372,
				intervalSeconds: 300, // 5 minutes
				kWh:             3.4444444444444444e-07,
			}
		}(),
		func() *testcase {
			return &testcase{
				name:            "error invalid 0 J",
				joules:          0,
				intervalSeconds: 300, // 5 minutes
				kWh:             0,
				err:             "energy consumption is 0",
			}
		}(),
	} {
		t.Run(test.name, func(t *testing.T) {
			actualKWh, err := convertJoulesTokWh(test.joules, test.intervalSeconds)
			assert.Equalf(t, test.kWh, actualKWh, "Result should be: %v, got: %v", test.kWh, actualKWh)
			if test.err != "" {
				assert.Errorf(t, err, test.err)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func TestGetRegionFromInstance(t *testing.T) {
	type testcase struct {
		name     string
		instance string
		region   string
		hasErr   bool
		expErr   string
	}

	for _, test := range []*testcase{
		func() *testcase {
			return &testcase{
				name:     "valid aws instance",
				instance: "ip-10-12-12-154.eu-central-1.compute.internal",
				region:   "eu-central-1",
				hasErr:   false,
			}
		}(),
		func() *testcase {
			return &testcase{
				name:     "valid gce instance",
				instance: "gke-gc0-europe-west1-default-f0c26727-1irq",
				region:   "europe-west1",
				hasErr:   false,
			}
		}(),
		func() *testcase {
			return &testcase{
				name:     "valid gce instance us-central1",
				instance: "gke-gc0-apps-us-central1-default-f0c26727-1irq",
				region:   "us-central1",
				hasErr:   false,
			}
		}(),
		func() *testcase {
			return &testcase{
				name:     "invalid instance",
				instance: "invalid-instance",
				region:   "",
				hasErr:   true,
				expErr:   "invalid instance",
			}
		}(),
		func() *testcase {
			return &testcase{
				name:     "invalid gcp instance missing zone",
				instance: "gke-gc0-apps-europe-west-medium-nodes-f09525f4-uokn",
				region:   "",
				hasErr:   true,
				expErr:   "invalid gcp instance",
			}
		}(),
	} {
		t.Run(test.name, func(t *testing.T) {
			actualRegion, err := getRegionFromInstance(test.instance)
			assert.Equalf(t, test.region, actualRegion, "Result should be: %v, got: %v", test.region, actualRegion)
			if test.hasErr {
				assert.Errorf(t, err, test.expErr)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}
