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
  } from "$lib";
  import { goto } from "$app/navigation";
  import { ChevronDownOutline } from "flowbite-svelte-icons";

  const isFrontend = process.env.NODE_ENV === "development";

  let dropdownCollectionOpen = false;
  let collections = [];
  let initCollections = [];
  let documents = [];

  let botName = "";
  let description = "";
  let fileImage = null;
  let fileImageName = null;
  let fileImageResult = null;
  let selectedplatforms = null;
  let selectedCollection = "";
  let additionalPrompt = "";
  let needAuth = false;
  let modalOpen = false;

  let platforms = [];

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

  let screenWidth;
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
    let token = sessionStorage.getItem("token");

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
    let token = sessionStorage.getItem("token");
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
          }
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
          console.log(response);
        }
      } catch (error) {
        console.log(error);
      }
    }
  }

  async function createBot() {
    let token = sessionStorage.getItem("token");

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

  onMount(() => {
    getCollections();
    getModels();
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
        text={$_("documentCollection")}
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
                : $_("chooseCollection")}
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
          <Toggle checked={needAuth} onclick={() => (needAuth = !needAuth)} class="w-max"
            >{$_("requiresLogin")}</Toggle
          >
        </div>
      </div>
    </div>
  </div>
</div>

<Modal title={$_("chooseSource")} bind:open={modalOpen} autoclose outsideclose>
  <div class="flex flex-col justify-center gap-4">
    <div class="flex flex-col justify-center mt-3 gap-3">
      <Typography preset="Text sm/Medium" text={$_("recentCollections")} color="quaterary" />

      <SelectedCard
        selectionData={collections}
        bind:selectedData={selectedCollection}
        numGrid={2}
        isRadio
      />
    </div>
  </div>
  <svelte:fragment slot="footer">
    <Button on:click={() => {}} class="w-full" disabled={!selectedCollection}>{$_("select")}</Button
    >
  </svelte:fragment>
</Modal>
