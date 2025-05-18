import { writable } from 'svelte/store';

export const auth = writable({
  loggedIn: false,
  role: null,
  name: null
});

export function checkAuth() {
  const token = localStorage.getItem('token');
  if (!token) return;
  
  try {
    const payload = JSON.parse(atob(token.split('.')[1]));
    auth.set({
      loggedIn: true,
      role: payload.role,
      name: payload.name
    });
  } catch {
    localStorage.removeItem('token');
  }
}
