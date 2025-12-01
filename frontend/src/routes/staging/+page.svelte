<script>
  import { onMount, onDestroy } from "svelte";
  import { _ } from "svelte-i18n";
  import { Typography } from "$lib";
  import { getIcons } from "$lib";
  import LaodingLottie from "$lib/components/LaodingLottie.svelte";
  import StagingLottie from "$lib/components/system_start_lottie.json";

  let RefreshIcon = getIcons("refreshCw05");

  let innerWidth = 0;
  let refreshReady = false;
  let refreshTimeout;

  function updateWidth() {
    innerWidth = window.innerWidth;
  }

  function startRefreshTimer() {
    refreshTimeout = setTimeout(() => {
      refreshReady = true;
    }, 120000); // 1 minute = 60000 ms
  }

  onMount(() => {
    updateWidth();
    window.addEventListener("resize", updateWidth);
    startRefreshTimer();
  });

  onDestroy(() => {
    window.removeEventListener("resize", updateWidth);
  });

  // $: console.log(`Current innerWidth: ${innerWidth}`);
</script>

<div class="flex flex-col justify-center items-center h-screen md:gap-4">
  {#if innerWidth > 768}
    <LaodingLottie width={360} animationData={StagingLottie} />
    <Typography preset="Display xl/Semibold" text={$_("theSystemIsStarting")} />
    <Typography
      preset="Text xl/Regular"
      color="tertiary"
      text={$_("justAMomentWhileWeGetThingsReady")}
    />
  {:else}
    <LaodingLottie width={240} animationData={StagingLottie} />
    <Typography preset="Display xs/Semibold" text={$_("theSystemIsStarting")} />
    <Typography
      preset="Text md/Regular"
      color="tertiary"
      text={$_("justAMomentWhileWeGetThingsReady")}
    />
  {/if}
  {#if refreshReady}
    <button
      class="flex gap-2 items-center px-4 py-2 bg-button-secondary-light dark:bg-button-secondary-dark rounded-lg border border-secondary-light dark:border-secondary-dark hover:bg-button-secondary-hover-light dark:hover:bg-button-secondary-hover-dark"
      on:click={() => {
        refreshReady = false;
        startRefreshTimer();
        window.location.reload();
      }}
    >
      <RefreshIcon svgClass="w-5 h-5 text-fg-secondary-light dark:text-fg-secondary-dark" />
      <Typography preset="Text md/Semibold" text="Refresh" color="fg-secondary" />
    </button>
  {/if}
</div>
