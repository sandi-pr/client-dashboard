<script lang="ts">
  import { Input } from 'flowbite-svelte';
  let {
    value,
    readonly = false
  }: {
    value: string;
    readonly?: boolean;
  } = $props();
</script>

<Input class="display-block" type="date" bind:value={value} readonly={readonly} />
