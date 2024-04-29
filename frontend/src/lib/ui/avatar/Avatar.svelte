<script lang="ts">
  import { Avatar } from 'flowbite-svelte';

  let { src, size, color, border, rounded, onclick }: {
    src?: string,
    size?: "none" | "xs" | "sm" | "lg" | "xl" | "md",
    color?: string,
    border?: boolean,
    rounded?: boolean,
    onclick?: () => void
  } = $props();
</script>

<Avatar src={src} size={size} color={color} border={border} rounded={rounded} onclick={onclick} />

<!--
  @component
  ## Config
  @prop src: string; // Path to the avatar image.
  @prop size: 'none' | 'xs' | 'sm' | 'lg' | 'xl' | 'md'; // Size of the avatar.
  @prop color: 'primary' | 'alternative' | 'blue' | 'dark' | 'green' | 'light' | 'primary' | 'purple' | 'red' | 'yellow' | 'none'; // Color of the avatar.
  @prop border: string | boolean; // Border style of the avatar.
  @prop rounded: boolean; // Whether the avatar should have rounded corners.
  @prop onclick: () => void; // Click event handler for the avatar.
-->