<script>
  import { onMount } from "svelte";
  import { _ } from "svelte-i18n";
  import { format } from "date-fns";
  import { page } from "$app/stores";
  import { goto } from "$app/navigation";
  import { Button, Modal, Card } from "flowbite-svelte";
  import { ExclamationCircleSolid, ArrowUpOutline } from "flowbite-svelte-icons";

  import {
    call,
    check,
    getIcons,
    triggerToast,
    triggerLoading,
    deleteCollectionAPIConnector,
    updateCollectionAPIConnector,
  } from "$lib";
  import {
    TabBar,
    Divider,
    InputRw,
    TabDetail,
    DataTable,
    StateCard,
    Typography,
    ToastContainer,
  } from "$lib";

  const MoveIcon = getIcons("move");
  const Copy01Icon = getIcons("copy01");
  const Save01Icon = getIcons("save01");
  const Trash01Icon = getIcons("trash01");
  const ChevronDownIcon = getIcons("chevronDown");
  const PenToolPlusIcon = getIcons("penToolPlus");
  const DotsVerticalIcon = getIcons("dotsVertical");
  const IllustrationEmpty = getIcons("illustrationEmpty");
  const IllustrationProgress = getIcons("illustrationProgress");

  const isFrontend = process.env.NODE_ENV === "development";

  let tmpId = $state(null);
  let textarea = $state(null);
  let collectionId = $state(null);
  let groupedFiles = $state(null);
  let tmpSelectedCollection = $state(null);

  let popupModal = $state(false);
  let updateModal = $state(false);
  let deleteModal = $state(false);
  let gettingDocs = $state(false);
  let dropdownMore = $state(false);
  let gettingDatas = $state(false);
  let isMobileView = $state(false);
  let buttonDisabled = $state(true);
  let checkboxDelete = $state(false);
  let nameIsAvailable = $state(true);
  let buttonDisabledCheckAction = $state(true);
  let dropdownAuthorizationType = $state(false);

  let libraryId = $state(0);

  let search = $state("");
  let tmpCollectionId = $state("");
  let tmpDocumentTitle = $state("");
  let tmpSelectedAuthorizationType = $state("No Auth");

  let files = $state([]);
  let collectionList = $state([]);
  let collectionListDatas = $state([]);
  let childCollectionListDatas = $state([]);
  let apisData = $state([
    // {
    //   id: 1,
    //   // checkbox: {
    //   //   checked: false,
    //   // },
    //   name: "GetUserAPI",
    //   method: {
    //     type: "success",
    //     name: "GET",
    //   },
    //   apiLink: "https://hexacode.ai/api/character?name",
    //   dateUploaded: "Jan 4, 2022",
    //   actions: {
    //     objects: [
    //       {
    //         href: "#",
    //         name: "View",
    //       },
    //     ],
    //   },
    // },
  ]);

  let pagination = $state({ page: 1, pageSize: 20, pageCount: 1, total: 0 });
  let paginationCollection = $state({ page: 1, pageSize: 20, pageCount: 1, total: 0 });
  let tabDatas = $state({
    selected: "My API",
    data: [
      {
        name: "My API",
      },
      {
        name: "Authorization",
      },
      {
        name: "Information",
      },
    ],
  });
  let tabDatasCollections = $state({
    selected: "Collections",
    data: [
      {
        name: "Collections",
      },
    ],
  });
  let authorizationTypeList = $state({
    selected: "No Auth",
    data: [
      { name: "No Auth" },
      { name: "Basic Auth" },
      { name: "Bearer Token" },
    ],
  });
  let authBasicAuth = $state({
    username: "",
    password: "",
  });
  let authBearerToken = $state({
    token: "",
  });
  let collectionDatas = $state({
    parentId: 0,
    title: "",
    description: "",
    createdAt: "",
  });

  const handleGetCollection = async () => {
    let token = localStorage.getItem("token");

    if (token !== null) {
      let response = await $call.call("client-dashboard", "get-collection-api", {
        URLQuery: {},
        URLParams: {
          id: collectionId,
        },
        Params: {},
        Headers: {
          Authorization: "Bearer " + token,
        },
      });

      // console.log("collection", response);

      collectionDatas = {
        title: response.result.data.attributes.name,
        description: response.result.data.attributes?.meta?.description || "",
        createdAt: format(new Date(response.result.data.attributes.createdAt), "dd-MM-yyyy"),
        parentId: response.result.data.attributes.parent_id,
      };

      authorizationTypeList.selected = response.result.data.attributes.auth_type || "No Auth";
      authBasicAuth.username = response.result.data.attributes.username || "";
      authBasicAuth.password = response.result.data.attributes.password || "";
      authBearerToken.token = response.result.data.attributes.token || "";
    }
  };

  const getDocuments = async () => {
    let token = localStorage.getItem("token");

    if (token !== null) {
      gettingDocs = true;

      if (search !== "") {
        pagination.page = 1;
      }

      try {
        let response = await $call.call("client-dashboard", "get-apis", {
          URLQuery: {},
          URLParams: {},
          Params: {
            page: pagination.page,
            page_size: pagination.pageSize,
            parent_id: parseInt(collectionId),
            search: search,
          },
          Headers: {
            Authorization: "Bearer " + token,
          },
        });

        collectionListDatas = [];

        if (response.status !== 200) {
          console.log(response);
          gettingDocs = false;
          if (response.status === 400) {
            console.error("Invalid parameters:", response);
            triggerToast($_("invalidParameters"), true, "error");
          } else if (response.status === 401) {
            console.error("Unauthorized access:", response);
            triggerToast($_("unauthorized"), true, "error");
          } else if (response.status === 500) {
            console.error("Server error:", response);
            triggerToast($_("serverError"), true, "error");
          } else {
            console.error("Unexpected status code:", response);
            triggerToast(`${$_("unexpectedError")}: ${response.status}`, true, "error");
          }
          return;
        }

        // console.log("documets", response.result.data);

        pagination = response.result.meta.pagination;

        let link = isFrontend
          ? "/api_connector/collections/api/detail"
          : await $call.getPage("client-dashboard", "api-connector-collections-api-detail");

        let promise = response.result.data.map(async (collection) => {
          const params = new URLSearchParams(window.location.search);

          params.delete("documentId");

          const queryParams = {};
          for (const [key, value] of params.entries()) {
            queryParams[key] = value;
          }
          const queryString = new URLSearchParams(queryParams).toString();

          return {
            id: collection.id,
            // checkbox: {
            //   checked: false,
            // },
            name: collection.attributes.name,
            method: {
              type: "success",
              name: collection.attributes.method,
            },
            apiLink: collection.attributes.url,
            dateUploaded: format(new Date(collection.attributes.createdAt), "dd-MM-yyyy"),
            actions: {
              objects: [
                {
                  href: link + "?apiId=" + collection.id + "&" + queryString,
                  name: "View",
                },
              ],
            },
            // href: link + "?documentId=" + collection.id + "&" + queryString,
          };
        });

        apisData = await Promise.all(promise);

        // console.log("collectionListDatas", collectionListDatas);
      } catch (error) {
        console.error("Error during documents get:", error);
        gettingDocs = false;
        triggerToast($_("failedGetDocuments"), true, "error");
      }
      gettingDocs = false;
    } else {
      triggerToast($_("unauthorized"), true, "error");
      gettingDocs = false;
      return;
    }
  };

  const deleteDocument = async (id) => {
    triggerLoading(true);
    popupModal = false;

    try {
      let token = localStorage.getItem("token");

      if (token === null) {
        const error = new Error("Unauthorized");
        error.status = 401;
        throw error;
      }

      const response = await $call.call("client-dashboard", "delete-library-document", {
        URLQuery: {},
        URLParams: {
          id: id ? id : tmpId.toString(),
        },
        Params: {},
        Headers: {
          Authorization: "Bearer " + token,
        },
      });

      if (response?.status !== 200) {
        const error = new Error("Failed to delete document");
        error.response = response;
        error.status = response.status;
        throw error;
      }
    } catch (error) {
      errorHandling($_, error, "delete");
    } finally {
      if (!id) {
        tmpId = null;
        triggerToast($_("successDeleteDocument"), true);
        getDocuments();
        triggerLoading(false);
      }
    }
  };

  const deleteDocuments = async () => {
    let datasDelete = collectionListDatas.filter((data) => data.checked === true);

    for (const data of datasDelete) {
      await deleteDocument(data.id);
    }

    triggerToast($_("successDeleteDocuments"), true);
    getDocuments();
    triggerLoading(false);
  };

  const updateDocument = async () => {
    let datasUpdate = collectionListDatas.filter((data) => data.checked === true);
    updateModal = false;
    triggerLoading(true);

    for (const data of datasUpdate) {
      const id = data.id;
      let token = localStorage.getItem("token");
      if (token !== null) {
        let response;
        response = await $call.call("client-dashboard", "update-library-document-parent", {
          URLQuery: {},
          URLParams: {},
          Params: {
            document_id: parseInt(id),
            parent_id: parseInt(tmpSelectedCollection.id),
          },
          Headers: {
            Authorization: "Bearer " + token,
          },
        });
        if (response.status !== 200) {
          console.log(response);
          triggerLoading(false);
          if (response.status === 400) {
            console.error("Invalid parameters:", response);
            triggerToast($_("invalidParameters"), true, "error");
          } else if (response.status === 401) {
            console.error("Unauthorized access:", response);
            triggerToast($_("unauthorized"), true, "error");
          } else if (response.status === 500) {
            console.error("Server error:", response);
            triggerToast($_("serverError"), true, "error");
          } else {
            console.error("Unexpected status code:", response);
            triggerToast(`${$_("unexpectedError")}: ${response.status}`, true, "error");
          }
          return;
        }

        triggerLoading(false);
        // console.log(response);
        triggerToast("successUpdateDocument", true);
        getDocuments();
      } else {
        triggerToast($_("unauthorized"), true, "error");
        return;
      }
    }
  };

  const getCollections = async (isChild = false, subCollections = false) => {
    let token = localStorage.getItem("token");

    if (token !== null) {
      gettingDatas = true;
      const aggregatedData = [];

      paginationCollection.page = 1;

      do {
        let response = await $call.call("client-dashboard", "get-apis", {
          URLQuery: {},
          URLParams: {},
          Params: {
            page: paginationCollection.page,
            page_size: paginationCollection.pageSize,
            parent_id: parseInt(
              isChild ? collectionId : subCollections ? tmpSelectedCollection.id : libraryId,
            ),
          },
          Headers: {
            Authorization: "Bearer " + token,
          },
        });

        if (response.status !== 200) {
          console.log(response);
          gettingDatas = false;
          if (response.status === 400) {
            console.error("Invalid parameters:", response);
            triggerToast($_("invalidParameters"), true, "error");
          } else if (response.status === 401) {
            console.error("Unauthorized access:", response);
            triggerToast($_("unauthorized"), true, "error");
          } else if (response.status === 500) {
            console.error("Server error:", response);
            triggerToast($_("serverError"), true, "error");
          } else {
            console.error("Unexpected status code:", response);
            triggerToast(`${$_("unexpectedError")}: ${response.status}`, true, "error");
          }
          return;
        }

        aggregatedData.push(...response.result.data);

        paginationCollection = response.result.meta.pagination;

        paginationCollection.page++;
      } while (paginationCollection.page <= paginationCollection.pageCount);

      let link;

      if (isChild) {
        link = isFrontend
          ? "/handwriting/collections/documents/sub_documents"
          : await $call.getPage("client-dashboard", "handwriting-sub-documents");
      } else {
        link = isFrontend
          ? "/handwriting/collections/documents"
          : await $call.getPage("client-dashboard", "handwriting-documents");
      }

      const params = new URLSearchParams(window.location.search);
      params.delete("collectionId");

      const queryParams = {};
      for (const [key, value] of params.entries()) {
        queryParams[key] = value;
      }
      const queryString = new URLSearchParams(queryParams).toString();

      const promise = aggregatedData.map(async (collection) => ({
        id: collection.id,
        collectionName: collection.attributes.name,
        description: collection.attributes?.meta?.description || "",
        type: collection.attributes.type,
        href:
          link +
          "?subCollectionId=" +
          collection.id +
          "&" +
          "collectionId=" +
          collectionId +
          "&" +
          queryString,
      }));

      if (isChild) {
        childCollectionList = await Promise.all(promise);
        // console.log("ischild", childCollectionList);
      } else if (subCollections) {
        childCollectionListDatas = await Promise.all(promise);
        // console.log("subcollection", childCollectionListDatas);
        if (childCollectionListDatas.length > 0) {
          if (tabDatasCollections.data.length < 2) {
            tabDatasCollections.data = [
              ...tabDatasCollections.data,
              {
                name: "Sub Collections",
              },
            ];
          }
        } else {
          tabDatasCollections.data = tabDatasCollections.data.filter(
            (item) => item.name !== "Sub Collections",
          );
        }
      } else {
        // console.log("collections", collectionList);
        collectionList = await Promise.all(promise);
      }
      gettingDatas = false;
    } else {
      triggerToast($_("unauthorized"), true, "error");
      gettingDatas = false;
      return;
    }
  };

  const handleDeleteCollection = async () => {
    deleteModal = false;
    gettingDatas = true;

    const result = await deleteCollectionAPIConnector($_, $call, tmpCollectionId);

    if (result) {
      goto(
        isFrontend
          ? "api_connector/collections"
          : await $call.getPage("client-dashboard", "api-connector-collections"),
      );
    } else {
      gettingDatas = false;
    }
  };

  const handleUpdateCollection = async () => {
    gettingDatas = true;

    console.log("test", authorizationTypeList.selected);

    const result = await updateCollectionAPIConnector(
      $_,
      $call,
      collectionDatas.parentId,
      collectionId,
      collectionDatas.title,
      collectionDatas.description,
      authorizationTypeList.selected,
      authBasicAuth.username,
      authBasicAuth.password,
      authBearerToken.token,
    );

    if (result) {
      handleGetCollection();
    } else {
      gettingDatas = false;
    }
  };

  const checkSingle = (index) => {
    collectionListDatas = collectionListDatas.map((data, i) => {
      if (i === index) {
        return {
          ...data,
          checked: !data.checked,
        };
      }
      return data;
    });
    // console.log("cek single", collectionListDatas);
  };

  const checkAll = () => {
    let findCheckAll = true;
    for (const data of collectionListDatas) {
      if (data.checked === false) {
        findCheckAll = false;
        break;
      }
    }

    if (findCheckAll) {
      collectionListDatas = collectionListDatas.map((data) => {
        return {
          ...data,
          checked: false,
        };
      });
    } else {
      collectionListDatas = collectionListDatas.map((data) => {
        return {
          ...data,
          checked: true,
        };
      });
    }
  };

  const nextPage = (evt) => {
    if (pagination.pageCount > pagination.page) pagination.page = pagination.page + 1;
    getDocuments();
    // console.log("next page");
  };

  const prevPage = (evt) => {
    if (pagination.page > 1) pagination.page = pagination.page - 1;
    getDocuments();
    // console.log("prev page");
  };

  const autosize = () => {
    textarea.style.height = "auto"; // Reset height
    textarea.style.height = `${textarea.scrollHeight}px`; // Set height to scroll height

    textarea.addEventListener("input", (e) => {
      e.target.style.height = "auto";
      e.target.style.height = `${e.target.scrollHeight}px`;
    });
  };

  const selectCollection = (item, subcollection = false) => {
    tmpSelectedCollection = item;
    if (!subcollection) {
      getCollections(false, true);
    }
  };

  const checkMobileView = () => {
    isMobileView = window.innerWidth <= 768;
  };

  const scrollToTop = () => {
    // console.log("Scrolling to top using scrollIntoView");

    // Use an existing element at the top of the page
    const topElement = document.querySelector(".top-element");
    if (topElement) {
      topElement.scrollIntoView({ behavior: "smooth" });
    } else {
      console.error("Top element not found");
    }
  };

  const handleDropdownFocusLoss = ({ relatedTarget, currentTarget }) => {
    if (relatedTarget instanceof HTMLElement && currentTarget.contains(relatedTarget)) return;
    dropdownMore = false;
    dropdownAuthorizationType = false;
  };

  onMount(() => {
    if (check("documentAi.documentLibrary:view") === false) {
      window.history.back();
      return;
    }

    const unsubscribe = page.subscribe(async ($page) => {
      collectionId = $page.url.searchParams.get("collectionId");
      if (collectionId === null || collectionId === undefined) {
        goto(await $call.getPage("client-dashboard", "api-connector"));
        return;
      }

      await handleGetCollection();
      await getDocuments();
    });

    checkMobileView();
    window.addEventListener("resize", checkMobileView);

    return () => {
      unsubscribe();
      window.removeEventListener("resize", checkMobileView);
    };
  });

  $effect(() => {
    if (textarea) {
      autosize();
    }
  });

  $effect(() => {
    if (groupedFiles != null) {
      buttonDisabled = false;
    } else {
      buttonDisabled = true;
    }
  });

  $effect(() => {
    buttonDisabledCheckAction = !collectionListDatas.some((data) => data.checked === true);
    if (files.length > 0) {
      hasDocOrDocx = files.some((file) => file.ext === ".doc" || file.ext === ".docx");
      buttonDisabled =
        files.length === 0 || files.length === 1
          ? tmpDocumentTitle === "" || !nameIsAvailable
          : false;
    } else if (groupedFiles != null) {
      buttonDisabled = false;
    } else {
      buttonDisabled = true;
    }
  });
