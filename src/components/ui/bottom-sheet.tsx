'use client';

import { cn } from '@/lib/utils';

interface BottomSheetProps {
    isOpen: boolean;
    onClose: () => void;
    title?: string;
    children: React.ReactNode;
}

export function BottomSheet({ isOpen, onClose, title, children }: BottomSheetProps) {
    if (!isOpen) return null;

    return (
        <>
            <div className="bottom-sheet-overlay" onClick={onClose} />
            <div className="bottom-sheet">
                <div className="sticky top-0 bg-surface-800 border-b border-surface-700 px-4 pt-3 pb-3 z-10">
                    <div className="mx-auto mb-3 h-1 w-10 rounded-full bg-surface-600" />
                    {title && (
                        <div className="flex items-center justify-between">
                            <h3 className="text-lg font-semibold text-white">{title}</h3>
                            <button
                                onClick={onClose}
                                className="p-1 text-surface-400 hover:text-white transition-colors"
                            >
                                <svg className="w-5 h-5" fill="none" viewBox="0 0 24 24" strokeWidth={2} stroke="currentColor">
                                    <path strokeLinecap="round" strokeLinejoin="round" d="M6 18 18 6M6 6l12 12" />
                                </svg>
                            </button>
                        </div>
                    )}
                </div>
                <div className="p-4">{children}</div>
            </div>
        </>
    );
}
