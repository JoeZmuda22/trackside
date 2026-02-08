require('dotenv').config({ path: '.env.local' });
const { PrismaClient } = require('@prisma/client');
const prisma = new PrismaClient();

async function listAllTracks() {
  try {
    const tracks = await prisma.track.findMany({
      select: {
        id: true,
        name: true,
        location: true,
        imageUrl: true,
        createdAt: true,
      },
    });

    if (tracks.length === 0) {
      console.log('❌ No tracks found in database');
      process.exit(0);
    }

    console.log('\n📍 All Tracks in Database:\n');
    tracks.forEach((track, index) => {
      console.log(`${index + 1}. ${track.name}`);
      console.log(`   Location: ${track.location}`);
      console.log(`   Image: ${track.imageUrl ? '✅ ' + track.imageUrl : '❌ Missing'}`);
      console.log(`   ID: ${track.id}`);
      console.log(`   Created: ${track.createdAt}`);
      console.log('');
    });

    process.exit(0);
  } catch (error) {
    console.error('❌ Error:', error.message);
    process.exit(1);
  } finally {
    await prisma.$disconnect();
  }
}

listAllTracks();
