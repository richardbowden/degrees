# 40 Degrees Frontend

## Stack
Next.js 15+ App Router, TypeScript, Tailwind CSS.

## API
Backend REST API at http://localhost:8080/api/v1 (Go gRPC gateway).
Session auth via httpOnly cookies, always use credentials: include.
All prices are int64 cents from the API, format as AUD client-side.
Errors are RFC 7807 Problem Details.

## Patterns
- Server components for public pages (SSR for SEO)
- Client components for interactive features (cart, forms, calendar)
- Single API client in lib/api.ts, all requests go through it
- Types in lib/types.ts matching API response shapes
- Format helpers in lib/format.ts (cents to dollars, date formatting)
- Next.js middleware for route protection (check session cookie)

## Style
- Tailwind utility classes only, no custom CSS files
- Functional first, polish later. Keep it clean and readable.
- No component library. Simple, self-contained components.
- No dark mode for now.
- Mobile responsive but desktop-first layout.

## Rules
- No localStorage or sessionStorage
- Use React state (useState, useReducer) for client state
- Forms use server actions where possible, otherwise POST via api client
- Loading states on all async operations
- Error boundaries on route segments
