import axios from 'axios';

const API_URL = import.meta.env.VITE_API_URL || '/api/v1';
console.log('API_URL:', API_URL);

export const setApiKey = (key: string) => {
    console.log('Setting API Key:', key);
    localStorage.setItem('server_moni_key', key);
};

export const getApiKey = () => {
    const key = localStorage.getItem('server_moni_key');
    console.log('Getting API Key:', key ? 'Found' : 'Not Found');
    return key;
};

export const client = axios.create({
    baseURL: API_URL,
});

client.interceptors.request.use((config) => {
    const key = getApiKey();
    if (key) {
        console.log('Attaching Authorization Header');
        config.headers.Authorization = `Bearer ${key.trim()}`;
    } else {
        console.warn('No API Key found for request');
    }
    return config;
});

export interface System {
    id: number;
    name: string;
    url: string;
    api_key: string;
    created_at: string;
}

export const fetchSystems = async (): Promise<System[]> => {
    const response = await client.get('/systems');
    return response.data || [];
};

export const addSystem = async (name: string, url: string, apiKey: string) => {
    const response = await client.post('/systems', { name, url, api_key: apiKey });
    return response.data;
};

export const deleteSystem = async (id: number) => {
    await client.delete(`/systems/${id}`);
};

export const fetchMetrics = async (systemId?: string | number) => {
    const params = systemId ? { system_id: systemId } : {};
    const response = await client.get('/metrics', { params });
    return response.data;
};

export const fetchStatus = async () => {
    // This might fail if we are not connected to a local agent directly
    // But for the dashboard, we might want to check backend status
    // For now, let's just return true
    return { status: 'ok' };
};

export const verifyKey = async (_key: string) => {
    // Legacy verification, might not be needed for User Token
    return true;
};

export const login = async (email: string, password: string) => {
    const response = await client.post('/auth/login', { email, password });
    if (response.data.token) {
        setApiKey(response.data.token);
    }
    return response.data;
};

export const register = async (email: string, password: string) => {
    const response = await client.post('/auth/register', { email, password });
    return response.data;
};

export const logout = async () => {
    try {
        await client.post('/auth/logout');
    } catch (e) {
        // ignore
    }
    localStorage.removeItem('server_moni_key');
};
