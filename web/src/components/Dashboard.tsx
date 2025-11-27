import React, { useEffect, useState } from 'react';
import { fetchMetrics, fetchServers, type ServerSummary } from '../utils/api';
import type { SystemMetrics } from '../types';
import { Activity, Cpu, HardDrive, Server, Layers, Settings, LogOut, LayoutDashboard, Box } from 'lucide-react';
import clsx from 'clsx';
import { OverviewSection } from './sections/OverviewSection';
import { CpuSection } from './sections/CpuSection';
import { MemorySection } from './sections/MemorySection';
import { DiskSection } from './sections/DiskSection';
import { NetworkSection } from './sections/NetworkSection';
import { ProcessSection } from './sections/ProcessSection';
import { DockerSection } from './sections/DockerSection';
import { SettingsSection } from './sections/SettingsSection';
import { SystemSection } from './sections/SystemSection';

interface DashboardProps {
    onLogout: () => void;
}

const TABS = [
    { id: 'overview', label: 'Overview', icon: LayoutDashboard },
    { id: 'cpu', label: 'CPU', icon: Cpu },
    { id: 'memory', label: 'RAM', icon: Layers },
    { id: 'network', label: 'Network', icon: Activity },
    { id: 'disk', label: 'Disks', icon: HardDrive },
    { id: 'processes', label: 'Processes', icon: Server },
    { id: 'docker', label: 'Docker', icon: Box },
    { id: 'settings', label: 'Settings', icon: Settings },
];

export const Dashboard: React.FC<DashboardProps> = ({ onLogout }) => {
    const [metrics, setMetrics] = useState<SystemMetrics | null>(null);
    const [servers, setServers] = useState<ServerSummary[]>([]);
    const [selectedServer, setSelectedServer] = useState<string>('local');
    const [activeTab, setActiveTab] = useState('overview');
    const [refreshRate, setRefreshRate] = useState(2000);
    const [error, setError] = useState('');

    useEffect(() => {
        const loadServers = async () => {
            try {
                const list = await fetchServers();
                setServers(list);
            } catch (e) {
                console.error("Failed to fetch servers", e);
            }
        };
        loadServers();
        const interval = setInterval(loadServers, 5000);
        return () => clearInterval(interval);
    }, []);

    useEffect(() => {
        const load = async () => {
            try {
                const data = await fetchMetrics(selectedServer);
                setMetrics(data);
                setError('');
            } catch (err) {
                setError('Connection lost');
            }
        };

        load();
        const interval = setInterval(load, refreshRate);
        return () => clearInterval(interval);
    }, [refreshRate, selectedServer]);

    return (
        <div className="flex flex-col h-screen bg-background text-gray-100 overflow-hidden font-sans selection:bg-primary/30">
            {/* Top Navigation Bar */}
            <header className="bg-surface/80 backdrop-blur-md border-b border-white/5 z-20">
                <div className="px-6 py-3 flex items-center justify-between">
                    <div className="flex items-center space-x-3">
                        <div className="bg-primary/10 p-2 rounded-lg border border-primary/20 shadow-glow-blue">
                            <Activity className="text-primary w-5 h-5" />
                        </div>
                        <h1 className="text-lg font-bold tracking-wide text-white font-sans">SERVER MONI</h1>
                    </div>

                    <div className="flex items-center space-x-4">
                        {/* Server Selector */}
                        <div className="flex items-center space-x-2 bg-background/50 rounded-lg px-3 py-1.5 border border-white/5">
                            <span className="text-xs text-gray-500 uppercase font-semibold tracking-wider">Server</span>
                            <select
                                className="bg-transparent text-sm focus:outline-none text-accent font-medium cursor-pointer"
                                value={selectedServer}
                                onChange={(e) => setSelectedServer(e.target.value)}
                            >
                                <option value="local">Localhost</option>
                                {servers.map(s => (
                                    <option key={s.id} value={s.id}>
                                        {s.hostname} ({s.id})
                                    </option>
                                ))}
                            </select>
                        </div>

                        <div className="flex items-center space-x-2 bg-background/50 rounded-lg px-3 py-1.5 border border-white/5">
                            <span className="text-xs text-gray-500 uppercase font-semibold tracking-wider">Refresh</span>
                            <select
                                className="bg-transparent text-sm focus:outline-none text-primary font-medium cursor-pointer"
                                value={refreshRate}
                                onChange={(e) => setRefreshRate(Number(e.target.value))}
                            >
                                <option value={1000}>1s</option>
                                <option value={2000}>2s</option>
                                <option value={5000}>5s</option>
                                <option value={10000}>10s</option>
                            </select>
                        </div>
                        <button onClick={onLogout} className="text-gray-400 hover:text-white transition-colors p-2 hover:bg-white/5 rounded-lg">
                            <LogOut size={18} />
                        </button>
                    </div>
                </div>

                {/* Tabs */}
                <div className="px-6 flex items-center space-x-2 overflow-x-auto no-scrollbar">
                    {TABS.map((tab) => {
                        const isActive = activeTab === tab.id;
                        return (
                            <button
                                key={tab.id}
                                onClick={() => setActiveTab(tab.id)}
                                className={clsx(
                                    "px-4 py-3 text-sm font-medium border-b-2 transition-all duration-200 whitespace-nowrap flex items-center space-x-2",
                                    isActive
                                        ? "border-primary text-primary bg-primary/5"
                                        : "border-transparent text-gray-400 hover:text-gray-200 hover:bg-white/5"
                                )}
                            >
                                <tab.icon size={16} />
                                <span>{tab.label}</span>
                            </button>
                        );
                    })}
                </div>
            </header>

            {/* Main Content */}
            <main className="flex-1 overflow-y-auto p-6 bg-background relative">
                {error && (
                    <div className="absolute top-0 left-0 w-full bg-danger/90 backdrop-blur text-white text-center text-sm py-1 z-50 font-medium">
                        {error} - Retrying...
                    </div>
                )}

                {activeTab === 'settings' ? (
                    <SettingsSection />
                ) : metrics ? (
                    <div className="max-w-7xl mx-auto space-y-6 pb-10">
                        {activeTab === 'overview' && <OverviewSection metrics={metrics} />}
                        {activeTab === 'cpu' && <CpuSection metrics={metrics} />}
                        {activeTab === 'memory' && <MemorySection metrics={metrics} />}
                        {activeTab === 'disk' && <DiskSection metrics={metrics} />}
                        {activeTab === 'network' && <NetworkSection metrics={metrics} />}
                        {activeTab === 'processes' && <ProcessSection metrics={metrics} />}
                        {activeTab === 'docker' && <DockerSection metrics={metrics} />}
                        {activeTab === 'system' && <SystemSection metrics={metrics} />}
                    </div>
                ) : (
                    <div className="flex flex-col items-center justify-center h-full text-gray-500 space-y-4">
                        <div className="w-12 h-12 border-4 border-primary/30 border-t-primary rounded-full animate-spin" />
                        <p className="font-medium animate-pulse">Connecting to server...</p>
                    </div>
                )}
            </main>
        </div>
    );
};
