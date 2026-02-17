'use client';

import { useState, useRef } from 'react';
import { useRouter } from 'next/navigation';
import { useAuth } from '@/components/providers/auth-provider';
import { api } from '@/lib/api';
import Link from 'next/link';
import { cn } from '@/lib/utils';

const EVENT_OPTIONS = [
    { value: 'AUTOCROSS', label: 'Autocross', description: 'Cone courses & parking lots' },
    { value: 'ROADCOURSE', label: 'Road Course', description: 'Time attack & racing' },
    { value: 'DRIFT', label: 'Drift', description: 'Sideways action' },
    { value: 'DRAG', label: 'Drag', description: 'Straight-line speed' },
];

export default function NewTrackPage() {
    const router = useRouter();
    const { user } = useAuth();
    const fileInputRef = useRef<HTMLInputElement>(null);

    const [name, setName] = useState('');
    const [location, setLocation] = useState('');
    const [description, setDescription] = useState('');
    const [eventTypes, setEventTypes] = useState<string[]>([]);
    const [imageUrl, setImageUrl] = useState('');
    const [imagePreview, setImagePreview] = useState<string | null>(null);
    const [uploading, setUploading] = useState(false);
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState('');

    if (!user) {
        return (
            <div className="page-container flex flex-col items-center justify-center min-h-[60vh]">
                <p className="text-surface-400 mb-4">You need to sign in to upload a track</p>
                <Link href="/login" className="btn-primary">Sign In</Link>
            </div>
        );
    }

    const toggleEventType = (value: string) => {
        setEventTypes((prev) =>
            prev.includes(value)
                ? prev.filter((v) => v !== value)
                : [...prev, value]
        );
    };

    const handleImageUpload = async (e: React.ChangeEvent<HTMLInputElement>) => {
        const file = e.target.files?.[0];
        if (!file) return;

        // Show preview
        const reader = new FileReader();
        reader.onload = (e) => setImagePreview(e.target?.result as string);
        reader.readAsDataURL(file);

        setUploading(true);
        try {
            const formData = new FormData();
            formData.append('file', file);

            const res = await api('/api/upload', {
                method: 'POST',
                body: formData,
            });

            if (!res.ok) {
                const data = await res.json();
                setError(data.error || 'Upload failed');
                setImagePreview(null);
                return;
            }

            const data = await res.json();
            setImageUrl(data.imageUrl);
        } catch {
            setError('Image upload failed');
            setImagePreview(null);
        } finally {
            setUploading(false);
        }
    };

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        setError('');

        if (eventTypes.length === 0) {
            setError('Select at least one event type');
            return;
        }

        setLoading(true);

        try {
            const res = await api('/api/tracks', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({
                    name,
                    location,
                    description: description || undefined,
                    imageUrl: imageUrl || undefined,
                    eventTypes,
                }),
            });

            if (!res.ok) {
                const data = await res.json();
                setError(data.error || 'Failed to create track');
                return;
            }

            const track = await res.json();
            router.push(`/tracks/${track.id}`);
        } catch {
            setError('Something went wrong');
        } finally {
            setLoading(false);
        }
    };

    return (
        <div className="page-container">
            <div className="flex items-center gap-3 mb-6">
                <button onClick={() => router.back()} className="btn-ghost p-2">
                    <svg className="w-5 h-5" fill="none" viewBox="0 0 24 24" strokeWidth={2} stroke="currentColor">
                        <path strokeLinecap="round" strokeLinejoin="round" d="M15.75 19.5 8.25 12l7.5-7.5" />
                    </svg>
                </button>
                <h1 className="page-title mb-0">Upload Track</h1>
            </div>

            <form onSubmit={handleSubmit} className="space-y-5">
                {error && (
                    <div className="rounded-lg bg-red-500/10 border border-red-500/20 p-3 text-sm text-red-400">
                        {error}
                    </div>
                )}

                {/* Track Image Upload */}
                <div>
                    <label className="label">Track Layout Image</label>
                    <div
                        onClick={() => fileInputRef.current?.click()}
                        className={cn(
                            'relative flex flex-col items-center justify-center rounded-xl border-2 border-dashed transition-colors cursor-pointer',
                            imagePreview
                                ? 'border-brand-500/50 bg-brand-500/5'
                                : 'border-surface-600 bg-surface-800 hover:border-surface-500'
                        )}
                    >
                        {imagePreview ? (
                            <div className="relative w-full">
                                <img
                                    src={imagePreview}
                                    alt="Track layout preview"
                                    className="w-full rounded-xl object-contain max-h-48"
                                />
                                {uploading && (
                                    <div className="absolute inset-0 flex items-center justify-center bg-surface-900/70 rounded-xl">
                                        <div className="h-8 w-8 animate-spin rounded-full border-2 border-brand-500/30 border-t-brand-500" />
                                    </div>
                                )}
                            </div>
                        ) : (
                            <div className="py-8 text-center">
                                <svg className="mx-auto w-10 h-10 text-surface-500 mb-2" fill="none" viewBox="0 0 24 24" strokeWidth={1} stroke="currentColor">
                                    <path strokeLinecap="round" strokeLinejoin="round" d="m2.25 15.75 5.159-5.159a2.25 2.25 0 0 1 3.182 0l5.159 5.159m-1.5-1.5 1.409-1.409a2.25 2.25 0 0 1 3.182 0l2.909 2.909M3.75 21h16.5A2.25 2.25 0 0 0 22.5 18.75V5.25A2.25 2.25 0 0 0 20.25 3H3.75A2.25 2.25 0 0 0 1.5 5.25v13.5A2.25 2.25 0 0 0 3.75 21Z" />
                                </svg>
                                <p className="text-sm text-surface-400">Tap to upload track layout</p>
                                <p className="text-xs text-surface-500 mt-1">JPEG, PNG, WebP or SVG • Max 10MB</p>
                            </div>
                        )}
                    </div>
                    <input
                        ref={fileInputRef}
                        type="file"
                        accept="image/jpeg,image/png,image/webp,image/svg+xml"
                        onChange={handleImageUpload}
                        className="hidden"
                    />
                </div>

                {/* Track Name */}
                <div>
                    <label htmlFor="name" className="label">Track Name</label>
                    <input
                        id="name"
                        type="text"
                        value={name}
                        onChange={(e) => setName(e.target.value)}
                        className="input-field"
                        placeholder="e.g. Laguna Seca"
                        required
                    />
                </div>

                {/* Location */}
                <div>
                    <label htmlFor="location" className="label">Location</label>
                    <input
                        id="location"
                        type="text"
                        value={location}
                        onChange={(e) => setLocation(e.target.value)}
                        className="input-field"
                        placeholder="e.g. Monterey, CA"
                        required
                    />
                </div>

                {/* Description */}
                <div>
                    <label htmlFor="description" className="label">Description (optional)</label>
                    <textarea
                        id="description"
                        value={description}
                        onChange={(e) => setDescription(e.target.value)}
                        className="input-field min-h-[80px] resize-none"
                        placeholder="Tell drivers what this track is like..."
                        rows={3}
                    />
                </div>

                {/* Event Types */}
                <div>
                    <label className="label">Event Types</label>
                    <p className="text-xs text-surface-500 mb-2">What kind of events does this track host?</p>
                    <div className="space-y-2">
                        {EVENT_OPTIONS.map((option) => (
                            <button
                                key={option.value}
                                type="button"
                                onClick={() => toggleEventType(option.value)}
                                className={cn(
                                    'w-full flex items-center gap-3 rounded-lg border p-3 transition-colors text-left',
                                    eventTypes.includes(option.value)
                                        ? 'border-brand-500 bg-brand-500/10'
                                        : 'border-surface-600 bg-surface-800 hover:border-surface-500'
                                )}
                            >
                                <div
                                    className={cn(
                                        'flex h-5 w-5 items-center justify-center rounded border-2 transition-colors flex-shrink-0',
                                        eventTypes.includes(option.value)
                                            ? 'border-brand-500 bg-brand-500'
                                            : 'border-surface-500'
                                    )}
                                >
                                    {eventTypes.includes(option.value) && (
                                        <svg className="w-3 h-3 text-white" fill="none" viewBox="0 0 24 24" strokeWidth={3} stroke="currentColor">
                                            <path strokeLinecap="round" strokeLinejoin="round" d="m4.5 12.75 6 6 9-13.5" />
                                        </svg>
                                    )}
                                </div>
                                <div>
                                    <p className="font-medium text-surface-100">{option.label}</p>
                                    <p className="text-xs text-surface-500">{option.description}</p>
                                </div>
                            </button>
                        ))}
                    </div>
                </div>

                <button
                    type="submit"
                    disabled={loading || uploading}
                    className="btn-primary w-full"
                >
                    {loading ? (
                        <div className="h-5 w-5 animate-spin rounded-full border-2 border-white/30 border-t-white" />
                    ) : (
                        <>
                            <svg className="w-5 h-5" fill="none" viewBox="0 0 24 24" strokeWidth={2} stroke="currentColor">
                                <path strokeLinecap="round" strokeLinejoin="round" d="M3 16.5v2.25A2.25 2.25 0 0 0 5.25 21h13.5A2.25 2.25 0 0 0 21 18.75V16.5m-13.5-9L12 3m0 0 4.5 4.5M12 3v13.5" />
                            </svg>
                            Upload Track
                        </>
                    )}
                </button>
            </form>
        </div>
    );
}