</script>

<div class="top-element"></div>
<!-- Add this element at the top of the page -->
<div class="flex dark:bg-primary-dark overflow-auto md:overflow-y-hidden md:hover:overflow-y-auto">
  {#if isMobileView}
    <button
      class="fixed bottom-5 right-5 z-40 bg-brand-primary-light dark:bg-brand-primary-dark text-white rounded-full w-12 h-12 flex items-center justify-center shadow-lg"
      onclick={scrollToTop}
    >
      <ArrowUpOutline class="w-6 h-6" />
    </button>
  {/if}
  <div class="w-full pb-10 min-h-[30rem]">
    <div class="flex flex-col m-4">
      <div class="flex flex-col">
        <!-- User Information -->
        <div class="flex flex-col md:flex-row md:items-center justify-between gap-2 mb-4">
          <div class="flex flex-col gap-2">
            <div class="flex gap-4">
              <Typography
                preset="Display sm/Semibold"
                color="primary"
                text={collectionDatas.title}
                propClass="fluid-xl line-clamp-2"
              />
              <button
                onclick={() => {
                  navigator.clipboard
                    .writeText(collectionDatas.title)
                    .then(() => {
                      // console.log("Text copied to clipboard");
                      triggerToast("copyToClipboard");
                    })
                    .catch((err) => {
                      console.error("Failed to copy: ", err);
                    });
                }}
                class="p-2 border-2 rounded-lg border-secondary-light dark:border-secondary-dark"
              >
                <Copy01Icon
                  svgClass="w-5 h-5 text-fg-secondary-light dark:text-fg-secondary-dark"
                />
              </button>
            </div>
            {#if collectionDatas?.description !== "" && collectionDatas?.description != null && collectionDatas?.description !== undefined}
              <Typography
                preset="Text md/Medium"
                color="tertiary"
                text={collectionDatas.description}
              />
            {/if}
            <Typography
              preset="Text sm/Medium"
              color="quarternary"
              text={`${$_("createOn")} ${collectionDatas.createdAt}`}
            />
          </div>

          <div class="flex relative gap-3" onfocusout={handleDropdownFocusLoss}>
            {#if !buttonDisabledCheckAction}
              <Button
                on:click={() => {
                  popupModal = true;
                  checkboxDelete = true;
                  // id
                }}
                disabled={!check("documentAi.documentLibrary.document:delete")}
                class="p-2 md:px-3 md:py-2 gap-2 text-error-primary-light dark:text-error-primary-dark hover:text-error-primary-light dark:hover:text-error-primary-dark border-secondary-light dark:border-secondary-dark hover:bg-secondary-light hover:dark:bg-secondary-dark"
                outline
              >
                <Trash01Icon svgClass="w-5 h-5" />
                <div class="hidden md:block">{$_("delete")}</div>
              </Button>
              <Button
                on:click={() => {
                  updateModal = true;
                  checkboxUpdate = true;
                  paginationCollection.page = 1;
                  // id
                }}
                disabled={!check("documentAi.documentLibrary.document:update")}
                class="p-2 md:px-3 md:py-2 gap-2 text-fg-secondary-light dark:text-fg-secondary-dark border-secondary-light dark:border-secondary-dark hover:bg-secondary-light hover:dark:bg-secondary-dark hover:text-fg-secondary-light dark:hover:text-fg-secondary-dark"
                outline
              >
                <MoveIcon svgClass="w-5 h-5" />
                <div class="hidden md:block">{$_("move")}</div>
              </Button>
            {/if}
            <Button
              outline
              class="min-w-max gap-1.5"
              disabled={!check("documentAi.documentLibrary.document:create")}
              on:click={async () => {
                let link = isFrontend
                  ? "/api_connector/collections/api/create" + `?collectionId=${collectionId}`
                  : (await $call.getPage(
                      "client-dashboard",
                      "api-connector-collections-api-create",
                    )) + `?collectionId=${collectionId}`;
                goto(link);
              }}
            >
              <PenToolPlusIcon svgClass="w-5 h-5 text-quinary-light dark:text-quinary-dark" />
              <div class="hidden md:block">{$_("addNewApi")}</div>
            </Button>
            {#if check("documentAi.documentLibrary.collection:delete")}
              <div onfocusout={handleDropdownFocusLoss} class=" flex items-center relative">
                <Button on:click={() => (dropdownMore = !dropdownMore)} color="light" class="p-2.5">
                  <DotsVerticalIcon svgClass="w-5 h-5 text-quinary-light dark:text-quinary-dark" />
                </Button>
                <div
                  class={`z-10 ${dropdownMore ? "absolute top-10 right-0" : "hidden"} bg-white divide-y divide-gray-100 rounded-lg shadow w-max dark:bg-gray-700`}
                >
                  <ul class="py-2 grow text-sm">
                    {#if check("documentAi.documentLibrary.collection:delete")}
                      <button
                        class="flex p-2 w-full text-secondary-light dark:text-secondary-dark hover:text-secondary-hover-light dark:hover:text-secondary-hover-dark hover:bg-primary-hover-light dark:hover:bg-primary-hover-dark"
                        onclick={() => {
                          dropdownMore = !dropdownMore;
                          deleteModal = true;
                        }}
                      >
                        {$_("delete")}
                      </button>
                    {/if}
                  </ul>
                </div>
              </div>
            {/if}
          </div>
        </div>

        <TabBar
          {tabDatas}
          changeTabbar={(tabName) => {
            tabDatas.selected = tabName;
          }}
        />
        <Divider class="mb-4" />

        {#if tabDatas.selected === "My API"}
          {#if apisData.length > 0}
            <DataTable
              bind:datas={apisData}
              next={nextPage}
              previous={prevPage}
              {checkAll}
              {checkSingle}
            />
          {:else if gettingDatas || gettingDocs}
            <StateCard
              illustration={IllustrationProgress}
              description={$_("fetchingDataPleaseWait")}
            />
          {:else}
            <StateCard illustration={IllustrationEmpty} description={$_("youDontHaveAnyDataYet")} />
          {/if}
        {:else if tabDatas.selected === "Authorization"}
          <div class="grid grid-cols-12 mt-4 gap-4">
            <Typography
              preset="Text md/Semibold"
              color="secondary"
              text={$_("authType")}
              propClass="col-span-3"
            />
            <div
              onfocusout={handleDropdownFocusLoss}
              class=" flex flex-col gap-2 relative col-span-9"
            >
              <Typography
                preset="Text sm/Medium"
                color="secondary"
                text={$_("authorizationType")}
              />
              <Button
                on:click={() => (dropdownAuthorizationType = !dropdownAuthorizationType)}
                class="flex justify-between bg-primary-light dark:bg-primary-dark border border-primary-light dark:border-primary-dark hover:bg-primary-light dark:hover:bg-primary-dark hover:border-primary-light dark:hover:border-primary-dark focus:ring-0"
              >
                <Typography preset="Text md/Reguler" text={authorizationTypeList.selected} />
                <ChevronDownIcon svgClass="w-5 h-5 text-quinary-light dark:text-quinary-dark" />
              </Button>
              <Typography
                preset="Text sm/Reguler"
                color="tertiary"
                text={$_("theAuthorizationHeaderWillBeAutomaticallyGeneratedWhenYouSendTheRequest")}
              />
              <div
                class={`z-10 ${dropdownAuthorizationType ? "absolute top-[5rem] left-0" : "hidden"} bg-white divide-y divide-gray-100 rounded-lg shadow w-max dark:bg-gray-700`}
              >
                <ul class="py-2 grow text-sm">
                  {#each authorizationTypeList.data as auth}
                    <button
                      class="flex p-2 w-full text-secondary-light dark:text-secondary-dark hover:text-secondary-hover-light dark:hover:text-secondary-hover-dark hover:bg-primary-hover-light dark:hover:bg-primary-hover-dark"
                      onclick={() => {
                        dropdownAuthorizationType = !dropdownAuthorizationType;
                        authorizationTypeList.selected = auth.name;
                        authBasicAuth.username = "";
                        authBasicAuth.password = "";
                        authBearerToken.token = "";
                      }}
                    >
                      {$_(auth.name)}
                    </button>
                  {/each}
                </ul>
              </div>
            </div>

            {#if authorizationTypeList.selected === "Basic Auth"}
              <Divider className="col-span-12 my-4" />
              <Typography
                preset="Text md/Semibold"
                color="secondary"
                text={$_("token")}
                propClass="col-span-3"
              />
              <div class="col-span-9 flex flex-col gap-2">
                <Typography preset="Text sm/Medium" color="secondary" text={$_("username")} />
                <InputRw bind:value={authBasicAuth.username} />
              </div>
              <div class="col-span-3"></div>
              <div class="col-span-9 flex flex-col gap-2">
                <Typography preset="Text sm/Medium" color="secondary" text={$_("password")} />
                <InputRw bind:value={authBasicAuth.password} type="password" />
              </div>
            {:else if authorizationTypeList.selected === "Bearer Token"}
              <Divider className="col-span-12 my-4" />
              <Typography
                preset="Text md/Semibold"
                color="secondary"
                text={$_("token")}
                propClass="col-span-3"
              />
              <div class="col-span-9 flex flex-col gap-2">
                <Typography preset="Text sm/Medium" color="secondary" text={$_("token")} />
                <InputRw bind:value={authBearerToken.token} />
              </div>
            {/if}
            <div class="col-span-12 flex justify-end">
              <Button
                class="min-w-max gap-1.5"
                disabled={authorizationTypeList.selected === "No Auth" ||
                authorizationTypeList.selected === "Inherit Auth from Parent"
                  ? authorizationTypeList.selected === tmpSelectedAuthorizationType
                  : authorizationTypeList.selected === "Basic Auth"
                    ? authBasicAuth.username === "" && authBasicAuth.password === ""
                    : authorizationTypeList.selected === "Bearer Token"
                      ? authBearerToken.token === ""
                      : false}
                on:click={async () => {
                  handleUpdateCollection();
                }}
              >
                <Save01Icon svgClass="w-5 h-5 text-fg-primary-light dark:text-fg-primary-dark" />
                <div class="hidden md:block">{$_("saveChanges")}</div>
              </Button>
            </div>
          </div>
        {:else if tabDatas.selected === "Information"}
          <div class="grid grid-cols-12 mt-4 gap-4">
            <Typography
              preset="Text md/Semibold"
              color="secondary"
              text={$_("generalInformation")}
              propClass="col-span-3"
            />
            <div class="flex flex-col col-span-9 gap-2">
              <Typography preset="Text sm/Medium" color="secondary" text={$_("collectionTitle")} />
              <InputRw value={collectionDatas.title} />
            </div>
            <div class="col-span-3"></div>
            <div class="flex flex-col col-span-9 gap-2">
              <Typography preset="Text sm/Medium" color="secondary" text={$_("description")} />
              <InputRw value={collectionDatas.description} />
            </div>
          </div>
        {:else}
          Hello
        {/if}
      </div>
    </div>
  </div>
  <ToastContainer />
</div>

<Modal bind:open={popupModal} size="xs" autoclose>
  <div class="text-center">
    <ExclamationCircleSolid class="mx-auto mb-4 text-red-600 w-12 h-12" />
    <h3 class="mb-5 text-lg font-normal text-gray-500 dark:text-gray-400">
      {$_("areYouSureYouWantToDelete")}
      {checkboxDelete ? $_("theseDocuments") : $_("thisDocument")}
    </h3>
    <Button onclick={checkboxDelete ? deleteDocuments : deleteDocument} color="red" class="me-2"
      >{$_("yesImSure")}</Button
    >
    <Button color="alternative">{$_("noCancel")}</Button>
  </div>
</Modal>

<Modal bind:open={updateModal}>
  <svelte:fragment slot="header">
    <div class="flex flex-col">
      <Typography preset="Text lg/Semibold" text={$_("moveDocument")} color="primary" />
      <Typography
        preset="Text sm/Regular"
        text={$_("whereDoYouWantToMoveThisDocument")}
        color="tertiary"
      />
    </div>
  </svelte:fragment>

  <div class="flex flex-col">
    <div class="bg-white dark:bg-gray-800 sticky -top-5">
      <TabDetail
        spaceTop="mt-0"
        {tabDatasCollections}
        on:tabChange={(e) => (tabDatasCollections.selected = e.detail)}
      />
    </div>
    {#if tabDatasCollections.selected === "Collections"}
      <div class="flex flex-col justify-center mt-3 gap-3">
        <div class="flex flex-col gap-6">
          {#each collectionList as item, index (item)}
            <Card
              class="{tmpSelectedCollection === item
                ? 'border-brand-solid-light dark:border-brand-solid-dark ring-brand-solid-light dark:ring-brand-solid-dark'
                : ''}bg-hc-graylight-50 hover:bg-hc-brand-50 cursor-pointer gap-3"
              size="xl"
              on:click={() => selectCollection(item)}
            >
              <div class="flex flex-col gap-0.5">
                <Typography
                  preset="Text md/Semibold"
                  color={tmpSelectedCollection === item ? "brand-primary" : "primary"}
                  text={item.collectionName}
                />
                {#if item.description !== "" && item.description != null && item.description !== undefined}
                  <Typography
                    preset="Text xs/Regular"
                    color={tmpSelectedCollection === item ? "brand-secondary" : "secondary"}
                    text={item.description}
                  />
                {/if}
              </div>
            </Card>
          {/each}
        </div>
      </div>
    {:else if tabDatasCollections.selected === "Sub Collections"}
      <div class="flex flex-col justify-center mt-3 gap-3">
        <div class="flex flex-col gap-6">
          {#each childCollectionListDatas as item, index (item)}
            <Card
              class="{tmpSelectedCollection === item
                ? 'border-brand-solid-light dark:border-brand-solid-dark ring-brand-solid-light dark:ring-brand-solid-dark'
                : ''}bg-hc-graylight-50 hover:bg-hc-brand-50 cursor-pointer gap-3"
              size="xl"
              on:click={() => selectCollection(item, true)}
            >
              <div class="flex flex-col gap-0.5">
                <Typography
                  preset="Text md/Semibold"
                  color={tmpSelectedCollection === item ? "brand-primary" : "primary"}
                  text={item.collectionName}
                />
                {#if item.description !== "" && item.description != null && item.description !== undefined}
                  <Typography
                    preset="Text xs/Regular"
                    color={tmpSelectedCollection === item ? "brand-secondary" : "secondary"}
                    text={item.description}
                  />
                {/if}
              </div>
            </Card>
          {/each}
        </div>
      </div>
    {/if}
  </div>
  <svelte:fragment slot="footer">
    <Button on:click={updateDocument} class="w-full" disabled={!tmpSelectedCollection}
      >{$_("move")}</Button
    >
  </svelte:fragment>
</Modal>

<Modal bind:open={deleteModal} size="xs" autoclose>
  <div class="text-center">
    <ExclamationCircleSolid class="mx-auto mb-4 text-red-600 w-12 h-12" />
    <h3 class="mb-5 text-lg font-normal text-gray-500 dark:text-gray-400">
      {$_("areYouSureYouWantToDelete")}
      {$_("collection")}
      "{collectionDatas.title}"?
    </h3>
    <Button
      onclick={() => {
        handleDeleteCollection();
        deleteModal = false;
      }}
      color="red"
      class="me-2"
    >
      {$_("yesImSure")}
    </Button>
    <Button color="alternative" onclick={() => (deleteModal = false)}>
      {$_("noCancel")}
    </Button>
  </div>
</Modal>
