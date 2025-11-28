import React, { useState, useEffect } from 'react';
import { Save, Shield, Server, Key, Plus, Trash2, ExternalLink } from 'lucide-react';
import { setApiKey as saveApiKey, getApiKey, fetchSystems, addSystem, deleteSystem, type System } from '../../utils/api';

export const SettingsSection: React.FC = () => {
    const [apiKey, setApiKey] = useState('');
    const [saved, setSaved] = useState(false);
    const [systems, setSystems] = useState<System[]>([]);
    const [newSystem, setNewSystem] = useState({ name: '', url: '', apiKey: '' });
    const [loading, setLoading] = useState(false);

    useEffect(() => {
        const key = getApiKey();
        if (key) setApiKey(key);
        loadSystems();
    }, []);

    const loadSystems = async () => {
        try {
            const list = await fetchSystems();
            setSystems(list);
        } catch (e) {
            console.error("Failed to load systems", e);
        }
    };

    const handleSaveKey = () => {
        saveApiKey(apiKey);
        setSaved(true);
        setTimeout(() => setSaved(false), 2000);
    };

    const handleAddSystem = async () => {
        if (!newSystem.name || !newSystem.url || !newSystem.apiKey) return;
        setLoading(true);
        try {
            await addSystem(newSystem.name, newSystem.url, newSystem.apiKey);
            setNewSystem({ name: '', url: '', apiKey: '' });
            await loadSystems();
        } catch (e: any) {
            console.error(e);
            const msg = e.response?.data?.error || e.message || 'Failed to add system';
            alert(`Error: ${msg}`);
        } finally {
            setLoading(false);
        }
    };

    const handleDeleteSystem = async (id: number) => {
        if (!confirm('Are you sure you want to remove this system?')) return;
        try {
            await deleteSystem(id);
            await loadSystems();
        } catch (e) {
            alert('Failed to delete system');
        }
    };

    return (
        <div className="max-w-4xl mx-auto space-y-8 pb-10">
            {/* API Key Config */}
            <div className="bg-[#1e293b]/50 backdrop-blur-sm p-8 rounded-2xl border border-gray-800/50 shadow-xl">
                <div className="flex items-center space-x-3 mb-6">
                    <div className="p-2 bg-blue-500/10 rounded-lg">
                        <Key className="text-blue-400 w-6 h-6" />
                    </div>
                    <div>
                        <h2 className="text-xl font-bold text-white">Dashboard Access</h2>
                        <p className="text-gray-400 text-sm">Your master API key for this dashboard</p>
                    </div>
                </div>

                <div className="flex space-x-4">
                    <div className="relative flex-1">
                        <input
                            type="password"
                            value={apiKey}
                            onChange={(e) => setApiKey(e.target.value)}
                            className="w-full bg-gray-900/50 border border-gray-700 text-white rounded-xl px-4 py-3 focus:outline-none focus:ring-2 focus:ring-blue-500/50 focus:border-blue-500 transition-all font-mono"
                            placeholder="Enter your Master API key..."
                        />
                        <div className="absolute right-3 top-3 text-gray-500">
                            <Shield size={20} />
                        </div>
                    </div>
                    <button
                        onClick={handleSaveKey}
                        className="bg-blue-600 hover:bg-blue-500 text-white font-medium px-6 rounded-xl transition-all flex items-center space-x-2 shadow-lg shadow-blue-600/20"
                    >
                        <Save size={20} />
                        <span>{saved ? 'Saved!' : 'Save'}</span>
                    </button>
                </div>
            </div>

            {/* Systems Management */}
            <div className="bg-[#1e293b]/50 backdrop-blur-sm p-8 rounded-2xl border border-gray-800/50 shadow-xl">
                <div className="flex items-center justify-between mb-6">
                    <div className="flex items-center space-x-3">
                        <div className="p-2 bg-purple-500/10 rounded-lg">
                            <Server className="text-purple-400 w-6 h-6" />
                        </div>
                        <div>
                            <h2 className="text-xl font-bold text-white">Monitored Systems</h2>
                            <p className="text-gray-400 text-sm">Manage your connected agents</p>
                        </div>
                    </div>
                </div>

                {/* Add New System */}
                <div className="bg-gray-900/50 p-6 rounded-xl border border-gray-700/50 mb-8">
                    <h3 className="text-sm font-semibold text-gray-300 mb-4 flex items-center">
                        <Plus size={16} className="mr-2" /> Add New System
                    </h3>
                    <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
                        <input
                            type="text"
                            placeholder="System Name (e.g. Prod DB)"
                            value={newSystem.name}
                            onChange={(e) => setNewSystem({ ...newSystem, name: e.target.value })}
                            className="bg-black/30 border border-gray-700 rounded-lg px-3 py-2 text-sm text-white focus:outline-none focus:border-purple-500"
                        />
                        <input
                            type="text"
                            placeholder="Agent URL (e.g. http://1.2.3.4:8080)"
                            value={newSystem.url}
                            onChange={(e) => setNewSystem({ ...newSystem, url: e.target.value })}
                            className="bg-black/30 border border-gray-700 rounded-lg px-3 py-2 text-sm text-white focus:outline-none focus:border-purple-500"
                        />
                        <input
                            type="password"
                            placeholder="Agent API Key"
                            value={newSystem.apiKey}
                            onChange={(e) => setNewSystem({ ...newSystem, apiKey: e.target.value })}
                            className="bg-black/30 border border-gray-700 rounded-lg px-3 py-2 text-sm text-white focus:outline-none focus:border-purple-500"
                        />
                    </div>
                    <div className="mt-4 flex justify-end">
                        <button
                            onClick={handleAddSystem}
                            disabled={loading}
                            className="bg-purple-600 hover:bg-purple-500 text-white px-6 py-2 rounded-lg text-sm font-medium transition-colors disabled:opacity-50"
                        >
                            {loading ? 'Adding...' : 'Add System'}
                        </button>
                    </div>
                </div>

                {/* Agent Helper Script */}
                <div className="bg-blue-500/10 p-6 rounded-xl border border-blue-500/20 mb-8">
                    <h3 className="text-sm font-semibold text-blue-400 mb-2 flex items-center">
                        <Server size={16} className="mr-2" /> Easy Setup Script
                    </h3>
                    <p className="text-gray-400 text-xs mb-3">
                        Run this command on your agent server to instantly get the URL and API Key:
                    </p>
                    <div className="bg-black/30 p-3 rounded-lg border border-blue-500/10 flex items-center justify-between group">
                        <code className="font-mono text-xs text-blue-300 break-all select-all">
                            curl -sL {window.location.origin}/get-key.sh | bash
                        </code>
                        <button
                            onClick={() => {
                                navigator.clipboard.writeText(`curl -sL ${window.location.origin}/get-key.sh | bash`);
                                alert('Command copied to clipboard!');
                            }}
                            className="ml-4 text-gray-500 hover:text-white transition-colors"
                            title="Copy to clipboard"
                        >
                            <Save className="w-4 h-4" />
                        </button>
                    </div>
                </div>

                {/* Systems List */}
                <div className="space-y-3">
                    {systems.length === 0 ? (
                        <div className="text-center py-8 text-gray-500">
                            No systems added yet. Add one above!
                        </div>
                    ) : (
                        systems.map((sys) => (
                            <div key={sys.id} className="flex items-center justify-between bg-gray-800/30 p-4 rounded-xl border border-white/5 hover:border-white/10 transition-colors">
                                <div className="flex items-center space-x-4">
                                    <div className="w-2 h-2 bg-green-500 rounded-full" />
                                    <div>
                                        <h4 className="font-medium text-white">{sys.name}</h4>
                                        <div className="flex items-center space-x-2 text-xs text-gray-500">
                                            <span>{sys.url}</span>
                                            <span className="text-gray-700">•</span>
                                            <span>Added {new Date(sys.created_at).toLocaleDateString()}</span>
                                        </div>
                                    </div>
                                </div>
                                <div className="flex items-center space-x-3">
                                    <a
                                        href={`${sys.url}/api/v1/ping`}
                                        target="_blank"
                                        rel="noreferrer"
                                        className="p-2 text-gray-400 hover:text-blue-400 transition-colors"
                                        title="Test Connection"
                                    >
                                        <ExternalLink size={18} />
                                    </a>
                                    <button
                                        onClick={() => handleDeleteSystem(sys.id)}
                                        className="p-2 text-gray-400 hover:text-red-400 transition-colors"
                                        title="Remove System"
                                    >
                                        <Trash2 size={18} />
                                    </button>
                                </div>
                            </div>
                        ))
                    )}
                </div>
            </div>
        </div>
    );
};
