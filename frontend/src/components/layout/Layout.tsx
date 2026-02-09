import { Outlet, Link, useNavigate, useLocation } from 'react-router-dom';
import { useAuthStore } from '../../stores/authStore';
import { LayoutDashboard, Code2, Dumbbell, Target, LogOut, Flame } from 'lucide-react';

export default function Layout() {
  const { user, logout } = useAuthStore();
  const navigate = useNavigate();
  const location = useLocation();

  const handleLogout = () => {
    logout();
    navigate('/login');
  };

  const navItems = [
    { path: '/', label: 'Dashboard', icon: LayoutDashboard },
    { path: '/problems', label: 'Problems', icon: Code2 },
    { path: '/practice', label: 'Practice', icon: Target },
    { path: '/training-plans', label: 'Training Plans', icon: Dumbbell },
  ];

  const isActive = (path: string) => {
    if (path === '/') {
      return location.pathname === '/';
    }
    return location.pathname.startsWith(path);
  };

  return (
    <div className="min-h-screen bg-gradient-to-br from-gray-50 to-gray-100">
      {/* Navigation Bar */}
      <nav className="bg-white shadow-sm border-b border-gray-200">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex justify-between h-16">
            {/* Logo and Nav Links */}
            <div className="flex items-center space-x-8">
              <Link
                to="/"
                className="flex items-center gap-2 text-2xl font-bold text-primary-600 hover:text-primary-700 transition-colors"
              >
                <Flame className="w-7 h-7" />
                Algoholic
              </Link>

              <div className="hidden md:flex items-center space-x-1">
                {navItems.map((item) => {
                  const Icon = item.icon;
                  const active = isActive(item.path);
                  return (
                    <Link
                      key={item.path}
                      to={item.path}
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

            {/* User Info and Logout */}
            <div className="flex items-center space-x-4">
              {user && (
                <div className="flex items-center gap-3">
                  {/* Streak Display */}
                  <div className="hidden sm:flex items-center gap-2 px-3 py-1.5 bg-orange-50 border border-orange-200 rounded-full">
                    <Flame className="w-4 h-4 text-orange-600" />
                    <span className="text-sm font-semibold text-orange-700">
                      {user.current_streak_days || 0}
                    </span>
                    <span className="text-xs text-orange-600">day streak</span>
                  </div>

                  {/* Username */}
                  <div className="hidden sm:block text-sm">
                    <span className="text-gray-600">Welcome,</span>{' '}
                    <span className="font-semibold text-gray-900">{user.username}</span>
                  </div>

                  {/* Logout Button */}
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

        {/* Mobile Navigation */}
        <div className="md:hidden border-t border-gray-200">
          <div className="flex justify-around py-2">
            {navItems.map((item) => {
              const Icon = item.icon;
              const active = isActive(item.path);
              return (
                <Link
                  key={item.path}
                  to={item.path}
                  className={`flex flex-col items-center gap-1 px-3 py-2 rounded-lg text-xs font-medium transition-colors ${
                    active
                      ? 'text-primary-600'
                      : 'text-gray-600'
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

      {/* Main Content */}
      <main className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <Outlet />
      </main>

      {/* Footer */}
      <footer className="bg-white border-t border-gray-200 mt-auto">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-6">
          <div className="flex flex-col sm:flex-row justify-between items-center gap-4">
            <p className="text-sm text-gray-600">
              © 2024 Algoholic. Master DSA for FAANG Interviews.
            </p>
            <div className="flex items-center gap-6 text-sm text-gray-600">
              <a href="#" className="hover:text-primary-600 transition-colors">About</a>
              <a href="#" className="hover:text-primary-600 transition-colors">Help</a>
              <a href="#" className="hover:text-primary-600 transition-colors">Privacy</a>
            </div>
          </div>
        </div>
      </footer>
    </div>
  );
}
