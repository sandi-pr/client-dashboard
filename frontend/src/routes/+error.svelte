<script>
  import { onMount } from "svelte";
  import { Typography, Button } from "$lib";
  import { page } from "$app/stores";

  import Logo404 from "$lib/image/404.svg";
  import Logo502 from "$lib/image/502.png";

  const errorList = {
    404: "Oops! Page Not Found",
    500: "Internal server error",
  };

  const errorSubtitle = {
    404: "Uh oh! This page has gone on a little adventure. Try searching for what you need or head back to the homepage",
    500: "Something went wrong on our end. Please try again later.",
  };

  let status = $page.status;
  let referrer = "";

  onMount(() => {
    referrer = $page.url.href !== document.referrer ? document.referrer : "";
    // console.log("referrer", document.referrer);
  });
  console.log("page", $page);
</script>

<svelte:head>
  <title>{status} - {errorList[status]}</title>
</svelte:head>

<main
  class="flex flex-col items-center justify-center bg-primary-light dark:bg-primary-dark w-full h-screen"
>
  {#if status === 404}
    <img src={Logo404} alt="404" class="w-[40rem] h-[40rem]" />
  {:else}
    <img src={Logo502} alt="502" class="w-[40rem] h-[40rem]" />
  {/if}
  <Typography text={errorList[status]} preset="Display xl/semibold" propClass="m-8" />
  {#if status === 404}
    <Typography text={errorSubtitle[status]} preset="Text md/regular" />
  {/if}
  <div class="flex gap-4 mt-8">
    {#if status === 404}
      <Button
        text="Go back"
        funcButton={() => (window.location.href = referrer || "/")}
        color="alternative"
        leftIcon="true"
        classesSvg="arrowLeft"
      />
      <Button text="Take me home" funcButton={() => (window.location.href = "/")} />
    {:else}
      <Button
        text="Refresh"
        funcButton={() => (window.location.href = "/")}
        color="alternative"
        leftIcon="true"
        classesSvg="refresh"
      />
    {/if}
  </div>
</main>
