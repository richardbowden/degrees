'use client';

import { useState, useEffect, useCallback } from 'react';
import { api, type ApiError } from '@/lib/api';
import type { DetailingService, DetailingServiceOption, ServiceCategory, VehicleCategory, ServicePriceTier } from '@/lib/types';
import { formatPrice } from '@/lib/format';

function slugify(s: string): string {
  return s.toLowerCase().replace(/[^a-z0-9]+/g, '-').replace(/(^-|-$)/g, '');
}

interface ServiceForm {
  categoryId: string;
  name: string;
  slug: string;
  description: string;
  shortDesc: string;
  basePriceDollars: string;
  durationMinutes: string;
  isActive: boolean;
  sortOrder: string;
}

const EMPTY_FORM: ServiceForm = {
  categoryId: '',
  name: '',
  slug: '',
  description: '',
  shortDesc: '',
  basePriceDollars: '',
  durationMinutes: '',
  isActive: true,
  sortOrder: '0',
};

interface OptionForm {
  name: string;
  description: string;
  priceDollars: string;
  isActive: boolean;
  sortOrder: string;
}

const EMPTY_OPTION: OptionForm = {
  name: '',
  description: '',
  priceDollars: '',
  isActive: true,
  sortOrder: '0',
};

type StatusFilter = 'all' | 'active' | 'disabled';

