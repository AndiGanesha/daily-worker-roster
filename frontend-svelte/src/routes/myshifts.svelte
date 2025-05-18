<script>
  import { onMount } from 'svelte';
  import ShiftList from '../components/ShiftList.svelte';
  import Navbar from '../components/Navbar.svelte';
  import { api } from '../lib/api.js';

  let shifts = [];
  let error = '';

  onMount(async () => {
    try {
      shifts = await api('/myshifts');
    } catch (err) {
      error = err.message;
    }
  });
</script>

<Navbar />

<h2>My Assigned Shifts</h2>
{#if error}
  <div class="error">{error}</div>
{:else}
  <ShiftList {shifts} />
{/if}

<style>
  .error { color: red; }
</style>
