<!-- 
  CR3 Client Dashboard
  Package: gitlab.hexacode.org/public-repo/js/hexacomponents-svelte
  Maintainer:
  * Sandi Permana <sandy@hexacode.org>
  * Taqie Fadlillah <taqief@hexacode.org>
  * Dean Ramhan <dean@hexacode.org>
  
  Copyright (C) 2024 Hexacode Teknologi Indonesia
  All Rights Reserved.
-->
<script>
  import { Typography } from "$lib";
  let { brandName = '', href = '/', brandLogo = '', imgClass = 'h-6 me-2 sm:h-8' } = $props();

  let site = {
    name: brandName,
    href: href,
    img: brandLogo,
  };

  let aClass = 'flex items-center mb-5 w-max';

</script>

<a href={site.href} class={aClass}>
  <img src={site.img} class={imgClass} alt={site.name} />
  <Typography preset="Text lg/Bold" text={site.name} color="primary" />
</a>


