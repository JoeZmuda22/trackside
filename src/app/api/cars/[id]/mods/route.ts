import { NextResponse } from 'next/server';
import { auth } from '@/lib/auth';
import { prisma } from '@/lib/db';
import { carModSchema } from '@/lib/validations';

// POST /api/cars/[id]/mods — add a mod to a car
export async function POST(
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
        const validated = carModSchema.parse({ ...body, carId: id });

        const mod = await prisma.carMod.create({
            data: {
                name: validated.name,
                category: validated.category,
                notes: validated.notes,
                carId: id,
            },
        });

        return NextResponse.json(mod, { status: 201 });
    } catch (error: any) {
        if (error?.name === 'ZodError') {
            return NextResponse.json({ error: 'Validation failed', details: error.errors }, { status: 400 });
        }
        return NextResponse.json({ error: 'Internal server error' }, { status: 500 });
    }
}
