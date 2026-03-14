import { create } from 'zustand';
import type { User } from '../domains/auth/entities/User';

interface AuthState {
    user: User | null;
    token: string | null;
    login: (user: User, token: string) => void;
    logout: () => void;
}

function loadUser(): User | null {
    try {
        const raw = localStorage.getItem('user');
        return raw ? JSON.parse(raw) : null;
    } catch {
        return null;
    }
}

export const useAuthStore = create<AuthState>((set) => ({
    user: loadUser(),
    token: localStorage.getItem('token'),
    login: (user, token) => {
        localStorage.setItem('token', token);
        localStorage.setItem('user', JSON.stringify(user));
        set({ user, token });
    },
    logout: () => {
        localStorage.removeItem('token');
        localStorage.removeItem('user');
        set({ user: null, token: null });
    },
}));
