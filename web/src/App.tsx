import { useEffect, useState } from 'react';
import { getApiKey, getAuthMode } from './utils/api';
import { Setup } from './components/Setup';
import { Dashboard } from './components/Dashboard';
import { Login } from './components/Login';
import { Register } from './components/Register';
import { AuthCallback } from './components/AuthCallback';
import { Loader2 } from 'lucide-react';

function App() {
  const [isAuthenticated, setIsAuthenticated] = useState(false);
  const [loading, setLoading] = useState(true);
  const [saasMode, setSaasMode] = useState(false);
  const [authView, setAuthView] = useState<'login' | 'register'>('login');
  const [isCallback, setIsCallback] = useState(false);

  useEffect(() => {
    checkAuth();
  }, []);

  const checkAuth = async () => {
    try {
      // Check for OAuth callback
      if (window.location.search.includes('code=')) {
        setIsCallback(true);
        // Don't check key yet, let callback handle it
      } else {
        const key = getApiKey();
        if (key) {
          setIsAuthenticated(true);
        }
      }

      // Fetch Auth Mode
      const config = await getAuthMode();
      setSaasMode(config.saas_mode);
    } catch (e) {
      console.error("Failed to fetch auth config", e);
    } finally {
      setLoading(false);
    }
  };

  if (loading) {
    return (
      <div className="min-h-screen bg-gray-900 flex items-center justify-center">
        <Loader2 className="text-blue-500 animate-spin" size={48} />
      </div>
    );
  }

  if (isCallback) {
    return <AuthCallback onLogin={() => { setIsCallback(false); setIsAuthenticated(true); }} />;
  }

  if (isAuthenticated) {
    return <Dashboard onLogout={() => setIsAuthenticated(false)} />;
  }

  // SaaS Mode Flow
  if (saasMode) {
    return authView === 'login' ? (
      <Login onLogin={() => setIsAuthenticated(true)} onSwitchToRegister={() => setAuthView('register')} />
    ) : (
      <Register onRegisterSuccess={() => setAuthView('login')} onSwitchToLogin={() => setAuthView('login')} />
    );
  }

  // Self-Hosted Mode Flow
  return <Setup onComplete={() => setIsAuthenticated(true)} />;
}

export default App;
