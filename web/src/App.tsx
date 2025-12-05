import { useEffect, useState } from 'react';
import { getApiKey } from './utils/api';
import Auth from './components/Auth';
import { Dashboard } from './components/Dashboard';

function App() {
  const [isAuthenticated, setIsAuthenticated] = useState(false);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const key = getApiKey();
    if (key) {
      setIsAuthenticated(true);
    }
    setLoading(false);
  }, []);

  if (loading) return null;

  return (
    <div className="min-h-screen bg-gray-900 text-white font-sans">
      {isAuthenticated ? (
        <Dashboard onLogout={() => setIsAuthenticated(false)} />
      ) : (
        <Auth onAuth={() => setIsAuthenticated(true)} />
      )}
    </div>
  );
}

export default App;
