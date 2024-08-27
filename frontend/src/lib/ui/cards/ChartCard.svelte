<script>
  import { Card } from "flowbite-svelte";
  import { Typography, TextButton, showAI } from "$lib";
  import Chart from "$lib/components/Chart.svelte";
  import { _ } from "svelte-i18n";

  let { chartAction = [], chartList = [] } = $props();
</script>

<div
  class={`grid ${$showAI ? "grid-cols-1 xl:grid-cols-2" : "grid-cols-1 md:grid-cols-2 xl:grid-cols-3"} gap-4 items-center justify-center`}
>
  {#each chartAction as { href, dataName, value, change, changeType }, index}
    <Card
      horizontal="true"
      class="grid grid-cols-12 items-center !px-4 !py-0 min-h-28 justify-between bg-primary-light dark:bg-primary-dark border-secondary-light dark:border-secondary-dark shadow-sm"
      size="none"
    >
      <div class="space-y-2 col-span-6 sm:col-span-7 md:col-span-5">
        <div class="flex items-center gap-2">
          <Typography preset="Text sm/Semibold" text={dataName} color="secondary" />
          <Typography
            preset="Text sm/Semibold"
            text={`${changeType === "decrease" ? "-" : "+"}${change}`}
            color={`${changeType === "decrease" ? "error-primary" : "success-primary"}`}
          />
        </div>
        <Typography preset="Text xl/Semibold" text={value} />
        <TextButton {href} text={$_("seeAll")} />
      </div>
      <div class="col-span-6 sm:col-span-5 md:col-span-7">
        {#each chartList as chart, i}
          {#if i === index}
            <Chart dataChart={chart} />
          {/if}
        {/each}
      </div>
    </Card>
  {/each}
</div>
