<script>
  import { onMount } from 'svelte';
  import { api } from '../lib/api.js';

  let shifts = [];
  let pendingRequests = [];
  let error = '';
  let newShift = { date: '', start: '', end: '', role: '', location: '' };

  // For editing
  let editingShift = null;
  let editShiftData = {};

  // Load all shifts and pending requests
  async function loadAll() {
    error = '';
    try {
      const allShifts = await api('/shifts');
      shifts = Array.isArray(allShifts) ? allShifts : [];
    } catch (err) {
      error += `\nError loading shifts: ${err.message}`;
    }
    try {
      const reqs = await api('/admin/requests');
      pendingRequests = Array.isArray(reqs) ? reqs : [];
    } catch (err) {
      error += `\nError loading requests: ${err.message}`;
    }
  }

  onMount(loadAll);

  // Create a new shift
  async function createShift() {
    error = '';
    try {
      await api('/admin/shifts', 'POST', newShift);
      newShift = { date: '', start: '', end: '', role: '', location: '' };
      await loadAll();
    } catch (err) {
      error = `Create error: ${err.message}`;
    }
  }

  // Edit a shift
  function startEdit(shift) {
    editingShift = shift.id;
    editShiftData = { ...shift };
  }
  async function saveEdit() {
    try {
      await api(`/admin/shifts/${editingShift}`, 'PUT', editShiftData);
      editingShift = null;
      await loadAll();
    } catch (err) {
      error = `Edit error: ${err.message}`;
    }
  }
  function cancelEdit() {
    editingShift = null;
  }

  // Delete a shift
  async function deleteShift(id) {
    if (!confirm('Delete this shift?')) return;
    try {
      await api(`/admin/shifts/${id}`, 'DELETE');
      await loadAll();
    } catch (err) {
      error = `Delete error: ${err.message}`;
    }
  }

  // Approve/reject a request
  async function handleRequest(id, action) {
    try {
      await api(`/requests/${id}/${action}`, 'POST');
      await loadAll();
    } catch (err) {
      error = `Request action error: ${err.message}`;
    }
  }
</script>

<h1>Admin Dashboard</h1>

{#if error}
  <pre style="color: red">{error}</pre>
{/if}

<!-- Create Shift Form -->
<h2>Create New Shift</h2>
<form on:submit|preventDefault={createShift}>
  <input type="date" bind:value={newShift.date} required>
  <input type="time" bind:value={newShift.start} required>
  <input type="time" bind:value={newShift.end} required>
  <input type="text" bind:value={newShift.role} placeholder="Role" required>
  <input type="text" bind:value={newShift.location} placeholder="Location">
  <button type="submit">Create Shift</button>
</form>

<!-- All Shifts List -->
<h2>All Shifts</h2>
<table>
  <thead>
    <tr>
      <th>Date</th><th>Start</th><th>End</th><th>Role</th><th>Location</th><th>Actions</th>
    </tr>
  </thead>
  <tbody>
    {#each shifts as shift}
      {#if editingShift === shift.id}
        <tr>
          <td><input type="date" bind:value={editShiftData.date}></td>
          <td><input type="time" bind:value={editShiftData.start}></td>
          <td><input type="time" bind:value={editShiftData.end}></td>
          <td><input type="text" bind:value={editShiftData.role}></td>
          <td><input type="text" bind:value={editShiftData.location}></td>
          <td>
            <button on:click={saveEdit}>Save</button>
            <button on:click={cancelEdit}>Cancel</button>
          </td>
        </tr>
      {:else}
        <tr>
          <td>{shift.date}</td>
          <td>{shift.start}</td>
          <td>{shift.end}</td>
          <td>{shift.role}</td>
          <td>{shift.location}</td>
          <td>
            <button on:click={() => startEdit(shift)}>Edit</button>
            <button on:click={() => deleteShift(shift.id)}>Delete</button>
          </td>
        </tr>
      {/if}
    {/each}
  </tbody>
</table>

<!-- Pending Requests -->
<h2>Pending Requests</h2>
<table>
  <thead>
    <tr>
      <th>User</th>
      <th>Date</th>
      <th>Start</th>
      <th>End</th>
      <th>Role</th>
      <th>Location</th>
      <th>Status</th>
      <th>Action</th>
    </tr>
  </thead>
  <tbody>
    {#each pendingRequests as req}
      <tr>
        <td>{req.user_name}</td>
        <td>{req.date}</td>
        <td>{req.start}</td>
        <td>{req.end}</td>
        <td>{req.role}</td>
        <td>{req.location}</td>
        <td>{req.status}</td>
        <td>
          <button on:click={() => handleRequest(req.request_id, 'approve')}>Approve</button>
          <button on:click={() => handleRequest(req.request_id, 'reject')}>Reject</button>
        </td>
      </tr>
    {/each}
  </tbody>
</table>

<style>
  table { border-collapse: collapse; margin-top: 1rem; }
  th, td { border: 1px solid #ccc; padding: 0.5rem 1rem; }
  form input { margin: 0 0.5rem 0.5rem 0; }
</style>
