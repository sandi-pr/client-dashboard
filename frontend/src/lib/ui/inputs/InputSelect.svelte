<script lang="ts">
  import { Select } from 'flowbite-svelte';
  let {
    value,
    options,
    placeholder
  }: {
    value: string;
    options: { value: string; name: string }[];
    placeholder: string | undefined;
  } = $props();
</script>

<Select
  items={options}
  class="display-block"
  bind:value={value}
  placeholder={placeholder}></Select>
