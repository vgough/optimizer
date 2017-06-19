#ifndef _OPTIMIZER_H_
#define _OPTIMIZER_H_

#ifdef __cplusplus
extern "C" {
#endif // __cplusplus

typedef struct {
    int num_params;
    int num_residuals;

    double relative_step_size; // 0 for default

    bool progress_to_stdout;
    bool print_summary; // Show result summary at end.

    void *internal;
} ceres_config;

void init_ceres();

ceres_config *create_config();
void delete_config(ceres_config *config);
void add_lower_bound(ceres_config *config, int param_num, double lower);
void add_upper_bound(ceres_config *config, int param_num, double upper);

// optimize runs the optimizer loop.  Returns when optimization is complete.
// cost_fn_is is an internal ID of a registered cost function.  This is passed
// in callbacks to GO, so that no function pointers are used.
//
// The parameter block is mutated and will hold the adjusted values when
// optimize returns.
int optimize(int cost_fn_id, ceres_config *config, double *params);

#ifdef __cplusplus
}
#endif // __cplusplus

#endif // _OPTIMIZER_H_
