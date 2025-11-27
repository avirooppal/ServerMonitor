import { useState, useEffect } from 'react';

export function useMetricHistory<T>(value: T, maxPoints: number = 60) {
    const [history, setHistory] = useState<{ time: number; value: T }[]>([]);

    useEffect(() => {
        setHistory(prev => {
            const newHistory = [...prev, { time: Date.now(), value }];
            if (newHistory.length > maxPoints) {
                return newHistory.slice(newHistory.length - maxPoints);
            }
            return newHistory;
        });
    }, [value, maxPoints]);

    return history;
}
