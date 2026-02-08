import { NextResponse } from 'next/server';
import { auth } from '@/lib/auth';
import { prisma } from '@/lib/db';
import { lapRecordSchema } from '@/lib/validations';

// GET /api/lapbook — get current user's lap records
export async function GET(request: Request) {
    const session = await auth();
    if (!session?.user?.id) {
        return NextResponse.json({ error: 'Unauthorized' }, { status: 401 });
    }

    const { searchParams } = new URL(request.url);
    const trackId = searchParams.get('trackId');
    const eventType = searchParams.get('eventType');
    const carId = searchParams.get('carId');

    const where: any = {
        driverId: session.user.id,
    };

    if (trackId) where.trackId = trackId;
    if (carId) where.carId = carId;
    if (eventType) {
        where.trackEvent = { eventType };
    }

    const lapRecords = await prisma.lapRecord.findMany({
        where,
        include: {
            track: { select: { id: true, name: true, location: true } },
            trackEvent: true,
            car: { select: { id: true, make: true, model: true, year: true } },
        },
        orderBy: { createdAt: 'desc' },
    });

    return NextResponse.json(lapRecords);
}

// POST /api/lapbook — create a new lap record
export async function POST(request: Request) {
    const session = await auth();
    if (!session?.user?.id) {
        return NextResponse.json({ error: 'Unauthorized' }, { status: 401 });
    }

    try {
        const body = await request.json();
        const validated = lapRecordSchema.parse(body);

        // Verify car belongs to user
        const car = await prisma.car.findFirst({
            where: { id: validated.carId, userId: session.user.id },
        });

        if (!car) {
            return NextResponse.json({ error: 'Car not found' }, { status: 404 });
        }

        // Verify track exists
        const track = await prisma.track.findUnique({
            where: { id: validated.trackId },
        });

        if (!track) {
            return NextResponse.json({ error: 'Track not found' }, { status: 404 });
        }

        const lapRecord = await prisma.lapRecord.create({
            data: {
                lapTime: validated.lapTime,
                conditions: validated.conditions,
                notes: validated.notes,
                tirePressureFL: validated.tirePressureFL,
                tirePressureFR: validated.tirePressureFR,
                tirePressureRL: validated.tirePressureRL,
                tirePressureRR: validated.tirePressureRR,
                fuelLevel: validated.fuelLevel,
                camberFL: validated.camberFL,
                camberFR: validated.camberFR,
                camberRL: validated.camberRL,
                camberRR: validated.camberRR,
                casterFL: validated.casterFL,
                casterFR: validated.casterFR,
                toeFL: validated.toeFL,
                toeFR: validated.toeFR,
                toeRL: validated.toeRL,
                toeRR: validated.toeRR,
                trackId: validated.trackId,
                trackEventId: validated.trackEventId || null,
                carId: validated.carId,
                driverId: session.user.id,
            },
            include: {
                track: { select: { id: true, name: true, location: true } },
                trackEvent: true,
                car: { select: { id: true, make: true, model: true, year: true } },
            },
        });

        return NextResponse.json(lapRecord, { status: 201 });
    } catch (error: any) {
        if (error?.name === 'ZodError') {
            return NextResponse.json({ error: 'Validation failed', details: error.errors }, { status: 400 });
        }
        console.error('Lap record creation error:', error);
        return NextResponse.json({ error: 'Internal server error' }, { status: 500 });
    }
}
