<script lang="ts">
  import { P, Textarea } from "flowbite-svelte";
  import { twMerge } from "tailwind-merge";
  import { Typography, toCamelCase } from "$lib";
  import { _ } from "svelte-i18n";

  export let screenWidth: number;
  export let datas: {
    grid: number;
    title: string;
    titleGrid: number;
    innerGrid: number;
    data: {
      name: string;
      value: string;
      grid: number;
      disabled: boolean;
    }[];
  };

  const inputClass = (data) => {
    const baseGridClass =
      datas.innerGrid > data.grid
        ? screenWidth > 1024
          ? `col-span-${data.grid}`
          : `col-span-${datas.innerGrid}`
        : `col-span-${data.grid}`;

    const disabledClasses = data.disabled
      ? "[&>textarea]:text-disabled-light [&>textarea]:dark:text-disabled-dark [&>textarea]:bg-disabled-subtle-light [&>textarea]:dark:bg-disabled-subtle-dark [&>textarea]:border-disabled-subtle-light [&>textarea]:dark:border-disabled-subtle-dark"
      : "[&>*]:bg-primary-light [&>*]:dark:bg-primary-dark [&>*]:border-primary-light [&>*]:dark:border-primary-dark";

    return twMerge(baseGridClass, disabledClasses, "flex flex-col gap-2");
  };
</script>

<div class={`grid gap-4 ${screenWidth > 1024 ? "grid-cols-" + datas.grid : "grid-cols-1"}`}>
  <Typography
    preset="Text md/Semibold"
    text={$_(toCamelCase(datas.title))}
    color="secondary"
    propClass={`${screenWidth > 1024 ? "col-span-" + datas.titleGrid : "col-span-1"}`}
  />

  <div
    class={`col-span-${datas.grid - datas.titleGrid} grid gap-4 ${screenWidth > 1024 ? "grid-cols-" + datas.innerGrid : "grid-cols-1"}`}
  >
    {#each datas.data as data (data)}
      <div class={inputClass(data)}>
        <Typography text={$_(toCamelCase(data.name))} preset="Text sm/Medium" color="secondary" />
        <Textarea
          class="resize-none"
          rows={screenWidth > 1024 ? "1" : "2"}
          disabled={data.disabled}
          value={data.value}
        />
      </div>
    {/each}
  </div>
</div>

<!--
@component
## Config
@prop screenWidth: number; // Screen width (Needed for responsive design).
@prop datas: object; // Data for the information form.
-->
