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

document.addEventListener('DOMContentLoaded', () => {
    const isDark = getThemePreference()
    if (isDark) {
        document.documentElement.classList.add('dark');
    }
})

function getThemePreference() {
    const cookies = document.cookie.split(';');
    for (let cookie of cookies) {
        const [name, value] = cookie.trim().split('=');
        if (name === 'theme') {
            return value === 'dark';
        }
    }
    return window.matchMedia('(prefers-color-scheme: dark)').matches;
}
