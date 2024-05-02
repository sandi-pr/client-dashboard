<script lang="ts">
  import { Breadcrumb, BreadcrumbItem } from "flowbite-svelte";
  import { onMount } from "svelte";

  const {
    activeUrl = "/",
    linkClass = "ms-1 text-sm font-medium text-gray-700 hover:text-gray-900 md:ms-2 dark:text-gray-400 dark:hover:text-white",
    spanClass = "ms-1 text-sm font-medium text-gray-500 md:ms-2 dark:text-gray-400",
  } = $props();
  let accumulatedPath = $state("");
  let params = $state("");
  let queryString = $state("");
  let activeUrlArray: any[] = $state([]);
  let ready = $state(false);

  onMount(() => {
    params = new URLSearchParams(window.location.search).toString();
    queryString = params.toString() ? "?" + params.toString() : "";
    activeUrlArray = activeUrl
      .split("/")
      .filter((url) => url !== "")
      .map((url, index, array) => {
        accumulatedPath += `/${url}`;
        let href = accumulatedPath;
        if (index === array.length - 1) {
          href += queryString;
        }
        return { name: url, href: href };
      });
    if (activeUrlArray.some((item) => item.name === "client-dashboard")) {
      activeUrlArray = activeUrlArray.filter((item) => item.name !== "client-dashboard");
    }
    ready = true;
  });
</script>

<Breadcrumb aria-label="Default breadcrumb">
  <BreadcrumbItem href="/" home>Home</BreadcrumbItem>
  {#if activeUrl !== "/"}
    {#each activeUrlArray as item}
      <li class={"inline-flex items-center"}>
        <svg
          class="w-6 h-6 text-gray-400 rtl:-scale-x-100"
          fill="currentColor"
          viewBox="0 0 20 20"
          xmlns="http://www.w3.org/2000/svg"
        >
          <path
            fill-rule="evenodd"
            d="M7.293 14.707a1 1 0 010-1.414L10.586 10 7.293 6.707a1 1 0 011.414-1.414l4 4a1 1 0 010 1.414l-4 4a1 1 0 01-1.414 0z"
            clip-rule="evenodd"
          />
        </svg>
        {#if item.href}
          <button
            type="button"
            class={linkClass}
            onclick={() => {
              window.location = item.href;
            }}
          >
            {item.name}
          </button>
        {:else}
          <span class={spanClass}>
            {item.name}
          </span>
        {/if}
      </li>
    {/each}
  {/if}
</Breadcrumb>

<!--
@component
## Config
@prop activeUrl: string; // current active url
@prop linkClass: string; // class for the link
@prop spanClass: string; // class for the span
-->
