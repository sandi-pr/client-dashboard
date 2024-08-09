<script lang="ts">
  import { Card, Radio } from "flowbite-svelte";
  import { Typography } from "$lib";

  let { selectionData, selectedData = $bindable(""), numGrid = 1, isFlex = false, isRadio = false, widthFlex = "min-w-80" }: { selectionData: any, selectedData?: any, numGrid?: any, isFlex?: any, isRadio?: any, widthFlex?: any } = $props();
  const selectData = (index: any) => {
    selectedData = "selectedData_" + index;
    console.log(selectedData);
  };
</script>

<div class={`${isFlex ? "flex !justify-start" : `grid grid-cols-${numGrid}`} w-full gap-6 items-center justify-center`}>
  {#each selectionData as data, i}
    <Card
      horizontal={true}
      class={`justify-between ${isFlex ? `gap-3 ${widthFlex}` : "w-full"} !py-6 !pr-4 bg-primary-light dark:bg-primary-dark hover:bg-primary-hover-light dark:hover:bg-primary-hover-dark ${
        selectedData === `selectedData_${i}`
          ? "border-brand-solid-light dark:border-brand-solid-dark"
          : "border-secondary-light dark:border-secondary-dark"
      } shadow-sm hover:no-underline cursor-pointer`}
      size="none"
      on:click={() => selectData(i)}
    >
      <div class="flex items-start gap-4">
        {#if data.icon}
          <svelte:component this={data.icon} svgClass={`w-7 h-7 ${selectedData === `selectedData_${i}` ? "text-fg-brand-primary-alt-light dark:text-fg-brand-primary-alt-dark" : "text-utility-gray-500-light dark:text-utility-gray-500-dark"}`} />
        {/if}
        <div class="flex-col">
          <div class="flex flex-col gap-[2px]">
            <Typography preset="Text md/Semibold" text={data.name} color={`${selectedData === `selectedData_${data.value}` ? "brand-primary" : "primary"}`} />
            {#if data.description}
              <Typography preset="Text xs/Regular" text={data.description} color={`${selectedData === `selectedData_${data.value}` ? "brand-secondary" : "secondary"}`} />
            {/if}
          </div>
        </div>
      </div>
      {#if isRadio}
        <Radio name="companyType" value={`selectedData_${i}`} bind:group={selectedData}></Radio>
      {/if}
    </Card>
  {/each}
</div>