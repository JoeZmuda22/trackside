'use client';

import { useState, useEffect, useCallback } from 'react';
import Link from 'next/link';
import { cn, getEventLabel } from '@/lib/utils';
import { api } from '@/lib/api';
import { Badge } from '@/components/ui/badge';
import { StarRating } from '@/components/ui/star-rating';
import { LoadingSpinner, EmptyState } from '@/components/ui/loading';

interface Track {
    id: string;
    name: string;
    location: string;
    description: string | null;
    imageUrl: string | null;
    avgRating: number;
    events: { id: string; eventType: string }[];
    uploadedBy: { id: string; name: string | null };
    _count: { reviews: number; zones: number; lapRecords: number };
}

const EVENT_FILTERS = [
    { value: '', label: 'All' },
    { value: 'AUTOCROSS', label: 'Autocross' },
    { value: 'ROADCOURSE', label: 'Road Course' },
    { value: 'DRIFT', label: 'Drift' },
    { value: 'DRAG', label: 'Drag' },
];

export default function HomePage() {
    const [tracks, setTracks] = useState<Track[]>([]);
    const [loading, setLoading] = useState(true);
    const [search, setSearch] = useState('');
    const [eventFilter, setEventFilter] = useState('');

    const fetchTracks = useCallback(async () => {
        setLoading(true);
        const params = new URLSearchParams();
        if (search) params.set('search', search);
        if (eventFilter) params.set('eventType', eventFilter);

        try {
            const res = await api(`/api/tracks?${params}`);
            if (!res.ok) {
                const errorText = await res.text();
                console.error('API Error:', res.status, errorText);
                setTracks([]);
                return;
            }
            const data = await res.json();
            setTracks(Array.isArray(data) ? data : []);
        } catch (error) {
            console.error('Failed to fetch tracks:', error);
            setTracks([]);
        } finally {
            setLoading(false);
        }
    }, [search, eventFilter]);

    useEffect(() => {
        fetchTracks();
    }, [fetchTracks]);

    return (
        <div className="page-container">
            {/* Header */}
            <div className="flex items-center justify-between mb-6">
                <div>
                    <h1 className="text-2xl font-bold text-white">Tracks</h1>
                    <p className="text-sm text-surface-400">
                        {tracks.length} track{tracks.length !== 1 ? 's' : ''} available
                    </p>
                </div>
                <Link href="/tracks/new" className="btn-primary text-sm">
                    <svg className="w-4 h-4" fill="none" viewBox="0 0 24 24" strokeWidth={2} stroke="currentColor">
                        <path strokeLinecap="round" strokeLinejoin="round" d="M12 4.5v15m7.5-7.5h-15" />
                    </svg>
                    Add Track
                </Link>
            </div>

            {/* Search */}
            <div className="relative mb-4">
                <svg
                    className="absolute left-3 top-1/2 h-5 w-5 -translate-y-1/2 text-surface-500"
                    fill="none"
                    viewBox="0 0 24 24"
                    strokeWidth={1.5}
                    stroke="currentColor"
                >
                    <path strokeLinecap="round" strokeLinejoin="round" d="m21 21-5.197-5.197m0 0A7.5 7.5 0 1 0 5.196 5.196a7.5 7.5 0 0 0 10.607 10.607Z" />
                </svg>
                <input
                    type="text"
                    value={search}
                    onChange={(e) => setSearch(e.target.value)}
                    placeholder="Search tracks or locations..."
                    className="input-field pl-10"
                />
            </div>

            {/* Event Type Filter Chips */}
            <div className="flex gap-2 mb-6 overflow-x-auto pb-1 -mx-4 px-4">
                {EVENT_FILTERS.map((filter) => (
                    <button
                        key={filter.value}
                        onClick={() => setEventFilter(filter.value)}
                        className={cn(
                            'flex-shrink-0 rounded-full px-4 py-2 text-sm font-medium transition-colors',
                            eventFilter === filter.value
                                ? 'bg-brand-600 text-white'
                                : 'bg-surface-800 text-surface-300 hover:bg-surface-700'
                        )}
                    >
                        {filter.label}
                    </button>
                ))}
            </div>

            {/* Track List */}
            {loading ? (
                <LoadingSpinner />
            ) : tracks.length === 0 ? (
                <EmptyState
                    icon={
                        <svg className="w-12 h-12" fill="none" viewBox="0 0 24 24" strokeWidth={1} stroke="currentColor">
                            <path strokeLinecap="round" strokeLinejoin="round" d="M9 6.75V15m6-6v8.25m.503 3.498 4.875-2.437c.381-.19.622-.58.622-1.006V4.82c0-.836-.88-1.38-1.628-1.006l-3.869 1.934c-.317.159-.69.159-1.006 0L9.503 3.252a1.125 1.125 0 0 0-1.006 0L3.622 5.689C3.24 5.88 3 6.27 3 6.695V19.18c0 .836.88 1.38 1.628 1.006l3.869-1.934c.317-.159.69-.159 1.006 0l4.994 2.497c.317.158.69.158 1.006 0Z" />
                        </svg>
                    }
                    title="No tracks found"
                    description={search ? 'Try a different search term' : 'Be the first to upload a track!'}
                    action={
                        <Link href="/tracks/new" className="btn-primary text-sm">
                            Upload a Track
                        </Link>
                    }
                />
            ) : (
                <div className="space-y-3">
                    {tracks.map((track) => (
                        <Link key={track.id} href={`/tracks/${track.id}`}>
                            <div className="card hover:border-surface-600 transition-colors active:bg-surface-750">
                                {track.imageUrl && (
                                    <div className="relative -mx-4 -mt-4 mb-3 h-32 overflow-hidden rounded-t-xl">
                                        <img
                                            src={track.imageUrl}
                                            alt={track.name}
                                            className="h-full w-full object-cover"
                                        />
                                        <div className="absolute inset-0 bg-gradient-to-t from-surface-800/80 to-transparent" />
                                    </div>
                                )}

                                <div className="flex items-start justify-between">
                                    <div className="flex-1 min-w-0">
                                        <h3 className="font-semibold text-white truncate">{track.name}</h3>
                                        <p className="text-sm text-surface-400 flex items-center gap-1 mt-0.5">
                                            <svg className="w-3.5 h-3.5 flex-shrink-0" fill="none" viewBox="0 0 24 24" strokeWidth={1.5} stroke="currentColor">
                                                <path strokeLinecap="round" strokeLinejoin="round" d="M15 10.5a3 3 0 1 1-6 0 3 3 0 0 1 6 0Z" />
                                                <path strokeLinecap="round" strokeLinejoin="round" d="M19.5 10.5c0 7.142-7.5 11.25-7.5 11.25S4.5 17.642 4.5 10.5a7.5 7.5 0 1 1 15 0Z" />
                                            </svg>
                                            <span className="truncate">{track.location}</span>
                                        </p>
                                    </div>
                                    <div className="flex flex-col items-end gap-1 ml-3">
                                        <StarRating value={Math.round(track.avgRating)} readonly size="sm" />
                                        <span className="text-xs text-surface-500">
                                            {track._count.reviews} review{track._count.reviews !== 1 ? 's' : ''}
                                        </span>
                                    </div>
                                </div>

                                {/* Event Type Badges */}
                                <div className="flex flex-wrap gap-1.5 mt-3">
                                    {track.events.map((event) => (
                                        <Badge key={event.id} variant="event" value={event.eventType}>
                                            {getEventLabel(event.eventType)}
                                        </Badge>
                                    ))}
                                </div>

                                {/* Stats */}
                                <div className="flex items-center gap-4 mt-3 pt-3 border-t border-surface-700">
                                    <span className="text-xs text-surface-500 flex items-center gap-1">
                                        <svg className="w-3.5 h-3.5" fill="none" viewBox="0 0 24 24" strokeWidth={1.5} stroke="currentColor">
                                            <path strokeLinecap="round" strokeLinejoin="round" d="M7.5 8.25h9m-9 3H12m-9.75 1.51c0 1.6 1.123 2.994 2.707 3.227 1.087.16 2.185.283 3.293.369V21l4.076-4.076a1.526 1.526 0 0 1 1.037-.443 48.282 48.282 0 0 0 5.68-.494c1.584-.233 2.707-1.626 2.707-3.228V6.741c0-1.602-1.123-2.995-2.707-3.228A48.394 48.394 0 0 0 12 3c-2.392 0-4.744.175-7.043.513C3.373 3.746 2.25 5.14 2.25 6.741v6.018Z" />
                                        </svg>
                                        {track._count.zones} zone{track._count.zones !== 1 ? 's' : ''}
                                    </span>
                                    <span className="text-xs text-surface-500 flex items-center gap-1">
                                        <svg className="w-3.5 h-3.5" fill="none" viewBox="0 0 24 24" strokeWidth={1.5} stroke="currentColor">
                                            <path strokeLinecap="round" strokeLinejoin="round" d="M12 6v6h4.5m4.5 0a9 9 0 1 1-18 0 9 9 0 0 1 18 0Z" />
                                        </svg>
                                        {track._count.lapRecords} lap{track._count.lapRecords !== 1 ? 's' : ''}
                                    </span>
                                    <span className="text-xs text-surface-500 ml-auto">
                                        by {track.uploadedBy.name || 'Unknown'}
                                    </span>
                                </div>
                            </div>
                        </Link>
                    ))}
                </div>
            )}
        </div>
    );
}
