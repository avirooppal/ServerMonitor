import React, { useState, useEffect } from 'react';
import { Save, Shield, Server, Key } from 'lucide-react';
import { setApiKey as saveApiKey, getApiKey } from '../../utils/api';

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
