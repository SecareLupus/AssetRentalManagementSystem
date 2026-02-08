import React, { createContext, useState, useContext, useEffect } from 'react';
import axios from 'axios';

const AuthContext = createContext();

export const AuthProvider = ({ children }) => {
    const [user, setUser] = useState(null);
    const [token, setToken] = useState(() => {
        const storedToken = localStorage.getItem('rms_token');
        if (storedToken) {
            axios.defaults.headers.common['Authorization'] = `Bearer ${storedToken}`;
        }
        return storedToken;
    });
    const [loading, setLoading] = useState(true);

    useEffect(() => {
        if (token) {
            localStorage.setItem('rms_token', token);
            // Setup global axios auth header
            axios.defaults.headers.common['Authorization'] = `Bearer ${token}`;
        } else {
            localStorage.removeItem('rms_token');
            delete axios.defaults.headers.common['Authorization'];
        }
    }, [token]);

    // Check if token is valid (could add verify endpoint call here)
    useEffect(() => {
        const storedUser = localStorage.getItem('rms_user');
        if (storedUser && token) {
            try {
                setUser(JSON.parse(storedUser));
            } catch (e) {
                console.error("Failed to parse stored user", e);
            }
        }
        setLoading(false);
    }, [token]);

    const login = async (username, password) => {
        try {
            const response = await axios.post('/v1/auth/login', { username, password });
            const { token: newToken, user: userData } = response.data;

            // Set header IMMEDIATELY to prevent race condition with subsequent requests
            axios.defaults.headers.common['Authorization'] = `Bearer ${newToken}`;

            setToken(newToken);
            setUser(userData);
            localStorage.setItem('rms_token', newToken); // Also set manually here for safety
            localStorage.setItem('rms_user', JSON.stringify(userData));
            return true;
        } catch (error) {
            console.error("Login failed", error);
            throw error;
        }
    };

    const logout = () => {
        setToken(null);
        setUser(null);
        localStorage.removeItem('rms_user');
        window.location.href = '/';
    };

    // Axios interceptor for 401s (automatic logout on expiry)
    useEffect(() => {
        const interceptor = axios.interceptors.response.use(
            response => response,
            error => {
                if (error.response?.status === 401) {
                    logout();
                }
                return Promise.reject(error);
            }
        );
        return () => axios.interceptors.response.eject(interceptor);
    }, []);

    const value = {
        user,
        token,
        loading,
        login,
        logout,
        isAuthenticated: !!token
    };

    return (
        <AuthContext.Provider value={value}>
            {!loading && children}
        </AuthContext.Provider>
    );
};

export const useAuth = () => useContext(AuthContext);
