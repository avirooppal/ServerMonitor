import axios from 'axios';

const API_URL = import.meta.env.VITE_API_URL || (import.meta.env.PROD ? '/api/v1' : 'http://localhost:8080/api/v1');

export const setApiKey = (key: string) => {
    localStorage.setItem('server_moni_key', key);
};

export const getApiKey = () => {
    return localStorage.getItem('server_moni_key');
};

const client = axios.create({
    baseURL: API_URL,
});

client.interceptors.request.use((config) => {
    const key = getApiKey();
    if (key) {
        config.headers.Authorization = `Bearer ${key}`;
    }
    return config;
});

export const fetchMetrics = async () => {
    const response = await client.get('/metrics');
    return response.data;
};

export const fetchStatus = async () => {
    const response = await client.get('/status');
    return response.data;
};

export const verifyKey = async (key: string) => {
    // We can use a dummy call or a specific endpoint
    // For now, let's try to fetch status with the key
    const testClient = axios.create({
        baseURL: API_URL,
        headers: { Authorization: `Bearer ${key}` }
    });
    await testClient.get('/status');
    return true;
};
