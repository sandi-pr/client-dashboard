<script>
  import { _ } from "svelte-i18n";
  import { twMerge } from "tailwind-merge";
  import { Typography, Badge } from "$lib";
  import { getIcons } from "$lib";

  let { icon, title, description, onclick, comingSoon } = $props();

  let PropsIcon = getIcons(icon);
  let ArrowRightIcon = getIcons("arrowRight");

  let cardClass = twMerge(
    "p-6",
    "border rounded-xl border-primary-light dark:border-primary-dark",
    "bg-primary-alt-light dark:bg-primary-alt-dark hover:bg-primary-hover-light dark:hover:bg-primary-hover-dark active:bg-secondary-light dark:active:bg-secondary-dark",
  );
</script>

<button {onclick} class={cardClass}>
  <div class="flex flex-col gap-2">
    {#if comingSoon}
      <Badge type="brand" text="Soon" />
    {:else}
      <div class="py-[11px]"></div>
    {/if}
    <Typography preset="Text lg/Semibold" text={title} color="primary" />
    <Typography preset="Text sm/Regular" text={description} color="secondary" />
    <div class="flex justify-between items-end">
      <div class="flex items-center gap-2">
        <Typography preset="Text sm/Semibold" text={$_("tryNow")} color="secondary" />
        <ArrowRightIcon svgClass="w-5 h-5 text-secondary-light dark:text-secondary-dark" />
      </div>
      <PropsIcon svgClass="w-[7.5rem] h-[7.5rem]" />
    </div>
  </div>
</button>

<!--
@component
## Config
@prop href: string; // URL to navigate when the card is clicked.
@prop leftIcon: string; // Path to the left icon of the card.
@prop icon: string; // Path to the icon of the card.
@prop title: string; // Title of the card.
@prop description: string; // Description of the card.
@prop button: { name: string, onclick?: () => void }; // Button configuration for the card.
-->
