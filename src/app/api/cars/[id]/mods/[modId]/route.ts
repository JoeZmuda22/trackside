import { NextResponse } from 'next/server';
import { auth } from '@/lib/auth';
import { prisma } from '@/lib/db';

// DELETE /api/cars/[id]/mods/[modId] — delete a mod
export async function DELETE(
    _request: Request,
    { params }: { params: Promise<{ id: string; modId: string }> }
) {
    const session = await auth();
    if (!session?.user?.id) {
        return NextResponse.json({ error: 'Unauthorized' }, { status: 401 });
    }

    const { id, modId } = await params;

    // Verify car belongs to user
    const car = await prisma.car.findFirst({
        where: { id, userId: session.user.id },
    });

    if (!car) {
        return NextResponse.json({ error: 'Car not found' }, { status: 404 });
    }

    const mod = await prisma.carMod.findFirst({
        where: { id: modId, carId: id },
    });

    if (!mod) {
        return NextResponse.json({ error: 'Mod not found' }, { status: 404 });
    }

    await prisma.carMod.delete({ where: { id: modId } });

    return NextResponse.json({ success: true });
}
