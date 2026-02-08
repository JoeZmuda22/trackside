import { NextResponse } from 'next/server';
import { auth } from '@/lib/auth';
import { prisma } from '@/lib/db';
import { profileSchema } from '@/lib/validations';

// GET /api/profile — get current user's profile
export async function GET() {
    const session = await auth();
    if (!session?.user?.id) {
        return NextResponse.json({ error: 'Unauthorized' }, { status: 401 });
    }

    const user = await prisma.user.findUnique({
        where: { id: session.user.id },
        select: {
            id: true,
            name: true,
            email: true,
            experience: true,
            image: true,
            createdAt: true,
            cars: {
                include: { mods: true },
                orderBy: { createdAt: 'desc' },
            },
            _count: {
                select: {
                    trackReviews: true,
                    lapRecords: true,
                    tracks: true,
                    zoneTips: true,
                },
            },
        },
    });

    if (!user) {
        return NextResponse.json({ error: 'User not found' }, { status: 404 });
    }

    return NextResponse.json(user);
}

// PUT /api/profile — update current user's profile
export async function PUT(request: Request) {
    const session = await auth();
    if (!session?.user?.id) {
        return NextResponse.json({ error: 'Unauthorized' }, { status: 401 });
    }

    try {
        const body = await request.json();
        const validated = profileSchema.parse(body);

        const user = await prisma.user.update({
            where: { id: session.user.id },
            data: {
                name: validated.name,
                experience: validated.experience,
            },
            select: {
                id: true,
                name: true,
                email: true,
                experience: true,
            },
        });

        return NextResponse.json(user);
    } catch (error: any) {
        if (error?.name === 'ZodError') {
            return NextResponse.json({ error: 'Validation failed', details: error.errors }, { status: 400 });
        }
        return NextResponse.json({ error: 'Internal server error' }, { status: 500 });
    }
}
