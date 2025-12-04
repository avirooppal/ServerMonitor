import React from 'react';
import type { SystemMetrics } from '../../types';
import { Box, Play, Square, AlertCircle, Cpu, Layers } from 'lucide-react';
import clsx from 'clsx';

interface DockerSectionProps {
    metrics: SystemMetrics;
    systemId: number;
}

export const DockerSection: React.FC<DockerSectionProps> = ({ metrics }) => {
    const containers = metrics.containers || [];

    const getStatusColor = (state: string) => {
        switch (state.toLowerCase()) {
            case 'running': return 'text-success';
            case 'exited': return 'text-gray-500';
            case 'paused': return 'text-warning';
            case 'restarting': return 'text-info';
            case 'dead': return 'text-danger';
            default: return 'text-gray-400';
        }
    };

    const getStatusIcon = (state: string) => {
        switch (state.toLowerCase()) {
            case 'running': return <Play size={14} className="fill-current" />;
            case 'exited': return <Square size={14} className="fill-current" />;
            default: return <AlertCircle size={14} />;
        }
    };

    return (
        <div className="space-y-6">
            <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
                <div className="bg-surface border border-white/5 rounded-xl p-4 flex items-center justify-between">
                    <div>
                        <p className="text-gray-400 text-sm font-medium">Total Containers</p>
                        <p className="text-2xl font-bold text-white mt-1">{containers.length}</p>
                    </div>
                    <div className="p-3 bg-primary/10 rounded-lg">
                        <Box size={24} className="text-primary" />
                    </div>
                </div>
                <div className="bg-surface border border-white/5 rounded-xl p-4 flex items-center justify-between">
                    <div>
                        <p className="text-gray-400 text-sm font-medium">Running</p>
                        <p className="text-2xl font-bold text-success mt-1">
                            {containers.filter(c => c.state === 'running').length}
                        </p>
                    </div>
                    <div className="p-3 bg-success/10 rounded-lg">
                        <Play size={24} className="text-success fill-current" />
                    </div>
                </div>
                <div className="bg-surface border border-white/5 rounded-xl p-4 flex items-center justify-between">
                    <div>
                        <p className="text-gray-400 text-sm font-medium">Stopped</p>
                        <p className="text-2xl font-bold text-gray-400 mt-1">
                            {containers.filter(c => c.state !== 'running').length}
                        </p>
                    </div>
                    <div className="p-3 bg-white/5 rounded-lg">
                        <Square size={24} className="text-gray-400 fill-current" />
                    </div>
                </div>
            </div>

            <div className="bg-surface border border-white/5 rounded-xl overflow-hidden">
                <div className="px-6 py-4 border-b border-white/5 flex items-center justify-between">
                    <h3 className="text-lg font-semibold text-white flex items-center gap-2">
                        <Box size={20} className="text-primary" />
                        Containers
                    </h3>
                    <span className="text-xs text-gray-500 bg-white/5 px-2 py-1 rounded">
                        {containers.length} Total
                    </span>
                </div>

                <div className="overflow-x-auto">
                    <table className="w-full text-left border-collapse">
                        <thead>
                            <tr className="bg-white/5 text-gray-400 text-xs uppercase tracking-wider">
                                <th className="px-6 py-3 font-medium">Name / ID</th>
                                <th className="px-6 py-3 font-medium">Image</th>
                                <th className="px-6 py-3 font-medium">State</th>
                                <th className="px-6 py-3 font-medium text-right">CPU</th>
                                <th className="px-6 py-3 font-medium text-right">Memory</th>
                            </tr>
                        </thead>
                        <tbody className="divide-y divide-white/5">
                            {containers.length === 0 ? (
                                <tr>
                                    <td colSpan={5} className="px-6 py-8 text-center text-gray-500">
                                        No containers found.
                                    </td>
                                </tr>
                            ) : (
                                containers.map((container) => (
                                    <tr key={container.id} className="hover:bg-white/5 transition-colors group">
                                        <td className="px-6 py-4">
                                            <div className="flex flex-col">
                                                <span className="font-medium text-white group-hover:text-primary transition-colors">
                                                    {container.name.replace(/^\//, '')}
                                                </span>
                                                <span className="text-xs text-gray-500 font-mono">
                                                    {container.id.substring(0, 12)}
                                                </span>
                                            </div>
                                        </td>
                                        <td className="px-6 py-4">
                                            <span className="text-sm text-gray-300 font-mono bg-white/5 px-2 py-1 rounded">
                                                {container.image}
                                            </span>
                                        </td>
                                        <td className="px-6 py-4">
                                            <div className={clsx("flex items-center gap-2 text-sm font-medium", getStatusColor(container.state))}>
                                                {getStatusIcon(container.state)}
                                                <span className="capitalize">{container.state}</span>
                                            </div>
                                            <div className="text-xs text-gray-500 mt-0.5">
                                                {container.status}
                                            </div>
                                        </td>
                                        <td className="px-6 py-4 text-right">
                                            <div className="flex items-center justify-end gap-2 text-sm text-gray-300">
                                                <Cpu size={14} className="text-gray-500" />
                                                {container.cpu_percent.toFixed(2)}%
                                            </div>
                                        </td>
                                        <td className="px-6 py-4 text-right">
                                            <div className="flex items-center justify-end gap-2 text-sm text-gray-300">
                                                <Layers size={14} className="text-gray-500" />
                                                {(container.memory_usage / 1024 / 1024).toFixed(1)} MB
                                            </div>
                                            <div className="text-xs text-gray-500 mt-0.5">
                                                Limit: {(container.memory_limit / 1024 / 1024 / 1024).toFixed(1)} GB
                                            </div>
                                        </td>
                                    </tr>
                                ))
                            )}
                        </tbody>
                    </table>
                </div>
            </div>
        </div>
    );
};
