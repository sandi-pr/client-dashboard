<script>
  import { _ } from "svelte-i18n";
  import { Button } from "flowbite-svelte";
  import { windowSize } from "$lib";
  import ArrowLeftIcon from "$lib/icon/arrow-left.svelte";
  import ArrowRightIcon from "$lib/icon/arrow-right.svelte";

  export let configs = {};
  export let goToPage = () => {};
  export let next = () => {};
  export let previous = () => {};

  $: width = $windowSize.width;

  const getPaginationNumbers = (current, total) => {
    if (total < 7) {
      return Array.from({ length: total }, (_, i) => i + 1);
    }

    if (current <= 3) {
      return [1, 2, 3, "...", total - 2, total - 1, total];
    }

    if (current >= total - 2) {
      return [1, "...", total - 3, total - 2, total - 1, total];
    }

    return [1, "...", current - 1, current, current + 1, "...", total];
  };
</script>

{#if configs?.pagination && Object.keys(configs.pagination).length > 0}
  <div class="flex justify-between items-center !border-t-0 !rounded-t-none rounded-b-lg mt-4">
    <div class="mx-4 text-sm text-hc-graylight-500 dark:text-gray-400">
      {$_("showing")}
      <span>
        {configs.pagination.page === 1
          ? configs.pagination.total === 0
            ? 0
            : 1
          : (configs.pagination.page - 1) * configs.pagination.pageSize + 1}</span
      >
      -
      <span>
        {configs.pagination.page * configs.pagination.pageSize >= configs.pagination.total
          ? configs.pagination.total
          : configs.pagination.page * configs.pagination.pageSize}</span
      >
      {$_("of")}
      <span>{configs.pagination.total}</span>
    </div>
    <div class="flex">
      {#if configs.pagination.total > configs.pagination.pageSize}
        <div class="flex justify-end mr-2">
          <Button
            disabled={configs.pagination.page === 1}
            on:click={previous}
            color="light"
            class="fluid-sm gap-2 rounded-l-lg border-0 rounded-r-none px-4 py-2 focus:ring-0"
          >
            <ArrowLeftIcon
              svgClass="text-fg-secondary-light dark:text-fg-secondary-dark w-4 h-4 md:w-5 md:h-5"
            />
            {$_("previous")}
          </Button>

          {#if width > 768}
            {#each getPaginationNumbers(configs.pagination.page, configs.pagination.pageCount) as item, idx}
              {#if item === "..."}
                <Button
                  disabled
                  color="light"
                  class="px-4 py-2 border-0 rounded-none cursor-default">...</Button
                >
              {:else}
                <Button
                  color={item === configs.pagination.page ? "primary" : "light"}
                  class={`px-4 py-2 border-0 rounded-none focus:ring-0 ${item === configs.pagination.page ? "font-bold" : ""}`}
                  on:click={() => {
                    if (item !== configs.pagination.page) {
                      goToPage(item);
                    }
                  }}
                >
                  {item}
                </Button>
              {/if}
            {/each}
          {/if}

          <Button
            disabled={configs.pagination.page * configs.pagination.pageSize >=
              configs.pagination.total}
            on:click={next}
            color="light"
            class="fluid-sm gap-2 rounded-r-lg border-0 rounded-l-none px-4 py-2 focus:ring-0"
          >
            {$_("next")}
            <ArrowRightIcon
              svgClass="text-fg-secondary-light dark:text-fg-secondary-dark w-4 h-4 md:w-5 md:h-5"
            />
          </Button>
        </div>
      {/if}
    </div>
  </div>
{/if}
