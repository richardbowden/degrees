// Server-side: call the backend directly. Client-side: go through Next.js rewrite proxy (same origin, no CORS).
const isServer = typeof window === 'undefined';
const API_BASE = isServer
  ? (process.env.BACKEND_URL || 'http://127.0.0.1:8080') + '/api/v1'
  : '/api/v1';

export interface ApiError {
  status: number;
  title: string;
  detail: string;
  errors?: Record<string, string[]>;
}

export function getCartSession(): string {
  if (isServer) return '';
  const m = document.cookie.match(/(?:^|;\s*)cart_session=([^;]*)/);
  return m ? decodeURIComponent(m[1]) : '';
}

export function setCartSession(token: string): void {
  if (isServer) return;
  document.cookie = `cart_session=${encodeURIComponent(token)};path=/;max-age=${60 * 60 * 24 * 30};samesite=lax`;
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

  // Send cart session header on cart requests
  const cartToken = opts?.cartSession ?? (path.startsWith('/cart') ? getCartSession() : '');
  if (cartToken) headers['X-Cart-Session'] = cartToken;

  const res = await (opts?.fetch ?? fetch)(`${API_BASE}${path}`, {
    method: opts?.method ?? 'GET',
    headers,
    body: opts?.body ? JSON.stringify(opts.body) : undefined,
    ...(!isServer && { credentials: 'include' as RequestCredentials }),
  });

  if (!res.ok) {
    const body = await res.json().catch(() => ({ status: res.status, title: 'Error', detail: res.statusText }));
    // Handle both RFC 7807 (detail) and gRPC-gateway (message) error formats
    const detail = body.detail || body.message || res.statusText;
    throw { status: body.status ?? res.status, title: body.title || 'Error', detail, errors: body.errors } as ApiError;
  }

  if (res.status === 204) return undefined as T;

  const data = await res.json();

  // Auto-save cart session token when backend returns one
  if (data?.cart?.sessionToken) {
    setCartSession(data.cart.sessionToken);
  }

  return data;
}
