<script>
  import { onMount } from 'svelte';
  import Navbar from '../components/Navbar.svelte';
  import { api } from '../lib/api.js';

  let requests = [];
  let error = '';

  onMount(async () => {
    try {
      requests = await api('/requests');
    } catch (err) {
      error = err.message;
    }
  });
</script>

<Navbar />

<h2>My Shift Requests</h2>
{#if error}
  <div class="error">{error}</div>
{:else if requests.length === 0}
  <div>No requests yet.</div>
{:else}
  <table>
    <thead>
      <tr>
        <th>Request ID</th>
        <th>Shift ID</th>
        <th>Status</th>
      </tr>
    </thead>
    <tbody>
      {#each requests as req}
        <tr>
          <td>{req.id}</td>
          <td>{req.shift_id}</td>
          <td>{req.status}</td>
        </tr>
      {/each}
    </tbody>
  </table>
{/if}

<style>
  .error { color: red; }
  table { border-collapse: collapse; margin-top: 1rem; }
  th, td { border: 1px solid #ccc; padding: 0.5rem 1rem; }
</style>
