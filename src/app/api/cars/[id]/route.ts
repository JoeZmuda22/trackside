import { NextResponse } from 'next/server';
import { auth } from '@/lib/auth';
import { prisma } from '@/lib/db';
import { carSchema } from '@/lib/validations';

// PUT /api/cars/[id] — update a car
export async function PUT(
    request: Request,
    { params }: { params: Promise<{ id: string }> }
) {
    const session = await auth();
    if (!session?.user?.id) {
        return NextResponse.json({ error: 'Unauthorized' }, { status: 401 });
    }

    const { id } = await params;

    const car = await prisma.car.findFirst({
        where: { id, userId: session.user.id },
    });

    if (!car) {
        return NextResponse.json({ error: 'Car not found' }, { status: 404 });
    }

    try {
        const body = await request.json();
        const validated = carSchema.parse(body);

        const updated = await prisma.car.update({
            where: { id },
            data: validated,
            include: { mods: true },
        });

        return NextResponse.json(updated);
    } catch (error: any) {
        if (error?.name === 'ZodError') {
            return NextResponse.json({ error: 'Validation failed', details: error.errors }, { status: 400 });
        }
        return NextResponse.json({ error: 'Internal server error' }, { status: 500 });
    }
}

// DELETE /api/cars/[id] — delete a car
export async function DELETE(
    _request: Request,
    { params }: { params: Promise<{ id: string }> }
) {
    const session = await auth();
    if (!session?.user?.id) {
        return NextResponse.json({ error: 'Unauthorized' }, { status: 401 });
    }

    const { id } = await params;

    const car = await prisma.car.findFirst({
        where: { id, userId: session.user.id },
    });

    if (!car) {
        return NextResponse.json({ error: 'Car not found' }, { status: 404 });
    }

    await prisma.car.delete({ where: { id } });

    return NextResponse.json({ success: true });
}
