<script>
  import { Divider, Typography, StateCard, Icon, call, SelectedCard, datasourceLibraryId, Pagination } from "$lib";
  import { _ } from "svelte-i18n";
  import Markdown from "svelte-exmarkdown";
  import { gfmPlugin } from "svelte-exmarkdown/gfm";
  import UserIcon from "$lib/icon/user-01.svg";
  import File06Icon from "$lib/icon/file-06.svelte";
  import SendIcon from "$lib/icon/send-01.svg";
  import ChevronDownIcon from "$lib/icon/chevron-down.svelte";
  import DataTableListCollection from "$lib/components/DataTableListCollection.svelte";
  import IllustrationEmpty from "$lib/icon/illustration_empty.svelte";
  import IllustrationProgress from "$lib/icon/illustration_progress.svelte";
  import {
    Button,
    Modal,
    Table,
    TableBody,
    TableBodyCell,
    TableBodyRow,
    TableHead,
    TableHeadCell,
  } from "flowbite-svelte";
  import { onMount } from "svelte";
  import { id } from "date-fns/locale";

  const isFrontend = process.env.NODE_ENV === "development";

  const plugins = [gfmPlugin()];

  let messageInput = $state("");
  let textarea = $state(null);
  let modalOpen = $state(false);
  let allDataCollection = $state([]);
  let currentArrayData = $state([]);
  let datasChecked = $state([
    {
      id: 1,
      checked: false,
    },
    {
      id: 2,
      checked: false,
    },
  ]);
  let libraryId = $state(null);
  let modeChat = $state("Code");
  let tmpDatasChecked = $state({});
  let configs = $state({ hideSearch: true });
  let searchTerm = $state("");
  let gettingDatas = $state(false);
  let innerWidth = $state(0);
  let winInnerWidth = $derived(innerWidth);
  let isSelectedLibrary = $state(false);
  let collections = $state([]);
  let initCollections = $state([]);
  let selectedCollection = $state("");
  let selectedDocument = $state("");
  let configsChecked = $state({});
  let tmpSelectedCollection = $state("");
  let oldSelectedCollection = $state("");
  let search = $state(false);
  let checkAll = $state(false);
  let selectCsv = $state(true);

  let paginationModal = $state({ page: 1, pageSize: 3, pageCount: 0, total: 0 });
  let pagination = $state({ page: 1, pageSize: 10, pageCount: 0, total: 0 });
  let paginationCollection = $state({ page: 1, pageSize: 20, pageCount: 0, total: 0 });
  let paginationNew = $state({ page: 1, pageSize: 24, pageCount: 1, total: 0 });
  let csvMeta = $state({ page: 1, pageSize: 10, pageCount: 2, total: 0 });

  let dataConversationChats = $state([]);

  let documentPath = $state("");
  let documentTitle = $state("CSV Insight");

  let ws;
  let current_session_user = $state("");
  let current_session_id = $state("");
  let image = $state("");
  let message = $state("");
  let isTyping = $state(false);
  let fileDetail = $state({});

  function connectWebSocket() {
    var base_url =
      (location.protocol == "https:" ? "wss:" : "ws:") +
      "//" +
      location.hostname +
      (location.port == "" ? "" : ":" + location.port);

    // if (need_auth) {
    if (isFrontend) {
      ws = new WebSocket(
        `http://192.168.18.15/insight/c/${encodeURIComponent(localStorage.getItem("token"))}`,
      );
    } else {
      ws = new WebSocket(
        `${base_url}/insight/c/${encodeURIComponent(localStorage.getItem("token"))}`,
      );
    }
    // } else {
    //   ws = new WebSocket(`${base_url}/hexa-ai-rag/pub/c/${localStorage.getItem("botId")}`);
    // }

    ws.onopen = function (event) {
      console.log("WebSocket connected", event);
      Joining();
    };

    ws.onerror = function (error) {
      console.error("WebSocket error occurred:", error);
    };

    ws.onmessage = function (event) {
      console.log("Received WebSocket message:", event.data);
      const data = JSON.parse(event.data);

      if (data?.file_detail && data.file_detail.name !== "" && data.file_detail.local_path !== "") {
        fileDetail = data.file_detail;
        // console.log("fileDetail", fileDetail);
      }

      current_session_user = data.session_user;

      if (data.session_id && data.session_id !== "dummy") {
        current_session_id = data.session_id;
      }

      isTyping = data.action.is_typing;

      if (data.type === "content") {
        message = data.content.message;
        image = data.content.image;

        // if (image){
        //   image = base64toBlob(image, "image/png");
        // }

        if (data.content.role === "assistant") {
          dataConversationChats = [
            ...dataConversationChats,
            {
              id: dataConversationChats.length + 1,
              message: message,
              role: data.content.role,
              image: image,
            }
          ]
        }


        console.log("masuk berapa");

        console.log("cekkk dataConversationChats", dataConversationChats);

        // console.log("data.content", data.content);
      }
    };

    ws.onclose = () => {
      console.log("WebSocket disconnected.");
      // Opsi: coba reconnect otomatis jika diperlukan
      setTimeout(() => {
        connectWebSocket();
      }, 3000);
    };
  }

  const SendMessage = (message) => {
    console.log("SendMessage", message);
    ws.send(JSON.stringify(message));
  };

  const Joining = (session_user) => {
    console.log("Joining", session_user);
    if (session_user !== null && session_user !== undefined) {
      console.log("Joining existing");
      SendMessage({
        type: "joining",
        joining: {
          session_user: session_user,
        },
      });
    } else {
      console.log("Joining new");
      SendMessage({
        type: "joining",
        // joining: {
        //   session_user: session_user,
        // },
      });
    }
  };

  // const Leaving = (session_user) => {
  //   SendMessage({
  //     type: "leaving",
  //     session_user: session_user,
  //   });
  // };

  const Select = (session_user, path) => {
    SendMessage({
      type: "select",
      session_user: session_user,
      local_path: path,
    });
  };

  const Message = (session_user, session_id, message) => {
    SendMessage({
      type: "content",
      session_id: session_id,
      session_user: session_user,
      file_detail: fileDetail,
      content: {
        role: "user",
        message: message,
      },
      mode: modeChat,
    });
  };

  const autosize = () => {
    textarea.addEventListener("input", (e) => {
      setTimeout(() => {
        textarea.style.height = "auto";
        let scHeight = e.target.scrollHeight + 2; // Tambahkan sedikit untuk padding
        textarea.style.height = `${scHeight}px`;
      }, 0);
    });
  };

  const resetTextarea = () => {
    textarea.style.height = "auto";
  };

  async function sendMessageInput() {
    dataConversationChats = [];
    dataConversationChats = [
      ...dataConversationChats,
      {
        id: dataConversationChats.length + 1,
        role: "user",
        message: messageInput,
      },
    ];

    console.log("dataConversationChats", dataConversationChats);
    Message(current_session_user, current_session_id, messageInput);
    messageInput = "";
    resetTextarea();
    console.log("messageInput", messageInput);
  }

  const handleKeyDown = (event) => {
    if (winInnerWidth > 768) {
      if (event.key === "Enter") {
        if (!event.shiftKey) {
          event.preventDefault();
          if (messageInput.trim() !== "") {
            onclickSend();
          }
        }
      }
    }
  };

  const onclickSend = () => {
    if (messageInput.trim() !== "") {
      sendMessageInput();
    }
  };

  async function next() {
    if (csvMeta.page >= csvMeta.pageCount) return;
    csvMeta.page += 1;
    getCsvData();
  }

  async function previous() {
    if (csvMeta.page > 1) {
      csvMeta.page -= 1;
      getCsvData();
    }
  }

  async function nextCollectionDoc() {
    if (pagination.page >= pagination.pageCount) return;
    pagination.page += 1;
    getCollectionSub(selectedCollection);
  }

  async function previousCollectionDoc() {
    if (pagination.page > 1) {
      pagination.page -= 1;
      getCollectionSub(selectedCollection);
    }
  }

  async function handleCheckAll(allChecked) {
    datasChecked = datasChecked.map((data) => ({
      id: data.id,
      checked: allChecked ? false : true,
    }));

    // console.log("Updated datasChecked (handleCheckAll):", datasChecked);
  }

  async function handleCheck(id) {
    if (!Array.isArray(datasChecked)) {
      datasChecked = [];
    }

    datasChecked = datasChecked.map((data) => {
      if (data.id === id) {
        return { id: data.id, checked: !data.checked };
      }
      return data;
    });

    let pageNow = paginationModal.page;
    if (tmpDatasChecked[pageNow] === undefined) {
      tmpDatasChecked = {
        ...tmpDatasChecked,
        [pageNow]: datasChecked,
      };
    } else {
      tmpDatasChecked[pageNow] = datasChecked;
    }
    // console.log("datasChecked setelah pembaruan:", datasChecked);
  }

  function areAllUnchecked() {
    // Gabungkan data dari halaman sekarang dengan data yang tersimpan di tmpDatasChecked
    const allDataModal = [...datasChecked, ...Object.values(tmpDatasChecked).flat()];

    // Cek apakah semua data tidak ada yang checked
    return allDataModal.every((item) => item.checked === false);
  }

  const underscoreToSpace = (text) => {
    return text.replace(/_/g, " ");
  };

  async function getCollections() {
    let token = localStorage.getItem("token");

    if (token !== null) {
      let pageNow = 1;
      initCollections = [];
      while (true) {
        let response = await $call.call("client-dashboard", "get-library-collections", {
          URLQuery: {},
          URLParams: {},
          Params: {
            page: pageNow,
            page_size: 25,
            parent_id: parseInt(libraryId),
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
          initCollections = initCollections.concat(response.result.data);
          pageNow++;
        } else {
          break;
        }
      }

      console.log("libraries", initCollections);

      collections = initCollections.map((collection) => {
        return {
          value: collection.id,
          id: collection.id,
          name: collection.attributes.name,
          description: collection.attributes.meta?.description,
        };
      });
    }
  }

  async function getCollectionSub(id) {
    let token = localStorage.getItem("token");
    if (token !== null) {
      if (!search) {
        gettingDatas = true;
      }
      let response = await $call.call("client-dashboard", "find-library-by-parent-id", {
        URLQuery: {},
        URLParams: {},
        Params: {
          page: pagination.page,
          page_size: pagination.pageSize,
          parent_id: parseInt(id),
          search: searchTerm,
        },
        Headers: {
          Authorization: "Bearer " + token,
        },
      });

      if (response.status !== 200) {
        console.log(response);
        gettingDatas = false;
        return;
      }

      console.log("sub collections", response);

      if (response.result.data.length > 0) {
        pagination = response.result.meta.pagination;
        configsChecked = {
          search: true,
          searchHandler: handleSearch,
          pagination: response.result.meta.pagination,
          space: "p-4",
          // actionInLeft: true,
          // actions: [
          //   {
          //     id: "deselect_all",
          //     name: "Deselect All",
          //     onclick: () => {
          //       handleCheckAllCollection(false);
          //     },
          //   },
          //   {
          //     id: "select_all",
          //     name: "Select All",
          //     type: "brand-inverse",
          //     onclick: () => {
          //       handleCheckAllCollection(true);
          //     },
          //   },
          // ],
        };
        let promises = response.result.data.map(async (document) => {
          const sub_collections = await getSubDocuments(document.id);
          return {
            id: document.id,
            name: document.attributes.name,
            type: document.attributes.type,
            path: document.attributes.path,
            sub_collections,
          };
        });
        allDataCollection = await Promise.all(promises);
        console.log("all data", allDataCollection);
        datasChecked = allDataCollection.map((data) => {
          return {
            id: data.id,
            checked: false,
            sub_collection: data.sub_collections.map((sub) => ({
              id: sub.id,
              checked: false,
            })),
          };
        });

        if (checkAll) {
          datasChecked = allDataCollection.map((data) => {
            return {
              id: data.id,
              checked: true,
              sub_collection: data.sub_collections.map((sub) => ({
                id: sub.id,
                checked: true,
              })),
            };
          });
        } else {
          if (tmpDatasChecked[pagination.page] === undefined) {
            if (selectedCollection === tmpSelectedCollection) {
              tmpSelectedCollection = null;
              datasChecked = datasChecked.map((data) => {
                if (collectionsChecked.includes(data.id)) {
                  data.checked = true;
                }

                // Perbarui status checked untuk setiap sub_collection
                data.sub_collection = data.sub_collection.map((sub) => {
                  if (collectionsChecked.includes(sub.id)) {
                    sub.checked = true;
                  }
                  return sub;
                });

                return data;
              });
            } else {
              datasChecked = allDataCollection.map((data) => {
                return { id: data.id, checked: false };
              });
            }
          } else {
            datasChecked = tmpDatasChecked[pagination.page];
          }
        }
      } else {
        pagination = response.result.meta.pagination;
        allDataCollection = [];
        datasChecked = [];
        configsChecked = {
          search: true,
          searchHandler: handleSearch,
          pagination: response.result.meta.pagination,
          space: "p-4",
          actionInLeft: false,
          // actionInLeft: true,
          // actions: [
            // {
            //   id: "deselect_all",
            //   name: "Deselect All",
            //   onclick: () => {
            //     handleCheckAllCollection(false);
            //   },
            // },
            // {
            //   id: "select_all",
            //   name: "Select All",
            //   type: "brand-inverse",
            //   onclick: () => {
            //     handleCheckAllCollection(true);
            //   },
            // },
          // ],
          // actions: [
          //   {
          //     id: "delete",
          //     name: "Delete",
          //     icon: "trash",
          //     type: "error",
          //     disabled: datasChecked.find((data) => data.checked) ? false : true,
          //     onclick: () => {
          //       if (datasChecked.find((data) => data.checked)) {
          //         popupModal = true;
          //       }
          //     },
          //   },
          //   {
          //     id: "add_new",
          //     name: "Add New",
          //     icon: "plus",
          //     type: "brand-inverse",
          //     onclick: () => {
          //       selectTypeCollection("Documents");
          //       addDocumentsModal = true;
          //       addOtherCategory = false;
          //     },
          //   },
          // ],
        };
      }

      gettingDatas = false;
      search = false;
    }
  }

  async function getSubCollections(id) {
    let token = localStorage.getItem("token");

    if (token !== null) {
      gettingDatas = true;
      const aggregatedData = [];

      paginationCollection.page = 1;

      do {
        let response = await $call.call("client-dashboard", "get-library-collections", {
          URLQuery: {},
          URLParams: {},
          Params: {
            page: paginationCollection.page,
            page_size: paginationCollection.pageSize,
            parent_id: parseInt(id),
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

      const promise = aggregatedData.map(async (collection) => ({
        id: collection.id,
        collection_name: collection.attributes.name,
      }));

      const subCollectionList = await Promise.all(promise);
      console.log("subCollectionList", subCollectionList);
      gettingDatas = false;
      return subCollectionList;
    } else {
      triggerToast($_("unauthorized"), true, "error");
      gettingDatas = false;
      return;
    }
  }

  async function getSubDocuments(id) {
    let token = localStorage.getItem("token");

    if (token !== null) {
      gettingDatas = true;
      const aggregatedData = [];

      paginationCollection.page = 1;

      do {
        let response = await $call.call("client-dashboard", "get-library-documents", {
          URLQuery: {},
          URLParams: {},
          Params: {
            page: paginationCollection.page,
            page_size: paginationCollection.pageSize,
            parent_id: parseInt(id),
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

      const promise = aggregatedData.map(async (collection) => ({
        id: collection.id,
        collection_name: collection.attributes.name,
        path: collection.attributes.path,
      }));

      const subCollectionList = await Promise.all(promise);
      console.log("subDocumentList", subCollectionList);
      gettingDatas = false;
      return subCollectionList;
    } else {
      triggerToast($_("unauthorized"), true, "error");
      gettingDatas = false;
      return;
    }
  }

  async function getCsvData() {
    gettingDatas = true;
    let token = localStorage.getItem("token");
    if (token !== null) {
      try {
        let response = await $call.call("client-dashboard", "get-library-document-csv", {
          URLQuery: {},
          URLParams: {
            id: selectedDocument.toString(),
          },
          Params: {
            page: csvMeta.page,
            page_size: csvMeta.pageSize,
          },
          Headers: {
            Authorization: "Bearer " + token,
          },
        });

        csvMeta.pageCount = response.result.result.total_pages;
        csvMeta.total = response.result.result.total_items;

        console.log("document csv response", response);

        documentPath = response.result.data.path;

        documentTitle = response.result.data.name;

        console.log("documentPath", documentPath);

        if (response.status !== 200) {
          console.log(response);
          triggerLoading(false);
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
          goto(await $call.getPage("client-dashboard", "document-library"));
          return;
        }

        csvMeta.pageCount = response.result.result.total_pages;
        csvMeta.total = response.result.result.total_items;

        currentArrayData = response.result.result.data;

        console.log("currentArrayData", currentArrayData);

        gettingDatas = false;
      } catch (error) {
        gettingDatas = false;
        console.error("Error during document get:", error);
        triggerToast($_("failedGetDocument"), true, "error");
      }
    } else {
      triggerToast($_("unauthorized"), true, "error");
      return;
    }
  }

  const handleSearch = async (value) => {
    search = true;
    pagination.page = 1;
    searchTerm = value;

    await getCollectionSub(selectedCollection);
  };

  $effect(() => {
    if (selectedCollection) {
      if (typeof selectedCollection === "number" && selectedCollection !== oldSelectedCollection) {
        oldSelectedCollection = selectedCollection;
        tmpDatasChecked = [];
        getCollectionSub(selectedCollection);
      }
    }
  });

  $effect(() => {
    if (selectedDocument) {
      console.log("selectedDocument", selectedDocument);
    }
  });

  onMount(() => {
    // if (allDataMeta) {
    //   currentArrayData = allData.slice(
    //     (allDataMeta.page - 1) * allDataMeta.pageSize,
    //     allDataMeta.page * allDataMeta.pageSize,
    //   );
    // }

    connectWebSocket();
    if (textarea) {
      autosize();
    }
    const resizeHandler = () => {
      innerWidth = globalThis.innerWidth;
    };
    // libraryId = localStorage.getItem("libraryId") || $datasourceLibraryId;
    libraryId = localStorage.getItem("libraryId" + "insights");
    globalThis.addEventListener("resize", resizeHandler);
    resizeHandler(); // call once initially to set the correct value
    // getContentLibrary();
    getCollections();
    return () => {
      globalThis.removeEventListener("resize", resizeHandler);
    };
  });
</script>

<div class="flex items-center justify-between mb-3 md:mb-8 pt-5 px-4">
  <div class="flex flex-col gap-2">
    <div class="flex gap-4">
      <Typography
        preset="Display sm/Semibold"
        color="primary"
        text={documentTitle}
        propClass="fluid-xl"
      />
    </div>
  </div>

  <div class="flex relative gap-3">
    <Button
      color="light"
      class="min-w-max gap-1.5 px-3.5"
      on:click={() => {
        modalOpen = true;
      }}
    >
      <File06Icon svgClass="w-5 h-5 text-fg-brand-primary-light dark:text-fg-brand-primary-dark" />
      <div class="font-semibold hidden md:block text-fg-brand-primary-light dark:text-fg-brand-primary-dark">
        {$_("changeDataSource")}
      </div>
      <ChevronDownIcon
        svgClass="w-5 h-5 text-fg-brand-primary-light dark:text-fg-brand-primary-dark"
      />
    </Button>
  </div>
</div>

<Divider />

<div class="flex flex-col md:flex-row w-full pb-12">
  <!-- table preview -->
  <div class={`w-full md:w-1/2 pt-8 px-4 flex flex-col gap-6`}>
    <Typography preset="Text md/Semibold" text={$_("providedTable")} />
    {#if currentArrayData.length === 0}
      <StateCard illustration={IllustrationEmpty} description={$_("noDataAvailable")} funcButton={() => (modalOpen = true)} actionText={$_("selectDataSource")} />
    {:else}
      <div
        class="bg-gray-50 dark:bg-gray-700 border border-gray-200 dark:border-gray-700 rounded-xl pt-4"
      >
        <Table class="w-full">
          <TableHead>
            {#each Object.keys(currentArrayData[0]) as key}
              {#if key === "id"}
                <!-- empty -->
              {:else}
                <TableHeadCell>{underscoreToSpace(key)}</TableHeadCell>
              {/if}
            {/each}
          </TableHead>
          <TableBody tableBodyClass="divide-y">
            {#each currentArrayData as data}
              <TableBodyRow class="border-0 !border-gray-200 dark:!border-gray-700">
                {#each Object.entries(data) as [key, value]}
                  {#if key === "id"}
                    <!-- empty -->
                  {:else}
                    <TableBodyCell>{value}</TableBodyCell>
                  {/if}
                {/each}
              </TableBodyRow>
            {/each}
          </TableBody>
        </Table>
        <Pagination pagination={csvMeta} previousPage={previous} nextPage={next} />
        <!-- <div class="flex justify-between items-center p-3 rounded-b-xl">
          <Typography
            preset="Text sm/Regular"
            text={`Showing ${allDataMeta.page * allDataMeta.pageSize - allDataMeta.pageSize + 1} to ${
              allDataMeta.page * allDataMeta.pageSize
            } of ${allDataMeta.total} entries`}
          />
          <Pagination
            on:next={next}
            on:previous={previous}
            activeClass="bg-hc-graylight-100 dark:bg-hc-graydark-700 dark:border-[#374151] dark:text-hc-graydark-300"
          >
            <span slot="prev">{$_("previous")}</span>
            <span slot="next">{$_("next")}</span>
          </Pagination>
        </div> -->
      </div>
    {/if}
  </div>
  <!-- chat -->
  <div class="w-full md:w-1/2 pt-8 px-4 flex flex-col gap-4">
    <Typography preset="Text md/Semibold" text={$_("csvInsight")} />
    {#each dataConversationChats as chat}
      <div class="flex gap-4 p-4">
        {#if chat.role === "assistant"}
          <div
            class="w-8 h-7 md:w-11 md:h-11 flex justify-center items-center bg-fg-brand-primary-light dark:bg-fg-brand-primary-dark rounded-full"
          >
            <Icon name="lightbulb04" color="primary-onbrand" />
          </div>
        {:else}
          <div
            class="w-8 h-7 md:w-11 md:h-11 flex justify-center items-center bg-tertiary-light dark:bg-tertiary-dark rounded-full border border-primary-light dark:border-primary-dark"
          >
            <img src={UserIcon} alt="User" class="w-4 h-4 md:w-6 md:h-6" />
          </div>
        {/if}
        <div class="w-5/6">
          <div class="text-wrap font-medium fluid-base md:fluid-base text-tertiary-light dark:text-primary-dark">
            {#if chat.role === "assistant"}
              <pre class="whitespace-pre-wrap break-words overflow-auto p-2">{chat.message}</pre>
            {:else}
              <Markdown md={chat.message} {plugins} />
            {/if}
            {#if chat.image}
              <img src={chat.image} alt="Image" class="w-full" />
            {/if}
          </div>
        </div>
      </div>
    {/each}

    {#if isTyping}
    <div class="flex gap-4 p-4">
      <div
        class="w-8 h-7 md:w-11 md:h-11 flex justify-center items-center bg-fg-brand-primary-light dark:bg-fg-brand-primary-dark rounded-full"
      >
        <Icon name="lightbulb04" color="primary-onbrand" />
      </div>
      <div class="text-wrap font-medium fluid-base md:fluid-base text-tertiary-light dark:text-primary-dark">
        <pre class="whitespace-pre-wrap break-words overflow-auto p-2">Typing...</pre>
      </div>
    </div>
    {/if}

    <form onsubmit={messageInput !== "" ? sendMessageInput : () => {}} class="w-full p-6">
      <div class="flex items-center">
        <div class="w-full flex gap-4">
          <!-- svelte-ignore event_directive_deprecated -->
          <textarea
            bind:this={textarea}
            bind:value={messageInput}
            onkeydown={handleKeyDown}
            disabled={isTyping || currentArrayData.length === 0}
            id="chat"
            rows="1"
            class="block resize-none p-2.5 w-full fluid-sm text-brand-primary-light dark:text-primary-dark bg-white dark:bg-primary-dark rounded-3xl border border-utility-brand-500-light dark:border-utility-brand-500-dark focus:border-utility-brand-500-light dark:focus:border-utility-brand-500-dark focus:ring-utility-brand-400-light dark:focus:ring-utility-brand-400-light dark:placeholder-gray-400"
            placeholder="Send a message..."
          ></textarea>
          <Button
            disabled={!messageInput || isTyping || currentArrayData.length === 0}
            onclick={onclickSend}
            class="rounded-full h-11 w-11 !p-3"
          >
            <img src={SendIcon} alt="SendIcon" class="w-5 h-5" />
          </Button>
        </div>
      </div>
    </form>
  </div>
</div>

<Modal bind:open={modalOpen} outsideclose>
  <svelte:fragment slot="header">
    <div class="flex gap-4 me-6 grow w-full">
      <div
        class="grow pt-2 border-t-4 border-fg-brand-primary-alt-light dark:border-fg-brand-primary-alt-dark"
      >
        <Typography
          preset="Text md/Bold"
          text={$_("chooseCollection")}
          color={!isSelectedLibrary ? "brand-secondary" : "secondary"}
        />
        <Typography
          preset="Text sm/Regular"
          text={$_("chooseCollectionAsDataSource")}
          color={!isSelectedLibrary ? "brand-tertiary" : "tertiary"}
        />
      </div>

      <div
        class={`grow pt-2 border-t-4 ${isSelectedLibrary ? "border-fg-brand-primary-alt-light dark:border-fg-brand-primary-alt-dark" : "border-utility-gray-200-light dark:border-utility-gray-200-dark"}`}
      >
        <Typography
          preset="Text md/Bold"
          text={$_("chooseCsvfile")}
          color={isSelectedLibrary ? "brand-secondary" : "secondary"}
        />
        <Typography
          preset="Text sm/Regular"
          text={$_("chooseTheOneYouWantToInclude")}
          color={isSelectedLibrary ? "brand-tertiary" : "tertiary"}
        />
      </div>
    </div>
  </svelte:fragment>
  {#if !isSelectedLibrary}
    <div class="flex flex-col justify-center gap-4">
      <Typography preset="Text sm/Medium" text={$_("latestDocumentLibrary")} color="quaterary" />

      <SelectedCard
        selectionData={collections}
        bind:selectedData={selectedCollection}
        numGrid={2}
        isRadio
      />
    </div>
  {:else if allDataCollection.length > 0 || searchTerm !== ""}
    {#if !gettingDatas}
      <DataTableListCollection
        bind:datas={allDataCollection}
        bind:datasChecked
        bind:configs={configsChecked}
        bind:gettingDatas
        bind:next={nextCollectionDoc}
        bind:previous={previousCollectionDoc}
        bind:handleCheckAll
        bind:handleCheck
        bind:selectedData={selectedDocument}
        bind:searchTerm
        {selectCsv}
      />
    {:else}
      <StateCard illustration={IllustrationProgress} description={$_("fetchingDataPleaseWait")} />
    {/if}
  {:else if gettingDatas}
    <StateCard illustration={IllustrationProgress} description={$_("fetchingDataPleaseWait")} />
  {:else}
    <StateCard illustration={IllustrationEmpty} description={$_("youDontHaveAnyDataYet")} />
  {/if}
  <svelte:fragment slot="footer">
    {#if !isSelectedLibrary}
      <Button
        on:click={() => {
          isSelectedLibrary = true;
        }}
        class="w-full"
        disabled={!selectedCollection}>{$_("next")}</Button
      >
    {:else}
      <div class="flex gap-4 w-full">
        <button
          onclick={() => {
            isSelectedLibrary = false;
          }}
          class="px-5">{$_("goBack")}</button
        >
        <Button
          on:click={async () => {
            modalOpen = false;
            await getCsvData();
            Select(current_session_user, documentPath);
          }}
          class="grow"
          disabled={selectedDocument === ""}>{$_("finish")}</Button
        >
      </div>
    {/if}
  </svelte:fragment>
</Modal>
