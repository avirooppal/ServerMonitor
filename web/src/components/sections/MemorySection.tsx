import React from 'react';
import type { SystemMetrics } from '../../types';
import { useMetricHistory } from '../../hooks/useMetricHistory';
import { SimpleLineChart } from '../charts/SimpleLineChart';

interface MemorySectionProps {
    metrics: SystemMetrics;
}

export const MemorySection: React.FC<MemorySectionProps> = ({ metrics }) => {
    const ramHistory = useMetricHistory(metrics.memory.usedPercent);
    const swapHistory = useMetricHistory(metrics.swap.usedPercent);

    const formatBytes = (bytes: number) => {
        if (bytes === 0) return '0 B';
        const k = 1024;
        const sizes = ['B', 'KB', 'MB', 'GB', 'TB'];
        const i = Math.floor(Math.log(bytes) / Math.log(k));
        return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
    };

    return (
        <div className="space-y-6">
            <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                {/* RAM Chart */}
                <div className="bg-gray-900 p-6 rounded-xl border border-gray-800 shadow-sm">
                    <h3 className="text-lg font-semibold mb-4 text-gray-200">RAM Usage</h3>
                    <div className="flex justify-between items-end mb-4">
                        <div>
                            <span className="text-4xl font-bold text-purple-500">{metrics.memory.usedPercent.toFixed(1)}%</span>
                            <div className="text-sm text-gray-400 mt-1">
                                {formatBytes(metrics.memory.used)} / {formatBytes(metrics.memory.total)}
                            </div>
                        </div>
                    </div>
                    <div className="grid grid-cols-2 gap-4 mb-4 text-sm">
                        <div className="bg-gray-800/50 p-2 rounded border border-gray-700/50">
                            <span className="text-gray-500 block text-xs uppercase">Buffers</span>
                            <span className="text-white font-mono">{formatBytes(metrics.memory.buffers)}</span>
                        </div>
                        <div className="bg-gray-800/50 p-2 rounded border border-gray-700/50">
                            <span className="text-gray-500 block text-xs uppercase">Cached</span>
                            <span className="text-white font-mono">{formatBytes(metrics.memory.cached)}</span>
                        </div>
                    </div>
                    <SimpleLineChart data={ramHistory.map(h => ({ time: h.time, value: h.value }))} dataKey="value" color="#a855f7" />
                </div>

                {/* Swap Chart */}
                <div className="bg-gray-900 p-6 rounded-xl border border-gray-800 shadow-sm">
                    <h3 className="text-lg font-semibold mb-4 text-gray-200">Swap Usage</h3>
                    <div className="flex justify-between items-end mb-4">
                        <div>
                            <span className="text-4xl font-bold text-orange-500">{metrics.swap.usedPercent.toFixed(1)}%</span>
                            <div className="text-sm text-gray-400 mt-1">
                                {formatBytes(metrics.swap.used)} / {formatBytes(metrics.swap.total)}
                            </div>
                        </div>
                    </div>
                    <SimpleLineChart data={swapHistory.map(h => ({ time: h.time, value: h.value }))} dataKey="value" color="#f97316" />
                </div>
            </div>
        </div>
    );
};
