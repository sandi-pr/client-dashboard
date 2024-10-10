<script>
  import { Card } from "flowbite-svelte";
  import { Typography, TextButton, showAI, FeaturedIcon, getIcons } from "$lib";
  import Chart from "$lib/components/Chart.svelte";
  import { _ } from "svelte-i18n";
  import { goto } from "$app/navigation";

  let { chartAction = [], chartList = [], isLibrary = false } = $props();
</script>

<div
  class={`grid ${$showAI ? "grid-cols-1 xl:grid-cols-2" : "grid-cols-1 md:grid-cols-2 xl:grid-cols-4"} gap-4 items-center justify-center`}
>
  {#each chartAction as { dataName, value, change, changeType, icon, title, description, href }, index}
    <Card
      class={`${isLibrary ? "w-full !py-4" : "grid grid-cols-12 min-h-28 !py-0"} items-center !px-4 justify-between bg-primary-light dark:bg-primary-dark border-secondary-light dark:border-secondary-dark shadow-sm`}
      size="none"
      on:click={isLibrary ? goto(href) : () => {}}
    >
      {#if isLibrary}
        <div class="flex w-full justify-start gap-3 hover:cursor-pointer">
          <div class="w-max">
            <FeaturedIcon classesSvg={getIcons(icon)} />
          </div>
          <div class="gap-[2px]">
            <Typography preset="Text md/Semibold" text={$_(title)} color="primary" />
            <Typography preset="Text sm/Regular" text={$_(description)} color="secondary" />
          </div>
        </div>
        <div class="col-span-12">
          {#each chartList as chart, i}
            {#if i === index}
              <Chart dataChart={chart} />
            {/if}
          {/each}
        </div>
        <div class="space-y-2 col-span-12">
          <div class="flex items-center gap-2">
            <Typography preset="Text sm/Semibold" text={dataName} color="secondary" />
            <Typography preset="Text sm/Semibold" text={value} />
            <Typography
              preset="Text sm/Semibold"
              text={`${changeType === "decrease" ? "" : "+"}${change}`}
              color={`${changeType === "decrease" ? "error-primary" : "success-primary"}`}
            />
          </div>
          <!-- <TextButton {href} text={$_("seeAll")} /> -->
        </div>
      {:else}
        <div class="space-y-2 col-span-6 sm:col-span-7 md:col-span-5">
          <div class="flex items-center gap-2">
            <Typography preset="Text sm/Semibold" text={dataName} color="secondary" />
            <Typography
              preset="Text sm/Semibold"
              text={`${changeType === "decrease" ? "" : "+"}${change}`}
              color={`${changeType === "decrease" ? "error-primary" : "success-primary"}`}
            />
          </div>
          <Typography preset="Text xl/Semibold" text={value} />
          <!-- <TextButton {href} text={$_("seeAll")} /> -->
        </div>
        <div class="col-span-6 sm:col-span-5 md:col-span-7">
          {#each chartList as chart, i}
            {#if i === index}
              <Chart dataChart={chart} />
            {/if}
          {/each}
        </div>
      {/if}
    </Card>
  {/each}
</div>
