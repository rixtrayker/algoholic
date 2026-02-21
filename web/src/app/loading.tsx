export default function Loading() {
  return (
    <div className="min-h-screen flex items-center justify-center bg-gradient-to-br from-slate-900 via-slate-800 to-slate-900">
      <div className="flex flex-col items-center gap-4">
        <div className="relative">
          <div className="w-16 h-16 border-4 border-primary-500/30 rounded-full animate-spin border-t-primary-500" />
          <div className="absolute inset-0 flex items-center justify-center">
            <div className="w-8 h-8 bg-primary-500 rounded-full animate-pulse" />
          </div>
        </div>
        <p className="text-gray-400 text-sm animate-pulse">Loading...</p>
      </div>
    </div>
  );
}
