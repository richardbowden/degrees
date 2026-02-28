'use client';

import { useState, useEffect, useCallback } from 'react';
import { useRouter } from 'next/navigation';
import { api, ApiError } from '@/lib/api';
import { VehicleSelect } from '@/components/vehicle-select';
import { BookingCalendar } from '@/components/booking-calendar';
import { TimeSlotGrid } from '@/components/time-slot-grid';
import { VehicleForm, VehicleFormData } from '@/components/vehicle-form';
import type { Vehicle, AvailableSlot, Booking, Cart, DetailingService } from '@/lib/types';

type Step = 'vehicle' | 'date' | 'time' | 'notes';

export default function CheckoutPage() {
  const router = useRouter();
  const [step, setStep] = useState<Step>('vehicle');
  const [vehicles, setVehicles] = useState<Vehicle[]>([]);
  const [cart, setCart] = useState<Cart | null>(null);
  const [selectedVehicle, setSelectedVehicle] = useState('');
  const [selectedDate, setSelectedDate] = useState<string | null>(null);
  const [selectedTime, setSelectedTime] = useState<string | null>(null);
  const [notes, setNotes] = useState('');
  const [slots, setSlots] = useState<AvailableSlot[]>([]);
  const [loading, setLoading] = useState(true);
  const [slotsLoading, setSlotsLoading] = useState(false);
  const [submitting, setSubmitting] = useState(false);
  const [showAddVehicle, setShowAddVehicle] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const [serviceDurations, setServiceDurations] = useState<Record<string, number>>({});

  useEffect(() => {
    async function load() {
      try {
        const [vehiclesData, cartData, catalogueData] = await Promise.all([
          api<{ vehicles: Vehicle[] }>('/me/vehicles'),
          api<{ cart: Cart }>('/cart'),
          api<{ services: DetailingService[] }>('/catalogue'),
        ]);
        setVehicles(vehiclesData.vehicles ?? []);
        setCart(cartData.cart);
        const durations: Record<string, number> = {};
        for (const svc of catalogueData.services ?? []) {
          durations[svc.id] = svc.durationMinutes;
        }
        setServiceDurations(durations);
      } catch (err: unknown) {
        const apiErr = err as ApiError;
        setError(apiErr.detail || 'Failed to load checkout data');
      } finally {
        setLoading(false);
      }
    }
    load();
  }, []);

  const totalDuration = cart?.items?.reduce((sum, item) => {
    const duration = serviceDurations[item.serviceId] || 60;
    return sum + duration * Number(item.quantity);
  }, 0) ?? 60;

  const fetchSlots = useCallback(async (date: string) => {
    setSlotsLoading(true);
    setSlots([]);
    setSelectedTime(null);
    try {
      const data = await api<{ slots: AvailableSlot[] }>(
        `/checkout/available-slots?date=${date}&durationMinutes=${totalDuration}`
      );
      setSlots(data.slots ?? []);
    } catch (err: unknown) {
      const apiErr = err as ApiError;
      setError(apiErr.detail || 'Failed to load time slots');
    } finally {
      setSlotsLoading(false);
    }
  }, [totalDuration]);

  async function handleAddVehicle(data: VehicleFormData) {
    setError(null);
    try {
      const res = await api<{ vehicle: Vehicle }>('/me/vehicles', {
        method: 'POST',
        body: data,
      });
      setVehicles(prev => [...prev, res.vehicle]);
      setShowAddVehicle(false);
      setSelectedVehicle(res.vehicle.id);
      setStep('date');
    } catch (err: unknown) {
      const apiErr = err as ApiError;
      setError(apiErr.detail || 'Failed to add vehicle');
    }
  }

  function handleVehicleSelect(vehicleId: string) {
    setSelectedVehicle(vehicleId);
    setStep('date');
  }

  function handleDateSelect(date: string) {
    setSelectedDate(date);
    setStep('time');
    fetchSlots(date);
  }

  function handleTimeSelect(time: string) {
    setSelectedTime(time);
    setStep('notes');
  }

  async function handleSubmit() {
    if (!selectedVehicle || !selectedDate || !selectedTime) return;
    setSubmitting(true);
    setError(null);
    try {
      const data = await api<{ booking: Booking }>('/checkout', {
        method: 'POST',
        body: {
          vehicleId: selectedVehicle,
          scheduledDate: selectedDate,
          scheduledTime: selectedTime,
          notes,
        },
      });
      router.push(`/checkout/deposit?booking_id=${data.booking.id}`);
    } catch (err: unknown) {
      const apiErr = err as ApiError;
      setError(apiErr.detail || 'Failed to create booking');
      setSubmitting(false);
    }
  }

  if (loading) {
    return (
      <div className="max-w-2xl mx-auto px-4 py-16 text-center">
        <p className="text-text-muted">Loading checkout...</p>
      </div>
    );
  }

  const steps: { key: Step; label: string; number: number }[] = [
    { key: 'vehicle', label: 'Vehicle', number: 1 },
    { key: 'date', label: 'Date', number: 2 },
    { key: 'time', label: 'Time', number: 3 },
    { key: 'notes', label: 'Confirm', number: 4 },
  ];

  const stepIndex = steps.findIndex(s => s.key === step);

  return (
    <div className="max-w-2xl mx-auto px-4 py-8">
      <h1 className="text-2xl font-bold mb-8">Checkout</h1>

      <div className="flex items-center gap-2 mb-8">
        {steps.map((s, i) => (
          <div key={s.key} className="flex items-center">
            <button
              type="button"
              onClick={() => i <= stepIndex && setStep(s.key)}
              disabled={i > stepIndex}
              className={`flex items-center gap-2 text-sm font-medium ${
                i <= stepIndex ? 'text-foreground' : 'text-text-muted'
              }`}
            >
              <span
                className={`w-6 h-6 rounded-full flex items-center justify-center text-xs ${
                  i < stepIndex
                    ? 'bg-green-500 text-white'
                    : i === stepIndex
                      ? 'bg-brand-500 text-white'
                      : 'bg-surface-raised text-text-muted'
                }`}
              >
                {i < stepIndex ? '\u2713' : s.number}
              </span>
              <span className="hidden sm:inline">{s.label}</span>
            </button>
            {i < steps.length - 1 && (
              <div className="w-8 h-px bg-border-subtle mx-2" />
            )}
          </div>
        ))}
      </div>

      {error && (
        <p className="text-red-400 text-sm mb-4">{error}</p>
      )}

      {step === 'vehicle' && (
        <div>
          <h2 className="text-lg font-semibold mb-4">Select Vehicle</h2>
          {showAddVehicle ? (
            <VehicleForm
              onSubmit={handleAddVehicle}
              onCancel={() => setShowAddVehicle(false)}
            />
          ) : (
            <>
              <VehicleSelect
                vehicles={vehicles}
                selectedId={selectedVehicle}
                onChange={handleVehicleSelect}
              />
              <button
                type="button"
                onClick={() => setShowAddVehicle(true)}
                className="mt-4 text-sm font-medium text-brand-400 underline hover:text-brand-500"
              >
                + Add a new vehicle
              </button>
            </>
          )}
        </div>
      )}

      {step === 'date' && (
        <div>
          <h2 className="text-lg font-semibold mb-4">Pick a Date</h2>
          <BookingCalendar
            selectedDate={selectedDate}
            onSelect={handleDateSelect}
          />
        </div>
      )}

      {step === 'time' && (
        <div>
          <h2 className="text-lg font-semibold mb-4">Pick a Time</h2>
          {slotsLoading ? (
            <p className="text-text-muted text-sm">Loading available slots...</p>
          ) : (
            <TimeSlotGrid
              slots={slots}
              selectedTime={selectedTime}
              onSelect={handleTimeSelect}
            />
          )}
        </div>
      )}

      {step === 'notes' && (
        <div>
          <h2 className="text-lg font-semibold mb-4">Any Special Requests?</h2>
          <textarea
            value={notes}
            onChange={e => setNotes(e.target.value)}
            rows={4}
            placeholder="Parking instructions, specific concerns, etc."
            className="w-full bg-surface-input border border-border-subtle rounded px-3 py-2 text-sm text-foreground focus:outline-none focus:ring-2 focus:ring-brand-500 mb-6"
          />
          <button
            type="button"
            onClick={handleSubmit}
            disabled={submitting}
            className="w-full btn-brand py-3"
          >
            {submitting ? 'Creating booking...' : 'Create Booking'}
          </button>
        </div>
      )}
    </div>
  );
}
