<script>
  import { _ } from "svelte-i18n";
  import { Button } from "flowbite-svelte";
  import { Typography, Divider, DataTableAIReview } from "$lib";

  export let title = "";
  export let isCSV = false;
  export let datas = [];
  export let configs = {};
  export let next = () => {};
  export let previous = () => {};
  export let goToPage = () => {};
  export let goBackPage = () => {};
</script>

<div class="flex flex-col px-4 pt-6 gap-6 pb-12">
  <div class="flex justify-between">
    <div class="flex flex-col gap-1">
      <Typography preset="Display sm/Semibold" text={title} color="primary" />
    </div>
    <Button class="h-min" color="light" outline on:click={goBackPage}>
      <Typography preset="Text md/Semibold" text={$_("goBack")} color="primary" />
    </Button>
  </div>
  <Divider />
  <div class="flex gap-3">
    <div class="flex flex-col grow gap-6">
      <div class="flex flex-col gap-6 w-full">
        {#if isCSV && datas.length > 0}
          <DataTableAIReview bind:datas bind:configs {next} {previous} {goToPage} />
        {/if}
      </div>
    </div>
  </div>
</div>
