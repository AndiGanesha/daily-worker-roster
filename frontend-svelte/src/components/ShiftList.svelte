<script>
  import RequestButton from './Requestbutton.svelte';
  export let shifts = [];
  export let showRequest = false;
</script>

{#if !shifts || shifts.length === 0}
  <div>No shifts found.</div>
{:else}
  <table>
    <!-- table headers -->
    <tbody>
      {#each shifts as shift}
        <tr>
          <td>{shift.date}</td>
          <td>{shift.start}</td>
          <td>{shift.end}</td>
          <td>{shift.role}</td>
          <td>{shift.location}</td>
          {#if showRequest}
            <td>
              <RequestButton shiftId={shift.id} />
            </td>
          {/if}
        </tr>
      {/each}
    </tbody>
  </table>
{/if}
