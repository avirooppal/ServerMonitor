import React from 'react';
import type { SystemMetrics } from '../../types';
import { CircularProgress } from '../CircularProgress';
import { ArrowDown, ArrowUp, HardDrive } from 'lucide-react';

interface OverviewSectionProps {
    metrics: SystemMetrics;
}

export const OverviewSection: React.FC<OverviewSectionProps> = ({ metrics }) => {
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
                    <span className="text-white font-medium font-mono">{formatUptime(metrics.host_info.uptime)}</span>
                </div>
                <div className="flex flex-col px-4 border-r border-white/5">
                    <span className="text-gray-500 uppercase text-xs font-semibold tracking-wider mb-1">OS / Arch</span>
                    <span className="text-white font-medium">{metrics.host_info.platform} ({metrics.host_info.kernelArch})</span>
                </div>
                <div className="flex flex-col px-4 border-r border-white/5">
                    <span className="text-gray-500 uppercase text-xs font-semibold tracking-wider mb-1">Kernel</span>
                    <span className="text-white font-medium">{metrics.host_info.kernelVersion}</span>
                </div>
                <div className="flex flex-col px-4 border-r border-white/5">
                    <span className="text-gray-500 uppercase text-xs font-semibold tracking-wider mb-1">Total RAM</span>
                    <span className="text-white font-medium font-mono">{formatBytes(metrics.memory.total)}</span>
                </div>
                <div className="flex flex-col px-4">
                    <span className="text-gray-500 uppercase text-xs font-semibold tracking-wider mb-1">Hostname</span>
                    <span className="text-primary font-medium">{metrics.host_info.hostname}</span>
                </div>
            </div>

            {/* Main Stats Row (Circular) */}
            <div className="bg-surface/50 backdrop-blur-sm p-8 rounded-2xl border border-white/5 shadow-xl">
                <div className="grid grid-cols-2 md:grid-cols-4 gap-8 justify-items-center">
                    <CircularProgress
                        value={metrics.cpu_total}
                        label="CPU Usage"
                        color="#3B82F6" // Primary
                    />
                    <CircularProgress
                        value={metrics.memory.usedPercent}
                        label="RAM Usage"
                        color="#8B5CF6" // Accent
                    />
                    <CircularProgress
                        value={metrics.swap.usedPercent}
                        label="Swap Usage"
                        color="#F59E0B" // Warning
                    />
                    <CircularProgress
                        value={metrics.disks[0]?.used_percent || 0}
                        label="Disk Usage"
                        color="#10B981" // Secondary
                    />
                </div>

                {/* Secondary Stats (Load, Network) */}
                <div className="mt-8 pt-8 border-t border-white/5 grid grid-cols-1 md:grid-cols-2 gap-8">
                    {/* Load Average */}
                    <div className="flex flex-col items-center">
                        <h3 className="text-gray-400 text-sm font-medium uppercase tracking-wider mb-4">Load Average</h3>
                        <div className="flex space-x-6">
                            <div className="flex flex-col items-center">
                                <span className="text-2xl font-bold text-white font-mono">{metrics.load_avg.load1.toFixed(2)}</span>
                                <span className="text-xs text-gray-500 mt-1">1 min</span>
                            </div>
                            <div className="flex flex-col items-center">
                                <span className="text-2xl font-bold text-white font-mono">{metrics.load_avg.load5.toFixed(2)}</span>
                                <span className="text-xs text-gray-500 mt-1">5 min</span>
                            </div>
                            <div className="flex flex-col items-center">
                                <span className="text-2xl font-bold text-white font-mono">{metrics.load_avg.load15.toFixed(2)}</span>
                                <span className="text-xs text-gray-500 mt-1">15 min</span>
                            </div>
                        </div>
                    </div>

                    {/* Network Status */}
                    <div className="flex flex-col items-center">
                        <h3 className="text-gray-400 text-sm font-medium uppercase tracking-wider mb-4">Network Status</h3>
                        <div className="flex space-x-8">
                            <div className="flex items-center space-x-3">
                                <div className="p-2 bg-secondary/10 rounded-lg">
                                    <ArrowDown className="text-secondary w-5 h-5" />
                                </div>
                                <div>
                                    <div className="text-xl font-bold text-white font-mono">{formatBytes(metrics.network.total_recv)}/s</div>
                                    <div className="text-xs text-gray-500">Download</div>
                                </div>
                            </div>
                            <div className="flex items-center space-x-3">
                                <div className="p-2 bg-primary/10 rounded-lg">
                                    <ArrowUp className="text-primary w-5 h-5" />
                                </div>
                                <div>
                                    <div className="text-xl font-bold text-white font-mono">{formatBytes(metrics.network.total_sent)}/s</div>
                                    <div className="text-xs text-gray-500">Upload</div>
                                </div>
                            </div>
                        </div>
                    </div>
                </div>
            </div>

            {/* Disk Usage Details */}
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
                {metrics.disks.map((disk) => (
                    <div key={disk.path} className="bg-surface/50 backdrop-blur-sm p-6 rounded-2xl border border-white/5 shadow-xl">
                        <div className="flex items-center justify-between mb-4">
                            <div className="flex items-center space-x-3">
                                <div className="p-2 bg-gray-800 rounded-lg">
                                    <HardDrive className="text-gray-400 w-5 h-5" />
                                </div>
                                <div>
                                    <h3 className="text-white font-medium">{disk.path}</h3>
                                    <p className="text-xs text-gray-500 font-mono">{formatBytes(disk.total)} Total</p>
                                </div>
                            </div>
                            <span className="text-2xl font-bold text-white font-mono">{Math.round(disk.used_percent)}%</span>
                        </div>

                        <div className="w-full bg-gray-800 rounded-full h-2 mb-4 overflow-hidden">
                            <div
                                className="bg-gradient-to-r from-primary to-accent h-2 rounded-full transition-all duration-500"
                                style={{ width: `${disk.used_percent}%` }}
                            />
                        </div>

                        <div className="flex justify-between text-sm">
                            <div className="flex flex-col">
                                <span className="text-gray-500 text-xs">Used</span>
                                <span className="text-white font-mono">{formatBytes(disk.used)}</span>
                            </div>
                            <div className="flex flex-col items-end">
                                <span className="text-gray-500 text-xs">Free</span>
                                <span className="text-white font-mono">{formatBytes(disk.free)}</span>
                            </div>
                        </div>
                    </div>
                ))}
            </div>
        </div>
    );
};
