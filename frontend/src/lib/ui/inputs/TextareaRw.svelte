<!--
  Copyright (C) 2024 - Hexacode TI
  All Rights Reserved.

  Author(s):
  * Sandi Permana <sandy@hexacode.org>

  This file is part of the "Crossroads" project.
-->

<script>
  import { Textarea } from "flowbite-svelte";
  import { twMerge } from "tailwind-merge";

  let { onEdit = false, value = $bindable(), placeholder = "", rows = 4 } = $props();

  const inputClass = () => {
    const disabledClass = onEdit
      ? ""
      : "cursor-default focus:ring-0 border-disabled-subtle-light dark:border-disabled-subtle-dark focus:border-disabled-subtle-light focus:dark:border-disabled-subtle-dark bg-disabled-subtle-light dark:bg-disabled-subtle-dark text-disabled-light dark:text-disabled-dark";

    return twMerge(
      "text-primary-light dark:text-primary-dark focus:ring-0 border-primary-light dark:border-primary-dark focus:border-primary-light focus:dark:border-primary-dark bg-primary-light dark:bg-primary-dark",
      disabledClass,
    );
  };
</script>

<Textarea readonly={!onEdit} bind:value {placeholder} {rows} class={inputClass()} />
