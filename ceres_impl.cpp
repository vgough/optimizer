#include "ceres/ceres.h"
#include "glog/logging.h"

#include "_cgo_export.h"

#include <list>

using ceres::DynamicNumericDiffCostFunction;
using ceres::NumericDiffOptions;
using ceres::CENTRAL;
using ceres::CostFunction;
using ceres::Problem;
using ceres::Solver;
using ceres::Solve;

struct bound {
  int param_num;
  double value;
};

struct config_internal {
  std::list<bound> lower;
  std::list<bound> upper;
};

ceres_config *create_config() {
  auto c = new (ceres_config);
  c->internal = static_cast<void *>(new (config_internal));
  return c;
}

void delete_config(ceres_config *config) {
  delete static_cast<config_internal*>(config->internal);
  delete config;
}

void add_lower_bound(ceres_config *config, int param_num, double lower) {
  bound b;
  b.param_num = param_num;
  b.value = lower;
  static_cast<config_internal *>(config->internal)->lower.push_back(b);
}

void add_upper_bound(ceres_config *config, int param_num, double upper) {
  bound b;
  b.param_num = param_num;
  b.value = upper;
  static_cast<config_internal *>(config->internal)->upper.push_back(b);
}

struct CostFunctor {
  CostFunctor(int id) : _id(id) {}
  bool operator()(double const *const *parameters, double *residuals) const {
    // ComputeResidual is a Go exported function.  Go doesn't have "const", so
    // cast that away to make the compiler happy.
    return ComputeResidual(_id, const_cast<double *>(parameters[0]), residuals);
  }

  int _id;
};

void init_ceres() {
  FLAGS_logtostderr = 1;
  google::InitGoogleLogging("optimize");
}

int optimize(int costFN, ceres_config *config, double *params) {
  Problem problem;

  NumericDiffOptions numeric_diff_options;
  if (config->relative_step_size > 1e-6) {
    numeric_diff_options.relative_step_size = config->relative_step_size;
  }

  auto cost_function = new DynamicNumericDiffCostFunction<CostFunctor>(
      new CostFunctor(costFN), ceres::TAKE_OWNERSHIP, numeric_diff_options);

  cost_function->AddParameterBlock(config->num_params);
  cost_function->SetNumResiduals(config->num_residuals);

  // Problem takes ownership of the cost function allocation.
  problem.AddResidualBlock(cost_function, NULL, params);

  auto internal = static_cast<config_internal*>(config->internal);
  for (auto it = internal->lower.begin(); it != internal->lower.end(); ++it) {
    problem.SetParameterLowerBound(params, it->param_num, it->value);
  }
  for (auto it = internal->upper.begin(); it != internal->upper.end(); ++it) {
    problem.SetParameterUpperBound(params, it->param_num, it->value);
  }

  // Run the solver.
  Solver::Options options;
  options.minimizer_progress_to_stdout = config->progress_to_stdout;
  Solver::Summary summary;
  Solve(options, &problem, &summary);

  if (config->print_summary) {
    std::cout << summary.FullReport() << "\n";
  }

  return 0;
}
