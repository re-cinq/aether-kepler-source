package kepler

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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
