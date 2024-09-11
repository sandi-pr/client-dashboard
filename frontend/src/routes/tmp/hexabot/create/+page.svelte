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
    // Dropdown,
    // DropdownItem,
  } from "flowbite-svelte";
  import { ChevronDownOutline } from "flowbite-svelte-icons";
  import WebIcon from "$lib/icon/server-06.svg";
  import WhatsappIcon from "$lib/icon/whatsapp.svg";
  // import TelegramIcon from "$lib/icon/telegram.svg";
  import LoadingSpinner from "$lib/components/LoadingSpinner.svelte";
  import { onMount } from "svelte";
  import { _ } from "svelte-i18n";
  import upload from "$lib/icon/upload-cloud-02.svg";
  import { call, Typography, SelectedCard, Dropdown } from "$lib";
  // import axios from "axios";
  // import { getTemplates } from "$lib/service/strapi";

  // export let data;

  let selectedplatforms = null;
  let dropdownCollectionOpen = false;
  let selectedCollection = "";
  let collections = [];
  let initCollections = [];
  let documents = [];

  let placement = "bottom-start";

  let platforms = [
    {
      id: 1,
      name: "LLama 3.1",
    },
  ];

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

  let fileImage = null;
  let fileImageName = null;
  let fileImageResult = null;

  // let product = {
  //   name: "CRM for Small Business",
  //   moduleId: "1454523",
  //   lastVersion: "1.2",
  //   status: { name: "Active", type: "success" },
  //   description:
  //     "Bank Jasa Jakarta is a bank that provides various financial services, including savings, loans, and investment services.",
  //   activePeriod: "Currently Running",
  //   chargingPeriod: "Monthly",
  // };

  // let componenDatas = [
  //   {
  //     id: 1,
  //     companyName: "CRM for Small Business",
  //     modulesUsed: "CRM",
  //     featuresUsed: "Profiles, Features 2, Features 3",
  //   },
  //   {
  //     id: 2,
  //     companyName: "CRM for Medium Business",
  //     modulesUsed: "CRM",
  //     featuresUsed: "Profiles, Features 2, Features 3",
  //   },
  //   {
  //     id: 3,
  //     companyName: "CRM for Larger Enterprises",
  //     modulesUsed: "CRM",
  //     featuresUsed: "Profiles, Features 2, Features 3",
  //   },
  // ];

  // let UsageDatas = [
  //   {
  //     id: 1,
  //     features: "Profiles",
  //     lowerLimit: "100 Users",
  //     upperLimit: "1000 Users",
  //     Pricing: "$0.25/user",
  //   },
  //   {
  //     id: 2,
  //     features: "Features 2",
  //     lowerLimit: "100 Users",
  //     upperLimit: "1000 Users",
  //     Pricing: "$0.25/user",
  //   },
  //   {
  //     id: 3,
  //     features: "Features 3",
  //     lowerLimit: "100 Users",
  //     upperLimit: "1000 Users",
  //     Pricing: "$0.25/user",
  //   },
  // ];

  // let formattedDate = "Created on 2021-08-12 14:00 PM";

  // let submenu = "information";
  let isLoading = false;
  let screenWidth;

  // let files = [];

  // function handleFileChange(event) {
  //   const selectedFiles = event.target.files;

  //   for (const file of selectedFiles) {
  //     if (file.size > 1 * 1024 * 1024) {
  //       alert(`File "${file.name}" size exceeds the 10 MB limit. Please choose a smaller file.`);
  //       return;
  //     }
  //     if (file) {
  //       const reader = new FileReader();
  //       reader.onloadend = () => {
  //         // pdf = reader.result;
  //         // base64 = reader.result.split(",")[1];
  //       };
  //       reader.readAsDataURL(file);
  //       files.push(file);
  //     }
  //   }
  //   console.log(files);
  //   files = [...files];
  // }

  // function formatBytes(bytes, decimals = 2) {
  //   if (bytes === 0) return "0 Bytes";

  //   const k = 1024;
  //   const dm = decimals < 0 ? 0 : decimals;
  //   const sizes = ["Bytes", "KB", "MB", "GB", "TB", "PB", "EB", "ZB", "YB"];

  //   const i = Math.floor(Math.log(bytes) / Math.log(k));

  //   return parseFloat((bytes / Math.pow(k, i)).toFixed(dm)) + " " + sizes[i];
  // }

  // function removeFile(index) {
  //   files = files.filter((_, i) => i !== index);
  // }

  // let templates = [];
  let dropdownIndustryOpen = false;
  let dropdownUseCaseOpen = false;
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
                  { label: "Field 1", required: true, type: "text", placeholder: "Enter Field 1", name: "field1" },
                  { label: "Field 2", required: false, type: "number", placeholder: "Enter Field 2", min: 0, max: 100, name: "field2" },
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

  let selectedIndustry = "";
  let selectedUseCase = "";

  // const getTemplate = async () => {
  //   const result = await getTemplates(data.data.jwt);

  //   if (result.data.length > 0) {
  //     templates = result.data;
  //     params = templates[0].attributes.parameters;
  //     const initialProduct = {};
  //     params[1].children[0].list_dict_format.forEach((field) => {
  //       initialProduct[field.name] = "";
  //     });
  //     products_and_services.push(initialProduct);
  //   }
  // };

  const addProduct = () => {
    const newProduct = {};
    params[1]?.children[0]?.list_dict_format.forEach((field) => {
      newProduct[field.name] = "";
    });
    products_and_services = [...products_and_services, newProduct];
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

      fileImage = files[0];
      fileImageName = files[0].name;
      // console.log(fileImageResult);
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

        // conver to base64
        fileImage = await readFileAsBase64(file);
        fileImageName = file.name;
        // console.log(fileImageResult);
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

        let page2 = 1;

        collections = initCollections.map((collection) => {
          let totalDocuments = documents.filter(
            (document) => document.attributes.collection_id === collection.id,
          ).length;
          return {
            id: collection.id,
            name: collection.attributes.collection_name,
            description: collection.attributes.description,
            totalDocuments: totalDocuments,
          };
        });
      }
    }
  }

  onMount(() => {
    getCollections();
  });

