<script lang="ts">
  import { Breadcrumb, BreadcrumbItem } from "flowbite-svelte";
  // import { CR3Call } from "hexacomponents-svelte";
  import { goto } from "$app/navigation";

  // const call = new CR3Call(globalThis.location.origin);

  // const templates = [
  //   {
  //     name: "home-page",
  //     path: "/",
  //     file_path: "template/index.html",
  //   },
  //   {
  //     name: "crm-overview",
  //     path: "/crm",
  //     file_path: "template/crm.html",
  //   },
  //   {
  //     name: "crm-person",
  //     path: "/crm/person",
  //     file_path: "template/crm/person.html",
  //   },
  //   {
  //     name: "crm-person-create",
  //     path: "/crm/person/create",
  //     file_path: "template/crm/person/create.html",
  //   },
  //   {
  //     name: "crm-person-detail",
  //     path: "/crm/person/single",
  //     file_path: "template/crm/person/detail.html",
  //   },
  //   {
  //     name: "crm-person-list",
  //     path: "/crm/person/list",
  //     file_path: "template/crm/person/list.html",
  //   },
  //   {
  //     name: "crm-product-service",
  //     path: "/crm/product_service",
  //     file_path: "template/crm/product_service.html",
  //   },
  //   {
  //     name: "crm-product-service-create",
  //     path: "/crm/product_service/create",
  //     file_path: "template/crm/product_service/create.html",
  //   },
  //   {
  //     name: "crm-product-service-detail",
  //     path: "/crm/product_service/detail",
  //     file_path: "template/crm/product_service/detail.html",
  //   },
  //   {
  //     name: "crm-users",
  //     path: "/crm/users",
  //     file_path: "template/crm/users.html",
  //   },
  //   {
  //     name: "crm-user-detail",
  //     path: "/crm/users/detail",
  //     file_path: "template/crm/users/detail.html",
  //   },
  //   {
  //     name: "crm-user-create",
  //     path: "/crm/users/create",
  //     file_path: "template/crm/users/create.html",
  //   },
  //   {
  //     name: "crm-categories",
  //     path: "/crm/categories",
  //     file_path: "template/crm/categories.html",
  //   },
  //   {
  //     name: "crm-category-create",
  //     path: "/crm/categories/create",
  //     file_path: "template/crm/categories/create.html",
  //   },
  //   {
  //     name: "crm-category-detail",
  //     path: "/crm/categories/detail",
  //     file_path: "template/crm/categories/detail.html",
  //   },
  //   {
  //     name: "crm-calendar",
  //     path: "/crm/calendar",
  //     file_path: "template/crm/calendar.html",
  //   },
  //   {
  //     name: "hexabot-overview",
  //     path: "/hexabot",
  //     file_path: "template/hexabot.html",
  //   },
  // ];

  let {
    activeUrl = "/",
    appUrl = "/",
    linkClass = "ms-1 text-sm font-medium text-gray-700 hover:text-gray-900 md:ms-2 dark:text-gray-400 dark:hover:text-white",
    spanClass = "ms-1 text-sm font-medium text-gray-500 md:ms-2 dark:text-gray-400",
  } = $props();
  let accumulatedPath = $state("");
  let params = $state(new URLSearchParams(window.location.search).toString());
  let queryString = $state(params.toString() ? "?" + params.toString() : "");
  let activeUrlArray: any[] = $state([]);
  let copyActiveUrl = $state("");
  let activeUrlModified = $state("");

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

  // onMount(() => {
  //   if (activeUrl) {
  //     const cleanAppUrl = appUrl.replace(/^\/|\/$/g, "");
  //     console.log(cleanAppUrl);
  //     activeUrl = activeUrl.replace(new RegExp(`^/${cleanAppUrl}`), "");
  //     console.log(activeUrl);
  //     activeUrlArray = activeUrl
  //       .split("/")
  //       .filter((url) => url !== "client-dashboard")
  //       .filter(Boolean)
  //       .map((url, index, array) => {
  //         accumulatedPath += `/${url}`;
  //         let href = accumulatedPath;
  //         if (index === array.length - 1) {
  //           href += queryString;
  //         }
  //         return { name: url, href: href };
  //       });
  //   }
  // });

  $effect(() => {
    if (activeUrl !== copyActiveUrl) {
      const cleanAppUrl = appUrl.replace(/^\/|\/$/g, "");
      copyActiveUrl = activeUrl;
      let cleanActiveUrl = activeUrl;
      activeUrlModified = cleanActiveUrl.replace(new RegExp(`^/${cleanAppUrl}`), "");
      activeUrlArray = activeUrlModified
        .split("/")
        .filter((url) => url !== "client-dashboard")
        .filter(Boolean)
        .map((url, index, array) => {
          accumulatedPath += `/${url}`;
          let href = accumulatedPath;
          if (index === array.length - 1) {
            href += queryString;
          }
          return { name: url, href: href };
        });
    }
  })
</script>

<Breadcrumb aria-label="Default breadcrumb">
  <BreadcrumbItem href="/" home></BreadcrumbItem>
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
            class={index === activeUrlArray.length - 1 ? `${linkClass} !font-bold` : linkClass}
            onclick={async () => {
                console.log(item.href);
                console.log(appUrl);
                goto(`${appUrl}${item.href}`);
            }}
          >
            {postProcessName(item.name)}
          </button>
        {:else}
          <span class={spanClass}>
            {postProcessName(item.name)}
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
