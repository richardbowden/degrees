'use client';

import { useState } from 'react';
import { api, ApiError } from '@/lib/api';
import type { CustomerProfile, Vehicle } from '@/lib/types';
import { VehicleForm, VehicleFormData } from '@/components/vehicle-form';

interface ProfileClientProps {
  initialProfile: CustomerProfile;
  initialVehicles: Vehicle[];
}

export function ProfileClient({ initialProfile, initialVehicles }: ProfileClientProps) {
  const [profile, setProfile] = useState(initialProfile);
  const [vehicles, setVehicles] = useState(initialVehicles);
  const [phone, setPhone] = useState(profile.phone);
  const [address, setAddress] = useState(profile.address);
  const [suburb, setSuburb] = useState(profile.suburb);
  const [postcode, setPostcode] = useState(profile.postcode);
  const [saving, setSaving] = useState(false);
  const [profileMessage, setProfileMessage] = useState('');
  const [profileError, setProfileError] = useState('');

  const [editingVehicle, setEditingVehicle] = useState<Vehicle | null>(null);
  const [showAddVehicle, setShowAddVehicle] = useState(false);
  const [vehicleLoading, setVehicleLoading] = useState<string | null>(null);
  const [vehicleError, setVehicleError] = useState('');

  async function saveProfile(e: React.FormEvent) {
    e.preventDefault();
    setSaving(true);
    setProfileMessage('');
    setProfileError('');
    try {
      const res = await api<{ profile: CustomerProfile }>('/me/profile', {
        method: 'PUT',
        body: { phone, address, suburb, postcode },
      });
      setProfile(res.profile);
      setProfileMessage('Profile updated successfully.');
    } catch (err: unknown) {
      const apiErr = err as ApiError;
      setProfileError(apiErr.detail || 'Failed to update profile.');
    } finally {
      setSaving(false);
    }
  }

  async function handleAddVehicle(data: VehicleFormData) {
    const res = await api<{ vehicle: Vehicle }>('/me/vehicles', {
      method: 'POST',
      body: data,
    });
    setVehicles(prev => [...prev, res.vehicle]);
    setShowAddVehicle(false);
  }

  async function handleEditVehicle(data: VehicleFormData) {
    if (!editingVehicle) return;
    const res = await api<{ vehicle: Vehicle }>(`/me/vehicles/${editingVehicle.id}`, {
      method: 'PUT',
      body: data,
    });
    setVehicles(prev => prev.map(v => v.id === editingVehicle.id ? res.vehicle : v));
    setEditingVehicle(null);
  }

  async function handleDeleteVehicle(id: string) {
    if (!confirm('Are you sure you want to delete this vehicle?')) return;
    setVehicleLoading(id);
    setVehicleError('');
    try {
      await api(`/me/vehicles/${id}`, { method: 'DELETE' });
      setVehicles(prev => prev.filter(v => v.id !== id));
    } catch (err: unknown) {
      const apiErr = err as ApiError;
      setVehicleError(apiErr.detail || 'Failed to delete vehicle.');
    } finally {
      setVehicleLoading(null);
    }
  }

  const inputClass = 'w-full border border-gray-300 rounded-md px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-gray-900';

  return (
    <div className="space-y-10">
      {/* Profile form */}
      <section>
        <h2 className="text-lg font-semibold text-gray-900 mb-4">Contact Details</h2>
        <form onSubmit={saveProfile} className="max-w-lg space-y-4">
          {profileMessage && <p className="text-sm text-green-600">{profileMessage}</p>}
          {profileError && <p className="text-sm text-red-600">{profileError}</p>}
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">Phone</label>
            <input
              type="tel"
              value={phone}
              onChange={e => setPhone(e.target.value)}
              className={inputClass}
            />
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">Address</label>
            <input
              type="text"
              value={address}
              onChange={e => setAddress(e.target.value)}
              className={inputClass}
            />
          </div>
          <div className="grid grid-cols-2 gap-4">
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">Suburb</label>
              <input
                type="text"
                value={suburb}
                onChange={e => setSuburb(e.target.value)}
                className={inputClass}
              />
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">Postcode</label>
              <input
                type="text"
                value={postcode}
                onChange={e => setPostcode(e.target.value)}
                className={inputClass}
              />
            </div>
          </div>
          <button
            type="submit"
            disabled={saving}
            className="px-4 py-2 bg-gray-900 text-white text-sm font-medium rounded-md hover:bg-gray-800 disabled:opacity-50"
          >
            {saving ? 'Saving...' : 'Save Profile'}
          </button>
        </form>
      </section>

      {/* Vehicles */}
      <section>
        <div className="flex items-center justify-between mb-4">
          <h2 className="text-lg font-semibold text-gray-900">Vehicles</h2>
          {!showAddVehicle && !editingVehicle && (
            <button
              onClick={() => setShowAddVehicle(true)}
              className="px-3 py-1.5 bg-gray-900 text-white text-sm font-medium rounded-md hover:bg-gray-800"
            >
              Add Vehicle
            </button>
          )}
        </div>

        {vehicleError && <p className="text-sm text-red-600 mb-4">{vehicleError}</p>}

        {showAddVehicle && (
          <div className="mb-6">
            <h3 className="text-sm font-medium text-gray-700 mb-3">New Vehicle</h3>
            <VehicleForm
              onSubmit={handleAddVehicle}
              onCancel={() => setShowAddVehicle(false)}
            />
          </div>
        )}

        {editingVehicle && (
          <div className="mb-6">
            <h3 className="text-sm font-medium text-gray-700 mb-3">Edit Vehicle</h3>
            <VehicleForm
              vehicle={editingVehicle}
              onSubmit={handleEditVehicle}
              onCancel={() => setEditingVehicle(null)}
            />
          </div>
        )}

        {vehicles.length === 0 && !showAddVehicle ? (
          <p className="text-gray-500 text-sm">No vehicles added yet.</p>
        ) : (
          <div className="space-y-3">
            {vehicles.map(vehicle => (
              <div
                key={vehicle.id}
                className="border border-gray-200 rounded-lg p-4 flex items-center justify-between"
              >
                <div>
                  <p className="font-medium text-gray-900">
                    {vehicle.year} {vehicle.make} {vehicle.model}
                    {vehicle.isPrimary && (
                      <span className="ml-2 text-xs bg-gray-100 text-gray-600 px-2 py-0.5 rounded-full">
                        Primary
                      </span>
                    )}
                  </p>
                  <p className="text-sm text-gray-500">
                    {vehicle.colour}
                    {vehicle.rego && ` \u00B7 ${vehicle.rego}`}
                    {vehicle.paintType && ` \u00B7 ${vehicle.paintType}`}
                  </p>
                  {vehicle.conditionNotes && (
                    <p className="text-xs text-gray-400 mt-1">{vehicle.conditionNotes}</p>
                  )}
                </div>
                <div className="flex gap-2">
                  <button
                    onClick={() => { setEditingVehicle(vehicle); setShowAddVehicle(false); }}
                    className="text-sm text-gray-600 hover:text-gray-900"
                  >
                    Edit
                  </button>
                  <button
                    onClick={() => handleDeleteVehicle(vehicle.id)}
                    disabled={vehicleLoading === vehicle.id}
                    className="text-sm text-red-600 hover:text-red-800 disabled:opacity-50"
                  >
                    {vehicleLoading === vehicle.id ? 'Deleting...' : 'Delete'}
                  </button>
                </div>
              </div>
            ))}
          </div>
        )}
      </section>
    </div>
  );
}
