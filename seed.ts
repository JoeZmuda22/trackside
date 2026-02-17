import { PrismaClient } from '@prisma/client';
import bcrypt from 'bcryptjs';

const prisma = new PrismaClient();

async function main() {
    console.log('🌱 Seeding database...');

    // Create demo user
    const passwordHash = await bcrypt.hash('password123', 12);

    const user = await prisma.user.upsert({
        where: { email: 'demo@trackside.com' },
        update: {},
        create: {
            name: 'Demo Driver',
            email: 'demo@trackside.com',
            passwordHash,
            experience: 'INTERMEDIATE',
        },
    });

    console.log('✅ Created demo user:', user.email);

    // Create a car for the demo user
    const car = await prisma.car.create({
        data: {
            make: 'Nissan',
            model: '350Z',
            year: 2006,
            userId: user.id,
            mods: {
                create: [
                    { name: 'BC Racing BR Coilovers', category: 'SUSPENSION' },
                    { name: 'Tomei Expreme Ti Exhaust', category: 'EXHAUST' },
                    { name: 'Z1 Motorsports Cold Air Intake', category: 'ENGINE' },
                    { name: 'Stoptech ST-40 Big Brake Kit', category: 'BRAKES' },
                    { name: 'Enkei RPF1 18x9.5', category: 'WHEELS_TIRES' },
                ],
            },
        },
    });

    console.log('✅ Created demo car:', car.year, car.make, car.model);

    // Create sample tracks
    const track1 = await prisma.track.create({
        data: {
            name: 'Laguna Seca',
            location: 'Monterey, CA',
            state: 'CA',
            latitude: 36.5754,
            longitude: -121.7627,
            description: 'Iconic road course featuring the famous Corkscrew turn. 2.238 miles, 11 turns, and significant elevation changes make this a challenging and rewarding track.',
            uploadedById: user.id,
            events: {
                create: [
                    { eventType: 'ROADCOURSE' },
                    { eventType: 'DRIFT' },
                ],
            },
            zones: {
                create: [
                    {
                        name: 'The Corkscrew (T8-T8A)',
                        description: 'Famous downhill left-right combo with 5.5 stories of elevation change. Blind entry — use the tree as a braking marker.',
                        posX: 65,
                        posY: 25,
                    },
                    {
                        name: 'Turn 2 (Andretti Hairpin)',
                        description: 'Tight left-hand hairpin. Late apex is key.',
                        posX: 30,
                        posY: 40,
                    },
                    {
                        name: 'Turn 5',
                        description: 'High-speed left sweeper heading uphill. Carry momentum.',
                        posX: 45,
                        posY: 60,
                    },
                ],
            },
        },
        include: { events: true, zones: true },
    });

    console.log('✅ Created track:', track1.name);

    const track2 = await prisma.track.create({
        data: {
            name: 'Atlanta Motorsports Park',
            location: 'Dawsonville, GA',
            state: 'GA',
            latitude: 34.3705,
            longitude: -84.1643,
            description: 'A 2-mile, 16-turn road course with 100ft of elevation change, designed by Hermann Tilke. Features a dedicated drift pad and drag strip.',
            uploadedById: user.id,
            events: {
                create: [
                    { eventType: 'ROADCOURSE' },
                    { eventType: 'DRIFT' },
                    { eventType: 'DRAG' },
                ],
            },
            zones: {
                create: [
                    {
                        name: 'Turn 1',
                        description: 'Fast right-hander after the main straight. Heavy braking zone.',
                        posX: 80,
                        posY: 30,
                    },
                    {
                        name: 'Turn 12 (Rollercoaster)',
                        description: 'Blind crest into a left-right combo. Commitment corner.',
                        posX: 35,
                        posY: 55,
                    },
                ],
            },
        },
        include: { events: true, zones: true },
    });

    console.log('✅ Created track:', track2.name);

    // Create sample zone tips
    for (const zone of track1.zones) {
        await prisma.zoneTip.create({
            data: {
                content: zone.name.includes('Corkscrew')
                    ? 'Use the big tree on the left as your turn-in point. Trust the line and commit — hesitation here is dangerous.'
                    : zone.name.includes('Turn 2')
                        ? 'Brake deep and trail brake in. The car will rotate naturally. Get on power early for the uphill section.'
                        : 'Stay wide and use all the road. The banking helps you carry more speed than you think.',
                conditions: 'DRY',
                zoneId: zone.id,
                authorId: user.id,
            },
        });
    }

    console.log('✅ Created zone tips');

    // Create sample reviews
    const gripEvent1 = track1.events.find((e: any) => e.eventType === 'ROADCOURSE')!;
    await prisma.trackReview.create({
        data: {
            rating: 5,
            content: 'Absolutely world-class track. The Corkscrew is every bit as intense as it looks on TV. The facility is well-maintained and the tech inspection is thorough. Will definitely be back.',
            conditions: 'DRY',
            trackId: track1.id,
            trackEventId: gripEvent1.id,
            authorId: user.id,
        },
    });

    console.log('✅ Created reviews');

    // Create sample lap records
    const gripEvent2 = track2.events.find((e: any) => e.eventType === 'ROADCOURSE')!;
    await prisma.lapRecord.create({
        data: {
            lapTime: '1:42.856',
            conditions: 'DRY',
            notes: 'Best time of the day. Car felt great after adjusting front camber.',
            tirePressureFL: 32.5,
            tirePressureFR: 32.5,
            tirePressureRL: 34.0,
            tirePressureRR: 34.0,
            fuelLevel: 50,
            camberFL: -2.5,
            camberFR: -2.5,
            camberRL: -1.8,
            camberRR: -1.8,
            casterFL: 5.2,
            casterFR: 5.2,
            toeFL: 0.1,
            toeFR: 0.1,
            toeRL: 0.15,
            toeRR: 0.15,
            trackId: track2.id,
            trackEventId: gripEvent2.id,
            carId: car.id,
            driverId: user.id,
        },
    });

    await prisma.lapRecord.create({
        data: {
            lapTime: '1:45.112',
            conditions: 'WET',
            notes: 'Started raining mid-session. Dropped tire pressure to help with wet grip.',
            tirePressureFL: 30.0,
            tirePressureFR: 30.0,
            tirePressureRL: 32.0,
            tirePressureRR: 32.0,
            fuelLevel: 40,
            camberFL: -2.5,
            camberFR: -2.5,
            camberRL: -1.8,
            camberRR: -1.8,
            trackId: track2.id,
            trackEventId: gripEvent2.id,
            carId: car.id,
            driverId: user.id,
        },
    });

    console.log('✅ Created lap records');
    console.log('');
    console.log('🏁 Seed complete! Login with:');
    console.log('   Email: demo@trackside.com');
    console.log('   Password: password123');
}

main()
    .then(async () => {
        await prisma.$disconnect();
    })
    .catch(async (e) => {
        console.error(e);
        await prisma.$disconnect();
        process.exit(1);
    });
