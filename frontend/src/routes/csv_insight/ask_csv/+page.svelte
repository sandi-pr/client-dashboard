<script>
  import { Divider, Typography, StateCard, Icon, call, SelectedCard } from "$lib";
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
    Pagination,
  } from "flowbite-svelte";
  import { onMount } from "svelte";

  const isFrontend = process.env.NODE_ENV === "development";

  const plugins = [gfmPlugin()];

  let messageInput = $state("");
  let textarea = $state(null);
  let modalOpen = $state(false);
  let allData = $state([
    {
      id: 1,
      customer_id: 1,
      name: "John Doe",
      rating: 5,
      date: "2021-08-01",
    },
    {
      id: 2,
      customer_id: 2,
      name: "Jane Doe",
      rating: 4,
      date: "2021-08-02",
    },
    {
      id: 3,
      customer_id: 3,
      name: "John Smith",
      rating: 3,
      date: "2021-08-03",
    },
    {
      id: 4,
      customer_id: 4,
      name: "Jane Smith",
      rating: 2,
      date: "2021-08-04",
    },
    {
      id: 5,
      customer_id: 5,
      name: "John Doe",
      rating: 1,
      date: "2021-08-05",
    },
    {
      id: 6,
      customer_id: 6,
      name: "Jane Doe",
      rating: 5,
      date: "2021-08-06",
    },
    {
      id: 7,
      customer_id: 7,
      name: "John Smith",
      rating: 4,
      date: "2021-08-07",
    },
    {
      id: 8,
      customer_id: 8,
      name: "Jane Smith",
      rating: 3,
      date: "2021-08-08",
    },
    {
      id: 9,
      customer_id: 9,
      name: "John Doe",
      rating: 2,
      date: "2021-08-09",
    },
    {
      id: 10,
      customer_id: 10,
      name: "Jane Doe",
      rating: 1,
      date: "2021-08-10",
    },
    {
      id: 11,
      customer_id: 11,
      name: "John Smith",
      rating: 5,
      date: "2021-08-11",
    },
    {
      id: 12,
      customer_id: 12,
      name: "Jane Smith",
      rating: 4,
      date: "2021-08-12",
    },
    {
      id: 13,
      customer_id: 13,
      name: "John Doe",
      rating: 3,
      date: "2021-08-13",
    },
  ]);
  let allDataCollection = $state([]);
  let currentArrayData = $state([]);
  let allDataMeta = $state({ page: 1, pageSize: 10, pageCount: 2, total: 13 });
  let allDataModal = $state([
    {
      id: 1,
      file_name: "quarterly_revenue.csv",
    },
    {
      id: 2,
      file_name: "ad_campaign_performance.csv",
    },
  ]);
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
  let configsChecked = $state({});
  let tmpSelectedCollection = $state("");
  let oldSelectedCollection = $state("");
  let search = $state(false);
  let checkAll = $state(false);
  let selectCsv = $state(true);

  let paginationModal = $state({ page: 1, pageSize: 3, pageCount: 0, total: 0 });
  let pagination = $state({ page: 1, pageSize: 10, pageCount: 0, total: 0 });
  let paginationCollection = $state({ page: 1, pageSize: 20, pageCount: 0, total: 0 });

  let dataConversationChatsDummy = $state([
    {
      id: 1,
      role: "user",
      message: `Bisakah Anda menjelaskan apa itu kiamat?`,
    },
    {
      id: 2,
      role: "assistant",
      message: `<think>
Okay, the user is asking "Apa itu kiamat," which translates to "What is kiamat?" I need to explain this concept based on Islamic teachings and any provided context.

From the existing answer, I see that kiamat refers to the Day of Judgment in Islam. It's crucial to elaborate on both the small and major aspects of kiamat—kiamat kecil and kiamat besar—as well as引用相关的 Quranic verses for clarity.

I should structure the response to include:
- Definition of kiamat
- Differences between kiamat kecil and kiamat syih
- Quranic references and scientific perspectives from the documents
- The significance of kiamat in Islamic belief

Make sure to avoid speculation and only use information from the provided context.
</think>

Kiamat adalah istilah yang merujuk pada hari akhir atau penghujung zaman di mana segala makhluk hidup akan menghadapi ajalnya dan kehidupan sebagaimana kita ketahui akan berakhir. Menurut ajaran Islam, kiamat adalah suatu peristiwa maha dahsyat yang tidak dapat dipungkiri sebagai bagian dari keyakinan muslim.

### 1. **Definisi Kiamat**
Kiamat dalam terminologi islam adalah hari dimana segala makhluk akan mati dan dihisab (diperiksa) oleh Allah SWT. Ini termasuk kiamat kecil dan kiamat besar:
   - **Kiamat Kecil (Kiamat Syiah):**
     Peristiwa-peristiwa yang terjadi sepanjang zaman, seperti bencana alam, erupsi gunung berapi, dan jatuhnya meteor, yang mengakibatkan kepunahan massal atau perubahan besar dalam kehidupan di muka bumi.
   - **Kiamat Besar (Kiamat Syih):**
     Hari akhir dimana seluruh alam semesta akan musnah dan segala makhluk hidup akan dimintahkan pertanggungjawaban.

### 2. **Penggambaran Kiamat dalam Al-Qur'an**
Al-Qur'an menggambarkan kiamat dengan gambaran yang sangat dramatis. Beberapa surah seperti Surah al-Qari'ah, Surah az-Zumar, dan Surah ad-Dukhan menceritakan tentang kehancuran bumi dan langit pada hari itu:

- **Surah al-Qari'ah (101):**
  "Danhari yang kiamatnya adalah hari yang menakutkan."

- **Surah az-Zumar (39:68):**
  "Danhari itu, tidak ada jiwa yang dapat memberikan hadiah kepada jiwa lain melain dari karena ridha Allah SWT."

### 3. **Kiamat dan Ilmu Pengetahuan**
Dari sisi ilmu pengetahuan, kiamat besar dapat disimak sebagai peristiwa dimana alam semesta ini akan berakhir sesuai dengan hadis-hadis nabi yang menyebutkan tentang hari kiamat.

### 4. **Peran Air di Surga dan Neraka**
Di dalam Al-Qur'an, air juga dijadikan simbol kehidupan yang melimpah di surga (Jannah) dan menjadi puniyah untuk muktabar di neraka. Ini tercermin dalam ayat-ayat seperti:

- **Surah As-Saffat (37:44):**
  "Dan Dia menjanjikan kepada mereka air yang tidak berkurang dan buah-buah yang tidak ada habisnya."

### 5. **Simbolisme Air**
Air di dalam Al-Qur'an juga menjadi simbol kebersihan, kehidupan, dan pertumbuhan. Ini tercermin dari ayat-ayat seperti:

- **Surah Al-Ma'ida (5:6):**
  "Dan jika kamu tidak menemukan air, maka bersuci dengan tanah yang baik."

### Kesimpulan
Kiamat adalah peristiwa akhir zaman yang ditandai dengan kehancuran segala makhluk dan alam semesta ini kembali kepada Sang Pencipta. Air menjadi bagian penting dalam konteks kehidupan di dunia dan surga, menegaskan betapa besar nikmat Allah SWT terhadap umat-Nya.

Reference: [file...](/bakrie/document_detail?documentId=3&page=1 "file")[[P1](/bakrie/document_detail?documentId=3&page=1 "file (Page 1)"), [P5](/bakrie/document_detail?documentId=3&page=5 "file (Page 5)")], [PT_BUMI_RESOURC......](/bakrie/document_detail?documentId=5&page=1 "PT_BUMI_RESOURCES_TBK")[[P1](/bakrie/document_detail?documentId=5&page=1 "PT_BUMI_RESOURCES_TBK (Page 1)"), [P2](/bakrie/document_detail?documentId=5&page=2 "PT_BUMI_RESOURCES_TBK (Page 2)"), [P3](/bakrie/document_detail?documentId=5&page=3 "PT_BUMI_RESOURCES_TBK (Page 3)"), [P4](/bakrie/document_detail?documentId=5&page=4 "PT_BUMI_RESOURCES_TBK (Page 4)"), [P5](/bakrie/document_detail?documentId=5&page=5 "PT_BUMI_RESOURCES_TBK (Page 5)")]`,
    },
  ]);

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
    dataConversationChatsDummy = [
      ...dataConversationChatsDummy,
      {
        id: dataConversationChatsDummy.length + 1,
        role: "user",
        message: messageInput,
      },
    ];
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
    if (paginationModal.page >= paginationModal.pageCount) return;

    let pageNow = paginationModal.page;
    if (tmpDatasChecked[pageNow] === undefined) {
      tmpDatasChecked = {
        ...tmpDatasChecked,
        [pageNow]: datasChecked,
      };
    } else {
      tmpDatasChecked[pageNow] = datasChecked;
    }

    paginationModal.page = paginationModal.page + 1;
    // TODO: create function getCsv();
  }

  async function previous() {
    if (paginationModal.page <= 1) return;

    let pageNow = paginationModal.page;
    if (tmpDatasChecked[pageNow] === undefined) {
      tmpDatasChecked = {
        ...tmpDatasChecked,
        [pageNow]: datasChecked,
      };
    } else {
      tmpDatasChecked[pageNow] = datasChecked;
    }

    paginationModal.page = paginationModal.page - 1;
    // TODO: create function getCsv();
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
        let response = await $call.call("client-dashboard", "get-librarys", {
          URLQuery: {},
          URLParams: {},
          Params: {
            page: pageNow,
            page_size: 25,
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
      let response = await $call.call("client-dashboard", "get-library-collections", {
        URLQuery: {},
        URLParams: {},
        Params: {
          page: pagination.page,
          page_size: pagination.pageSize,
          parent_id: id,
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

      console.log("collections", response);

      if (response.result.data.length > 0) {
        pagination = response.result.meta.pagination;
        configsChecked = {
          search: true,
          searchHandler: handleSearch,
          pagination: response.result.meta.pagination,
          space: "p-4",
          actionInLeft: true,
          actions: [
            {
              id: "deselect_all",
              name: "Deselect All",
              onclick: () => {
                handleCheckAllCollection(false);
              },
            },
            {
              id: "select_all",
              name: "Select All",
              type: "brand-inverse",
              onclick: () => {
                handleCheckAllCollection(true);
              },
            },
          ],
        };
        let promises = response.result.data.map(async (document) => {
          const sub_collections = await getSubCollections(document.id);
          return {
            id: document.id,
            collection_name: document.attributes.name,
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
          actionInLeft: true,
          actions: [
            {
              id: "deselect_all",
              name: "Deselect All",
              onclick: () => {
                handleCheckAllCollection(false);
              },
            },
            {
              id: "select_all",
              name: "Select All",
              type: "brand-inverse",
              onclick: () => {
                handleCheckAllCollection(true);
              },
            },
          ],
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

  onMount(() => {
    if (allDataMeta) {
      currentArrayData = allData.slice(
        (allDataMeta.page - 1) * allDataMeta.pageSize,
        allDataMeta.page * allDataMeta.pageSize,
      );
    }
    if (textarea) {
      autosize();
    }
    const resizeHandler = () => {
      innerWidth = globalThis.innerWidth;
    };
    globalThis.addEventListener("resize", resizeHandler);
    resizeHandler(); // call once initially to set the correct value
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
        text="Customer Feedback Data"
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
      <div class="font-semibold text-fg-brand-primary-light dark:text-fg-brand-primary-dark">
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
  <div class="w-full md:w-1/2 pt-8 px-4 flex flex-col gap-6">
    <Typography preset="Text md/Semibold" text={$_("providedTable")} />
    {#if currentArrayData.length === 0}
      <StateCard illustration={IllustrationEmpty} description={$_("noDataAvailable")} />
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
        <div class="flex justify-between items-center p-3 rounded-b-xl">
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
        </div>
      </div>
    {/if}
  </div>
  <!-- chat -->
  <div class="w-full md:w-1/2 pt-8 px-4 flex flex-col gap-4">
    <Typography preset="Text md/Semibold" text={$_("csvInsight")} />
    {#each dataConversationChatsDummy as chat}
      <div class="flex gap-4 p-4">
        {#if chat.role === "assistant"}
          <div
            class="w-8 h-7 md:w-12 md:h-11 flex justify-center items-center bg-fg-brand-primary-light dark:bg-fg-brand-primary-dark rounded-full"
          >
            <Icon name="lightbulb04" color="primary-onbrand" />
          </div>
        {:else}
          <div
            class="w-8 h-7 md:w-12 md:h-11 flex justify-center items-center bg-tertiary-light dark:bg-tertiary-dark rounded-full border border-primary-light dark:border-primary-dark"
          >
            <img src={UserIcon} alt="User" class="w-4 h-4 md:w-6 md:h-6" />
          </div>
        {/if}
        <div class="w-full">
          <div
            class="prose max-w-none font-medium fluid-base md:fluid-base text-tertiary-light dark:text-primary-dark"
          >
            <Markdown md={chat.message} {plugins} />
          </div>
        </div>
      </div>
    {/each}

    <form onsubmit={messageInput !== "" ? sendMessageInput : () => {}} class="w-full p-6">
      <div class="flex items-center">
        <div class="w-full flex gap-4">
          <!-- svelte-ignore event_directive_deprecated -->
          <textarea
            bind:this={textarea}
            bind:value={messageInput}
            onkeydown={handleKeyDown}
            id="chat"
            rows="1"
            class="block resize-none p-2.5 w-full fluid-sm text-brand-primary-light dark:text-primary-dark bg-white dark:bg-primary-dark rounded-3xl border border-utility-brand-500-light dark:border-utility-brand-500-dark focus:border-utility-brand-500-light dark:focus:border-utility-brand-500-dark focus:ring-utility-brand-400-light dark:focus:ring-utility-brand-400-light dark:placeholder-gray-400"
            placeholder="Send a message..."
          ></textarea>
          <Button
            disabled={!messageInput}
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
          text={$_("chooseDocumentsLibrary")}
          color={!isSelectedLibrary ? "brand-secondary" : "secondary"}
        />
        <Typography
          preset="Text sm/Regular"
          text={$_("documentLibraryAsThePrimarySource")}
          color={!isSelectedLibrary ? "brand-tertiary" : "tertiary"}
        />
      </div>

      <div
        class={`grow pt-2 border-t-4 ${isSelectedLibrary ? "border-fg-brand-primary-alt-light dark:border-fg-brand-primary-alt-dark" : "border-utility-gray-200-light dark:border-utility-gray-200-dark"}`}
      >
        <Typography
          preset="Text md/Bold"
          text={$_("chooseWhichCollections")}
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
        bind:next
        bind:previous
        bind:handleCheckAll
        bind:handleCheck
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
          on:click={() => {
            modalOpen = false;
          }}
          class="grow"
          disabled={areAllUnchecked()}>{$_("finish")}</Button
        >
      </div>
    {/if}
  </svelte:fragment>
</Modal>
