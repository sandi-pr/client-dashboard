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
  import { twMerge } from "tailwind-merge";
  export let className = "";

  let defaultClass = "w-full border-b border-hc-graylight-200 dark:border-hc-graydark-800";

  let classes = twMerge(defaultClass, className);
</script>

<div class={classes}></div>
