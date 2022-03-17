package functional_test

import (
	"testing"

	functional "github.com/standoffvenus/functional/pkg"
	"github.com/stretchr/testify/assert"
)

func TestPanics(t *testing.T) {
	t.Run("Chain", func(t *testing.T) {
		panics(t, "FirstArgumentNil", func() { functional.Chain[int](nil) })
		panics(t, "NotFirstArgumentNil", func() {
			f := func(int) int { return 0 }
			functional.Chain(f, f, nil, f)
		})
	})

	t.Run("Compose", func(t *testing.T) {
		g := func(string) int { return 0 }
		f := func(int) float32 { return 0 }

		panics(t, "FNil", func() { functional.Compose[string, int, float32](nil, g) })
		panics(t, "GNil", func() { functional.Compose[string, int, float32](f, nil) })
	})

	panics(t, "FilterNil", func() { functional.Filter[int](nil, nil) })
	panics(t, "MapNil", func() { functional.Map[int, int](nil, nil) })
	panics(t, "ReduceNil", func() { functional.Reduce[int, int](nil, nil) })
}

func panics(t *testing.T, subTestName string, f func()) {
	t.Run(subTestName, func(t *testing.T) {
		assert.Panics(t, f)
	})
}
