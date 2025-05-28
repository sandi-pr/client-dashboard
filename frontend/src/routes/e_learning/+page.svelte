<script>
  import { writable } from "svelte/store";
  const isFrontend = process.env.NODE_ENV === "development";

  let iframeSrc = isFrontend
    ? writable(`http://158.140.178.96:9191`)
    : writable(`http://${window.location.hostname}:9292`);
  let ref;
</script>

<iframe bind:this={ref} title="e-learning" class="w-full h-full" src={$iframeSrc}></iframe>
