import { NextResponse } from 'next/server';
import { auth } from '@/lib/auth';
import { prisma } from '@/lib/db';
import { trackReviewSchema } from '@/lib/validations';

// POST /api/tracks/[id]/reviews — add a review to a track
export async function POST(
    request: Request,
    { params }: { params: Promise<{ id: string }> }
) {
    const session = await auth();
    if (!session?.user?.id) {
        return NextResponse.json({ error: 'Unauthorized' }, { status: 401 });
    }

    const { id } = await params;

    const track = await prisma.track.findUnique({ where: { id } });
    if (!track) {
        return NextResponse.json({ error: 'Track not found' }, { status: 404 });
    }

    try {
        const body = await request.json();
        const validated = trackReviewSchema.parse({ ...body, trackId: id });

        const review = await prisma.trackReview.create({
            data: {
                rating: validated.rating,
                content: validated.content,
                conditions: validated.conditions,
                trackId: id,
                trackEventId: validated.trackEventId || null,
                authorId: session.user.id,
            },
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
        });

        return NextResponse.json(review, { status: 201 });
    } catch (error: any) {
        if (error?.name === 'ZodError') {
            return NextResponse.json({ error: 'Validation failed', details: error.errors }, { status: 400 });
        }
        return NextResponse.json({ error: 'Internal server error' }, { status: 500 });
    }
}
