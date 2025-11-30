import React, { useEffect, useState } from 'react';
import { type System, fetchMetrics } from '../../utils/api';
import { Search, Server } from 'lucide-react';
import clsx from 'clsx';

interface ServicesSectionProps {
    systems: System[];
}

interface SystemStatus {
    id: number;
    status: 'ok' | 'error' | 'loading';
    uptime: string;
    lastCheck: string;
    details: string;
}

export const ServicesSection: React.FC<ServicesSectionProps> = ({ systems }) => {
    const [statuses, setStatuses] = useState<Record<number, SystemStatus>>({});
    const [searchTerm, setSearchTerm] = useState('');

    useEffect(() => {
        const checkAllSystems = async () => {
            const newStatuses: Record<number, SystemStatus> = {};

            for (const system of systems) {
                try {
                    // We fetch metrics to check status. 
                    // In a real app, we might have a lightweight /ping endpoint.
                    const data = await fetchMetrics(system.id);

                    // Format uptime
                    const seconds = data.host_info.uptime;
                    const days = Math.floor(seconds / 86400);
                    const hours = Math.floor((seconds % 86400) / 3600);
                    const minutes = Math.floor((seconds % 3600) / 60);
                    let uptimeStr = "";
                    if (days > 0) uptimeStr = `${days} days, ${hours}h`;
                    else uptimeStr = `${hours}h ${minutes}m`;

                    newStatuses[system.id] = {
                        id: system.id,
                        status: 'ok',
                        uptime: uptimeStr,
                        lastCheck: new Date().toLocaleTimeString(),
                        details: `Online - ${data.host_info.platform} ${data.host_info.platformVersion}`
                    };
                } catch (e) {
                    newStatuses[system.id] = {
                        id: system.id,
                        status: 'error',
                        uptime: '-',
                        lastCheck: new Date().toLocaleTimeString(),
                        details: 'Connection failed'
                    };
                }
            }
            setStatuses(newStatuses);
        };

        checkAllSystems();
        const interval = setInterval(checkAllSystems, 30000); // Check every 30s
        return () => clearInterval(interval);
    }, [systems]);

    const filteredSystems = systems.filter(s =>
        s.name.toLowerCase().includes(searchTerm.toLowerCase()) ||
        s.url.toLowerCase().includes(searchTerm.toLowerCase())
    );

    return (
        <div className="space-y-6">
            {/* Header / Filter */}
            <div className="flex flex-col md:flex-row md:items-center justify-between gap-4">
                <div>
                    <h2 className="text-xl font-bold text-white flex items-center gap-2">
                        <Server className="w-5 h-5 text-primary" />
                        All Checks
                    </h2>
                    <p className="text-gray-400 text-sm">Overview of all your configured monitors</p>
                </div>

                <div className="flex items-center space-x-2 bg-surface/50 border border-white/10 rounded-lg px-3 py-2 w-full md:w-64">
                    <Search size={16} className="text-gray-500" />
                    <input
                        type="text"
                        placeholder="Search..."
                        className="bg-transparent border-none focus:outline-none text-sm text-white w-full placeholder-gray-500"
                        value={searchTerm}
                        onChange={(e) => setSearchTerm(e.target.value)}
                    />
                </div>
            </div>

            {/* List */}
            <div className="bg-surface/30 backdrop-blur-sm border border-white/5 rounded-xl overflow-hidden">
                {/* Table Header */}
                <div className="grid grid-cols-12 gap-4 px-6 py-3 bg-white/5 text-xs font-semibold text-gray-400 uppercase tracking-wider">
                    <div className="col-span-2 md:col-span-1">Status</div>
                    <div className="col-span-4 md:col-span-3">Details</div>
                    <div className="col-span-6 md:col-span-4 hidden md:block">Notifications</div>
                    <div className="col-span-6 md:col-span-4 text-right md:text-left">Uptime</div>
                </div>

                <div className="divide-y divide-white/5">
                    {filteredSystems.map(system => {
                        const status = statuses[system.id] || { status: 'loading', uptime: '...', lastCheck: '...', details: 'Checking...' };

                        return (
                            <div key={system.id} className="grid grid-cols-12 gap-4 px-6 py-4 items-center hover:bg-white/5 transition-colors">
                                {/* Status */}
                                <div className="col-span-2 md:col-span-1">
                                    {status.status === 'loading' ? (
                                        <div className="w-12 h-6 bg-gray-700 rounded animate-pulse" />
                                    ) : (
                                        <span className={clsx(
                                            "inline-flex items-center px-2.5 py-0.5 rounded text-xs font-bold uppercase tracking-wide",
                                            status.status === 'ok' ? "bg-green-500/20 text-green-400" : "bg-red-500/20 text-red-400"
                                        )}>
                                            {status.status === 'ok' ? 'OK' : 'DOWN'}
                                        </span>
                                    )}
                                </div>

                                {/* Details */}
                                <div className="col-span-4 md:col-span-3">
                                    <div className="font-medium text-white">{system.name}</div>
                                    <div className="text-xs text-gray-500 truncate">{status.details}</div>
                                    <div className="text-xs text-gray-600 mt-0.5">Last check at {status.lastCheck}</div>
                                </div>

                                {/* Notifications (Placeholder) */}
                                <div className="col-span-6 md:col-span-4 hidden md:block">
                                    <div className="flex items-center space-x-2 text-xs text-gray-400 bg-white/5 px-2 py-1 rounded w-fit">
                                        <span>admin@example.com</span>
                                    </div>
                                </div>

                                {/* Uptime */}
                                <div className="col-span-6 md:col-span-4 flex flex-col items-end md:items-start">
                                    <div className="flex items-center justify-between w-full mb-1">
                                        <span className="text-xs text-gray-400">Current Uptime</span>
                                        <span className="text-xs font-mono text-white">{status.uptime}</span>
                                    </div>
                                    {/* Uptime Bar Visualization */}
                                    <div className="flex space-x-0.5 w-full h-4">
                                        {[...Array(20)].map((_, i) => (
                                            <div
                                                key={i}
                                                className={clsx(
                                                    "flex-1 rounded-sm",
                                                    status.status === 'ok' ? "bg-green-500" : "bg-red-500",
                                                    status.status === 'loading' && "bg-gray-700 animate-pulse"
                                                )}
                                                style={{ opacity: status.status === 'ok' ? 0.6 + (i / 40) : 1 }} // Slight gradient effect
                                            />
                                        ))}
                                    </div>
                                    <div className="text-[10px] text-gray-500 mt-1 w-full text-right md:text-left">
                                        0 incidents
                                    </div>
                                </div>
                            </div>
                        );
                    })}

                    {filteredSystems.length === 0 && (
                        <div className="px-6 py-8 text-center text-gray-500">
                            No services found matching your search.
                        </div>
                    )}
                </div>
            </div>
        </div>
    );
};
