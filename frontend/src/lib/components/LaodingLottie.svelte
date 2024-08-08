<script>
  import { onMount } from 'svelte';
  import lottie from "lottie-web";
  import LottieHexa from '$lib/components/hexa_lottie.json';

  export let animationData = LottieHexa;
  let container;

  onMount(() => {
    lottie.loadAnimation({
      container,
      animationData,
      renderer: 'svg',
      loop: true,
      autoplay: true,
    });
  });
</script>

<div bind:this={container} style="width: 57px;"></div>