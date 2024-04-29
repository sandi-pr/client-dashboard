<script lang="ts">
  import { Badge } from 'flowbite-svelte';

  let { color, large, dismissable, onclick, closeIcon }: {
    color?: "none" | "red" | "yellow" | "green" | "indigo" | "purple" | "pink" | "blue" | "dark" | "primary" | undefined,
    large?: boolean,
    dismissable?: boolean,
    onclick?: () => void,
    closeIcon?: string
  } = $props();
</script>

<Badge color={color} large={large} dismissable={dismissable}>
  Default
  {#if dismissable}
    <!-- <button slot="close-button" let:close on:click={onclick} type="button" class="inline-flex items-center rounded-full p-0.5 my-0.5 ms-1.5 -me-1.5 text-sm bg-primary-500 dark:bg-primary-400 text-white dark:text-primary-800 hover:bg-primary-900 hover:text-white dark:hover:bg-red-900 dark:hover:text-yellow-300" aria-label="Remove">
      <img src={closeIcon} class="h-3.5 w-3.5" alt="Close" />
    </button> -->
  {/if}
</Badge>

<!--
  @component
  ## Config
  @prop color: "none" | "red" | "yellow" | "green" | "indigo" | "purple" | "pink" | "blue" | "dark" | "primary" | undefined; // Color of the badge.
  @prop large: boolean; // Whether the badge should be large.
  @prop dismissable: boolean; // Whether the badge should be dismissable.
  @prop onclick: () => void; // Click event handler for the badge.
  @prop closeIcon: string; // Path to the close icon for dismissable badge.
-->
