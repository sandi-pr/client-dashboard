<script>
  import { _ } from "svelte-i18n";
  let { onEdit, checked = $bindable(false), text } = $props();
</script>

<label
  class={`inline-flex items-center ${!onEdit ? "cursor-default" : "cursor-pointer"}`}
>
  <input
    disabled={!onEdit}
    type="checkbox"
    value=""
    class="hidden peer"
    bind:checked={checked}
  />
  <div
    class="relative w-11 h-6 bg-tertiary-light dark:bg-tertiary-dark peer-focus:outline-none peer-focus:ring-4 peer-focus:ring-utility-brand-200-light dark:peer-focus:ring-utility-brand-200-dark rounded-full peer peer-checked:after:translate-x-full rtl:peer-checked:after:-translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:start-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all dark:border-gray-600 peer-checked:bg-brand-solid-light"
  ></div>
  <span class="ms-3 text-sm font-medium text-gray-900 dark:text-gray-300"
    >{$_(text)}</span
  >
</label>