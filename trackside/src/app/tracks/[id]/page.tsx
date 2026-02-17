'use client';

import { useState, useEffect, use, useRef } from 'react';
import { useAuth } from '@/components/providers/auth-provider';
import { api } from '@/lib/api';
import Link from 'next/link';
import { cn, getEventLabel, getExperienceLabel } from '@/lib/utils';
import { Badge } from '@/components/ui/badge';
import { StarRating } from '@/components/ui/star-rating';
import { BottomSheet } from '@/components/ui/bottom-sheet';
import { LoadingSpinner, EmptyState } from '@/components/ui/loading';

interface TrackZone {
    id: string;
    name: string;
    description: string | null;
    posX: number;
    posY: number;
    tips: {
        id: string;
        content: string;
        conditions: string | null;
        createdAt: string;
        author: { id: string; name: string | null };
    }[];
}

interface TrackReview {
    id: string;
    rating: number;
    content: string | null;
    conditions: string;
    createdAt: string;
    author: {
        id: string;
        name: string | null;
        experience: string;
        cars: { make: string; model: string; year: number }[];
    };
    trackEvent: { eventType: string } | null;
}

interface Track {
    id: string;
    name: string;
    location: string;
    description: string | null;
    imageUrl: string | null;
    avgRating: number;
    uploadedBy: { id: string; name: string | null; experience: string };
    events: { id: string; eventType: string }[];
    zones: TrackZone[];
    reviews: TrackReview[];
    _count: { reviews: number; zones: number; lapRecords: number };
}

type Tab = 'map' | 'reviews' | 'info';

