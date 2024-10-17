<script>
  import { _ } from "svelte-i18n";
  import { Button, Modal, Input, Alert, Select, P } from "flowbite-svelte";

  import { Typography, ToastContainer, StateCard, triggerToast, InputRw, TextareaRw } from "$lib";
  import TabBar from "$lib/components/TabBar.svelte";

  import PlusIcon from "$lib/icon/plus.svg";
  import Copy01Icon from "$lib/icon/copy-01.svelte";
  import DotsVerticalIcon from "$lib/icon/dots-vertical.svelte";
  import FolderPlusIcon from "$lib/icon/folder-plus.svelte";
  import Stars03Icon from "$lib/icon/stars-03.svelte";
  import Folder from "$lib/icon/folder.svelte";
  import IllustrationEmpty from "$lib/icon/illustration_empty.svelte";
  import IllustrationProgress from "$lib/icon/illustration_progress.svelte";

  let addCollectionModal = false;
  let tmpCollectionName = "";
  let tmpDescription = "";
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
          <div class="grid grid-cols-4 gap-4">
            {#each collectionListDatas as collection}
              <div
                class="flex justify-between items-center p-4 rounded-md border-[1px] bg-primary-light dark:bg-primary-dark border-secondary-light dark:border-secondary-dark mt-4 cursor-pointer"
              >
                <div class="flex items-center gap-2">
                  <Folder svgClass="w-5 h-5 text-quarterary-light dark:text-quarterary-dark" />
                  <Typography
                    preset="Text sm/Semibold"
                    color="secondary"
                    text={collection.collectionName}
                  />
                </div>
                <DotsVerticalIcon
                  svgClass="w-5 h-5 text-quarterary-light dark:text-quarterary-dark cursor-pointer"
                />
              </div>
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
