require('dotenv').config({ path: '.env.local' });
const { PrismaClient } = require('@prisma/client');
const prisma = new PrismaClient();

async function deleteTracksByIds() {
  try {
    const idsToDelete = [
      'cmlbf0wjw0009hgenxb3gqalb', // Laguna Seca (with image)
      'cmlbf0wjy000ghgen9vulrjes', // Atlanta Motorsports Park (with image)
    ];

    console.log('🗑️  Deleting tracks...\n');

    for (const id of idsToDelete) {
      const track = await prisma.track.findUnique({
        where: { id },
      });

      if (track) {
        await prisma.track.delete({
          where: { id },
        });
        console.log(`✅ Deleted: ${track.name} (${track.location})`);
      } else {
        console.log(`⚠️  Track ${id} not found`);
      }
    }

    console.log('\n✨ Cleanup complete!');
    process.exit(0);
  } catch (error) {
    console.error('❌ Error:', error.message);
    process.exit(1);
  } finally {
    await prisma.$disconnect();
  }
}

deleteTracksByIds();
