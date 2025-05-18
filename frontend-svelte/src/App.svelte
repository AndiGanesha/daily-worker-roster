<script>
  import Login from './routes/login.svelte';
  import Dashboard from './routes/dashboard.svelte';
  import AdminDashboard from './routes/admin-dashboard.svelte';

  // Helper to decode JWT and get user info
  function getUserFromToken() {
    const token = localStorage.getItem('token');
    if (!token) return null;
    try {
      const payload = JSON.parse(atob(token.split('.')[1]));
      return payload;
    } catch {
      return null;
    }
  }

  let user = getUserFromToken();
  let page = user
    ? (user.role === 'admin' ? 'admin' : 'dashboard')
    : 'login';

  // Handle login event from Login.svelte
  function handleLogin(event) {
    user = getUserFromToken();
    page = user && user.role === 'admin' ? 'admin' : 'dashboard';
  }

  function handleLogout() {
    localStorage.removeItem('token');
    user = null;
    page = 'login';
  }

  // Optional: Protect against manual navigation
  $: if (page === 'dashboard' && (!user || user.role !== 'worker')) {
    page = 'login';
  }
  $: if (page === 'admin' && (!user || user.role !== 'admin')) {
    page = 'login';
  }
</script>

<!-- Simple navigation bar -->
<nav>
  {#if user}
    <button on:click={handleLogout}>Logout</button>
    {#if user.role === 'admin'}
      <button on:click={() => page = 'admin'}>Admin Dashboard</button>
    {:else}
      <button on:click={() => page = 'dashboard'}>Dashboard</button>
    {/if}
  {/if}
</nav>

<!-- Page routing -->
{#if page === 'login'}
  <Login on:login={handleLogin} />
{:else if page === 'dashboard'}
  <Dashboard />
{:else if page === 'admin'}
  <AdminDashboard />
{/if}

<style>
  nav {
    background: #f3f3f3;
    padding: 1rem;
    display: flex;
    gap: 1rem;
    justify-content: flex-end;
  }
  button {
    padding: 0.5rem 1rem;
    font-size: 1rem;
  }
</style>
