'use client';

export function BokehBackground() {
  return (
    <div className="fixed inset-0 z-0 overflow-hidden pointer-events-none" aria-hidden="true">
      <div
        className="absolute top-1/4 -left-32 w-96 h-96 rounded-full opacity-15 blur-3xl"
        style={{
          background: 'radial-gradient(circle, #ff6b00, transparent 70%)',
          animation: 'bokeh-float 20s ease-in-out infinite',
        }}
      />
      <div
        className="absolute top-2/3 -right-32 w-80 h-80 rounded-full opacity-10 blur-3xl"
        style={{
          background: 'radial-gradient(circle, #ff8c3a, transparent 70%)',
          animation: 'bokeh-float 25s ease-in-out infinite 5s',
        }}
      />
      <div
        className="absolute -top-16 left-1/2 w-72 h-72 rounded-full opacity-8 blur-3xl"
        style={{
          background: 'radial-gradient(circle, #ff6b00, transparent 70%)',
          animation: 'bokeh-float 30s ease-in-out infinite 10s',
        }}
      />
    </div>
  );
}
