import axios from 'axios';

// Client-Side System Management (localStorage)

export interface System {
    id: string;
    name: string;
    url: string;
    api_key: string; // This is the AGENT_SECRET
    created_at: string;
}

const STORAGE_KEY = 'server_moni_systems';

export const getSystems = (): System[] => {
    const data = localStorage.getItem(STORAGE_KEY);
    if (!data) return [];
    try {
        return JSON.parse(data);
    } catch (e) {
        console.error("Failed to parse systems from localStorage", e);
        return [];
    }
};

export const saveSystem = (system: System) => {
    const systems = getSystems();
    systems.push(system);
    localStorage.setItem(STORAGE_KEY, JSON.stringify(systems));
};

export const deleteSystem = (id: string) => {
    const systems = getSystems().filter(s => s.id !== id);
    localStorage.setItem(STORAGE_KEY, JSON.stringify(systems));
};

// Direct Agent Communication

export const fetchMetrics = async (system: System) => {
    try {
        const response = await axios.get(`${system.url}/api/v1/metrics`, {
            headers: {
                'Authorization': `Bearer ${system.api_key}`
            },
            timeout: 5000
        });
        return response.data;
    } catch (error) {
        console.error(`Failed to fetch metrics from ${system.name}:`, error);
        throw error;
    }
};

// Helper to generate random token
export const generateToken = () => {
    return Math.random().toString(36).substring(2) + Math.random().toString(36).substring(2);
};

