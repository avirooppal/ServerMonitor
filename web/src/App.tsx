import { useEffect, useState } from 'react';
import { getApiKey } from './utils/api';
import { Login } from './components/Login';
import { Register } from './components/Register';
import { Dashboard } from './components/Dashboard';

function App() {
  const [isAuthenticated, setIsAuthenticated] = useState(false);
  const [view, setView] = useState<'login' | 'register' | 'dashboard'>('login');
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const key = getApiKey();
    if (key) {
      setIsAuthenticated(true);
      setView('dashboard');
    }
    setLoading(false);
  }, []);

  const handleLogin = () => {
    setIsAuthenticated(true);
    setView('dashboard');
  };

  const handleLogout = () => {
    setIsAuthenticated(false);
    setView('login');
  };

  if (loading) return null;

  return (
    <div className="min-h-screen bg-gray-900 text-white font-sans">
      {isAuthenticated ? (
        <Dashboard onLogout={handleLogout} />
      ) : view === 'login' ? (
        <Login onLogin={handleLogin} onSwitchToRegister={() => setView('register')} />
      ) : (
        <Register onRegister={() => setView('login')} onSwitchToLogin={() => setView('login')} />
      )}
    </div>
  );
}

export default App;
