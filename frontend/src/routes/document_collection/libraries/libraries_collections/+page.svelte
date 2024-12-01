<script>
  import { _ } from "svelte-i18n";
  import { Button, Modal, Input, Alert, Select, P } from "flowbite-svelte";

  import { call } from "$lib";
  import {
    Typography,
    ToastContainer,
    StateCard,
    triggerToast,
    InputRw,
    TextareaRw,
    LibraryCard,
    check,
  } from "$lib";
  import TabBar from "$lib/components/TabBar.svelte";
  import { onMount, onDestroy } from "svelte";
  import { page } from "$app/stores";
  import { goto } from "$app/navigation";

  import { ExclamationCircleSolid } from "flowbite-svelte-icons";
  import Copy01Icon from "$lib/icon/copy-01.svelte";
  import DotsVerticalIcon from "$lib/icon/dots-vertical.svelte";
  import FolderPlusIcon from "$lib/icon/folder-plus.svelte";
  import Stars03Icon from "$lib/icon/stars-03.svelte";
  import Folder from "$lib/icon/folder.svelte";
  import IllustrationEmpty from "$lib/icon/illustration_empty.svelte";
  import IllustrationProgress from "$lib/icon/illustration_progress.svelte";
  import { format } from "date-fns";

  let addCollectionModal = $state(false);
  let editModal = $state(false);
  let editLibraryModal = $state(false);
  let deleteModal = $state(false);
  let deleteLibraryModal = $state(false);
  let tmpCollectionName = $state("");
  let tmpDescription = $state("");
  let tmpLibraryName = $state("");
  let tmpDescLibrary = $state("");
  let tmpCollectionId = $state("");
  let libraryId = $state("");
  let dropdownOpen = $state(false);

  let pagination = $state({ page: 1, pageSize: 20, pageCount: 1, total: 0 });

  let collectionDatas = $state({
    title: "",
    description: "",
    createdAt: "",
  });

  let tabDatas = $state({
    selected: "My Collections",
    data: [
      {
        name: "My Collections",
      },
      {
        name: "Information",
      },
    ],
  });
  let collectionListDatas = $state([]);
  let informationDatas = {
    id: "1",
    collectionName: "Bakrie’s Handwriting",
    totalDocuments: 2,
    type: "",
    status: "ready",
  };

  async function getLibrary() {
    let token = localStorage.getItem("token");

    if (token !== null) {
      let response = await $call.call("client-dashboard", "get-library", {
        URLQuery: {},
        URLParams: {
          id: libraryId,
        },
        Params: {},
        Headers: {
          Authorization: "Bearer " + token,
        },
      });

      console.log("library", response);

      collectionDatas = {
        title: response.result.attributes.name,
        description: response.result.attributes.meta.description,
        createdAt: format(new Date(response.result.attributes.createdAt), "dd-MM-yyyy"),
      };
    }
  }

  async function getCollections() {
    let token = localStorage.getItem("token");

    if (token !== null) {
      getLibrary();

      let response = await $call.call("client-dashboard", "get-library-collections", {
        URLQuery: {},
        URLParams: {},
        Params: {
          page: pagination.page,
          page_size: pagination.pageSize,
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

      // id: "4",
      // collectionName: "Study Scriptures of Qur’an Tafsir",
      // totalDocuments: 2,
      // type: "",
      // status: "reading",

      let link = await $call.getPage("client-dashboard", "document-library-documents");

      let promise = response.result.data.map(async (collection) => {
        const params = new URLSearchParams(window.location.search);

        params.delete("collectionId");

        const queryParams = {};
        for (const [key, value] of params.entries()) {
          queryParams[key] = value;
        }
        const queryString = new URLSearchParams(queryParams).toString();
        return {
          id: collection.id,
          collectionName: collection.attributes.name,
          description: collection.attributes.meta.description,
          type: collection.attributes.type,
          href: link + "?collectionId=" + collection.id + "&" + queryString,
        };
      });

      collectionListDatas = await Promise.all(promise);
    }
  }

  async function createCollection() {
    let token = localStorage.getItem("token");

    if (token !== null) {
      let response = await $call.call("client-dashboard", "create-library-collection", {
        URLQuery: {},
        URLParams: {},
        Params: {
          collection_name: tmpCollectionName,
          description: tmpDescription,
          parent_id: parseInt(libraryId),
          // bakrie
          bot_id: localStorage.getItem("botIdStrapi").toString(),
        },
        Headers: {
          Authorization: "Bearer " + token,
        },
      });

      if (response.status !== 200) {
        console.log(response);
        triggerToast("error");
        return;
      }

      console.log(response);

      // RAG Proccessing

      // let response2 = await $call.call("hexa-ai-rag", "add-document-by-collection-id", {
      //   URLQuery: {},
      //   URLParams: {},
      //   Params: {
      //     collection_id: response.result.data.id.toString(),
      //   },
      //   Headers: {
      //     Authorization: "Bearer " + token,
      //   },
      // });

      // if (response2.status !== 200) {
      //   console.log(response2);
      //   triggerToast("error");
      //   return;
      // }

      // console.log(response2);

      tmpCollectionName = "";
      tmpDescription = "";

      triggerToast("success");

      addCollectionModal = false;
      getCollections();
      // goto(await $call.getPage("client-dashboard", "document-library"));
    }
  }

  async function editLibraryName() {
    let token = localStorage.getItem("token");

    if (token !== null) {
      try {
        let response = await $call.call("client-dashboard", "update-library", {
          URLQuery: {},
          URLParams: {
            id: libraryId.toString(),
          },
          Params: {
            library_name: tmpLibraryName,
            description: tmpDescLibrary,
          },
          Headers: {
            Authorization: "Bearer " + token,
          },
        });

        if (response.status !== 200) {
          console.log(response);
          return;
        }

        tmpLibraryName = "";
        tmpDescLibrary = "";

        editLibraryModal = false;
        triggerToast("updated");
        getCollections();
      } catch (error) {
        console.log(error);
        triggerToast("error");
      }
    } else {
      console.log("Token not found");
    }
  }

  async function editCollectionName() {
    let token = localStorage.getItem("token");

    if (token !== null) {
      try {
        let response = await $call.call("client-dashboard", "update-library-collection", {
          URLQuery: {},
          URLParams: {
            id: tmpCollectionId.toString(),
          },
          Params: {
            collection_name: tmpCollectionName,
            description: tmpDescription,
          },
          Headers: {
            Authorization: "Bearer " + token,
          },
        });

        if (response.status !== 200) {
          console.log(response);
          return;
        }

        tmpCollectionName = "";
        tmpDescription = "";

        editModal = false;
        triggerToast("updated");
        getCollections();
      } catch (error) {
        console.log(error);
        triggerToast("error");
      }
    } else {
      console.log("Token not found");
    }
  }

  async function deleteLibrary() {
    let token = localStorage.getItem("token");

    if (token !== null) {
      try {
        // for (const document of documents) {
        //   if (document.attributes.collection_id === tmpCollectionId) {
        //     let response = await $call.call("client-dashboard", "delete-document", {
        //       URLQuery: {},
        //       URLParams: {
        //         id: document.id.toString(),
        //       },
        //       Params: {},
        //       Headers: {
        //         Authorization: "Bearer " + token,
        //       },
        //     });
        //     if (response.status !== 200) {
        //       console.log(response);
        //       triggerToast("error");
        //       return;
        //     }
        //   }
        // }
        let response2 = await $call.call("client-dashboard", "delete-library", {
          URLQuery: {},
          URLParams: {
            id: libraryId.toString(),
          },
          Params: {},
          Headers: {
            Authorization: "Bearer " + token,
          },
        });

        if (response2.status !== 200) {
          console.log(response2);
          triggerToast("error");
          return;
        }

        console.log("ayam", response2);

        triggerToast("deleted");

        getCollections();
        deleteLibraryModal = false;
      } catch (error) {
        triggerToast("error");
        console.log(error);
      }
    } else {
      console.log("Token not found");
    }
  }

  async function deleteDocumentCollection() {
    let token = localStorage.getItem("token");

    if (token !== null) {
      try {
        let response2 = await $call.call("client-dashboard", "delete-library-collection", {
          URLQuery: {},
          URLParams: {
            id: tmpCollectionId.toString(),
          },
          Params: {},
          Headers: {
            Authorization: "Bearer " + token,
          },
        });

        if (response2.status !== 200) {
          console.log(response2);
          triggerToast("error");
          return;
        }

        console.log("ayam", response2);

        triggerToast("deleted");

        getCollections();
        deleteModal = false;
      } catch (error) {
        triggerToast("error");
        console.log(error);
      }
    } else {
      console.log("Token not found");
    }
  }

  const handleDropdownFocusLoss = ({ relatedTarget, currentTarget }) => {
    if (relatedTarget instanceof HTMLElement && currentTarget.contains(relatedTarget)) return;
    dropdownOpen = false;
  };

  onMount(() => {
    if (check("documentAi.documentLibrary:view") === false) {
      window.history.back();
      return;
    }
    const unsubscribe = page.subscribe(async ($page) => {
      libraryId = $page.url.searchParams.get("libraryId");
      if (libraryId === null) {
        goto(await $call.getPage("client-dashboard", "document-library"));
        return;
      }
      getCollections();
    });
    return () => unsubscribe();
  });

  // $effect(() => {
  //   const unsubscribe = page.subscribe(async ($page) => {
  //     libraryId = $page.url.searchParams.get("libraryId");
  //     if (libraryId === null) {
  //       goto(await $call.getPage("client-dashboard", "document-library"));
  //       return;
  //     }
  //     getCollections();
  //   });
  //   onDestroy(() => unsubscribe());
  // });
</script>

<div class="flex dark:bg-primary-dark">
  <div class="w-full pb-10">
    <div class="flex flex-col m-4">
      <div class="flex flex-col">
        <!-- User Information -->
        <div class="flex items-center justify-between mb-3 md:mb-8">
          <div class="flex flex-col gap-2">
            <div class="flex gap-4">
              <Typography
                preset="Display sm/Semibold"
                color="primary"
                text={collectionDatas.title}
                propClass="fluid-xl"
              />
              <button
                onclick={() => {
                  navigator.clipboard
                    .writeText(collectionDatas.title)
                    .then(() => {
                      console.log("Text copied to clipboard");
                      triggerToast("copyToClipboard");
                    })
                    .catch((err) => {
                      console.error("Failed to copy: ", err);
                    });
                }}
              >
                <Copy01Icon svgClass="w-5 h-5 text-quinary-light dark:text-quinary-dark" />
              </button>
            </div>
            <P
              class="text-base font-medium leading-normal text-hc-graylight-500 dark:text-hc-graydark-400"
              >{collectionDatas.description}</P
            >
            <P
              class="text-sm font-medium leading-normal text-hc-graylight-500 dark:text-hc-graydark-400"
              >{$_("createOn")}{` ${collectionDatas.createdAt}`}</P
            >
          </div>

          <div class="flex gap-2">
            <Button
              disabled={!check("documentAi.documentLibrary.collection:create")}
              on:click={() => {
                addCollectionModal = true;
                tmpCollectionName = "";
                tmpDescription = "";
              }}
              class="p-2 md:px-3 md:py-2 gap-2"
              outline
            >
              <FolderPlusIcon svgClass="w-5 h-5 text-quinary-light dark:text-quinary-dark" />
              <div class="hidden md:block">{$_("addCollection")}</div>
            </Button>
            {#if check("documentAi.documentLibrary.collection:edit") || check("documentAi.documentLibrary.collection:delete")}
              <div onfocusout={handleDropdownFocusLoss} class=" flex items-center relative">
                <Button on:click={() => (dropdownOpen = !dropdownOpen)} color="light" class="p-2">
                  <DotsVerticalIcon svgClass="w-5 h-5 text-quinary-light dark:text-quinary-dark" />
                </Button>
                <div
                  class={`z-10 ${dropdownOpen ? "absolute top-10 right-0" : "hidden"} bg-white divide-y divide-gray-100 rounded-lg shadow w-max dark:bg-gray-700`}
                >
                  <ul class="py-2 grow text-sm">
                    {#if check("documentAi.documentLibrary.collection:edit")}
                      <button
                        class="flex p-2 w-full text-secondary-light dark:text-secondary-dark hover:text-secondary-hover-light dark:hover:text-secondary-hover-dark hover:bg-primary-hover-light dark:hover:bg-primary-hover-dark"
                        onclick={() => {
                          dropdownOpen = !dropdownOpen;
                          tmpLibraryName = collectionDatas.title;
                          tmpDescLibrary = collectionDatas.description;
                          editLibraryModal = true;
                        }}>{$_("editLibraryName")}</button
                      >
                    {/if}
                    {#if check("documentAi.documentLibrary.collection:delete")}
                      <button
                        class="flex p-2 w-full text-secondary-light dark:text-secondary-dark hover:text-secondary-hover-light dark:hover:text-secondary-hover-dark hover:bg-primary-hover-light dark:hover:bg-primary-hover-dark"
                        onclick={() => {
                          dropdownOpen = !dropdownOpen;
                          deleteLibraryModal = true;
                        }}>{$_("delete")}</button
                      >
                    {/if}
                  </ul>
                </div>
              </div>
            {/if}
            <!-- <Button
              on:click={() => {
                // addLibraryModal = true;
                // tmpCollectionName = "";
                // tmpDescription = "";
              }}
              color="light"
              class="p-2"
            >
              <DotsVerticalIcon svgClass="w-5 h-5 text-quinary-light dark:text-quinary-dark" />
            </Button> -->
          </div>
        </div>

        <TabBar {tabDatas} on:tabChange={(e) => (tabDatas.selected = e.detail)} />

        {#if tabDatas.selected === "My Collections" && collectionListDatas.length > 0}
          <div class="grid sm:grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4 mt-4">
            {#each collectionListDatas as collection}
              <LibraryCard
                collectionName={collection.collectionName}
                usedFor="libraries_collections"
                functButton={() => {
                  // selectCollection({
                  //   id: collection.id,
                  //   name: collection.collectionName,
                  //   description: collection.description,
                  //   totalDocuments: collection.totalDocuments,
                  // });
                }}
                functAddDocuments={() => {
                  tmpCollectionId = collection.id;
                  // tmpHref = collection.href;
                  // categories = collection.categories;
                  // selectTypeCollection("Documents");
                  // categories = collection.categories;
                  // otherCategory = "";
                  // addOtherCategory = false;
                  addDocumentsModal = true;
                }}
                functEdit={() => {
                  tmpCollectionId = collection.id;
                  tmpCollectionName = collection.collectionName;
                  tmpDescription = collection.description;
                  editModal = true;
                }}
                functDelete={() => {
                  tmpCollectionId = collection.id;
                  deleteModal = true;
                }}
                funcDetailCollection={() => {
                  goto(collection.href);
                }}
              />
            {/each}
          </div>
        {:else if tabDatas.selected === "Information" && Object.keys(informationDatas).length > 0}
          <div class="grid grid-cols-12 mt-4 gap-2">
            <div class="col-span-4">
              <P class="font-semibold text-secondary-light dark:text-secondary-dark">
                {$_("generalInformation")}
              </P>
            </div>
            <div class="col-span-8">
              <P
                class="text-sm font-medium leading-tight text-secondary-light dark:text-secondary-dark"
                >{$_("collectionTitle")}</P
              >
            </div>
            <div class="col-span-4"></div>
            <div class="col-span-8">
              <InputRw onEdit={false} bind:value={collectionDatas.title} />
            </div>
            <div class="col-span-4"></div>
            <div class="col-span-8">
              <P
                class="text-sm font-medium leading-tight text-secondary-light dark:text-secondary-dark"
                >{$_("description")}</P
              >
            </div>
            <div class="col-span-4"></div>
            <div class="col-span-8">
              <TextareaRw onEdit={false} bind:value={collectionDatas.description} />
            </div>
          </div>
        {:else}
          <StateCard illustration={IllustrationEmpty} description={$_("youDontHaveAnyDataYet")} />
        {/if}
      </div>
    </div>
  </div>
  <ToastContainer />
</div>

<Modal title={$_("createNewCollection")} bind:open={addCollectionModal} size="xs" autoclose>
  <div class="flex flex-col gap-4">
    <div class="flex flex-col gap-1.5">
      <Typography preset="Text sm/medium" color="secondary" text={`${$_("collectionName")}*`} />
      <Input bind:value={tmpCollectionName} placeholder={$_("whatIsYourCollectionTitle")} />
    </div>
    <div class="flex flex-col gap-1.5">
      <Typography preset="Text sm/medium" color="secondary" text={`${$_("description")}*`} />
      <Input bind:value={tmpDescription} placeholder={$_("egDescription")} />
    </div>
  </div>
  <svelte:fragment slot="footer">
    <Button
      onclick={createCollection}
      disabled={!tmpCollectionName || !tmpDescription}
      class="w-full">{$_("addNewCollection")}</Button
    >
  </svelte:fragment>
</Modal>

<Modal title={$_("editCollectionName")} bind:open={editModal} size="xs" autoclose>
  <div class="flex flex-col gap-4">
    <div class="flex flex-col gap-1.5">
      <Typography preset="Text sm/medium" color="secondary" text={`${$_("collectionName")}*`} />
      <Input bind:value={tmpCollectionName} placeholder={$_("whatIsYourCollectionTitle")} />
    </div>
    <div class="flex flex-col gap-1.5">
      <Typography preset="Text sm/medium" color="secondary" text={$_("description")} />
      <Input bind:value={tmpDescription} placeholder={$_("egDescription")} />
    </div>
  </div>
  <svelte:fragment slot="footer">
    <Button
      onclick={editCollectionName}
      disabled={!tmpCollectionName || !tmpDescription}
      class="me-2 w-full">{$_("save")}</Button
    >
  </svelte:fragment>
</Modal>

<Modal title={$_("editLibraryName")} bind:open={editLibraryModal} size="xs" autoclose>
  <div class="flex flex-col gap-4">
    <div class="flex flex-col gap-1.5">
      <Typography preset="Text sm/medium" color="secondary" text={`${$_("libraryName")}*`} />
      <Input bind:value={tmpLibraryName} placeholder={$_("whatIsYourLibraryTitle")} />
    </div>
    <div class="flex flex-col gap-1.5">
      <Typography preset="Text sm/medium" color="secondary" text={`${$_("description")}*`} />
      <Input bind:value={tmpDescLibrary} placeholder={$_("egDescription")} />
    </div>
  </div>
  <svelte:fragment slot="footer">
    <Button
      onclick={editLibraryName}
      disabled={!tmpLibraryName || !tmpDescLibrary}
      class="me-2 w-full">{$_("save")}</Button
    >
  </svelte:fragment>
</Modal>

<Modal bind:open={deleteModal} size="xs" autoclose>
  <div class="text-center">
    <ExclamationCircleSolid class="mx-auto mb-4 text-red-600 w-12 h-12" />
    <h3 class="mb-5 text-lg font-normal text-gray-500 dark:text-gray-400">
      {$_("areYouSureYouWantToDelete")}
    </h3>
    <Button onclick={deleteDocumentCollection} color="red" class="me-2">{$_("yesImSure")}</Button>
    <Button color="alternative" onclick={() => (deleteModal = false)}>{$_("noCancel")}</Button>
  </div>
</Modal>

<Modal bind:open={deleteLibraryModal} size="xs" autoclose>
  <div class="text-center">
    <ExclamationCircleSolid class="mx-auto mb-4 text-red-600 w-12 h-12" />
    <h3 class="mb-5 text-lg font-normal text-gray-500 dark:text-gray-400">
      {$_("areYouSureYouWantToDelete")}
    </h3>
    <Button onclick={deleteLibrary} color="red" class="me-2">{$_("yesImSure")}</Button>
    <Button color="alternative" onclick={() => (deleteLibraryModal = false)}
      >{$_("noCancel")}</Button
    >
  </div>
</Modal>
