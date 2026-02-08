'use client';

import { useState, useEffect } from 'react';
import { useSession, signOut } from 'next-auth/react';
import Link from 'next/link';
import { cn, getExperienceLabel } from '@/lib/utils';
import { LoadingSpinner } from '@/components/ui/loading';

interface Profile {
    id: string;
    name: string | null;
    email: string;
    experience: string;
    image: string | null;
    createdAt: string;
    cars: any[];
    _count: {
        trackReviews: number;
        lapRecords: number;
        tracks: number;
        zoneTips: number;
    };
}

const EXPERIENCE_OPTIONS = [
    { value: 'BEGINNER', label: 'Beginner', desc: 'New to motorsports' },
    { value: 'INTERMEDIATE', label: 'Intermediate', desc: 'A few track days' },
    { value: 'ADVANCED', label: 'Advanced', desc: 'Regular competitor' },
    { value: 'PRO', label: 'Pro', desc: 'Professional driver' },
];

export default function ProfilePage() {
    const { data: session } = useSession();
    const [profile, setProfile] = useState<Profile | null>(null);
    const [loading, setLoading] = useState(true);
    const [editing, setEditing] = useState(false);
    const [name, setName] = useState('');
    const [experience, setExperience] = useState('BEGINNER');
    const [saving, setSaving] = useState(false);

    useEffect(() => {
        if (session) fetchProfile();
    }, [session]);

    const fetchProfile = async () => {
        try {
            const res = await fetch('/api/profile');
            if (res.ok) {
                const data = await res.json();
                setProfile(data);
                setName(data.name || '');
                setExperience(data.experience);
            }
        } catch (error) {
            console.error('Failed to fetch profile:', error);
        } finally {
            setLoading(false);
        }
    };

    const handleSave = async () => {
        setSaving(true);
        try {
            const res = await fetch('/api/profile', {
                method: 'PUT',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ name, experience }),
            });

            if (res.ok) {
                setEditing(false);
                fetchProfile();
            }
        } catch (error) {
            console.error('Failed to update profile:', error);
        } finally {
            setSaving(false);
        }
    };

    if (!session) {
        return (
            <div className="page-container flex flex-col items-center justify-center min-h-[60vh]">
                <div className="text-center">
                    <div className="mx-auto mb-4 flex h-16 w-16 items-center justify-center rounded-full bg-surface-800">
                        <svg className="w-8 h-8 text-surface-500" fill="none" viewBox="0 0 24 24" strokeWidth={1.5} stroke="currentColor">
                            <path strokeLinecap="round" strokeLinejoin="round" d="M15.75 6a3.75 3.75 0 1 1-7.5 0 3.75 3.75 0 0 1 7.5 0ZM4.501 20.118a7.5 7.5 0 0 1 14.998 0A17.933 17.933 0 0 1 12 21.75c-2.676 0-5.216-.584-7.499-1.632Z" />
                        </svg>
                    </div>
                    <h2 className="text-xl font-bold text-white mb-2">Sign in to view profile</h2>
                    <p className="text-sm text-surface-400 mb-6">Manage your cars, experience, and track history</p>
                    <Link href="/login" className="btn-primary">Sign In</Link>
                </div>
            </div>
        );
    }

    if (loading) return <div className="page-container"><LoadingSpinner /></div>;

    return (
        <div className="page-container">
            <h1 className="page-title">Profile</h1>

            {/* Profile Card */}
            <div className="card mb-4">
                {editing ? (
                    <div className="space-y-4">
                        <div>
                            <label className="label">Driver Name</label>
                            <input
                                type="text"
                                value={name}
                                onChange={(e) => setName(e.target.value)}
                                className="input-field"
                                placeholder="Your name"
                            />
                        </div>
                        <div>
                            <label className="label">Experience Level</label>
                            <div className="grid grid-cols-2 gap-2">
                                {EXPERIENCE_OPTIONS.map((opt) => (
                                    <button
                                        key={opt.value}
                                        type="button"
                                        onClick={() => setExperience(opt.value)}
                                        className={cn(
                                            'rounded-lg border p-3 text-left transition-colors',
                                            experience === opt.value
                                                ? 'border-brand-500 bg-brand-500/10'
                                                : 'border-surface-600 bg-surface-800'
                                        )}
                                    >
                                        <p className="text-sm font-medium text-surface-100">{opt.label}</p>
                                        <p className="text-xs text-surface-500">{opt.desc}</p>
                                    </button>
                                ))}
                            </div>
                        </div>
                        <div className="flex gap-2">
                            <button onClick={() => setEditing(false)} className="btn-secondary flex-1 text-sm">Cancel</button>
                            <button onClick={handleSave} disabled={saving} className="btn-primary flex-1 text-sm">
                                {saving ? 'Saving...' : 'Save'}
                            </button>
                        </div>
                    </div>
                ) : (
                    <div>
                        <div className="flex items-center gap-3">
                            <div className="flex h-12 w-12 items-center justify-center rounded-full bg-brand-600 text-white text-lg font-bold">
                                {profile?.name?.[0]?.toUpperCase() || '?'}
                            </div>
                            <div className="flex-1">
                                <h2 className="font-semibold text-white">{profile?.name || 'Unknown'}</h2>
                                <p className="text-sm text-surface-400">{profile?.email}</p>
                            </div>
                            <button onClick={() => setEditing(true)} className="btn-ghost p-2 text-surface-400">
                                <svg className="w-5 h-5" fill="none" viewBox="0 0 24 24" strokeWidth={1.5} stroke="currentColor">
                                    <path strokeLinecap="round" strokeLinejoin="round" d="m16.862 4.487 1.687-1.688a1.875 1.875 0 1 1 2.652 2.652L10.582 16.07a4.5 4.5 0 0 1-1.897 1.13L6 18l.8-2.685a4.5 4.5 0 0 1 1.13-1.897l8.932-8.931Zm0 0L19.5 7.125M18 14v4.75A2.25 2.25 0 0 1 15.75 21H5.25A2.25 2.25 0 0 1 3 18.75V8.25A2.25 2.25 0 0 1 5.25 6H10" />
                                </svg>
                            </button>
                        </div>
                        <div className="mt-3 inline-flex items-center gap-1.5 rounded-full bg-surface-700 px-3 py-1 text-xs font-medium text-surface-300">
                            <svg className="w-3 h-3" fill="none" viewBox="0 0 24 24" strokeWidth={2} stroke="currentColor">
                                <path strokeLinecap="round" strokeLinejoin="round" d="M3 3v1.5M3 21v-6m0 0 2.77-.693a9 9 0 0 1 6.208.682l.108.054a9 9 0 0 0 6.086.71l3.114-.732a48.524 48.524 0 0 1-.005-10.499l-3.11.732a9 9 0 0 1-6.085-.711l-.108-.054a9 9 0 0 0-6.208-.682L3 4.5M3 15V4.5" />
                            </svg>
                            {getExperienceLabel(profile?.experience || '')}
                        </div>
                    </div>
                )}
            </div>

            {/* Stats */}
            <div className="grid grid-cols-2 gap-3 mb-4">
                <div className="card text-center">
                    <p className="text-2xl font-bold text-brand-500">{profile?._count.tracks || 0}</p>
                    <p className="text-xs text-surface-400 mt-1">Tracks Uploaded</p>
                </div>
                <div className="card text-center">
                    <p className="text-2xl font-bold text-brand-500">{profile?._count.trackReviews || 0}</p>
                    <p className="text-xs text-surface-400 mt-1">Reviews</p>
                </div>
                <div className="card text-center">
                    <p className="text-2xl font-bold text-brand-500">{profile?._count.lapRecords || 0}</p>
                    <p className="text-xs text-surface-400 mt-1">Lap Records</p>
                </div>
                <div className="card text-center">
                    <p className="text-2xl font-bold text-brand-500">{profile?._count.zoneTips || 0}</p>
                    <p className="text-xs text-surface-400 mt-1">Zone Tips</p>
                </div>
            </div>

            {/* Quick Links */}
            <div className="space-y-2 mb-4">
                <Link href="/garage" className="card flex items-center gap-3 hover:border-surface-600 transition-colors">
                    <div className="flex h-10 w-10 items-center justify-center rounded-lg bg-surface-700">
                        <svg className="w-5 h-5 text-surface-300" fill="none" viewBox="0 0 24 24" strokeWidth={1.5} stroke="currentColor">
                            <path strokeLinecap="round" strokeLinejoin="round" d="M11.42 15.17 17.25 21A2.652 2.652 0 0 0 21 17.25l-5.877-5.877M11.42 15.17l2.496-3.03c.317-.384.74-.626 1.208-.766M11.42 15.17l-4.655 5.653a2.548 2.548 0 1 1-3.586-3.586l6.837-5.63m5.108-.233c.55-.164 1.163-.188 1.743-.14a4.5 4.5 0 0 0 4.486-6.336l-3.276 3.277a3.004 3.004 0 0 1-2.25-2.25l3.276-3.276a4.5 4.5 0 0 0-6.336 4.486c.091 1.076-.071 2.264-.904 2.95l-.102.085m-1.745 1.437L5.909 7.5H4.5L2.25 3.75l1.5-1.5L7.5 4.5v1.409l4.26 4.26m-1.745 1.437 1.745-1.437m6.615 8.206L15.75 15.75M4.867 19.125h.008v.008h-.008v-.008Z" />
                        </svg>
                    </div>
                    <div className="flex-1">
                        <p className="font-medium text-surface-100">My Garage</p>
                        <p className="text-xs text-surface-500">{profile?.cars.length || 0} car{(profile?.cars.length || 0) !== 1 ? 's' : ''}</p>
                    </div>
                    <svg className="w-5 h-5 text-surface-500" fill="none" viewBox="0 0 24 24" strokeWidth={1.5} stroke="currentColor">
                        <path strokeLinecap="round" strokeLinejoin="round" d="m8.25 4.5 7.5 7.5-7.5 7.5" />
                    </svg>
                </Link>
                <Link href="/lapbook" className="card flex items-center gap-3 hover:border-surface-600 transition-colors">
                    <div className="flex h-10 w-10 items-center justify-center rounded-lg bg-surface-700">
                        <svg className="w-5 h-5 text-surface-300" fill="none" viewBox="0 0 24 24" strokeWidth={1.5} stroke="currentColor">
                            <path strokeLinecap="round" strokeLinejoin="round" d="M12 6v6h4.5m4.5 0a9 9 0 1 1-18 0 9 9 0 0 1 18 0Z" />
                        </svg>
                    </div>
                    <div className="flex-1">
                        <p className="font-medium text-surface-100">My Lap Book</p>
                        <p className="text-xs text-surface-500">{profile?._count.lapRecords || 0} recorded laps</p>
                    </div>
                    <svg className="w-5 h-5 text-surface-500" fill="none" viewBox="0 0 24 24" strokeWidth={1.5} stroke="currentColor">
                        <path strokeLinecap="round" strokeLinejoin="round" d="m8.25 4.5 7.5 7.5-7.5 7.5" />
                    </svg>
                </Link>
            </div>

            {/* Sign Out */}
            <button
                onClick={() => signOut({ callbackUrl: '/login' })}
                className="btn-outline w-full text-sm text-red-400 border-red-500/30 hover:bg-red-500/10"
            >
                Sign Out
            </button>
        </div>
    );
}
