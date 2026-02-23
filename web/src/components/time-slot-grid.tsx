'use client';

import { AvailableSlot } from '@/lib/types';
import { formatTime } from '@/lib/format';

interface TimeSlotGridProps {
  slots: AvailableSlot[];
  selectedTime: string | null;
  onSelect: (time: string) => void;
}

export function TimeSlotGrid({ slots, selectedTime, onSelect }: TimeSlotGridProps) {
  if (slots.length === 0) {
    return <p className="text-sm text-text-muted">No available time slots for this date.</p>;
  }

  return (
    <div className="grid grid-cols-3 sm:grid-cols-4 gap-2">
      {slots.map(slot => {
        const isSelected = slot.time === selectedTime;
        return (
          <button
            key={slot.time}
            type="button"
            onClick={() => onSelect(slot.time)}
            className={`px-3 py-2 rounded text-sm border ${
              isSelected
                ? 'bg-brand-500 text-white border-brand-500'
                : 'border-border-subtle text-text-secondary hover:border-brand-400'
            }`}
          >
            {formatTime(slot.time)}
          </button>
        );
      })}
    </div>
  );
}
