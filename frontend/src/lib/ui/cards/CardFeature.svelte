<script>
  import { twMerge } from "tailwind-merge";
  import { Typography, FeaturedIcon, Button } from "$lib";
  import { getIcons } from "$lib";

  let { cardData = {}, propClass = "" } = $props();

  let baseStyle = [
    "flex",
    "flex-col",
    "gap-4",
    "w-full",
    "border",
    "border-primary-light",
    "dark:border-primary-dark",
    "rounded-lg",
    "items-center",
    "justify-center",
    "bg-primary-light",
    "dark:bg-primary-dark",
  ].join(" ");

  let contentCenterClasses = [
    "p-4",
    "flex",
    "flex-col",
    "w-full",
    "items-center",
    "justify-center",
    "gap-2 md:gap-4",
  ].join(" ");

  let classes = twMerge(baseStyle, propClass);
</script>

<div class={classes} onclick={cardData?.onClick}>
  <div class={contentCenterClasses}>
    {#if cardData?.iconData}
      <!-- <Icon {...cardData.iconData} /> -->
      <FeaturedIcon classesSvg={getIcons(cardData.iconData.icon)} />
    {/if}

    {#if cardData?.contentCenter}
      {#each cardData.contentCenter as contentArray}
        {#each contentArray as content}
          <div
            class={`flex w-full gap-4 items-center ${
              "justify" in content ? content.justify : "justify-center"
            }`}
          >
            {#each Object.keys(content) as key}
              {#if key === "icon"}
                <FeaturedIcon classesSvg={getIcons(content.icon.icon)} />
              {/if}

              {#if key === "text"}
                <Typography {...content.text} />
              {/if}
            {/each}
          </div>
        {/each}
      {/each}
    {/if}
  </div>
  {#if cardData?.actionFooter}
    <div
      class="flex flex-col gap-2 justify-center w-full bg-secondary-light dark:bg-secondary-dark rounded-b-lg border-t border-primary-light dark:border-primary-dark p-8"
    >
      {#each cardData.actionFooter as action}
        <Button {...action.button} />
      {/each}
    </div>
  {/if}
</div>
