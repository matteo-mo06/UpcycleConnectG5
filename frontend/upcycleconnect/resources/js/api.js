import axios from 'axios';

const api = axios.create({
    baseURL: 'http://localhost:8084',
    headers: {
        'Content-Type': 'application/json',
    },
});

export function setAuthToken(token) {
    api.defaults.headers.common['Authorization'] = `Bearer ${token}`;
}

export function clearAuthToken() {
    delete api.defaults.headers.common['Authorization'];
}

export default api;
