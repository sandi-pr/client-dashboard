<script>
  import { onMount } from "svelte";
  import LottieHexa from "$lib/components/hexa_lottie.json";

  export let animationData = LottieHexa;
  export let width = 57; // Default width in pixels
  let container;

  onMount(async () => {
    if (typeof window !== "undefined") {
      const lottie = await import("lottie-web");
      lottie.default.loadAnimation({
        container,
        animationData,
        renderer: "svg",
        loop: true,
        autoplay: true,
      });
    }
  });
</script>

<div bind:this={container} style="width: {width}px;"></div>
