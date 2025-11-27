import React from 'react';
import type { SystemMetrics } from '../../types';
import clsx from 'clsx';

interface DiskSectionProps {
    metrics: SystemMetrics;
}

export const DiskSection: React.FC<DiskSectionProps> = ({ metrics }) => {
    const formatBytes = (bytes: number) => {
        if (bytes === 0) return '0 B';
        const k = 1024;
        const sizes = ['B', 'KB', 'MB', 'GB', 'TB'];
        const i = Math.floor(Math.log(bytes) / Math.log(k));
        return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
    };

    return (
        <div className="space-y-6">
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
                {metrics.disks.map((disk, index) => (
                    <div key={index} className="bg-gray-900 p-6 rounded-xl border border-gray-800 shadow-sm">
                        <h3 className="text-lg font-semibold mb-2 text-gray-200 truncate" title={disk.path}>{disk.path}</h3>
                        <div className="text-sm text-gray-500 mb-4">Total: {formatBytes(disk.total)}</div>

                        <div className="relative pt-1">
                            <div className="flex mb-2 items-center justify-between">
                                <div>
                                    <span className="text-xs font-semibold inline-block py-1 px-2 uppercase rounded-full text-blue-200 bg-blue-900">
                                        Used
                                    </span>
                                </div>
                                <div className="text-right">
                                    <span className="text-xs font-semibold inline-block text-blue-200">
                                        {disk.used_percent.toFixed(1)}%
                                    </span>
                                </div>
                            </div>
                            <div className="overflow-hidden h-2 mb-4 text-xs flex rounded bg-gray-700">
                                <div
                                    style={{ width: `${disk.used_percent}%` }}
                                    className={clsx(
                                        "shadow-none flex flex-col text-center whitespace-nowrap text-white justify-center transition-all duration-500",
                                        disk.used_percent > 90 ? "bg-red-500" : disk.used_percent > 70 ? "bg-yellow-500" : "bg-blue-500"
                                    )}
                                ></div>
                            </div>
                        </div>

                        <div className="grid grid-cols-2 gap-4 mt-4 text-sm">
                            <div>
                                <div className="text-gray-500">Used</div>
                                <div className="font-mono text-white">{formatBytes(disk.used)}</div>
                            </div>
                            <div>
                                <div className="text-gray-500">Free</div>
                                <div className="font-mono text-white">{formatBytes(disk.free)}</div>
                            </div>
                        </div>

                        {/* I/O Rates */}
                        <div className="mt-4 pt-4 border-t border-gray-800 grid grid-cols-2 gap-4 text-sm">
                            <div>
                                <div className="text-gray-500">Read</div>
                                <div className="font-mono text-green-400">{formatBytes(disk.read_rate)}/s</div>
                            </div>
                            <div>
                                <div className="text-gray-500">Write</div>
                                <div className="font-mono text-red-400">{formatBytes(disk.write_rate)}/s</div>
                            </div>
                        </div>
                    </div>
                ))}
            </div>
        </div>
    );
};
