<script>
  import { Button } from "flowbite-svelte";
  import { _ } from "svelte-i18n";
  import ArrowLeftIcon from "$lib/icon/arrow-left.svelte";
  import ArrowRightIcon from "$lib/icon/arrow-right.svelte";

  let { pagination, previousPage, nextPage } = $props();
</script>

<div class="flex justify-between !border-0">
  {#if pagination.total > 0}
    <div class="m-4 text-sm text-hc-graylight-500 dark:text-gray-400">
      {$_("showing")}
      <span>{pagination.page === 1 ? 1 : (pagination.page - 1) * pagination.pageSize + 1}</span>
      -
      <span
        >{pagination.page * pagination.pageSize >= pagination.total
          ? pagination.total
          : pagination.page * pagination.pageSize}</span
      >
      {$_("of")}
      <span>{pagination.total}</span>
    </div>

    <div class="m-4 text-sm text-hc-graylight-500 dark:text-gray-400">
      {$_("page")}
      <span>{pagination.page}</span>
      {$_("of")}
      <span>{pagination.pageCount}</span>
    </div>
  {/if}
</div>
{#if pagination.pageCount > 1}
  <div class="flex justify-end !border-0">
    <Button
      disabled={pagination.page === 1}
      on:click={previousPage}
      color="light"
      class="fluid-sm gap-1 md:gap-2 rounded-l-lg rounded-r-none px-3 py-2 md:px-4 md:py-2.5"
    >
      <ArrowLeftIcon
        svgClass="text-fg-secondary-light dark:text-fg-secondary-dark w-4 h-4 md:w-5 md:h-5"
      />
      {$_("previous")}
    </Button>
    <Button
      disabled={pagination.page * pagination.pageSize >= pagination.total}
      on:click={nextPage}
      color="light"
      class="fluid-sm gap-1 md:gap-2 rounded-r-lg rounded-l-none px-3 py-2 md:px-4 md:py-2.5"
    >
      {$_("next")}
      <ArrowRightIcon
        svgClass="text-fg-secondary-light dark:text-fg-secondary-dark w-4 h-4 md:w-5 md:h-5"
      />
    </Button>
  </div>
{/if}
