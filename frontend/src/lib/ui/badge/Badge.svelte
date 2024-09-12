<script lang="ts">
  import { Badge } from 'flowbite-svelte';

  let { text, type = "brand" }: {
    text: string;
    type?: string;
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
</script>

<p
class={`py-[2px] px-2 border-[1px] text-xs font-medium rounded-xl text-center w-max h-max ${classStatusType(type)}`}
>
  {text}
</p>

<!--
  @component
  ## Config
  @prop type: "success" | "error" | "warning" | "info" | "brand"; // Color of the badge.
  @prop text: string; // Text of the badge.
-->
