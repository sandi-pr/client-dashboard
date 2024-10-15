<script lang="ts">
  import { Badge } from 'flowbite-svelte';

  let { text, type = "brand", icon = null, iconClass = null }: {
    text: string;
    type?: string;
    icon?: any;
    iconClass?: string;
  } = $props();

  const classStatusType = (type) => {
    if (type === "success") {
      return "bg-utility-success-50-light dark:bg-utility-success-50-dark border-utility-success-200-light dark:border-utility-success-200-dark text-utility-success-700-light dark:text-utility-success-700-dark";
    } else if (type === "error") {
      return "bg-utility-error-50-light dark:bg-utility-error-50-dark border-utility-error-200-light dark:border-utility-error-200-dark text-utility-error-700-light dark:text-utility-error-700-dark";
    } else if (type === "warning") {
      return "bg-utility-warning-50-light dark:bg-utility-warning-50-dark border-utility-warning-200-light dark:border-utility-warning-200-dark text-utility-warning-700-light dark:text-utility-warning-700-dark";
    } else if (type === "info") {
      return "bg-utility-gray-50-light dark:bg-utility-gray-50-dark border-utility-gray-200-light dark:border-utility-gray-200-dark text-utility-gray-700-light dark:text-utility-gray-700-dark";
    } else if (type === "brand") {
      return "bg-utility-brand-50-light dark:bg-utility-brand-50-dark border-utility-brand-200-light dark:border-utility-brand-200-dark text-utility-brand-700-light dark:text-utility-brand-700-dark";
    } else if (type === "sky") {
      return "text-sky-700 dark:text-sky-300 bg-sky-50 dark:bg-sky-950 border-sky-200 dark:border-sky-800";
    }
  };

  const getIconClass = (type) => {
    if (type === "success") {
      return "text-utility-success-500-light dark:text-utility-success-500-dark";
    } else if (type === "error") {
      return "text-utility-error-500-light dark:text-utility-error-500-dark";
    } else if (type === "warning") {
      return "text-utility-warning-500-light dark:text-utility-warning-500-dark";
    } else if (type === "info") {
      return "text-utility-gray-500-light dark:text-utility-gray-500-dark";
    } else {
      return "text-utility-brand-500-light dark:text-utility-brand-500-dark";
    }
  }
</script>

<p
class={`flex items-center gap-1 py-[2px] px-2 border-[1px] text-xs font-medium rounded-xl text-center w-max ${classStatusType(type)}`}
>
  {#if icon}
    <!-- svelte-ignore svelte_component_deprecated -->
    <svelte:component this={icon} svgClass={getIconClass(type) + " w-3 h-3"} />
  {/if}
  {text}
</p>

<!--
  @component
  ## Config
  @prop type: "success" | "error" | "warning" | "info" | "brand"; // Color of the badge.
  @prop text: string; // Text of the badge.
-->
