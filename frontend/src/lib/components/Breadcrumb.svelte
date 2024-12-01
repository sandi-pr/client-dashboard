<!-- <script>
  import { Breadcrumb } from "flowbite-svelte";
  import { goto } from "$app/navigation";
  import { onMount } from "svelte";
  import { CR3Call } from "hexacomponents-svelte";
  import { getIcons } from "$lib";
  import { toCamelCase } from "$lib/utils/toCamelCase.js";
  import { _ } from "svelte-i18n";
  import { derived, writable } from "svelte/store";
  import { breadcrumbLastName } from "$lib/stores/breadcrumbLastName.js";

  // Stores to manage the URL state
  const activeUrlStore = writable("/");
  const appUrlStore = writable("/");

  let {
    activeUrl = "/",
    appUrl = "/",
    linkClass = "ms-1 text-sm font-medium text-gray-700 hover:text-gray-900 md:ms-2 dark:text-gray-400 dark:hover:text-white",
    spanClass = "ms-1 text-sm font-medium text-gray-500 md:ms-2 dark:text-gray-400",
  } = $props();

  // Track active URL
  activeUrlStore.set(activeUrl);
  appUrlStore.set(appUrl);

  const baseUrl = writable("");
  const params = writable(new URLSearchParams(window.location.search).toString());
  const queryString = derived(params, ($params) => ($params ? "?" + $params : ""));
  const activeUrlArray = writable([]);
  const copyActiveUrl = writable("");

  import { call } from "$lib";

  const postProcessName = (name) => {
    let newName = name.replace(/-/g, " ").replace(/_/g, " ");

    if (name === "crm") {
      return name.toUpperCase();
    } else if (name.split(" ").length > 1) {
      return newName
        .split(" ")
        .map((word) => word.charAt(0).toUpperCase() + word.slice(1))
        .join(" ");
    } else {
      return newName.charAt(0).toUpperCase() + newName.slice(1);
    }
  };

  onMount(async () => {
    const result = await $call.getPage("client-dashboard", "home-page");
    baseUrl.set(result.slice(0, result.lastIndexOf("/")));
  });

  $effect(() => {
    activeUrlStore.subscribe(($activeUrl) => {
      if ($activeUrl !== $copyActiveUrl) {
        console.log($activeUrl);
        const cleanAppUrl = appUrl.replace(/^\/|\/$/g, "");
        copyActiveUrl.set($activeUrl);
        let cleanActiveUrl = $activeUrl;
        let activeUrlModified = cleanActiveUrl.replace(new RegExp(`^/${cleanAppUrl}`), "");
        activeUrlArray.set(
          activeUrlModified
            .split("/")
            .filter((url) => url !== "client-dashboard")
            .filter(Boolean)
            .map((url, index, array) => {
              let accumulatedPath = "";
              for (let i = 0; i <= index; i++) {
                accumulatedPath += `/${array[i]}`;
              }
              let href = accumulatedPath;
              if (index === array.length - 1) {
                href += $queryString;
              }
              return { name: url, href: href };
            })
        );
      }

      if ($breadcrumbLastName) {
        // change the last name of the breadcrumb with the one from the store
        // console.log($breadcrumbLastName);
        $activeUrlArray.update((array) => {
          array[array.length - 1].name = $breadcrumbLastName;
          return array;
        });
      }
    });
  });
</script>

<Breadcrumb aria-label="Default breadcrumb">
  <li class="inline-flex items-center">
    <button
      type="button"
      class={linkClass}
      onclick={() => {
        goto($baseUrl);
      }}
    >
      <svelte:component
        this={getIcons("home")}
        svgClass="w-5 h-5 text-gray-700 dark:text-gray-400"
      />
    </button>
  </li>

  {#if $activeUrlStore !== "/"}
    {#each $activeUrlArray as item, index (item.href)}
      <li class="inline-flex items-center">
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
            class={index === $activeUrlArray.length - 1 ? `${linkClass} !font-bold` : linkClass}
            onclick={() => {
              goto(`${$baseUrl}${item.href}`);
            }}
          >
            {$_(toCamelCase(postProcessName(item.name)))}
          </button>
        {:else}
          <span class={spanClass}>
            {$_(toCamelCase(postProcessName(item.name)))}
          </span>
        {/if}
      </li>
    {/each}
  {/if}
</Breadcrumb> -->
<script>
  import { Breadcrumb, BreadcrumbItem } from "flowbite-svelte";
  import { twMerge } from "tailwind-merge";
  export let linkClass =
    "ms-1 text-sm font-medium text-gray-700 hover:text-gray-900 md:ms-2 dark:text-gray-400 dark:hover:text-white";
  export let spanClass = "ms-1 text-sm font-medium text-gray-500 md:ms-2 dark:text-gray-400";

  export let activeUrl = "/";
  let accumulatedPath = "";
  let params = new URLSearchParams(window.location.search);
  let queryString = params.toString() ? "?" + params.toString() : "";
  $: activeUrlArray = activeUrl
    .split("/")
    .filter((url) => url !== "")
    // .filter((url) => url !== "" && url !== "libraries") //bakrie
    .map((url, index, array) => {
      console.log("=================", url);
      accumulatedPath += `/${url}`;
      let href = accumulatedPath;
      if (index === array.length - 1) {
        href += queryString;
      }
      return { name: url, href: href };
    });
</script>

<Breadcrumb aria-label="Default breadcrumb">
  <BreadcrumbItem href="/" home>Home</BreadcrumbItem>
  {#if activeUrl !== "/"}
    {#each activeUrlArray as item}
      <li class={twMerge("inline-flex items-center", $$props.class)}>
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
          <a
            class={twMerge(linkClass, $$props.classLink)}
            tabindex="0"
            role="button"
            on:click={() => {
              const params = new URLSearchParams(window.location.search);
              const queryParams = {};
              for (const [key, value] of params.entries()) {
                queryParams[key] = value;
              }
              console.log(queryParams);
              // window.location = item.href;
            }}
          >
            {item.name}
          </a>
        {:else}
          <span class={twMerge(spanClass, $$props.classSpan)}>
            {item.name}
          </span>
        {/if}
      </li>
    {/each}
  {/if}
</Breadcrumb>
