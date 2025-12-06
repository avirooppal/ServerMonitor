import React, { useEffect, useState } from 'react';
import { fetchMetrics, fetchSystems, addSystem, type System } from '../utils/api';
import type { SystemMetrics } from '../types';
import { Activity, Cpu, HardDrive, Server, Layers, Settings, LogOut, LayoutDashboard, Box, Shield } from 'lucide-react';
import clsx from 'clsx';
import { OverviewSection } from './sections/OverviewSection';
import { CpuSection } from './sections/CpuSection';
import { MemorySection } from './sections/MemorySection';
import { DiskSection } from './sections/DiskSection';
import { NetworkSection } from './sections/NetworkSection';
import { ProcessSection } from './sections/ProcessSection';
import { DockerSection } from './sections/DockerSection';
import { SettingsSection } from './sections/SettingsSection';
import { ServicesSection } from './sections/ServicesSection';
import SecuritySection from './sections/SecuritySection';
import DiskAnalysisSection from './sections/DiskAnalysisSection';

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
    { id: 'docker', label: 'Containers', icon: Box },
    { id: 'security', label: 'Security', icon: Shield },
    { id: 'disk-analysis', label: 'Disk Usage', icon: HardDrive },
    { id: 'settings', label: 'Settings', icon: Settings },
];

