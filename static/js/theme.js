// Theme toggle functionality
document.addEventListener('htmx:afterSettle', function(evt) {
    if (evt.detail.pathInfo.requestPath.startsWith('/toggle-theme')) {
        const isDark = evt.detail.pathInfo.parameters.dark === 'true';
        if (isDark) {
            document.documentElement.classList.add('dark');
        } else {
            document.documentElement.classList.remove('dark');
        }
    }
});

// Check system preference on load
if (window.matchMedia('(prefers-color-scheme: dark)').matches) {
    document.documentElement.classList.add('dark');
}