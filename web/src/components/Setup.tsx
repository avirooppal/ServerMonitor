import React, { useState } from 'react';
import { verifyKey, setApiKey } from '../utils/api';
import { Key, Server, ShieldCheck, Loader2 } from 'lucide-react';

interface SetupProps {
    onComplete: () => void;
}

export const Setup: React.FC<SetupProps> = ({ onComplete }) => {
    const [key, setKey] = useState('');
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState('');

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        setLoading(true);
        setError('');

        try {
            const isValid = await verifyKey(key);
            if (isValid) {
                setApiKey(key);
                onComplete();
            } else {
                setError('Invalid API Key');
            }
        } catch (err) {
            setError('Failed to verify key. Check your connection.');
        } finally {
            setLoading(false);
        }
    };

    return (
        <div className="min-h-screen flex items-center justify-center bg-background p-4">
            <div className="max-w-md w-full bg-surface border border-white/10 rounded-2xl shadow-2xl overflow-hidden">
                <div className="p-8">
                    <div className="flex justify-center mb-6">
                        <div className="p-4 bg-primary/10 rounded-full">
                            <Server size={48} className="text-primary" />
                        </div>
                    </div>

                    <h1 className="text-2xl font-bold text-center text-white mb-2">Welcome to ServerMoni</h1>
                    <p className="text-gray-400 text-center mb-8">Enter your Master API Key to access the dashboard.</p>

                    <form onSubmit={handleSubmit} className="space-y-6">
                        <div className="space-y-2">
                            <label className="text-sm font-medium text-gray-300 flex items-center gap-2">
                                <Key size={16} /> Master API Key
                            </label>
                            <input
                                type="password"
                                value={key}
                                onChange={(e) => setKey(e.target.value)}
                                className="w-full bg-background border border-white/10 rounded-lg px-4 py-3 text-white focus:outline-none focus:border-primary focus:ring-1 focus:ring-primary transition-colors"
                                placeholder="Enter key..."
                                required
                            />
                        </div>

                        {error && (
                            <div className="p-3 bg-danger/10 border border-danger/20 rounded-lg flex items-center gap-2 text-danger text-sm">
                                <ShieldCheck size={16} />
                                {error}
                            </div>
                        )}

                        <button
                            type="submit"
                            disabled={loading}
                            className="w-full bg-primary hover:bg-primary-hover text-white font-semibold py-3 rounded-lg transition-all duration-200 flex items-center justify-center gap-2 disabled:opacity-50 disabled:cursor-not-allowed"
                        >
                            {loading ? (
                                <>
                                    <Loader2 size={20} className="animate-spin" /> Verifying...
                                </>
                            ) : (
                                'Connect to Dashboard'
                            )}
                        </button>
                    </form>
                </div>

                <div className="bg-background/50 p-4 text-center border-t border-white/5">
                    <p className="text-xs text-gray-500">
                        Don't have a key? Check your server logs: <br />
                        <code className="text-primary">docker exec server-moni-dashboard-1 cat data/api_key.txt</code>
                    </p>
                </div>
            </div>
        </div>
    );
};
