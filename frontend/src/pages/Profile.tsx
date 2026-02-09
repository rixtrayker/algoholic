import { useQuery } from '@tanstack/react-query';
import { activityAPI, userAPI } from '../lib/api';
import type { UserStats } from '../lib/api';
import { Flame, TrendingUp, Clock, Target, Calendar, Award } from 'lucide-react';
import { format, subDays, eachDayOfInterval, startOfYear } from 'date-fns';

interface ActivityData {
  date: string;
  problems_count: number;
  questions_count: number;
  study_time_seconds: number;
  streak: number;
}

interface ActivityStats {
  total_days: number;
  total_problems: number;
  total_questions: number;
  total_study_time_seconds: number;
  current_streak: number;
  longest_streak: number;
  average_per_day: number;
  most_productive_day: string;
  most_productive_date: string;
}

interface PracticeHistoryItem {
  date: string;
  problems_count: number;
  questions_count: number;
  study_time_seconds: number;
  total_attempts: number;
  correct_attempts: number;
  accuracy_rate: number;
}

export default function Profile() {
  const { data: stats, isLoading: statsLoading } = useQuery<UserStats>({
    queryKey: ['user-stats'],
    queryFn: userAPI.getStats,
  });

  const { data: activityData, isLoading: activityLoading } = useQuery<ActivityData[]>({
    queryKey: ['activity-chart'],
    queryFn: () => activityAPI.getActivityChart(365),
  });

  const { data: activityStats } = useQuery<ActivityStats>({
    queryKey: ['activity-stats'],
    queryFn: activityAPI.getActivityStats,
  });

  const { data: history } = useQuery<PracticeHistoryItem[]>({
    queryKey: ['practice-history'],
    queryFn: () => activityAPI.getPracticeHistory(30),
  });

  const formatTime = (seconds: number) => {
    const hours = Math.floor(seconds / 3600);
    const minutes = Math.floor((seconds % 3600) / 60);
    if (hours > 0) {
      return `${hours}h ${minutes}m`;
    }
    return `${minutes}m`;
  };

  // Generate commitment chart data
  const generateChartData = () => {
    const today = new Date();
    const startDate = subDays(today, 364);
    const allDays = eachDayOfInterval({ start: startDate, end: today });

    const activityMap = new Map<string, number>();
    activityData?.forEach((activity) => {
      const total = activity.problems_count + activity.questions_count;
      activityMap.set(activity.date, total);
    });

    return allDays.map((date) => {
      const dateStr = format(date, 'yyyy-MM-dd');
      return {
        date: dateStr,
        count: activityMap.get(dateStr) || 0,
      };
    });
  };

  const chartData = activityData ? generateChartData() : [];

  // Group chart data by week
  const getWeeks = () => {
    const weeks: Array<Array<{ date: string; count: number }>> = [];
    let currentWeek: Array<{ date: string; count: number }> = [];

    chartData.forEach((day, index) => {
      const dayOfWeek = new Date(day.date).getDay();

      if (index === 0 && dayOfWeek !== 0) {
        // Fill in empty days at the start
        for (let i = 0; i < dayOfWeek; i++) {
          currentWeek.push({ date: '', count: 0 });
        }
      }

      currentWeek.push(day);

      if (currentWeek.length === 7) {
        weeks.push(currentWeek);
        currentWeek = [];
      }
    });

    if (currentWeek.length > 0) {
      // Fill in remaining days
      while (currentWeek.length < 7) {
        currentWeek.push({ date: '', count: 0 });
      }
      weeks.push(currentWeek);
    }

    return weeks;
  };

  const getColor = (count: number) => {
    if (count === 0) return 'bg-gray-100';
    if (count < 3) return 'bg-green-200';
    if (count < 6) return 'bg-green-400';
    if (count < 10) return 'bg-green-600';
    return 'bg-green-800';
  };

  if (statsLoading || activityLoading) {
    return (
      <div className="flex items-center justify-center min-h-screen">
        <div className="text-lg">Loading profile...</div>
      </div>
    );
  }

  const weeks = getWeeks();

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-3xl font-bold text-gray-900">Your Profile</h1>
        <p className="text-gray-600 mt-2">Track your progress and commitment</p>
      </div>

      {/* Stats Grid */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
        <div className="card">
          <div className="flex items-center gap-3">
            <div className="p-3 bg-orange-100 rounded-lg">
              <Flame className="w-6 h-6 text-orange-600" />
            </div>
            <div>
              <p className="text-sm text-gray-600">Current Streak</p>
              <p className="text-2xl font-bold text-gray-900">
                {activityStats?.current_streak || 0} days
              </p>
            </div>
          </div>
        </div>

        <div className="card">
          <div className="flex items-center gap-3">
            <div className="p-3 bg-blue-100 rounded-lg">
              <Target className="w-6 h-6 text-blue-600" />
            </div>
            <div>
              <p className="text-sm text-gray-600">Total Problems</p>
              <p className="text-2xl font-bold text-gray-900">
                {stats?.problems_solved || 0}
              </p>
            </div>
          </div>
        </div>

        <div className="card">
          <div className="flex items-center gap-3">
            <div className="p-3 bg-purple-100 rounded-lg">
              <Clock className="w-6 h-6 text-purple-600" />
            </div>
            <div>
              <p className="text-sm text-gray-600">Study Time</p>
              <p className="text-2xl font-bold text-gray-900">
                {formatTime(activityStats?.total_study_time_seconds || 0)}
              </p>
            </div>
          </div>
        </div>

        <div className="card">
          <div className="flex items-center gap-3">
            <div className="p-3 bg-green-100 rounded-lg">
              <Award className="w-6 h-6 text-green-600" />
            </div>
            <div>
              <p className="text-sm text-gray-600">Accuracy</p>
              <p className="text-2xl font-bold text-gray-900">
                {stats?.accuracy_rate?.toFixed(1) || 0}%
              </p>
            </div>
          </div>
        </div>
      </div>

      {/* Commitment Chart */}
      <div className="card">
        <div className="flex items-center gap-2 mb-4">
          <Calendar className="w-5 h-5 text-gray-700" />
          <h2 className="text-xl font-bold text-gray-900">Commitment Chart</h2>
        </div>

        <div className="overflow-x-auto">
          <div className="inline-block min-w-full">
            <div className="flex gap-1">
              <div className="flex flex-col gap-1 text-xs text-gray-600 mr-2">
                <div style={{ height: '12px' }}></div>
                <div style={{ height: '12px' }}>Mon</div>
                <div style={{ height: '12px' }}></div>
                <div style={{ height: '12px' }}>Wed</div>
                <div style={{ height: '12px' }}></div>
                <div style={{ height: '12px' }}>Fri</div>
                <div style={{ height: '12px' }}></div>
              </div>

              <div className="flex gap-1">
                {weeks.map((week, weekIndex) => (
                  <div key={weekIndex} className="flex flex-col gap-1">
                    {week.map((day, dayIndex) => (
                      <div
                        key={`${weekIndex}-${dayIndex}`}
                        className={`w-3 h-3 rounded-sm ${day.date ? getColor(day.count) : 'bg-transparent'}`}
                        title={day.date ? `${day.date}: ${day.count} activities` : ''}
                      />
                    ))}
                  </div>
                ))}
              </div>
            </div>

            <div className="flex items-center gap-2 mt-4 text-xs text-gray-600">
              <span>Less</span>
              <div className="flex gap-1">
                <div className="w-3 h-3 bg-gray-100 rounded-sm" />
                <div className="w-3 h-3 bg-green-200 rounded-sm" />
                <div className="w-3 h-3 bg-green-400 rounded-sm" />
                <div className="w-3 h-3 bg-green-600 rounded-sm" />
                <div className="w-3 h-3 bg-green-800 rounded-sm" />
              </div>
              <span>More</span>
            </div>
          </div>
        </div>

        {/* Additional stats */}
        {activityStats && (
          <div className="grid grid-cols-1 md:grid-cols-3 gap-4 mt-6 pt-6 border-t">
            <div>
              <p className="text-sm text-gray-600">Longest Streak</p>
              <p className="text-xl font-bold text-gray-900">{activityStats.longest_streak} days</p>
            </div>
            <div>
              <p className="text-sm text-gray-600">Active Days</p>
              <p className="text-xl font-bold text-gray-900">{activityStats.total_days} days</p>
            </div>
            <div>
              <p className="text-sm text-gray-600">Avg per Day</p>
              <p className="text-xl font-bold text-gray-900">
                {activityStats.average_per_day.toFixed(1)} activities
              </p>
            </div>
          </div>
        )}
      </div>

      {/* Practice History */}
      <div className="card">
        <div className="flex items-center gap-2 mb-4">
          <TrendingUp className="w-5 h-5 text-gray-700" />
          <h2 className="text-xl font-bold text-gray-900">Recent Practice History</h2>
        </div>

        <div className="overflow-x-auto">
          <table className="min-w-full divide-y divide-gray-200">
            <thead className="bg-gray-50">
              <tr>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  Date
                </th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  Problems
                </th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  Questions
                </th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  Attempts
                </th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  Accuracy
                </th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  Study Time
                </th>
              </tr>
            </thead>
            <tbody className="bg-white divide-y divide-gray-200">
              {history && history.length > 0 ? (
                history.map((item) => (
                  <tr key={item.date} className="hover:bg-gray-50">
                    <td className="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">
                      {format(new Date(item.date), 'MMM dd, yyyy')}
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-700">
                      {item.problems_count}
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-700">
                      {item.questions_count}
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-700">
                      {item.correct_attempts}/{item.total_attempts}
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm">
                      <span
                        className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${
                          item.accuracy_rate >= 70
                            ? 'bg-green-100 text-green-800'
                            : item.accuracy_rate >= 50
                            ? 'bg-yellow-100 text-yellow-800'
                            : 'bg-red-100 text-red-800'
                        }`}
                      >
                        {item.accuracy_rate.toFixed(1)}%
                      </span>
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-700">
                      {formatTime(item.study_time_seconds)}
                    </td>
                  </tr>
                ))
              ) : (
                <tr>
                  <td colSpan={6} className="px-6 py-8 text-center text-gray-500">
                    No practice history yet. Start practicing to see your progress!
                  </td>
                </tr>
              )}
            </tbody>
          </table>
        </div>
      </div>
    </div>
  );
}
