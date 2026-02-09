import { useState } from 'react';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { questionsAPI, Question } from '../lib/api';
import { Clock, CheckCircle, XCircle, Lightbulb, ArrowRight } from 'lucide-react';

export default function Practice() {
  const [selectedAnswer, setSelectedAnswer] = useState('');
  const [startTime, setStartTime] = useState(Date.now());
  const [result, setResult] = useState<any>(null);
  const queryClient = useQueryClient();

  const { data: question, refetch, isLoading } = useQuery({
    queryKey: ['random-question'],
    queryFn: () => questionsAPI.getRandomQuestion(),
  });

  const submitMutation = useMutation({
    mutationFn: (answer: string) =>
      questionsAPI.submitAnswer(
        question!.question_id,
        { answer },
        Math.floor((Date.now() - startTime) / 1000)
      ),
    onSuccess: (data) => {
      setResult(data);
      // Invalidate stats to refresh dashboard
      queryClient.invalidateQueries({ queryKey: ['user-stats'] });
    },
  });

  const handleSubmit = () => {
    if (!selectedAnswer) return;
    submitMutation.mutate(selectedAnswer);
  };

  const handleNext = () => {
    setResult(null);
    setSelectedAnswer('');
    setStartTime(Date.now());
    refetch();
  };

  if (isLoading) {
    return <div className="text-center py-12">Loading question...</div>;
  }

  if (!question) {
    return (
      <div className="card text-center py-12">
        <p className="text-gray-500">No questions available. Please try again later.</p>
      </div>
    );
  }

  const timeTaken = Math.floor((Date.now() - startTime) / 1000);

  return (
    <div className="max-w-4xl mx-auto space-y-6">
      <div className="flex items-center justify-between">
        <h1 className="text-3xl font-bold text-gray-900">Practice</h1>
        {!result && (
          <div className="flex items-center gap-2 text-gray-600">
            <Clock className="w-5 h-5" />
            <span className="text-lg font-medium">
              {Math.floor(timeTaken / 60)}:{(timeTaken % 60).toString().padStart(2, '0')}
            </span>
          </div>
        )}
      </div>

      {/* Question Card */}
      <div className="card">
        <div className="mb-4">
          <div className="flex items-center justify-between mb-3">
            <span className="text-sm text-gray-500">
              Question #{question.question_id}
            </span>
            <span className="text-sm bg-gray-100 text-gray-700 px-3 py-1 rounded-full">
              {question.question_type}
            </span>
          </div>
          <h2 className="text-xl font-semibold text-gray-900 mb-2">
            {question.question_text}
          </h2>
          <div className="w-full bg-gray-200 rounded-full h-1.5 mb-2">
            <div
              className="bg-primary-600 h-1.5 rounded-full"
              style={{ width: `${question.difficulty_score}%` }}
            />
          </div>
          <p className="text-xs text-gray-600">
            Difficulty: {question.difficulty_score?.toFixed(0)}/100
          </p>
        </div>

        {!result ? (
          /* Answer Selection */
          <div className="space-y-4">
            <div className="space-y-3">
              {question.answer_options && Object.entries(question.answer_options).map(([key, option]: [string, any]) => (
                <label
                  key={key}
                  className={`flex items-start p-4 border-2 rounded-lg cursor-pointer transition-all ${
                    selectedAnswer === key
                      ? 'border-primary-500 bg-primary-50'
                      : 'border-gray-200 hover:border-gray-300 hover:bg-gray-50'
                  }`}
                >
                  <input
                    type="radio"
                    name="answer"
                    value={key}
                    checked={selectedAnswer === key}
                    onChange={(e) => setSelectedAnswer(e.target.value)}
                    className="mt-1 mr-3 text-primary-600 focus:ring-primary-500"
                  />
                  <div className="flex-1">
                    <span className="font-medium text-gray-900">{key}.</span>{' '}
                    <span className="text-gray-700">
                      {typeof option === 'string' ? option : option.text || JSON.stringify(option)}
                    </span>
                  </div>
                </label>
              ))}
            </div>

            <div className="flex items-center gap-3 pt-4">
              <button
                onClick={handleSubmit}
                disabled={!selectedAnswer || submitMutation.isPending}
                className="btn-primary flex-1 flex items-center justify-center gap-2"
              >
                {submitMutation.isPending ? (
                  'Submitting...'
                ) : (
                  <>
                    Submit Answer
                    <ArrowRight className="w-4 h-4" />
                  </>
                )}
              </button>
              <button className="btn-secondary flex items-center gap-2">
                <Lightbulb className="w-4 h-4" />
                Hint
              </button>
            </div>
          </div>
        ) : (
          /* Result Display */
          <div className="space-y-4">
            {/* Correct/Incorrect Banner */}
            <div
              className={`p-4 rounded-lg border-2 ${
                result.is_correct
                  ? 'bg-green-50 border-green-200'
                  : 'bg-red-50 border-red-200'
              }`}
            >
              <div className="flex items-center gap-3">
                {result.is_correct ? (
                  <CheckCircle className="w-6 h-6 text-green-600" />
                ) : (
                  <XCircle className="w-6 h-6 text-red-600" />
                )}
                <div>
                  <p className={`font-semibold text-lg ${
                    result.is_correct ? 'text-green-800' : 'text-red-800'
                  }`}>
                    {result.is_correct ? 'Correct!' : 'Incorrect'}
                  </p>
                  <p className={`text-sm ${
                    result.is_correct ? 'text-green-700' : 'text-red-700'
                  }`}>
                    {result.is_correct
                      ? 'Great job! You got it right.'
                      : `The correct answer was: ${result.correct_answer}`}
                  </p>
                </div>
              </div>
            </div>

            {/* Explanation */}
            {result.explanation && (
              <div className="p-4 bg-blue-50 border-2 border-blue-200 rounded-lg">
                <div className="flex items-start gap-2">
                  <Lightbulb className="w-5 h-5 text-blue-600 mt-0.5 flex-shrink-0" />
                  <div>
                    <p className="font-semibold text-blue-900 mb-1">Explanation</p>
                    <p className="text-blue-800 text-sm leading-relaxed">
                      {result.explanation}
                    </p>
                  </div>
                </div>
              </div>
            )}

            {/* Stats */}
            <div className="grid grid-cols-2 gap-4 pt-2">
              <div className="p-3 bg-gray-50 rounded-lg">
                <p className="text-xs text-gray-600 mb-1">Time Taken</p>
                <p className="text-lg font-semibold text-gray-900">
                  {Math.floor(timeTaken / 60)}m {timeTaken % 60}s
                </p>
              </div>
              <div className="p-3 bg-gray-50 rounded-lg">
                <p className="text-xs text-gray-600 mb-1">Points Earned</p>
                <p className="text-lg font-semibold text-primary-600">
                  +{result.points_earned || 0}
                </p>
              </div>
            </div>

            {/* Next Question Button */}
            <button
              onClick={handleNext}
              className="btn-primary w-full flex items-center justify-center gap-2"
            >
              Next Question
              <ArrowRight className="w-4 h-4" />
            </button>
          </div>
        )}
      </div>

      {/* Quick Tips */}
      {!result && (
        <div className="card bg-gray-50 border-gray-200">
          <h3 className="text-sm font-semibold text-gray-700 mb-2">Tips</h3>
          <ul className="text-xs text-gray-600 space-y-1">
            <li>Read the question carefully before selecting an answer</li>
            <li>Consider edge cases and time/space complexity</li>
            <li>Use hints if you're stuck, but try to solve it first</li>
          </ul>
        </div>
      )}
    </div>
  );
}
