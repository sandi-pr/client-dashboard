<!--
CR3 Client Dashboard
Package: gitlab.hexacode.org/cr3/client-dashboard
Maintainer: Azzis Arswendo <azzis@hexacode.org>

Copyright (C) 2024 Hexacode Teknologi Indonesia
All Rights Reserved
-->

<script>
  import "../app.css";
  import "@fontsource/plus-jakarta-sans";
  import { page } from "$app/stores";
  import { CR3Call } from "hexacomponents-svelte";
  import { Pane, Splitpanes } from "svelte-splitpanes";
  import { Datepicker } from "flowbite-svelte";
  import {
    Navbar,
    Sidebar,
    ToastContainer,
    setCall,
    callReady,
    triggerShowAI,
    showAI,
    triggerSelectDropdown,
    triggerSelectCollection,
    triggerReload,
    triggerToast,
  } from "$lib";
  import { onMount } from "svelte";
  import { locale, waitLocale, isLoading, _ } from "svelte-i18n";
  import "$lib/service/i18n";
  import SideDrawer from "$lib/components/SideDrawer.svelte";
  import LoadingSpinner from "$lib/components/LoadingSpinner.svelte";
  import LoadingContainer from "$lib/components/LoadingContainer.svelte";
  import { goto } from "$app/navigation";
  import { set_pdfjs_context } from "$lib/utils/vite.ts";
  import { sidebarLargeClosed, toggleSidebarLargeClosed } from "$lib/utils/sidebarLargeClosed";
  import { setPrimaryColor } from "$lib/utils/colorTheme";
  import { purpleColors, greenColors, blueColors } from "$lib/utils/presetColors";

  set_pdfjs_context();

  let innerWidth = $state(0);
  let activeUrl = $state($page.url.pathname);
  let appUrl = "/"; // change to application url to make the breadcrumb not display the application name
  let userData = $state({ username: "", email: "" });
  let authority = $state("");
  let brand = $state("");
  let logoUrl = $state("");
  const isFrontend = process.env.NODE_ENV === "development";
  let hideNav = $state(false);
  let ready = $state(false);

  const call = isFrontend
    ? new CR3Call("http://localhost")
    : new CR3Call(globalThis.location.origin);

  const hideLayoutList = $state(["/workflow/update", "/workflow/create", "/hexabot/webchat"]);

  const hideScrollList = $state(["/document/document_collection/chat"]);

  const hideNavbarList = $state(["/client-dashboard/datasets", "/datasets"]);

  const hideScroll = () => {
    return !hideScrollList.some((path) => {
      try {
        if (path.endsWith("/*")) {
          return globalThis.location.pathname.startsWith(path.slice(0, -1));
        } else {
          return globalThis.location.pathname === path;
        }
      } catch (error) {
        console.log(error);
        if (path.endsWith("/*")) {
          return window.location.pathname.startsWith(path.slice(0, -1));
        } else {
          return window.location.pathname === path;
        }
      }
    });
  };

  const hideLayout = () => {
    return !hideLayoutList.some((path) => {
      try {
        if (path.endsWith("/*")) {
          return globalThis.location.pathname.startsWith(path.slice(0, -1));
        } else {
          return globalThis.location.pathname === path;
        }
      } catch (error) {
        console.log(error);
        if (path.endsWith("/*")) {
          return window.location.pathname.startsWith(path.slice(0, -1));
        } else {
          return window.location.pathname === path;
        }
      }
    });
  };

  const hideNavbar = () => {
    hideNav = hideNavbarList.some((path) => {
      const currentPath = globalThis?.location?.pathname || window.location.pathname;

      if (path.endsWith("/*")) {
        return currentPath.startsWith(path.slice(0, -2));
      } else {
        return currentPath === path;
      }
    });
    console.log(hideNav);
  };

  const changeShowAI = () => {
    triggerShowAI(!$showAI);
    console.log($showAI);
    if (!$showAI) {
      triggerReload(true);
    }
    // isShowAI = !isShowAI;
    localStorage.setItem("showAI", $showAI ? "true" : "false");
  };

  const getModules = async () => {
    return [
      {
        icon: "users",
        name: "crmOverview",
        desc: "descCrm",
        status: "coming soon",
        children: [
          {
            href: isFrontend ? "/crm" : await call.getPage("client-dashboard", "crm-overview"),
            name: "overview",
          },
          {
            href: isFrontend ? "/crm/person" : await call.getPage("client-dashboard", "crm-person"),
            name: "persons",
          },
          {
            href: isFrontend
              ? "/crm/product_service"
              : await call.getPage("client-dashboard", "crm-product-service"),
            name: "product&Services",
          },
          {
            href: isFrontend
              ? "/crm/categories"
              : await call.getPage("client-dashboard", "crm-categories"),
            name: "categories",
          },
          {
            href: isFrontend
              ? "/crm/order_invoice"
              : await call.getPage("client-dashboard", "crm-order-invoice"),
            name: "order&Invoice",
          },
          {
            href: isFrontend
              ? "/crm/calendar"
              : await call.getPage("client-dashboard", "crm-calendar"),
            name: "calendar",
          },
        ],
      },
      {
        icon: "box",
        name: "documentCollection",
        children: [
          {
            href: isFrontend
              ? "/document_collection"
              : await call.getPage("client-dashboard", "document-collection-overview"),
            name: "overview",
          },
          // {
          //   href: isFrontend
          //     ? "/document/document_management"
          //     : await call.getPage("client-dashboard", "document-management"),
          //   name: "documentManagement",
          // },
          {
            href: isFrontend
              ? "/document_collection/collections"
              : await call.getPage("client-dashboard", "document-collection"),
            name: "collections",
          },
          // {
          //   href: isFrontend
          //     ? "/document/contract"
          //     : await call.getPage("client-dashboard", "document-contract"),
          //   name: "contract",
          // },
          // {
          //   href: isFrontend
          //     ? "/document/finetuning"
          //     : await call.getPage("client-dashboard", "document-finetuning"),
          //   name: "finetuning",
          // },
        ],
      },
      {
        icon: "data",
        name: "datasets",
        href: isFrontend ? "/datasets" : await call.getPage("client-dashboard", "datasets"),
        children: [],
      },
      {
        icon: "messageChat",
        name: "hexabot",
        href: isFrontend ? "/hexabot" : await call.getPage("client-dashboard", "hexabot"),
        children: [],
      },
      // {
      //   icon: "fileAttachment",
      //   name: "documentSigning",
      //   children: [
      //     {
      //       href: isFrontend
      //         ? "/document_signing"
      //         : await call.getPage("client-dashboard", "document-signing-overview"),
      //       name: "overview",
      //     },
      //   ],
      // },
      {
        icon: "userSquare",
        name: "HR",
        children: [
          {
            href: isFrontend ? "/hr" : await call.getPage("client-dashboard", "hr-overview"),
            name: "overview",
          },
          {
            href: isFrontend
              ? "/hr/attendance"
              : await call.getPage("client-dashboard", "hr-attendance"),
            name: "attendance",
          },
          {
            href: isFrontend
              ? "/hr/employee"
              : await call.getPage("client-dashboard", "hr-employee"),
            name: "employee",
          },
          {
            href: isFrontend
              ? "/hr/recruitment"
              : await call.getPage("client-dashboard", "hr-recruitment"),
            name: "recruitment",
          },
          // {
          //   href: isFrontend
          //     ? "/hr/job_level"
          //     : await call.getPage("client-dashboard", "hr-joblevel"),
          //   name: "jobLevel",
          // },
          // {
          //   href: isFrontend
          //     ? "/hr/job_position"
          //     : await call.getPage("client-dashboard", "hr-jobposition"),
          //   name: "jobPosition",
          // },
        ],
      },
      {
        icon: "fileShield",
        name: "contractReviewer",
        children: [
          {
            href: isFrontend
              ? "/contract_reviewer"
              : await call.getPage("client-dashboard", "contract-reviewer-overview"),
            name: "overview",
          },
          {
            href: isFrontend
              ? "/contract_reviewer/project"
              : await call.getPage("client-dashboard", "contract-project-list"),
            name: "project",
          },
          {
            href: isFrontend
              ? "/contract_reviewer/case"
              : await call.getPage("client-dashboard", "contract-case"),
            name: "case",
          },
        ],
      },
    ];
  };

  const logout = async () => {
    sessionStorage.clear();
    localStorage.clear();

    try {
      if (!isFrontend) {
        await call.goToPage("client-login", "home-page");
      }
    } catch (error) {
      console.log(error);
    }
  };

  const reset_database = async () => {
    try {
      await call.call("client-dashboard", "reset", {
        URLQuery: {},
        URLParams: {},
        Params: {},
        Headers: {
          Authorization: "Bearer " + sessionStorage.getItem("token"),
        },
      });

      await logout();
    } catch (error) {
      console.error(error);
    }
  };

  export const load = async () => {
    if (browser) {
      locale.set(globalThis.navigator.language);
    }
    await waitLocale();
  };

  async function CheckSetup() {
    try {
      var resp_setup = await call.call("client-login", "is-need-setup", {
        URLQuery: {},
        URLParams: {},
        Params: {},
        Headers: {},
      });
      console.log(resp_setup);
      if (resp_setup.status === 200) {
        if (resp_setup.result.is_need_setup === true) {
          if (sessionStorage.getItem("token") !== null) sessionStorage.removeItem("token");
          if (sessionStorage.getItem("user") !== null) sessionStorage.removeItem("user");
          await call.goToPage("hexapod-setup", "home-page");
          return;
        }
      }
      let token = sessionStorage.getItem("token");
      if (token !== null) {
        let resp = await call.call("client-login", "is-token-valid", {
          URLQuery: {},
          URLParams: {},
          Params: {
            token: token,
          },
          Headers: {},
        });

        if (resp.status === 200) {
          if (resp.result.confirmed === true && resp.result.blocked === false) {
            console.log("is token valid");
            if (resp.result.authority === "User") {
              globalThis.location.href = await call.getPage("lajnah", "home-page");
              return;
            }
            sessionStorage.setItem("user", JSON.stringify(resp.result));
            userData = resp.result;
            authority = resp.result.authority;
          }
        } else if (resp.status === 401) {
          console.log("is token invalid");
          localStorage.setItem("redirect", globalThis.location.href);
          goto(await call.getPage("client-login", "home-page"));
          return;
        } else {
          console.log("unhandled response");
        }
      } else {
        localStorage.setItem("redirect", globalThis.location.href);
        goto(await call.getPage("client-login", "home-page"));
        return;
      }
    } catch (error) {
      // console.log(error);
      userData = { username: "Hexa Staff", email: "hexacode@hexacode.org" };
    }
  }

  async function getModelsTrocr() {
    let token = sessionStorage.getItem("token");
    try {
      let resp = await call.call("client-dashboard", "get-model-trocr", {
        URLQuery: {},
        URLParams: {},
        Params: {},
        Headers: {
          Authorization: "Bearer " + token,
        },
      });

      if (resp === null) {
        sessionStorage.setItem("model_trocr", "trocr:latest");
        return;
      }

      if (resp.status !== 200) {
        console.log(resp);
        return;
      }

      if (resp.status === 200) {
        if (resp.result.model_trocr !== null) {
          sessionStorage.setItem("model_trocr", resp.result.model_trocr);
        } else {
          sessionStorage.setItem("model_trocr", "trocr:latest");
        }
      }
    } catch (error) {
      console.log(error);
    }
  }

  async function getIndustryId() {
    let token = sessionStorage.getItem("token");
    if (token) {
      const response = await call.call("client-dashboard", "get-company-type", {
        URLQuery: {},
        URLParams: {},
        Params: {},
        Headers: {
          Authorization: `Bearer ${token}`,
        },
      });

      if (response.status !== 200) {
        console.log(response);
        return;
      }

      sessionStorage.setItem(
        "industry_id",
        parseInt(response.result.data.attributes.company_industry),
      );

      if (response.result.data.attributes.sidebar_title !== null) {
        sessionStorage.setItem("sidebar_title", response.result.data.attributes.sidebar_title);
        brand = response.result.data.attributes.sidebar_title;
      } else {
        sessionStorage.setItem("sidebar_title", "Hexapod");
        brand = "Hexapod";
      }

      if (response.result.data.attributes.company_logo.data !== null) {
        let ext = response.result.data.attributes.company_logo.data.attributes.ext;

        let fileType;
        if (ext === ".png") {
          fileType = "image/png";
        } else if (ext === ".jpg") {
          fileType = "image/jpg";
        } else if (ext === ".jpeg") {
          fileType = "image/jpeg";
        } else if (ext === ".svg") {
          fileType = "image/svg+xml";
        } else {
          fileType = "image/png";
        }

        const blob = base64toBlob(response.result.company_logo, fileType);
        logoUrl = URL.createObjectURL(blob);
      }

      let theme_color = response.result.data.attributes.theme_color;

      if (theme_color !== null) {
        if (theme_color === "purple") {
          setPrimaryColor(purpleColors);
        } else if (theme_color === "green") {
          setPrimaryColor(greenColors);
        } else if (theme_color === "blue") {
          setPrimaryColor(blueColors);
        }
      } else {
        setPrimaryColor(purpleColors);
      }
    }
  }

  function base64toBlob(base64Data, fileType) {
    const byteCharacters = atob(base64Data);
    const byteArrays = [];

    for (let offset = 0; offset < byteCharacters.length; offset += 512) {
      const slice = byteCharacters.slice(offset, offset + 512);

      const byteNumbers = new Array(slice.length);
      for (let i = 0; i < slice.length; i++) {
        byteNumbers[i] = slice.charCodeAt(i);
      }

      const byteArray = new Uint8Array(byteNumbers);
      byteArrays.push(byteArray);
    }

    const blob = new Blob(byteArrays, { type: fileType });
    return blob;
  }

  onMount(async () => {
    if (sessionStorage.getItem("token") === null && isFrontend) {
      alert("Please insert token to session storage and refresh the page");
      return;
    }

    await CheckSetup();
    if (localStorage.getItem("lang") !== null) {
      locale.set(localStorage.getItem("lang"));
    }

    setCall(call);

    await getModelsTrocr();
    await getIndustryId();
    hideNavbar();

    const resizeHandler = () => {
      innerWidth = globalThis.innerWidth;
    };
    globalThis.addEventListener("resize", resizeHandler);
    resizeHandler(); // call once initially to set the correct value

    ready = true;
    return () => {
      globalThis.removeEventListener("resize", resizeHandler);
    };
  });

  let winInnerWidth = $derived(innerWidth);
  if (typeof localStorage !== "undefined") {
    triggerShowAI(localStorage.getItem("showAI") === "true" ? true : false);
    toggleSidebarLargeClosed(
      localStorage.getItem("sidebarLargeClosed") === "true" && globalThis.innerWidth > 1024
        ? true
        : false,
    );
    if (localStorage.getItem("selectedDropdown") === "collections") {
      triggerSelectDropdown(
        localStorage.getItem("selectedDropdown") === "collections" ? "collections" : "general",
      );
      if (localStorage.getItem("selectedCollection") !== null) {
        triggerSelectCollection(JSON.parse(localStorage.getItem("selectedCollection")));
      }
    }
  }

  $effect(() => {
    console.log(activeUrl);
    if (activeUrl !== $page.url.pathname) {
      activeUrl = $page.url.pathname;
      hideNavbar();
    }
  });
