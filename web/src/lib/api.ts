const API_BASE = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api/v1';

export interface ApiError {
  status: number;
  title: string;
  detail: string;
  errors?: Record<string, string[]>;
}

export async function api<T = unknown>(path: string, opts?: {
  method?: string;
  body?: unknown;
  token?: string;
  cartSession?: string;
  fetch?: typeof globalThis.fetch;
}): Promise<T> {
  const headers: Record<string, string> = { 'Content-Type': 'application/json' };
  if (opts?.token) headers['Authorization'] = `Bearer ${opts.token}`;
  if (opts?.cartSession) headers['X-Cart-Session'] = opts.cartSession;

  const res = await (opts?.fetch ?? fetch)(`${API_BASE}${path}`, {
    method: opts?.method ?? 'GET',
    headers,
    body: opts?.body ? JSON.stringify(opts.body) : undefined,
    credentials: 'include',
  });

  if (!res.ok) {
    const body = await res.json().catch(() => ({ status: res.status, title: 'Error', detail: res.statusText }));
    // Handle both RFC 7807 (detail) and gRPC-gateway (message) error formats
    const detail = body.detail || body.message || res.statusText;
    throw { status: body.status ?? res.status, title: body.title || 'Error', detail, errors: body.errors } as ApiError;
  }

  if (res.status === 204) return undefined as T;
  return res.json();
}
