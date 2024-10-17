<script>
  import { onMount } from "svelte";
  import { _ } from "svelte-i18n";
  import { Button, Modal, Input, Alert, Select, P } from "flowbite-svelte";

  import { Typography, ToastContainer, StateCard } from "$lib";
  import DataTableLibraryCollection from "$lib/components/DataTableLibraryCollection.svelte";
  import { SearchOutline } from "flowbite-svelte-icons";
  import Trash01Icon from "$lib/icon/trash-01.svelte";
  import FilePlus02Icon from "$lib/icon/file-plus-02.svelte";
  import IllustrationEmpty from "$lib/icon/illustration_empty.svelte";
  import IllustrationProgress from "$lib/icon/illustration_progress.svelte";

  let collectionDatas = {
    title: "My Collections",
    description:
      "You can manage your collections here. You can add, edit, and delete your collections. You can also see the details of your collections.",
    createdAt: "2021-08-12",
  };
  let collectionListDatas = [
    {
      id: "1",
      filename: "Bakrie’s Handwriting.png",
      fileext: "png",
      label: ["label", "lebal", "lebel", "lebil", "lebol", "lebul", "lebzl"],
      type: "Handwriting",
      filesize: 1000000,
      createdAt: "2021-08-12",
    },
    {
      id: "2",
      filename: "Study Scriptures of Qur’an Tafsir.jpg",
      fileext: "pdf",
      label: ["label", "lebal"],
      type: "Handwriting",
      filesize: 1000000,
      createdAt: "2021-08-12",
    },
    {
      id: "3",
      filename: "Study Scriptures of Qur’an Tafsir.jpg",
      fileext: "jpeg",
      label: ["label", "lebal"],
      type: "Handwriting",
      filesize: 100000000,
      createdAt: "2021-08-12",
    },
    {
      id: "4",
      filename: "Study Scriptures of Qur’an Tafsir.jpg",
      fileext: "doc",
      label: ["label", "lebal"],
      type: "Handwriting",
      filesize: 1000000,
      createdAt: "2021-08-12",
    },
    {
      id: "5",
      filename: "Study Scriptures of Qur’an Tafsir.jpg",
      fileext: "docx",
      label: ["label", "lebal"],
      type: "Handwriting",
      filesize: 1000000,
      createdAt: "2021-08-12",
    },
    {
      id: "6",
      filename: "Study Scriptures of Qur’an Tafsir.jpg",
      fileext: "jpg",
      label: ["label", "lebal"],
      type: "Handwriting",
      filesize: 1000000,
      createdAt: "2021-08-12",
    },
    {
      id: "7",
      filename: "Study Scriptures of Qur’an Tafsir.jpg",
      fileext: "jpg",
      label: ["label", "lebal"],
      type: "Handwriting",
      filesize: 1500000,
      createdAt: "2021-08-12",
    },
    {
      id: "8",
      filename: "Study Scriptures of Qur’an Tafsir.jpg",
      fileext: "jpg",
      label: ["label", "lebal"],
      type: "Handwriting",
      filesize: 42000000,
      createdAt: "2021-08-12",
    },
    {
      id: "9",
      filename: "Study Scriptures of Qur’an Tafsir.jpg",
      fileext: "jpg",
      label: ["label", "lebal"],
      type: "Handwriting",
      filesize: 1000000,
      createdAt: "2021-08-12",
    },
    {
      id: "10",
      filename: "Study Scriptures of Qur’an Tafsir.jpg",
      fileext: "jpg",
      label: ["label", "lebal"],
      type: "Handwriting",
      filesize: 1000000,
      createdAt: "2021-08-12",
    },
  ];
  let metadata = {
    page: 1,
    pageSize: 10,
    total: 100,
  };

  const checkStateAdder = () => {
    collectionListDatas = collectionListDatas.map((data) => {
      return {
        ...data,
        checked: false,
      };
    });
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
    console.log(collectionListDatas);
  };

  const checkAll = () => {
    collectionListDatas = collectionListDatas.map((data) => {
      return {
        ...data,
        checked: !data.checked,
      };
    });
  };

  const nextPage = (evt) => {
    console.log("next page");
  };

  const prevPage = (evt) => {
    console.log("prev page");
  };

  const handleSearch = (evt) => {
    console.log("clicked");
  };

  $: checkAllState = collectionListDatas.every((data) => data.checked === true);

  onMount(() => {
    if (collectionListDatas.length > 0) {
      checkStateAdder();
      console.log(collectionListDatas);
    }
  });
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
            </div>
          </div>

          <div class="flex gap-4">
            <Input
              id="search"
              type="search"
              class="bg-white dark:bg-hc-graydark-950"
              placeholder={$_("search")}
              on:change={handleSearch}
            >
              <SearchOutline
                slot="left"
                class="w-5 h-5 text-quinary-light dark:text-quinary-dark"
              />
            </Input>
            <Button
              on:click={() => {
                // addLibraryModal = true;
                // tmpCollectionName = "";
                // tmpDescription = "";
              }}
              class="p-2 md:px-3 md:py-2 gap-2 text-fg-secondary-light dark:text-fg-secondary-dark border-secondary-light dark:border-secondary-dark hover:bg-secondary-light hover:dark:bg-secondary-dark hover:text-fg-secondary-light dark:hover:text-fg-secondary-dark"
              outline
            >
              <Trash01Icon svgClass="w-5 h-5" />
              <div class="hidden md:block">{$_("delete")}</div>
            </Button>
            <Button
              on:click={() => {
                // addLibraryModal = true;
                // tmpCollectionName = "";
                // tmpDescription = "";
              }}
              class="p-2 md:px-3 md:py-2 gap-2 whitespace-nowrap"
            >
              <FilePlus02Icon svgClass="w-5 h-5 text-quinary-light dark:text-quinary-dark" />
              <div class="hidden md:block">{$_("addNew")}</div>
            </Button>
          </div>
        </div>

        <div
          class="mx-4 sm:mx-0 border-t-[1px] border-hc-graylight-300 dark:border-hc-graydark-700"
        ></div>

        {#if collectionListDatas.length > 0}
          <DataTableLibraryCollection
            {collectionListDatas}
            {metadata}
            {checkAllState}
            {checkSingle}
            {checkAll}
            next={nextPage}
            previous={prevPage}
          />
        {:else}
          <StateCard illustration={IllustrationEmpty} description={$_("youDontHaveAnyDataYet")} />
        {/if}
      </div>
    </div>
  </div>
  <ToastContainer />
</div>
