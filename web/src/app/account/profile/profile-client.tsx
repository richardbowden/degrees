'use client';

import { useState } from 'react';
import Link from 'next/link';
import { api, ApiError } from '@/lib/api';
import type { CustomerProfile } from '@/lib/types';

interface ProfileClientProps {
  initialProfile: CustomerProfile;
}

export function ProfileClient({ initialProfile }: ProfileClientProps) {
  const [profile, setProfile] = useState(initialProfile);
  const [phone, setPhone] = useState(profile.phone);
  const [address, setAddress] = useState(profile.address);
  const [suburb, setSuburb] = useState(profile.suburb);
  const [postcode, setPostcode] = useState(profile.postcode);
  const [saving, setSaving] = useState(false);
  const [profileMessage, setProfileMessage] = useState('');
  const [profileError, setProfileError] = useState('');

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

  const inputClass = 'w-full bg-white/5 border border-border-subtle rounded-md px-3 py-2 text-sm text-white focus:outline-none focus:ring-2 focus:ring-brand-500';

  return (
    <div className="space-y-10">
      <section>
        <h2 className="text-lg font-semibold text-white mb-4">Contact Details</h2>
        <form onSubmit={saveProfile} className="max-w-lg space-y-4">
          {profileMessage && <p className="text-sm text-green-400">{profileMessage}</p>}
          {profileError && <p className="text-sm text-red-400">{profileError}</p>}
          <div>
            <label className="block text-sm font-medium text-text-secondary mb-1">Phone</label>
            <input
              type="tel"
              value={phone}
              onChange={e => setPhone(e.target.value)}
              className={inputClass}
            />
          </div>
          <div>
            <label className="block text-sm font-medium text-text-secondary mb-1">Address</label>
            <input
              type="text"
              value={address}
              onChange={e => setAddress(e.target.value)}
              className={inputClass}
            />
          </div>
          <div className="grid grid-cols-2 gap-4">
            <div>
              <label className="block text-sm font-medium text-text-secondary mb-1">Suburb</label>
              <input
                type="text"
                value={suburb}
                onChange={e => setSuburb(e.target.value)}
                className={inputClass}
              />
            </div>
            <div>
              <label className="block text-sm font-medium text-text-secondary mb-1">Postcode</label>
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
            className="btn-brand px-4 py-2 text-sm font-medium rounded-md disabled:opacity-50"
          >
            {saving ? 'Saving...' : 'Save Profile'}
          </button>
        </form>
      </section>

      <section>
        <Link
          href="/account/vehicles"
          className="text-sm font-medium text-brand-400 hover:text-brand-400 underline"
        >
          Manage your vehicles &rarr;
        </Link>
      </section>
    </div>
  );
}
