const API_BASE = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';

export function getToken(): string | null {
    if (typeof window === 'undefined') return null;
    return localStorage.getItem('token');
}

export function setToken(token: string): void {
    localStorage.setItem('token', token);
}

export function clearToken(): void {
    localStorage.removeItem('token');
}

/**
 * Wrapper around fetch() that prepends the Go backend base URL
 * and attaches the JWT Authorization header when available.
 */
export async function api(path: string, options: RequestInit = {}): Promise<Response> {
    const token = getToken();
    const headers = new Headers(options.headers);

    if (token && !headers.has('Authorization')) {
        headers.set('Authorization', `Bearer ${token}`);
    }

    // Don't set Content-Type for FormData (browser sets it with boundary)
    if (options.body instanceof FormData) {
        headers.delete('Content-Type');
    }

    return fetch(`${API_BASE}${path}`, {
        ...options,
        headers,
    });
}
