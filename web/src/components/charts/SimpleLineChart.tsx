import React from 'react';
import { AreaChart, Area, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer } from 'recharts';

interface SimpleLineChartProps {
    data: any[];
    dataKey: string;
    color?: string;
    height?: number;
    yDomain?: [number | string, number | string];
}

export const SimpleLineChart: React.FC<SimpleLineChartProps> = ({
    data,
    dataKey,
    color = "#3b82f6",
    height = 200,
    yDomain = [0, 100]
}) => {
    return (
        <div style={{ width: '100%', height }}>
            <ResponsiveContainer>
                <AreaChart data={data}>
                    <defs>
                        <linearGradient id={`color${dataKey}`} x1="0" y1="0" x2="0" y2="1">
                            <stop offset="5%" stopColor={color} stopOpacity={0.4} />
                            <stop offset="95%" stopColor={color} stopOpacity={0.05} />
                        </linearGradient>
                    </defs>
                    <CartesianGrid strokeDasharray="3 3" stroke="#1e293b" vertical={false} />
                    <XAxis dataKey="time" hide />
                    <YAxis
                        domain={yDomain}
                        stroke="#64748b"
                        fontSize={11}
                        tickFormatter={(val) => `${Number(val).toFixed(0)}`}
                        axisLine={false}
                        tickLine={false}
                        width={30}
                    />
                    <Tooltip
                        contentStyle={{ backgroundColor: '#0f172a', borderColor: '#334155', color: '#f8fafc', borderRadius: '6px', fontSize: '12px' }}
                        itemStyle={{ color: '#f8fafc' }}
                        labelStyle={{ display: 'none' }}
                        cursor={{ stroke: '#475569', strokeWidth: 1, strokeDasharray: '4 4' }}
                    />
                    <Area
                        type="monotone"
                        dataKey={dataKey}
                        stroke={color}
                        strokeWidth={2}
                        fillOpacity={1}
                        fill={`url(#color${dataKey})`}
                        isAnimationActive={false}
                    />
                </AreaChart>
            </ResponsiveContainer>
        </div>
    );
};
