export function Footer() {
  return (
    <footer className="border-t border-gray-200 bg-gray-50 py-8 px-4">
      <div className="max-w-6xl mx-auto text-center text-sm text-gray-500">
        <p className="font-semibold text-gray-700 mb-1">40 Degrees Car Detailing</p>
        <p>Premium mobile detailing &middot; Perth, Western Australia</p>
        <p className="mt-2">&copy; {new Date().getFullYear()} 40 Degrees Car Detailing. All rights reserved.</p>
      </div>
    </footer>
  );
}
