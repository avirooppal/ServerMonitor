import React, { useState, useEffect } from 'react';
import { Server, Trash2, ExternalLink } from 'lucide-react';
import { fetchSystems, deleteSystem, type System } from '../../utils/api';

export const SettingsSection: React.FC = () => {
    const [systems, setSystems] = useState<System[]>([]);

    useEffect(() => {
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

                {/* Systems List */}
                <div className="space-y-3">
                    {systems.length === 0 ? (
                        <div className="text-center py-8 text-gray-500">
                            No systems added yet. Use the "+ Add Server" button in the dashboard header.
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
