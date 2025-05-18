<script>
  import { api } from '../lib/api';
  
  export let shiftId;
  let loading = false;
  let error = '';
  let requested = false; // <--- Track if this shift is already requested

  async function handleRequest() {
    loading = true;
    try {
      await api('/requests', 'POST', { shift_id: shiftId });
      error = '';
      requested = true; // <--- Set to true after success
    } catch (err) {
      error = err.message;
    }
    loading = false;
  }
</script>

<button on:click={handleRequest} disabled={loading || requested}>
  {loading ? 'Requesting...' : 'Request Shift'}
</button>
{#if error}<div class="error">{error}</div>{/if}

<style>
  .error { color: red; font-size: 0.8rem; }
</style>
