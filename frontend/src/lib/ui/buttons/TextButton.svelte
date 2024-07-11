<script>
  import { A } from "flowbite-svelte";

  let { href, text } = $props();
</script>

<A
  {href}
  class="text-sm font-semibold text-fg-brand-primary-light dark:text-fg-brand-primary-dark font-jakarta leading-tight"
  >{text}</A
>