export default function TrackDetailPage({ params }: { params: Promise<{ id: string }> }) {
    const { id } = use(params);
    const { user } = useAuth();
    const [track, setTrack] = useState<Track | null>(null);
    const [loading, setLoading] = useState(true);
    const [activeTab, setActiveTab] = useState<Tab>('map');
    const [selectedZone, setSelectedZone] = useState<TrackZone | null>(null);
    const [showAddZone, setShowAddZone] = useState(false);
    const [showAddReview, setShowAddReview] = useState(false);
    const [showAddTip, setShowAddTip] = useState(false);
    const [eventFilter, setEventFilter] = useState('');

    // Zone form state
    const [newZoneName, setNewZoneName] = useState('');
    const [newZoneDesc, setNewZoneDesc] = useState('');
    const [newZonePos, setNewZonePos] = useState<{ x: number; y: number } | null>(null);
    const [newZoneEventType, setNewZoneEventType] = useState('');

    // Review form state
    const [reviewRating, setReviewRating] = useState(0);
    const [reviewContent, setReviewContent] = useState('');
    const [reviewConditions, setReviewConditions] = useState<'DRY' | 'WET'>('DRY');
    const [reviewEventId, setReviewEventId] = useState('');

    // Tip form state
    const [tipContent, setTipContent] = useState('');
    const [tipConditions, setTipConditions] = useState<'DRY' | 'WET' | ''>('');

    // Image upload
    const imageInputRef = useRef<HTMLInputElement>(null);
    const [uploadingImage, setUploadingImage] = useState(false);
    const [images, setImages] = useState<any[]>([]);
    const [imageUrl, setImageUrl] = useState('');
    const [imageCaption, setImageCaption] = useState('');
    const [showImageUpload, setShowImageUpload] = useState(false);

    // Zone editing
    const [editingZoneId, setEditingZoneId] = useState<string | null>(null);
    const [editZoneName, setEditZoneName] = useState('');
    const [editZoneDesc, setEditZoneDesc] = useState('');

    const [submitting, setSubmitting] = useState(false);

    useEffect(() => {
        fetchTrack();
        fetchImages();
    }, [id, eventFilter]);

    const fetchImages = async () => {
        try {
            const res = await api(`/api/tracks/${id}/images`, { cache: 'no-store' });
            if (res.ok) {
                const data = await res.json();
                setImages(data);
            }
        } catch (error) {
            console.error('Failed to fetch images:', error);
        }
    };

    const fetchTrack = async () => {
        setLoading(true);
        try {
            const url = new URL(`/api/tracks/${id}`, window.location.origin);
            if (eventFilter) {
                url.searchParams.set('eventType', eventFilter);
            }
            const res = await api(url.pathname + url.search, { cache: 'no-store' });
            if (res.ok) {
                const data = await res.json();
                setTrack(data);
            }
        } catch (error) {
            console.error('Failed to fetch track:', error);
        } finally {
            setLoading(false);
        }
    };

    const handleImageClick = (e: React.MouseEvent<HTMLDivElement>) => {
        if (!showAddZone) return;
        const rect = e.currentTarget.getBoundingClientRect();
        const x = ((e.clientX - rect.left) / rect.width) * 100;
        const y = ((e.clientY - rect.top) / rect.height) * 100;
        setNewZonePos({ x, y });
    };

    const handleAddZone = async (e: React.FormEvent) => {
        e.preventDefault();
        if (!newZonePos) return;
        setSubmitting(true);

        try {
            const res = await api(`/api/tracks/${id}/zones`, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({
                    name: newZoneName,
                    description: newZoneDesc || undefined,
                    posX: newZonePos.x,
                    posY: newZonePos.y,
                    eventType: newZoneEventType || undefined,
                }),
            });

            if (res.ok) {
                setShowAddZone(false);
                setNewZoneName('');
                setNewZoneDesc('');
                setNewZonePos(null);
                setNewZoneEventType('');
                fetchTrack();
            }
        } catch (error) {
            console.error('Failed to add zone:', error);
        } finally {
            setSubmitting(false);
        }
    };

    const handleAddReview = async (e: React.FormEvent) => {
        e.preventDefault();
        if (reviewRating === 0) return;
        setSubmitting(true);

        try {
            const res = await api(`/api/tracks/${id}/reviews`, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({
                    rating: reviewRating,
                    content: reviewContent || undefined,
                    conditions: reviewConditions,
                    trackEventId: reviewEventId || undefined,
                }),
            });

            if (res.ok) {
                setShowAddReview(false);
                setReviewRating(0);
                setReviewContent('');
                setReviewConditions('DRY');
                setReviewEventId('');
                fetchTrack();
            }
        } catch (error) {
            console.error('Failed to add review:', error);
        } finally {
            setSubmitting(false);
        }
    };

    const handleAddTip = async (e: React.FormEvent) => {
        e.preventDefault();
        if (!selectedZone) return;
        setSubmitting(true);

        try {
            const res = await api(`/api/tracks/${id}/zones/${selectedZone.id}/tips`, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({
                    content: tipContent,
                    conditions: tipConditions || undefined,
                }),
            });

            if (res.ok) {
                setShowAddTip(false);
                setTipContent('');
                setTipConditions('');
                fetchTrack();
                // Refresh selected zone data
                const updatedTrack = await (await api(`/api/tracks/${id}`)).json();
                setTrack(updatedTrack);
                const updatedZone = updatedTrack.zones.find((z: TrackZone) => z.id === selectedZone.id);
                if (updatedZone) setSelectedZone(updatedZone);
            }
        } catch (error) {
            console.error('Failed to add tip:', error);
        } finally {
            setSubmitting(false);
        }
    };

    const handleAddImage = async (e: React.FormEvent) => {
        e.preventDefault();
        if (!imageUrl.trim()) return;
        setUploadingImage(true);

        try {
            const res = await api(`/api/tracks/${id}/images`, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({
                    url: imageUrl,
                    caption: imageCaption || undefined,
                }),
            });

            if (res.ok) {
                setImageUrl('');
                setImageCaption('');
                setShowImageUpload(false);
                fetchImages();
            }
        } catch (error) {
            console.error('Failed to add image:', error);
        } finally {
            setUploadingImage(false);
        }
    };

    const handleDeleteImage = async (imageId: string) => {
        try {
            const res = await api(`/api/tracks/${id}/images?imageId=${imageId}`, {
                method: 'DELETE',
            });

            if (res.ok) {
                fetchImages();
            }
        } catch (error) {
            console.error('Failed to delete image:', error);
        }
    };

    const handleEditZone = async (zoneId: string) => {
        if (!editZoneName.trim()) return;
        setSubmitting(true);

        try {
            const res = await api(`/api/tracks/${id}/zones/${zoneId}`, {
                method: 'PATCH',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({
                    name: editZoneName,
                    description: editZoneDesc || undefined,
                }),
            });

            if (res.ok) {
                setEditingZoneId(null);
                setEditZoneName('');
                setEditZoneDesc('');
                fetchTrack();
                // Update selected zone if it's the one being edited
                if (selectedZone?.id === zoneId) {
                    const updatedZone = track?.zones.find((z) => z.id === zoneId);
                    if (updatedZone) setSelectedZone(updatedZone);
                }
            }
        } catch (error) {
            console.error('Failed to update zone:', error);
        } finally {
            setSubmitting(false);
        }
    };

    const isOwner = user?.id === track?.uploadedBy?.id;

    if (loading) return <div className="page-container"><LoadingSpinner /></div>;
    if (!track) {
        return (
            <div className="page-container">
                <EmptyState title="Track not found" description="This track doesn't exist or was removed." />
            </div>
        );
    }

    const filteredReviews = eventFilter
        ? track.reviews.filter((r) => r.trackEvent?.eventType === eventFilter)
        : track.reviews;

    return (
        <div className="pb-20">
            {/* Header Image */}
            <div className="relative">
                {track.imageUrl ? (
                    <div className="relative h-48 w-full">
                        <img src={track.imageUrl} alt={track.name} className="h-full w-full object-cover" />
                        <div className="absolute inset-0 bg-gradient-to-t from-surface-900 via-surface-900/30 to-transparent" />
                    </div>
                ) : (
                    <div className="h-32 bg-surface-800" />
                )}

                {/* Back button */}
                <Link
                    href="/"
                    className="absolute top-4 left-4 flex h-9 w-9 items-center justify-center rounded-full bg-surface-900/70 backdrop-blur-sm text-white"
                >
                    <svg className="w-5 h-5" fill="none" viewBox="0 0 24 24" strokeWidth={2} stroke="currentColor">
                        <path strokeLinecap="round" strokeLinejoin="round" d="M15.75 19.5 8.25 12l7.5-7.5" />
                    </svg>
                </Link>
            </div>

            {/* Track Info */}
            <div className="px-4 -mt-6 relative z-10">
                <div className="flex items-end justify-between">
                    <div>
                        <h1 className="text-2xl font-bold text-white">{track.name}</h1>
                        <p className="text-sm text-surface-400 flex items-center gap-1 mt-1">
                            <svg className="w-3.5 h-3.5" fill="none" viewBox="0 0 24 24" strokeWidth={1.5} stroke="currentColor">
                                <path strokeLinecap="round" strokeLinejoin="round" d="M15 10.5a3 3 0 1 1-6 0 3 3 0 0 1 6 0Z" />
                                <path strokeLinecap="round" strokeLinejoin="round" d="M19.5 10.5c0 7.142-7.5 11.25-7.5 11.25S4.5 17.642 4.5 10.5a7.5 7.5 0 1 1 15 0Z" />
                            </svg>
                            {track.location}
                        </p>
                    </div>
                    <div className="flex flex-col items-end">
                        <StarRating value={Math.round(track.avgRating)} readonly size="sm" />
                        <span className="text-xs text-surface-500 mt-0.5">
                            {track.avgRating.toFixed(1)} ({track._count.reviews})
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
            </div>

            {/* Event Type Filter — only show if track has multiple events */}
            {track.events.length > 1 && (
                <div className="px-4 pb-3 border-b border-surface-700">
                    <p className="text-xs font-medium text-surface-400 mb-2">Filter by event type:</p>
                    <div className="flex flex-wrap gap-2">
                        <button
                            onClick={() => setEventFilter('')}
                            className={cn(
                                'px-3 py-1 rounded-full text-xs font-medium transition-colors',
                                !eventFilter
                                    ? 'bg-brand-600 text-white'
                                    : 'bg-surface-800 text-surface-300 hover:bg-surface-700'
                            )}
                        >
                            All Events
                        </button>
                        {track.events.map((event) => (
                            <button
                                key={event.id}
                                onClick={() => setEventFilter(event.eventType)}
                                className={cn(
                                    'px-3 py-1 rounded-full text-xs font-medium transition-colors',
                                    eventFilter === event.eventType
                                        ? 'bg-brand-600 text-white'
                                        : 'bg-surface-800 text-surface-300 hover:bg-surface-700'
                                )}
                            >
                                {getEventLabel(event.eventType)}
                            </button>
                        ))}
                    </div>
                </div>
            )}

            {/* Tabs */}
            <div className="flex border-b border-surface-700 mt-4 px-4">
                {(['map', 'reviews', 'info'] as Tab[]).map((tab) => (
                    <button
                        key={tab}
                        onClick={() => setActiveTab(tab)}
                        className={cn(
                            'flex-1 py-3 text-sm font-medium border-b-2 transition-colors capitalize',
                            activeTab === tab
                                ? 'border-brand-500 text-brand-500'
                                : 'border-transparent text-surface-400 hover:text-surface-300'
                        )}
                    >
                        {tab === 'map' ? `Zones (${track.zones.length})${eventFilter ? ` - ${getEventLabel(eventFilter)}` : ''}` : tab === 'reviews' ? `Reviews (${track._count.reviews})` : 'Info'}
                    </button>
                ))}
            </div>

            <div className="px-4 pt-4">
                {/* MAP TAB — Interactive Zone Map */}
                {activeTab === 'map' && (
                    <div>
                        {track.imageUrl ? (
                            <div className="relative rounded-xl overflow-hidden border border-surface-700">
                                <div
                                    className={cn(
                                        'relative',
                                        showAddZone && 'cursor-crosshair'
                                    )}
                                    onClick={handleImageClick}
                                >
                                    <img
                                        src={track.imageUrl}
                                        alt={track.name}
                                        className="w-full object-contain"
                                    />
                                    {/* Zone pins */}
                                    {track.zones.map((zone) => (
                                        <button
                                            key={zone.id}
                                            onClick={(e) => {
                                                e.stopPropagation();
                                                setSelectedZone(zone);
                                            }}
                                            className="absolute -translate-x-1/2 -translate-y-1/2 group"
                                            style={{ left: `${zone.posX}%`, top: `${zone.posY}%` }}
                                        >
                                            <div className="flex h-7 w-7 items-center justify-center rounded-full bg-brand-600 text-white text-xs font-bold shadow-lg ring-2 ring-white/20 transition-transform active:scale-110">
                                                {zone.tips.length}
                                            </div>
                                            <div className="absolute -bottom-6 left-1/2 -translate-x-1/2 whitespace-nowrap text-[10px] font-medium text-surface-300 bg-surface-900/80 px-1.5 py-0.5 rounded opacity-0 group-hover:opacity-100 transition-opacity">
                                                {zone.name}
                                            </div>
                                        </button>
                                    ))}

                                    {/* New zone pin preview */}
                                    {showAddZone && newZonePos && (
                                        <div
                                            className="absolute -translate-x-1/2 -translate-y-1/2 pointer-events-none"
                                            style={{ left: `${newZonePos.x}%`, top: `${newZonePos.y}%` }}
                                        >
                                            <div className="flex h-7 w-7 items-center justify-center rounded-full bg-green-500 text-white text-xs font-bold shadow-lg ring-2 ring-green-300/50 animate-pulse">
                                                +
                                            </div>
                                        </div>
                                    )}
                                </div>

                                {/* Change image button — only for track owner */}
                                {isOwner && (
                                    <button
                                        onClick={() => imageInputRef.current?.click()}
                                        disabled={uploadingImage}
                                        className="absolute top-2 right-2 flex items-center gap-1.5 rounded-full bg-surface-900/70 backdrop-blur-sm px-3 py-1.5 text-xs font-medium text-white hover:bg-surface-900/90 transition-colors"
                                    >
                                        <svg className="w-3.5 h-3.5" fill="none" viewBox="0 0 24 24" strokeWidth={2} stroke="currentColor">
                                            <path strokeLinecap="round" strokeLinejoin="round" d="M6.827 6.175A2.31 2.31 0 0 1 5.186 7.23c-.38.054-.757.112-1.134.175C2.999 7.58 2.25 8.507 2.25 9.574V18a2.25 2.25 0 0 0 2.25 2.25h15A2.25 2.25 0 0 0 21.75 18V9.574c0-1.067-.75-1.994-1.802-2.169a47.865 47.865 0 0 0-1.134-.175 2.31 2.31 0 0 1-1.64-1.055l-.822-1.316a2.192 2.192 0 0 0-1.736-1.039 48.774 48.774 0 0 0-5.232 0 2.192 2.192 0 0 0-1.736 1.039l-.821 1.316Z" />
                                            <path strokeLinecap="round" strokeLinejoin="round" d="M16.5 12.75a4.5 4.5 0 1 1-9 0 4.5 4.5 0 0 1 9 0Z" />
                                        </svg>
                                        {uploadingImage ? 'Uploading...' : 'Change Image'}
                                    </button>
                                )}
                            </div>
                        ) : (
                            <div className="card text-center py-10">
                                <svg className="w-12 h-12 text-surface-600 mx-auto mb-3" fill="none" viewBox="0 0 24 24" strokeWidth={1} stroke="currentColor">
                                    <path strokeLinecap="round" strokeLinejoin="round" d="m2.25 15.75 5.159-5.159a2.25 2.25 0 0 1 3.182 0l5.159 5.159m-1.5-1.5 1.409-1.409a2.25 2.25 0 0 1 3.182 0l2.909 2.909M3.75 21h16.5A2.25 2.25 0 0 0 22.5 19.5V4.5A2.25 2.25 0 0 0 20.25 3H3.75A2.25 2.25 0 0 0 1.5 4.5v15A2.25 2.25 0 0 0 3.75 21Z" />
                                </svg>
                                <p className="text-surface-400 text-sm mb-1">No track layout image yet</p>
                                {isOwner ? (
                                    <button
                                        onClick={() => imageInputRef.current?.click()}
                                        disabled={uploadingImage}
                                        className="btn-primary text-sm mt-3"
                                    >
                                        {uploadingImage ? 'Uploading...' : '📸 Upload Track Image'}
                                    </button>
                                ) : (
                                    <p className="text-surface-500 text-xs">The track owner can upload an image</p>
                                )}
                            </div>
                        )}

                        {/* Action buttons */}
                        <div className="flex gap-2 mt-3">
                            {user && (
                                <button
                                    onClick={() => {
                                        setShowAddZone(!showAddZone);
                                        setNewZonePos(null);
                                    }}
                                    className={cn(
                                        'flex-1 text-sm',
                                        showAddZone ? 'btn-danger' : 'btn-secondary'
                                    )}
                                >
                                    {showAddZone ? 'Cancel' : '+ Add Zone'}
                                </button>
                            )}
                            <button
                                onClick={() => setShowImageUpload(!showImageUpload)}
                                className="btn-secondary text-sm px-3"
                                title="Add track photos"
                            >
                                📸
                            </button>
                        </div>

                        {/* Add zone form */}
                        {showAddZone && (
                            <form onSubmit={handleAddZone} className="card mt-3 space-y-3">
                                <p className="text-sm text-surface-300">
                                    {newZonePos
                                        ? '✓ Pin placed! Name this zone:'
                                        : 'Tap on the track image to place a zone pin'}
                                </p>
                                {newZonePos && (
                                    <>
                                        <input
                                            type="text"
                                            value={newZoneName}
                                            onChange={(e) => setNewZoneName(e.target.value)}
                                            placeholder="Zone name (e.g. Turn 3)"
                                            className="input-field"
                                            required
                                        />
                                        <textarea
                                            value={newZoneDesc}
                                            onChange={(e) => setNewZoneDesc(e.target.value)}
                                            placeholder="Description (optional)"
                                            className="input-field min-h-[60px] resize-none"
                                        />
                                        {track.events.length > 1 && (
                                            <select
                                                value={newZoneEventType}
                                                onChange={(e) => setNewZoneEventType(e.target.value)}
                                                className="input-field"
                                            >
                                                <option value="">All Event Types</option>
                                                {track.events.map((event) => (
                                                    <option key={event.id} value={event.eventType}>
                                                        {getEventLabel(event.eventType)}
                                                    </option>
                                                ))}
                                            </select>
                                        )}
                                        <button type="submit" disabled={submitting} className="btn-primary w-full">
                                            {submitting ? 'Adding...' : 'Add Zone'}
                                        </button>
                                    </>
                                )}
                            </form>
                        )}

                        {/* Image upload form */}
                        {showImageUpload && (
                            <form onSubmit={handleAddImage} className="card mt-3 space-y-3">
                                <h3 className="section-title">Add Track Photo</h3>
                                <input
                                    type="url"
                                    value={imageUrl}
                                    onChange={(e) => setImageUrl(e.target.value)}
                                    placeholder="Paste image URL (https://...)"
                                    className="input-field"
                                    required
                                />
                                <input
                                    type="text"
                                    value={imageCaption}
                                    onChange={(e) => setImageCaption(e.target.value)}
                                    placeholder="Caption (optional)"
                                    className="input-field"
                                />
                                <div className="flex gap-2">
                                    <button type="submit" disabled={uploadingImage || !imageUrl.trim()} className="flex-1 btn-primary text-sm">
                                        {uploadingImage ? 'Adding...' : 'Add Photo'}
                                    </button>
                                    <button
                                        type="button"
                                        onClick={() => setShowImageUpload(false)}
                                        className="flex-1 btn-secondary text-sm"
                                    >
                                        Cancel
                                    </button>
                                </div>
                            </form>
                        )}

                        {/* Image gallery */}
                        {images.length > 0 && (
                            <div className="card mt-3">
                                <h3 className="section-title mb-3">Track Photos ({images.length})</h3>
                                <div className="grid grid-cols-2 gap-2">
                                    {images.map((img) => (
                                        <div key={img.id} className="relative group">
                                            <img
                                                src={img.url}
                                                alt={img.caption || 'Track photo'}
                                                className="w-full h-24 object-cover rounded-lg"
                                            />
                                            {img.caption && (
                                                <p className="text-xs text-surface-300 mt-1 truncate">{img.caption}</p>
                                            )}
                                            {user && (
                                                <button
                                                    onClick={() => handleDeleteImage(img.id)}
                                                    className="absolute -top-2 -right-2 p-1 bg-red-600 text-white rounded-full opacity-0 group-hover:opacity-100 transition-opacity"
                                                    title="Delete photo"
                                                >
                                                    <svg className="w-3 h-3" fill="none" viewBox="0 0 24 24" strokeWidth={3} stroke="currentColor">
                                                        <path strokeLinecap="round" strokeLinejoin="round" d="M6 18L18 6M6 6l12 12" />
                                                    </svg>
                                                </button>
                                            )}
                                        </div>
                                    ))}
                                </div>
                            </div>
                        )}

                        {/* Zone list */}
                        <div className="mt-4 space-y-2">
                            <h3 className="section-title">All Zones</h3>
                            {track.zones.length === 0 ? (
                                <p className="text-sm text-surface-500">No zones added yet. Be the first!</p>
                            ) : (
                                track.zones.map((zone) => (
                                    <button
                                        key={zone.id}
                                        onClick={() => setSelectedZone(zone)}
                                        className="card w-full text-left hover:border-surface-600 transition-colors"
                                    >
                                        <div className="flex items-center justify-between">
                                            <div>
                                                <h4 className="font-medium text-surface-100">{zone.name}</h4>
                                                {zone.description && (
                                                    <p className="text-xs text-surface-500 mt-0.5">{zone.description}</p>
                                                )}
                                            </div>
                                            <div className="flex items-center gap-1 text-surface-400">
                                                <svg className="w-4 h-4" fill="none" viewBox="0 0 24 24" strokeWidth={1.5} stroke="currentColor">
                                                    <path strokeLinecap="round" strokeLinejoin="round" d="M7.5 8.25h9m-9 3H12m-9.75 1.51c0 1.6 1.123 2.994 2.707 3.227 1.087.16 2.185.283 3.293.369V21l4.076-4.076a1.526 1.526 0 0 1 1.037-.443 48.282 48.282 0 0 0 5.68-.494c1.584-.233 2.707-1.626 2.707-3.228V6.741c0-1.602-1.123-2.995-2.707-3.228A48.394 48.394 0 0 0 12 3c-2.392 0-4.744.175-7.043.513C3.373 3.746 2.25 5.14 2.25 6.741v6.018Z" />
                                                </svg>
                                                <span className="text-xs">{zone.tips.length}</span>
                                            </div>
                                        </div>
                                    </button>
                                ))
                            )}
                        </div>
                    </div>
                )}

                {/* REVIEWS TAB */}
                {activeTab === 'reviews' && (
                    <div>
                        {/* Event filter */}
                        <div className="flex gap-2 mb-4 overflow-x-auto pb-1">
                            <button
                                onClick={() => setEventFilter('')}
                                className={cn(
                                    'flex-shrink-0 rounded-full px-3 py-1.5 text-xs font-medium transition-colors',
                                    !eventFilter ? 'bg-brand-600 text-white' : 'bg-surface-800 text-surface-300'
                                )}
                            >
                                All
                            </button>
                            {track.events.map((event) => (
                                <button
                                    key={event.id}
                                    onClick={() => setEventFilter(event.eventType)}
                                    className={cn(
                                        'flex-shrink-0 rounded-full px-3 py-1.5 text-xs font-medium transition-colors',
                                        eventFilter === event.eventType
                                            ? 'bg-brand-600 text-white'
                                            : 'bg-surface-800 text-surface-300'
                                    )}
                                >
                                    {getEventLabel(event.eventType)}
                                </button>
                            ))}
                        </div>

                        {user && (
                            <button
                                onClick={() => setShowAddReview(true)}
                                className="btn-primary w-full mb-4 text-sm"
                            >
                                Write a Review
                            </button>
                        )}

                        {filteredReviews.length === 0 ? (
                            <EmptyState
                                title="No reviews yet"
                                description="Be the first to review this track!"
                            />
                        ) : (
                            <div className="space-y-3">
                                {filteredReviews.map((review) => (
                                    <div key={review.id} className="card">
                                        <div className="flex items-start justify-between">
                                            <div>
                                                <p className="font-medium text-surface-100">{review.author.name}</p>
                                                <p className="text-xs text-surface-500">
                                                    {getExperienceLabel(review.author.experience)}
                                                    {review.author.cars.length > 0 && (
                                                        <> • {review.author.cars[0].year} {review.author.cars[0].make} {review.author.cars[0].model}</>
                                                    )}
                                                </p>
                                            </div>
                                            <StarRating value={review.rating} readonly size="sm" />
                                        </div>
                                        <div className="flex gap-1.5 mt-2">
                                            <Badge variant="condition" value={review.conditions}>
                                                {review.conditions === 'WET' ? '🌧 Wet' : '☀️ Dry'}
                                            </Badge>
                                            {review.trackEvent && (
                                                <Badge variant="event" value={review.trackEvent.eventType}>
                                                    {getEventLabel(review.trackEvent.eventType)}
                                                </Badge>
                                            )}
                                        </div>
                                        {review.content && (
                                            <p className="mt-2 text-sm text-surface-300">{review.content}</p>
                                        )}
                                        <p className="mt-2 text-xs text-surface-600">
                                            {new Date(review.createdAt).toLocaleDateString()}
                                        </p>
                                    </div>
                                ))}
                            </div>
                        )}
                    </div>
                )}

                {/* INFO TAB */}
                {activeTab === 'info' && (
                    <div className="space-y-4">
                        {track.description && (
                            <div className="card">
                                <h3 className="section-title">About</h3>
                                <p className="text-sm text-surface-300">{track.description}</p>
                            </div>
                        )}
                        <div className="card">
                            <h3 className="section-title">Details</h3>
                            <dl className="space-y-2 text-sm">
                                <div className="flex justify-between">
                                    <dt className="text-surface-400">Uploaded by</dt>
                                    <dd className="text-surface-200">{track.uploadedBy.name}</dd>
                                </div>
                                <div className="flex justify-between">
                                    <dt className="text-surface-400">Total Reviews</dt>
                                    <dd className="text-surface-200">{track._count.reviews}</dd>
                                </div>
                                <div className="flex justify-between">
                                    <dt className="text-surface-400">Zones Mapped</dt>
                                    <dd className="text-surface-200">{track._count.zones}</dd>
                                </div>
                                <div className="flex justify-between">
                                    <dt className="text-surface-400">Lap Records</dt>
                                    <dd className="text-surface-200">{track._count.lapRecords}</dd>
                                </div>
                            </dl>
                        </div>

                        {/* Log a lap link */}
                        <Link
                            href={`/lapbook/new?trackId=${track.id}`}
                            className="btn-secondary w-full text-sm"
                        >
                            <svg className="w-4 h-4" fill="none" viewBox="0 0 24 24" strokeWidth={2} stroke="currentColor">
                                <path strokeLinecap="round" strokeLinejoin="round" d="M12 6v6h4.5m4.5 0a9 9 0 1 1-18 0 9 9 0 0 1 18 0Z" />
                            </svg>
                            Log a Lap Time
                        </Link>
                    </div>
                )}
            </div>

            {/* Zone Tips Bottom Sheet */}
            <BottomSheet
                isOpen={!!selectedZone && !showAddTip}
                onClose={() => setSelectedZone(null)}
                title={selectedZone?.name}
            >
                {selectedZone && (
                    <div>
                        {/* Zone name editor */}
                        {editingZoneId === selectedZone.id ? (
                            <div className="mb-4 p-3 bg-surface-700/30 rounded-lg space-y-2">
                                <input
                                    type="text"
                                    value={editZoneName}
                                    onChange={(e) => setEditZoneName(e.target.value)}
                                    placeholder="Zone name"
                                    className="input-field"
                                    autoFocus
                                />
                                <textarea
                                    value={editZoneDesc}
                                    onChange={(e) => setEditZoneDesc(e.target.value)}
                                    placeholder="Description (optional)"
                                    className="input-field min-h-[60px] resize-none"
                                />
                                <div className="flex gap-2">
                                    <button
                                        onClick={() => handleEditZone(selectedZone.id)}
                                        disabled={submitting || !editZoneName.trim()}
                                        className="flex-1 btn-primary text-sm"
                                    >
                                        {submitting ? 'Saving...' : 'Save'}
                                    </button>
                                    <button
                                        onClick={() => setEditingZoneId(null)}
                                        className="flex-1 btn-secondary text-sm"
                                    >
                                        Cancel
                                    </button>
                                </div>
                            </div>
                        ) : (
                            <div className="mb-4 flex items-start justify-between">
                                <div className="flex-1">
                                    {selectedZone.description && (
                                        <p className="text-sm text-surface-400 mb-2">{selectedZone.description}</p>
                                    )}
                                </div>
                                {user && (
                                    <button
                                        onClick={() => {
                                            setEditingZoneId(selectedZone.id);
                                            setEditZoneName(selectedZone.name);
                                            setEditZoneDesc(selectedZone.description || '');
                                        }}
                                        className="ml-2 p-2 text-surface-400 hover:text-surface-200 transition-colors"
                                        title="Edit zone name"
                                    >
                                        <svg className="w-4 h-4" fill="none" viewBox="0 0 24 24" strokeWidth={2} stroke="currentColor">
                                            <path strokeLinecap="round" strokeLinejoin="round" d="m16.862 4.487l1.687-1.688a1.875 1.875 0 112.652 2.652L10.582 16.07a4.5 4.5 0 01-1.897 1.13L6 18l.8-2.685a4.5 4.5 0 011.13-1.897l8.932-8.931zm0 0L19.5 7.125M18 14v4.75A2.25 2.25 0 0115.75 21H5.25A2.25 2.25 0 013 18.75V8.25A2.25 2.25 0 015.25 6H10" />
                                        </svg>
                                    </button>
                                )}
                            </div>
                        )}

                        {user && (
                            <button
                                onClick={() => setShowAddTip(true)}
                                className="btn-primary w-full text-sm mb-4"
                            >
                                + Add Tip
                            </button>
                        )}

                        {selectedZone.tips.length === 0 ? (
                            <p className="text-sm text-surface-500 text-center py-4">
                                No tips yet. Share your knowledge!
                            </p>
                        ) : (
                            <div className="space-y-3">
                                {selectedZone.tips.map((tip) => (
                                    <div key={tip.id} className="rounded-lg bg-surface-700/50 p-3">
                                        <div className="flex items-center justify-between mb-1">
                                            <span className="text-xs font-medium text-surface-300">
                                                {tip.author.name}
                                            </span>
                                            {tip.conditions && (
                                                <Badge variant="condition" value={tip.conditions}>
                                                    {tip.conditions === 'WET' ? '🌧 Wet' : '☀️ Dry'}
                                                </Badge>
                                            )}
                                        </div>
                                        <p className="text-sm text-surface-200">{tip.content}</p>
                                        <p className="text-xs text-surface-600 mt-1">
                                            {new Date(tip.createdAt).toLocaleDateString()}
                                        </p>
                                    </div>
                                ))}
                            </div>
                        )}
                    </div>
                )}
            </BottomSheet>

            {/* Add Tip Bottom Sheet */}
            <BottomSheet
                isOpen={showAddTip}
                onClose={() => setShowAddTip(false)}
                title={`Add Tip for ${selectedZone?.name}`}
            >
                <form onSubmit={handleAddTip} className="space-y-3">
                    <div>
                        <label className="label">Your Tip</label>
                        <textarea
                            value={tipContent}
                            onChange={(e) => setTipContent(e.target.value)}
                            placeholder="Share your advice for this section..."
                            className="input-field min-h-[80px] resize-none"
                            rows={3}
                            required
                        />
                    </div>
                    <div>
                        <label className="label">Conditions (optional)</label>
                        <div className="flex gap-2">
                            {(['', 'DRY', 'WET'] as const).map((cond) => (
                                <button
                                    key={cond}
                                    type="button"
                                    onClick={() => setTipConditions(cond)}
                                    className={cn(
                                        'flex-1 rounded-lg px-3 py-2.5 text-sm font-medium transition-colors',
                                        tipConditions === cond
                                            ? 'bg-brand-600 text-white'
                                            : 'bg-surface-700 text-surface-300'
                                    )}
                                >
                                    {cond === '' ? 'Any' : cond === 'DRY' ? '☀️ Dry' : '🌧 Wet'}
                                </button>
                            ))}
                        </div>
                    </div>
                    <button type="submit" disabled={submitting} className="btn-primary w-full">
                        {submitting ? 'Posting...' : 'Post Tip'}
                    </button>
                </form>
            </BottomSheet>

            {/* Add Review Bottom Sheet */}
            <BottomSheet
                isOpen={showAddReview}
                onClose={() => setShowAddReview(false)}
                title="Write a Review"
            >
                <form onSubmit={handleAddReview} className="space-y-4">
                    <div>
                        <label className="label">Rating</label>
                        <StarRating value={reviewRating} onChange={setReviewRating} size="lg" />
                    </div>
                    <div>
                        <label className="label">Conditions</label>
                        <div className="flex gap-2">
                            {(['DRY', 'WET'] as const).map((cond) => (
                                <button
                                    key={cond}
                                    type="button"
                                    onClick={() => setReviewConditions(cond)}
                                    className={cn(
                                        'flex-1 rounded-lg px-3 py-2.5 text-sm font-medium transition-colors',
                                        reviewConditions === cond
                                            ? 'bg-brand-600 text-white'
                                            : 'bg-surface-700 text-surface-300'
                                    )}
                                >
                                    {cond === 'DRY' ? '☀️ Dry' : '🌧 Wet'}
                                </button>
                            ))}
                        </div>
                    </div>
                    <div>
                        <label className="label">Event Type (optional)</label>
                        <select
                            value={reviewEventId}
                            onChange={(e) => setReviewEventId(e.target.value)}
                            className="input-field"
                        >
                            <option value="">General</option>
                            {track.events.map((event) => (
                                <option key={event.id} value={event.id}>
                                    {getEventLabel(event.eventType)}
                                </option>
                            ))}
                        </select>
                    </div>
                    <div>
                        <label className="label">Review (optional)</label>
                        <textarea
                            value={reviewContent}
                            onChange={(e) => setReviewContent(e.target.value)}
                            placeholder="Share your experience at this track..."
                            className="input-field min-h-[80px] resize-none"
                            rows={3}
                        />
                    </div>
                    <button type="submit" disabled={submitting || reviewRating === 0} className="btn-primary w-full">
                        {submitting ? 'Posting...' : 'Submit Review'}
                    </button>
                </form>
            </BottomSheet>
        </div>
    );
}
