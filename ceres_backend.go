package optimizer

// #cgo pkg-config: eigen3
// #cgo LDFLAGS: -lceres -lglog
// #include <stdbool.h>
// #include "ceres_impl.h"
import "C"
import (
	"errors"
	"log"
	"sync"
	"unsafe"
)

// CeresOpt is an option for the ceres backend.
type CeresOpt func(c *C.ceres_config) error

var ceresInit sync.Once

// WithRelativeStepSize sets the relative step size for initial parameter
// adjustments.  The backend may modify the step size itself.
func WithRelativeStepSize(rel float64) CeresOpt {
	return func(cfg *C.ceres_config) error {
		cfg.relative_step_size = C.double(rel)
		return nil
	}
}

// WithVerboseOutput requests verbose logging / output from the backend..
func WithVerboseOutput(verbose bool) CeresOpt {
	return func(cfg *C.ceres_config) error {
		cfg.progress_to_stdout = C.bool(verbose)
		cfg.print_summary = C.bool(verbose)
		return nil
	}
}

// WithLowerBound sets a lower bound on a parameter.
func WithLowerBound(paramNum int, lower float64) CeresOpt {
	return func(cfg *C.ceres_config) error {
		if paramNum < 0 || paramNum >= int(cfg.num_params) {
			return errors.New("parameter out of bounds")
		}
		C.add_lower_bound(cfg, C.int(paramNum), C.double(lower))
		return nil
	}
}

// WithUpperBound sets an upper bound on a parameter.
func WithUpperBound(paramNum int, upper float64) CeresOpt {
	return func(cfg *C.ceres_config) error {
		if paramNum < 0 || paramNum >= int(cfg.num_params) {
			return errors.New("parameter out of bounds")
		}
		C.add_upper_bound(cfg, C.int(paramNum), C.double(upper))
		return nil
	}
}

func Adjust(params []float64, numResiduals int, fn CostFN, opts ...CeresOpt) error {
	ceresInit.Do(func() {
		C.init_ceres()
	})

	numParams := len(params)
	var err error
	id := callbackRegistry.register(func(p *C.double, r *C.double) C.bool {
		// Map inputs to Go slices.
		pin := (*[1 << 20]float64)(unsafe.Pointer(p))[:numParams:numParams]
		out := (*[1 << 20]float64)(unsafe.Pointer(r))[:numResiduals:numResiduals]
		err = fn(pin, out)
		return err == nil
	})
	defer callbackRegistry.deregister(id)

	// Create ceres backend config.
	ccfg := C.create_config()
	ccfg.num_params = C.int(numParams)
	ccfg.num_residuals = C.int(numResiduals)
	for _, o := range opts {
		if err := o(ccfg); err != nil {
			return err
		}
	}
	defer C.delete_config(ccfg)

	// Pass control to C++ functions.
	code := C.optimize(C.int(id),
		ccfg,
		(*C.double)(&params[0]),
	)
	if code != 0 {
		return errors.New("optimization failed")
	}
	return nil
}

// ComputeResidual is exported to C for evaluation of parameter fit.
//export ComputeResidual
func ComputeResidual(id C.int, params *C.double, residuals *C.double) C.bool {
	fn := callbackRegistry.lookup(id)
	if fn == nil {
		log.Printf("no callback registered for id %v", id)
		return false
	}

	return fn(params, residuals)
}
