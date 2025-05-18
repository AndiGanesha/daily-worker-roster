<script>
  import { onMount } from 'svelte';
  import ShiftList from '../components/ShiftList.svelte';
  import Navbar from '../components/Navbar.svelte';
  import { api } from '../lib/api.js';

  let shifts = [];
  let error = '';
  let refresh = false;

  async function fetchShifts() {
    try {
      shifts = await api('/shifts');
      error = '';
    } catch (err) {
      error = err.message;
    }
  }

  onMount(fetchShifts);

  // This will be called by RequestButton when a request is made
  function handleRequested() {
    fetchShifts();
  }
</script>

<Navbar />

<h2>Available Shifts</h2>
{#if error}
  <div class="error">{error}</div>
{:else}
  <ShiftList {shifts} showRequestButton={true} on:requested={handleRequested} />
{/if}

<style>
  .error { color: red; }
</style>
