import { useQuery } from '@tanstack/react-query';
import { userAPI } from '../lib/api';
import { BarChart, TrendingUp, Flame, Target } from 'lucide-react';

export default function Dashboard() {
  const { data: stats, isLoading } = useQuery({
    queryKey: ['user-stats'],
    queryFn: userAPI.getStats,
  });

  if (isLoading) {
    return <div className="text-center py-12">Loading your stats...</div>;
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <h1 className="text-3xl font-bold text-gray-900">Dashboard</h1>
        <div className="flex items-center gap-2 text-orange-600">
          <Flame className="w-6 h-6" />
          <span className="text-2xl font-bold">{stats?.current_streak_days || 0}</span>
          <span className="text-sm">day streak</span>
        </div>
      </div>

      {/* Stats Grid */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
        <div className="card bg-gradient-to-br from-blue-50 to-blue-100 border-blue-200">
          <div className="flex items-center justify-between mb-2">
            <h3 className="text-sm font-medium text-gray-700">Total Attempts</h3>
            <BarChart className="w-5 h-5 text-blue-600" />
          </div>
          <p className="text-3xl font-bold text-blue-700">
            {stats?.total_attempts || 0}
          </p>
          <p className="text-xs text-gray-600 mt-1">
            {stats?.correct_attempts || 0} correct
          </p>
        </div>

        <div className="card bg-gradient-to-br from-green-50 to-green-100 border-green-200">
          <div className="flex items-center justify-between mb-2">
            <h3 className="text-sm font-medium text-gray-700">Accuracy Rate</h3>
            <Target className="w-5 h-5 text-green-600" />
          </div>
          <p className="text-3xl font-bold text-green-700">
            {stats?.accuracy_rate?.toFixed(1) || 0}%
          </p>
          <p className="text-xs text-gray-600 mt-1">
            Keep it above 70%!
          </p>
        </div>

        <div className="card bg-gradient-to-br from-purple-50 to-purple-100 border-purple-200">
          <div className="flex items-center justify-between mb-2">
            <h3 className="text-sm font-medium text-gray-700">Problems Solved</h3>
            <TrendingUp className="w-5 h-5 text-purple-600" />
          </div>
          <p className="text-3xl font-bold text-purple-700">
            {stats?.problems_solved || 0}
          </p>
          <p className="text-xs text-gray-600 mt-1">
            of {stats?.problems_attempted || 0} attempted
          </p>
        </div>

        <div className="card bg-gradient-to-br from-orange-50 to-orange-100 border-orange-200">
          <div className="flex items-center justify-between mb-2">
            <h3 className="text-sm font-medium text-gray-700">Study Time</h3>
            <Flame className="w-5 h-5 text-orange-600" />
          </div>
          <p className="text-3xl font-bold text-orange-700">
            {Math.floor((stats?.total_study_time_seconds || 0) / 3600)}h
          </p>
          <p className="text-xs text-gray-600 mt-1">
            {Math.floor(((stats?.total_study_time_seconds || 0) % 3600) / 60)}m total
          </p>
        </div>
      </div>

      {/* Topics Grid */}
      <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
        {/* Strong Topics */}
        <div className="card">
          <h3 className="text-lg font-semibold mb-4 flex items-center gap-2">
            <span className="text-green-600">✓</span>
            Strong Topics
          </h3>
          {stats?.strong_topics && stats.strong_topics.length > 0 ? (
            <ul className="space-y-2">
              {stats.strong_topics.map((topic) => (
                <li
                  key={topic}
                  className="flex items-center p-2 bg-green-50 rounded-lg"
                >
                  <span className="text-green-600 mr-3 font-bold">✓</span>
                  <span className="text-gray-800">{topic}</span>
                </li>
              ))}
            </ul>
          ) : (
            <p className="text-gray-500 text-sm">
              Complete more questions to see your strengths!
            </p>
          )}
        </div>

        {/* Weak Topics */}
        <div className="card">
          <h3 className="text-lg font-semibold mb-4 flex items-center gap-2">
            <span className="text-red-600">!</span>
            Areas to Improve
          </h3>
          {stats?.weak_topics && stats.weak_topics.length > 0 ? (
            <ul className="space-y-2">
              {stats.weak_topics.map((topic) => (
                <li
                  key={topic}
                  className="flex items-center p-2 bg-red-50 rounded-lg"
                >
                  <span className="text-red-600 mr-3 font-bold">!</span>
                  <span className="text-gray-800">{topic}</span>
                </li>
              ))}
            </ul>
          ) : (
            <p className="text-gray-500 text-sm">
              No weak areas detected yet. Keep practicing!
            </p>
          )}
        </div>
      </div>

      {/* Quick Actions */}
      <div className="card bg-primary-50 border-primary-200">
        <h3 className="text-lg font-semibold mb-3">Quick Actions</h3>
        <div className="flex flex-wrap gap-3">
          <a href="/practice" className="btn-primary">
            Start Practicing
          </a>
          <a href="/problems" className="btn-secondary">
            Browse Problems
          </a>
          <a href="/training-plans" className="btn-secondary">
            View Training Plans
          </a>
        </div>
      </div>
    </div>
  );
}
