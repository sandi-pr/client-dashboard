<!--
  Copyright (C) 2024 - Hexacode TI
  All Rights Reserved.

  Author(s):
  * Dean Ramhan <dean@hexacode.org>
  * Taqie Fadlillah <taqief@hexacode.org>
  * Sandi Permana <sandy@hexacode.org>

  This file is part of the "Crossroads" project.
-->

<script>
  import {
    P,
    Card,
    Input,
    Button,
    Checkbox,
    Textarea,
    Avatar,
    Modal,
    Toggle,
  } from "flowbite-svelte";
  import { onMount } from "svelte";
  import { _ } from "svelte-i18n";
  import Divider from "$lib/components/Divider.svelte";
  import {
    call,
    SelectedCard,
    UploadFile,
    triggerToast,
    Typography,
    InputRw,
    TextareaRw,
    StateCard,
  } from "$lib";
  import DataTableListCollection from "$lib/components/DataTableListCollection.svelte";
  import { goto } from "$app/navigation";
  import { ChevronDownOutline } from "flowbite-svelte-icons";
  import IllustrationEmpty from "$lib/icon/illustration_empty.svelte";
  import IllustrationProgress from "$lib/icon/illustration_progress.svelte";
    import { tr } from "date-fns/locale";

  const isFrontend = process.env.NODE_ENV === "development";

  let dropdownCollectionOpen = false;
  let collections = $state([]);
  let initCollections = $state([]);
  let documents = $state([]);

  let botName = $state("");
  let description = $state("");
  let fileImage = $state(null);
  let fileImageName = $state(null);
  let fileImageResult = $state(null);
  let selectedplatforms = $state(null);
  let selectedCollection = $state("");
  let additionalPrompt = $state("");
  let needAuth = $state(false);
  let modalOpen = $state(false);
  let isSelectedLibrary = $state(false);
  let allData = $state([]);
  let searchTerm = $state("");
  let search = $state(false);
  let page = $state(1);
  let datasChecked = $state([]);
  let configs = $state({});
  let gettingDatas = $state(false);
  let pagination = $state({ page: 1, pageSize: 10, pageCount: 0, total: 0 });

  let platforms = $state([]);

  let selectedDays = [
    {
      name: "Monday",
      label: "M",
      selected: false,
    },
    {
      name: "Tuesday",
      label: "T",
      selected: false,
    },
    {
      name: "Wednesday",
      label: "W",
      selected: true,
    },
    {
      name: "Thursday",
      label: "T",
      selected: true,
    },
    {
      name: "Friday",
      label: "F",
      selected: true,
    },
    {
      name: "Saturday",
      label: "S",
      selected: false,
    },
    {
      name: "Sunday",
      label: "S",
      selected: false,
    },
  ];

  let screenWidth = $state(0);
  let params = [];
  let products_and_services = [];

  let templates = [
    {
      id: 1,
      attributes: {
        name: "Template 1",
        parameters: [
          {
            children: [
              {
                label: "Parameter 1",
                required: true,
                type: "text",
                placeholder: "Enter Parameter 1",
                enum: ["Option 1", "Option 2", "Option 3"],
                description: "Description for Parameter 1",
              },
            ],
          },
          {
            label: "Label 1",
            children: [
              {
                list_dict_format: [
                  {
                    label: "Field 1",
                    required: true,
                    type: "text",
                    placeholder: "Enter Field 1",
                    name: "field1",
                  },
                  {
                    label: "Field 2",
                    required: false,
                    type: "number",
                    placeholder: "Enter Field 2",
                    min: 0,
                    max: 100,
                    name: "field2",
                  },
                ],
              },
            ],
          },
        ],
      },
    },
    {
      id: 2,
      attributes: {
        name: "Template 2",
        parameters: [
          {
            children: [
              {
                label: "Parameter 2",
                required: false,
                type: "number",
                placeholder: "Enter Parameter 2",
                enum: ["Option A", "Option B", "Option C"],
                description: "Description for Parameter 2",
              },
            ],
          },
        ],
      },
    },
    {
      id: 3,
      attributes: {
        name: "Template 3",
        parameters: [
          {
            children: [
              {
                label: "Parameter 3",
                required: true,
                type: "text",
                placeholder: "Enter Parameter 3",
                enum: ["Option X", "Option Y", "Option Z"],
                description: "Description for Parameter 3",
              },
            ],
          },
        ],
      },
    },
  ];

  async function dummyGetCollections() {
    let token = localStorage.getItem("token");
    if (token !== null) {
      if (!search) {
        gettingDatas = true;
      }
      let response = await $call.call("client-dashboard", "get-documents-by-collectionid", {
        URLQuery: {},
        URLParams: {},
        Params: {
          page: page,
          page_size: pagination.pageSize,
          collection_id: '137',
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

      if (response.result.data.length > 0) {
        pagination = response.result.meta.pagination;
        configs = {
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
                handleCheckAll(true);
              },
            },
            {
              id: "select_all",
              name: "Select All",
              type: "brand-inverse",
              onclick: () => {
                handleCheckAll(false);
              },
            },
          ],
        };

        let promises = response.result.data.map(async (document) => {
          return {
            id: document.id,
            collection_name: document.attributes.document_title,
          };
        });

        allData = await Promise.all(promises);
        console.log("all data", allData);
        datasChecked = allData.map((data) => {
          return { id: data.id, checked: false };
        });
      } else {
        allData = [];
        datasChecked = [];
        configs = {
          search: true,
          searchHandler: handleSearch,
          headers: "File Uploaded",
          pagination: response.result.meta.pagination,
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

  const handleSearch = async (value) => {
    search = true;
    page = 1;
    searchTerm = value;

    // await getDocuments();
  };

  function handleDragOver(event) {
    event.preventDefault();
    dropzoneActive = true;
  }

  function handleDragLeave() {
    dropzoneActive = false;
  }

  async function handleDrop(event) {
    event.preventDefault();
    dropzoneActive = false;

    const files = event.dataTransfer.files;

    if (files[0] && (files[0].type === "image/jpeg" || files[0].type === "image/png")) {
      fileImageResult = await new Promise((resolve, reject) => {
        const reader = new FileReader();

        reader.onload = () => resolve(reader.result);
        reader.onerror = reject;

        reader.readAsDataURL(files[0]);
      });

      fileImage = await readFileAsBase64(files[0]);
      fileImageName = files[0].name;
    } else {
      alert("Invalid file. Please select a valid .jpg or .png file.");
    }
  }

  const handleInputFile = () => {
    const inputElement = document.createElement("input");
    inputElement.type = "file";
    inputElement.accept = "image/jpeg, image/png";

    inputElement.addEventListener("change", async (event) => {
      const file = event.target.files[0];

      if (file && (file.type === "image/jpeg" || file.type === "image/png")) {
        fileImageResult = await new Promise((resolve, reject) => {
          const reader = new FileReader();

          reader.onload = () => resolve(reader.result);
          reader.onerror = reject;

          reader.readAsDataURL(file);
        });

        fileImage = await readFileAsBase64(file);
        fileImageName = file.name;
      } else {
        alert("Invalid file. Please select a valid .jpg or .png file.");
      }
    });

    inputElement.click();
  };

  async function readFileAsBase64(file) {
    return new Promise((resolve, reject) => {
      const reader = new FileReader();

      reader.onload = () => {
        const base64Data = reader.result.split(",")[1];
        resolve(base64Data);
      };
      reader.onerror = reject;

      reader.readAsDataURL(file);
    });
  }

  async function getCollections() {
    let token = localStorage.getItem("token");

    if (token !== null) {
      let page = 1;
      initCollections = [];
      while (true) {
        let response = await $call.call("client-dashboard", "get-document-collections", {
          URLQuery: {},
          URLParams: {},
          Params: {
            page: page,
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
          page++;
        } else {
          break;
        }
      }

      collections = initCollections.map((collection) => {
        let totalDocuments = documents.filter(
          (document) => document.attributes.collection_id === collection.id,
        ).length;
        return {
          value: collection.id,
          id: collection.id,
          name: collection.attributes.collection_name,
          description: collection.attributes.description,
          totalDocuments: totalDocuments,
        };
      });
    }
  }

  async function getModels() {
    let token = localStorage.getItem("token");
    if (token) {
      try {
        let response = await $call.call("hexa-ai-rag", "get-models", {
          URLQuery: {},
          URLParams: {},
          Params: {},
          Headers: {
            Authorization: "Bearer " + token,
          },
        });

        console.log("models", response);

        if (response === null) {
          if (isFrontend) {
            platforms = [
              {
                id: 1,
                name: "Llama3.1",
                model: "llama3.1:8b",
                description: "Llama3.1 Develop by Meta",
                create_at: "2024-07-23",
              },
            ];
          } else {
            triggerLoading(false);
            triggerToast("error");
          }
          return;
        }

        if (response.status === 200) {
          let i = 0;
          let promise = response.result.models.map((model) => {
            return {
              id: i++,
              name: model.name,
              model: model.model,
              description: model.description,
              create_at: model.create_at,
            };
          });

          platforms = await Promise.all(promise);
        } else {
          triggerToast("error");
          console.log(response);
        }
      } catch (error) {
        console.log(error);
      }
    }
  }

  async function createBot() {
    let token = localStorage.getItem("token");

    if (token !== null) {
      try {
        const response = await $call.call("client-dashboard", "create-bot", {
          URLQuery: {},
          URLParams: {},
          Params: {
            bot_name: botName,
            description: description,
            bot_logo: JSON.stringify({ name: fileImageName, base64Data: fileImage }),
            ai_engine: platforms.find((platform) => platform.id === selectedplatforms).model,
            collection_id: selectedCollection,
            additional_prompt: additionalPrompt,
            need_auth: needAuth,
          },
          Headers: {
            Authorization: "Bearer " + token,
          },
        });

        console.log(response);

        if (response.status !== 200) {
          console.log(response);
          return;
        }

        goto(await $call.getPage("client-dashboard", "hexabot"));
        triggerToast("created");
      } catch (error) {
        console.log(error);
      }
    } else {
      console.log("no token");
      return;
    }
  }

  async function handleCheckAll(allChecked) {
    datasChecked = datasChecked.map((data) => {
      return { id: data.id, checked: allChecked ? false : true };
    });

    console.log("datasChecked", datasChecked.length);
  }

  async function handleCheck(id) {
    let datasCheckedCopy = datasChecked;
    datasCheckedCopy = datasCheckedCopy.map((data) => {
      if (data.id === id) {
        return { id: data.id, checked: !data.checked };
      }
      return data;
    });
    datasChecked = datasCheckedCopy;
  }

  function areAllUnchecked() {
    return datasChecked.every(item => item.checked === false);
  }

  async function next() {
    if (page >= pagination.pageCount) return;
    page++;
    // getDocuments();
  }

  async function previous() {
    if (page <= 1) return;
    page--;
    // getDocuments();
  }

  onMount(() => {
    getCollections();
    getModels();
    dummyGetCollections();
  });
</script>

<svelte:window bind:innerWidth={screenWidth} />

<div class="m-4 pb-12">
  <div class="flex justify-between">
    <Typography preset="Display sm/Semibold" text={$_("createNewBots")} />
    <div class="flex text-sm gap-3">
      <Button
        on:click={async () => {
          goto(await $call.getPage("client-dashboard", "hexabot"));
        }}
        color="light">{$_("cancel")}</Button
      >
      <Button
        disabled={botName === "" || selectedplatforms === "" || selectedCollection === ""}
        on:click={() => {
          createBot();
        }}
        class="bg-hc-brand-600 hover:bg-hc-brand-700 dark:bg-hc-brand-600 dark:hover:bg-hc-brand-700"
        >{$_("create")}</Button
      >
    </div>
  </div>

  <!-- Content -->
  <div class="mt-6 space-y-6">
    <div class="md:flex">
      <Typography
        preset="Text md/Semibold"
        text={$_("botsInformation")}
        propClass="w-full md:w-1/4"
      />
      <div class="w-full md:w-3/4 grid grid-cols-1 md:grid-cols-4 gap-3">
        <div class="flex flex-col gap-1.5 col-span-4">
          <Typography preset="Text sm/Medium" text={$_("botsName")} color="secondary" />
          <InputRw bind:value={botName} onEdit />
        </div>
        <div class="flex flex-col gap-1.5 col-span-4">
          <Typography preset="Text sm/Medium" text={$_("descriptionOptional")} color="secondary" />
          <TextareaRw bind:value={description} onEdit />
        </div>
        <div class="col-span-4">
          <UploadFile
            title={$_("botsProfilePicture")}
            handleFileUploadWrapper={() => handleInputFile()}
            isPreview
            srcAvatar={fileImageName ? fileImageResult : ""}
          />
        </div>
      </div>
    </div>
    <Divider />
    <div class="md:flex">
      <Typography
        preset="Text md/Semibold"
        text={$_("aiEngines")}
        color="secondary"
        propClass="w-full md:w-1/4"
      />
      <div class="w-full md:w-3/4">
        <SelectedCard
          selectionData={platforms}
          bind:selectedData={selectedplatforms}
          numGrid={3}
          isRadio
        />
      </div>
    </div>
    <Divider />
    <div class="md:flex">
      <Typography
        preset="Text md/Semibold"
        text={$_("knowledgeSource")}
        color="secondary"
        propClass="w-full md:w-1/4"
      />
      <div class="w-full md:w-3/4">
        <div class="flex items-center gap-4">
          <Button
            on:click={() => (modalOpen = true)}
            color="light"
            class={`w-max rounded-lg border border-primary-light dark:border-primary-dark bg-primary-light dark:bg-primary-dark gap-1.5 flex justify-between items-center py-2.5 px-3`}
          >
            <Typography
              preset="Text md/Semibold"
              text={selectedCollection
                ? collections.find((item) => item.value === selectedCollection)?.name
                : $_("configure")}
              color="secondary"
            />
            <ChevronDownOutline
              class="w-3 h-3 ms-2 text-fg-quaternary-light dark:text-fg-quaternary-dark"
            />
          </Button>
        </div>
      </div>
    </div>
    <Divider />
    <div class="md:flex">
      <Typography
        preset="Text md/Semibold"
        text={$_("additionalPrompts")}
        color="secondary"
        propClass="w-full md:w-1/4"
      />
      <div class="w-full md:w-3/4 flex flex-col gap-1.5">
        <Typography preset="Text sm/Medium" text={$_("botsRole")} color="secondary" />
        <TextareaRw bind:value={additionalPrompt} onEdit />
      </div>
    </div>
    <Divider />
    <div class="md:flex">
      <Typography
        preset="Text md/Semibold"
        text={$_("additionalSettings")}
        color="secondary"
        propClass="w-full md:w-1/4"
      />
      <div class="w-full md:w-3/4 grid grid-cols-1 md:grid-cols-4 gap-4">
        <div class="col-span-4">
          <label class="inline-flex items-center cursor-pointer">
            <input type="checkbox" value="" class="hidden peer" bind:checked={needAuth} />
            <div
              class="relative w-11 h-6 bg-tertiary-light dark:bg-tertiary-dark peer-focus:outline-none peer-focus:ring-4 peer-focus:ring-utility-brand-200-light dark:peer-focus:ring-utility-brand-200-dark rounded-full peer peer-checked:after:translate-x-full rtl:peer-checked:after:-translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:start-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all dark:border-gray-600 peer-checked:bg-brand-solid-light"
            ></div>
            <span class="ms-3 text-sm font-medium text-gray-900 dark:text-gray-300"
              >{$_("requiresLogin")}</span
            >
          </label>
        </div>
      </div>
    </div>
  </div>
</div>

<Modal bind:open={modalOpen} outsideclose>
  <svelte:fragment slot="header">
    <div class="flex gap-4 me-6 grow w-full">
      <div class="grow pt-2 border-t-4 border-fg-brand-primary-alt-light dark:border-fg-brand-primary-alt-dark">
        <Typography preset="Text md/Bold" text={$_("chooseDocumentsLibrary")} color={!isSelectedLibrary ? "brand-secondary" : "secondary"} />
        <Typography preset="Text sm/Regular" text={$_("documentLibraryAsThePrimarySource")} color={!isSelectedLibrary ? "brand-tertiary" : "tertiary"} />
      </div>

      <div class={`grow pt-2 border-t-4 ${isSelectedLibrary ? "border-fg-brand-primary-alt-light dark:border-fg-brand-primary-alt-dark" : "border-utility-gray-200-light dark:border-utility-gray-200-dark"}`}>
        <Typography preset="Text md/Bold" text={$_("chooseWhichCollections")} color={isSelectedLibrary ? "brand-secondary" : "secondary"} />
        <Typography preset="Text sm/Regular" text={$_("chooseTheOneYouWantToInclude")} color={isSelectedLibrary ? "brand-tertiary" : "tertiary"} />
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
  {:else}
    {#if allData.length > 0 || searchTerm !== ""}
      {#if !gettingDatas}
        <DataTableListCollection
          bind:datas={allData}
          bind:datasChecked
          bind:configs
          bind:gettingDatas
          bind:next
          bind:previous
          bind:handleCheckAll
          bind:handleCheck
        />
      {:else}
        <StateCard
          illustration={IllustrationProgress}
          description={$_("fetchingDataPleaseWait")}
        />
      {/if}
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
  {/if}
  <svelte:fragment slot="footer">
    {#if !isSelectedLibrary}
    <Button on:click={() => {isSelectedLibrary = true}} class="w-full" disabled={!selectedCollection}>{$_("next")}</Button>
    {:else}
    <div class="flex gap-4 w-full">
      <button onclick={() => {isSelectedLibrary = false}} class="px-5">{$_("goBack")}</button>
      <Button on:click={() => {}} class="grow" disabled={areAllUnchecked()}>{$_("finish")}</Button>
    </div>
    {/if}
  </svelte:fragment>
</Modal>
