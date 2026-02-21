'use client';

import Link from 'next/link';
import Image from 'next/image';
import { usePathname, useRouter } from 'next/navigation';
import { useAuthStore } from '@/stores/authStore';
import { LayoutDashboard, Code2, Dumbbell, Target, LogOut, Flame, User, List } from 'lucide-react';

const navItems = [
  { path: '/dashboard', label: 'Dashboard', icon: LayoutDashboard },
  { path: '/dashboard/problems', label: 'Problems', icon: Code2 },
  { path: '/dashboard/practice', label: 'Practice', icon: Target },
  { path: '/dashboard/training-plans', label: 'Training Plans', icon: Dumbbell },
  { path: '/dashboard/lists', label: 'Lists', icon: List },
  { path: '/dashboard/profile', label: 'Profile', icon: User },
];

export default function DashboardLayout({ children }: { children: React.ReactNode }) {
  const { user, logout } = useAuthStore();
  const router = useRouter();
  const pathname = usePathname();

  const handleLogout = () => {
    logout();
    router.push('/login');
  };

  const isActive = (path: string) => {
    if (path === '/dashboard') {
      return pathname === '/dashboard';
    }
    return pathname?.startsWith(path);
  };

  return (
    <div className="min-h-screen bg-gradient-to-br from-gray-50 to-gray-100">
      <nav className="bg-white shadow-sm border-b border-gray-200">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex justify-between h-16">
            <div className="flex items-center space-x-8">
              <Link
                href="/dashboard"
                className="flex items-center gap-2 text-2xl font-bold text-primary-600 hover:text-primary-700 transition-colors"
              >
                <Image src="/icon.svg" alt="Algoholic" width={32} height={32} className="w-8 h-8" />
                <span>Algoholic</span>
              </Link>

              <div className="hidden md:flex items-center space-x-1">
                {navItems.map((item) => {
                  const Icon = item.icon;
                  const active = isActive(item.path);
                  return (
                    <Link
                      key={item.path}
                      href={item.path}
                      className={`flex items-center gap-2 px-4 py-2 rounded-lg text-sm font-medium transition-colors ${
                        active
                          ? 'bg-primary-50 text-primary-700'
                          : 'text-gray-600 hover:text-gray-900 hover:bg-gray-50'
                      }`}
                    >
                      <Icon className="w-4 h-4" />
                      {item.label}
                    </Link>
                  );
                })}
              </div>
            </div>

            <div className="flex items-center space-x-4">
              {user && (
                <div className="flex items-center gap-3">
                  <div className="hidden sm:flex items-center gap-2 px-3 py-1.5 bg-orange-50 border border-orange-200 rounded-full">
                    <Flame className="w-4 h-4 text-orange-600" />
                    <span className="text-sm font-semibold text-orange-700">
                      {user.current_streak_days || 0}
                    </span>
                    <span className="text-xs text-orange-600">day streak</span>
                  </div>

                  <div className="hidden sm:block text-sm">
                    <span className="text-gray-600">Welcome,</span>{' '}
                    <span className="font-semibold text-gray-900">{user.username}</span>
                  </div>

                  <button
                    onClick={handleLogout}
                    className="btn-secondary flex items-center gap-2"
                  >
                    <LogOut className="w-4 h-4" />
                    <span className="hidden sm:inline">Logout</span>
                  </button>
                </div>
              )}
            </div>
          </div>
        </div>

        <div className="md:hidden border-t border-gray-200">
          <div className="flex justify-around py-2">
            {navItems.map((item) => {
              const Icon = item.icon;
              const active = isActive(item.path);
              return (
                <Link
                  key={item.path}
                  href={item.path}
                  className={`flex flex-col items-center gap-1 px-3 py-2 rounded-lg text-xs font-medium transition-colors ${
                    active ? 'text-primary-600' : 'text-gray-600'
                  }`}
                >
                  <Icon className="w-5 h-5" />
                  {item.label}
                </Link>
              );
            })}
          </div>
        </div>
      </nav>

      <main className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">{children}</main>

      <footer className="bg-white border-t border-gray-200 mt-auto">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-6">
          <div className="flex flex-col sm:flex-row justify-between items-center gap-4">
            <div className="flex items-center gap-2">
              <Image src="/icon.svg" alt="Algoholic" width={20} height={20} />
              <p className="text-sm text-gray-600">© 2024 Algoholic. Master DSA for FAANG Interviews.</p>
            </div>
            <div className="flex items-center gap-6 text-sm text-gray-600">
              <a href="https://github.com" target="_blank" rel="noopener noreferrer" className="hover:text-primary-600 transition-colors">
                GitHub
              </a>
              <a href="#" className="hover:text-primary-600 transition-colors">
                Help
              </a>
              <a href="#" className="hover:text-primary-600 transition-colors">
                Privacy
              </a>
            </div>
          </div>
        </div>
      </footer>
    </div>
  );
}
