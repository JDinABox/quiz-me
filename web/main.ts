import Alpine from 'alpinejs';

window.Alpine = Alpine;
window.Alpine.data('app', () => {
	return {
		theme: {
			dark: true,
			userToggledDark: false
		},
		init() {
			this.theme.userToggledDark = localStorage.userToggledDark === 'true';
			if (this.theme.userToggledDark) {
				const setDark = localStorage.dark ?? '';
				if (setDark.length > 0) {
					this.theme.dark = setDark === 'true';
					return;
				}
			}
			this.theme.dark = window.matchMedia('(prefers-color-scheme: dark)').matches;
		},
		toggleTheme() {
			this.theme.userToggledDark = true;
			localStorage.userToggledDark = 'true';
			this.theme.dark = !this.theme.dark;
			localStorage.dark = String(this.theme.dark);
		}
	};
});
window.Alpine.start();
