<script lang="ts">
  import { Search } from 'flowbite-svelte';
  let {
    value,
    placeholder
  }: {
    value: string;
    placeholder: string | undefined;
  } = $props();
</script>

<Search placeholder={placeholder} bind:value={value}></Search>
