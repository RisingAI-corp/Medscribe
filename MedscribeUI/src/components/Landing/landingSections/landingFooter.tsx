function LandingFooter() {
  const currentYear = new Date().getFullYear();
  return (
    <footer className="w-full py-4 px-6 bg-gray-50 border-t border-gray-100">
      <div className="max-w-screen-xl mx-auto flex flex-col md:flex-row justify-between items-center text-sm text-gray-500">
        <div className="mb-2 md:mb-0 md:ml-auto">
          <p>Copyright © {currentYear} RisingAI</p>
        </div>
      </div>
    </footer>
  );
}

export default LandingFooter;
