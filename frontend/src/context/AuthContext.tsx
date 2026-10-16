import { createContext, useContext, useEffect, useState, type ReactNode } from 'react';
import { getToken, setToken, request } from '../api/client';
import type { AuthResponse, AuthUser } from '../api/types';

interface AuthContextValue {
  user: AuthUser | null;
  isAuthenticated: boolean;
  login: (username: string, password: string) => Promise<void>;
  signup: (username: string, password: string) => Promise<void>;
  logout: () => void;
}

const AuthContext = createContext<AuthContextValue | undefined>(undefined);

const USER_KEY = 'pokedex_user';

function loadStoredUser(): AuthUser | null {
  try {
    const raw = localStorage.getItem(USER_KEY);
    return raw ? (JSON.parse(raw) as AuthUser) : null;
  } catch {
    return null;
  }
}

function storeUser(user: AuthUser | null) {
  try {
    if (user) localStorage.setItem(USER_KEY, JSON.stringify(user));
    else localStorage.removeItem(USER_KEY);
  } catch {
    // ignore storage failures
  }
}

export function AuthProvider({ children }: { children: ReactNode }) {
  const [user, setUser] = useState<AuthUser | null>(() => (getToken() ? loadStoredUser() : null));

  // Single choke point for session expiry: client.ts dispatches this event
  // on any 401, so every page reacts the same way without checking status
  // codes itself.
  useEffect(() => {
    function handleExpired() {
      setUser(null);
      storeUser(null);
    }
    window.addEventListener('auth:expired', handleExpired);
    return () => window.removeEventListener('auth:expired', handleExpired);
  }, []);

  function applyAuth(data: AuthResponse) {
    setToken(data.token);
    storeUser(data.user);
    setUser(data.user);
  }

  async function signup(username: string, password: string) {
    const data = await request<AuthResponse>('/api/auth/signup', {
      method: 'POST',
      body: JSON.stringify({ username, password }),
    });
    applyAuth(data);
  }

  async function login(username: string, password: string) {
    const data = await request<AuthResponse>('/api/auth/login', {
      method: 'POST',
      body: JSON.stringify({ username, password }),
    });
    applyAuth(data);
  }

  function logout() {
    setToken(null);
    storeUser(null);
    setUser(null);
  }

  return (
    <AuthContext.Provider value={{ user, isAuthenticated: !!user, login, signup, logout }}>
      {children}
    </AuthContext.Provider>
  );
}

export function useAuth(): AuthContextValue {
  const ctx = useContext(AuthContext);
  if (!ctx) throw new Error('useAuth must be used within AuthProvider');
  return ctx;
}
