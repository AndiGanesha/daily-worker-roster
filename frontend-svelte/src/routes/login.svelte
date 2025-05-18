<script>
  import { api } from '../lib/api.js';
  import { createEventDispatcher } from 'svelte';

  let username = '';
  let password = '';
  let error = '';
  const dispatch = createEventDispatcher();

  async function handleLogin() {
    try {
      const response = await api('/login', 'POST', { username, password });
      localStorage.setItem('token', response.token);
      const payload = JSON.parse(atob(response.token.split('.')[1]));
      // You can pass the role to the parent if you want
      dispatch('login', { role: payload.role });
    } catch (err) {
      error = err.message;
    }
  }
</script>

<div class="login-container">
  <h2>Worker Login</h2>
  <form on:submit|preventDefault={handleLogin}>
    <input type="text" bind:value={username} placeholder="Username" required>
    <input type="password" bind:value={password} placeholder="Password" required>
    <button type="submit">Login</button>
  </form>
  {#if error}<div class="error">{error}</div>{/if}
</div>

<style>
  .login-container { max-width: 400px; margin: 2rem auto; }
  .error { color: red; margin-top: 1rem; }
</style>
