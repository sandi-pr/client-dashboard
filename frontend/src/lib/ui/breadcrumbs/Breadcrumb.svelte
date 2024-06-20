<script>
  import { Breadcrumb, BreadcrumbItem } from "flowbite-svelte";
  import { goto } from "$app/navigation";
  import { onMount } from "svelte";
  import { CR3Call } from "hexacomponents-svelte";
  import { getIcons } from "$lib";
  import { toCamelCase } from "$lib/utils/toCamelCase.js";
  import { _ } from "svelte-i18n";

  let {
    activeUrl = "/",
    appUrl = "/",
    linkClass = "ms-1 text-sm font-medium text-gray-700 hover:text-gray-900 md:ms-2 dark:text-gray-400 dark:hover:text-white",
    spanClass = "ms-1 text-sm font-medium text-gray-500 md:ms-2 dark:text-gray-400",
  } = $props();
  let baseUrl = $state("");
  let params = $state(new URLSearchParams(window.location.search).toString());
  let queryString = $state(params.toString() ? "?" + params.toString() : "");
  let activeUrlArray = $state([]);
  let copyActiveUrl = $state("");
  let activeUrlModified = $state("");

  const call = new CR3Call(globalThis.location.origin);

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

  onMount(async () => {
    const result = await call.getPage("client-dashboard", "home-page");
    baseUrl = result.slice(0, result.lastIndexOf("/"));
  });

  $effect(() => {
    if (activeUrl !== copyActiveUrl) {
      console.log(activeUrl);
      const cleanAppUrl = appUrl.replace(/^\/|\/$/g, "");
      copyActiveUrl = activeUrl;
      let cleanActiveUrl = activeUrl;
      activeUrlModified = cleanActiveUrl.replace(
        new RegExp(`^/${cleanAppUrl}`),
        "",
      );
      activeUrlArray = activeUrlModified
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
            href += queryString;
          }
          return { name: url, href: href };
        });
    }
  });
</script>

<Breadcrumb aria-label="Default breadcrumb">
  <!-- <BreadcrumbItem href="/" home></BreadcrumbItem> -->
  <li class={"inline-flex items-center"}>
    <button
      type="button"
      class={linkClass}
      onclick={() => {
        // window.location = `${appUrl}${item.href}`;
        // goto(`${appUrl}${item.href}`);
        // console.log(item.href);
        // globalThis.location.href = "/client-dashboard" + item.href;
        goto(`${baseUrl}`);
      }}
      ><svelte:component
        this={getIcons("home")}
        svgClass="w-5 h-5 text-gray-700 dark:text-gray-400"
      />
    </button>
  </li>
  {#if activeUrl !== "/"}
    {#each activeUrlArray as item, index (item.href)}
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
            class={index === activeUrlArray.length - 1
              ? `${linkClass} !font-bold`
              : linkClass}
            onclick={() => {
              // window.location = `${appUrl}${item.href}`;
              // goto(`${appUrl}${item.href}`);
              console.log(item.href);
              // globalThis.location.href = "/client-dashboard" + item.href;
              goto(`${baseUrl}${item.href}`);
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
</Breadcrumb>

<!--
@component
## Config
@prop activeUrl: string; // current active url
@prop linkClass: string; // class for the link
@prop spanClass: string; // class for the span
-->