</script>

<svelte:window bind:innerWidth={screenWidth} />

<!-- <LoadingSpinner {isLoading} /> -->

<div class="m-4">
  <div class="flex justify-between">
    <div class="ml-4">
      <P size="3xl" weight="semibold">Create New Bots</P>
    </div>
    <div class="flex text-sm gap-3">
      <Button on:click={() => {}} color="light">Cancel</Button>
      <Button
        on:click={() => {}}
        class="bg-hc-brand-600 hover:bg-hc-brand-700 dark:bg-hc-brand-600 dark:hover:bg-hc-brand-700"
        >Finish</Button
      >
    </div>
  </div>

  <!-- Content -->
  <div class="mx-4 mt-8 pb-4">
    <div class="md:flex">
      <P size="lg" weight="semibold" class="w-full md:w-1/4 pb-2">Bots Information</P>

      <div class="w-full md:w-3/4 grid grid-cols-1 md:grid-cols-4 gap-4">
        <div
          class="flex flex-col gap-1 col-span-4 [&>*]:bg-primary-light [&>*]:dark:bg-primary-dark [&>textarea]:text-primary-light dark:[&>textarea]:text-primary-dark"
        >
          <P size="sm" weight="medium" class="text-secondary-light dark:text-secondary-dark"
            >Bots Name*</P
          >
          <Textarea
            class="resize-none"
            rows={screenWidth > 1024 ? "1" : "2"}
            unWrappedClass="p-2.5 text-sm border border-primary-light dark:border-primary-dark focus:ring-primary-500 focus:border-primary-500 dark:focus:ring-primary-500 dark:focus:border-primary-500"
          />
        </div>
        <div
          class="space-y-1 col-span-4 [&>*]:bg-primary-light [&>*]:dark:bg-primary-dark [&>textarea]:text-primary-light dark:[&>textarea]:text-primary-dark"
        >
          <P size="sm" weight="medium" class="text-secondary-light dark:text-secondary-dark"
            >Description (Opsional)</P
          >
          <Textarea class="resize-none" />
        </div>
        <div class="flex flex-col col-span-4 w-full space-y-1">
          <div class="">
            <P size="sm" weight="medium" class="text-secondary-light dark:text-secondary-dark"
              >Bot’s Profile Picture</P
            >
          </div>
          <div class="flex w-full">
            <Avatar class="size-12" src={fileImageName ? fileImageResult : ""} />
            <div class="grow ml-3">
              <Button
                ondragover={handleDragOver}
                ondragleave={handleDragLeave}
                ondrop={handleDrop}
                onclick={handleInputFile}
                color="light"
                class={`w-full grow rounded-xl border border-primary-light dark:border-primary-dark bg-primary-light dark:bg-primary-dark py-4 flex flex-col justify-start items-center`}
              >
                <div
                  class="rounded-xl shadow border border-primary-light dark:border-primary-dark justify-center items-center inline-flex mb-3 pointer-events-none"
                >
                  <img
                    src={upload}
                    alt="mail.svg"
                    class="p-4 dark:invert pointer-events-none"
                  />
                </div>
                <div class="flex space-x-1">
                  <P class="text-hc-brand-700 pointer-events-none">{$_("clickToUpload")}</P><P
                    class="pointer-events-none">{$_("orDragAndDrop")}</P
                  >
                </div>
                <P class="pointer-events-none">PNG or SVG file (max. 1MB)</P>
              </Button>
            </div>
          </div>
        </div>
      </div>
    </div>
    <div class="border border-b-1 my-4"></div>
    <div class="w-full flex">
      <P size="lg" weight="semibold" class="w-full md:w-1/4 pb-2">AI Engines</P>
      <div class="w-full md:w-3/4">
        <SelectedCard
          selectionData={platforms}
          bind:selectedData={selectedplatforms}
          isFlex
          isRadio
          widthFlex="w-max"
        />
      </div>
    </div>
    <div class="border border-b-1 my-4"></div>
    <div class="md:flex">
      <P size="lg" weight="semibold" class="w-full md:w-1/4 pb-2">Document Collection</P>
      <div class="w-full md:w-3/4 flex flex-col justify-center gap-2">
        <div class="flex items-center gap-4">
          <div class="w-max">
            <Dropdown
              dropdownOpen={dropdownCollectionOpen}
              bind:selectedItem={selectedCollection}
              items={collections}
              placeholder={"Choose Collection"}
              isOverlay={true}
            />
          </div>
        </div>
      </div>
    </div>
    <div class="border border-b-1 my-4"></div>
    <div class="md:flex">
      <P size="lg" weight="semibold" class="w-full md:w-1/4 pb-2">Additional Prompts</P>
      <div class="w-full md:w-3/4 flex flex-col justify-center gap-2">
        <P size="sm" weight="medium" class="text-secondary-light dark:text-secondary-dark"
          >Additional Prompts*</P
        >
        <Textarea class="resize-none" rows={10} />
      </div>
    </div>
  </div>
</div>
