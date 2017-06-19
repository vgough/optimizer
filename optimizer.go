package optimizer

// CostFN is the signature of a cost function which computes residuals from the
// given parameters.
//
// The residuals array should be filled out by the cost function.  The optimizer
// attempts to minimize the magnitude of the residuals.
//
// The optimizer will attempt to minimize the Euclidean norm (2-norm) of
// residual values, meaning that residuals closest to 0.0 are prefered.
// For the optimizer to work well, the residuals should vary smoothly with
// changes in parameters.
type CostFN func(params []float64, residuals []float64) error
