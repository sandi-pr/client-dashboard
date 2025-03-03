<script>
  import { _ } from "svelte-i18n";
  import { Button, Modal, Input, Alert, Select, P, Helper } from "flowbite-svelte";

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
    bakrieLibraryId,
    isBakrie,
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
  let tmpOldName = $state("");
  let tmpDescription = $state("");
  let tmpLibraryName = $state("");
  let tmpDescLibrary = $state("");
  let tmpCollectionId = $state("");
  let libraryId = $state("");
  let dropdownOpen = $state(false);
  let nameIsAvailable = $state(true);
  let gettingDatas = $state(false);

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
        description: response.result.attributes?.meta?.description || "",
        createdAt: format(new Date(response.result.attributes.createdAt), "dd-MM-yyyy"),
      };
    }
  }

  async function getCollections() {
    let token = localStorage.getItem("token");

    if (token !== null) {
      gettingDatas = true;
      try {
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

        // id: "4",
        // collectionName: "Study Scriptures of Qur’an Tafsir",
        // totalDocuments: 2,
        // type: "",
        // status: "reading",

        let link = await $call.getPage("client-dashboard", "handwriting-documents");

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
            description: collection.attributes?.meta?.description || "",
            type: collection.attributes.type,
            href: link + "?collectionId=" + collection.id + "&" + queryString,
          };
        });

        collectionListDatas = await Promise.all(promise);
        gettingDatas = false;
      } catch (error) {
        console.error("Error during collection get:", error);
        gettingDatas = false;
        triggerToast($_("failedGetCollection"), true, "error");
      }
    } else {
      triggerToast($_("unauthorized"), true, "error");
      return;
    }
  }

  async function createCollection() {
    let token = localStorage.getItem("token");

    if (token !== null) {
      try {
        let dataSend = {};
        if ($isBakrie) {
          dataSend = {
            collection_name: tmpCollectionName,
            description: tmpDescription,
            parent_id: parseInt(libraryId),
            // bakrie
            bot_id: localStorage.getItem("botIdStrapi").toString(),
          };
        } else {
          dataSend = {
            collection_name: tmpCollectionName,
            description: tmpDescription,
            parent_id: parseInt(libraryId),
          };
        }

        let response = await $call.call("client-dashboard", "create-library-collection", {
          URLQuery: {},
          URLParams: {},
          Params: dataSend,
          Headers: {
            Authorization: "Bearer " + token,
          },
        });

        if (response.status !== 200) {
          console.log(response);
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

        // console.log(response);

        tmpCollectionName = "";
        tmpDescription = "";
        addCollectionModal = false;

        triggerToast($_("successCreateCollection"), true);
        getCollections();
        // goto(await $call.getPage("client-dashboard", "document-library"));
      } catch (error) {
        console.error("Error during collection creation:", error);
        triggerToast($_("failedCreateCollection"), true, "error");
      }
    } else {
      triggerToast($_("unauthorized"), true, "error");
      return;
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

        tmpLibraryName = "";
        tmpDescLibrary = "";

        editLibraryModal = false;
        triggerToast($_("successUpdateLibrary"), true);
        getCollections();
      } catch (error) {
        console.error("Error during library editing:", error);
        triggerToast($_("failedEditLibrary"), true, "error");
      }
    } else {
      triggerToast($_("unauthorized"), true, "error");
      return;
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

        tmpCollectionName = "";
        tmpDescription = "";

        editModal = false;
        triggerToast($_("successUpdateCollection"), true);
        getCollections();
      } catch (error) {
        console.error("Error during collection editing:", error);
        triggerToast($_("failedEditCollection"), true, "error");
      }
    } else {
      triggerToast($_("unauthorized"), true, "error");
      return;
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
        let response = await $call.call("client-dashboard", "delete-library", {
          URLQuery: {},
          URLParams: {
            id: libraryId.toString(),
          },
          Params: {},
          Headers: {
            Authorization: "Bearer " + token,
          },
        });

        if (response.status !== 200) {
          console.log(response);
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

        console.log("ayam", response);

        triggerToast($_("successDeleteLibrary"), true);

        getCollections();
        deleteLibraryModal = false;
      } catch (error) {
        console.error("Error during library deletion:", error);
        triggerToast($_("failedDeleteLibrary"), true, "error");
      }
    } else {
      triggerToast($_("unauthorized"), true, "error");
      return;
    }
  }

  async function deleteDocumentCollection() {
    let token = localStorage.getItem("token");

    if (token !== null) {
      try {
        let response = await $call.call("client-dashboard", "delete-library-collection", {
          URLQuery: {},
          URLParams: {
            id: tmpCollectionId.toString(),
          },
          Params: {},
          Headers: {
            Authorization: "Bearer " + token,
          },
        });

        if (response.status !== 200) {
          console.log(response);
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

        console.log("ayam", response);

        triggerToast($_("successDeleteCollection"), true);

        getCollections();
        deleteModal = false;
      } catch (error) {
        console.error("Error during collection deletion:", error);
        triggerToast($_("failedDeleteCollection"), true, "error");
      }
    } else {
      triggerToast($_("unauthorized"), true, "error");
      return;
    }
  }

  const checkLibraryName = async (name) => {
    let token = localStorage.getItem("token");
    if (!name || !token) {
      return null;
    }

    try {
      let resp = await $call.call("client-dashboard", "is-library-name-available", {
        URLQuery: {},
        URLParams: {},
        Params: {
          name: name,
        },
        Headers: {
          Authorization: "Bearer " + token,
        },
      });

      console.log("cek resp", resp);

      if (resp === null) {
        throw new Error("Server sibuk");
      }

      if (resp.status === 200) {
        return resp.result.available; // Mengembalikan status (true: tersedia, false: tidak tersedia)
      } else {
        console.error("Error memeriksa name:", resp);
        return null;
      }
    } catch (error) {
      console.error("Terjadi kesalahan:", error);
      return null;
    }
  };

  const handleInput = async (event) => {
    const name = event.target.value;
    const status = await checkLibraryName(name.trim());

    if (tmpOldName === name) {
      nameIsAvailable = true;
    } else if (status === false) {
      nameIsAvailable = false;
    } else if (status === true) {
      nameIsAvailable = true;
    }
  };

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
        // goto(await $call.getPage("client-dashboard", "document-library"));
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
        <div class="flex items-start justify-between mb-3 md:mb-8">
          <div class="flex flex-col gap-1.5">
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
                nameIsAvailable = true;
                tmpCollectionName = "";
                tmpDescription = "";
              }}
              class="p-2 md:px-3 md:py-2 gap-2"
              outline
            >
              <FolderPlusIcon svgClass="w-5 h-5 text-quinary-light dark:text-quinary-dark" />
              <div class="hidden md:block">{$_("addCollection")}</div>
            </Button>
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
                  tmpOldName = document.collectionName;
                  tmpDescription = collection.description;
                  nameIsAvailable = true;
                  editModal = true;
                }}
                functDelete={() => {
                  tmpCollectionId = collection.id;
                  tmpCollectionName = collection.collectionName;
                  deleteModal = true;
                }}
                funcDetailCollection={() => {
                  goto(collection.href);
                }}
                {libraryId}
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
        {:else if gettingDatas}
          <StateCard
            illustration={IllustrationProgress}
            description={$_("fetchingDataPleaseWait")}
          />
        {:else}
          <StateCard
            illustration={IllustrationEmpty}
            description={$_("youDontHaveAnyDataYet")}
          />
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
      <Input
        bind:value={tmpCollectionName}
        on:input={handleInput}
        placeholder={$_("whatIsYourCollectionTitle")}
      />
      {#if !nameIsAvailable}
        <Helper color="red">
          <span class="font-medium">{$_("nameAlreadyTaken")}</span>
        </Helper>
      {/if}
    </div>
    <div class="flex flex-col gap-1.5">
      <Typography preset="Text sm/medium" color="secondary" text={`${$_("description")}`} />
      <Input bind:value={tmpDescription} placeholder={$_("egDescription")} />
    </div>
  </div>
  <svelte:fragment slot="footer">
    <Button
      onclick={createCollection}
      disabled={!tmpCollectionName || !nameIsAvailable}
      class="w-full">{$_("addNewCollection")}</Button
    >
  </svelte:fragment>
</Modal>

<Modal title={$_("editCollectionName")} bind:open={editModal} size="xs" autoclose>
  <div class="flex flex-col gap-4">
    <div class="flex flex-col gap-1.5">
      <Typography preset="Text sm/medium" color="secondary" text={`${$_("collectionName")}*`} />
      <Input
        bind:value={tmpCollectionName}
        on:input={handleInput}
        placeholder={$_("whatIsYourCollectionTitle")}
      />
      {#if !nameIsAvailable}
        <Helper color="red">
          <span class="font-medium">{$_("nameAlreadyTaken")}</span>
        </Helper>
      {/if}
    </div>
    <div class="flex flex-col gap-1.5">
      <Typography preset="Text sm/medium" color="secondary" text={`${$_("description")}`} />
      <Input bind:value={tmpDescription} placeholder={$_("egDescription")} />
    </div>
  </div>
  <svelte:fragment slot="footer">
    <Button
      onclick={editCollectionName}
      disabled={!tmpCollectionName || !nameIsAvailable}
      class="me-2 w-full">{$_("save")}</Button
    >
  </svelte:fragment>
</Modal>

<Modal title={$_("editLibraryName")} bind:open={editLibraryModal} size="xs" autoclose>
  <div class="flex flex-col gap-4">
    <div class="flex flex-col gap-1.5">
      <Typography preset="Text sm/medium" color="secondary" text={`${$_("libraryName")}*`} />
      <Input
        bind:value={tmpLibraryName}
        on:input={handleInput}
        placeholder={$_("whatIsYourLibraryTitle")}
      />
      {#if !nameIsAvailable}
        <Helper color="red">
          <span class="font-medium">{$_("nameAlreadyTaken")}</span>
        </Helper>
      {/if}
    </div>
    <div class="flex flex-col gap-1.5">
      <Typography preset="Text sm/medium" color="secondary" text={`${$_("description")}`} />
      <Input bind:value={tmpDescLibrary} placeholder={$_("egDescription")} />
    </div>
  </div>
  <svelte:fragment slot="footer">
    <Button
      onclick={editLibraryName}
      disabled={!tmpLibraryName || !nameIsAvailable}
      class="me-2 w-full">{$_("save")}</Button
    >
  </svelte:fragment>
</Modal>

<Modal bind:open={deleteModal} size="xs" autoclose>
  <div class="text-center">
    <ExclamationCircleSolid class="mx-auto mb-4 text-red-600 w-12 h-12" />
    <h3 class="mb-5 text-lg font-normal text-gray-500 dark:text-gray-400">
      {$_("areYouSureYouWantToDelete")}
      {$_("collection")}
      "{tmpCollectionName}"?
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
      {$_("library")}
      {collectionDatas.title}?
    </h3>
    <Button onclick={deleteLibrary} color="red" class="me-2">{$_("yesImSure")}</Button>
    <Button color="alternative" onclick={() => (deleteLibraryModal = false)}
      >{$_("noCancel")}</Button
    >
  </div>
</Modal>
