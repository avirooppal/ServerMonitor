import axios from 'axios';

const API_URL = import.meta.env.VITE_API_URL || (import.meta.env.PROD ? '/api/v1' : 'http://localhost:8080/api/v1');

export const setApiKey = (key: string) => {
    localStorage.setItem('server_moni_key', key);
};

export const getApiKey = () => {
    return localStorage.getItem('server_moni_key');
};

export const client = axios.create({
    baseURL: API_URL,
});

client.interceptors.request.use((config) => {
    const key = getApiKey();
    if (key) {
        config.headers.Authorization = `Bearer ${key.trim()}`;
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

export const verifyKey = async (key: string) => {
    const testClient = axios.create({
        baseURL: API_URL,
        headers: { Authorization: `Bearer ${key}` }
    });
    await testClient.post('/verify-key');
    return true;
};

export const login = async (username: string, password: string): Promise<string> => {
    const response = await client.post('/login', { username, password });
    return response.data.api_key;
};

export const register = async (username: string, password: string): Promise<string> => {
    const response = await client.post('/register', { username, password });
    return response.data.api_key;
};
