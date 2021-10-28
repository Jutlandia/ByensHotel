document.addEventListener('DOMContentLoaded', () => {
	const navbarBurger = document.querySelector('.navbar-burger');
	const navbarMenu = document.querySelector('.navbar-menu');

	// removes the class is-active on initial load because is-active is active by default
	// for users that has JavaScript turned off
	if (navbarBurger.classList.contains('is-active') && navbarMenu.classList.contains('is-active')) {
		navbarBurger.classList.remove('is-active');
		navbarMenu.classList.remove('is-active');
	}

	navbarBurger.addEventListener('click', () => {
		navbarBurger.classList.toggle('is-active');
		navbarMenu.classList.toggle('is-active');
	});
});
