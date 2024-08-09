<!--
  Copyright (C) 2024 - Hexacode TI
  All Rights Reserved.

  Author(s):
  * Sandi Permana <sandy@hexacode.org>

  This file is part of the "Crossroads" project.
-->

<script>
  import { Textarea } from "flowbite-svelte";

  let { onEdit = false, value = $bindable(), placeholder = "", rows = 4 } = $props();
</script>

<Textarea
  readonly={!onEdit}
  bind:value={value}
  placeholder={placeholder}
  rows={rows}
  class={`bg-primary-light dark:bg-secondary-dark ${!onEdit ? "cursor-default focus:ring-0 focus:!border-gray-200 dark:focus:!border-gray-600" : "focus:border-primary-500 focus:ring-primary-500 dark:focus:border-primary-500 dark:focus:ring-primary-500"}`}
/>
