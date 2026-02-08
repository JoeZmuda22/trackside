import { NextResponse } from 'next/server';
import { auth } from '@/lib/auth';
import { prisma } from '@/lib/db';
import { trackSchema } from '@/lib/validations';

// GET /api/tracks — list all tracks with filters
export async function GET(request: Request) {
    const { searchParams } = new URL(request.url);
    const search = searchParams.get('search') || '';
    const eventType = searchParams.get('eventType');

    const where: any = {
        status: 'APPROVED',
    };

    if (search) {
        where.OR = [
            { name: { contains: search, mode: 'insensitive' } },
            { location: { contains: search, mode: 'insensitive' } },
        ];
    }

    if (eventType) {
        where.events = {
            some: { eventType: eventType },
        };
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

    // Calculate average rating for each track
    const tracksWithRating = await Promise.all(
        tracks.map(async (track) => {
            const avgRating = await prisma.trackReview.aggregate({
                where: { trackId: track.id },
                _avg: { rating: true },
            });
            return {
                ...track,
                avgRating: avgRating._avg.rating || 0,
            };
        })
    );

    return NextResponse.json(tracksWithRating);
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
