import React, { useState, useEffect } from 'react';
import { Save, Shield, Server, Key } from 'lucide-react';
import { setApiKey as saveApiKey, getApiKey, client } from '../../utils/api';

export const SettingsSection: React.FC = () => {
    const [apiKey, setApiKey] = useState('');
    const [saved, setSaved] = useState(false);

    useEffect(() => {
        const key = getApiKey();
        if (key) setApiKey(key);
    }, []);

    const handleSave = () => {
        saveApiKey(apiKey);
        setSaved(true);
        setTimeout(() => setSaved(false), 2000);
    };

    return (
        <div className="max-w-2xl mx-auto space-y-8">
            <div className="bg-[#1e293b]/50 backdrop-blur-sm p-8 rounded-2xl border border-gray-800/50 shadow-xl">
                <div className="flex items-center space-x-3 mb-6">
                    <div className="p-2 bg-blue-500/10 rounded-lg">
                        <Key className="text-blue-400 w-6 h-6" />
                    </div>
                    <div>
                        <h2 className="text-xl font-bold text-white">API Configuration</h2>
                        <p className="text-gray-400 text-sm">Manage your connection credentials</p>
                    </div>
                </div>

                <div className="space-y-4">
                    <div>
                        <label className="block text-sm font-medium text-gray-300 mb-2">
                            Server API Key
                        </label>
                        <div className="relative">
                            <input
                                type="password"
                                value={apiKey}
                                onChange={(e) => setApiKey(e.target.value)}
                                className="w-full bg-gray-900/50 border border-gray-700 text-white rounded-xl px-4 py-3 focus:outline-none focus:ring-2 focus:ring-blue-500/50 focus:border-blue-500 transition-all font-mono"
                                placeholder="Enter your API key..."
                            />
                            <div className="absolute right-3 top-3 text-gray-500">
                                <Shield size={20} />
                            </div>
                        </div>
                        <p className="mt-2 text-xs text-gray-500">
                            Your API key is stored locally in your browser.
                        </p>
                    </div>

                    <button
                        onClick={handleSave}
                        className="w-full bg-blue-600 hover:bg-blue-500 text-white font-medium py-3 rounded-xl transition-all flex items-center justify-center space-x-2 shadow-lg shadow-blue-600/20"
                    >
                        <Save size={20} />
                        <span>{saved ? 'Saved Successfully!' : 'Save Configuration'}</span>
                    </button>
                </div>
            </div>

            <div className="bg-[#1e293b]/50 backdrop-blur-sm p-8 rounded-2xl border border-gray-800/50 shadow-xl">
                <div className="flex items-center space-x-3 mb-6">
                    <div className="p-2 bg-purple-500/10 rounded-lg">
                        <Server className="text-purple-400 w-6 h-6" />
                    </div>
                    <div>
                        <h2 className="text-xl font-bold text-white">Connect New Server</h2>
                        <p className="text-gray-400 text-sm">Deploy an agent to monitor another system</p>
                    </div>
                </div>

                <div className="space-y-4">
                    <div className="bg-gray-900/50 p-4 rounded-xl border border-gray-700/50">
                        <p className="text-sm text-gray-400 mb-2">1. Run this command on your target server:</p>
                        <div className="font-mono text-xs text-green-400 break-all bg-black/30 p-3 rounded-lg border border-white/5 select-all">
                            curl -sL {window.location.origin}/install.sh | bash -s -- -server {window.location.origin}
                        </div>
                        <p className="text-xs text-gray-500 mt-2">
                            The script will generate a unique <strong>Agent Token</strong>. Copy it.
                        </p>
                    </div>

                    <div className="bg-gray-900/50 p-4 rounded-xl border border-gray-700/50">
                        <p className="text-sm text-gray-400 mb-2">2. Paste the Agent Token here to link it:</p>
                        <div className="flex space-x-2">
                            <input
                                type="text"
                                placeholder="Paste Agent Token..."
                                className="flex-1 bg-black/30 border border-gray-700 rounded-lg px-3 py-2 text-sm text-white focus:outline-none focus:border-purple-500"
                                id="agent-token-input"
                            />
                            <button
                                onClick={async () => {
                                    const input = document.getElementById('agent-token-input') as HTMLInputElement;
                                    const token = input.value.trim();
                                    if (!token) return;

                                    try {
                                        await client.post('/agents', { token, name: 'New Agent' });
                                        input.value = '';
                                        alert('Agent linked successfully!');
                                    } catch (e) {
                                        alert('Failed to link agent');
                                    }
                                }}
                                className="bg-purple-600 hover:bg-purple-500 text-white px-4 py-2 rounded-lg text-sm font-medium transition-colors"
                            >
                                Link
                            </button>
                        </div>
                    </div>
                </div>
            </div>

            <div className="bg-[#1e293b]/50 backdrop-blur-sm p-8 rounded-2xl border border-gray-800/50 shadow-xl">
                <div className="flex items-center space-x-3 mb-6">
                    <div className="p-2 bg-green-500/10 rounded-lg">
                        <Server className="text-green-400 w-6 h-6" />
                    </div>
                    <div>
                        <h2 className="text-xl font-bold text-white">Connection Status</h2>
                        <p className="text-gray-400 text-sm">Current server connectivity</p>
                    </div>
                </div>

                <div className="flex items-center space-x-4 p-4 bg-green-500/10 border border-green-500/20 rounded-xl">
                    <div className="w-3 h-3 bg-green-500 rounded-full animate-pulse" />
                    <div>
                        <p className="text-green-400 font-medium">Connected</p>
                        <p className="text-green-500/70 text-sm">Receiving real-time updates</p>
                    </div>
                </div>
            </div>
        </div>
    );
};
