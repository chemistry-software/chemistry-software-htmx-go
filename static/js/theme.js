// Theme toggle functionality
document.addEventListener('htmx:afterSettle', function(evt) {
    if (evt.detail.requestConfig.path.startsWith('/toggle-theme')) {
        const responseURL = evt.detail.xhr.responseURL;
        const urlParams = new URL(responseURL).searchParams;
        const isDark = urlParams.get('dark') === 'true';
        if (isDark) {
            document.documentElement.classList.add('dark');
        } else {
            document.documentElement.classList.remove('dark');
        }
    }
});