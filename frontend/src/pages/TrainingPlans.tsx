import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { trainingPlansAPI, TrainingPlan } from '../lib/api';
import { Calendar, Target, TrendingUp, Play, Pause, Trash2 } from 'lucide-react';

export default function TrainingPlans() {
  const queryClient = useQueryClient();

  const { data: plans, isLoading } = useQuery({
    queryKey: ['training-plans'],
    queryFn: trainingPlansAPI.getPlans,
  });

  const pauseMutation = useMutation({
    mutationFn: (planId: number) => trainingPlansAPI.pausePlan(planId),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['training-plans'] });
    },
  });

  const resumeMutation = useMutation({
    mutationFn: (planId: number) => trainingPlansAPI.resumePlan(planId),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['training-plans'] });
    },
  });

  const deleteMutation = useMutation({
    mutationFn: (planId: number) => trainingPlansAPI.deletePlan(planId),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['training-plans'] });
    },
  });

  const handlePause = (planId: number) => {
    pauseMutation.mutate(planId);
  };

  const handleResume = (planId: number) => {
    resumeMutation.mutate(planId);
  };

  const handleDelete = (planId: number) => {
    if (window.confirm('Are you sure you want to delete this training plan?')) {
      deleteMutation.mutate(planId);
    }
  };

  const getStatusColor = (status: string) => {
    switch (status.toLowerCase()) {
      case 'active':
        return 'bg-green-100 text-green-700 border-green-200';
      case 'paused':
        return 'bg-yellow-100 text-yellow-700 border-yellow-200';
      case 'completed':
        return 'bg-blue-100 text-blue-700 border-blue-200';
      default:
        return 'bg-gray-100 text-gray-700 border-gray-200';
    }
  };

  if (isLoading) {
    return <div className="text-center py-12">Loading training plans...</div>;
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <h1 className="text-3xl font-bold text-gray-900">Training Plans</h1>
        <button className="btn-primary">
          Create New Plan
        </button>
      </div>

      {/* Plans Grid */}
      {plans && plans.length > 0 ? (
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
          {plans.map((plan: TrainingPlan) => (
            <div key={plan.plan_id} className="card hover:shadow-lg transition-shadow">
              {/* Header */}
              <div className="flex items-start justify-between mb-4">
                <div className="flex-1">
                  <h3 className="text-xl font-semibold text-gray-900 mb-1">
                    {plan.name}
                  </h3>
                  {plan.description && (
                    <p className="text-sm text-gray-600 line-clamp-2">
                      {plan.description}
                    </p>
                  )}
                </div>
                <span className={`text-xs px-3 py-1 rounded-full border ${getStatusColor(plan.status)}`}>
                  {plan.status}
                </span>
              </div>

              {/* Progress Bar */}
              <div className="mb-4">
                <div className="flex items-center justify-between mb-2">
                  <span className="text-sm font-medium text-gray-700">Progress</span>
                  <span className="text-sm font-semibold text-primary-600">
                    {plan.progress_percentage?.toFixed(0) || 0}%
                  </span>
                </div>
                <div className="w-full bg-gray-200 rounded-full h-3">
                  <div
                    className="bg-gradient-to-r from-primary-500 to-primary-600 h-3 rounded-full transition-all duration-300"
                    style={{ width: `${plan.progress_percentage || 0}%` }}
                  />
                </div>
              </div>

              {/* Stats */}
              <div className="grid grid-cols-2 gap-3 mb-4">
                <div className="flex items-center gap-2 p-3 bg-gray-50 rounded-lg">
                  <Calendar className="w-4 h-4 text-gray-600" />
                  <div>
                    <p className="text-xs text-gray-600">Started</p>
                    <p className="text-sm font-semibold text-gray-900">
                      {new Date(plan.start_date).toLocaleDateString()}
                    </p>
                  </div>
                </div>
                <div className="flex items-center gap-2 p-3 bg-gray-50 rounded-lg">
                  <Target className="w-4 h-4 text-gray-600" />
                  <div>
                    <p className="text-xs text-gray-600">Daily Goal</p>
                    <p className="text-sm font-semibold text-gray-900">
                      {plan.questions_per_day} questions
                    </p>
                  </div>
                </div>
              </div>

              {/* Actions */}
              <div className="flex items-center gap-2">
                {plan.status.toLowerCase() === 'active' ? (
                  <>
                    <button className="btn-primary flex-1 flex items-center justify-center gap-2">
                      <Play className="w-4 h-4" />
                      Continue
                    </button>
                    <button
                      onClick={() => handlePause(plan.plan_id)}
                      disabled={pauseMutation.isPending}
                      className="btn-secondary flex items-center gap-2"
                    >
                      <Pause className="w-4 h-4" />
                      Pause
                    </button>
                  </>
                ) : plan.status.toLowerCase() === 'paused' ? (
                  <>
                    <button
                      onClick={() => handleResume(plan.plan_id)}
                      disabled={resumeMutation.isPending}
                      className="btn-primary flex-1 flex items-center justify-center gap-2"
                    >
                      <Play className="w-4 h-4" />
                      Resume
                    </button>
                    <button
                      onClick={() => handleDelete(plan.plan_id)}
                      disabled={deleteMutation.isPending}
                      className="btn-secondary flex items-center gap-2 text-red-600 hover:bg-red-50"
                    >
                      <Trash2 className="w-4 h-4" />
                    </button>
                  </>
                ) : (
                  <button className="btn-secondary flex-1" disabled>
                    Completed
                  </button>
                )}
              </div>
            </div>
          ))}
        </div>
      ) : (
        /* Empty State */
        <div className="card text-center py-12">
          <TrendingUp className="w-16 h-16 text-gray-400 mx-auto mb-4" />
          <h3 className="text-lg font-semibold text-gray-900 mb-2">
            No Training Plans Yet
          </h3>
          <p className="text-gray-600 mb-6">
            Create a personalized training plan to track your progress and stay motivated.
          </p>
          <button className="btn-primary">
            Create Your First Plan
          </button>
        </div>
      )}

      {/* Info Card */}
      <div className="card bg-primary-50 border-primary-200">
        <h3 className="text-lg font-semibold text-primary-900 mb-2">
          About Training Plans
        </h3>
        <p className="text-sm text-primary-800 mb-3">
          Training plans help you stay consistent with daily practice goals tailored to your skill level.
          Track your progress, build streaks, and master DSA concepts systematically.
        </p>
        <ul className="text-sm text-primary-700 space-y-1">
          <li>• Adaptive difficulty based on your performance</li>
          <li>• Daily question goals to build consistency</li>
          <li>• Focus on specific topics or patterns</li>
          <li>• Spaced repetition for better retention</li>
        </ul>
      </div>
    </div>
  );
}
