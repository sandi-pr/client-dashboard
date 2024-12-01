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
    triggerShowNotif,
    showNotif,
    showAI,
    triggerSelectDropdown,
    triggerSelectCollection,
    triggerReload,
    triggerToast,
    industryData,
    notifications,
    notificationsUnread,
    notificationsRead,
    notificationsNewCount,
    newNotification,
    paginationAll,
    paginationRead,
    paginationUnread,
    connectWebSocket,
    disconnectWebSocket,
    setLibraryIds,
    setCollectionIds,
    bakrieLibraryId,
    handwritingCollectionId,
    onProgress,
  } from "$lib";
  import { onDestroy, onMount } from "svelte";
  import { locale, waitLocale, isLoading, _ } from "svelte-i18n";
  import "$lib/service/i18n";
  import SideDrawer from "$lib/components/SideDrawer.svelte";
  import NotifDrawer from "$lib/components/NotifDrawer.svelte";
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
  let gettingNewNotif = $state(false);

  const call = isFrontend
    ? new CR3Call("http://192.168.18.14")
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
    if ($showNotif) {
      triggerShowNotif(false);
    }
    // isShowAI = !isShowAI;
    localStorage.setItem("showAI", $showAI ? "true" : "false");
  };

  const changeShowNotif = () => {
    if ($showAI) {
      triggerShowAI(false);
    }
    if (!$showNotif === true) {
      hasNewNotif = false;
    }
    triggerShowNotif(!$showNotif);
    triggerReload(true);
    localStorage.setItem("showNotif", $showNotif ? "true" : "false");

    console.log($showNotif);
  };

  const getModules = async () => {
    return [
      // {
      //   icon: "briefcase",
      //   name: "crmOverview",
      //   children: [
      //     // {
      //     //   href: isFrontend ? "/crm" : await call.getPage("client-dashboard", "crm-overview"),
      //     //   name: "overview",
      //     // },
      //     {
      //       icon: "user",
      //       href: isFrontend ? "/crm/person" : await call.getPage("client-dashboard", "crm-person"),
      //       name: "persons",
      //     },
      //     {
      //       icon: "faceHappy",
      //       href: isFrontend ? "/hr" : await call.getPage("client-dashboard", "hr-overview"),
      //       name: "hR",
      //     },
      //     {
      //       icon: "package",
      //       href: isFrontend
      //         ? "/crm/product_service"
      //         : await call.getPage("client-dashboard", "crm-product-service"),
      //       name: "product&Services",
      //     },
      //     // {
      //     //   href: isFrontend
      //     //     ? "/crm/categories"
      //     //     : await call.getPage("client-dashboard", "crm-categories"),
      //     //   name: "categories",
      //     // },
      //     {
      //       icon: "receipt",
      //       href: isFrontend
      //         ? "/crm/order_invoice"
      //         : await call.getPage("client-dashboard", "crm-order-invoice"),
      //       name: "order&Invoice",
      //     },
      //     {
      //       icon: "calendar",
      //       href: isFrontend
      //         ? "/crm/calendar"
      //         : await call.getPage("client-dashboard", "crm-calendar"),
      //       name: "calendar",
      //     },
      //   ],
      // },
      // await checkLibrariesIndustry(sessionStorage.getItem("industry_id")),
      {
        icon: "stars",
        name: "aiManagement",
        children: [
          // {
          //   icon: "messageChat",
          //   name: "hexabot",
          //   href: isFrontend ? "/hexabot" : await call.getPage("client-dashboard", "hexabot"),
          // },
          {
            icon: "archive",
            name: "documentLibrary",
            href: isFrontend
              ? "/document_collection"
              : await call.getPage("client-dashboard", "document-collection-overview"),
          },
          // {
          //   icon: "archive",
          //   name: "documentCollection",
          //   href: isFrontend
          //     ? "/document_collection/collections"
          //     : await call.getPage("client-dashboard", "document-collection-overview"),
          // },
          // {
          //   icon: "fileCheck",
          //   name: "contractReviewer",
          //   href: isFrontend
          //     ? "/contract_reviewer"
          //     : await call.getPage("client-dashboard", "contract-reviewer-overview"),
          // },
          // {
          //   icon: "data",
          //   name: "datasets",
          //   href: isFrontend ? "/datasets" : await call.getPage("client-dashboard", "datasets"),
          // },
        ],
      },
      // {
      //   icon: "box",
      //   name: "documentCollection",
      //   children: [
      //     {
      //       href: isFrontend
      //         ? "/document_collection"
      //         : await call.getPage("client-dashboard", "document-collection-overview"),
      //       name: "overview",
      //     },
      //     // {
      //     //   href: isFrontend
      //     //     ? "/document/document_management"
      //     //     : await call.getPage("client-dashboard", "document-management"),
      //     //   name: "documentManagement",
      //     // },
      //     {
      //       href: isFrontend
      //         ? "/document_collection/collections"
      //         : await call.getPage("client-dashboard", "document-library"),
      //       name: "collections",
      //     },
      //     // {
      //     //   href: isFrontend
      //     //     ? "/document/contract"
      //     //     : await call.getPage("client-dashboard", "document-contract"),
      //     //   name: "contract",
      //     // },
      //     // {
      //     //   href: isFrontend
      //     //     ? "/document/finetuning"
      //     //     : await call.getPage("client-dashboard", "document-finetuning"),
      //     //   name: "finetuning",
      //     // },
      //   ],
      // },
      // {
      //   icon: "data",
      //   name: "datasets",
      //   href: isFrontend ? "/datasets" : await call.getPage("client-dashboard", "datasets"),
      //   children: [],
      // },
      // {
      //   icon: "messageChat",
      //   name: "hexabot",
      //   href: isFrontend ? "/hexabot" : await call.getPage("client-dashboard", "hexabot"),
      //   children: [],
      // },
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
      // {
      //   icon: "userSquare",
      //   name: "HR",
      //   children: [
      //     {
      //       href: isFrontend ? "/hr" : await call.getPage("client-dashboard", "hr-overview"),
      //       name: "overview",
      //     },
      //     {
      //       href: isFrontend
      //         ? "/hr/attendance"
      //         : await call.getPage("client-dashboard", "hr-attendance"),
      //       name: "attendance",
      //     },
      //     {
      //       href: isFrontend
      //         ? "/hr/employee"
      //         : await call.getPage("client-dashboard", "hr-employee"),
      //       name: "employee",
      //     },
      //     {
      //       href: isFrontend
      //         ? "/hr/recruitment"
      //         : await call.getPage("client-dashboard", "hr-recruitment"),
      //       name: "recruitment",
      //     },
      //     // {
      //     //   href: isFrontend
      //     //     ? "/hr/job_level"
      //     //     : await call.getPage("client-dashboard", "hr-joblevel"),
      //     //   name: "jobLevel",
      //     // },
      //     // {
      //     //   href: isFrontend
      //     //     ? "/hr/job_position"
      //     //     : await call.getPage("client-dashboard", "hr-jobposition"),
      //     //   name: "jobPosition",
      //     // },
      //   ],
      // },
      // {
      //   icon: "fileShield",
      //   name: "contractReviewer",
      //   children: [
      //     {
      //       href: isFrontend
      //         ? "/contract_reviewer"
      //         : await call.getPage("client-dashboard", "contract-reviewer-overview"),
      //       name: "overview",
      //     },
      //     {
      //       href: isFrontend
      //         ? "/contract_reviewer/project"
      //         : await call.getPage("client-dashboard", "contract-project-list"),
      //       name: "project",
      //     },
      //     {
      //       href: isFrontend
      //         ? "/contract_reviewer/case"
      //         : await call.getPage("client-dashboard", "contract-case"),
      //       name: "case",
      //     },
      //   ],
      // },
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
          Authorization: "Bearer " + localStorage.getItem("token"),
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
          if (localStorage.getItem("token") !== null) localStorage.removeItem("token");
          if (localStorage.getItem("user") !== null) localStorage.removeItem("user");
          await call.goToPage("hexapod-setup", "home-page");
          return;
        }
      }
      let token = localStorage.getItem("token");
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
            // if (resp.result.authority === "User") {
            //   globalThis.location.href = await call.getPage("lajnah", "home-page");
            //   return;
            // }
            localStorage.setItem("user", JSON.stringify(resp.result));
            await getPermissions(resp.result.authority.toLowerCase().replace(/\s+/g, "_"));
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
    let token = localStorage.getItem("token");
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

  async function checkLibrariesIndustry(industryId) {
    console.log("checkLibrariesIndustry", industryId);
    const industry = industryData.find((item) => item.id === parseInt(industryId));
    if (!industry) return;

    let industryLibrary = industry.points.map((point) => ({
      name: point.title,
      desc: point.desc,
    }));

    let initIndustryLibrary = industry.points.map((point) => ({
      href: "",
      name: point.title,
      icon: "bookClosed",
    }));

    for (const point of industryLibrary) {
      let token = localStorage.getItem("token");
      if (token) {
        const response = await call.call("client-dashboard", "find-create-library-by-name", {
          URLQuery: {},
          URLParams: {},
          Params: {
            library_name: point.name,
            description: point.desc,
          },
          Headers: {
            Authorization: `Bearer ${token}`,
          },
        });

        if (response.status !== 200) {
          console.log(response);
          return;
        }

        let link =
          (await call.getPage("client-dashboard", "document-library-collections")) +
          "?libraryId=" +
          response.result.id;

        initIndustryLibrary.find((item) => item.name === point.name).href = link;
      }
    }

    const useCase = {
      icon: "stars",
      name: "documentAi",
      children: [...initIndustryLibrary],
    };

    console.log("usecase", useCase);

    return useCase;
  }

  async function checkLibrariesBakrie() {
    let token = localStorage.getItem("token");
    if (token) {
      const response = await call.call("client-dashboard", "find-create-library-by-name", {
        URLQuery: {},
        URLParams: {},
        Params: {
          library_name: "Bakrie Library",
          description: "Bakrie Library Collection",
        },
        Headers: {
          Authorization: `Bearer ${token}`,
        },
      });

      if (response.status !== 200) {
        console.log("status not 200", response);
        return;
      }

      console.log("created bakrie library", response);

      localStorage.setItem("libraryIdWebchatBakrie", response.result.id);
      setLibraryIds(response.result.id);
    }
  }

  async function checkCollectionsBakrie() {
    const bakrieCollection = [
      {
        name: "Documents",
        desc: "collection documents",
      },
      {
        name: "Images",
        desc: "collection images",
      },
      {
        name: "Document Images",
        desc: "collection document images",
      },
      {
        name: "Handwriting",
        desc: "collection handwriting",
      },
      {
        name: "Document Audio",
        desc: "collection audio",
      },
    ];

    const libraryIds = {
      documentsCollectionIdValue: null,
      imagesCollectionIdValue: null,
      documentImagesCollectionIdValue: null,
      handwritingCollectionIdValue: null,
      audioCollectionIdValue: null,
    };

    for (const point of bakrieCollection) {
      let token = localStorage.getItem("token");
      if (token) {
        const response = await call.call(
          "client-dashboard",
          "find-create-library-collection-by-name",
          {
            URLQuery: {},
            URLParams: {},
            Params: {
              parent_id: $bakrieLibraryId,
              collection_name: point.name,
              description: point.desc,
            },
            Headers: {
              Authorization: `Bearer ${token}`,
            },
          },
        );

        if (response.status !== 200) {
          console.log("status not 200", response);
          return;
        }

        // Save the collection ID based on the collection name
        switch (point.name) {
          case "Documents":
            localStorage.setItem("collectionIdDocumentsBakrie", response.result.id);
            libraryIds.documentsCollectionIdValue = response.result.id;
            break;
          case "Images":
            localStorage.setItem("collectionIdImagesBakrie", response.result.id);
            libraryIds.imagesCollectionIdValue = response.result.id;
            break;
          case "Document Images":
            localStorage.setItem("collectionIdDocumentImagesBakrie", response.result.id);
            libraryIds.documentImagesCollectionIdValue = response.result.id;
            break;
          case "Handwriting":
            localStorage.setItem("collectionIdHandwritingBakrie", response.result.id);
            libraryIds.handwritingCollectionIdValue = response.result.id;
            break;
          case "Document Audio":
            localStorage.setItem("collectionIdAudioBakrie", response.result.id);
            libraryIds.audioCollectionIdValue = response.result.id;
            break;
        }
      }
    }

    setCollectionIds(libraryIds);
  }

  async function createBakrieBot() {
    let token = localStorage.getItem("token");

    if (token !== null) {
      const resp = await call.call("client-dashboard", "get-bot-by-name", {
        URLQuery: {},
        URLParams: {},
        Params: {
          name: "Bakrie’s Assistant",
        },
        Headers: {
          Authorization: "Bearer " + localStorage.getItem("token"),
        },
      });

      console.log("------", resp);

      if (resp.status !== 404) {
        console.log("404", resp);
        localStorage.setItem("botIdStrapi", resp.result.id);
        // gettingData = false;
        return;
      }

      let dataSend = {
        bot_name: "Bakrie’s Assistant",
        description: "",
        ai_engine: "llama3.1:8b",
        library_id: $bakrieLibraryId,
        collection_id: "all",
        additional_prompt: "",
        need_auth: false,
      };

      // console.log("dataSend", dataSend);

      try {
        const response = await call.call("client-dashboard", "create-bot", {
          URLQuery: {},
          URLParams: {},
          Params: dataSend,
          Headers: {
            Authorization: "Bearer " + token,
          },
        });

        console.log("createBakrieBot", response);

        if (response.status !== 200) {
          console.log(response);
          return;
        }

        localStorage.setItem("botIdStrapi", response.result.data.id);
      } catch (error) {
        console.log(error);
      }
    } else {
      console.log("no token");
      return;
    }
  }

  async function getIndustryId() {
    let token = localStorage.getItem("token");
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

      // await checkLibrariesIndustry(response.result.data.attributes.company_industry);

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

  async function getPermissions(role) {
    let token = localStorage.getItem("token");
    if (token) {
      const response = await call.call("module-manager", "get-role", {
        URLQuery: {},
        URLParams: {
          role: role,
        },
        Params: {},
        Headers: {
          Authorization: `Bearer ${token}`,
        },
      });

      if (response.status !== 200) {
        console.log(response);
        return;
      }

      console.log("itu", response.result.role.permissions);
      localStorage.setItem("permissions", JSON.stringify(response.result.role.permissions));
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

  async function getNotifications(refresh = false) {
    let token = localStorage.getItem("token");

    if (refresh) {
      notifications.set([]);
      paginationAll.set({
        page: 1,
        page_size: 10,
        page_count: 0,
        total: 0,
      });
    }

    try {
      const response = await call.call("client-dashboard", "get-notifications", {
        URLQuery: {},
        URLParams: {},
        Params: {
          page: $paginationAll.page,
          page_size: $paginationAll.page_size,
        },
        Headers: {
          Authorization: "Bearer " + token,
        },
      });

      if (response.status !== 200) {
        console.log(response);
        return;
      }

      if (response.result.data.length > 0) {
        notifications.update((data) => [...response.result.data, ...data]);
      }
    } catch (error) {
      console.log(error);
    }
  }

  async function getNotificationsUnread(refresh = false) {
    let token = localStorage.getItem("token");

    if (refresh) {
      notificationsUnread.set([]);
      paginationUnread.set({
        page: 1,
        page_size: 10,
        page_count: 0,
        total: 0,
      });
    }

    try {
      const response = await call.call("client-dashboard", "get-notifications-unread", {
        URLQuery: {},
        URLParams: {},
        Params: {
          page: $paginationUnread.page,
          page_size: $paginationUnread.page_size,
        },
        Headers: {
          Authorization: "Bearer " + token,
        },
      });

      if (response.status !== 200) {
        console.log(response);
        return;
      }

      if (response.result.data.length > 0) {
        notificationsUnread.update((data) => [...response.result.data, ...data]);
        notificationsNewCount.set(response.result.meta.pagination.total);
      }
    } catch (error) {
      console.log(error);
    }
  }

  async function getNotificationsRead() {
    let token = localStorage.getItem("token");

    try {
      const response = await call.call("client-dashboard", "get-notifications-read", {
        URLQuery: {},
        URLParams: {},
        Params: {
          page: $paginationRead.page,
          page_size: $paginationRead.page_size,
        },
        Headers: {
          Authorization: "Bearer " + token,
        },
      });

      if (response.status !== 200) {
        console.log(response);
        return;
      }

      if (response.result.data.length > 0) {
        notificationsRead.update((data) => [...response.result.data, ...data]);
      }
    } catch (error) {
      console.log(error);
    }
  }

  onMount(async () => {
    if (localStorage.getItem("token") === null && isFrontend) {
      const token = prompt("Please insert your token:", "");
      if (token) {
        localStorage.setItem("token", token);
        location.reload();
      } else {
        return;
      }
    }

    setCall(call);

    await CheckSetup();
    await checkLibrariesBakrie();
    await checkCollectionsBakrie();
    await createBakrieBot();
    // await getPermissions();
    if (localStorage.getItem("lang") !== null) {
      locale.set(localStorage.getItem("lang"));
    }

    await getModelsTrocr();
    await getIndustryId();
    hideNavbar();

    getNotifications();
    getNotificationsUnread();
    getNotificationsRead();

    connectWebSocket(localStorage.getItem("token"), isFrontend);

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

  let hasNewNotif = $state(false);

  // $effect(() => {
  //   if ($notificationsNewCount > 0) {
  //     hasNewNotif = true;
  //   }
  // });

  // $effect(() => {
  //   if ($notificationsNewCount < 1) {
  //     hasNewNotif = false;
  //   }
  // });

  $effect(async () => {
    if ($newNotification) {
      if (!gettingNewNotif) {
        gettingNewNotif = true;
        hasNewNotif = true;
        await getNotifications(true);
        await getNotificationsUnread(true);
        newNotification.set(false);
        gettingNewNotif = false;
      }
    }
  });

  // onDestroy(() => {
  //   disconnectWebSocket();
  // });
</script>

<svelte:head>
  <title>{brand}</title>
</svelte:head>

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
        class={`w-full h-screen ${winInnerWidth > 1024 ? "pt-16 md:pt-20" : $showAI ? "pt-0" : "pt-16"} lg:pt-0 flex border-l border-secondary-light dark:border-secondary-dark ${!$sidebarLargeClosed ? "lg:pl-80" : "lg:pl-20"} transition-all duration-300 ease-in-out`}
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
                    {changeShowNotif}
                    {appUrl}
                    {reset_database}
                    bind:authority
                    bind:hasNewNotif
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
          {#if $showNotif}
            <NotifDrawer innerWidth={winInnerWidth} {changeShowNotif} />
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
