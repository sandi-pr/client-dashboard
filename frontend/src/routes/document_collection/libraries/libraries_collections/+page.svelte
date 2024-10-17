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
  } from "$lib";
  import TabBar from "$lib/components/TabBar.svelte";

  import { ExclamationCircleSolid } from "flowbite-svelte-icons";
  import Copy01Icon from "$lib/icon/copy-01.svelte";
  import DotsVerticalIcon from "$lib/icon/dots-vertical.svelte";
  import FolderPlusIcon from "$lib/icon/folder-plus.svelte";
  import Stars03Icon from "$lib/icon/stars-03.svelte";
  import Folder from "$lib/icon/folder.svelte";
  import IllustrationEmpty from "$lib/icon/illustration_empty.svelte";
  import IllustrationProgress from "$lib/icon/illustration_progress.svelte";

  let addCollectionModal = false;
  let editModal = false;
  let deleteModal = false;
  let tmpCollectionName = "";
  let tmpDescription = "";
  let tmpCollectionId = "";
  let collectionDatas = {
    title: "My Collections",
    description:
      "You can manage your collections here. You can add, edit, and delete your collections. You can also see the details of your collections.",
    createdAt: "2021-08-12",
  };
  let tabDatas = {
    selected: "My Collections",
    data: [
      {
        name: "My Collections",
      },
      {
        name: "Information",
      },
    ],
  };
  let collectionListDatas = [
    {
      id: "1",
      collectionName: "Bakrie’s Handwriting",
      totalDocuments: 2,
      type: "",
      status: "ready",
    },
    {
      id: "2",
      collectionName: "Study Scriptures of Qur’an Tafsir",
      totalDocuments: 2,
      type: "",
      status: "reading",
    },
    {
      id: "3",
      collectionName: "Study Scriptures of Qur’an Tafsir",
      totalDocuments: 2,
      type: "",
      status: "reading",
    },
    {
      id: "4",
      collectionName: "Study Scriptures of Qur’an Tafsir",
      totalDocuments: 2,
      type: "",
      status: "reading",
    },
    {
      id: "5",
      collectionName: "Study Scriptures of Qur’an Tafsir",
      totalDocuments: 2,
      type: "",
      status: "reading",
    },
    {
      id: "6",
      collectionName: "Study Scriptures of Qur’an Tafsir",
      totalDocuments: 2,
      type: "",
      status: "reading",
    },
    {
      id: "7",
      collectionName: "Study Scriptures of Qur’an Tafsir",
      totalDocuments: 2,
      type: "",
      status: "reading",
    },
    {
      id: "8",
      collectionName: "Study Scriptures of Qur’an Tafsir",
      totalDocuments: 2,
      type: "",
      status: "reading",
    },
    {
      id: "9",
      collectionName: "Study Scriptures of Qur’an Tafsir",
      totalDocuments: 2,
      type: "",
      status: "reading",
    },
    {
      id: "10",
      collectionName: "Study Scriptures of Qur’an Tafsir",
      totalDocuments: 2,
      type: "",
      status: "reading",
    },
  ];
  let informationDatas = {
    id: "1",
    collectionName: "Bakrie’s Handwriting",
    totalDocuments: 2,
    type: "",
    status: "ready",
  };

  const createCollection = (evt) => {
    console.log("clicked");
  };

  async function editCollectionName() {
    let token = localStorage.getItem("token");

    if (token !== null) {
      try {
        let response = await $call.call("client-dashboard", "update-document-collection", {
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

  async function deleteDocumentCollection() {
    let token = localStorage.getItem("token");

    if (token !== null) {
      try {
        for (const document of documents) {
          if (document.attributes.collection_id === tmpCollectionId) {
            let response = await $call.call("client-dashboard", "delete-document", {
              URLQuery: {},
              URLParams: {
                id: document.id.toString(),
              },
              Params: {},
              Headers: {
                Authorization: "Bearer " + token,
              },
            });
            if (response.status !== 200) {
              console.log(response);
              triggerToast("error");
              return;
            }
          }
        }
        let response2 = await $call.call("client-dashboard", "delete-document-collection", {
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

          <div class="flex gap-4">
            <Button
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
            <Button
              on:click={() => {
                // addLibraryModal = true;
                // tmpCollectionName = "";
                // tmpDescription = "";
              }}
              class="p-2 md:px-3 md:py-2 gap-2"
            >
              <Stars03Icon svgClass="w-5 h-5 text-quinary-light dark:text-quinary-dark" />
              <div class="hidden md:block">{$_("startReading")}</div>
            </Button>
            <Button
              on:click={() => {
                // addLibraryModal = true;
                // tmpCollectionName = "";
                // tmpDescription = "";
              }}
              class="p-2 md:px-3 md:py-2"
            >
              <DotsVerticalIcon svgClass="w-5 h-5 text-quinary-light dark:text-quinary-dark" />
            </Button>
          </div>
        </div>

        <TabBar {tabDatas} on:tabChange={(e) => (tabDatas.selected = e.detail)} />

        {#if tabDatas.selected === "My Collections" && collectionListDatas.length > 0}
          <div class="grid sm:grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4 mt-4">
            {#each collectionListDatas as collection}
              <LibraryCard
                collectionName={collection.collectionName}
                totalDocuments={collection.totalDocuments}
                status={collection.status}
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
                  // tmpCollectionId = collection.id;
                  // tmpHref = collection.href;
                  // categories = collection.categories;
                  // selectTypeCollection("Documents");
                  // categories = collection.categories;
                  // otherCategory = "";
                  // addOtherCategory = false;
                  // addDocumentsModal = true;
                }}
                functEdit={() => {
                  // tmpCollectionId = collection.id;
                  // tmpCollectionName = collection.collectionName;
                  // tmpDescription = collection.description;
                  editModal = true;
                }}
                functDelete={() => {
                  // tmpCollectionId = collection.id;
                  deleteModal = true;
                }}
                funcDetailCollection={() => {
                  // goto(collection.href);
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
      <Typography preset="Text sm/medium" color="secondary" text={$_("description")} />
      <Input bind:value={tmpDescription} placeholder={$_("egDescription")} />
    </div>
  </div>
  <svelte:fragment slot="footer">
    <Button onclick={createCollection} disabled={!tmpCollectionName} class="w-full"
      >{$_("addNewCollection")}</Button
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
    <Button onclick={editCollectionName} class="me-2 w-full">{$_("save")}</Button>
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