</script>

{#if $isLoading}
  <!-- <LoadingSpinner /> -->
{:else if winInnerWidth > 0 && $callReady === true && ready === true}
  {#if hideLayout()}
    <div class="lg:flex">
      <div class="z-20">
        <Sidebar
          getModule={getModules}
          bind:name={userData.username}
          bind:email={userData.email}
          bind:brand
          bind:logoUrl
          isFrontend
          innerWidth={winInnerWidth}
          {logout}
        />
      </div>
      <div
        class={`w-full h-screen pt-16 md:pt-20 lg:pt-0 flex border-l border-secondary-light dark:border-secondary-dark ${!$sidebarLargeClosed ? "lg:pl-80" : "lg:pl-20"} transition-all duration-300 ease-in-out`}
      >
        <!-- Navbar -->
        <Splitpanes class="[&>*]:!bg-primary-light [&>*]:dark:!bg-primary-dark [&>*]:!border-none">
          <Pane minSize={innerWidth > 1024 ? 80 : 0}>
            <div class="relative transition-all duration-300 h-screen">
              {#if !hideNav}
                <div class="absolute w-full top-0 z-10">
                  <Navbar
                    {activeUrl}
                    {$showAI}
                    {changeShowAI}
                    {appUrl}
                    {reset_database}
                    bind:authority
                  />
                </div>
              {/if}
              <div class={`h-screen ${!hideNav ? "pt-12 lg:pt-20 pb-12 lg:pb-0" : ""}`}>
                {#if hideScroll()}
                  <div class="h-full overflow-y-auto">
                    <slot />
                  </div>
                {:else}
                  <div class="h-full">
                    <slot />
                  </div>
                {/if}
                <ToastContainer />
              </div>
              <LoadingContainer />
            </div>
          </Pane>
          <!-- Hexabot Sidebar -->
          {#if $showAI}
            <SideDrawer innerWidth={winInnerWidth} {changeShowAI} />
          {/if}
        </Splitpanes>
      </div>
    </div>
  {:else}
    <slot />
  {/if}
{:else}{/if}

<style>
  /* Hide scrollbar for Chrome, Safari and Opera */
  .no-scrollbar::-webkit-scrollbar {
    display: none;
  }

  /* Hide scrollbar for IE, Edge and Firefox */
  .no-scrollbar {
    -ms-overflow-style: none; /* IE and Edge */
    scrollbar-width: none; /* Firefox */
  }
</style>
