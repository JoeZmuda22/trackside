import { z } from 'zod';

// ─── Auth ───────────────────────────────────────────────────────────────────────

export const registerSchema = z.object({
    name: z.string().min(2, 'Name must be at least 2 characters'),
    email: z.string().email('Invalid email address'),
    password: z.string().min(8, 'Password must be at least 8 characters'),
    confirmPassword: z.string(),
}).refine((data) => data.password === data.confirmPassword, {
    message: 'Passwords do not match',
    path: ['confirmPassword'],
});

export const loginSchema = z.object({
    email: z.string().email('Invalid email address'),
    password: z.string().min(1, 'Password is required'),
});

// ─── Profile ────────────────────────────────────────────────────────────────────

export const profileSchema = z.object({
    name: z.string().min(2, 'Name must be at least 2 characters'),
    experience: z.enum(['BEGINNER', 'INTERMEDIATE', 'ADVANCED', 'PRO']),
});

export const carSchema = z.object({
    make: z.string().min(1, 'Make is required'),
    model: z.string().min(1, 'Model is required'),
    year: z.number().int().min(1900).max(2030),
});

export const carModSchema = z.object({
    name: z.string().min(1, 'Mod name is required'),
    category: z.enum([
        'ENGINE', 'SUSPENSION', 'AERO', 'BRAKES', 'WHEELS_TIRES',
        'DRIVETRAIN', 'EXHAUST', 'INTERIOR', 'EXTERIOR', 'ELECTRONICS', 'OTHER',
    ]),
    notes: z.string().optional(),
    carId: z.string(),
});

// ─── Tracks ─────────────────────────────────────────────────────────────────────

export const trackSchema = z.object({
    name: z.string().min(2, 'Track name is required'),
    location: z.string().min(2, 'Location is required'),
    description: z.string().optional(),
    imageUrl: z.string().optional(),
    eventTypes: z.array(z.enum(['DRIFT', 'DRAG', 'GRIP'])).min(1, 'Select at least one event type'),
});

export const trackZoneSchema = z.object({
    name: z.string().min(1, 'Zone name is required'),
    description: z.string().optional(),
    posX: z.number().min(0).max(100),
    posY: z.number().min(0).max(100),
    trackId: z.string(),
});

export const zoneTipSchema = z.object({
    content: z.string().min(1, 'Tip content is required'),
    conditions: z.enum(['DRY', 'WET']).optional(),
    zoneId: z.string(),
});

// ─── Reviews ────────────────────────────────────────────────────────────────────

export const trackReviewSchema = z.object({
    rating: z.number().int().min(1).max(5),
    content: z.string().optional(),
    conditions: z.enum(['DRY', 'WET']),
    trackId: z.string(),
    trackEventId: z.string().optional(),
});

// ─── Lap Records ────────────────────────────────────────────────────────────────

export const lapRecordSchema = z.object({
    lapTime: z.string().min(1, 'Lap time is required'),
    conditions: z.enum(['DRY', 'WET']),
    notes: z.string().optional(),
    tirePressureFL: z.number().positive().optional(),
    tirePressureFR: z.number().positive().optional(),
    tirePressureRL: z.number().positive().optional(),
    tirePressureRR: z.number().positive().optional(),
    fuelLevel: z.number().min(0).optional(),
    camberFL: z.number().optional(),
    camberFR: z.number().optional(),
    camberRL: z.number().optional(),
    camberRR: z.number().optional(),
    casterFL: z.number().optional(),
    casterFR: z.number().optional(),
    toeFL: z.number().optional(),
    toeFR: z.number().optional(),
    toeRL: z.number().optional(),
    toeRR: z.number().optional(),
    trackId: z.string(),
    trackEventId: z.string().optional(),
    carId: z.string(),
});

// ─── Types ──────────────────────────────────────────────────────────────────────

export type RegisterInput = z.infer<typeof registerSchema>;
export type LoginInput = z.infer<typeof loginSchema>;
export type ProfileInput = z.infer<typeof profileSchema>;
export type CarInput = z.infer<typeof carSchema>;
export type CarModInput = z.infer<typeof carModSchema>;
export type TrackInput = z.infer<typeof trackSchema>;
export type TrackZoneInput = z.infer<typeof trackZoneSchema>;
export type ZoneTipInput = z.infer<typeof zoneTipSchema>;
export type TrackReviewInput = z.infer<typeof trackReviewSchema>;
export type LapRecordInput = z.infer<typeof lapRecordSchema>;
