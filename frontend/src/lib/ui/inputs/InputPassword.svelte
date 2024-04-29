<script lang="ts">
    import { Input } from 'flowbite-svelte';
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

<Input class="display-block" type="password" bind:value={value} placeholder={placeholder} required={required} readonly={readonly} />