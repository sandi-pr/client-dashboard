<script>
    import { writable } from "svelte/store";

    let iframeSrc = writable(`http://${window.location.hostname}:9292`);
    let ref;
    
</script>

<iframe bind:this={ref} title="e-learning" class="w-full h-full" src={$iframeSrc}></iframe>
