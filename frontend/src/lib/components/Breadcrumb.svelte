<script>
  import { Breadcrumb, BreadcrumbItem } from "flowbite-svelte";

  export let activeUrl = "/";
  let accumulatedPath = "";
  let params = new URLSearchParams(window.location.search);
  let queryString = params.toString() ? "?" + params.toString() : "";
  $: activeUrlArray = activeUrl
    .split("/")
    .filter((url) => url !== "")
    .map((url) => {
      accumulatedPath += `/${url}`;
      return { name: url, href: accumulatedPath + queryString };
    });
</script>

<Breadcrumb aria-label="Default breadcrumb">
  <BreadcrumbItem href="/" home>Home</BreadcrumbItem>
  {#if activeUrl !== "/"}
    {#each activeUrlArray as item}
      <BreadcrumbItem href={item.href} active={item.href === activeUrl}>{item.name}</BreadcrumbItem>
    {/each}
  {/if}
</Breadcrumb>
