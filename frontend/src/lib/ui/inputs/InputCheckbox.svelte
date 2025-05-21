<script lang="ts">
  import { Checkbox } from 'flowbite-svelte';
  let {
    checked = $bindable(false),
    required
  }: {
    checked: boolean;
    required?: boolean;
  } = $props();
</script>

<Checkbox class="display-block" bind:checked={checked} required={required} on:click={() => checked = !checked}/>