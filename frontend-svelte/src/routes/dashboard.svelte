<script>
  import ShiftList from '../components/ShiftList.svelte';
  import { api } from '../lib/api.js';
  import { onMount } from 'svelte';

  let assignedShifts = [];
  let availableShifts = [];
  let error = '';

  onMount(async () => {
  try {
    let assigned = await api('/myshifts');
    assignedShifts = Array.isArray(assigned) ? assigned : [];
  } catch (err) {
    assignedShifts = [];
    error = "Error fetching assigned shifts: " + err.message;
  }

  try {
    let available = await api('/shifts');
    availableShifts = Array.isArray(available) ? available : [];
  } catch (err) {
    availableShifts = [];
    error += "\nError fetching available shifts: " + err.message;
  }
});
</script>

<div class="dashboard">
  <h1>Welcome to Your Dashboard</h1>
  
  <h2>Assigned Shifts</h2>
  <ShiftList shifts={assignedShifts}/>
  
  <h2>Available Shifts</h2>
  <ShiftList shifts={availableShifts} showRequest />
</div>
