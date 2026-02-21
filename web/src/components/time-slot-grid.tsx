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
    return <p className="text-sm text-gray-500">No available time slots for this date.</p>;
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
                ? 'bg-gray-900 text-white border-gray-900'
                : 'border-gray-300 text-gray-700 hover:border-gray-500'
            }`}
          >
            {formatTime(slot.time)}
          </button>
        );
      })}
    </div>
  );
}
