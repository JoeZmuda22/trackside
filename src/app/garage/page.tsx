'use client';

import { useState, useEffect } from 'react';
import { useAuth } from '@/components/providers/auth-provider';
import { api } from '@/lib/api';
import Link from 'next/link';
import { cn } from '@/lib/utils';
import { BottomSheet } from '@/components/ui/bottom-sheet';
import { LoadingSpinner, EmptyState } from '@/components/ui/loading';

interface CarMod {
    id: string;
    name: string;
    category: string;
    notes: string | null;
}

interface Car {
    id: string;
    make: string;
    model: string;
    year: number;
    mods: CarMod[];
}

const MOD_CATEGORIES = [
    'ENGINE', 'SUSPENSION', 'AERO', 'BRAKES', 'WHEELS_TIRES',
    'DRIVETRAIN', 'EXHAUST', 'INTERIOR', 'EXTERIOR', 'ELECTRONICS', 'OTHER',
];

const MOD_CATEGORY_LABELS: Record<string, string> = {
    ENGINE: 'Engine',
    SUSPENSION: 'Suspension',
    AERO: 'Aero',
    BRAKES: 'Brakes',
    WHEELS_TIRES: 'Wheels & Tires',
    DRIVETRAIN: 'Drivetrain',
    EXHAUST: 'Exhaust',
    INTERIOR: 'Interior',
    EXTERIOR: 'Exterior',
    ELECTRONICS: 'Electronics',
    OTHER: 'Other',
};

