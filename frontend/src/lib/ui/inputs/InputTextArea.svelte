<script lang="ts">
  import { Textarea } from 'flowbite-svelte';
  let {
    value,
    placeholder,
    required,
    readonly = false
  }: {
    value: string;
    placeholder: string | undefined;
    required?: boolean;
    readonly?: boolean;
  } = $props();
</script>

<Textarea class="display-block" bind:value {placeholder} {required} {readonly} />
