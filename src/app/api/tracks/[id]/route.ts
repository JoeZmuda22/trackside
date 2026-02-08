import { NextResponse } from 'next/server';
import { prisma } from '@/lib/db';
import { auth } from '@/lib/auth';

// GET /api/tracks/[id] — get a single track with all details
export async function GET(
    _request: Request,
    { params }: { params: Promise<{ id: string }> }
) {
    const { id } = await params;

    const track = await prisma.track.findUnique({
        where: { id },
        include: {
            events: true,
            uploadedBy: { select: { id: true, name: true, experience: true } },
            zones: {
                include: {
                    tips: {
                        include: {
                            author: { select: { id: true, name: true } },
                        },
                        orderBy: { createdAt: 'desc' },
                    },
                },
            },
            reviews: {
                include: {
                    author: {
                        select: {
                            id: true,
                            name: true,
                            experience: true,
                            cars: { select: { make: true, model: true, year: true } },
                        },
                    },
                    trackEvent: true,
                },
                orderBy: { createdAt: 'desc' },
            },
            _count: {
                select: { reviews: true, zones: true, lapRecords: true },
            },
        },
    });

    if (!track) {
        return NextResponse.json({ error: 'Track not found' }, { status: 404 });
    }

    // Calculate average rating
    const avgRating = await prisma.trackReview.aggregate({
        where: { trackId: id },
        _avg: { rating: true },
    });

    return NextResponse.json({
        ...track,
        avgRating: avgRating._avg.rating || 0,
    });
}

// PATCH /api/tracks/[id] — update track (e.g. image, name, description)
export async function PATCH(
    request: Request,
    { params }: { params: Promise<{ id: string }> }
) {
    const session = await auth();
    if (!session?.user?.id) {
        return NextResponse.json({ error: 'Unauthorized' }, { status: 401 });
    }

    const { id } = await params;

    // Verify the track exists and belongs to the user
    const track = await prisma.track.findUnique({
        where: { id },
        select: { uploadedById: true },
    });

    if (!track) {
        return NextResponse.json({ error: 'Track not found' }, { status: 404 });
    }

    if (track.uploadedById !== session.user.id) {
        return NextResponse.json({ error: 'You can only edit your own tracks' }, { status: 403 });
    }

    const body = await request.json();
    const { imageUrl, name, description, location } = body;

    const updated = await prisma.track.update({
        where: { id },
        data: {
            ...(imageUrl !== undefined && { imageUrl }),
            ...(name !== undefined && { name }),
            ...(description !== undefined && { description }),
            ...(location !== undefined && { location }),
        },
    });

    return NextResponse.json(updated);
}
