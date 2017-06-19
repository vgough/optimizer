package optimizer

import (
	"fmt"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOptimizer(t *testing.T) {
	assert := assert.New(t)

	// This is a port of Ceres' HelloWorld example, which minimizes the equation
	// `f(x)=0.5 * (10-x) ^ 2`
	//
	// The optimizer will attempt to minimize the 'residual' or error term
	// if we were to formulate as `x: 0 = 0.5 * (10-x) ^ 2
	//
	// The solution is `x=10`, since `f(x)` is a simple parabola:
	//      https://www.wolframalpha.com/input/?i=0.5+*+(10-x)+%5E+2
	//
	// See also: http://ceres-solver.org/nnls_tutorial.html#hello-world
	//
	// In this example, the initial value of x is given as 0.5, which should
	// become 10.0 after adjustment.
	//
	// By default, the ceres backend is optimized for non-linear search, so it
	// takes multiple steps to find the answer, which can bee seen by running
	// the test in verbose mode.
	params := []float64{0.5}
	err := Adjust(params, 1, func(params []float64, residuals []float64) error {
		residuals[0] = 0.5 * math.Pow(10.0-params[0], 2)
		return nil
	}, WithVerboseOutput(true))

	assert.NoError(err)
	assert.InDelta(10.0, params[0], 1e-3)
}

func ExampleAdjust() {
	params := []float64{0.5}
	err := Adjust(params, 1, func(params []float64, residuals []float64) error {
		residuals[0] = 0.5 * math.Pow(10.0-params[0], 2)
		return nil
	})
	if err != nil {
		fmt.Printf("failed: %s", err)
	}
}
