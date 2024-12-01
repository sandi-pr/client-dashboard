<script>
  import { Breadcrumb, BreadcrumbItem } from "flowbite-svelte";
  import { goto } from "$app/navigation";
  import { onMount } from "svelte";
  import { CR3Call } from "hexacomponents-svelte";
  import { getIcons } from "$lib";
  import { toCamelCase } from "$lib/utils/toCamelCase.js";
  import ChevronRightIcon from "$lib/icon/chevron-right.svelte";
  import { _ } from "svelte-i18n";

  let {
    activeUrl = "/",
    appUrl = "/",
    linkClass = "md:ms-1 fluid-xs md:fluid-sm line-clamp-1 font-medium text-gray-700 hover:text-gray-900 md:ms-2 dark:text-gray-400 dark:hover:text-white",
    spanClass = "md:ms-1 fluid-xs md:fluid-sm line-clamp-1 font-medium text-gray-500 md:ms-2 dark:text-gray-400",
  } = $props();
  let baseUrl = $state("");
  let params = $state(new URLSearchParams(window.location.search).toString());
  let queryString = $state(params.toString() ? "?" + params.toString() : "");
  let activeUrlArray = $state([]);
  let copyActiveUrl = $state("");
  let activeUrlModified = $state("");
  let detail = $state(false);
  let title = $state("");
  let isMobile = $state(false);

  import { call } from "$lib";

  const postProcessName = (name) => {
    let newName = name.replace(/-/g, " ");
    newName = newName.replace(/_/g, " ");

    if (name === "crm") {
      return name.toUpperCase();
    } else if (name.split(" ").length > 1) {
      return newName
        .split(" ")
        .map((word) => {
          return word.charAt(0).toUpperCase() + word.slice(1);
        })
        .join(" ");
    } else {
      return newName.charAt(0).toUpperCase() + newName.slice(1);
    }
  };

  const isDetailPage = (url) => {
    if (url.split("/").pop() === "detail") {
      console.log("detail");
      detail = true;
      const params = new URLSearchParams(window.location.search);
      title = params.get("title");
    } else {
      console.log("not detail");
      detail = false;
    }
  };

  const updateIsMobile = () => {
    isMobile = globalThis.innerWidth <= 768; // Assume mobile screen is less than or equal to 768px
  };

  onMount(() => {
    updateIsMobile();
    globalThis.addEventListener("resize", updateIsMobile);

    return () => {
      globalThis.removeEventListener("resize", updateIsMobile);
    };
  });

  onMount(async () => {
    const result = await $call.getPage("client-dashboard", "home-page");
    baseUrl = result.slice(0, result.lastIndexOf("/"));
  });

  $effect(() => {
    if (activeUrl !== copyActiveUrl) {
      console.log(activeUrl);
      isDetailPage(activeUrl);
      const cleanAppUrl = appUrl.replace(/^\/|\/$/g, "");
      copyActiveUrl = activeUrl;
      let cleanActiveUrl = activeUrl;
      activeUrlModified = cleanActiveUrl.replace(new RegExp(`^/${cleanAppUrl}`), "");
      activeUrlArray = activeUrlModified
        .split("/")
        .filter((url) => url !== "client-dashboard" && url !== "libraries")
        .filter(Boolean)
        .map((url, index, array) => {
          let accumulatedPath = "";
          for (let i = 0; i <= index; i++) {
            accumulatedPath += `/${array[i]}`;
          }
          let href = accumulatedPath;
          if (index === array.length - 1) {
            href += queryString;
          }
          return { name: url, href: href };
        });

      console.log("activeUrlArray", activeUrlArray);
    }
  });
</script>

<Breadcrumb aria-label="Default breadcrumb" class="">
  <!-- <BreadcrumbItem href="/" home></BreadcrumbItem> -->
  <li class={"inline-flex items-center"}>
    <button
      type="button"
      class={linkClass}
      onclick={() => {
        goto(`${baseUrl}`);
      }}
      ><svelte:component
        this={getIcons("home")}
        svgClass="w-4 h-4 md:w-5 md:h-5 text-gray-700 dark:text-gray-400"
      />
    </button>
  </li>
  {#if activeUrl !== "/"}
    {#each activeUrlArray as item, index (item.href)}
      <!-- Conditional rendering for mobile screen -->
      {#if isMobile && (index === 0 || index === activeUrlArray.length - 1)}
        <li class={"inline-flex items-center"}>
          <ChevronRightIcon svgClass="w-4 h-4 text-fg-senary-light dark:text-fg-senary-dark" />

          {#if item.href}
            <button
              type="button"
              disabled={index === activeUrlArray.length - 1}
              class={index === activeUrlArray.length - 1 ? `${linkClass} !font-bold` : linkClass}
              onclick={() => {
                const params = new URLSearchParams(window.location.search);
                const queryParams = {};
                for (const [key, value] of params.entries()) {
                  queryParams[key] = value;
                }
                const queryString = new URLSearchParams(queryParams).toString();
                goto(`${baseUrl}${item.href}?${queryString}`);
              }}
            >
              {#if detail && title !== null && index === activeUrlArray.length - 1}
                {title}
              {:else}
                {$_(toCamelCase(postProcessName(item.name)))}
              {/if}
            </button>
          {:else}
            <span class={spanClass}>
              {$_(toCamelCase(postProcessName(item.name)))}
            </span>
          {/if}
        </li>

        <!-- Render ellipsis between the first and the second-to-last item -->
        {#if isMobile && activeUrlArray.length > 1 && index === 0}
          <li class="inline-flex items-center">
            <ChevronRightIcon svgClass="w-4 h-4 text-fg-senary-light dark:text-fg-senary-dark" />
            <span class={spanClass}>...</span>
          </li>
        {/if}
      {:else if !isMobile}
        <!-- Render all items for non-mobile screens -->
        <li class={"inline-flex items-center"}>
          <ChevronRightIcon svgClass="w-4 h-4 text-fg-senary-light dark:text-fg-senary-dark" />

          {#if item.href}
            <button
              type="button"
              disabled={index === activeUrlArray.length - 1}
              class={index === activeUrlArray.length - 1 ? `${linkClass} !font-bold` : linkClass}
              onclick={() => {
                const params = new URLSearchParams(window.location.search);
                const queryParams = {};
                for (const [key, value] of params.entries()) {
                  queryParams[key] = value;
                }
                const queryString = new URLSearchParams(queryParams).toString();
                goto(`${baseUrl}${item.href}?${queryString}`);
              }}
            >
              {#if detail && title !== null && index === activeUrlArray.length - 1}
                {title}
              {:else}
                {$_(toCamelCase(postProcessName(item.name)))}
              {/if}
            </button>
          {:else}
            <span class={spanClass}>
              {$_(toCamelCase(postProcessName(item.name)))}
            </span>
          {/if}
        </li>
      {/if}
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
