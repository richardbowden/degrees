export function Footer() {
  return (
    <footer className="relative z-10 border-t border-border-subtle bg-surface-raised py-8 px-4">
      <div className="max-w-6xl mx-auto text-center text-sm text-text-muted">
        <p className="font-semibold text-text-secondary mb-1">40 Degrees Car Detailing</p>
        <p>Premium mobile detailing &middot; Perth, Western Australia</p>
        <p className="mt-2">&copy; {new Date().getFullYear()} 40 Degrees Car Detailing. All rights reserved.</p>
      </div>
    </footer>
  );
}
