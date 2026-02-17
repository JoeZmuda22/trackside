import { NextResponse } from 'next/server';
import { auth } from '@/lib/auth';
import { prisma } from '@/lib/db';
import { trackSchema } from '@/lib/validations';

// GET /api/tracks — list all tracks with filters
export async function GET(request: Request) {
    // Ensure fresh data on every request
    const { searchParams } = new URL(request.url);
    const search = searchParams.get('search') || '';
    const eventType = searchParams.get('eventType');
    const state = searchParams.get('state');

    try {
        const where: any = {
            status: 'APPROVED',
        };

        if (search) {
            where.OR = [
                { name: { contains: search } },
                { location: { contains: search } },
            ];
        }

        if (eventType) {
            where.events = {
                some: { eventType: eventType },
            };
        }

        if (state) {
            where.state = state.toUpperCase();
        }

        const tracks = await prisma.track.findMany({
            where,
            include: {
                events: true,
                uploadedBy: { select: { id: true, name: true } },
                _count: {
                    select: { reviews: true, zones: true, lapRecords: true },
                },
            },
            orderBy: { createdAt: 'desc' },
        });

        // Get all ratings in a single query
        const allRatings = await prisma.trackReview.groupBy({
            by: ['trackId'],
            _avg: { rating: true },
        });

        // Create a map for fast lookup
        const ratingMap = new Map(
            allRatings.map((r) => [r.trackId, r._avg.rating || 0])
        );

        // Add ratings to tracks
        const tracksWithRating = tracks.map((track) => ({
            ...track,
            avgRating: ratingMap.get(track.id) || 0,
        }));

        return NextResponse.json(tracksWithRating);
    } catch (error: any) {
        console.error('Tracks GET error:', error);
        return NextResponse.json(
            { error: 'Failed to fetch tracks', details: error.message },
            { status: 500 }
        );
    }
}

// POST /api/tracks — upload a new track
export async function POST(request: Request) {
    const session = await auth();
    if (!session?.user?.id) {
        return NextResponse.json({ error: 'Unauthorized' }, { status: 401 });
    }

    try {
        const body = await request.json();
        const validated = trackSchema.parse(body);

        const track = await prisma.track.create({
            data: {
                name: validated.name,
                location: validated.location,
                description: validated.description,
                imageUrl: validated.imageUrl,
                uploadedById: session.user.id,
                events: {
                    create: validated.eventTypes.map((et) => ({ eventType: et })),
                },
            },
            include: {
                events: true,
                uploadedBy: { select: { id: true, name: true } },
            },
        });

        return NextResponse.json(track, { status: 201 });
    } catch (error: any) {
        if (error?.name === 'ZodError') {
            return NextResponse.json({ error: 'Validation failed', details: error.errors }, { status: 400 });
        }
        console.error('Track creation error:', error);
        return NextResponse.json({ error: 'Internal server error' }, { status: 500 });
    }
}
