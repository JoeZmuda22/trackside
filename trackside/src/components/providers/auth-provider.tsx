'use client';

import { createContext, useContext, useState, useEffect, useCallback } from 'react';
import { api, setToken as storeToken, clearToken, getToken } from '@/lib/api';

interface User {
    id: string;
    email: string;
    name: string | null;
}

interface AuthContextType {
    user: User | null;
    loading: boolean;
    login: (email: string, password: string) => Promise<{ error?: string }>;
    logout: () => void;
    updateUser: (user: User) => void;
}

const AuthContext = createContext<AuthContextType>({
    user: null,
    loading: true,
    login: async () => ({ error: 'Not initialized' }),
    logout: () => {},
    updateUser: () => {},
});

export function useAuth() {
    return useContext(AuthContext);
}

export function AuthProvider({ children }: { children: React.ReactNode }) {
    const [user, setUser] = useState<User | null>(null);
    const [loading, setLoading] = useState(true);

    // On mount, check for stored auth
    useEffect(() => {
        const token = getToken();
        const storedUser = localStorage.getItem('user');
        if (token && storedUser) {
            try {
                setUser(JSON.parse(storedUser));
            } catch {
                clearToken();
                localStorage.removeItem('user');
            }
        }
        setLoading(false);
    }, []);

    const login = useCallback(async (email: string, password: string) => {
        try {
            const res = await api('/api/auth/login', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ email, password }),
            });

            if (!res.ok) {
                const data = await res.json().catch(() => ({}));
                return { error: data.error || 'Invalid email or password' };
            }

            const data = await res.json();
            storeToken(data.token);
            const userData: User = {
                id: data.user.id,
                email: data.user.email,
                name: data.user.name,
            };
            localStorage.setItem('user', JSON.stringify(userData));
            setUser(userData);
            return {};
        } catch {
            return { error: 'Something went wrong' };
        }
    }, []);

    const logout = useCallback(() => {
        clearToken();
        localStorage.removeItem('user');
        setUser(null);
    }, []);

    const updateUser = useCallback((updated: User) => {
        localStorage.setItem('user', JSON.stringify(updated));
        setUser(updated);
    }, []);

    return (
        <AuthContext.Provider value={{ user, loading, login, logout, updateUser }}>
            {children}
        </AuthContext.Provider>
    );
}
