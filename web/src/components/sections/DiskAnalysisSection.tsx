import React, { useEffect, useState } from 'react';
import { HardDrive, Folder, Activity } from 'lucide-react';
import axios from 'axios';
import { type System } from '../../utils/api';

interface FolderSize {
    path: string;
    size: number;
}

interface DiskAnalysisSectionProps {
    system: System;
}

interface DiskHistory {
    timestamp: string;
    used_percent: number;
    total: number;
    used: number;
}

const DiskAnalysisSection: React.FC<DiskAnalysisSectionProps> = ({ system }) => {
    const [folders, setFolders] = useState<FolderSize[]>([]);
    const [history, setHistory] = useState<DiskHistory[]>([]);
    const [loading, setLoading] = useState(true);

    useEffect(() => {
        const fetchData = async () => {
            try {
                const config = {
                    headers: { 'Authorization': `Bearer ${system.api_key}` }
                };

                // Fetch Folders
                const folderRes = await axios.get(`${system.url}/api/v1/disk/usage`, {
                    ...config,
                    params: { path: '/' } // Default path
                });
                setFolders(folderRes.data || []);

                // Fetch History
                const historyRes = await axios.get(`${system.url}/api/v1/disk/history`, config);
                setHistory(historyRes.data || []);

            } catch (error) {
                console.error("Failed to fetch disk data", error);
            } finally {
                setLoading(false);
            }
        };

        fetchData();
    }, [system]);

    const formatBytes = (bytes: number) => {
        if (bytes === 0) return '0 B';
        const k = 1024;
        const sizes = ['B', 'KB', 'MB', 'GB', 'TB'];
        const i = Math.floor(Math.log(bytes) / Math.log(k));
        return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
    };

    if (loading) return <div className="p-4 text-gray-400">Loading disk analysis...</div>;

    return (
        <div className="bg-gray-800 rounded-lg p-6 border border-gray-700">
            <div className="flex items-center gap-2 mb-4">
                <HardDrive className="w-5 h-5 text-purple-400" />
                <h3 className="text-lg font-medium text-white">Top Folders by Size</h3>
            </div>

            <div className="space-y-3">
                {folders.map((folder, i) => (
                    <div key={i} className="bg-gray-900/50 p-3 rounded flex items-center justify-between">
                        <div className="flex items-center gap-3 overflow-hidden">
                            <Folder className="w-4 h-4 text-gray-500 flex-shrink-0" />
                            <span className="text-gray-300 text-sm truncate font-mono" title={folder.path}>
                                {folder.path}
                            </span>
                        </div>
                        <span className="text-purple-400 font-bold text-sm whitespace-nowrap ml-4">
                            {formatBytes(folder.size)}
                        </span>
                    </div>
                ))}
                {folders.length === 0 && (
                    <div className="text-center text-gray-500 py-8">
                        No disk usage data available yet. (Calculated every 15 mins)
                    </div>
                )}
            </div>

            {/* Disk History */}
            <div className="mt-6 pt-6 border-t border-gray-700">
                <div className="flex items-center gap-2 mb-4">
                    <Activity className="w-5 h-5 text-blue-400" />
                    <h3 className="text-lg font-medium text-white">Daily Disk Usage History</h3>
                </div>

                <div className="overflow-x-auto">
                    <table className="w-full text-left text-sm text-gray-400">
                        <thead className="bg-gray-900/50 text-gray-300 uppercase text-xs">
                            <tr>
                                <th className="px-4 py-3">Date</th>
                                <th className="px-4 py-3">Usage %</th>
                                <th className="px-4 py-3">Used</th>
                                <th className="px-4 py-3">Total</th>
                            </tr>
                        </thead>
                        <tbody className="divide-y divide-gray-700">
                            {history.map((h, i) => (
                                <tr key={i} className="hover:bg-gray-700/30">
                                    <td className="px-4 py-3 whitespace-nowrap">{new Date(h.timestamp).toLocaleDateString()} {new Date(h.timestamp).toLocaleTimeString()}</td>
                                    <td className="px-4 py-3">
                                        <div className="flex items-center gap-2">
                                            <div className="w-16 h-2 bg-gray-700 rounded-full overflow-hidden">
                                                <div
                                                    className="h-full bg-blue-500 rounded-full"
                                                    style={{ width: `${h.used_percent}%` }}
                                                />
                                            </div>
                                            <span className="text-white font-mono">{h.used_percent.toFixed(1)}%</span>
                                        </div>
                                    </td>
                                    <td className="px-4 py-3 font-mono">{formatBytes(h.used)}</td>
                                    <td className="px-4 py-3 font-mono">{formatBytes(h.total)}</td>
                                </tr>
                            ))}
                            {history.length === 0 && (
                                <tr>
                                    <td colSpan={4} className="px-4 py-8 text-center text-gray-500">
                                        No history available yet. (Snapshots taken daily)
                                    </td>
                                </tr>
                            )}
                        </tbody>
                    </table>
                </div>
            </div>
        </div>
    );
};

export default DiskAnalysisSection;
