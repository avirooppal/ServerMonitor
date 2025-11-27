import React from 'react';

interface CircularProgressProps {
    value: number;
    label: string;
    color: string;
    size?: number;
    strokeWidth?: number;
}

export const CircularProgress: React.FC<CircularProgressProps> = ({
    value,
    label,
    color,
    size = 120,
    strokeWidth = 10
}) => {
    const radius = (size - strokeWidth) / 2;
    const circumference = radius * 2 * Math.PI;
    const offset = circumference - (value / 100) * circumference;

    return (
        <div className="flex flex-col items-center justify-center">
            <div className="relative" style={{ width: size, height: size }}>
                {/* Background Circle */}
                <svg className="w-full h-full transform -rotate-90">
                    <circle
                        cx={size / 2}
                        cy={size / 2}
                        r={radius}
                        stroke="#1e293b" // Dark gray background
                        strokeWidth={strokeWidth}
                        fill="transparent"
                    />
                    {/* Progress Circle */}
                    <circle
                        cx={size / 2}
                        cy={size / 2}
                        r={radius}
                        stroke={color}
                        strokeWidth={strokeWidth}
                        fill="transparent"
                        strokeDasharray={circumference}
                        strokeDashoffset={offset}
                        strokeLinecap="round"
                        className="transition-all duration-1000 ease-out"
                    />
                </svg>
                {/* Center Text */}
                <div className="absolute inset-0 flex flex-col items-center justify-center">
                    <span className="text-2xl font-bold text-white">{Math.round(value)}%</span>
                </div>
            </div>
            <span className="mt-2 text-sm text-gray-400 font-medium uppercase tracking-wider">{label}</span>
        </div>
    );
};
