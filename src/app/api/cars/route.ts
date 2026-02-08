import { NextResponse } from 'next/server';
import { auth } from '@/lib/auth';
import { prisma } from '@/lib/db';
import { carSchema } from '@/lib/validations';

// GET /api/cars — get current user's cars
export async function GET() {
    const session = await auth();
    if (!session?.user?.id) {
        return NextResponse.json({ error: 'Unauthorized' }, { status: 401 });
    }

    const cars = await prisma.car.findMany({
        where: { userId: session.user.id },
        include: { mods: true },
        orderBy: { createdAt: 'desc' },
    });

    return NextResponse.json(cars);
}

// POST /api/cars — add a new car
export async function POST(request: Request) {
    const session = await auth();
    if (!session?.user?.id) {
        return NextResponse.json({ error: 'Unauthorized' }, { status: 401 });
    }

    try {
        const body = await request.json();
        const validated = carSchema.parse(body);

        const car = await prisma.car.create({
            data: {
                ...validated,
                userId: session.user.id,
            },
            include: { mods: true },
        });

        return NextResponse.json(car, { status: 201 });
    } catch (error: any) {
        if (error?.name === 'ZodError') {
            return NextResponse.json({ error: 'Validation failed', details: error.errors }, { status: 400 });
        }
        return NextResponse.json({ error: 'Internal server error' }, { status: 500 });
    }
}
