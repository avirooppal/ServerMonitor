import React, { useState, useMemo } from 'react';
import type { SystemMetrics } from '../../types';
import { Search, ArrowUpDown } from 'lucide-react';

interface ProcessSectionProps {
    metrics: SystemMetrics;
}

type SortField = 'cpu' | 'mem' | 'name' | 'pid';

export const ProcessSection: React.FC<ProcessSectionProps> = ({ metrics }) => {
    const [searchTerm, setSearchTerm] = useState('');
    const [sortField, setSortField] = useState<SortField>('cpu');
    const [sortDesc, setSortDesc] = useState(true);

    const filteredProcesses = useMemo(() => {
        if (!metrics.processes) return [];
        let procs = [...metrics.processes];

        // Filter
        if (searchTerm) {
            const term = searchTerm.toLowerCase();
            procs = procs.filter(p =>
                (p.name || '').toLowerCase().includes(term) ||
                (p.username || '').toLowerCase().includes(term) ||
                p.pid.toString().includes(term)
            );
        }

        // Sort
        procs.sort((a, b) => {
            let valA = a[sortField];
            let valB = b[sortField];

            // Handle string comparison for name
            if (typeof valA === 'string' && typeof valB === 'string') {
                return sortDesc ? valB.localeCompare(valA) : valA.localeCompare(valB);
            }

            // Handle number comparison
            if (valA < valB) return sortDesc ? 1 : -1;
            if (valA > valB) return sortDesc ? -1 : 1;
            return 0;
        });

        return procs;
    }, [metrics.processes, searchTerm, sortField, sortDesc]);

    const handleSort = (field: SortField) => {
        if (sortField === field) {
            setSortDesc(!sortDesc);
        } else {
            setSortField(field);
            setSortDesc(true);
        }
    };

    return (
        <div className="space-y-6">
            <div className="bg-gray-900 rounded-xl border border-gray-800 shadow-sm overflow-hidden">
                <div className="px-6 py-4 border-b border-gray-800 flex flex-col md:flex-row justify-between items-center gap-4">
                    <div className="flex items-center space-x-2">
                        <h3 className="text-lg font-semibold text-gray-200">Active Processes</h3>
                        <span className="text-sm text-gray-500">({metrics.processes?.length || 0} total)</span>
                    </div>

                    <div className="relative w-full md:w-64">
                        <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-500 w-4 h-4" />
                        <input
                            type="text"
                            placeholder="Search process..."
                            className="w-full bg-gray-800 border border-gray-700 text-gray-200 text-sm rounded-lg pl-9 pr-3 py-2 focus:outline-none focus:border-blue-500 transition-colors"
                            value={searchTerm}
                            onChange={(e) => setSearchTerm(e.target.value)}
                        />
                    </div>
                </div>

                <div className="overflow-x-auto">
                    <table className="w-full text-left text-sm text-gray-400">
                        <thead className="bg-gray-800 text-gray-200 uppercase font-medium">
                            <tr>
                                <th className="px-6 py-3 cursor-pointer hover:bg-gray-700/50 transition-colors" onClick={() => handleSort('pid')}>
                                    <div className="flex items-center space-x-1">
                                        <span>PID</span>
                                        <ArrowUpDown size={12} className={sortField === 'pid' ? 'text-blue-400' : 'text-gray-600'} />
                                    </div>
                                </th>
                                <th className="px-6 py-3 cursor-pointer hover:bg-gray-700/50 transition-colors" onClick={() => handleSort('name')}>
                                    <div className="flex items-center space-x-1">
                                        <span>Name</span>
                                        <ArrowUpDown size={12} className={sortField === 'name' ? 'text-blue-400' : 'text-gray-600'} />
                                    </div>
                                </th>
                                <th className="px-6 py-3">User</th>
                                <th className="px-6 py-3 text-right cursor-pointer hover:bg-gray-700/50 transition-colors" onClick={() => handleSort('cpu')}>
                                    <div className="flex items-center justify-end space-x-1">
                                        <span>CPU %</span>
                                        <ArrowUpDown size={12} className={sortField === 'cpu' ? 'text-blue-400' : 'text-gray-600'} />
                                    </div>
                                </th>
                                <th className="px-6 py-3 text-right cursor-pointer hover:bg-gray-700/50 transition-colors" onClick={() => handleSort('mem')}>
                                    <div className="flex items-center justify-end space-x-1">
                                        <span>Mem %</span>
                                        <ArrowUpDown size={12} className={sortField === 'mem' ? 'text-blue-400' : 'text-gray-600'} />
                                    </div>
                                </th>
                            </tr>
                        </thead>
                        <tbody className="divide-y divide-gray-800">
                            {filteredProcesses.map((proc) => (
                                <tr key={proc.pid} className="hover:bg-gray-800/50 transition-colors">
                                    <td className="px-6 py-4 font-mono">{proc.pid}</td>
                                    <td className="px-6 py-4 font-medium text-white">{proc.name}</td>
                                    <td className="px-6 py-4">{proc.username}</td>
                                    <td className="px-6 py-4 text-right font-mono text-blue-400">{proc.cpu.toFixed(1)}%</td>
                                    <td className="px-6 py-4 text-right font-mono text-purple-400">{proc.mem.toFixed(1)}%</td>
                                </tr>
                            ))}
                            {filteredProcesses.length === 0 && (
                                <tr>
                                    <td colSpan={5} className="px-6 py-8 text-center text-gray-500">
                                        No processes found matching "{searchTerm}"
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
