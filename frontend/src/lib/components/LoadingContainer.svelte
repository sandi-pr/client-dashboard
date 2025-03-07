<!--
  Copyright (C) 2024 - Hexacode TI
  All Rights Reserved.

  Author(s):
  * Dean Ramhan <dean@hexacode.org>
  * Taqie Fadlillah <taqief@hexacode.org>
  * Sandi Permana <sandy@hexacode.org>

  This file is part of the "Crossroads" project.
-->

<script>
  import { Spinner } from "flowbite-svelte";
  import { toastLoading } from '$lib/utils/loading';
  import { _ } from 'svelte-i18n';
</script>

{#if $toastLoading}
<div
  class="absolute bottom-0 left-0 right-0 top-0 z-50 flex items-center justify-center bg-white dark:bg-black backdrop-blur-sm bg-opacity-40"
>
  <div class="rounded-lg p-8">
    <div class="flex justify-center">
      <Spinner class="h-14 w-14" />
    </div>
    <div class="mt-4 text-center">
      <h2
        class="bg-gradient-to-tr from-black to-gray-500 bg-clip-text text-lg font-semibold text-transparent"
      >
        {$_("loading")}
      </h2>
      <p class=" text-gray-600">{$_("thisMayTakeAFewSecondsPleaseWait")}</p>
    </div>
  </div>
</div>
{/if}
