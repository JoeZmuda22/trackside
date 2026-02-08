import type { Metadata, Viewport } from 'next';
import './globals.css';
import { BottomNav } from '@/components/layout/bottom-nav';
import { AuthProvider } from '@/components/providers/auth-provider';

export const metadata: Metadata = {
    title: 'Trackside — Motorsports Track Reviews',
    description: 'Review tracks, log laps, and share tips with fellow drivers.',
    manifest: '/manifest.json',
};

export const viewport: Viewport = {
    width: 'device-width',
    initialScale: 1,
    maximumScale: 1,
    userScalable: false,
    themeColor: '#0f172a',
};

export default function RootLayout({
    children,
}: {
    children: React.ReactNode;
}) {
    return (
        <html lang="en">
            <body className="font-sans">
                <AuthProvider>
                    <main className="min-h-screen pb-20">
                        {children}
                    </main>
                    <BottomNav />
                </AuthProvider>
            </body>
        </html>
    );
}
