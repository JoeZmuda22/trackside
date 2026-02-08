'use client';

import { cn } from '@/lib/utils';

export function LoadingSpinner({ className }: { className?: string }) {
    return (
        <div className={cn('flex items-center justify-center py-8', className)}>
            <div className="h-8 w-8 animate-spin rounded-full border-2 border-surface-600 border-t-brand-500" />
        </div>
    );
}

export function EmptyState({
    icon,
    title,
    description,
    action,
}: {
    icon?: React.ReactNode;
    title: string;
    description?: string;
    action?: React.ReactNode;
}) {
    return (
        <div className="flex flex-col items-center justify-center py-12 px-4 text-center">
            {icon && <div className="mb-4 text-surface-500">{icon}</div>}
            <h3 className="text-lg font-semibold text-surface-300">{title}</h3>
            {description && (
                <p className="mt-1 text-sm text-surface-500 max-w-xs">{description}</p>
            )}
            {action && <div className="mt-4">{action}</div>}
        </div>
    );
}
