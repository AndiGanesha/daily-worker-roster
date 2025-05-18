const API_BASE = import.meta.env.VITE_API_BASE || 'http://localhost:8080';

export async function api(endpoint, method = 'GET', data = null) {
  const opts = {
    method,
    headers: { 'Content-Type': 'application/json' },
  };
  const token = localStorage.getItem('token');
  if (token) opts.headers['Authorization'] = `Bearer ${token}`;
  if (data) opts.body = JSON.stringify(data);

  const response = await fetch(API_BASE + endpoint, opts);

  if (!response.ok) {
    throw new Error(await response.text() || response.statusText);
  }

  const text = await response.text();
  if (!text) return null;
  try {
    return JSON.parse(text);
  } catch {
    return text;
  }
}
