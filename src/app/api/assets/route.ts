import { NextResponse } from 'next/server';
import { readdir } from 'fs/promises';
import { join } from 'path';

export async function GET() {
  try {
    const assetsPath = join(process.cwd(), 'public', 'assets');
    const files = await readdir(assetsPath);
    
    // Filter for image files
    const imageExtensions = ['.jpg', '.jpeg', '.png', '.gif', '.webp', '.svg'];
    const imageFiles = files.filter(file => 
      imageExtensions.some(ext => file.toLowerCase().endsWith(ext))
    );
    
    // Return the image paths
    const imagePaths = imageFiles.map(file => `/assets/${file}`);
    
    return NextResponse.json({ images: imagePaths });
  } catch (error) {
    console.error('Error reading assets directory:', error);
    return NextResponse.json(
      { error: 'Failed to load assets' },
      { status: 500 }
    );
  }
}