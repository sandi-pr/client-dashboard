<script lang="ts">
  import { getIcons, Typography, FeaturedIcon } from "$lib";
  import { Card, Button } from "flowbite-svelte";

  let {
    type,
    href,
    icon,
    image,
    alt,
    title,
    description,
    subtext,
    onclick
  }: {
    type?: string;
    href?: string;
    icon?: string;
    image?: string;
    alt?: string;
    title: string;
    description: string;
    subtext?: string;
    onclick?: () => void;
  } = $props();
</script>

<button>
  <Card {href} on:click={onclick} horizontal={true} class={`md:items-center w-full ${type === "ledger" ? "!p-6 border-secondary-light dark:border-secondary-dark" : "min-h-28 !p-4 border-brand-light dark:border-primary-dark gap-3"} bg-primary-light dark:bg-primary-dark hover:bg-primary-hover-light dark:hover:bg-primary-hover-dark shadow-sm hover:no-underline`} size="none">
    {#if icon}
      <FeaturedIcon classesSvg={getIcons(icon)} />
    {/if}
    {#if type === "ledger"}
      <div class="gap-2">
        <div class="gap-1">
          <Typography preset="Text sm/Semibold" text={title} color="primary" />
          <Typography preset="Text sm/Regular" text={description} color="tertiary" />
        </div>
        <Typography preset="Text sm/Medium" text={subtext} color="tertiary" />
      </div>
    {:else}
      <div class="gap-[2px]">
        <Typography preset="Text md/Semibold" text={title} color="primary" />
        <Typography preset="Text sm/Regular" text={description} color="secondary" />
      </div>
    {/if}
  </Card>
</button>

<!--
@component
## Config
@prop href: string; // URL to navigate when the card is clicked.
@prop leftIcon: string; // Path to the left icon of the card.
@prop icon: string; // Path to the icon of the card.
@prop image: string; // Path to the image of the card.
@prop alt: string; // Alternative text for the image or icon of the card.
@prop title: string; // Title of the card.
@prop description: string; // Description of the card.
@prop button: { name: string, onclick?: () => void }; // Button configuration for the card.
-->
