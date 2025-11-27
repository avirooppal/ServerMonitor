import React from 'react';
import type { SystemMetrics } from '../../types';

interface SystemSectionProps {
    metrics: SystemMetrics;
}

export const SystemSection: React.FC<SystemSectionProps> = ({ metrics }) => {
    const formatUptime = (seconds: number) => {
        const d = Math.floor(seconds / (3600 * 24));
        const h = Math.floor((seconds % (3600 * 24)) / 3600);
        const m = Math.floor((seconds % 3600) / 60);
        return `${d}d ${h}h ${m}m`;
    };

    const info = metrics.host_info;

    return (
        <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
            <div className="bg-gray-900 p-6 rounded-xl border border-gray-800 shadow-sm">
                <h3 className="text-lg font-semibold mb-6 text-gray-200">System Information</h3>
                <div className="space-y-4">
                    <div className="flex justify-between border-b border-gray-800 pb-2">
                        <span className="text-gray-500">Hostname</span>
                        <span className="text-white font-mono">{info.hostname}</span>
                    </div>
                    <div className="flex justify-between border-b border-gray-800 pb-2">
                        <span className="text-gray-500">OS</span>
                        <span className="text-white">{info.os} {info.platform} {info.platformVersion}</span>
                    </div>
                    <div className="flex justify-between border-b border-gray-800 pb-2">
                        <span className="text-gray-500">Kernel</span>
                        <span className="text-white font-mono">{info.kernelVersion}</span>
                    </div>
                    <div className="flex justify-between border-b border-gray-800 pb-2">
                        <span className="text-gray-500">Uptime</span>
                        <span className="text-green-400 font-mono">{formatUptime(info.uptime)}</span>
                    </div>
                    <div className="flex justify-between border-b border-gray-800 pb-2">
                        <span className="text-gray-500">Architecture</span>
                        <span className="text-white font-mono">{info.kernelArch}</span>
                    </div>
                </div>
            </div>
        </div>
    );
};
