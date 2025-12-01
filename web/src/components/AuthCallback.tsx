import React, { useEffect, useState } from 'react';
import { githubLogin } from '../utils/api';
import { Loader2, AlertCircle } from 'lucide-react';

interface AuthCallbackProps {
    onLogin: () => void;
}

export const AuthCallback: React.FC<AuthCallbackProps> = ({ onLogin }) => {
    const [error, setError] = useState('');

    useEffect(() => {
        const code = new URLSearchParams(window.location.search).get('code');
        if (code) {
            handleGitHubCallback(code);
        } else {
            setError('No authorization code found');
        }
    }, []);

    const handleGitHubCallback = async (code: string) => {
        try {
            await githubLogin(code);
            // Remove code from URL
            window.history.replaceState({}, document.title, window.location.pathname);
            onLogin();
        } catch (err: any) {
            setError(err.response?.data?.error || 'GitHub login failed');
        }
    };

    if (error) {
        return (
            <div className="min-h-screen flex items-center justify-center bg-gray-900 text-white">
                <div className="bg-red-500/10 border border-red-500/20 rounded-lg p-6 flex flex-col items-center max-w-md text-center">
                    <AlertCircle size={48} className="text-red-500 mb-4" />
                    <h2 className="text-xl font-bold mb-2">Authentication Failed</h2>
                    <p className="text-gray-400 mb-4">{error}</p>
                    <a href="/" className="bg-gray-800 hover:bg-gray-700 px-4 py-2 rounded-lg transition-colors">
                        Return to Login
                    </a>
                </div>
            </div>
        );
    }

    return (
        <div className="min-h-screen flex flex-col items-center justify-center bg-gray-900 text-white">
            <Loader2 size={48} className="text-blue-500 animate-spin mb-4" />
            <h2 className="text-xl font-medium">Authenticating with GitHub...</h2>
        </div>
    );
};
