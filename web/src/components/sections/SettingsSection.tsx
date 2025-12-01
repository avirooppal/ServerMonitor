import React, { useState, useEffect } from 'react';
import { Server, Plus, Trash2, ExternalLink, Copy, Check } from 'lucide-react';
import { getSystems, saveSystem, deleteSystem, generateToken, type System } from '../../utils/api';

export const SettingsSection: React.FC = () => {
    const [systems, setSystems] = useState<System[]>([]);
    const [newSystem, setNewSystem] = useState({ name: '', url: '', apiKey: '' });
    const [generatedToken, setGeneratedToken] = useState('');
    const [copied, setCopied] = useState(false);
    const [useSSL, setUseSSL] = useState(false);
    const [domain, setDomain] = useState('');

    useEffect(() => {
        loadSystems();
        generateNewToken();
    }, []);

    const loadSystems = () => {
        setSystems(getSystems());
    };

    const generateNewToken = () => {
        const token = generateToken();
        setGeneratedToken(token);
        setNewSystem(prev => ({ ...prev, apiKey: token }));
    };

    const handleCopyCommand = () => {
        let cmd = `curl -sL https://raw.githubusercontent.com/avirooppal/ServerMonitor/main/web/public/setup.sh | bash -s -- ${generatedToken}`;
        if (useSSL) {
            const d = domain || '<YOUR_VPS_IP>.nip.io';
            cmd += ` ${d}`;
        }
        navigator.clipboard.writeText(cmd);
        setCopied(true);
        setTimeout(() => setCopied(false), 2000);
    };

    const handleAddSystem = () => {
        if (!newSystem.name || !newSystem.url || !newSystem.apiKey) return;

        const system: System = {
            id: Date.now().toString(),
            name: newSystem.name.trim(),
            url: newSystem.url.trim().replace(/\/$/, ''), // Remove trailing slash
            api_key: newSystem.apiKey.trim(),
            created_at: new Date().toISOString()
        };

        saveSystem(system);
        setNewSystem({ name: '', url: '', apiKey: '' });
        generateNewToken(); // Reset for next
        loadSystems();
    };

    const handleDeleteSystem = (id: string) => {
        if (!confirm('Are you sure you want to remove this system?')) return;
        deleteSystem(id);
        loadSystems();
    };

    return (
        <div className="max-w-4xl mx-auto space-y-8 pb-10">
            {/* Add New System Flow */}
            <div className="bg-[#1e293b]/50 backdrop-blur-sm p-8 rounded-2xl border border-gray-800/50 shadow-xl">
                <div className="flex items-center space-x-3 mb-6">
                    <div className="p-2 bg-purple-500/10 rounded-lg">
                        <Plus className="text-purple-400 w-6 h-6" />
                    </div>
                    <div>
                        <h2 className="text-xl font-bold text-white">Add New Server</h2>
                        <p className="text-gray-400 text-sm">Install the agent and connect it to your dashboard</p>
                    </div>
                </div>

                {/* Step 1: Install Agent */}
                <div className="mb-8">
                    <div className="flex items-center justify-between mb-2">
                        <h3 className="text-sm font-semibold text-blue-400 uppercase tracking-wider">Step 1: Install Agent</h3>
                        <div className="flex items-center space-x-2">
                            <label className="text-xs text-gray-400 flex items-center space-x-2 cursor-pointer">
                                <input
                                    type="checkbox"
                                    checked={useSSL}
                                    onChange={(e) => setUseSSL(e.target.checked)}
                                    className="rounded border-gray-600 bg-gray-800 text-purple-500 focus:ring-purple-500/50"
                                />
                                <span>Enable SSL (HTTPS)</span>
                            </label>
                        </div>
                    </div>

                    {useSSL && (
                        <div className="mb-3">
                            <input
                                type="text"
                                placeholder="Enter Domain or VPS IP (e.g. 1.2.3.4)"
                                value={domain}
                                onChange={(e) => setDomain(e.target.value)}
                                className="w-full bg-black/30 border border-gray-700 rounded-lg px-3 py-2 text-sm text-white focus:outline-none focus:border-blue-500"
                            />
                            <p className="text-[10px] text-gray-500 mt-1">
                                No domain? Enter your IP, and we'll use <b>{domain || '1.2.3.4'}.nip.io</b> for SSL.
                            </p>
                        </div>
                    )}

                    <div className="bg-black/30 p-4 rounded-xl border border-blue-500/10 flex items-center justify-between group relative overflow-hidden">
                        <code className="font-mono text-sm text-blue-300 break-all select-all z-10">
                            curl -sL https://raw.githubusercontent.com/avirooppal/ServerMonitor/main/web/public/setup.sh | bash -s -- {generatedToken} {useSSL ? (domain.match(/^\d+\.\d+\.\d+\.\d+$/) ? `${domain}.nip.io` : (domain || '<YOUR_DOMAIN>')) : ''}
                        </code>
                        <button
                            onClick={handleCopyCommand}
                            className="ml-4 text-gray-400 hover:text-white transition-colors p-2 hover:bg-white/5 rounded-lg z-10"
                            title="Copy to clipboard"
                        >
                            {copied ? <Check className="w-5 h-5 text-green-400" /> : <Copy className="w-5 h-5" />}
                        </button>
                    </div>
                    <p className="text-xs text-gray-500 mt-2">
                        Run this command on your VPS. {useSSL ? 'It will set up Caddy for HTTPS.' : 'It will start the agent on HTTP port 8080.'}
                    </p>
                </div>

                {/* Step 2: Connect */}
                <div>
                    <h3 className="text-sm font-semibold text-purple-400 mb-4 uppercase tracking-wider">Step 2: Connect</h3>
                    <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                        <div>
                            <label className="block text-xs text-gray-500 mb-1">System Name</label>
                            <input
                                type="text"
                                placeholder="e.g. Production DB"
                                value={newSystem.name}
                                onChange={(e) => setNewSystem({ ...newSystem, name: e.target.value })}
                                className="w-full bg-black/30 border border-gray-700 rounded-lg px-3 py-2 text-sm text-white focus:outline-none focus:border-purple-500 transition-colors"
                            />
                        </div>
                        <div>
                            <label className="block text-xs text-gray-500 mb-1">Agent URL</label>
                            <input
                                type="text"
                                placeholder="e.g. http://1.2.3.4:8080"
                                value={newSystem.url}
                                onChange={(e) => setNewSystem({ ...newSystem, url: e.target.value })}
                                className="w-full bg-black/30 border border-gray-700 rounded-lg px-3 py-2 text-sm text-white focus:outline-none focus:border-purple-500 transition-colors"
                            />
                        </div>
                    </div>
                    <div className="mt-4 flex justify-end">
                        <button
                            onClick={handleAddSystem}
                            className="bg-purple-600 hover:bg-purple-500 text-white px-6 py-2 rounded-lg text-sm font-medium transition-colors shadow-lg shadow-purple-600/20 flex items-center space-x-2"
                        >
                            <Server size={16} />
                            <span>Connect Server</span>
                        </button>
                    </div>
                </div>
            </div>

            {/* Systems List */}
            <div className="bg-[#1e293b]/50 backdrop-blur-sm p-8 rounded-2xl border border-gray-800/50 shadow-xl">
                <div className="flex items-center justify-between mb-6">
                    <div className="flex items-center space-x-3">
                        <div className="p-2 bg-blue-500/10 rounded-lg">
                            <Server className="text-blue-400 w-6 h-6" />
                        </div>
                        <div>
                            <h2 className="text-xl font-bold text-white">Monitored Systems</h2>
                            <p className="text-gray-400 text-sm">Manage your connected agents</p>
                        </div>
                    </div>
                </div>

                <div className="space-y-3">
                    {systems.length === 0 ? (
                        <div className="text-center py-8 text-gray-500 border border-dashed border-gray-700 rounded-xl">
                            No systems added yet. Follow the steps above to add one!
                        </div>
                    ) : (
                        systems.map((sys) => (
                            <div key={sys.id} className="flex items-center justify-between bg-gray-800/30 p-4 rounded-xl border border-white/5 hover:border-white/10 transition-colors group">
                                <div className="flex items-center space-x-4">
                                    <div className="w-2 h-2 bg-green-500 rounded-full shadow-[0_0_8px_rgba(34,197,94,0.5)]" />
                                    <div>
                                        <h4 className="font-medium text-white group-hover:text-blue-400 transition-colors">{sys.name}</h4>
                                        <div className="flex items-center space-x-2 text-xs text-gray-500">
                                            <span className="font-mono">{sys.url}</span>
                                            <span className="text-gray-700">•</span>
                                            <span>Added {new Date(sys.created_at).toLocaleDateString()}</span>
                                        </div>
                                    </div>
                                </div>
                                <div className="flex items-center space-x-3 opacity-50 group-hover:opacity-100 transition-opacity">
                                    <a
                                        href={`${sys.url}/api/v1/ping`}
                                        target="_blank"
                                        rel="noreferrer"
                                        className="p-2 text-gray-400 hover:text-blue-400 transition-colors bg-white/5 rounded-lg"
                                        title="Test Connection"
                                    >
                                        <ExternalLink size={16} />
                                    </a>
                                    <button
                                        onClick={() => handleDeleteSystem(sys.id)}
                                        className="p-2 text-gray-400 hover:text-red-400 transition-colors bg-white/5 rounded-lg hover:bg-red-500/10"
                                        title="Remove System"
                                    >
                                        <Trash2 size={16} />
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
