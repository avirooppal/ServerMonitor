import React from 'react';
import type { SystemMetrics } from '../../types';
import { CircularProgress } from '../CircularProgress';
import { CheckCircle, ArrowRight } from 'lucide-react';

interface OverviewSectionProps {
    metrics: SystemMetrics;
    onNavigate: (tab: string) => void;
}

export const OverviewSection: React.FC<OverviewSectionProps> = ({ metrics, onNavigate }) => {
    const formatBytes = (bytes: number) => {
        if (bytes === 0) return '0 B';
        const k = 1024;
        const sizes = ['B', 'KB', 'MB', 'GB', 'TB'];
        const i = Math.floor(Math.log(bytes) / Math.log(k));
        return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
    };

    const formatUptime = (seconds: number) => {
        const days = Math.floor(seconds / 86400);
        const hours = Math.floor((seconds % 86400) / 3600);
        const minutes = Math.floor((seconds % 3600) / 60);
        if (days > 0) return `${days} days, ${hours}h`;
        return `${hours}h ${minutes}m`;
    };

    return (
        <div className="space-y-6">
            {/* Top Info Banner */}
            <div className="bg-surface/50 backdrop-blur-sm p-4 rounded-2xl border border-white/5 shadow-xl flex flex-wrap justify-between items-center text-sm">
                <div className="flex flex-col px-4 border-r border-white/5">
                    <span className="text-gray-500 uppercase text-xs font-semibold tracking-wider mb-1">System Uptime</span>
                    <span className="text-white font-medium font-mono">{metrics.host_info ? formatUptime(metrics.host_info.uptime) : '-'}</span>
                </div>
                <div className="flex flex-col px-4 border-r border-white/5">
                    <span className="text-gray-500 uppercase text-xs font-semibold tracking-wider mb-1">OS / Arch</span>
                    <span className="text-white font-medium">{metrics.host_info ? `${metrics.host_info.platform} (${metrics.host_info.kernelArch})` : '-'}</span>
                </div>
                <div className="flex flex-col px-4 border-r border-white/5">
                    <span className="text-gray-500 uppercase text-xs font-semibold tracking-wider mb-1">Kernel</span>
                    <span className="text-white font-medium">{metrics.host_info?.kernelVersion || '-'}</span>
                </div>
                <div className="flex flex-col px-4 border-r border-white/5">
                    <span className="text-gray-500 uppercase text-xs font-semibold tracking-wider mb-1">Total RAM</span>
                    <span className="text-white font-medium font-mono">{metrics.memory ? formatBytes(metrics.memory.total) : '-'}</span>
                </div>
                <div className="flex flex-col px-4">
                    <span className="text-gray-500 uppercase text-xs font-semibold tracking-wider mb-1">Hostname</span>
                    <span className="text-primary font-medium">{metrics.host_info?.hostname || '-'}</span>
                </div>
            </div>

            {/* Main Stats Row (Circular) */}
            <div className="bg-surface/50 backdrop-blur-sm p-8 rounded-2xl border border-white/5 shadow-xl">
                <div className="grid grid-cols-2 md:grid-cols-4 gap-8 justify-items-center">
                    <CircularProgress
                        value={metrics.cpu_total || 0}
                        label="CPU Usage"
                        color="#3B82F6" // Primary
                    />
                    <CircularProgress
                        value={metrics.memory?.usedPercent || 0}
                        label="RAM Usage"
                        color="#8B5CF6" // Accent
                    />
                    <CircularProgress
                        value={metrics.swap?.usedPercent || 0}
                        label="Swap Usage"
                        color="#F59E0B" // Warning
                    />
                    <CircularProgress
                        value={metrics.disks?.[0]?.used_percent || 0}
                        label="I/O Usage"
                        color="#10B981" // Secondary
                    />
                </div>

                {/* Secondary Stats (Load, Network, Alerts) */}
                <div className="mt-8 pt-8 border-t border-white/5 grid grid-cols-1 md:grid-cols-3 gap-8">
                    {/* Load Average */}
                    <div className="bg-gray-900/50 p-6 rounded-xl border border-white/5">
                        <h3 className="text-gray-400 text-xs font-bold uppercase tracking-wider mb-6">Load Status</h3>
                        <div className="grid grid-cols-3 gap-4 text-center">
                            <div>
                                <div className="text-xs text-gray-500 uppercase mb-1">Load 1</div>
                                <div className="text-2xl font-bold text-white font-mono">{metrics.load_avg?.load1?.toFixed(2) || '-'}</div>
                            </div>
                            <div>
                                <div className="text-xs text-gray-500 uppercase mb-1">Load 5</div>
                                <div className="text-2xl font-bold text-white font-mono">{metrics.load_avg?.load5?.toFixed(2) || '-'}</div>
                            </div>
                            <div>
                                <div className="text-xs text-gray-500 uppercase mb-1">Load 15</div>
                                <div className="text-2xl font-bold text-white font-mono">{metrics.load_avg?.load15?.toFixed(2) || '-'}</div>
                            </div>
                        </div>
                    </div>

                    {/* Network Status */}
                    <div className="bg-gray-900/50 p-6 rounded-xl border border-white/5">
                        <h3 className="text-gray-400 text-xs font-bold uppercase tracking-wider mb-6">Network Status</h3>
                        <div className="grid grid-cols-1 gap-4 text-center">
                            <div>
                                <div className="text-xs text-gray-500 uppercase mb-1">Total Throughput</div>
                                <div className="text-2xl font-bold text-white font-mono">
                                    {metrics.network ? formatBytes(metrics.network.total_recv + metrics.network.total_sent).split(' ')[0] : '0'}
                                    <span className="text-sm text-gray-600 ml-1">
                                        {metrics.network ? formatBytes(metrics.network.total_recv + metrics.network.total_sent).split(' ')[1] : 'B'}/s
                                    </span>
                                </div>
                                <div className="flex justify-center gap-4 mt-2 text-xs text-gray-500">
                                    <span>↓ {metrics.network ? formatBytes(metrics.network.total_recv) : '0 B'}/s</span>
                                    <span>↑ {metrics.network ? formatBytes(metrics.network.total_sent) : '0 B'}/s</span>
                                </div>
                            </div>
                        </div>
                    </div>

                    {/* Alerts Card */}
                    <div className="bg-gray-900/50 p-6 rounded-xl border border-white/5 flex flex-col justify-between">
                        <h3 className="text-gray-400 text-xs font-bold uppercase tracking-wider mb-6">Alerts</h3>
                        <div>
                            <div className="text-xs text-gray-500 uppercase mb-1">Status</div>
                            <div className="flex items-center space-x-2">
                                <span className="text-lg font-bold text-blue-400">All Systems</span>
                                <CheckCircle className="w-5 h-5 text-blue-500" />
                            </div>
                        </div>
                        <button
                            onClick={() => onNavigate('services')}
                            className="mt-4 w-full py-2 bg-white/5 hover:bg-white/10 rounded-lg text-xs font-medium text-gray-300 transition-colors flex items-center justify-center gap-2"
                        >
                            View Details <ArrowRight size={14} />
                        </button>
                    </div>
                </div>
            </div>

            {/* Partitions (Circular) */}
            <div className="bg-surface/50 backdrop-blur-sm p-6 rounded-2xl border border-white/5 shadow-xl">
                <h3 className="text-gray-400 text-xs font-bold uppercase tracking-wider mb-6">Partitions</h3>
                <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-8">
                    {(!metrics.disks || metrics.disks.length === 0) && (
                        <div className="col-span-full text-center text-gray-400 py-4">No partitions found</div>
                    )}
                    {metrics.disks?.map((disk) => (
                        <div key={disk.path} className="flex items-center space-x-4 bg-gray-900/30 p-4 rounded-xl border border-white/5">
                            <div className="relative w-16 h-16 flex-shrink-0">
                                <svg className="w-full h-full transform -rotate-90">
                                    <circle
                                        cx="32"
                                        cy="32"
                                        r="28"
                                        stroke="#1e293b"
                                        strokeWidth="6"
                                        fill="transparent"
                                    />
                                    <circle
                                        cx="32"
                                        cy="32"
                                        r="28"
                                        stroke={disk.used_percent > 80 ? "#ef4444" : "#f59e0b"} // Orange/Red based on usage
                                        strokeWidth="6"
                                        fill="transparent"
                                        strokeDasharray={2 * Math.PI * 28}
                                        strokeDashoffset={2 * Math.PI * 28 * (1 - disk.used_percent / 100)}
                                        strokeLinecap="round"
                                        className="transition-all duration-1000 ease-out"
                                    />
                                </svg>
                                <div className="absolute inset-0 flex items-center justify-center">
                                    <span className="text-xs font-bold text-white">{Math.round(disk.used_percent)}%</span>
                                </div>
                            </div>
                            <div className="flex flex-col overflow-hidden">
                                <span className="text-white font-medium truncate" title={disk.path}>{disk.path}</span>
                                <span className="text-xs text-gray-500">{formatBytes(disk.used)} / {formatBytes(disk.total)}</span>
                            </div>
                        </div>
                    ))}
                </div>
            </div>
        </div>
    );
};