export function ServicesClient({ token }: { token: string }) {
  const [categories, setCategories] = useState<ServiceCategory[]>([]);
  const [services, setServices] = useState<DetailingService[]>([]);
  const [vehicleCategories, setVehicleCategories] = useState<VehicleCategory[]>([]);
  const [statusFilter, setStatusFilter] = useState<StatusFilter>('all');
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');

  // Service form state
  const [showForm, setShowForm] = useState(false);
  const [editingId, setEditingId] = useState<string | null>(null);
  const [form, setForm] = useState<ServiceForm>(EMPTY_FORM);
  const [formLoading, setFormLoading] = useState(false);
  const [formError, setFormError] = useState('');

  // Option form state
  const [showOptionFor, setShowOptionFor] = useState<string | null>(null);
  const [optionForm, setOptionForm] = useState<OptionForm>(EMPTY_OPTION);
  const [optionLoading, setOptionLoading] = useState(false);
  const [optionError, setOptionError] = useState('');

  // Price tier state
  const [showTiersFor, setShowTiersFor] = useState<string | null>(null);
  const [tierPrices, setTierPrices] = useState<Record<string, string>>({});
  const [tierLoading, setTierLoading] = useState(false);
  const [tierError, setTierError] = useState('');

  // Delete state
  const [deletingId, setDeletingId] = useState<string | null>(null);

  // Toggle active state
  const [togglingId, setTogglingId] = useState<string | null>(null);

  async function handleToggleActive(svc: DetailingService) {
    setTogglingId(svc.id);
    try {
      await api(`/admin/services/${svc.id}`, {
        method: 'PUT',
        body: {
          id: svc.id,
          categoryId: svc.categoryId,
          name: svc.name,
          slug: svc.slug,
          description: svc.description,
          shortDesc: svc.shortDesc,
          basePrice: svc.basePrice,
          durationMinutes: svc.durationMinutes,
          isActive: !svc.isActive,
          sortOrder: svc.sortOrder,
        },
        token,
      });
      await fetchData();
    } catch {
      alert('Failed to update service');
    } finally {
      setTogglingId(null);
    }
  }

  const fetchData = useCallback(async () => {
    setLoading(true);
    setError('');
    try {
      const [catRes, svcRes] = await Promise.all([
        api<{ categories: ServiceCategory[] }>('/catalogue/categories', { token }),
        api<{ services: DetailingService[] }>('/admin/services', { token }),
      ]);
      setCategories(catRes.categories ?? []);
      setServices(svcRes.services ?? []);

      // Vehicle categories are optional - don't break the page if they fail
      try {
        const vcRes = await api<{ vehicleCategories: VehicleCategory[] }>('/catalogue/vehicle-categories', { token });
        setVehicleCategories(vcRes.vehicleCategories ?? []);
      } catch {
        setVehicleCategories([]);
      }
    } catch {
      setError('Failed to load services');
    } finally {
      setLoading(false);
    }
  }, [token]);

  useEffect(() => {
    fetchData();
  }, [fetchData]);

  function openAddForm() {
    setEditingId(null);
    setForm({ ...EMPTY_FORM, categoryId: categories[0]?.id ?? '' });
    setFormError('');
    setShowForm(true);
  }

  function openEditForm(svc: DetailingService) {
    setEditingId(svc.id);
    setForm({
      categoryId: svc.categoryId,
      name: svc.name,
      slug: svc.slug,
      description: svc.description,
      shortDesc: svc.shortDesc,
      basePriceDollars: (Number(svc.basePrice) / 100).toFixed(2),
      durationMinutes: String(svc.durationMinutes),
      isActive: svc.isActive,
      sortOrder: String(svc.sortOrder),
    });
    setFormError('');
    setShowForm(true);
  }

  function updateForm(field: keyof ServiceForm, value: string | boolean) {
    setForm(prev => {
      const next = { ...prev, [field]: value };
      if (field === 'name' && !editingId) {
        next.slug = slugify(value as string);
      }
      return next;
    });
  }

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault();
    setFormLoading(true);
    setFormError('');
    const body = {
      categoryId: form.categoryId,
      name: form.name,
      slug: form.slug,
      description: form.description,
      shortDesc: form.shortDesc,
      basePrice: Math.round(parseFloat(form.basePriceDollars || '0') * 100),
      durationMinutes: parseInt(form.durationMinutes || '0', 10),
      isActive: form.isActive,
      sortOrder: parseInt(form.sortOrder || '0', 10),
    };

    try {
      if (editingId) {
        await api(`/admin/services/${editingId}`, { method: 'PUT', body: { id: editingId, ...body }, token });
      } else {
        await api('/admin/services', { method: 'POST', body, token });
      }
      setShowForm(false);
      await fetchData();
    } catch (err) {
      setFormError((err as ApiError)?.detail || 'Failed to save service');
    } finally {
      setFormLoading(false);
    }
  }

  async function handleDelete(id: string) {
    if (!confirm('Delete this service?')) return;
    setDeletingId(id);
    try {
      await api(`/admin/services/${id}`, { method: 'DELETE', token });
      await fetchData();
    } catch {
      alert('Failed to delete service');
    } finally {
      setDeletingId(null);
    }
  }

  async function handleAddOption(serviceId: string, e: React.FormEvent) {
    e.preventDefault();
    setOptionLoading(true);
    setOptionError('');
    try {
      await api(`/admin/services/${serviceId}/options`, {
        method: 'POST',
        body: {
          serviceId: serviceId,
          name: optionForm.name,
          description: optionForm.description,
          price: Math.round(parseFloat(optionForm.priceDollars || '0') * 100),
          isActive: optionForm.isActive,
          sortOrder: parseInt(optionForm.sortOrder || '0', 10),
        },
        token,
      });
      setShowOptionFor(null);
      setOptionForm(EMPTY_OPTION);
      await fetchData();
    } catch (err) {
      setOptionError((err as ApiError)?.detail || 'Failed to add option');
    } finally {
      setOptionLoading(false);
    }
  }

  function openTiers(svc: DetailingService) {
    const prices: Record<string, string> = {};
    for (const vc of vehicleCategories) {
      const tier = (svc.priceTiers ?? []).find((t: ServicePriceTier) => t.vehicleCategoryId === vc.id);
      prices[vc.id] = tier ? (Number(tier.price) / 100).toFixed(2) : '';
    }
    setTierPrices(prices);
    setTierError('');
    setShowTiersFor(svc.id);
  }

  async function handleSaveTiers(serviceId: string) {
    setTierLoading(true);
    setTierError('');
    try {
      const tiers = vehicleCategories
        .filter(vc => tierPrices[vc.id] && parseFloat(tierPrices[vc.id]) > 0)
        .map(vc => ({
          vehicleCategoryId: vc.id,
          price: Math.round(parseFloat(tierPrices[vc.id]) * 100),
        }));
      await api(`/admin/services/${serviceId}/price-tiers`, {
        method: 'PUT',
        body: { serviceId, tiers },
        token,
      });
      setShowTiersFor(null);
      await fetchData();
    } catch (err) {
      setTierError((err as ApiError)?.detail || 'Failed to save price tiers');
    } finally {
      setTierLoading(false);
    }
  }

  // Filter services by status
  const filteredServices = services.filter(s => {
    if (statusFilter === 'active') return s.isActive;
    if (statusFilter === 'disabled') return !s.isActive;
    return true;
  });

  // Group services by category
  const grouped = categories.map(cat => ({
    category: cat,
    services: filteredServices.filter(s => s.categoryId === cat.id),
  }));
  const uncategorized = filteredServices.filter(s => !categories.some(c => c.id === s.categoryId));

  if (loading) return <p className="text-sm text-text-muted">Loading services...</p>;
  if (error && services.length === 0) return <p className="text-red-400 text-sm">{error}</p>;

  return (
    <div>
      <div className="flex justify-between items-center mb-6">
        <div className="flex items-center gap-4">
          <p className="text-sm text-text-muted">
            {filteredServices.length} of {services.length} services
          </p>
          <div className="flex gap-1 bg-white/5 rounded-md p-0.5">
            {(['all', 'active', 'disabled'] as StatusFilter[]).map(f => (
              <button
                key={f}
                onClick={() => setStatusFilter(f)}
                className={`px-3 py-1 rounded text-xs font-medium capitalize transition-colors ${
                  statusFilter === f
                    ? 'bg-white/15 text-white'
                    : 'text-text-muted hover:text-white'
                }`}
              >
                {f}
              </button>
            ))}
          </div>
        </div>
        <button
          onClick={openAddForm}
          className="btn-brand px-4 py-2 rounded-md text-sm font-medium"
        >
          Add Service
        </button>
      </div>

      {/* Service Form */}
      {showForm && (
        <div className="glass-card p-5 mb-6">
          <h2 className="text-lg font-semibold text-white mb-4">
            {editingId ? 'Edit Service' : 'Add Service'}
          </h2>
          <form onSubmit={handleSubmit} className="space-y-4">
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div>
                <label className="block text-xs font-medium text-text-muted mb-1">Category</label>
                <select
                  value={form.categoryId}
                  onChange={e => updateForm('categoryId', e.target.value)}
                  className="w-full bg-white/5 border border-border-subtle rounded-md px-3 py-1.5 text-sm text-white"
                  required
                >
                  <option value="">Select category</option>
                  {categories.map(c => (
                    <option key={c.id} value={c.id}>{c.name}</option>
                  ))}
                </select>
              </div>
              <div>
                <label className="block text-xs font-medium text-text-muted mb-1">Name</label>
                <input
                  type="text"
                  value={form.name}
                  onChange={e => updateForm('name', e.target.value)}
                  className="w-full bg-white/5 border border-border-subtle rounded-md px-3 py-1.5 text-sm text-white"
                  required
                />
              </div>
              <div>
                <label className="block text-xs font-medium text-text-muted mb-1">Slug</label>
                <input
                  type="text"
                  value={form.slug}
                  onChange={e => updateForm('slug', e.target.value)}
                  className="w-full bg-white/5 border border-border-subtle rounded-md px-3 py-1.5 text-sm text-white"
                  required
                />
              </div>
              <div>
                <label className="block text-xs font-medium text-text-muted mb-1">Short Description</label>
                <input
                  type="text"
                  value={form.shortDesc}
                  onChange={e => updateForm('shortDesc', e.target.value)}
                  className="w-full bg-white/5 border border-border-subtle rounded-md px-3 py-1.5 text-sm text-white"
                />
              </div>
              <div>
                <label className="block text-xs font-medium text-text-muted mb-1">Base Price (AUD)</label>
                <input
                  type="number"
                  step="0.01"
                  min="0"
                  value={form.basePriceDollars}
                  onChange={e => updateForm('basePriceDollars', e.target.value)}
                  className="w-full bg-white/5 border border-border-subtle rounded-md px-3 py-1.5 text-sm text-white"
                  required
                />
              </div>
              <div>
                <label className="block text-xs font-medium text-text-muted mb-1">Duration (minutes)</label>
                <input
                  type="number"
                  min="0"
                  value={form.durationMinutes}
                  onChange={e => updateForm('durationMinutes', e.target.value)}
                  className="w-full bg-white/5 border border-border-subtle rounded-md px-3 py-1.5 text-sm text-white"
                  required
                />
              </div>
              <div>
                <label className="block text-xs font-medium text-text-muted mb-1">Sort Order</label>
                <input
                  type="number"
                  value={form.sortOrder}
                  onChange={e => updateForm('sortOrder', e.target.value)}
                  className="w-full bg-white/5 border border-border-subtle rounded-md px-3 py-1.5 text-sm text-white"
                />
              </div>
              <div className="flex items-end">
                <label className="flex items-center gap-2 text-sm text-text-secondary">
                  <input
                    type="checkbox"
                    checked={form.isActive}
                    onChange={e => updateForm('isActive', e.target.checked)}
                    className="rounded border-border-subtle"
                  />
                  Active
                </label>
              </div>
            </div>
            <div>
              <label className="block text-xs font-medium text-text-muted mb-1">Description</label>
              <p className="text-xs text-text-muted mb-1">Supports Markdown: **bold**, *italic*, - bullet lists</p>
              <textarea
                value={form.description}
                onChange={e => updateForm('description', e.target.value)}
                rows={3}
                className="w-full bg-white/5 border border-border-subtle rounded-md px-3 py-2 text-sm text-white"
              />
            </div>
            {formError && <p className="text-sm text-red-400">{formError}</p>}
            <div className="flex gap-2">
              <button
                type="submit"
                disabled={formLoading}
                className="btn-brand px-4 py-2 rounded-md text-sm font-medium disabled:opacity-50"
              >
                {formLoading ? 'Saving...' : editingId ? 'Update Service' : 'Create Service'}
              </button>
              <button
                type="button"
                onClick={() => setShowForm(false)}
                className="px-4 py-2 border border-border-subtle rounded-md text-sm text-text-secondary hover:bg-white/5"
              >
                Cancel
              </button>
            </div>
          </form>
        </div>
      )}

      {/* Services grouped by category */}
      {grouped.map(({ category, services: catServices }) => (
        <div key={category.id} className="mb-8">
          <h2 className="text-lg font-semibold text-white mb-3">{category.name}</h2>
          {catServices.length === 0 ? (
            <p className="text-sm text-text-muted">No services in this category.</p>
          ) : (
            <div className="space-y-3">
              {catServices.map(svc => (
                <div key={svc.id} className="glass-card p-4">
                  <div className="flex items-start justify-between">
                    <div className="flex-1">
                      <div className="flex items-center gap-2">
                        <h3 className="font-medium text-white">{svc.name}</h3>
                        {!svc.isActive && (
                          <span className="text-xs bg-white/10 text-text-muted px-1.5 py-0.5 rounded">Inactive</span>
                        )}
                      </div>
                      {svc.shortDesc && <p className="text-sm text-text-muted mt-0.5">{svc.shortDesc}</p>}
                      <div className="flex gap-4 mt-2 text-sm text-text-secondary">
                        <span>{formatPrice(svc.basePrice)}</span>
                        <span>{svc.durationMinutes} mins</span>
                      </div>

                      {/* Options */}
                      {svc.options && svc.options.length > 0 && (
                        <div className="mt-3 pl-4 border-l-2 border-white/10">
                          <p className="text-xs font-medium text-text-muted mb-1">Options</p>
                          {svc.options.map((opt: DetailingServiceOption) => (
                            <div key={opt.id} className="flex items-center gap-3 text-sm text-text-secondary py-0.5">
                              <span>{opt.name}</span>
                              <span className="text-text-muted">+{formatPrice(opt.price)}</span>
                              {!opt.isActive && <span className="text-xs text-text-muted">(inactive)</span>}
                            </div>
                          ))}
                        </div>
                      )}

                      {/* Add Option Form */}
                      {showOptionFor === svc.id && (
                        <form onSubmit={(e) => handleAddOption(svc.id, e)} className="mt-3 p-3 bg-white/5 rounded-md space-y-3">
                          <div className="grid grid-cols-2 gap-3">
                            <input
                              type="text"
                              value={optionForm.name}
                              onChange={e => setOptionForm(prev => ({ ...prev, name: e.target.value }))}
                              placeholder="Option name"
                              className="bg-white/5 border border-border-subtle rounded-md px-3 py-1.5 text-sm text-white"
                              required
                            />
                            <input
                              type="number"
                              step="0.01"
                              min="0"
                              value={optionForm.priceDollars}
                              onChange={e => setOptionForm(prev => ({ ...prev, priceDollars: e.target.value }))}
                              placeholder="Price (AUD)"
                              className="bg-white/5 border border-border-subtle rounded-md px-3 py-1.5 text-sm text-white"
                              required
                            />
                          </div>
                          <input
                            type="text"
                            value={optionForm.description}
                            onChange={e => setOptionForm(prev => ({ ...prev, description: e.target.value }))}
                            placeholder="Description"
                            className="w-full bg-white/5 border border-border-subtle rounded-md px-3 py-1.5 text-sm text-white"
                          />
                          <div className="flex gap-3 items-center">
                            <input
                              type="number"
                              value={optionForm.sortOrder}
                              onChange={e => setOptionForm(prev => ({ ...prev, sortOrder: e.target.value }))}
                              placeholder="Sort"
                              className="w-20 bg-white/5 border border-border-subtle rounded-md px-3 py-1.5 text-sm text-white"
                            />
                            <label className="flex items-center gap-1.5 text-sm text-text-secondary">
                              <input
                                type="checkbox"
                                checked={optionForm.isActive}
                                onChange={e => setOptionForm(prev => ({ ...prev, isActive: e.target.checked }))}
                                className="rounded border-border-subtle"
                              />
                              Active
                            </label>
                          </div>
                          {optionError && <p className="text-sm text-red-400">{optionError}</p>}
                          <div className="flex gap-2">
                            <button
                              type="submit"
                              disabled={optionLoading}
                              className="btn-brand px-3 py-1.5 rounded-md text-sm disabled:opacity-50"
                            >
                              {optionLoading ? 'Adding...' : 'Add Option'}
                            </button>
                            <button
                              type="button"
                              onClick={() => { setShowOptionFor(null); setOptionForm(EMPTY_OPTION); setOptionError(''); }}
                              className="px-3 py-1.5 border border-border-subtle rounded-md text-sm text-text-secondary hover:bg-white/5"
                            >
                              Cancel
                            </button>
                          </div>
                        </form>
                      )}

                      {/* Price Tiers */}
                      {svc.priceTiers && svc.priceTiers.length > 0 && showTiersFor !== svc.id && (
                        <div className="mt-3 pl-4 border-l-2 border-brand-500/30 space-y-0.5">
                          <p className="text-xs font-medium text-text-muted mb-1">Price Tiers</p>
                          {svc.priceTiers.map((t: ServicePriceTier) => (
                            <div key={t.vehicleCategoryId} className="flex items-center gap-3 text-sm text-text-secondary py-0.5">
                              <span>{t.categoryName}</span>
                              <span className="text-brand-400 font-semibold">{formatPrice(t.price)}</span>
                            </div>
                          ))}
                        </div>
                      )}

                      {showTiersFor === svc.id && (
                        <div className="mt-3 p-3 bg-white/5 rounded-md space-y-3">
                          <p className="text-xs font-medium text-text-muted">Price per Vehicle Category (AUD)</p>
                          {vehicleCategories.length === 0 ? (
                            <p className="text-sm text-text-muted">No vehicle categories. Create some first.</p>
                          ) : (
                            <div className="space-y-2">
                              {vehicleCategories.map(vc => (
                                <div key={vc.id} className="flex items-center gap-3">
                                  <span className="text-sm text-text-secondary w-40">{vc.name}</span>
                                  <input
                                    type="number"
                                    step="0.01"
                                    min="0"
                                    value={tierPrices[vc.id] ?? ''}
                                    onChange={e => setTierPrices(prev => ({ ...prev, [vc.id]: e.target.value }))}
                                    placeholder="Leave blank for base price"
                                    className="flex-1 bg-white/5 border border-border-subtle rounded-md px-3 py-1.5 text-sm text-white"
                                  />
                                </div>
                              ))}
                            </div>
                          )}
                          {tierError && <p className="text-sm text-red-400">{tierError}</p>}
                          <div className="flex gap-2">
                            <button
                              onClick={() => handleSaveTiers(svc.id)}
                              disabled={tierLoading}
                              className="btn-brand px-3 py-1.5 rounded-md text-sm disabled:opacity-50"
                            >
                              {tierLoading ? 'Saving...' : 'Save Tiers'}
                            </button>
                            <button
                              onClick={() => setShowTiersFor(null)}
                              className="px-3 py-1.5 border border-border-subtle rounded-md text-sm text-text-secondary hover:bg-white/5"
                            >
                              Cancel
                            </button>
                          </div>
                        </div>
                      )}
                    </div>

                    {/* Actions */}
                    <div className="flex gap-2 ml-4 flex-shrink-0">
                      <button
                        onClick={() => handleToggleActive(svc)}
                        disabled={togglingId === svc.id}
                        className={`text-xs px-2 py-1 rounded border disabled:opacity-50 ${
                          svc.isActive
                            ? 'text-green-400 border-green-400/30 hover:bg-green-400/10'
                            : 'text-text-muted border-border-subtle hover:bg-white/5'
                        }`}
                      >
                        {togglingId === svc.id ? '...' : svc.isActive ? 'Enabled' : 'Disabled'}
                      </button>
                      <button
                        onClick={() => openTiers(svc)}
                        className="text-xs text-text-muted hover:text-white px-2 py-1 border border-border-subtle rounded"
                      >
                        Tiers
                      </button>
                      <button
                        onClick={() => { setShowOptionFor(svc.id); setOptionForm(EMPTY_OPTION); setOptionError(''); }}
                        className="text-xs text-text-muted hover:text-white px-2 py-1 border border-border-subtle rounded"
                      >
                        + Option
                      </button>
                      <button
                        onClick={() => openEditForm(svc)}
                        className="text-xs text-brand-400 hover:text-brand-500 px-2 py-1 border border-border-subtle rounded"
                      >
                        Edit
                      </button>
                      <button
                        onClick={() => handleDelete(svc.id)}
                        disabled={deletingId === svc.id}
                        className="text-xs text-red-400 hover:text-red-500 px-2 py-1 border border-border-subtle rounded disabled:opacity-50"
                      >
                        {deletingId === svc.id ? '...' : 'Delete'}
                      </button>
                    </div>
                  </div>
                </div>
              ))}
            </div>
          )}
        </div>
      ))}

      {uncategorized.length > 0 && (
        <div className="mb-8">
          <h2 className="text-lg font-semibold text-white mb-3">Uncategorized</h2>
          <div className="space-y-3">
            {uncategorized.map(svc => (
              <div key={svc.id} className="glass-card p-4">
                <div className="flex items-center justify-between">
                  <div>
                    <h3 className="font-medium text-white">{svc.name}</h3>
                    <p className="text-sm text-text-muted">{formatPrice(svc.basePrice)} - {svc.durationMinutes} mins</p>
                  </div>
                  <div className="flex gap-2">
                    <button
                      onClick={() => handleToggleActive(svc)}
                      disabled={togglingId === svc.id}
                      className={`text-xs px-2 py-1 rounded border disabled:opacity-50 ${
                        svc.isActive
                          ? 'text-green-400 border-green-400/30 hover:bg-green-400/10'
                          : 'text-text-muted border-border-subtle hover:bg-white/5'
                      }`}
                    >
                      {togglingId === svc.id ? '...' : svc.isActive ? 'Enabled' : 'Disabled'}
                    </button>
                    <button onClick={() => openEditForm(svc)} className="text-xs text-brand-400 hover:text-brand-500 px-2 py-1 border border-border-subtle rounded">Edit</button>
                    <button onClick={() => handleDelete(svc.id)} disabled={deletingId === svc.id} className="text-xs text-red-400 hover:text-red-500 px-2 py-1 border border-border-subtle rounded disabled:opacity-50">
                      {deletingId === svc.id ? '...' : 'Delete'}
                    </button>
                  </div>
                </div>
              </div>
            ))}
          </div>
        </div>
      )}
    </div>
  );
}
