<script lang="ts">
  import { Card, Radio } from "flowbite-svelte";
  import { Typography } from "$lib";

  let { selectionData, selectedData }: { selectionData: any, selectedData?: any } = $props();
  const selectData = (index: any) => {
    selectedData = "selectedData_" + index;
  };
</script>

<div class="grid grid-cols-1 gap-6 w-full items-center justify-center">
  {#each selectionData as { icon, name, desc }, i}
    <Card
      horizontal={true}
      class={`justify-between w-full !p-6 min-h-28 bg-primary-light dark:bg-primary-dark hover:bg-primary-hover-light dark:hover:bg-primary-hover-dark ${
        selectedData === `selectedData_${i}`
          ? "border-brand-solid-light dark:border-brand-solid-dark"
          : "border-secondary-light dark:border-secondary-dark"
      } shadow-sm hover:no-underline cursor-pointer`}
      size="none"
      on:click={() => selectData(i)}
    >
      <div class="flex items-start gap-4">
        <svelte:component this={icon} svgClass={`w-7 h-7 ${selectedData === `selectedData_${i}` ? "text-fg-brand-primary-alt-light dark:text-fg-brand-primary-alt-dark" : "text-utility-gray-500-light dark:text-utility-gray-500-dark"}`} />
        <div class="flex-col">
          <div class="flex flex-col gap-[2px]">
            <Typography preset="Text md/Semibold" text={name} color={`${selectedData === `selectedData_${i}` ? "brand-primary" : "primary"}`} />
            <Typography preset="Text xs/Regular" text={desc} color={`${selectedData === `selectedData_${i}` ? "brand-secondary" : "secondary"}`} />
          </div>
        </div>
      </div>
      <Radio name="companyType" value={`selectedData_${i}`} bind:group={selectedData}></Radio>
    </Card>
  {/each}
</div>