'use client';

import { cn, getConditionBadgeColor, getEventBadgeColor } from '@/lib/utils';

interface BadgeProps {
    children: React.ReactNode;
    variant?: 'default' | 'condition' | 'event';
    value?: string;
    className?: string;
}

export function Badge({ children, variant = 'default', value, className }: BadgeProps) {
    let colorClass = 'bg-surface-700 text-surface-300';

    if (variant === 'condition' && value) {
        colorClass = getConditionBadgeColor(value);
    } else if (variant === 'event' && value) {
        colorClass = getEventBadgeColor(value);
    }

    return (
        <span className={cn('badge', colorClass, className)}>
            {children}
        </span>
    );
}
