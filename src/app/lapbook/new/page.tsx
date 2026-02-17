'use client';

import { useState, useEffect, Suspense } from 'react';
import { useAuth } from '@/components/providers/auth-provider';
import { api } from '@/lib/api';
import { useRouter, useSearchParams } from 'next/navigation';
import Link from 'next/link';
import { cn, getEventLabel } from '@/lib/utils';
import { LoadingSpinner } from '@/components/ui/loading';

interface Track {
    id: string;
    name: string;
    location: string;
    events: { id: string; eventType: string }[];
}

interface Car {
    id: string;
    make: string;
    model: string;
    year: number;
}

function NewLapForm() {
    const router = useRouter();
    const searchParams = useSearchParams();
    const { user } = useAuth();
    const preselectedTrackId = searchParams.get('trackId') || '';

    const [tracks, setTracks] = useState<Track[]>([]);
    const [cars, setCars] = useState<Car[]>([]);
    const [loading, setLoading] = useState(true);
    const [submitting, setSubmitting] = useState(false);
    const [error, setError] = useState('');

    // Form state
    const [trackId, setTrackId] = useState(preselectedTrackId);
    const [carId, setCarId] = useState('');
    const [trackEventId, setTrackEventId] = useState('');
    const [lapTime, setLapTime] = useState('');
    const [conditions, setConditions] = useState<'DRY' | 'WET'>('DRY');
    const [notes, setNotes] = useState('');

    // Tire pressures
    const [tirePressureFL, setTirePressureFL] = useState('');
    const [tirePressureFR, setTirePressureFR] = useState('');
    const [tirePressureRL, setTirePressureRL] = useState('');
    const [tirePressureRR, setTirePressureRR] = useState('');

    // Fuel
    const [fuelLevel, setFuelLevel] = useState('');

    // Alignment - Camber
    const [camberFL, setCamberFL] = useState('');
    const [camberFR, setCamberFR] = useState('');
    const [camberRL, setCamberRL] = useState('');
    const [camberRR, setCamberRR] = useState('');

    // Alignment - Caster
    const [casterFL, setCasterFL] = useState('');
    const [casterFR, setCasterFR] = useState('');

    // Alignment - Toe
    const [toeFL, setToeFL] = useState('');
    const [toeFR, setToeFR] = useState('');
    const [toeRL, setToeRL] = useState('');
    const [toeRR, setToeRR] = useState('');

    const [showTirePressure, setShowTirePressure] = useState(false);
    const [showAlignment, setShowAlignment] = useState(false);

    useEffect(() => {
        if (user) {
            Promise.all([
                api('/api/tracks').then((r) => r.json()),
                api('/api/cars').then((r) => r.json()),
            ]).then(([tracksData, carsData]) => {
                setTracks(tracksData);
                setCars(carsData);
                setLoading(false);
            });
        }
    }, [user]);

    const selectedTrack = tracks.find((t) => t.id === trackId);

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        setError('');
        setSubmitting(true);

        const toNum = (v: string) => (v ? parseFloat(v) : undefined);

        try {
            const res = await api('/api/lapbook', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({
                    lapTime,
                    conditions,
                    notes: notes || undefined,
                    trackId,
                    trackEventId: trackEventId || undefined,
                    carId,
                    tirePressureFL: toNum(tirePressureFL),
                    tirePressureFR: toNum(tirePressureFR),
                    tirePressureRL: toNum(tirePressureRL),
                    tirePressureRR: toNum(tirePressureRR),
                    fuelLevel: toNum(fuelLevel),
                    camberFL: toNum(camberFL),
                    camberFR: toNum(camberFR),
                    camberRL: toNum(camberRL),
                    camberRR: toNum(camberRR),
                    casterFL: toNum(casterFL),
                    casterFR: toNum(casterFR),
                    toeFL: toNum(toeFL),
                    toeFR: toNum(toeFR),
                    toeRL: toNum(toeRL),
                    toeRR: toNum(toeRR),
                }),
            });

            if (!res.ok) {
                const data = await res.json();
                setError(data.error || 'Failed to log lap');
                return;
            }

            router.push('/lapbook');
        } catch {
            setError('Something went wrong');
        } finally {
            setSubmitting(false);
        }
    };

    if (!user) {
        return (
            <div className="page-container flex flex-col items-center justify-center min-h-[60vh]">
                <p className="text-surface-400 mb-4">Sign in to log laps</p>
                <Link href="/login" className="btn-primary">Sign In</Link>
            </div>
        );
    }

    if (loading) return <div className="page-container"><LoadingSpinner /></div>;

    if (cars.length === 0) {
        return (
            <div className="page-container flex flex-col items-center justify-center min-h-[60vh]">
                <p className="text-surface-400 mb-2">You need to add a car first</p>
                <Link href="/garage" className="btn-primary text-sm">Go to Garage</Link>
            </div>
        );
    }

    return (
        <div className="page-container">
            <div className="flex items-center gap-3 mb-6">
                <button onClick={() => router.back()} className="btn-ghost p-2">
                    <svg className="w-5 h-5" fill="none" viewBox="0 0 24 24" strokeWidth={2} stroke="currentColor">
                        <path strokeLinecap="round" strokeLinejoin="round" d="M15.75 19.5 8.25 12l7.5-7.5" />
                    </svg>
                </button>
                <h1 className="page-title mb-0">Log Lap Time</h1>
            </div>

            <form onSubmit={handleSubmit} className="space-y-5">
                {error && (
                    <div className="rounded-lg bg-red-500/10 border border-red-500/20 p-3 text-sm text-red-400">
                        {error}
                    </div>
                )}

                {/* Lap Time — big input */}
                <div>
                    <label className="label">Lap Time</label>
                    <input
                        type="text"
                        value={lapTime}
                        onChange={(e) => setLapTime(e.target.value)}
                        className="input-field text-2xl font-mono text-center"
                        placeholder="1:23.456"
                        required
                    />
                </div>

                {/* Track Selection */}
                <div>
                    <label className="label">Track</label>
                    <select
                        value={trackId}
                        onChange={(e) => {
                            setTrackId(e.target.value);
                            setTrackEventId('');
                        }}
                        className="input-field"
                        required
                    >
                        <option value="">Select a track</option>
                        {tracks.map((t) => (
                            <option key={t.id} value={t.id}>
                                {t.name} — {t.location}
                            </option>
                        ))}
                    </select>
                </div>

                {/* Event Type */}
                {selectedTrack && selectedTrack.events.length > 0 && (
                    <div>
                        <label className="label">Event Type</label>
                        <div className="flex gap-2 flex-wrap">
                            {selectedTrack.events.map((event) => (
                                <button
                                    key={event.id}
                                    type="button"
                                    onClick={() => setTrackEventId(trackEventId === event.id ? '' : event.id)}
                                    className={cn(
                                        'rounded-full px-4 py-2 text-sm font-medium transition-colors',
                                        trackEventId === event.id
                                            ? 'bg-brand-600 text-white'
                                            : 'bg-surface-800 text-surface-300'
                                    )}
                                >
                                    {getEventLabel(event.eventType)}
                                </button>
                            ))}
                        </div>
                    </div>
                )}

                {/* Car Selection */}
                <div>
                    <label className="label">Car</label>
                    <select
                        value={carId}
                        onChange={(e) => setCarId(e.target.value)}
                        className="input-field"
                        required
                    >
                        <option value="">Select your car</option>
                        {cars.map((c) => (
                            <option key={c.id} value={c.id}>
                                {c.year} {c.make} {c.model}
                            </option>
                        ))}
                    </select>
                </div>

                {/* Conditions */}
                <div>
                    <label className="label">Driving Conditions</label>
                    <div className="flex gap-2">
                        {(['DRY', 'WET'] as const).map((cond) => (
                            <button
                                key={cond}
                                type="button"
                                onClick={() => setConditions(cond)}
                                className={cn(
                                    'flex-1 rounded-lg px-3 py-3 text-sm font-medium transition-colors',
                                    conditions === cond
                                        ? 'bg-brand-600 text-white'
                                        : 'bg-surface-800 text-surface-300'
                                )}
                            >
                                {cond === 'DRY' ? '☀️ Dry' : '🌧 Wet'}
                            </button>
                        ))}
                    </div>
                </div>

                {/* Fuel Level */}
                <div>
                    <label className="label">Fuel Level</label>
                    <input
                        type="number"
                        inputMode="decimal"
                        value={fuelLevel}
                        onChange={(e) => setFuelLevel(e.target.value)}
                        className="input-field"
                        placeholder="e.g. 75 (%)"
                        step="0.1"
                    />
                </div>

                {/* Tire Pressure - Collapsible */}
                <div>
                    <button
                        type="button"
                        onClick={() => setShowTirePressure(!showTirePressure)}
                        className="flex items-center justify-between w-full py-2"
                    >
                        <span className="text-sm font-medium text-surface-300">Tire Pressure (PSI)</span>
                        <svg
                            className={cn('w-4 h-4 text-surface-400 transition-transform', showTirePressure && 'rotate-180')}
                            fill="none" viewBox="0 0 24 24" strokeWidth={2} stroke="currentColor"
                        >
                            <path strokeLinecap="round" strokeLinejoin="round" d="m19.5 8.25-7.5 7.5-7.5-7.5" />
                        </svg>
                    </button>
                    {showTirePressure && (
                        <div className="grid grid-cols-2 gap-3 mt-2">
                            <div>
                                <label className="text-xs text-surface-500">Front Left</label>
                                <input type="number" inputMode="decimal" value={tirePressureFL} onChange={(e) => setTirePressureFL(e.target.value)} className="input-field" placeholder="FL" step="0.1" />
                            </div>
                            <div>
                                <label className="text-xs text-surface-500">Front Right</label>
                                <input type="number" inputMode="decimal" value={tirePressureFR} onChange={(e) => setTirePressureFR(e.target.value)} className="input-field" placeholder="FR" step="0.1" />
                            </div>
                            <div>
                                <label className="text-xs text-surface-500">Rear Left</label>
                                <input type="number" inputMode="decimal" value={tirePressureRL} onChange={(e) => setTirePressureRL(e.target.value)} className="input-field" placeholder="RL" step="0.1" />
                            </div>
                            <div>
                                <label className="text-xs text-surface-500">Rear Right</label>
                                <input type="number" inputMode="decimal" value={tirePressureRR} onChange={(e) => setTirePressureRR(e.target.value)} className="input-field" placeholder="RR" step="0.1" />
                            </div>
                        </div>
                    )}
                </div>

                {/* Alignment - Collapsible */}
                <div>
                    <button
                        type="button"
                        onClick={() => setShowAlignment(!showAlignment)}
                        className="flex items-center justify-between w-full py-2"
                    >
                        <span className="text-sm font-medium text-surface-300">Alignment Settings</span>
                        <svg
                            className={cn('w-4 h-4 text-surface-400 transition-transform', showAlignment && 'rotate-180')}
                            fill="none" viewBox="0 0 24 24" strokeWidth={2} stroke="currentColor"
                        >
                            <path strokeLinecap="round" strokeLinejoin="round" d="m19.5 8.25-7.5 7.5-7.5-7.5" />
                        </svg>
                    </button>
                    {showAlignment && (
                        <div className="space-y-4 mt-2">
                            {/* Camber */}
                            <div>
                                <h4 className="text-xs font-medium text-surface-400 mb-2">CAMBER (degrees)</h4>
                                <div className="grid grid-cols-2 gap-3">
                                    <div>
                                        <label className="text-xs text-surface-500">FL</label>
                                        <input type="number" inputMode="decimal" value={camberFL} onChange={(e) => setCamberFL(e.target.value)} className="input-field" placeholder="-1.5" step="0.1" />
                                    </div>
                                    <div>
                                        <label className="text-xs text-surface-500">FR</label>
                                        <input type="number" inputMode="decimal" value={camberFR} onChange={(e) => setCamberFR(e.target.value)} className="input-field" placeholder="-1.5" step="0.1" />
                                    </div>
                                    <div>
                                        <label className="text-xs text-surface-500">RL</label>
                                        <input type="number" inputMode="decimal" value={camberRL} onChange={(e) => setCamberRL(e.target.value)} className="input-field" placeholder="-1.0" step="0.1" />
                                    </div>
                                    <div>
                                        <label className="text-xs text-surface-500">RR</label>
                                        <input type="number" inputMode="decimal" value={camberRR} onChange={(e) => setCamberRR(e.target.value)} className="input-field" placeholder="-1.0" step="0.1" />
                                    </div>
                                </div>
                            </div>

                            {/* Caster */}
                            <div>
                                <h4 className="text-xs font-medium text-surface-400 mb-2">CASTER (degrees)</h4>
                                <div className="grid grid-cols-2 gap-3">
                                    <div>
                                        <label className="text-xs text-surface-500">FL</label>
                                        <input type="number" inputMode="decimal" value={casterFL} onChange={(e) => setCasterFL(e.target.value)} className="input-field" placeholder="5.0" step="0.1" />
                                    </div>
                                    <div>
                                        <label className="text-xs text-surface-500">FR</label>
                                        <input type="number" inputMode="decimal" value={casterFR} onChange={(e) => setCasterFR(e.target.value)} className="input-field" placeholder="5.0" step="0.1" />
                                    </div>
                                </div>
                            </div>

                            {/* Toe */}
                            <div>
                                <h4 className="text-xs font-medium text-surface-400 mb-2">TOE (degrees, + = toe-in)</h4>
                                <div className="grid grid-cols-2 gap-3">
                                    <div>
                                        <label className="text-xs text-surface-500">FL</label>
                                        <input type="number" inputMode="decimal" value={toeFL} onChange={(e) => setToeFL(e.target.value)} className="input-field" placeholder="0.1" step="0.01" />
                                    </div>
                                    <div>
                                        <label className="text-xs text-surface-500">FR</label>
                                        <input type="number" inputMode="decimal" value={toeFR} onChange={(e) => setToeFR(e.target.value)} className="input-field" placeholder="0.1" step="0.01" />
                                    </div>
                                    <div>
                                        <label className="text-xs text-surface-500">RL</label>
                                        <input type="number" inputMode="decimal" value={toeRL} onChange={(e) => setToeRL(e.target.value)} className="input-field" placeholder="0.1" step="0.01" />
                                    </div>
                                    <div>
                                        <label className="text-xs text-surface-500">RR</label>
                                        <input type="number" inputMode="decimal" value={toeRR} onChange={(e) => setToeRR(e.target.value)} className="input-field" placeholder="0.1" step="0.01" />
                                    </div>
                                </div>
                            </div>
                        </div>
                    )}
                </div>

                {/* Notes */}
                <div>
                    <label className="label">Notes (optional)</label>
                    <textarea
                        value={notes}
                        onChange={(e) => setNotes(e.target.value)}
                        className="input-field min-h-[60px] resize-none"
                        placeholder="Track conditions, car feel, setup changes..."
                        rows={2}
                    />
                </div>

                <button
                    type="submit"
                    disabled={submitting}
                    className="btn-primary w-full"
                >
                    {submitting ? (
                        <div className="h-5 w-5 animate-spin rounded-full border-2 border-white/30 border-t-white" />
                    ) : (
                        'Save Lap Record'
                    )}
                </button>
            </form>
        </div>
    );
}

export default function NewLapPage() {
    return (
        <Suspense fallback={<div className="page-container"><LoadingSpinner /></div>}>
            <NewLapForm />
        </Suspense>
    );
}
