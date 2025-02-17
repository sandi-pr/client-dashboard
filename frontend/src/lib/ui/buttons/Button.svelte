<script lang="ts">
  import { twMerge } from "tailwind-merge";
  import { Button } from "flowbite-svelte";
  import { getIcons } from "$lib";

  // Ekstrak properti langsung dari $props()
  let {
    color,
    outline,
    width,
    spaceY,
    text,
    position,
    icon,
    classesSvg,
    leftIcon,
    rightIcon,
    alt,
    svgClassReplace,
    funcButton,
  }: {
    color?:
      | "primary"
      | "alternative"
      | "blue"
      | "dark"
      | "green"
      | "light"
      | "purple"
      | "red"
      | "yellow"
      | "none";
    outline?: boolean;
    width?: string;
    spaceY?: string;
    text?: string;
    position?: string;
    icon?: string;
    classesSvg?: string;
    leftIcon?: string;
    rightIcon?: string;
    alt?: string;
    svgClassReplace?: string;
    funcButton?: () => void;
  } = $props();

  let svgClass = $state(
    svgClassReplace ? svgClassReplace : "w-5 h-5 text-gray-700 dark:text-gray-400",
  );
</script>

<div class={`flex ${position ?? "justify-center"}`}>
  <Button
    on:click={funcButton}
    color={color ?? "primary"}
    {outline}
    class={`${width ?? "w-max"} ${spaceY ?? "mt-4"} ${text ? "" : "!px-[10px]"}`}
  >
    {#if leftIcon || icon}
      {#if classesSvg}
        <svelte:component this={getIcons(classesSvg)} {svgClass} />
      {:else}
        <img src={leftIcon ?? icon} alt={alt ?? ""} class={text ? "me-2" : ""} />
      {/if}
    {/if}
    {text ?? ""}
    {#if rightIcon}
      {#if classesSvg}
        <svelte:component this={getIcons(classesSvg)} {svgClass} />
      {:else}
        <img src={rightIcon ?? icon} alt={alt ?? ""} class={text ? "me-2" : ""} />
      {/if}
    {/if}
    <slot />
  </Button>
</div>

<!--
@component
## Config
@prop color: 'primary' | 'alternative' | 'blue' | 'dark' | 'green' | 'light' | 'primary' | 'purple' | 'red' | 'yellow' | 'none'; // Color of the button.
@prop outline: false | true; // Whether the button should have an outline.
@prop width: 'w-max'; // Width of the button.
@prop spaceY: 'mt-4'; // Vertical spacing of the button.
@prop text: string; // Text content of the button.
@prop position: 'justify-center'; // Position of the button within its container.
@prop icon: string; // Path to the button icon.
@prop leftIcon: string; // Path to the left icon of the button.
@prop rightIcon: string; // Path to the right icon of the button.
@prop alt: string; // Alternative text for the button icon.
-->
