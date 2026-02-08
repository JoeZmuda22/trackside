import { NextResponse } from 'next/server';
import bcrypt from 'bcryptjs';
import { prisma } from '@/lib/db';
import { registerSchema } from '@/lib/validations';

export async function POST(request: Request) {
    try {
        const body = await request.json();
        const validated = registerSchema.parse(body);

        // Check if user already exists
        const existingUser = await prisma.user.findUnique({
            where: { email: validated.email },
        });

        if (existingUser) {
            return NextResponse.json(
                { error: 'Email already registered' },
                { status: 409 }
            );
        }

        const passwordHash = await bcrypt.hash(validated.password, 12);

        const user = await prisma.user.create({
            data: {
                name: validated.name,
                email: validated.email,
                passwordHash,
            },
            select: {
                id: true,
                name: true,
                email: true,
            },
        });

        return NextResponse.json(user, { status: 201 });
    } catch (error: any) {
        if (error?.name === 'ZodError') {
            return NextResponse.json(
                { error: 'Validation failed', details: error.errors },
                { status: 400 }
            );
        }
        console.error('Registration error:', error);
        return NextResponse.json(
            { error: 'Internal server error' },
            { status: 500 }
        );
    }
}
