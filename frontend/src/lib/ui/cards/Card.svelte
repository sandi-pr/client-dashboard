<script lang="ts">
  import { getIcons, Typography, FeaturedIcon, Badge, Button } from "$lib";
  import { _ } from "svelte-i18n";
  import { Card } from "flowbite-svelte";

  let {
    type,
    href,
    icon,
    iconColor,
    image,
    alt,
    title,
    description,
    subtext,
    changeType,
    onclick,
    comingSoon,
  }: {
    type?: string;
    href?: string;
    icon?: string;
    iconColor?: string;
    image?: string;
    alt?: string;
    title: string;
    description: string;
    subtext?: string;
    changeType?: string;
    onclick?: () => void;
    comingSoon?: boolean;
  } = $props();
</script>

<button>
  <Card
    {href}
    on:click={onclick}
    horizontal={type === "alt" || type === "report" || type === "reportNoButton" ? false : true}
    class={`${type === "alt" ? "!p-0 min-h-[270px]" : type === "report" || type === "reportNoButton" ? "!p-0" : "!p-4 md:items-center min-h-28"} !gap-1 md:!gap-3 w-full ${type === "ledger" ? " border-secondary-light dark:border-secondary-dark" : "border-brand-light dark:border-primary-dark gap-3"} bg-primary-light dark:bg-primary-dark hover:bg-primary-hover-light dark:hover:bg-primary-hover-dark shadow-sm hover:no-underline`}
    size="none"
  >
    <div
      class={type === "alt" || type === "report" || type === "reportNoButton"
        ? "flex flex-col grow gap-2 md:gap-5 px-4 pt-4 md:px-6 md:pt-6"
        : "flex items-center !gap-1 md:!gap-3"}
    >
      {#if icon}
        <div class="flex gap-4 items-center">
          <FeaturedIcon
            classesSvg={getIcons(icon)}
            color={type === "alt"
              ? "alt"
              : type === "report" || type === "reportNoButton"
                ? "primary"
                : "brand"}
            {iconColor}
          />
          {#if type === "reportNoButton"}
            <Typography preset="Text md/Regular" text={description} color="tertiary" />
          {/if}
        </div>
      {/if}
      {#if type === "ledger"}
        <div class="gap-2">
          <div class="gap-1">
            <Typography preset="Text sm/Semibold" text={title} color="primary" />
            <Typography preset="Text sm/Regular" text={description} color="tertiary" />
          </div>
          <Typography preset="Text sm/Medium" text={subtext} color="tertiary" />
        </div>
      {:else if type === "report" || type === "reportNoButton"}
        <div class="gap-[2px]">
          {#if type === "report"}
            <Typography
              preset="Text sm/Medium"
              text={description}
              color="secondary"
              propClass="!leading-tight !pb-2"
            />
          {/if}
          <div class="flex justify-between items-center {type === 'reportNoButton' ? 'mb-4' : ''}">
            <Typography preset="Display sm/Semibold" text={title} color="primary" />
          </div>
          {#if type === "report"}
            <div class="flex gap-1 justify-end">
              <!-- svelte-ignore svelte_component_deprecated -->
              <svelte:component
                this={getIcons(changeType === "decrease" ? "trendDown" : "trendUp")}
              />
              <Typography
                preset="Text sm/Medium"
                text={subtext}
                color={changeType === "decrease" ? "error-primary" : "success-primary"}
                propClass="pe-1"
              />
              <Typography preset="Text sm/Medium" text={$_("vsLastProject")} color="tertiary" />
            </div>
          {/if}
        </div>
      {:else}
        <div class={type === "alt" ? "flex flex-col gap-2" : "gap-[2px]"}>
          <Typography
            preset={`${type === "alt" ? "Display sm/Semibold" : "Text md/Semibold"}`}
            text={title}
            color="primary"
          />
          <Typography
            preset={`${type === "alt" ? "Text sm/Medium" : "Text sm/Regular"}`}
            text={description}
            color="secondary"
            propClass="!leading-tight !pb-2"
          />
          {#if comingSoon}
            <Badge type="brand" text="Coming soon" />
          {/if}
        </div>
      {/if}
    </div>
    {#if type === "alt" || type === "report"}
      <div
        class={`flex ${type === "report" ? "items-end pr-4" : "pl-4"} pb-4 flex-col items-start border-t border-fg-secondary-light dark:border-fg-secondary-dark gap-2`}
      >
        <Button funcButton={onclick} color="light">
          {type === "report" ? $_("viewReport") : $_("openHere")}
        </Button>
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
