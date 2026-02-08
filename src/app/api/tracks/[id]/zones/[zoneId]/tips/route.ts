import { NextResponse } from 'next/server';
import { auth } from '@/lib/auth';
import { prisma } from '@/lib/db';
import { zoneTipSchema } from '@/lib/validations';

// POST /api/tracks/[id]/zones/[zoneId]/tips — add a tip to a zone
export async function POST(
    request: Request,
    { params }: { params: Promise<{ id: string; zoneId: string }> }
) {
    const session = await auth();
    if (!session?.user?.id) {
        return NextResponse.json({ error: 'Unauthorized' }, { status: 401 });
    }

    const { id, zoneId } = await params;

    // Verify zone belongs to track
    const zone = await prisma.trackZone.findFirst({
        where: { id: zoneId, trackId: id },
    });

    if (!zone) {
        return NextResponse.json({ error: 'Zone not found' }, { status: 404 });
    }

    try {
        const body = await request.json();
        const validated = zoneTipSchema.parse({ ...body, zoneId });

        const tip = await prisma.zoneTip.create({
            data: {
                content: validated.content,
                conditions: validated.conditions,
                zoneId,
                authorId: session.user.id,
            },
            include: {
                author: { select: { id: true, name: true } },
            },
        });

        return NextResponse.json(tip, { status: 201 });
    } catch (error: any) {
        if (error?.name === 'ZodError') {
            return NextResponse.json({ error: 'Validation failed', details: error.errors }, { status: 400 });
        }
        return NextResponse.json({ error: 'Internal server error' }, { status: 500 });
    }
}
