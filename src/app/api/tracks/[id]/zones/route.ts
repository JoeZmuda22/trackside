import { NextResponse } from 'next/server';
import { auth } from '@/lib/auth';
import { prisma } from '@/lib/db';
import { trackZoneSchema } from '@/lib/validations';

// POST /api/tracks/[id]/zones — add a zone to a track
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
        const validated = trackZoneSchema.parse({ ...body, trackId: id });

        const zone = await prisma.trackZone.create({
            data: {
                name: validated.name,
                description: validated.description,
                posX: validated.posX,
                posY: validated.posY,
                trackId: id,
            },
            include: {
                tips: {
                    include: { author: { select: { id: true, name: true } } },
                },
            },
        });

        return NextResponse.json(zone, { status: 201 });
    } catch (error: any) {
        if (error?.name === 'ZodError') {
            return NextResponse.json({ error: 'Validation failed', details: error.errors }, { status: 400 });
        }
        return NextResponse.json({ error: 'Internal server error' }, { status: 500 });
    }
}
