'use client';

import { useEffect } from 'react';
import { useRouter } from 'next/navigation';
import Image from 'next/image';

export default function Error({
  error,
  reset,
}: {
  error: Error & { digest?: string };
  reset: () => void;
}) {
  const router = useRouter();

  useEffect(() => {
    console.error(error);
  }, [error]);

  return (
    <div className="min-h-screen flex items-center justify-center bg-gradient-to-br from-slate-900 via-slate-800 to-slate-900 p-4">
      <div className="text-center">
        <Image src="/icon.svg" alt="Algoholic" width={64} height={64} className="w-16 h-16 mx-auto mb-6 opacity-50" />
        <h1 className="text-4xl font-bold text-white mb-4">Something went wrong</h1>
        <p className="text-gray-400 mb-8 max-w-md">
          An unexpected error occurred. Please try again or contact support if the problem persists.
        </p>
        <div className="flex gap-4 justify-center">
          <button
            onClick={() => reset()}
            className="px-6 py-3 bg-primary-500 hover:bg-primary-600 text-white rounded-lg transition-colors"
          >
            Try again
          </button>
          <button
            onClick={() => router.push('/dashboard')}
            className="px-6 py-3 bg-slate-700 hover:bg-slate-600 text-white rounded-lg transition-colors"
          >
            Go to Dashboard
          </button>
        </div>
      </div>
    </div>
  );
}
