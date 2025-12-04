import React, { useEffect, useState } from 'react';
import { Shield, AlertTriangle, Lock, Activity } from 'lucide-react';
import { client } from '../../utils/api';

interface Fail2BanStats {
    total_bans: number;
    bans_by_ip: Record<string, number>;
    jails: string[];
}

interface AuthLog {
    time: string;
    user: string;
    ip: string;
    message: string;
    success: boolean;
}

interface SecuritySectionProps {
    systemId: number;
}

const SecuritySection: React.FC<SecuritySectionProps> = ({ systemId }) => {
    const [fail2ban, setFail2ban] = useState<Fail2BanStats | null>(null);
    const [authLogs, setAuthLogs] = useState<AuthLog[]>([]);
    const [loading, setLoading] = useState(true);

    useEffect(() => {
        const fetchData = async () => {
            try {
                // Fetch Fail2Ban
                const f2bResponse = await client.get(`/systems/${systemId}/proxy`, {
                    params: { path: '/security/fail2ban' }
                });
                setFail2ban(f2bResponse.data);

                // Fetch Auth Logs
                const authResponse = await client.get(`/systems/${systemId}/proxy`, {
                    params: { path: '/security/logins' }
                });
                setAuthLogs(authResponse.data);

            } catch (error) {
                console.error("Failed to fetch security stats", error);
            } finally {
                setLoading(false);
            }
        };

        fetchData();
    }, [systemId]);

    if (loading) return <div className="p-4 text-gray-400">Loading security stats...</div>;

    return (
        <div className="space-y-6">
            {/* Fail2Ban Stats */}
            <div className="bg-gray-800 rounded-lg p-6 border border-gray-700">
                <div className="flex items-center gap-2 mb-4">
                    <Shield className="w-5 h-5 text-green-400" />
                    <h3 className="text-lg font-medium text-white">Fail2Ban Status</h3>
                </div>

                <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
                    <div className="bg-gray-900/50 p-4 rounded-lg">
                        <div className="text-gray-400 text-sm">Total Bans</div>
                        <div className="text-2xl font-bold text-red-400">{fail2ban?.total_bans || 0}</div>
                    </div>
                    <div className="bg-gray-900/50 p-4 rounded-lg">
                        <div className="text-gray-400 text-sm">Active Jails</div>
                        <div className="flex flex-wrap gap-2 mt-2">
                            {fail2ban?.jails?.map(jail => (
                                <span key={jail} className="px-2 py-1 bg-gray-800 text-xs rounded text-gray-300 border border-gray-700">
                                    {jail}
                                </span>
                            )) || <span className="text-gray-500 text-sm">None</span>}
                        </div>
                    </div>
                </div>
            </div>

            {/* Auth Logs */}
            <div className="bg-gray-800 rounded-lg p-6 border border-gray-700">
                <div className="flex items-center gap-2 mb-4">
                    <Lock className="w-5 h-5 text-blue-400" />
                    <h3 className="text-lg font-medium text-white">Recent Login Attempts</h3>
                </div>

                <div className="overflow-x-auto">
                    <table className="w-full text-left text-sm text-gray-400">
                        <thead className="bg-gray-900/50 text-gray-300 uppercase text-xs">
                            <tr>
                                <th className="px-4 py-3">Time</th>
                                <th className="px-4 py-3">User</th>
                                <th className="px-4 py-3">IP</th>
                                <th className="px-4 py-3">Status</th>
                            </tr>
                        </thead>
                        <tbody className="divide-y divide-gray-700">
                            {authLogs.map((log, i) => (
                                <tr key={i} className="hover:bg-gray-700/30">
                                    <td className="px-4 py-3 whitespace-nowrap">{new Date(log.time).toLocaleString()}</td>
                                    <td className="px-4 py-3">{log.user || 'unknown'}</td>
                                    <td className="px-4 py-3 font-mono text-xs">{log.ip || '-'}</td>
                                    <td className="px-4 py-3">
                                        {log.success ? (
                                            <span className="text-green-400 flex items-center gap-1"><Activity className="w-3 h-3" /> Success</span>
                                        ) : (
                                            <span className="text-red-400 flex items-center gap-1"><AlertTriangle className="w-3 h-3" /> Failed</span>
                                        )}
                                    </td>
                                </tr>
                            ))}
                            {authLogs.length === 0 && (
                                <tr>
                                    <td colSpan={4} className="px-4 py-8 text-center text-gray-500">No recent logs found</td>
                                </tr>
                            )}
                        </tbody>
                    </table>
                </div>
            </div>
        </div>
    );
};

export default SecuritySection;
