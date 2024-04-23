/*
 * CR3 Client Dashboard
 * Package: gitlab.hexacode.org/cr3/client-dashboard
 * Maintainer: Azzis Arswendo <azzis@hexacode.org>
 *
 * Copyright (C) 2024 Hexacode Teknologi Indonesia
 * All Rights Reserved
 */

import adapter from "@sveltejs/adapter-static";
import { vitePreprocess } from "@sveltejs/vite-plugin-svelte";
import conf from "../config.json" assert { type: "json" };

/** @type {import('@sveltejs/kit').Config} */
const config_production = {
  kit: {
    appDir: "public",
    paths: {
      base: "/" + conf.namespace,
      relative: false,
    },
    // See https://kit.svelte.dev/docs/adapters for more information about adapters.
    adapter: adapter({
      // default options are shown. On some platforms
      // these options are set automatically — see below
      pages: "template",
      assets: "template",
      fallback: undefined,
      precompress: false,
      strict: true,
    }),
    prerender: {
      crawl: true,
      entries: ["*", "/documents/1", "/person/1", "/workflow/1", "/workflow/update/1"],
    },
  },
  preprocess: vitePreprocess(),
};

/** @type {import('@sveltejs/kit').Config} */
const config_development = {
  kit: {
    // See https://kit.svelte.dev/docs/adapters for more information about adapters.
    adapter: adapter({
      // default options are shown. On some platforms
      // these options are set automatically — see below
      pages: "template",
      assets: "template",
      fallback: undefined,
      precompress: false,
      strict: true,
    }),
  },
  preprocess: vitePreprocess(),
};

var config = process.env.NODE_ENV === "production" ? config_production : config_development;

export default config;
