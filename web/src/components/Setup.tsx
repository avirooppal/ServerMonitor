import React, { useState } from 'react';
import { setApiKey, verifyKey } from '../utils/api';
import { Key } from 'lucide-react';

interface SetupProps {
    onComplete: () => void;
}

export const Setup: React.FC<SetupProps> = ({ onComplete }) => {
    const [key, setKeyInput] = useState('');
    const [error, setError] = useState('');
    const [loading, setLoading] = useState(false);

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        setLoading(true);
        setError('');

        try {
            await verifyKey(key);
            setApiKey(key);
            onComplete();
        } catch (err) {
            setError('Invalid API Key or Server Unreachable');
        } finally {
            setLoading(false);
        }
    };

    return (
        <div className="flex items-center justify-center h-screen bg-gray-900">
            <div className="bg-gray-800 p-8 rounded-lg shadow-lg max-w-md w-full border border-gray-700">
                <div className="flex justify-center mb-6">
                    <div className="p-3 bg-blue-600 rounded-full">
                        <Key className="w-8 h-8 text-white" />
                    </div>
                </div>
                <h2 className="text-2xl font-bold text-center text-white mb-6">Server Moni Setup</h2>
                <form onSubmit={handleSubmit} className="space-y-4">
                    <div>
                        <label className="block text-sm font-medium text-gray-400 mb-1">API Key</label>
                        <input
                            type="text"
                            value={key}
                            onChange={(e) => setKeyInput(e.target.value)}
                            className="w-full px-4 py-2 bg-gray-700 border border-gray-600 rounded text-white focus:outline-none focus:border-blue-500"
                            placeholder="Paste your API Key here"
                            required
                        />
                    </div>
                    {error && <p className="text-red-500 text-sm text-center">{error}</p>}
                    <button
                        type="submit"
                        disabled={loading}
                        className="w-full py-2 px-4 bg-blue-600 hover:bg-blue-700 text-white font-semibold rounded transition duration-200 disabled:opacity-50"
                    >
                        {loading ? 'Verifying...' : 'Connect'}
                    </button>
                </form>
                <div className="mt-6 p-4 bg-black/30 rounded-lg border border-gray-700/50">
                    <p className="text-xs text-gray-400 mb-2 text-center">To get your API Key, run this in your server terminal:</p>
                    <div className="bg-black/50 p-2 rounded border border-gray-700 text-center select-all cursor-text">
                        <code className="text-xs font-mono text-green-400">type agent_data\api_key.txt</code>
                    </div>
                </div>
            </div>
        </div>
    );
};
