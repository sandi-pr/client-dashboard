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
    setDatasourceLibraryIds,
    datasourceLibraryId,
    onProgress,
    isBakrie,
    check,
    PageComingSoon,
    setCookie,
    getCookie,
    isMoodle,
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
  import FingerprintJS from "@fingerprintjs/fingerprintjs";
  import CryptoJS from "crypto-js";

  set_pdfjs_context();

  let innerWidth = $state(0);
  let activeUrl = $state($page.url.pathname);
  let appUrl = "/";
  let userData = $state({ username: "", email: "" });
  let authority = $state("");
  let brand = $state("");
  let logoUrl = $state("");
  const isFrontend = process.env.NODE_ENV === "development";
  let hideNav = $state(false);
  let ready = $state(false);
  let gettingNewNotif = $state(false);

  const call = isFrontend
    ? new CR3Call("https://tank.hexacode.ai")
    : new CR3Call(globalThis.location.origin);

  let hideBakrie = $state(true);

  const hideLayoutListBakrie = isFrontend
    ? ["/document_collection", "/document_collection/libraries"]
    : ["/client-dashboard/documents", "/client-dashboard/documents/libraries"];

  const hideLayoutList = $state(["/workflow/update", "/workflow/create", "/hexabot/webchat"]);

  const hideScrollList = $state(["/document/document_collection/chat"]);

  const hideNavbarList = $state(["/client-dashboard/datasets", "/datasets"]);

  const pageComingSoon = $state(["/support"]);
  let isComingSoon = $state(false);

  const hideScroll = () => {
    return !hideScrollList.some((path) => {
      try {
        if (path.endsWith("/*")) {
          return globalThis.location.pathname.startsWith(path.slice(0, -1));
        } else {
          return globalThis.location.pathname === path;
        }
      } catch (error) {
        // console.log(error);
        if (path.endsWith("/*")) {
          return window.location.pathname.startsWith(path.slice(0, -1));
        } else {
          return window.location.pathname === path;
        }
      }
    });
  };

  const hideLayout = () => {
    return (
      !hideLayoutList.some((path) => {
        try {
          if (path.endsWith("/*")) {
            return globalThis.location.pathname.startsWith(path.slice(0, -1));
          } else {
            return globalThis.location.pathname === path;
          }
        } catch (error) {
          // console.log(error);
          if (path.endsWith("/*")) {
            return window.location.pathname.startsWith(path.slice(0, -1));
          } else {
            return window.location.pathname === path;
          }
        }
      }) && $page.status === 200
    );
  };

  const hideLayoutBakrie = () => {
    if ($isBakrie) {
      const result = !hideLayoutListBakrie.some((path) => {
        const currentPath = window.location.pathname.trim();
        const normalizedPath = path.trim();
        return currentPath === normalizedPath;
      });
      return result;
    } else {
      return true;
    }
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
    // console.log(hideNav);
  };

  const changeShowAI = async () => {
    if ($isBakrie) {
      let url = await call.getPage("bakrie", "home-page");
      goto(url);
      return;
    } else {
      // console.log($showAI);
      triggerShowAI(!$showAI);
      if (!$showAI) {
        triggerReload(true);
      }
      if ($showNotif) {
        triggerShowNotif(false);
      }
      // isShowAI = !isShowAI;
      localStorage.setItem("showAI", $showAI ? "true" : "false");
    }
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

    // console.log($showNotif);
  };

  const getModules = async () => {
    return [
      // Start of new sidebar
      {
        icon: "messageChatSquare",
        name: "assistant",
        href: isFrontend ? "/assistant" : await call.getPage("client-dashboard", "assistant"),
        children: [
          // {
          //   icon: "messageHeartSquare",
          //   name: "hexassistant",
          //   href: isFrontend
          //     ? "/hexassistant"
          //     : await call.getPage("client-dashboard", "hexassistant"),
          // },
          // {
          //   icon: "messageSmileSquare",
          //   name: "myBots",
          //   href: isFrontend ? "/my_bots" : await call.getPage("client-dashboard", "my-bots"),
          // },
          // {
          //   icon: "code01",
          //   name: "codingAssistant",
          //   href: isFrontend
          //     ? "/coding_assistant"
          //     : await call.getPage("client-dashboard", "coding-assistant"),
          // },
        ],
      },
      // {
      //   icon: "code01",
      //   name: "codingAssistant",
      //   commingSoon: true,
      //   children: [],
      //   href: isFrontend
      //     ? "/coding_assistant"
      //     : await call.getPage("client-dashboard", "coding-assistant"),
      // },
      {
        icon: "bookOpen",
        name: "documents",
        href: isFrontend
          ? "/documents"
          : await call.getPage("client-dashboard", "document-collection-overview"),
        children: [],
      },
      // {
      //   icon: "marketplace",
      //   name: "marketplace",
      //   commingSoon: true,
      //   divider: true,
      //   href: isFrontend
      //     ? "/marketplace"
      //     : await call.getPage("client-dashboard", "marketplace"),
      //     children: []
      // },
      // {
      //   icon: "graduationHat02",
      //   name: "hexatutor",
      //   children: [
      //     {
      //       icon: "graduationHat01",
      //       name: "myTutor",
      //       href: isFrontend ? "/my_tutor" : await call.getPage("client-dashboard", "my-tutor"),
      //     },
      //     {
      //       icon: "coding",
      //       name: "codingAssistant",
      //       href: isFrontend ? "/coding_assistant" : await call.getPage("client-dashboard", "coding-assistant"),
      //     },
      //   ],
      // },
      // {
      //   icon: "presentation",
      //   name: "Analytics",
      //   children: [
      //     {
      //       icon: "bookClosed",
      //       name: "dataSources",
      //       href: isFrontend
      //         ? "/data_sources"
      //         : await call.getPage("client-dashboard", "data-sources"),
      //     },
      //     {
      //       icon: "lightbulb",
      //       name: "insights",
      //       href: isFrontend ? "/csv_insight" : await call.getPage("client-dashboard", "csv-insight"),
      //     },
      //   ],
      // },
      // todo config
      {
        icon: "presentation",
        name: "analytics",
        children: [
          {
            icon: "bookClosed",
            name: "dataSources",
            href: isFrontend
              ? "/data_sources"
              : await call.getPage("client-dashboard", "data-sources"),
          },
          {
            icon: "lightbulb",
            name: "insights",
            href: isFrontend
              ? "/csv_insight"
              : await call.getPage("client-dashboard", "csv-insight"),
          },
        ],
      },
      //
      // {
      //   icon: "star04",
      //   name: "aiAnalyzer",
      //   children: [
      // todo config
      {
        icon: "penTool",
        name: "handwriting",
        children: [],
        href: isFrontend ? "/handwriting" : await call.getPage("client-dashboard", "handwriting"),
      },
      {
        icon: "recording02",
        name: "voice",
        children: [],
        href: isFrontend ? "/voice" : await call.getPage("client-dashboard", "voice"),
      },
      //
      // {
      //   icon: "image03",
      //   name: "memories",
      //   children: [],
      //   commingSoon: true,
      //   href: isFrontend ? "/memories" : await call.getPage("client-dashboard", "memories"),
      // },
      // {
      //   icon: "videoRecorder",
      //   name: "vision",
      //   children: [],
      //   commingSoon: true,
      //   href: isFrontend ? "/vision" : await call.getPage("client-dashboard", "vision"),
      // },
      //   ],
      // },
      // {
      //   icon: "aiDocument",
      //   name: "smartDocument",
      //   children: [
      // {
      //   icon: "fileLock",
      //   name: "contractReviewer",
      //   children: [],
      //   commingSoon: true,
      //   href: isFrontend
      //     ? "/contract_reviewer"
      //     : await call.getPage("client-dashboard", "contract-reviewer-overview"),
      // },
      //     {
      //       icon: "fileLock",
      //       name: "sopCompliance",
      //       href: isFrontend ? "/sop_compliance" : await call.getPage("client-dashboard", "sop-compliance"),
      //       children: []
      //     },
      // {
      //   icon: "fileCheck",
      //   name: "docsGenerator",
      //   children: [],
      //   commingSoon: true,
      //   href: isFrontend
      //     ? "/docs_generator"
      //     : await call.getPage("client-dashboard", "docs-generator"),
      // },
      //   ],
      // },
      {
        icon: "globe01",
        name: "llms",
        href: isFrontend ? "/llms" : await call.getPage("client-dashboard", "llms"),
        children: [
          // {
          //   icon: "dashboard01",
          //   name: "overview",
          //   href: isFrontend ? "/llms" : await call.getPage("client-dashboard", "llms"),
          //   children: [],
          // },
          // {
          //   icon: "dataflow04",
          //   name: "lists",
          //   href: isFrontend ? "/llms/list" : await call.getPage("client-dashboard", "dataflow"),
          //   children: [],
          // },
          // {
          //   icon: "meta",
          //   name: "Llama3.2",
          //   href: isFrontend
          //     ? "/llms?tab=llama"
          //     : await call.getPage("client-dashboard", "llms"),
          // },
          // {
          //   icon: "gemini",
          //   name: "Gemma 2",
          //   href: isFrontend
          //     ? "/llms?tab=gemma"
          //     : await call.getPage("client-dashboard", "llms"),
          // },
          // {
          //   icon: "deepseek",
          //   name: "DeepSeek R1",
          //   href: isFrontend
          //     ? "/llms?tab=deepseek"
          //     : await call.getPage("client-dashboard", "llms"),
          // }
        ],
      },
      // E-learning
      {
        icon: "graduationHat02",
        name: "eLearning",
        commingSoon: true,
        href: isFrontend ? "/e_learning" : await call.getPage("client-dashboard", "e-learning"),
        children: [],
      },
      {
        icon: "settings03",
        name: "configurations",
        children: [
          {
            icon: "messagePlusSquare",
            name: "bots",
            href: isFrontend ? "/bots" : await call.getPage("client-dashboard", "bots"),
            children: [],
          },
          {
            icon: "data",
            name: "finetuning",
            href: isFrontend
              ? "/finetunings"
              : await call.getPage("client-dashboard", "finetunings"),
          },
          {
            icon: "dataflow02",
            name: "apiConnector",
            href: isFrontend
              ? "/api_connector"
              : await call.getPage("client-dashboard", "api-connector"),
          },
        ],
      },
      // {
      //   icon: "phoneCall",
      //   name: "support",
      //   href: isFrontend ? "/support" : await call.getPage("client-dashboard", "support"),
      //   commingSoon: true,
      //   children: [],
      // },
      // End of new sidebar
    ].filter(Boolean);
  };

  const loadFingerprint = async () => {
    const fp = await FingerprintJS.load();

    const result = await fp.get();

    return result.visitorId;
  };

  const logout = async () => {
    const device_id = await loadFingerprint();
    var resp = await call.call("client-login", "logout", {
      URLQuery: {},
      URLParams: {},
      Params: {
        device_id: device_id,
      },
      Headers: {
        Authorization: "Bearer " + localStorage.getItem("token"),
      },
    });

    if (resp.status !== 200) {
      // console.log(resp);
      triggerToast("error");
      return;
    }

    sessionStorage.clear();
    localStorage.clear();

    try {
      if (!isFrontend) {
        await call.goToPage("client-login", "home-page");
      }
    } catch (error) {
      // console.log(error);
    }
  };

  // const reset_database = async () => {
  //   try {
  //     await call.call("client-dashboard", "reset", {
  //       URLQuery: {},
  //       URLParams: {},
  //       Params: {},
  //       Headers: {
  //         Authorization: "Bearer " + localStorage.getItem("token"),
  //       },
  //     });

  //     await logout();
  //   } catch (error) {
  //     console.error(error);
  //   }
  // };

  export const load = async () => {
    if (browser) {
      locale.set(globalThis.navigator.language);
    }
    await waitLocale();
  };

  async function CheckSetup() {
    var resp_setup = await call.call("client-login", "is-need-setup", {
      URLQuery: {},
      URLParams: {},
      Params: {},
      Headers: {},
    });
    // console.log("===============================");
    // console.log(resp_setup);
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
          // console.log("is token valid");
          localStorage.setItem("user", JSON.stringify(resp.result));
          await getPermissions(resp.result.authority.toLowerCase().replace(/\s+/g, "_"));
          userData = resp.result;
          authority = resp.result.authority;
        }
      } else if (resp.status === 401) {
        let email = localStorage.getItem("email") || "";
        let passwordTmp = localStorage.getItem("password") || "";

        if (
          email === "" ||
          passwordTmp === "" ||
          email === undefined ||
          passwordTmp === undefined ||
          email === null ||
          passwordTmp === null
        ) {
          // console.log("is token invalid");
          localStorage.setItem("redirect", globalThis.location.href);
          goto(await call.getPage("client-login", "home-page"));
          return;
        }
        let password =
          CryptoJS.AES.decrypt(passwordTmp, "hexacode-admin").toString(CryptoJS.enc.Utf8) || "";

        const device_id = await loadFingerprint();

        var resp3 = await call.call("client-login", "login", {
          URLQuery: {},
          URLParams: {},
          Params: {
            email: email,
            password: password,
            device_id: device_id,
          },
          Headers: {},
        });

        if (resp3.status === 200) {
          // login success
          if (resp3.result?.token !== null && resp3.result?.token !== undefined) {
            localStorage.setItem("token", resp3.result.token);

            var resp2 = await call.call("client-login", "is-token-valid", {
              URLQuery: {},
              URLParams: {},
              Params: {
                token: resp.result.token,
              },
              Headers: {},
            });

            if (resp2.status === 200) {
              if (resp2.result.confirmed === true && resp2.result.blocked === false) {
                localStorage.setItem("user", resp2.result);
                // globalThis.location.href = await call.getPage(
                //   resp.result.redirect_page.namespace,
                //   resp.result.redirect_page.page,
                // );
                location.reload();
                return;
              }
            }

            // globalThis.location.href = await call.getPage(
            //   resp.result.redirect_page.namespace,
            //   resp.result.redirect_page.page,
            // );
            location.reload();

            return;
          }
          // else if (resp3.result?.need_otp) {
          //   globalThis.location.href = await call.getPage(
          //     resp3.result.redirect_page.namespace,
          //     resp3.result.redirect_page.page,
          //   );

          //   return;
          // }
          else {
            // console.log("is token invalid");
            localStorage.setItem("redirect", globalThis.location.href);
            goto(await call.getPage("client-login", "home-page"));
            return;
          }
        } else {
          // console.log("is token invalid");
          localStorage.setItem("redirect", globalThis.location.href);
          goto(await call.getPage("client-login", "home-page"));
          return;
        }
      } else {
        // console.log("is token invalid");
        localStorage.setItem("redirect", globalThis.location.href);
        goto(await call.getPage("client-login", "home-page"));
        return;
      }
    } else {
      localStorage.setItem("redirect", globalThis.location.href);
      goto(await call.getPage("client-login", "home-page"));
      return;
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
        // console.log(resp);
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
      // console.log(error);
    }
  }

  async function checkLibrariesIndustry(industryId) {
    // console.log("checkLibrariesIndustry", industryId);
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
          // console.log(response);
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

    // console.log("usecase", useCase);

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
        // console.log("status not 200", response);
        return;
      }

      // console.log("created bakrie library", response);

      localStorage.setItem("libraryIdWebchatBakrie", response.result.id);
      setLibraryIds(response.result.id);
    }
  }

  // async function checkLibrariesDatasouce() {
  //   let token = localStorage.getItem("token");
  //   if (token) {
  //     const response = await call.call("client-dashboard", "find-create-library-by-name", {
  //       URLQuery: {},
  //       URLParams: {},
  //       Params: {
  //         library_name: "Data Source List",
  //         description: "Data Source Collection",
  //       },
  //       Headers: {
  //         Authorization: `Bearer ${token}`,
  //       },
  //     });

  //     if (response.status !== 200) {
  //       // console.log("status not 200", response);
  //       return;
  //     }

  //     // console.log("created datasource library", response);

  //     localStorage.setItem("libraryIdDatasource", response.result.id);
  //     setDatasourceLibraryIds(response.result.id);
  //   }
  // }

  async function checkCollectionsBakrie() {
    try {
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
        // {
        //   name: "Document Audio",
        //   desc: "collection audio",
        // },
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
            // console.log("status not 200", response);
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
    } catch (error) {
      // console.log("checkCollectionsBakrie", error);
    }

    // setCollectionIds(libraryIds);
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

      // console.log("------", resp);

      if (resp.status !== 404) {
        // console.log("404", resp);
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

      // // console.log("dataSend", dataSend);

      try {
        const response = await call.call("client-dashboard", "create-bot", {
          URLQuery: {},
          URLParams: {},
          Params: dataSend,
          Headers: {
            Authorization: "Bearer " + token,
          },
        });

        // console.log("createBakrieBot", response);

        if (response.status !== 200) {
          // console.log(response);
          return;
        }

        localStorage.setItem("botIdStrapi", response.result.data.id);
      } catch (error) {
        // console.log(error);
      }
    } else {
      // console.log("no token");
      return;
    }
  }

  // async function createVoiceBot() {
  //   let token = localStorage.getItem("token");

  //   if (token !== null) {
  //     const resp = await call.call("client-dashboard", "get-bot-by-name", {
  //       URLQuery: {},
  //       URLParams: {},
  //       Params: {
  //         name: "CSV Insights Bots",
  //       },
  //       Headers: {
  //         Authorization: "Bearer " + localStorage.getItem("token"),
  //       },
  //     });

  //     // console.log("------", resp);

  //     if (resp.status !== 404) {
  //       // console.log("bot already exists", resp);
  //       // localStorage.setItem("botIdStrapi", resp.result.id);
  //       // gettingData = false;
  //       return;
  //     }

  //     let dataSend = {
  //       bot_name: "CSV Insights Bots",
  //       description: "",
  //       ai_engine: "llama3.1:8b",
  //       library_id: $datasourceLibraryId,
  //       collection_id: "all",
  //       additional_prompt: "",
  //       need_auth: false,
  //       module_type: "insights",
  //     };

  //     // // console.log("dataSend", dataSend);

  //     try {
  //       const response = await call.call("client-dashboard", "create-bot", {
  //         URLQuery: {},
  //         URLParams: {},
  //         Params: dataSend,
  //         Headers: {
  //           Authorization: "Bearer " + token,
  //         },
  //       });

  //       // console.log("createVoiceBot", response);

  //       if (response.status !== 200) {
  //         // console.log(response);
  //         return;
  //       }

  //       // localStorage.setItem("botIdStrapi", response.result.data.id);
  //     } catch (error) {
  //       // console.log(error);
  //     }
  //   } else {
  //     // console.log("no token");
  //     return;
  //   }
  // }

  let coreModule = [
    {
      name_library: "Data Source List",
      name_collection: "My Data Source",
      bot_name: "CSV Insights Bots",
      module_type: "insights",
      additional_prompt: "",
    },
    {
      name_library: "Voice Analyzer",
      name_collection: "My Voice",
      bot_name: "Voice’s Assistant",
      module_type: "voice",
      additional_prompt:
        "You are bot assistant for Voice Analyzer. You will help me to analyze voice data.",
    },
    {
      name_library: "Handwriting Reader",
      name_collection: "My Handwriting",
      bot_name: "Handwriting’s Assistant",
      module_type: "handwriting",
      additional_prompt:
        "Kamu adalah bot yang menjawab pertanyaan user sesuai dengan dokumen yang ada",
    },
    {
      name_library: "API Connector",
      module_type: "apiconnector",
      name_collection: "My API Collection",
    },
  ];

  async function createCoreModule(
    libraryName,
    collectionName,
    botName,
    moduleType,
    additionalPrompt,
  ) {
    let token = localStorage.getItem("token");
    if (token) {
      const response = await call.call("client-dashboard", "find-create-library-by-name", {
        URLQuery: {},
        URLParams: {},
        Params: {
          library_name: libraryName,
          description: "This is a library for " + libraryName,
          module_type: moduleType,
        },
        Headers: {
          Authorization: `Bearer ${token}`,
        },
      });

      // console.log("------" + libraryName, response);

      if (response.status !== 200) {
        // console.log("status not 200", response);
        // console.log("11111111");
        return;
      }

      localStorage.setItem("libraryId" + moduleType, response.result.id);

      if (collectionName !== undefined) {
        // console.log("collectionName not null" + collectionName);
        const responseCollection = await call.call(
          "client-dashboard",
          "find-create-library-collection-by-name",
          {
            URLQuery: {},
            URLParams: {},
            Params: {
              parent_id: response.result.id,
              collection_name: collectionName,
              description: "This is a collection for " + libraryName,
            },
            Headers: {
              Authorization: `Bearer ${token}`,
            },
          },
        );

        localStorage.setItem("collectionId" + moduleType, responseCollection.result.id);
      } else {
        // console.log("collectionName null" + libraryName);
      }

      // console.log("checked library" + libraryName, response);
      // console.log("------");

      if (botName !== undefined) {
        // console.log("botName not null" + botName);
        const resp = await call.call("client-dashboard", "get-bot-by-name", {
          URLQuery: {},
          URLParams: {},
          Params: {
            name: botName,
          },
          Headers: {
            Authorization: "Bearer " + localStorage.getItem("token"),
          },
        });

        // console.log("------" + botName, resp);

        if (resp.status !== 404) {
          // console.log("bot already exists" + botName, resp);
          localStorage.setItem("botId" + moduleType, resp.result.attributes.bot_id);
          localStorage.setItem("botIdStrapi" + moduleType, resp.result.id);
          // gettingData = false;
          return;
        }

        let dataSend = {
          bot_name: botName,
          description: "",
          ai_engine: "llama3.1:8b",
          library_id: response.result.id,
          collection_id: "all",
          additional_prompt: additionalPrompt,
          need_auth: false,
          module_type: moduleType,
        };

        try {
          const response2 = await call.call("client-dashboard", "create-bot", {
            URLQuery: {},
            URLParams: {},
            Params: dataSend,
            Headers: {
              Authorization: "Bearer " + token,
            },
          });

          // console.log("createCustomBot" + botName, response2);
          localStorage.setItem("botId" + moduleType, response2.result.attributes.bot_id);

          if (response2.status !== 200) {
            // console.log(response2);
            return;
          }
        } catch (error) {
          // console.log(error);
        }
      } else {
        // console.log("botName null" + libraryName);
      }
    }
  }

  // async function createDemoBot() {
  //   let token = localStorage.getItem("token");

  //   if (token !== null) {
  //     const resp = await call.call("client-dashboard", "get-bot-by-name", {
  //       URLQuery: {},
  //       URLParams: {},
  //       Params: {
  //         name: "Try LLMs Bot",
  //       },
  //       Headers: {
  //         Authorization: "Bearer " + localStorage.getItem("token"),
  //       },
  //     });

  //     if (resp.status !== 404) {
  //       // console.log(resp);
  //       gettingData = false;
  //       return;
  //     }

  //     let dataSend = {
  //       bot_name: "Try LLMs Bot",
  //       description: "",
  //       ai_engine: "llama3.2:8b",
  //       library_id: parseInt(localStorage.getItem("libraryIdWebchatBakrie")),
  //       collection_id: "all",
  //       additional_prompt: "",
  //       need_auth: false,
  //     };

  //     // console.log("dataSend", dataSend);

  //     try {
  //       const response = await call.call("client-dashboard", "create-bot", {
  //         URLQuery: {},
  //         URLParams: {},
  //         Params: dataSend,
  //         Headers: {
  //           Authorization: "Bearer " + token,
  //         },
  //       });

  //       if (response.status !== 200) {
  //         // console.log(response);
  //         return;
  //       }
  //     } catch (error) {
  //       // console.log(error);
  //     }
  //   } else {
  //     // console.log("no token");
  //     return;
  //   }
  // }

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
        // console.log(response);
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
        sessionStorage.setItem("sidebar_title", "HexaCode.AI");
        brand = "HexaCode.AI";
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
        // console.log(response);
        return;
      }

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
        // console.log(response);
        return;
      }

      if (response.result.data.length > 0) {
        notifications.update((data) => {
          // Gabungkan data baru dengan data lama
          // const combined = [...response.result.data, ...data];

          // Hilangkan duplikat berdasarkan `id`
          // const unique = combined.filter(
          //   (item, index, self) =>
          //     index === self.findIndex((t) => t.attributes.title === item.attributes.title),
          // );

          return response.result.data;
        });

        // notifications.update((data) => [...response.result.data, ...data]);
      }
    } catch (error) {
      // console.log(error);
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
        // console.log(response);
        return;
      }

      if (response.result.data.length > 0) {
        notificationsUnread.update((data) => [...response.result.data, ...data]);
        notificationsNewCount.set(response.result.meta.pagination.total);
      }
    } catch (error) {
      // console.log(error);
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
        // console.log(response);
        return;
      }

      if (response.result.data.length > 0) {
        notificationsRead.update((data) => [...response.result.data, ...data]);
      }
    } catch (error) {
      // console.log(error);
    }
  }

  async function CheckBotLlm() {
    let token = localStorage.getItem("token");
    if (token) {
      try {
        const response = await call.call("client-dashboard", "find-create-library-llm", {
          URLQuery: {},
          URLParams: {},
          Params: {},
          Headers: {
            Authorization: "Bearer " + token,
          },
        });

        if (response.status !== 200) {
          // console.log(response);
          return;
        }

        // console.log("succes get library & colection:", response);

        if (
          response.result?.collection_id !== undefined &&
          response.result?.library_id !== undefined
        ) {
          let collectionId = response.result.collection_id;
          let libraryId = response.result.library_id;

          const response2 = await call.call("client-dashboard", "find-create-bot-llm", {
            URLQuery: {},
            URLParams: {},
            Params: {
              collection_id: collectionId.toString(),
              library_id: libraryId,
            },
            Headers: {
              Authorization: "Bearer " + token,
            },
          });

          if (response2.status !== 200) {
            // console.log(response2);
            return;
          }

          if (response2.result?.bot_id !== undefined) {
            localStorage.setItem("botLlmId", response2.result.bot_id);
          }
        }
      } catch (error) {
        // console.log("bot llm err:", error);
      }
    }
  }

  const checkPageComingSoon = () => {
    return pageComingSoon.some((path) => path === globalThis.location.pathname);
  };

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

    // if ($isMoodle) {
    //   // Example set cookie
    //   const cookies = "3581cd4947eef2b29950792f7e6389c8";
    //   setCookie("MoodleSession", cookies, 1);
    // }

    isComingSoon = checkPageComingSoon();

    setCall(call);

    await CheckSetup();
    if ($isBakrie) {
      await checkLibrariesBakrie();
      await checkCollectionsBakrie();
      await createBakrieBot();
    }
    // await getPermissions();
    if (localStorage.getItem("lang") !== null) {
      locale.set(localStorage.getItem("lang"));
    }

    // await checkLibrariesDatasouce();
    // await createVoiceBot();

    coreModule.forEach(async (module) => {
      // // console.log(1111);
      await createCoreModule(
        module.name_library,
        module.name_collection,
        module.bot_name,
        module.module_type,
        module.additional_prompt,
      );
    });

    await getModelsTrocr();
    await getIndustryId();
    await CheckBotLlm();
    // hideNavbar(); // Taqie - Function to hide navbar on load

    getNotifications();
    getNotificationsUnread();
    getNotificationsRead();

    connectWebSocket(localStorage.getItem("token"), isFrontend);

    // bakrie;
    if ($isBakrie) {
      if (!hideLayoutBakrie()) {
        const libraryIdWebchatBakrie = window.localStorage.getItem("libraryIdWebchatBakrie");
        let collectionList = isFrontend
          ? window.location.pathname.trim() === "/documents"
            ? "documents/libraries/libraries_collections/"
            : "libraries/libraries_collections/"
          : await call.getPage("client-dashboard", "document-library-collections");
        const redirectUrl = `${collectionList}?libraryId=${libraryIdWebchatBakrie}`;
        goto(redirectUrl);
      } else {
        hideBakrie = false;
      }
    }

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
    // triggerShowAI(localStorage.getItem("showAI") === "true" ? true : false);
    toggleSidebarLargeClosed(
      localStorage.getItem("sidebarLargeClosed") === "true" && globalThis.innerWidth > 1024
        ? true
        : false,
    );
    // if (localStorage.getItem("selectedDropdown") === "collections") {
    //   triggerSelectDropdown(
    //     localStorage.getItem("selectedDropdown") === "collections" ? "collections" : "general",
    //   );
    //   if (localStorage.getItem("selectedCollection") !== null) {
    //     triggerSelectCollection(JSON.parse(localStorage.getItem("selectedCollection")));
    //   }
    // }
  }

  $effect(() => {
    // console.log(activeUrl);
    if (activeUrl !== $page.url.pathname) {
      activeUrl = $page.url.pathname;
      isComingSoon = checkPageComingSoon();
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
  <!-- {:else if winInnerWidth > 0 && $callReady === true && ready === true && !hideBakrie} -->
  {#if hideLayout()}
    <div class="lg:flex">
      <div class="z-20">
        <Sidebar
          getModule={getModules}
          bind:name={userData.username}
          bind:email={userData.email}
          bind:brand
          bind:logoUrl
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
                  <Navbar bind:authority bind:hasNewNotif {activeUrl} {appUrl} {changeShowNotif} />
                </div>
              {/if}
              {#if isComingSoon}
                <PageComingSoon />
              {:else}
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
              {/if}
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
