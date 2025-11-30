import React, { useEffect, useState } from 'react';
import { HardDrive, Folder } from 'lucide-react';

interface FolderSize {
    path: string;
    size: number;
}

interface DiskAnalysisSectionProps {
    systemId: number;
    apiKey: string;
}

const DiskAnalysisSection: React.FC<DiskAnalysisSectionProps> = ({ systemId, apiKey }) => {
    const [folders, setFolders] = useState<FolderSize[]>([]);
    const [loading, setLoading] = useState(true);

    useEffect(() => {
        const fetchData = async () => {
            try {
                const headers = { Authorization: `Bearer ${apiKey}` };
                // Fetch Disk Usage
                // Using the proxy endpoint again.
                const response = await fetch(`/api/v1/systems/${systemId}/proxy?path=/disk/usage`, { headers });
                if (response.ok) {
                    setFolders(await response.json());
                }
            } catch (error) {
                console.error("Failed to fetch disk usage", error);
            } finally {
                setLoading(false);
            }
        };

        fetchData();
    }, [systemId, apiKey]);

    const formatBytes = (bytes: number) => {
        if (bytes === 0) return '0 B';
        const k = 1024;
        const sizes = ['B', 'KB', 'MB', 'GB', 'TB'];
        const i = Math.floor(Math.log(bytes) / Math.log(k));
        return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
    };

    if (loading) return <div className="p-4 text-gray-400">Analyzing disk usage...</div>;

    return (
        <div className="bg-gray-800 rounded-lg p-6 border border-gray-700">
            <div className="flex items-center gap-2 mb-6">
                <HardDrive className="w-5 h-5 text-purple-400" />
                <h3 className="text-lg font-medium text-white">Disk Space Analysis (Top Folders)</h3>
            </div>

            <div className="space-y-3">
                {folders.map((folder, i) => (
                    <div key={i} className="flex items-center justify-between p-3 bg-gray-900/50 rounded-lg hover:bg-gray-700/30 transition-colors">
                        <div className="flex items-center gap-3">
                            <Folder className="w-4 h-4 text-yellow-500" />
                            <span className="text-gray-300 font-mono text-sm">/{folder.path}</span>
                        </div>
                        <div className="flex items-center gap-4">
                            {/* Simple bar visualization */}
                            <div className="hidden md:block w-32 h-2 bg-gray-700 rounded-full overflow-hidden">
                                <div
                                    className="h-full bg-purple-500 rounded-full"
                                    style={{ width: `${Math.min((folder.size / (folders[0]?.size || 1)) * 100, 100)}%` }}
                                />
                            </div>
                            <span className="text-white font-medium min-w-[80px] text-right">
                                {formatBytes(folder.size)}
                            </span>
                        </div>
                    </div>
                ))}
                {folders.length === 0 && (
                    <div className="text-center text-gray-500 py-8">No disk usage data available</div>
                )}
            </div>
        </div>
    );
};

export default DiskAnalysisSection;
