'use client';

import { useState } from 'react';
import { signIn } from 'next-auth/react';
import { useRouter } from 'next/navigation';
import Link from 'next/link';

export default function LoginPage() {
    const router = useRouter();
    const [email, setEmail] = useState('');
    const [password, setPassword] = useState('');
    const [error, setError] = useState('');
    const [loading, setLoading] = useState(false);

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        setError('');
        setLoading(true);

        try {
            const result = await signIn('credentials', {
                email,
                password,
                redirect: false,
            });

            if (result?.error) {
                setError('Invalid email or password');
            } else {
                router.push('/');
                router.refresh();
            }
        } catch {
            setError('Something went wrong');
        } finally {
            setLoading(false);
        }
    };

    return (
        <div className="flex min-h-screen flex-col items-center justify-center px-4">
            <div className="w-full max-w-sm">
                {/* Logo */}
                <div className="mb-8 text-center">
                    <div className="mx-auto mb-3 flex h-14 w-14 items-center justify-center rounded-2xl bg-brand-600">
                        <svg className="w-8 h-8 text-white" fill="none" viewBox="0 0 24 24" strokeWidth={2} stroke="currentColor">
                            <path strokeLinecap="round" strokeLinejoin="round" d="M3 3v1.5M3 21v-6m0 0 2.77-.693a9 9 0 0 1 6.208.682l.108.054a9 9 0 0 0 6.086.71l3.114-.732a48.524 48.524 0 0 1-.005-10.499l-3.11.732a9 9 0 0 1-6.085-.711l-.108-.054a9 9 0 0 0-6.208-.682L3 4.5M3 15V4.5" />
                        </svg>
                    </div>
                    <h1 className="text-2xl font-bold text-white">Trackside</h1>
                    <p className="mt-1 text-sm text-surface-400">Sign in to your driver account</p>
                </div>

                <form onSubmit={handleSubmit} className="space-y-4">
                    {error && (
                        <div className="rounded-lg bg-red-500/10 border border-red-500/20 p-3 text-sm text-red-400">
                            {error}
                        </div>
                    )}

                    <div>
                        <label htmlFor="email" className="label">Email</label>
                        <input
                            id="email"
                            type="email"
                            value={email}
                            onChange={(e) => setEmail(e.target.value)}
                            className="input-field"
                            placeholder="driver@example.com"
                            required
                            autoComplete="email"
                        />
                    </div>

                    <div>
                        <label htmlFor="password" className="label">Password</label>
                        <input
                            id="password"
                            type="password"
                            value={password}
                            onChange={(e) => setPassword(e.target.value)}
                            className="input-field"
                            placeholder="••••••••"
                            required
                            autoComplete="current-password"
                        />
                    </div>

                    <button
                        type="submit"
                        disabled={loading}
                        className="btn-primary w-full"
                    >
                        {loading ? (
                            <div className="h-5 w-5 animate-spin rounded-full border-2 border-white/30 border-t-white" />
                        ) : (
                            'Sign In'
                        )}
                    </button>
                </form>

                <p className="mt-6 text-center text-sm text-surface-400">
                    Don&apos;t have an account?{' '}
                    <Link href="/register" className="font-medium text-brand-500 hover:text-brand-400">
                        Sign up
                    </Link>
                </p>
            </div>
        </div>
    );
}
