'use client';

import { useState, useEffect, useCallback } from 'react';
import { api, type ApiError } from '@/lib/api';
import type { DetailingService, DetailingServiceOption, ServiceCategory } from '@/lib/types';
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

export function ServicesClient({ token }: { token: string }) {
  const [categories, setCategories] = useState<ServiceCategory[]>([]);
  const [services, setServices] = useState<DetailingService[]>([]);
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

  // Delete state
  const [deletingId, setDeletingId] = useState<string | null>(null);

  const fetchData = useCallback(async () => {
    setLoading(true);
    setError('');
    try {
      const [catRes, svcRes] = await Promise.all([
        api<{ categories: ServiceCategory[] }>('/catalogue/categories', { token }),
        api<{ services: DetailingService[] }>('/catalogue', { token }),
      ]);
      setCategories(catRes.categories ?? []);
      setServices(svcRes.services ?? []);
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
      basePriceDollars: (svc.basePrice / 100).toFixed(2),
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

  // Group services by category
  const grouped = categories.map(cat => ({
    category: cat,
    services: services.filter(s => s.categoryId === cat.id),
  }));
  const uncategorized = services.filter(s => !categories.some(c => c.id === s.categoryId));

  if (loading) return <p className="text-sm text-gray-500">Loading services...</p>;
  if (error && services.length === 0) return <p className="text-red-600 text-sm">{error}</p>;

  return (
    <div>
      <div className="flex justify-between items-center mb-6">
        <p className="text-sm text-gray-500">{services.length} services across {categories.length} categories</p>
        <button
          onClick={openAddForm}
          className="bg-gray-900 text-white px-4 py-2 rounded-md text-sm font-medium hover:bg-gray-800"
        >
          Add Service
        </button>
      </div>

      {/* Service Form */}
      {showForm && (
        <div className="bg-white rounded-lg shadow-sm border border-gray-200 p-5 mb-6">
          <h2 className="text-lg font-semibold text-gray-900 mb-4">
            {editingId ? 'Edit Service' : 'Add Service'}
          </h2>
          <form onSubmit={handleSubmit} className="space-y-4">
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div>
                <label className="block text-xs font-medium text-gray-500 mb-1">Category</label>
                <select
                  value={form.categoryId}
                  onChange={e => updateForm('categoryId', e.target.value)}
                  className="w-full border border-gray-300 rounded-md px-3 py-1.5 text-sm"
                  required
                >
                  <option value="">Select category</option>
                  {categories.map(c => (
                    <option key={c.id} value={c.id}>{c.name}</option>
                  ))}
                </select>
              </div>
              <div>
                <label className="block text-xs font-medium text-gray-500 mb-1">Name</label>
                <input
                  type="text"
                  value={form.name}
                  onChange={e => updateForm('name', e.target.value)}
                  className="w-full border border-gray-300 rounded-md px-3 py-1.5 text-sm"
                  required
                />
              </div>
              <div>
                <label className="block text-xs font-medium text-gray-500 mb-1">Slug</label>
                <input
                  type="text"
                  value={form.slug}
                  onChange={e => updateForm('slug', e.target.value)}
                  className="w-full border border-gray-300 rounded-md px-3 py-1.5 text-sm"
                  required
                />
              </div>
              <div>
                <label className="block text-xs font-medium text-gray-500 mb-1">Short Description</label>
                <input
                  type="text"
                  value={form.shortDesc}
                  onChange={e => updateForm('shortDesc', e.target.value)}
                  className="w-full border border-gray-300 rounded-md px-3 py-1.5 text-sm"
                />
              </div>
              <div>
                <label className="block text-xs font-medium text-gray-500 mb-1">Base Price (AUD)</label>
                <input
                  type="number"
                  step="0.01"
                  min="0"
                  value={form.basePriceDollars}
                  onChange={e => updateForm('basePriceDollars', e.target.value)}
                  className="w-full border border-gray-300 rounded-md px-3 py-1.5 text-sm"
                  required
                />
              </div>
              <div>
                <label className="block text-xs font-medium text-gray-500 mb-1">Duration (minutes)</label>
                <input
                  type="number"
                  min="0"
                  value={form.durationMinutes}
                  onChange={e => updateForm('durationMinutes', e.target.value)}
                  className="w-full border border-gray-300 rounded-md px-3 py-1.5 text-sm"
                  required
                />
              </div>
              <div>
                <label className="block text-xs font-medium text-gray-500 mb-1">Sort Order</label>
                <input
                  type="number"
                  value={form.sortOrder}
                  onChange={e => updateForm('sortOrder', e.target.value)}
                  className="w-full border border-gray-300 rounded-md px-3 py-1.5 text-sm"
                />
              </div>
              <div className="flex items-end">
                <label className="flex items-center gap-2 text-sm text-gray-700">
                  <input
                    type="checkbox"
                    checked={form.isActive}
                    onChange={e => updateForm('isActive', e.target.checked)}
                    className="rounded border-gray-300"
                  />
                  Active
                </label>
              </div>
            </div>
            <div>
              <label className="block text-xs font-medium text-gray-500 mb-1">Description</label>
              <textarea
                value={form.description}
                onChange={e => updateForm('description', e.target.value)}
                rows={3}
                className="w-full border border-gray-300 rounded-md px-3 py-2 text-sm"
              />
            </div>
            {formError && <p className="text-sm text-red-600">{formError}</p>}
            <div className="flex gap-2">
              <button
                type="submit"
                disabled={formLoading}
                className="bg-gray-900 text-white px-4 py-2 rounded-md text-sm font-medium hover:bg-gray-800 disabled:opacity-50"
              >
                {formLoading ? 'Saving...' : editingId ? 'Update Service' : 'Create Service'}
              </button>
              <button
                type="button"
                onClick={() => setShowForm(false)}
                className="px-4 py-2 border border-gray-300 rounded-md text-sm text-gray-700 hover:bg-gray-50"
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
          <h2 className="text-lg font-semibold text-gray-900 mb-3">{category.name}</h2>
          {catServices.length === 0 ? (
            <p className="text-sm text-gray-500">No services in this category.</p>
          ) : (
            <div className="space-y-3">
              {catServices.map(svc => (
                <div key={svc.id} className="bg-white rounded-lg shadow-sm border border-gray-200 p-4">
                  <div className="flex items-start justify-between">
                    <div className="flex-1">
                      <div className="flex items-center gap-2">
                        <h3 className="font-medium text-gray-900">{svc.name}</h3>
                        {!svc.isActive && (
                          <span className="text-xs bg-gray-100 text-gray-500 px-1.5 py-0.5 rounded">Inactive</span>
                        )}
                      </div>
                      {svc.shortDesc && <p className="text-sm text-gray-500 mt-0.5">{svc.shortDesc}</p>}
                      <div className="flex gap-4 mt-2 text-sm text-gray-600">
                        <span>{formatPrice(svc.basePrice)}</span>
                        <span>{svc.durationMinutes} mins</span>
                      </div>

                      {/* Options */}
                      {svc.options && svc.options.length > 0 && (
                        <div className="mt-3 pl-4 border-l-2 border-gray-100">
                          <p className="text-xs font-medium text-gray-500 mb-1">Options</p>
                          {svc.options.map((opt: DetailingServiceOption) => (
                            <div key={opt.id} className="flex items-center gap-3 text-sm text-gray-600 py-0.5">
                              <span>{opt.name}</span>
                              <span className="text-gray-400">+{formatPrice(opt.price)}</span>
                              {!opt.isActive && <span className="text-xs text-gray-400">(inactive)</span>}
                            </div>
                          ))}
                        </div>
                      )}

                      {/* Add Option Form */}
                      {showOptionFor === svc.id && (
                        <form onSubmit={(e) => handleAddOption(svc.id, e)} className="mt-3 p-3 bg-gray-50 rounded-md space-y-3">
                          <div className="grid grid-cols-2 gap-3">
                            <input
                              type="text"
                              value={optionForm.name}
                              onChange={e => setOptionForm(prev => ({ ...prev, name: e.target.value }))}
                              placeholder="Option name"
                              className="border border-gray-300 rounded-md px-3 py-1.5 text-sm"
                              required
                            />
                            <input
                              type="number"
                              step="0.01"
                              min="0"
                              value={optionForm.priceDollars}
                              onChange={e => setOptionForm(prev => ({ ...prev, priceDollars: e.target.value }))}
                              placeholder="Price (AUD)"
                              className="border border-gray-300 rounded-md px-3 py-1.5 text-sm"
                              required
                            />
                          </div>
                          <input
                            type="text"
                            value={optionForm.description}
                            onChange={e => setOptionForm(prev => ({ ...prev, description: e.target.value }))}
                            placeholder="Description"
                            className="w-full border border-gray-300 rounded-md px-3 py-1.5 text-sm"
                          />
                          <div className="flex gap-3 items-center">
                            <input
                              type="number"
                              value={optionForm.sortOrder}
                              onChange={e => setOptionForm(prev => ({ ...prev, sortOrder: e.target.value }))}
                              placeholder="Sort"
                              className="w-20 border border-gray-300 rounded-md px-3 py-1.5 text-sm"
                            />
                            <label className="flex items-center gap-1.5 text-sm text-gray-600">
                              <input
                                type="checkbox"
                                checked={optionForm.isActive}
                                onChange={e => setOptionForm(prev => ({ ...prev, isActive: e.target.checked }))}
                                className="rounded border-gray-300"
                              />
                              Active
                            </label>
                          </div>
                          {optionError && <p className="text-sm text-red-600">{optionError}</p>}
                          <div className="flex gap-2">
                            <button
                              type="submit"
                              disabled={optionLoading}
                              className="bg-gray-900 text-white px-3 py-1.5 rounded-md text-sm hover:bg-gray-800 disabled:opacity-50"
                            >
                              {optionLoading ? 'Adding...' : 'Add Option'}
                            </button>
                            <button
                              type="button"
                              onClick={() => { setShowOptionFor(null); setOptionForm(EMPTY_OPTION); setOptionError(''); }}
                              className="px-3 py-1.5 border border-gray-300 rounded-md text-sm text-gray-700 hover:bg-gray-50"
                            >
                              Cancel
                            </button>
                          </div>
                        </form>
                      )}
                    </div>

                    {/* Actions */}
                    <div className="flex gap-2 ml-4 flex-shrink-0">
                      <button
                        onClick={() => { setShowOptionFor(svc.id); setOptionForm(EMPTY_OPTION); setOptionError(''); }}
                        className="text-xs text-gray-500 hover:text-gray-700 px-2 py-1 border border-gray-200 rounded"
                      >
                        + Option
                      </button>
                      <button
                        onClick={() => openEditForm(svc)}
                        className="text-xs text-blue-600 hover:text-blue-800 px-2 py-1 border border-gray-200 rounded"
                      >
                        Edit
                      </button>
                      <button
                        onClick={() => handleDelete(svc.id)}
                        disabled={deletingId === svc.id}
                        className="text-xs text-red-600 hover:text-red-800 px-2 py-1 border border-gray-200 rounded disabled:opacity-50"
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
          <h2 className="text-lg font-semibold text-gray-900 mb-3">Uncategorized</h2>
          <div className="space-y-3">
            {uncategorized.map(svc => (
              <div key={svc.id} className="bg-white rounded-lg shadow-sm border border-gray-200 p-4">
                <div className="flex items-center justify-between">
                  <div>
                    <h3 className="font-medium text-gray-900">{svc.name}</h3>
                    <p className="text-sm text-gray-500">{formatPrice(svc.basePrice)} - {svc.durationMinutes} mins</p>
                  </div>
                  <div className="flex gap-2">
                    <button onClick={() => openEditForm(svc)} className="text-xs text-blue-600 hover:text-blue-800 px-2 py-1 border border-gray-200 rounded">Edit</button>
                    <button onClick={() => handleDelete(svc.id)} disabled={deletingId === svc.id} className="text-xs text-red-600 hover:text-red-800 px-2 py-1 border border-gray-200 rounded disabled:opacity-50">
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
