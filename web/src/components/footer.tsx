export function Footer() {
  return (
    <footer className="relative z-10 border-t border-border-subtle bg-surface-raised py-8 px-4">
      <div className="max-w-6xl mx-auto text-center text-sm text-text-muted">
        <p className="font-semibold text-text-secondary mb-1">40 Degrees Car Detailing</p>
        <p>Premium mobile detailing &middot; Yanchep, down to Perth and beyond</p>
        <div className="flex flex-col sm:flex-row gap-4 justify-center items-center mt-3 text-text-muted">
          <a href="tel:0448263659" className="hover:text-white transition-colors">0448 263 659</a>
          <span className="hidden sm:inline">&middot;</span>
          <a href="mailto:detailing@40degrees.au" className="hover:text-white transition-colors">detailing@40degrees.au</a>
        </div>
        <p className="mt-3">&copy; {new Date().getFullYear()} 40 Degrees Car Detailing. All rights reserved.</p>
      </div>
    </footer>
  );
}
