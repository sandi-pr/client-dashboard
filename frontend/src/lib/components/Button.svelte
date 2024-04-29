<script>
    import { Button } from 'flowbite-svelte'

    export let color = 'blue'

</script>

<Button color={color}>
    <slot></slot>
</Button>