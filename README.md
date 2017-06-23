# optimizer

[![Go Report Card](https://goreportcard.com/badge/github.com/vgough/optimizer)](https://goreportcard.com/report/github.com/vgough/optimizer)
[![Apache 2.0 License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](LICENSE)

[optimizer](https://github.com/vgough/optimizer) - wrappers for numerical optimization.

## Goal

Provides a wrapper for [Ceres](http://ceres-solver.org/), which is a C++ library for solving large
optimization problems.  This package provides Go wrappers to enable optimization of Go functions.

## Limitations

* Numerical differentiation is the only supported mode.
* Limited to single parameter block (of any size).

## Example

```Go
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
    // takes multiple steps to find the answer, which can bee seen by requesting
    // verbose output.
    params := []float64{0.5}
    err := Adjust(params, 1, func(params []float64, residuals []float64) error {
        residuals[0] = 0.5 * math.Pow(10.0-params[0], 2)
        return nil
    }, WithVerboseOutput(true))
```

