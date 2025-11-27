import React from 'react';
import type { SystemMetrics } from '../../types';
import { useMetricHistory } from '../../hooks/useMetricHistory';
import { SimpleLineChart } from '../charts/SimpleLineChart';
import clsx from 'clsx';

interface CpuSectionProps {
    metrics: SystemMetrics;
}

export const CpuSection: React.FC<CpuSectionProps> = ({ metrics }) => {
    const history = useMetricHistory(metrics.cpu_total);
    const chartData = history.map(h => ({ time: h.time, value: h.value }));

    return (
        <div className="space-y-6">
            <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
                {/* Main CPU Chart */}
                <div className="md:col-span-2 bg-gray-900 p-6 rounded-xl border border-gray-800 shadow-sm">
                    <h3 className="text-lg font-semibold mb-4 text-gray-200">Total CPU Usage</h3>
                    <div className="flex items-end space-x-2 mb-4">
                        <span className="text-4xl font-bold text-blue-500">{metrics.cpu_total.toFixed(1)}%</span>
                        <span className="text-gray-500 mb-1">current</span>
                    </div>
                    <SimpleLineChart data={chartData} dataKey="value" color="#3b82f6" />
                </div>

                {/* Load Average */}
                <div className="bg-gray-900 p-6 rounded-xl border border-gray-800 shadow-sm">
                    <h3 className="text-lg font-semibold mb-4 text-gray-200">Load Average</h3>
                    <div className="space-y-4">
                        <div>
                            <div className="text-sm text-gray-400 mb-1">1 Min</div>
                            <div className="text-2xl font-mono text-white">{metrics.load_avg.load1.toFixed(2)}</div>
                        </div>
                        <div>
                            <div className="text-sm text-gray-400 mb-1">5 Min</div>
                            <div className="text-2xl font-mono text-white">{metrics.load_avg.load5.toFixed(2)}</div>
                        </div>
                        <div>
                            <div className="text-sm text-gray-400 mb-1">15 Min</div>
                            <div className="text-2xl font-mono text-white">{metrics.load_avg.load15.toFixed(2)}</div>
                        </div>
                    </div>
                </div>
            </div>

            {/* Per Core Usage */}
            <div className="bg-gray-900 p-6 rounded-xl border border-gray-800 shadow-sm">
                <h3 className="text-lg font-semibold mb-4 text-gray-200">Per Core Usage</h3>
                <div className="grid grid-cols-2 sm:grid-cols-4 md:grid-cols-6 lg:grid-cols-8 gap-4">
                    {metrics.cpu.map((usage, index) => (
                        <div key={index} className="bg-gray-800 p-3 rounded border border-gray-700">
                            <div className="text-xs text-gray-400 mb-2">Core {index}</div>
                            <div className="h-20 flex items-end bg-gray-700/50 rounded overflow-hidden relative">
                                <div
                                    className={clsx(
                                        "w-full transition-all duration-500",
                                        usage > 80 ? "bg-red-500" : usage > 50 ? "bg-yellow-500" : "bg-green-500"
                                    )}
                                    style={{ height: `${usage}%` }}
                                />
                            </div>
                            <div className="text-right text-sm font-mono mt-1">{usage.toFixed(0)}%</div>
                        </div>
                    ))}
                </div>
            </div>
        </div>
    );
};
