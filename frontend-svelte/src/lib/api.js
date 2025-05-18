export async function api(endpoint, method = 'GET', data = null) {
  const opts = {
    method,
    headers: {
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${localStorage.getItem('token')}`
    },
    body: data ? JSON.stringify(data) : undefined
  };

  const response = await fetch(`http://localhost:8000${endpoint}`, opts);
  
  if (!response.ok) {
    const error = await response.text();
    throw new Error(error);
  }
  
  // Only try to parse JSON if there is a body
  const text = await response.text();
  if (!text) return null; // No body, just return null or true
  try {
    return JSON.parse(text);
  } catch {
    return text; // Fallback: return raw text
  }
}
