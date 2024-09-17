<script lang="ts">
  import { Card, Radio } from "flowbite-svelte";
  import { Typography } from "$lib";
  import { datas } from "../../../routes/settings/dummy";

  let {
    selectionData,
    selectedData = $bindable(""),
    numGrid = 1,
    isFlex = false,
    isRadio = false,
    widthFlex = "min-w-80",
    spacing = "",
    widgetIndicator,
  }: {
    selectionData: any;
    selectedData?: any;
    numGrid?: any;
    isFlex?: any;
    isRadio?: any;
    widthFlex?: any;
    spacing?: any;
    widgetIndicator?: any;
  } = $props();
  const selectData = (id: any) => {
    selectedData = id;
    console.log(selectedData);
  };
</script>

<div
  class={`${isFlex ? "flex !justify-start" : `grid grid-cols-${numGrid}`} w-full gap-4 items-center justify-center`}
>
  {#each selectionData as data, i}
    <Card
      horizontal={true}
      class={`justify-between ${isFlex ? `gap-3 ${widthFlex}` : "w-full h-full"} ${spacing ? spacing : "!py-6 !pr-4"} bg-primary-light dark:bg-primary-dark hover:bg-primary-hover-light dark:hover:bg-primary-hover-dark ${
        selectedData === data.id
          ? "border-brand-solid-light dark:border-brand-solid-dark"
          : "border-secondary-light dark:border-secondary-dark"
      } shadow-sm hover:no-underline cursor-pointer`}
      size="none"
      on:click={() => selectData(data?.id)}
    >
      <div class="flex items-center gap-4">
        {#if data.icon}
          <svelte:component
            this={data.icon}
            svgClass={`w-7 h-7 ${selectedData === data.id ? "text-fg-brand-primary-alt-light dark:text-fg-brand-primary-alt-dark" : "text-utility-gray-500-light dark:text-utility-gray-500-dark"}`}
          />
        {/if}
        {#if widgetIndicator}
          <svelte:component this={widgetIndicator} color={data.value} size="lg" />
        {/if}
        <div class="flex-col">
          <div class="flex flex-col gap-[2px]">
            <Typography
              preset="Text md/Semibold"
              text={data.name}
              color={`${selectedData === data.id ? "brand-primary" : "primary"}`}
            />
            {#if data.description}
              <Typography
                preset="Text xs/Regular"
                text={data.description}
                color={`${selectedData === data.id ? "brand-secondary" : "secondary"}`}
              />
            {/if}
          </div>
        </div>
      </div>
      {#if isRadio}
        <Radio value={data.id} bind:group={selectedData}></Radio>
      {/if}
    </Card>
  {/each}
</div>
