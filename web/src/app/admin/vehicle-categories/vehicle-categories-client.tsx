'use client';

import { useState, useEffect, useCallback } from 'react';
import { api, type ApiError } from '@/lib/api';
import type { VehicleCategory } from '@/lib/types';

function slugify(s: string): string {
  return s.toLowerCase().replace(/[^a-z0-9]+/g, '-').replace(/(^-|-$)/g, '');
}

interface CategoryForm {
  name: string;
  slug: string;
  description: string;
  sortOrder: string;
}

const EMPTY_FORM: CategoryForm = {
  name: '',
  slug: '',
  description: '',
  sortOrder: '0',
};

export function VehicleCategoriesClient({ token }: { token: string }) {
  const [categories, setCategories] = useState<VehicleCategory[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');

  const [showForm, setShowForm] = useState(false);
  const [editingId, setEditingId] = useState<string | null>(null);
  const [form, setForm] = useState<CategoryForm>(EMPTY_FORM);
  const [formLoading, setFormLoading] = useState(false);
  const [formError, setFormError] = useState('');

  const [deletingId, setDeletingId] = useState<string | null>(null);

  const fetchData = useCallback(async () => {
    setLoading(true);
    setError('');
    try {
      const res = await api<{ vehicleCategories: VehicleCategory[] }>('/catalogue/vehicle-categories', { token });
      setCategories(res.vehicleCategories ?? []);
    } catch {
      setError('Failed to load vehicle categories');
    } finally {
      setLoading(false);
    }
  }, [token]);

  useEffect(() => {
    fetchData();
  }, [fetchData]);

  function openAddForm() {
    setEditingId(null);
    setForm(EMPTY_FORM);
    setFormError('');
    setShowForm(true);
  }

  function openEditForm(cat: VehicleCategory) {
    setEditingId(cat.id);
    setForm({
      name: cat.name,
      slug: cat.slug,
      description: cat.description,
      sortOrder: String(cat.sortOrder),
    });
    setFormError('');
    setShowForm(true);
  }

  function updateForm(field: keyof CategoryForm, value: string) {
    setForm(prev => {
      const next = { ...prev, [field]: value };
      if (field === 'name' && !editingId) {
        next.slug = slugify(value);
      }
      return next;
    });
  }

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault();
    setFormLoading(true);
    setFormError('');
    const body = {
      name: form.name,
      slug: form.slug,
      description: form.description,
      sortOrder: parseInt(form.sortOrder || '0', 10),
    };

    try {
      if (editingId) {
        await api(`/admin/vehicle-categories/${editingId}`, { method: 'PUT', body: { id: editingId, ...body }, token });
      } else {
        await api('/admin/vehicle-categories', { method: 'POST', body, token });
      }
      setShowForm(false);
      await fetchData();
    } catch (err) {
      setFormError((err as ApiError)?.detail || 'Failed to save vehicle category');
    } finally {
      setFormLoading(false);
    }
  }

  async function handleDelete(id: string) {
    if (!confirm('Delete this vehicle category?')) return;
    setDeletingId(id);
    try {
      await api(`/admin/vehicle-categories/${id}`, { method: 'DELETE', token });
      await fetchData();
    } catch {
      alert('Failed to delete vehicle category');
    } finally {
      setDeletingId(null);
    }
  }

  if (loading) return <p className="text-sm text-text-muted">Loading vehicle categories...</p>;
  if (error && categories.length === 0) return <p className="text-red-400 text-sm">{error}</p>;

  return (
    <div>
      <div className="flex justify-between items-center mb-6">
        <p className="text-sm text-text-muted">{categories.length} vehicle categories</p>
        <button onClick={openAddForm} className="btn-brand px-4 py-2 rounded-md text-sm font-medium">
          Add Category
        </button>
      </div>

      {showForm && (
        <div className="glass-card p-5 mb-6">
          <h2 className="text-lg font-semibold text-white mb-4">
            {editingId ? 'Edit Vehicle Category' : 'Add Vehicle Category'}
          </h2>
          <form onSubmit={handleSubmit} className="space-y-4">
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div>
                <label className="block text-xs font-medium text-text-muted mb-1">Name</label>
                <input
                  type="text"
                  value={form.name}
                  onChange={e => updateForm('name', e.target.value)}
                  className="w-full bg-white/5 border border-border-subtle rounded-md px-3 py-1.5 text-sm text-white"
                  placeholder="e.g. Sedan / Hatchback"
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
                <label className="block text-xs font-medium text-text-muted mb-1">Description</label>
                <input
                  type="text"
                  value={form.description}
                  onChange={e => updateForm('description', e.target.value)}
                  className="w-full bg-white/5 border border-border-subtle rounded-md px-3 py-1.5 text-sm text-white"
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
            </div>
            {formError && <p className="text-sm text-red-400">{formError}</p>}
            <div className="flex gap-2">
              <button
                type="submit"
                disabled={formLoading}
                className="btn-brand px-4 py-2 rounded-md text-sm font-medium disabled:opacity-50"
              >
                {formLoading ? 'Saving...' : editingId ? 'Update' : 'Create'}
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

      <div className="space-y-3">
        {categories.map(cat => (
          <div key={cat.id} className="glass-card p-4">
            <div className="flex items-center justify-between">
              <div>
                <h3 className="font-medium text-white">{cat.name}</h3>
                <p className="text-sm text-text-muted">
                  {cat.slug}{cat.description ? ` — ${cat.description}` : ''}
                </p>
              </div>
              <div className="flex gap-2">
                <button
                  onClick={() => openEditForm(cat)}
                  className="text-xs text-brand-400 hover:text-brand-500 px-2 py-1 border border-border-subtle rounded"
                >
                  Edit
                </button>
                <button
                  onClick={() => handleDelete(cat.id)}
                  disabled={deletingId === cat.id}
                  className="text-xs text-red-400 hover:text-red-500 px-2 py-1 border border-border-subtle rounded disabled:opacity-50"
                >
                  {deletingId === cat.id ? '...' : 'Delete'}
                </button>
              </div>
            </div>
          </div>
        ))}
        {categories.length === 0 && (
          <p className="text-sm text-text-muted">No vehicle categories yet. Add one to enable size-based pricing.</p>
        )}
      </div>
    </div>
  );
}
