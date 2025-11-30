import React, { useState } from 'react';
import type { SystemMetrics } from '../../types';
import { Box, Play, Square, FileText, X, Loader } from 'lucide-react';
import clsx from 'clsx';

interface DockerSectionProps {
    metrics: SystemMetrics;
    systemId: number;
    apiKey: string;
}

export const DockerSection: React.FC<DockerSectionProps> = ({ metrics, systemId, apiKey }) => {
    const containers = metrics.containers || [];
    const [selectedContainer, setSelectedContainer] = useState<string | null>(null);
    const [logs, setLogs] = useState<string>('');
    const [loadingLogs, setLoadingLogs] = useState(false);

    const formatBytes = (bytes: number) => {
        if (bytes === 0) return '0 B';
        const k = 1024;
        const sizes = ['B', 'KB', 'MB', 'GB', 'TB'];
        const i = Math.floor(Math.log(bytes) / Math.log(k));
        return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
    };

    const fetchLogs = async (containerId: string) => {
        setSelectedContainer(containerId);
        setLoadingLogs(true);
        setLogs('');
        try {
            const headers = { Authorization: `Bearer ${apiKey}` };
            // Proxy request to get logs
            const response = await fetch(`/api/v1/systems/${systemId}/proxy?path=/docker/containers/${containerId}/logs`, { headers });
            if (response.ok) {
                const text = await response.text();
                setLogs(text);
            } else {
                setLogs('Failed to fetch logs.');
            }
        } catch (error) {
            console.error("Failed to fetch logs", error);
            setLogs('Error fetching logs.');
        } finally {
            setLoadingLogs(false);
        }
    };

    return (
        <div className="space-y-6 relative">
            <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
                <div className="bg-[#1e293b]/50 backdrop-blur-sm p-6 rounded-2xl border border-gray-800/50 shadow-xl">
                    <div className="flex items-center space-x-3 mb-2">
                        <div className="p-2 bg-blue-500/10 rounded-lg">
                            <Box className="text-blue-400 w-5 h-5" />
                        </div>
                        <h3 className="text-gray-400 font-medium">Total Containers</h3>
                    </div>
                    <p className="text-3xl font-bold text-white">{containers.length}</p>
                </div>
                <div className="bg-[#1e293b]/50 backdrop-blur-sm p-6 rounded-2xl border border-gray-800/50 shadow-xl">
                    <div className="flex items-center space-x-3 mb-2">
                        <div className="p-2 bg-green-500/10 rounded-lg">
                            <Play className="text-green-400 w-5 h-5" />
                        </div>
                        <h3 className="text-gray-400 font-medium">Running</h3>
                    </div>
                    <p className="text-3xl font-bold text-white">
                        {containers.filter(c => c.state === 'running').length}
                    </p>
                </div>
                <div className="bg-[#1e293b]/50 backdrop-blur-sm p-6 rounded-2xl border border-gray-800/50 shadow-xl">
                    <div className="flex items-center space-x-3 mb-2">
                        <div className="p-2 bg-red-500/10 rounded-lg">
                            <Square className="text-red-400 w-5 h-5" />
                        </div>
                        <h3 className="text-gray-400 font-medium">Stopped</h3>
                    </div>
                    <p className="text-3xl font-bold text-white">
                        {containers.filter(c => c.state !== 'running').length}
                    </p>
                </div>
            </div>

            <div className="bg-[#1e293b]/50 backdrop-blur-sm rounded-2xl border border-gray-800/50 shadow-xl overflow-hidden">
                <div className="p-6 border-b border-gray-800/50">
                    <h3 className="text-lg font-semibold text-white">Container List</h3>
                </div>
                <div className="overflow-x-auto">
                    <table className="w-full text-left">
                        <thead>
                            <tr className="bg-gray-900/50 text-gray-400 text-sm uppercase tracking-wider">
                                <th className="px-6 py-4 font-medium">Name / ID</th>
                                <th className="px-6 py-4 font-medium">Image</th>
                                <th className="px-6 py-4 font-medium">State</th>
                                <th className="px-6 py-4 font-medium">CPU %</th>
                                <th className="px-6 py-4 font-medium">Mem Usage</th>
                                <th className="px-6 py-4 font-medium">Actions</th>
                            </tr>
                        </thead>
                        <tbody className="divide-y divide-gray-800/50">
                            {containers.map((container) => (
                                <tr key={container.id} className="hover:bg-gray-800/30 transition-colors">
                                    <td className="px-6 py-4">
                                        <div className="flex flex-col">
                                            <span className="font-medium text-white">{container.name}</span>
                                            <span className="text-xs text-gray-500 font-mono">{container.id.substring(0, 12)}</span>
                                        </div>
                                    </td>
                                    <td className="px-6 py-4 text-gray-300 font-mono text-sm">{container.image}</td>
                                    <td className="px-6 py-4">
                                        <span className={clsx(
                                            "inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium border",
                                            container.state === 'running'
                                                ? "bg-green-500/10 text-green-400 border-green-500/20"
                                                : "bg-gray-700/30 text-gray-400 border-gray-600/30"
                                        )}>
                                            {container.state}
                                        </span>
                                    </td>
                                    <td className="px-6 py-4">
                                        <div className="flex items-center space-x-2">
                                            <span className="text-white font-mono text-sm">
                                                {container.cpu_percent ? container.cpu_percent.toFixed(2) : '0.00'}%
                                            </span>
                                        </div>
                                    </td>
                                    <td className="px-6 py-4">
                                        <div className="flex flex-col">
                                            <span className="text-white font-mono text-sm">
                                                {formatBytes(container.memory_usage)}
                                            </span>
                                            <span className="text-xs text-gray-500">
                                                Limit: {formatBytes(container.memory_limit)}
                                            </span>
                                        </div>
                                    </td>
                                    <td className="px-6 py-4">
                                        <button
                                            onClick={() => fetchLogs(container.id)}
                                            className="text-blue-400 hover:text-blue-300 transition-colors p-2 hover:bg-blue-500/10 rounded-lg"
                                            title="View Logs"
                                        >
                                            <FileText size={18} />
                                        </button>
                                    </td>
                                </tr>
                            ))}
                            {containers.length === 0 && (
                                <tr>
                                    <td colSpan={6} className="px-6 py-12 text-center text-gray-500">
                                        No containers found or Docker not available.
                                    </td>
                                </tr>
                            )}
                        </tbody>
                    </table>
                </div>
            </div>

            {/* Logs Modal */}
            {selectedContainer && (
                <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/50 backdrop-blur-sm p-4">
                    <div className="bg-[#1e293b] w-full max-w-4xl rounded-2xl shadow-2xl border border-gray-700 flex flex-col max-h-[80vh]">
                        <div className="flex items-center justify-between p-4 border-b border-gray-700">
                            <h3 className="text-lg font-semibold text-white flex items-center gap-2">
                                <FileText className="w-5 h-5 text-blue-400" />
                                Container Logs
                            </h3>
                            <button
                                onClick={() => setSelectedContainer(null)}
                                className="text-gray-400 hover:text-white transition-colors"
                            >
                                <X size={20} />
                            </button>
                        </div>
                        <div className="p-4 flex-1 overflow-auto bg-black/30 font-mono text-sm text-gray-300 whitespace-pre-wrap">
                            {loadingLogs ? (
                                <div className="flex items-center justify-center h-40">
                                    <Loader className="w-8 h-8 text-blue-400 animate-spin" />
                                </div>
                            ) : (
                                logs || "No logs available."
                            )}
                        </div>
                    </div>
                </div>
            )}
        </div>
    );
};
