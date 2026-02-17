import { type ClassValue, clsx } from 'clsx';
import { twMerge } from 'tailwind-merge';

export function cn(...inputs: ClassValue[]) {
    return twMerge(clsx(inputs));
}

export function getConditionBadgeColor(condition: string): string {
    return condition === 'WET'
        ? 'bg-blue-100 text-blue-800'
        : 'bg-amber-100 text-amber-800';
}

export function getEventBadgeColor(eventType: string): string {
    switch (eventType) {
        case 'DRIFT':
            return 'bg-purple-100 text-purple-800';
        case 'DRAG':
            return 'bg-red-100 text-red-800';
        case 'GRIP':
            return 'bg-green-100 text-green-800';
        default:
            return 'bg-gray-100 text-gray-800';
    }
}

export function getExperienceLabel(level: string): string {
    switch (level) {
        case 'BEGINNER':
            return 'Beginner';
        case 'INTERMEDIATE':
            return 'Intermediate';
        case 'ADVANCED':
            return 'Advanced';
        case 'PRO':
            return 'Pro';
        default:
            return level;
    }
}

export function getEventLabel(eventType: string): string {
    switch (eventType) {
        case 'AUTOCROSS':
            return 'Autocross';
        case 'ROADCOURSE':
            return 'Road Course';
        case 'DRIFT':
            return 'Drift';
        case 'DRAG':
            return 'Drag';
        default:
            return eventType;
    }
}
