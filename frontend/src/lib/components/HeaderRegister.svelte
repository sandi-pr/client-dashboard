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
  import { DarkMode } from "flowbite-svelte";
  import logo from "../icon/Logomark.svg";
  import { _, locale } from "svelte-i18n";
  import { LogoBrand, LanguageSelector } from "hexacomponents-svelte";

  const changeLanguage = (lang) => {
    locale.set(lang);
  };
</script>

<div class="flex justify-between items-center w-full">
  <div class="flex items-center">
    <LogoBrand brandLogo={logo} brandName="Crossroad" href="/login" />
  </div>

  <div class="flex gap-4 items-center">
    <LanguageSelector {changeLanguage} />
    <DarkMode size="lg" />
  </div>
</div>