export const Dashboard: React.FC<DashboardProps> = ({ onLogout }) => {
    const [metrics, setMetrics] = useState<SystemMetrics | null>(null);
    const [systems, setSystems] = useState<System[]>([]);
    const [selectedSystemId, setSelectedSystemId] = useState<string>('');
    const [activeTab, setActiveTab] = useState('overview');
    const [refreshRate, setRefreshRate] = useState(2000);
    const [error, setError] = useState('');

    const [showAddModal, setShowAddModal] = useState(false);
    const [newSystemName, setNewSystemName] = useState('');
    const [newSystemCommand, setNewSystemCommand] = useState('');

    useEffect(() => {
        const loadSystems = async () => {
            try {
                const list = await fetchSystems();
                setSystems(list);

                // If list is empty, clear selection
                if (list.length === 0) {
                    setSelectedSystemId('');
                    return;
                }

                // If no selection or invalid selection, select first
                if (!selectedSystemId || !list.find(s => s.id.toString() === selectedSystemId)) {
                    setSelectedSystemId(list[0].id.toString());
                }
            } catch (e) {
                console.error("Failed to fetch systems", e);
            }
        };
        loadSystems();
        // Refresh systems list occasionally in case added from another tab
        const interval = setInterval(loadSystems, 10000);
        return () => clearInterval(interval);
    }, [selectedSystemId]);

    useEffect(() => {
        if (!selectedSystemId && activeTab !== 'settings') {
            return;
        }

        const load = async () => {
            if (!selectedSystemId) return;
            try {
                const data = await fetchMetrics(selectedSystemId);
                setMetrics(data);
                setError('');
            } catch (err) {
                setError('Connection lost');
                setMetrics(null);
            }
        };

        load();
        const interval = setInterval(load, refreshRate);
        return () => clearInterval(interval);
    }, [refreshRate, selectedSystemId, activeTab]);

    const [serverUrl, setServerUrl] = useState(() => localStorage.getItem('server_url') || 'http://107.150.20.37:8080');

    const handleAddSystem = async (e: React.FormEvent) => {
        e.preventDefault();
        try {
            // Save Server URL
            localStorage.setItem('server_url', serverUrl);

            // Generate random API Key
            const apiKey = Math.random().toString(36).substring(2) + Math.random().toString(36).substring(2);
            await addSystem(newSystemName, 'push', apiKey);

            // Generate Command
            // Remove trailing slash if present
            const cleanUrl = serverUrl.replace(/\/$/, '');
            const linuxCmd = `wget -q -O - https://raw.githubusercontent.com/avirooppal/ServerMonitor/main/downloads/install_agent_linux.sh | sudo bash -s ${apiKey} ${cleanUrl}`;
            const winCmd = `iwr https://raw.githubusercontent.com/avirooppal/ServerMonitor/main/scripts/install_agent_windows.ps1 -OutFile install.ps1; .\\install.ps1 -ServerUrl ${cleanUrl} -Token ${apiKey}`;

            setNewSystemCommand(`Linux:\n${linuxCmd}\n\nWindows (PowerShell):\n${winCmd}`);
            setNewSystemName('');
            // Refresh list
            const list = await fetchSystems();
            setSystems(list);
        } catch (err) {
            console.error(err);
            setError('Failed to add system');
        }
    };

    return (
        <div className="flex flex-col h-screen bg-background text-gray-100 overflow-hidden font-sans selection:bg-primary/30">
            {/* Top Navigation Bar */}
            <header className="bg-surface/80 backdrop-blur-md border-b border-white/5 z-20">
                <div className="px-6 py-3 flex items-center justify-between">
                    <div className="flex items-center space-x-4">
                        {/* Server Selector */}
                        <div className="flex items-center space-x-2 bg-background/50 rounded-lg px-3 py-1.5 border border-white/5">
                            <span className="text-xs text-gray-500 uppercase font-semibold tracking-wider">System</span>
                            <select
                                className="bg-transparent text-sm focus:outline-none text-accent font-medium cursor-pointer"
                                value={selectedSystemId}
                                onChange={(e) => setSelectedSystemId(e.target.value)}
                            >
                                {systems.length === 0 && <option value="">No Systems</option>}
                                {systems.map(s => (
                                    <option key={s.id} value={s.id}>
                                        {s.name}
                                    </option>
                                ))}
                            </select>
                        </div>

                        <button
                            onClick={() => setShowAddModal(true)}
                            className="bg-primary hover:bg-primary-hover text-white text-xs font-bold py-1.5 px-3 rounded transition-colors"
                        >
                            + Add Server
                        </button>

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
                    </div>

                    <div className="flex items-center space-x-4">
                        <button onClick={onLogout} className="text-gray-400 hover:text-white transition-colors p-2 hover:bg-white/5 rounded-lg">
                            <LogOut size={18} />
                        </button>
                    </div>
                </div>

                {/* Tabs */}
                <div className="px-6 pb-3 flex items-center space-x-2 overflow-x-auto no-scrollbar">
                    {TABS.map((tab) => {
                        const isActive = activeTab === tab.id;
                        return (
                            <button
                                key={tab.id}
                                onClick={() => setActiveTab(tab.id)}
                                className={clsx(
                                    "px-4 py-2 text-sm font-medium rounded-full transition-all duration-200 whitespace-nowrap flex items-center space-x-2",
                                    isActive
                                        ? "bg-primary text-white shadow-lg shadow-primary/25"
                                        : "bg-surface hover:bg-white/5 text-gray-400 hover:text-gray-200 border border-white/5"
                                )}
                            >
                                <tab.icon size={14} />
                                <span>{tab.label}</span>
                            </button>
                        );
                    })}
                </div>
            </header>

            {/* Main Content */}
            <main className="flex-1 overflow-y-auto p-6 bg-background relative">
                {error && activeTab !== 'settings' && (
                    <div className="absolute top-0 left-0 w-full bg-danger/90 backdrop-blur text-white text-center text-sm py-1 z-50 font-medium">
                        {error} - Retrying...
                    </div>
                )}

                {/* Add Server Modal */}
                {showAddModal && (
                    <div className="fixed inset-0 bg-black/80 backdrop-blur-sm z-50 flex items-center justify-center p-4">
                        <div className="bg-surface border border-white/10 rounded-xl max-w-lg w-full p-6 shadow-2xl">
                            <h2 className="text-xl font-bold text-white mb-4">Add New Server</h2>
                            {!newSystemCommand ? (
                                <form onSubmit={handleAddSystem} className="space-y-4">
                                    <div>
                                        <label className="block text-sm font-medium text-gray-400 mb-1">Server Name</label>
                                        <input
                                            type="text"
                                            value={newSystemName}
                                            onChange={(e) => setNewSystemName(e.target.value)}
                                            className="w-full bg-background border border-white/10 rounded px-3 py-2 text-white focus:outline-none focus:border-primary"
                                            placeholder="e.g. Production DB"
                                            required
                                        />
                                    </div>
                                    <div>
                                        <label className="block text-sm font-medium text-gray-400 mb-1">Backend URL (VPS)</label>
                                        <input
                                            type="url"
                                            value={serverUrl}
                                            onChange={(e) => setServerUrl(e.target.value)}
                                            className="w-full bg-background border border-white/10 rounded px-3 py-2 text-white focus:outline-none focus:border-primary"
                                            placeholder="http://your-vps-ip:8080"
                                            required
                                        />
                                        <p className="text-xs text-gray-500 mt-1">The URL where your backend is running.</p>
                                    </div>
                                    <div className="flex justify-end space-x-3 pt-4">
                                        <button
                                            type="button"
                                            onClick={() => setShowAddModal(false)}
                                            className="px-4 py-2 text-gray-400 hover:text-white transition-colors"
                                        >
                                            Cancel
                                        </button>
                                        <button
                                            type="submit"
                                            className="px-4 py-2 bg-primary hover:bg-primary-hover text-white rounded font-medium transition-colors"
                                        >
                                            Generate Command
                                        </button>
                                    </div>
                                </form>
                            ) : (
                                <div className="space-y-4">
                                    <p className="text-green-400 font-medium">Server Added! Run this command on your server:</p>
                                    <div className="bg-black/50 p-4 rounded border border-white/5 font-mono text-xs text-gray-300 overflow-x-auto whitespace-pre-wrap">
                                        {newSystemCommand}
                                    </div>
                                    <div className="flex justify-end pt-4">
                                        <button
                                            onClick={() => {
                                                setShowAddModal(false);
                                                setNewSystemCommand('');
                                            }}
                                            className="px-4 py-2 bg-primary hover:bg-primary-hover text-white rounded font-medium transition-colors"
                                        >
                                            Done
                                        </button>
                                    </div>
                                </div>
                            )}
                        </div>
                    </div>
                )}

                {activeTab === 'settings' ? (
                    <SettingsSection />
                ) : activeTab === 'services' ? (
                    <div className="max-w-7xl mx-auto pb-10">
                        <ServicesSection systems={systems} />
                    </div>
                ) : !selectedSystemId ? (
                    <div className="flex flex-col items-center justify-center h-full text-gray-500 space-y-4">
                        <Server size={48} className="text-gray-600" />
                        <p className="font-medium">No systems configured.</p>
                        <button
                            onClick={() => setShowAddModal(true)}
                            className="text-primary hover:underline"
                        >
                            Add your first server
                        </button>
                    </div>
                ) : metrics ? (
                    <div className="max-w-7xl mx-auto space-y-6 pb-10">
                        {activeTab === 'overview' && <OverviewSection metrics={metrics} onNavigate={setActiveTab} />}
                        {activeTab === 'cpu' && <CpuSection metrics={metrics} />}
                        {activeTab === 'memory' && <MemorySection metrics={metrics} />}
                        {activeTab === 'disk' && <DiskSection metrics={metrics} />}
                        {activeTab === 'network' && <NetworkSection metrics={metrics} />}
                        {activeTab === 'processes' && <ProcessSection metrics={metrics} />}
                        {activeTab === 'docker' && selectedSystemId && systems.find(s => s.id.toString() === selectedSystemId) && (
                            <DockerSection
                                metrics={metrics}
                                systemId={Number(selectedSystemId)}
                            />
                        )}
                        {activeTab === 'security' && selectedSystemId && systems.find(s => s.id.toString() === selectedSystemId) && (
                            <SecuritySection
                                systemId={Number(selectedSystemId)}
                            />
                        )}
                        {activeTab === 'disk-analysis' && selectedSystemId && systems.find(s => s.id.toString() === selectedSystemId) && (
                            <DiskAnalysisSection
                                systemId={Number(selectedSystemId)}
                            />
                        )}
                    </div>
                ) : (
                    <div className="flex flex-col items-center justify-center h-full text-gray-500 space-y-4">
                        <div className="w-12 h-12 border-4 border-primary/30 border-t-primary rounded-full animate-spin" />
                        <p className="font-medium animate-pulse">Connecting to system...</p>
                    </div>
                )}
            </main>
        </div>
    );
};
