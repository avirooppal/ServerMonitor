import React from 'react';
import type { SystemMetrics } from '../../types';
import { useMetricHistory } from '../../hooks/useMetricHistory';
import { SimpleLineChart } from '../charts/SimpleLineChart';

interface NetworkSectionProps {
    metrics: SystemMetrics;
}

export const NetworkSection: React.FC<NetworkSectionProps> = ({ metrics }) => {
    const rxHistory = useMetricHistory(metrics.network.total_recv);
    const txHistory = useMetricHistory(metrics.network.total_sent);

    const formatBits = (bytes: number) => {
        const bits = bytes * 8;
        if (bits === 0) return '0 bps';
        const k = 1000;
        const sizes = ['bps', 'Kbps', 'Mbps', 'Gbps'];
        const i = Math.floor(Math.log(bits) / Math.log(k));
        return parseFloat((bits / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
    };

    return (
        <div className="space-y-6">
            <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                {/* RX Chart */}
                <div className="bg-gray-900 p-6 rounded-xl border border-gray-800 shadow-sm">
                    <h3 className="text-lg font-semibold mb-4 text-gray-200">Incoming Traffic</h3>
                    <div className="flex items-end space-x-2 mb-4">
                        <span className="text-4xl font-bold text-green-500">{formatBits(metrics.network.total_recv)}</span>
                    </div>
                    <SimpleLineChart
                        data={rxHistory.map(h => ({ time: h.time, value: h.value * 8 }))}
                        dataKey="value"
                        color="#22c55e"
                        yDomain={[0, 'auto']}
                    />
                </div>

                {/* TX Chart */}
                <div className="bg-gray-900 p-6 rounded-xl border border-gray-800 shadow-sm">
                    <h3 className="text-lg font-semibold mb-4 text-gray-200">Outgoing Traffic</h3>
                    <div className="flex items-end space-x-2 mb-4">
                        <span className="text-4xl font-bold text-blue-500">{formatBits(metrics.network.total_sent)}</span>
                    </div>
                    <SimpleLineChart
                        data={txHistory.map(h => ({ time: h.time, value: h.value * 8 }))}
                        dataKey="value"
                        color="#3b82f6"
                        yDomain={[0, 'auto']}
                    />
                </div>
            </div>

            {/* Interfaces Table */}
            <div className="bg-gray-900 rounded-xl border border-gray-800 shadow-sm overflow-hidden">
                <div className="px-6 py-4 border-b border-gray-800">
                    <h3 className="text-lg font-semibold text-gray-200">Interfaces</h3>
                </div>
                <div className="overflow-x-auto">
                    <table className="w-full text-left text-sm text-gray-400">
                        <thead className="bg-gray-800 text-gray-200 uppercase font-medium">
                            <tr>
                                <th className="px-6 py-3">Name</th>
                                <th className="px-6 py-3 text-right">Receive Rate</th>
                                <th className="px-6 py-3 text-right">Send Rate</th>
                            </tr>
                        </thead>
                        <tbody className="divide-y divide-gray-800">
                            {metrics.network.interfaces.map((iface) => (
                                <tr key={iface.name} className="hover:bg-gray-800/50">
                                    <td className="px-6 py-4 font-medium text-white">{iface.name}</td>
                                    <td className="px-6 py-4 text-right font-mono text-green-400">{formatBits(iface.recv_rate)}</td>
                                    <td className="px-6 py-4 text-right font-mono text-blue-400">{formatBits(iface.sent_rate)}</td>
                                </tr>
                            ))}
                        </tbody>
                    </table>
                </div>
            </div>
        </div>
    );
};