export default function GaragePage() {
    const { user } = useAuth();
    const [cars, setCars] = useState<Car[]>([]);
    const [loading, setLoading] = useState(true);
    const [showAddCar, setShowAddCar] = useState(false);
    const [showAddMod, setShowAddMod] = useState(false);
    const [selectedCarId, setSelectedCarId] = useState<string | null>(null);
    const [expandedCar, setExpandedCar] = useState<string | null>(null);

    // Car form
    const [carMake, setCarMake] = useState('');
    const [carModel, setCarModel] = useState('');
    const [carYear, setCarYear] = useState(new Date().getFullYear());

    // Mod form
    const [modName, setModName] = useState('');
    const [modCategory, setModCategory] = useState('ENGINE');
    const [modNotes, setModNotes] = useState('');

    const [submitting, setSubmitting] = useState(false);

    useEffect(() => {
        if (user) fetchCars();
    }, [user]);

    const fetchCars = async () => {
        try {
            const res = await api('/api/cars');
            if (res.ok) {
                const data = await res.json();
                setCars(data);
            }
        } catch (error) {
            console.error('Failed to fetch cars:', error);
        } finally {
            setLoading(false);
        }
    };

    const handleAddCar = async (e: React.FormEvent) => {
        e.preventDefault();
        setSubmitting(true);

        try {
            const res = await api('/api/cars', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ make: carMake, model: carModel, year: carYear }),
            });

            if (res.ok) {
                setShowAddCar(false);
                setCarMake('');
                setCarModel('');
                setCarYear(new Date().getFullYear());
                fetchCars();
            }
        } catch (error) {
            console.error('Failed to add car:', error);
        } finally {
            setSubmitting(false);
        }
    };

    const handleDeleteCar = async (carId: string) => {
        if (!confirm('Delete this car and all its mods?')) return;

        try {
            await api(`/api/cars/${carId}`, { method: 'DELETE' });
            fetchCars();
        } catch (error) {
            console.error('Failed to delete car:', error);
        }
    };

    const handleAddMod = async (e: React.FormEvent) => {
        e.preventDefault();
        if (!selectedCarId) return;
        setSubmitting(true);

        try {
            const res = await api(`/api/cars/${selectedCarId}/mods`, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({
                    name: modName,
                    category: modCategory,
                    notes: modNotes || undefined,
                }),
            });

            if (res.ok) {
                setShowAddMod(false);
                setModName('');
                setModCategory('ENGINE');
                setModNotes('');
                fetchCars();
            }
        } catch (error) {
            console.error('Failed to add mod:', error);
        } finally {
            setSubmitting(false);
        }
    };

    const handleDeleteMod = async (carId: string, modId: string) => {
        try {
            await api(`/api/cars/${carId}/mods/${modId}`, { method: 'DELETE' });
            fetchCars();
        } catch (error) {
            console.error('Failed to delete mod:', error);
        }
    };

    if (!user) {
        return (
            <div className="page-container flex flex-col items-center justify-center min-h-[60vh]">
                <p className="text-surface-400 mb-4">Sign in to manage your garage</p>
                <Link href="/login" className="btn-primary">Sign In</Link>
            </div>
        );
    }

    if (loading) return <div className="page-container"><LoadingSpinner /></div>;

    return (
        <div className="page-container">
            <div className="flex items-center justify-between mb-6">
                <h1 className="page-title mb-0">My Garage</h1>
                <button onClick={() => setShowAddCar(true)} className="btn-primary text-sm">
                    + Add Car
                </button>
            </div>

            {cars.length === 0 ? (
                <EmptyState
                    icon={
                        <svg className="w-12 h-12" fill="none" viewBox="0 0 24 24" strokeWidth={1} stroke="currentColor">
                            <path strokeLinecap="round" strokeLinejoin="round" d="M11.42 15.17 17.25 21A2.652 2.652 0 0 0 21 17.25l-5.877-5.877M11.42 15.17l2.496-3.03c.317-.384.74-.626 1.208-.766M11.42 15.17l-4.655 5.653a2.548 2.548 0 1 1-3.586-3.586l6.837-5.63m5.108-.233c.55-.164 1.163-.188 1.743-.14a4.5 4.5 0 0 0 4.486-6.336l-3.276 3.277a3.004 3.004 0 0 1-2.25-2.25l3.276-3.276a4.5 4.5 0 0 0-6.336 4.486c.091 1.076-.071 2.264-.904 2.95l-.102.085m-1.745 1.437L5.909 7.5H4.5L2.25 3.75l1.5-1.5L7.5 4.5v1.409l4.26 4.26m-1.745 1.437 1.745-1.437m6.615 8.206L15.75 15.75M4.867 19.125h.008v.008h-.008v-.008Z" />
                        </svg>
                    }
                    title="Garage is empty"
                    description="Add your first car to get started"
                    action={
                        <button onClick={() => setShowAddCar(true)} className="btn-primary text-sm">
                            + Add Car
                        </button>
                    }
                />
            ) : (
                <div className="space-y-3">
                    {cars.map((car) => (
                        <div key={car.id} className="card">
                            <button
                                onClick={() => setExpandedCar(expandedCar === car.id ? null : car.id)}
                                className="w-full flex items-center justify-between"
                            >
                                <div className="text-left">
                                    <h3 className="font-semibold text-white">
                                        {car.year} {car.make} {car.model}
                                    </h3>
                                    <p className="text-xs text-surface-500">
                                        {car.mods.length} mod{car.mods.length !== 1 ? 's' : ''}
                                    </p>
                                </div>
                                <svg
                                    className={cn(
                                        'w-5 h-5 text-surface-400 transition-transform',
                                        expandedCar === car.id && 'rotate-180'
                                    )}
                                    fill="none"
                                    viewBox="0 0 24 24"
                                    strokeWidth={2}
                                    stroke="currentColor"
                                >
                                    <path strokeLinecap="round" strokeLinejoin="round" d="m19.5 8.25-7.5 7.5-7.5-7.5" />
                                </svg>
                            </button>

                            {expandedCar === car.id && (
                                <div className="mt-3 pt-3 border-t border-surface-700">
                                    {/* Mods list */}
                                    {car.mods.length === 0 ? (
                                        <p className="text-sm text-surface-500 mb-3">No mods yet</p>
                                    ) : (
                                        <div className="space-y-2 mb-3">
                                            {car.mods.map((mod) => (
                                                <div key={mod.id} className="flex items-center justify-between rounded-lg bg-surface-700/50 p-2.5">
                                                    <div>
                                                        <p className="text-sm font-medium text-surface-200">{mod.name}</p>
                                                        <p className="text-xs text-surface-500">
                                                            {MOD_CATEGORY_LABELS[mod.category] || mod.category}
                                                            {mod.notes && ` • ${mod.notes}`}
                                                        </p>
                                                    </div>
                                                    <button
                                                        onClick={() => handleDeleteMod(car.id, mod.id)}
                                                        className="p-1 text-surface-500 hover:text-red-400 transition-colors"
                                                    >
                                                        <svg className="w-4 h-4" fill="none" viewBox="0 0 24 24" strokeWidth={1.5} stroke="currentColor">
                                                            <path strokeLinecap="round" strokeLinejoin="round" d="m14.74 9-.346 9m-4.788 0L9.26 9m9.968-3.21c.342.052.682.107 1.022.166m-1.022-.165L18.16 19.673a2.25 2.25 0 0 1-2.244 2.077H8.084a2.25 2.25 0 0 1-2.244-2.077L4.772 5.79m14.456 0a48.108 48.108 0 0 0-3.478-.397m-12 .562c.34-.059.68-.114 1.022-.165m0 0a48.11 48.11 0 0 1 3.478-.397m7.5 0v-.916c0-1.18-.91-2.164-2.09-2.201a51.964 51.964 0 0 0-3.32 0c-1.18.037-2.09 1.022-2.09 2.201v.916m7.5 0a48.667 48.667 0 0 0-7.5 0" />
                                                        </svg>
                                                    </button>
                                                </div>
                                            ))}
                                        </div>
                                    )}

                                    <div className="flex gap-2">
                                        <button
                                            onClick={() => {
                                                setSelectedCarId(car.id);
                                                setShowAddMod(true);
                                            }}
                                            className="btn-secondary flex-1 text-sm"
                                        >
                                            + Add Mod
                                        </button>
                                        <button
                                            onClick={() => handleDeleteCar(car.id)}
                                            className="btn-ghost text-sm text-red-400 px-3"
                                        >
                                            Delete Car
                                        </button>
                                    </div>
                                </div>
                            )}
                        </div>
                    ))}
                </div>
            )}

            {/* Add Car Bottom Sheet */}
            <BottomSheet isOpen={showAddCar} onClose={() => setShowAddCar(false)} title="Add Car">
                <form onSubmit={handleAddCar} className="space-y-4">
                    <div>
                        <label className="label">Make</label>
                        <input
                            type="text"
                            value={carMake}
                            onChange={(e) => setCarMake(e.target.value)}
                            className="input-field"
                            placeholder="e.g. Nissan"
                            required
                        />
                    </div>
                    <div>
                        <label className="label">Model</label>
                        <input
                            type="text"
                            value={carModel}
                            onChange={(e) => setCarModel(e.target.value)}
                            className="input-field"
                            placeholder="e.g. 350Z"
                            required
                        />
                    </div>
                    <div>
                        <label className="label">Year</label>
                        <input
                            type="number"
                            value={carYear}
                            onChange={(e) => setCarYear(parseInt(e.target.value))}
                            className="input-field"
                            min={1900}
                            max={2030}
                            required
                        />
                    </div>
                    <button type="submit" disabled={submitting} className="btn-primary w-full">
                        {submitting ? 'Adding...' : 'Add Car'}
                    </button>
                </form>
            </BottomSheet>

            {/* Add Mod Bottom Sheet */}
            <BottomSheet isOpen={showAddMod} onClose={() => setShowAddMod(false)} title="Add Modification">
                <form onSubmit={handleAddMod} className="space-y-4">
                    <div>
                        <label className="label">Mod Name</label>
                        <input
                            type="text"
                            value={modName}
                            onChange={(e) => setModName(e.target.value)}
                            className="input-field"
                            placeholder="e.g. Coilovers"
                            required
                        />
                    </div>
                    <div>
                        <label className="label">Category</label>
                        <select
                            value={modCategory}
                            onChange={(e) => setModCategory(e.target.value)}
                            className="input-field"
                        >
                            {MOD_CATEGORIES.map((cat) => (
                                <option key={cat} value={cat}>
                                    {MOD_CATEGORY_LABELS[cat]}
                                </option>
                            ))}
                        </select>
                    </div>
                    <div>
                        <label className="label">Notes (optional)</label>
                        <input
                            type="text"
                            value={modNotes}
                            onChange={(e) => setModNotes(e.target.value)}
                            className="input-field"
                            placeholder="e.g. BC Racing BR Series"
                        />
                    </div>
                    <button type="submit" disabled={submitting} className="btn-primary w-full">
                        {submitting ? 'Adding...' : 'Add Mod'}
                    </button>
                </form>
            </BottomSheet>
        </div>
    );
}
