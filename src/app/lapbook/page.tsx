'use client';

import { useState, useEffect } from 'react';
import { useSession } from 'next-auth/react';
import Link from 'next/link';
import { cn, getEventLabel } from '@/lib/utils';
import { Badge } from '@/components/ui/badge';
import { LoadingSpinner, EmptyState } from '@/components/ui/loading';

interface LapRecord {
    id: string;
    lapTime: string;
    conditions: string;
    notes: string | null;
    tirePressureFL: number | null;
    tirePressureFR: number | null;
    tirePressureRL: number | null;
    tirePressureRR: number | null;
    fuelLevel: number | null;
    camberFL: number | null;
    camberFR: number | null;
    camberRL: number | null;
    camberRR: number | null;
    casterFL: number | null;
    casterFR: number | null;
    toeFL: number | null;
    toeFR: number | null;
    toeRL: number | null;
    toeRR: number | null;
    createdAt: string;
    track: { id: string; name: string; location: string };
    trackEvent: { eventType: string } | null;
    car: { id: string; make: string; model: string; year: number };
}

export default function LapBookPage() {
    const { data: session } = useSession();
    const [records, setRecords] = useState<LapRecord[]>([]);
    const [loading, setLoading] = useState(true);
    const [expandedRecord, setExpandedRecord] = useState<string | null>(null);

    useEffect(() => {
        if (session) fetchRecords();
    }, [session]);

    const fetchRecords = async () => {
        try {
            const res = await fetch('/api/lapbook');
            if (res.ok) {
                const data = await res.json();
                setRecords(data);
            }
        } catch (error) {
            console.error('Failed to fetch records:', error);
        } finally {
            setLoading(false);
        }
    };

    const handleDelete = async (id: string) => {
        if (!confirm('Delete this lap record?')) return;
        try {
            await fetch(`/api/lapbook/${id}`, { method: 'DELETE' });
            fetchRecords();
        } catch (error) {
            console.error('Failed to delete record:', error);
        }
    };

    if (!session) {
        return (
            <div className="page-container flex flex-col items-center justify-center min-h-[60vh]">
                <p className="text-surface-400 mb-4">Sign in to view your lap book</p>
                <Link href="/login" className="btn-primary">Sign In</Link>
            </div>
        );
    }

    if (loading) return <div className="page-container"><LoadingSpinner /></div>;

    return (
        <div className="page-container">
            <div className="flex items-center justify-between mb-6">
                <div>
                    <h1 className="page-title mb-0">My Lap Book</h1>
                    <p className="text-sm text-surface-400">{records.length} recorded lap{records.length !== 1 ? 's' : ''}</p>
                </div>
                <Link href="/lapbook/new" className="btn-primary text-sm">
                    + Log Lap
                </Link>
            </div>

            {records.length === 0 ? (
                <EmptyState
                    icon={
                        <svg className="w-12 h-12" fill="none" viewBox="0 0 24 24" strokeWidth={1} stroke="currentColor">
                            <path strokeLinecap="round" strokeLinejoin="round" d="M12 6v6h4.5m4.5 0a9 9 0 1 1-18 0 9 9 0 0 1 18 0Z" />
                        </svg>
                    }
                    title="No lap records yet"
                    description="Start logging your lap times with full telemetry data"
                    action={
                        <Link href="/lapbook/new" className="btn-primary text-sm">
                            Log Your First Lap
                        </Link>
                    }
                />
            ) : (
                <div className="space-y-3">
                    {records.map((record) => (
                        <div key={record.id} className="card">
                            <button
                                onClick={() => setExpandedRecord(expandedRecord === record.id ? null : record.id)}
                                className="w-full text-left"
                            >
                                <div className="flex items-start justify-between">
                                    <div>
                                        <p className="text-xl font-mono font-bold text-white">{record.lapTime}</p>
                                        <p className="text-sm text-surface-400 mt-0.5">{record.track.name}</p>
                                        <p className="text-xs text-surface-500">
                                            {record.car.year} {record.car.make} {record.car.model}
                                        </p>
                                    </div>
                                    <div className="flex flex-col items-end gap-1">
                                        <Badge variant="condition" value={record.conditions}>
                                            {record.conditions === 'WET' ? '🌧 Wet' : '☀️ Dry'}
                                        </Badge>
                                        {record.trackEvent && (
                                            <Badge variant="event" value={record.trackEvent.eventType}>
                                                {getEventLabel(record.trackEvent.eventType)}
                                            </Badge>
                                        )}
                                    </div>
                                </div>
                            </button>

                            {expandedRecord === record.id && (
                                <div className="mt-3 pt-3 border-t border-surface-700 space-y-3">
                                    {/* Tire Pressures */}
                                    {(record.tirePressureFL || record.tirePressureFR || record.tirePressureRL || record.tirePressureRR) && (
                                        <div>
                                            <h4 className="text-xs font-medium text-surface-400 mb-2">TIRE PRESSURE (PSI)</h4>
                                            <div className="grid grid-cols-2 gap-2">
                                                <div className="rounded-lg bg-surface-700/50 p-2 text-center">
                                                    <p className="text-[10px] text-surface-500">FL</p>
                                                    <p className="text-sm font-mono text-surface-200">{record.tirePressureFL ?? '—'}</p>
                                                </div>
                                                <div className="rounded-lg bg-surface-700/50 p-2 text-center">
                                                    <p className="text-[10px] text-surface-500">FR</p>
                                                    <p className="text-sm font-mono text-surface-200">{record.tirePressureFR ?? '—'}</p>
                                                </div>
                                                <div className="rounded-lg bg-surface-700/50 p-2 text-center">
                                                    <p className="text-[10px] text-surface-500">RL</p>
                                                    <p className="text-sm font-mono text-surface-200">{record.tirePressureRL ?? '—'}</p>
                                                </div>
                                                <div className="rounded-lg bg-surface-700/50 p-2 text-center">
                                                    <p className="text-[10px] text-surface-500">RR</p>
                                                    <p className="text-sm font-mono text-surface-200">{record.tirePressureRR ?? '—'}</p>
                                                </div>
                                            </div>
                                        </div>
                                    )}

                                    {/* Fuel */}
                                    {record.fuelLevel != null && (
                                        <div>
                                            <h4 className="text-xs font-medium text-surface-400 mb-1">FUEL LEVEL</h4>
                                            <p className="text-sm font-mono text-surface-200">{record.fuelLevel}</p>
                                        </div>
                                    )}

                                    {/* Alignment - Camber */}
                                    {(record.camberFL || record.camberFR || record.camberRL || record.camberRR) && (
                                        <div>
                                            <h4 className="text-xs font-medium text-surface-400 mb-2">CAMBER (°)</h4>
                                            <div className="grid grid-cols-2 gap-2">
                                                <div className="rounded-lg bg-surface-700/50 p-2 text-center">
                                                    <p className="text-[10px] text-surface-500">FL</p>
                                                    <p className="text-sm font-mono text-surface-200">{record.camberFL ?? '—'}</p>
                                                </div>
                                                <div className="rounded-lg bg-surface-700/50 p-2 text-center">
                                                    <p className="text-[10px] text-surface-500">FR</p>
                                                    <p className="text-sm font-mono text-surface-200">{record.camberFR ?? '—'}</p>
                                                </div>
                                                <div className="rounded-lg bg-surface-700/50 p-2 text-center">
                                                    <p className="text-[10px] text-surface-500">RL</p>
                                                    <p className="text-sm font-mono text-surface-200">{record.camberRL ?? '—'}</p>
                                                </div>
                                                <div className="rounded-lg bg-surface-700/50 p-2 text-center">
                                                    <p className="text-[10px] text-surface-500">RR</p>
                                                    <p className="text-sm font-mono text-surface-200">{record.camberRR ?? '—'}</p>
                                                </div>
                                            </div>
                                        </div>
                                    )}

                                    {/* Alignment - Caster */}
                                    {(record.casterFL || record.casterFR) && (
                                        <div>
                                            <h4 className="text-xs font-medium text-surface-400 mb-2">CASTER (°)</h4>
                                            <div className="grid grid-cols-2 gap-2">
                                                <div className="rounded-lg bg-surface-700/50 p-2 text-center">
                                                    <p className="text-[10px] text-surface-500">FL</p>
                                                    <p className="text-sm font-mono text-surface-200">{record.casterFL ?? '—'}</p>
                                                </div>
                                                <div className="rounded-lg bg-surface-700/50 p-2 text-center">
                                                    <p className="text-[10px] text-surface-500">FR</p>
                                                    <p className="text-sm font-mono text-surface-200">{record.casterFR ?? '—'}</p>
                                                </div>
                                            </div>
                                        </div>
                                    )}

                                    {/* Alignment - Toe */}
                                    {(record.toeFL || record.toeFR || record.toeRL || record.toeRR) && (
                                        <div>
                                            <h4 className="text-xs font-medium text-surface-400 mb-2">TOE (°)</h4>
                                            <div className="grid grid-cols-2 gap-2">
                                                <div className="rounded-lg bg-surface-700/50 p-2 text-center">
                                                    <p className="text-[10px] text-surface-500">FL</p>
                                                    <p className="text-sm font-mono text-surface-200">{record.toeFL ?? '—'}</p>
                                                </div>
                                                <div className="rounded-lg bg-surface-700/50 p-2 text-center">
                                                    <p className="text-[10px] text-surface-500">FR</p>
                                                    <p className="text-sm font-mono text-surface-200">{record.toeFR ?? '—'}</p>
                                                </div>
                                                <div className="rounded-lg bg-surface-700/50 p-2 text-center">
                                                    <p className="text-[10px] text-surface-500">RL</p>
                                                    <p className="text-sm font-mono text-surface-200">{record.toeRL ?? '—'}</p>
                                                </div>
                                                <div className="rounded-lg bg-surface-700/50 p-2 text-center">
                                                    <p className="text-[10px] text-surface-500">RR</p>
                                                    <p className="text-sm font-mono text-surface-200">{record.toeRR ?? '—'}</p>
                                                </div>
                                            </div>
                                        </div>
                                    )}

                                    {/* Notes */}
                                    {record.notes && (
                                        <div>
                                            <h4 className="text-xs font-medium text-surface-400 mb-1">NOTES</h4>
                                            <p className="text-sm text-surface-300">{record.notes}</p>
                                        </div>
                                    )}

                                    <div className="flex gap-2 pt-2">
                                        <p className="text-xs text-surface-600 flex-1">
                                            {new Date(record.createdAt).toLocaleDateString()}
                                        </p>
                                        <button
                                            onClick={() => handleDelete(record.id)}
                                            className="text-xs text-red-400 hover:text-red-300"
                                        >
                                            Delete
                                        </button>
                                    </div>
                                </div>
                            )}
                        </div>
                    ))}
                </div>
            )}
        </div>
    );
}
