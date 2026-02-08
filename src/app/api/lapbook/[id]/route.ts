import { NextResponse } from 'next/server';
import { auth } from '@/lib/auth';
import { prisma } from '@/lib/db';

// DELETE /api/lapbook/[id] — delete a lap record
export async function DELETE(
    _request: Request,
    { params }: { params: Promise<{ id: string }> }
) {
    const session = await auth();
    if (!session?.user?.id) {
        return NextResponse.json({ error: 'Unauthorized' }, { status: 401 });
    }

    const { id } = await params;

    const record = await prisma.lapRecord.findFirst({
        where: { id, driverId: session.user.id },
    });

    if (!record) {
        return NextResponse.json({ error: 'Record not found' }, { status: 404 });
    }

    await prisma.lapRecord.delete({ where: { id } });

    return NextResponse.json({ success: true });
}
