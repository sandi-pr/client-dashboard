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
  import { Chart } from "flowbite-svelte";

  export let dataChart;

  let options = {};
  let ready = false;

  if (dataChart) {
    options = dataChart;
    ready = true;
  }
</script>

{#if ready}
  <Chart {options} />
{/if}
