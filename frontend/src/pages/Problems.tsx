import { useQuery } from '@tanstack/react-query';
import { problemsAPI, Problem } from '../lib/api';
import { useState } from 'react';
import { Search, Filter } from 'lucide-react';

export default function Problems() {
  const [searchQuery, setSearchQuery] = useState('');
  const [minDifficulty, setMinDifficulty] = useState(0);
  const [maxDifficulty, setMaxDifficulty] = useState(100);
  const [selectedPattern, setSelectedPattern] = useState('');

  const { data, isLoading } = useQuery({
    queryKey: ['problems', minDifficulty, maxDifficulty, selectedPattern],
    queryFn: () => problemsAPI.getProblems({
      min_difficulty: minDifficulty,
      max_difficulty: maxDifficulty,
      pattern: selectedPattern || undefined,
      limit: 50,
    }),
  });

  const { data: searchResults } = useQuery({
    queryKey: ['problem-search', searchQuery],
    queryFn: () => problemsAPI.searchProblems(searchQuery),
    enabled: searchQuery.length > 2,
  });

  const problems = searchQuery.length > 2
    ? searchResults?.problems || []
    : data?.problems || [];

  const getDifficultyColor = (score: number) => {
    if (score < 30) return 'bg-green-100 text-green-700 border-green-200';
    if (score < 60) return 'bg-yellow-100 text-yellow-700 border-yellow-200';
    if (score < 80) return 'bg-orange-100 text-orange-700 border-orange-200';
    return 'bg-red-100 text-red-700 border-red-200';
  };

  if (isLoading) {
    return <div className="text-center py-12">Loading problems...</div>;
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <h1 className="text-3xl font-bold text-gray-900">Problems</h1>
        <div className="text-sm text-gray-600">
          {problems.length} problems found
        </div>
      </div>

      {/* Search and Filters */}
      <div className="card">
        <div className="space-y-4">
          {/* Search Bar */}
          <div className="relative">
            <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400 w-5 h-5" />
            <input
              type="text"
              placeholder="Search problems by title, pattern, or topic..."
              className="input pl-10"
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
            />
          </div>

          {/* Filters */}
          <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">
                <Filter className="inline w-4 h-4 mr-1" />
                Min Difficulty
              </label>
              <input
                type="range"
                min="0"
                max="100"
                value={minDifficulty}
                onChange={(e) => setMinDifficulty(Number(e.target.value))}
                className="w-full"
              />
              <div className="text-xs text-gray-600 mt-1">{minDifficulty}</div>
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">
                Max Difficulty
              </label>
              <input
                type="range"
                min="0"
                max="100"
                value={maxDifficulty}
                onChange={(e) => setMaxDifficulty(Number(e.target.value))}
                className="w-full"
              />
              <div className="text-xs text-gray-600 mt-1">{maxDifficulty}</div>
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">
                Pattern
              </label>
              <select
                className="input"
                value={selectedPattern}
                onChange={(e) => setSelectedPattern(e.target.value)}
              >
                <option value="">All Patterns</option>
                <option value="Two Pointers">Two Pointers</option>
                <option value="Sliding Window">Sliding Window</option>
                <option value="Binary Search">Binary Search</option>
                <option value="DFS">DFS</option>
                <option value="BFS">BFS</option>
                <option value="Dynamic Programming">Dynamic Programming</option>
                <option value="Backtracking">Backtracking</option>
                <option value="Greedy">Greedy</option>
              </select>
            </div>
          </div>
        </div>
      </div>

      {/* Problems List */}
      <div className="space-y-4">
        {problems.length === 0 ? (
          <div className="card text-center py-12">
            <p className="text-gray-500">No problems found. Try adjusting your filters.</p>
          </div>
        ) : (
          problems.map((problem: Problem) => (
            <div
              key={problem.problem_id}
              className="card hover:shadow-lg transition-shadow cursor-pointer"
            >
              <div className="flex items-start justify-between mb-3">
                <h3 className="text-xl font-semibold text-gray-900">
                  {problem.title}
                </h3>
                <span className={`text-sm px-3 py-1 rounded-full border ${getDifficultyColor(problem.difficulty_score)}`}>
                  {problem.difficulty_score?.toFixed(0)}
                </span>
              </div>

              <p className="text-gray-600 mb-4 line-clamp-2">
                {problem.description}
              </p>

              <div className="flex flex-wrap items-center gap-2">
                {problem.primary_pattern && (
                  <span className="text-xs bg-primary-100 text-primary-700 px-3 py-1 rounded-full">
                    {problem.primary_pattern}
                  </span>
                )}
                <span className="text-xs text-gray-500">
                  #{problem.problem_id}
                </span>
              </div>
            </div>
          ))
        )}
      </div>
    </div>
  );
}
