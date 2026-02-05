import React, { createContext, useState, useContext, useEffect } from 'react';
import axios from 'axios';

const DeveloperContext = createContext();

export const DeveloperProvider = ({ children }) => {
    const [isDevMode, setIsDevMode] = useState(true); // Default to true for phase 2
    const [apiLogs, setApiLogs] = useState([]);

    const addLog = (log) => {
        setApiLogs(prev => [log, ...prev].slice(0, 50)); // Keep last 50
    };

    const clearLogs = () => setApiLogs([]);

    // Setup Axios Interceptors for Global Logging
    useEffect(() => {
        const reqInterceptor = axios.interceptors.request.use(config => {
            config.metadata = { startTime: new Date() };
            return config;
        });

        const resInterceptor = axios.interceptors.response.use(
            response => {
                const duration = new Date() - response.config.metadata.startTime;
                if (isDevMode) {
                    addLog({
                        id: Date.now() + Math.random(),
                        timestamp: new Date(),
                        method: response.config.method.toUpperCase(),
                        url: response.config.url,
                        status: response.status,
                        duration,
                        requestData: response.config.data,
                        responseData: response.data,
                        type: 'success'
                    });
                }
                return response;
            },
            error => {
                const duration = new Date() - (error.config?.metadata?.startTime || new Date());
                if (isDevMode) {
                    addLog({
                        id: Date.now() + Math.random(),
                        timestamp: new Date(),
                        method: error.config?.method?.toUpperCase() || 'UNKNOWN',
                        url: error.config?.url || 'UNKNOWN',
                        status: error.response?.status || 0,
                        duration,
                        requestData: error.config?.data,
                        responseData: error.response?.data || error.message,
                        type: 'error'
                    });
                }
                return Promise.reject(error);
            }
        );

        return () => {
            axios.interceptors.request.eject(reqInterceptor);
            axios.interceptors.response.eject(resInterceptor);
        };
    }, [isDevMode]);

    return (
        <DeveloperContext.Provider value={{ isDevMode, setIsDevMode, apiLogs, clearLogs }}>
            {children}
        </DeveloperContext.Provider>
    );
};

export const useDeveloper = () => useContext(DeveloperContext);
