'use client';

import { useState } from 'react';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { trainingPlansAPI, topicsAPI } from '@/lib/api';
import type { TrainingPlan, Topic } from '@/lib/api';
import { Calendar, Target, TrendingUp, Play, Pause, Trash2, Plus, X } from 'lucide-react';
import toast from 'react-hot-toast';

export default function TrainingPlansPage() {
  const queryClient = useQueryClient();
  const [showCreateForm, setShowCreateForm] = useState(false);
  const [formData, setFormData] = useState({
    name: '',
    description: '',
    plan_type: 'custom',
    target_topics: [] as number[],
    target_patterns: [] as string[],
    duration_days: 30,
    questions_per_day: 5,
    difficulty_min: 30,
    difficulty_max: 70,
    adaptive_difficulty: true,
  });

  const { data: plansData, isLoading } = useQuery({
    queryKey: ['training-plans'],
    queryFn: trainingPlansAPI.getPlans,
  });

  const { data: topicsData } = useQuery({
    queryKey: ['topics'],
    queryFn: topicsAPI.getTopics,
  });

  const plans = plansData?.plans || [];
  const topics = topicsData?.topics || [];

  const createMutation = useMutation({
    mutationFn: trainingPlansAPI.createPlan,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['training-plans'] });
      toast.success('Training plan created!');
      setShowCreateForm(false);
      setFormData({
        name: '',
        description: '',
        plan_type: 'custom',
        target_topics: [],
        target_patterns: [],
        duration_days: 30,
        questions_per_day: 5,
        difficulty_min: 30,
        difficulty_max: 70,
        adaptive_difficulty: true,
      });
    },
    onError: () => toast.error('Failed to create plan'),
  });

  const pauseMutation = useMutation({
    mutationFn: trainingPlansAPI.pausePlan,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['training-plans'] });
      toast.success('Plan paused');
    },
  });

  const resumeMutation = useMutation({
    mutationFn: trainingPlansAPI.resumePlan,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['training-plans'] });
      toast.success('Plan resumed');
    },
  });

  const deleteMutation = useMutation({
    mutationFn: trainingPlansAPI.deletePlan,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['training-plans'] });
      toast.success('Plan deleted');
    },
  });

  const handleDelete = (planId: number) => {
    if (window.confirm('Are you sure you want to delete this training plan?')) {
      deleteMutation.mutate(planId);
    }
  };

  const handleCreateSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (!formData.name.trim()) {
      toast.error('Plan name is required');
      return;
    }
    createMutation.mutate(formData);
  };

  const toggleTopic = (topicId: number) => {
    setFormData((prev) => ({
      ...prev,
      target_topics: prev.target_topics.includes(topicId)
        ? prev.target_topics.filter((id) => id !== topicId)
        : [...prev.target_topics, topicId],
    }));
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
        <button onClick={() => setShowCreateForm(true)} className="btn-primary flex items-center gap-2">
          <Plus className="w-5 h-5" />
          Create New Plan
        </button>
      </div>

      {showCreateForm && (
        <div className="card">
          <div className="flex items-center justify-between mb-4">
            <h2 className="text-xl font-bold text-gray-900">Create New Training Plan</h2>
            <button onClick={() => setShowCreateForm(false)} className="text-gray-500 hover:text-gray-700">
              <X className="w-5 h-5" />
            </button>
          </div>

          <form onSubmit={handleCreateSubmit} className="space-y-4">
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">Plan Name *</label>
                <input
                  type="text"
                  value={formData.name}
                  onChange={(e) => setFormData({ ...formData, name: e.target.value })}
                  className="input"
                  placeholder="e.g., 30-Day DP Bootcamp"
                  required
                />
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">Plan Type</label>
                <select
                  value={formData.plan_type}
                  onChange={(e) => setFormData({ ...formData, plan_type: e.target.value })}
                  className="input"
                >
                  <option value="custom">Custom</option>
                  <option value="topic_focused">Topic Focused</option>
                  <option value="pattern_focused">Pattern Focused</option>
                  <option value="adaptive">Adaptive</option>
                </select>
              </div>
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">Description</label>
              <textarea
                value={formData.description}
                onChange={(e) => setFormData({ ...formData, description: e.target.value })}
                className="input"
                rows={2}
                placeholder="Optional description"
              />
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">Target Topics</label>
              <div className="flex flex-wrap gap-2">
                {topics.map((topic: Topic) => (
                  <button
                    key={topic.topic_id}
                    type="button"
                    onClick={() => toggleTopic(topic.topic_id)}
                    className={`px-3 py-1 rounded-full text-sm border transition-colors ${
                      formData.target_topics.includes(topic.topic_id)
                        ? 'bg-primary-100 border-primary-500 text-primary-700'
                        : 'bg-gray-50 border-gray-200 text-gray-700 hover:bg-gray-100'
                    }`}
                  >
                    {topic.name}
                  </button>
                ))}
              </div>
            </div>

            <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">Duration (days)</label>
                <input
                  type="number"
                  min="1"
                  max="365"
                  value={formData.duration_days}
                  onChange={(e) => setFormData({ ...formData, duration_days: Number(e.target.value) })}
                  className="input"
                />
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">Questions per Day</label>
                <input
                  type="number"
                  min="1"
                  max="20"
                  value={formData.questions_per_day}
                  onChange={(e) => setFormData({ ...formData, questions_per_day: Number(e.target.value) })}
                  className="input"
                />
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">Min Difficulty</label>
                <input
                  type="number"
                  min="0"
                  max="100"
                  value={formData.difficulty_min}
                  onChange={(e) => setFormData({ ...formData, difficulty_min: Number(e.target.value) })}
                  className="input"
                />
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">Max Difficulty</label>
                <input
                  type="number"
                  min="0"
                  max="100"
                  value={formData.difficulty_max}
                  onChange={(e) => setFormData({ ...formData, difficulty_max: Number(e.target.value) })}
                  className="input"
                />
              </div>
            </div>

            <div className="flex items-center gap-2">
              <input
                type="checkbox"
                id="adaptive"
                checked={formData.adaptive_difficulty}
                onChange={(e) => setFormData({ ...formData, adaptive_difficulty: e.target.checked })}
                className="w-4 h-4 text-primary-600"
              />
              <label htmlFor="adaptive" className="text-sm font-medium text-gray-700">
                Adaptive difficulty (adjusts based on performance)
              </label>
            </div>

            <div className="flex gap-3">
              <button type="submit" disabled={createMutation.isPending} className="btn-primary">
                {createMutation.isPending ? 'Creating...' : 'Create Plan'}
              </button>
              <button type="button" onClick={() => setShowCreateForm(false)} className="btn-secondary">
                Cancel
              </button>
            </div>
          </form>
        </div>
      )}

      {plans.length > 0 ? (
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
          {plans.map((plan: TrainingPlan) => (
            <div key={plan.plan_id} className="card hover:shadow-lg transition-shadow">
              <div className="flex items-start justify-between mb-4">
                <div className="flex-1">
                  <h3 className="text-xl font-semibold text-gray-900 mb-1">{plan.name}</h3>
                  {plan.description && (
                    <p className="text-sm text-gray-600 line-clamp-2">{plan.description}</p>
                  )}
                </div>
                <span className={`text-xs px-3 py-1 rounded-full border ${getStatusColor(plan.status)}`}>
                  {plan.status}
                </span>
              </div>

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
                    <p className="text-sm font-semibold text-gray-900">{plan.questions_per_day} questions</p>
                  </div>
                </div>
              </div>

              <div className="flex items-center gap-2">
                {plan.status.toLowerCase() === 'active' ? (
                  <>
                    <button className="btn-primary flex-1 flex items-center justify-center gap-2">
                      <Play className="w-4 h-4" />
                      Continue
                    </button>
                    <button
                      onClick={() => pauseMutation.mutate(plan.plan_id)}
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
                      onClick={() => resumeMutation.mutate(plan.plan_id)}
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
        <div className="card text-center py-12">
          <TrendingUp className="w-16 h-16 text-gray-400 mx-auto mb-4" />
          <h3 className="text-lg font-semibold text-gray-900 mb-2">No Training Plans Yet</h3>
          <p className="text-gray-600 mb-6">
            Create a personalized training plan to track your progress and stay motivated.
          </p>
          <button onClick={() => setShowCreateForm(true)} className="btn-primary">
            Create Your First Plan
          </button>
        </div>
      )}

      <div className="card bg-primary-50 border-primary-200">
        <h3 className="text-lg font-semibold text-primary-900 mb-2">About Training Plans</h3>
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
