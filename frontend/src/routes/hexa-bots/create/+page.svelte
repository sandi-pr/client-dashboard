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
    // Dropdown,
    // DropdownItem,
  } from "flowbite-svelte";
  import { ChevronDownOutline } from "flowbite-svelte-icons";
  import WebIcon from "$lib/icon/server-06.svg";
  // import WhatsappIcon from "$lib/icon/whatsapp.svg";
  // import TelegramIcon from "$lib/icon/telegram.svg";
  import LoadingSpinner from "$lib/components/LoadingSpinner.svelte";
  import { onMount } from "svelte";
  // import axios from "axios";
  // import { getTemplates } from "$lib/service/strapi";

  // export let data;

  let placement = "bottom-start";

  let platforms = [
    {
      name: "Web",
      icon: WebIcon,
      description: "A web-based chat interface for your bot, built with SvelteKit.",
      selected: true,
      disabled: true,
    },
    // {
    //   name: "Whatsapp",
    //   icon: WhatsappIcon,
    //   description: "A mobile messaging app for your bot, widely used for communication.",
    //   selected: false,
    //   disabled: true,
    // },
    // {
    //   name: "Telegram",
    //   icon: TelegramIcon,
    //   description: "A cloud-based instant messaging platform for your bot.",
    //   selected: false,
    //   disabled: true,
    // },
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

  onMount(() => {
    // getTemplate();
  });

  let templates = [];
  let dropdownIndustryOpen = false;
  let dropdownUseCaseOpen = false;
  let params = [];
  let products_and_services = [];

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
</script>

<svelte:window bind:innerWidth={screenWidth} />

<LoadingSpinner {isLoading} />

<div class="m-4 lg:ml-80">
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
  <div class="mx-4 mt-8">
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
        <div
          class="space-y-1 col-span-4 sm:col-span-2 [&>*]:bg-primary-light [&>*]:dark:bg-primary-dark [&>textarea]:text-primary-light dark:[&>textarea]:text-primary-dark"
        >
          <P size="sm" weight="medium" class="text-secondary-light dark:text-secondary-dark"
            >Industry*</P
          >
          <Button
            class="w-full flex justify-between border-primary-light dark:border-primary-dark"
            color="light"
            >{templates.length > 0 ? templates[0].attributes.name : ""}<ChevronDownOutline
              class="w-3 h-3 ms-2 text-fg-quaternary-light dark:text-fg-quaternary-dark"
            /></Button
          >
          <!-- <Dropdown containerClass="w-10/12 sm:w-1/4" {placement} bind:open={dropdownIndustryOpen}>
            {#each templates as template}
              <DropdownItem
                on:click={() => (
                  (dropdownIndustryOpen = false), (params = template.attributes.parameters)
                )}>{template.attributes.name}</DropdownItem
              >
            {/each}
          </Dropdown> -->
          <P size="sm" class="!mt-1" color="text-tertiary-light dark:text-tertiary-dark"
            >This will set bots base configuration</P
          >
        </div>
        <div
          class="space-y-1 col-span-4 sm:col-span-2 [&>*]:bg-primary-light [&>*]:dark:bg-primary-dark [&>textarea]:text-primary-light dark:[&>textarea]:text-primary-dark"
        >
          <P size="sm" weight="medium" class="text-secondary-light dark:text-secondary-dark"
            >{params[0]?.children[0]?.label}{params[0]?.children[0]?.required ? "*" : ""}</P
          >
          <Button
            class="w-full flex justify-between border-primary-light dark:border-primary-dark"
            color="light"
            >{params[0]?.children?.length > 0
              ? params[0]?.children[0]?.enum[0]
              : ""}<ChevronDownOutline
              class="w-3 h-3 ms-2 text-fg-quaternary-light dark:text-fg-quaternary-dark"
            /></Button
          >
          <!-- <Dropdown containerClass="w-10/12 sm:w-1/4" {placement} bind:open={dropdownUseCaseOpen}>
            {#each params[0]?.children[0]?.enum as useCase}
              <DropdownItem on:click={() => (dropdownUseCaseOpen = false)}>{useCase}</DropdownItem>
            {/each}
          </Dropdown> -->
          <P size="sm" class="!mt-1" color="text-tertiary-light dark:text-tertiary-dark"
            >{params[0]?.children[0]?.description}</P
          >
        </div>
        {#if params[0]?.children.length > 1}
          {#each params[0].children as param, index (param)}
            {#if index !== 0}
              <div
                class="flex flex-col gap-1 col-span-4 [&>*]:bg-primary-light [&>*]:dark:bg-primary-dark [&>textarea]:text-primary-light dark:[&>textarea]:text-primary-dark"
              >
                <P size="sm" weight="medium" class="text-secondary-light dark:text-secondary-dark"
                  >{param.label}{param.required ? "*" : ""}</P
                >
                <Input
                  type={param.type}
                  placeholder={param.placeholder}
                  defaultClass={param.required}
                  min={param.min}
                  max={param.max}
                />
              </div>
            {/if}
          {/each}
        {/if}
      </div>
    </div>
    <div class="border border-b-1 my-4"></div>
    <div class="md:flex">
      <P size="lg" weight="semibold" class="w-full md:w-1/4 pb-2">Product & Services</P>
      <div class="w-full md:w-3/4 flex flex-col justify-center gap-2">
        {#if params[1]?.children.length > 0}
          <P>{params[1]?.label !== undefined ? params[1]?.label : "Products & Services"}</P>
          <div class="grid grid-cols-1 lg:grid-cols-3 gap-4 w-full items-center justify-center">
            {#each products_and_services as product, index}
              <Card>
                {#each params[1]?.children[0].list_dict_format as field}
                  <div class="flex flex-col gap-1 col-span-4">
                    <P
                      size="sm"
                      weight="medium"
                      class="text-secondary-light dark:text-secondary-dark"
                    >
                      {field.label}{field.required ? "*" : ""}
                    </P>
                    <Input
                      type={field.type}
                      placeholder={field.placeholder}
                      defaultClass={field.required}
                      min={field.min}
                      max={field.max}
                      bind:value={product[field.name]}
                    />
                  </div>
                {/each}
              </Card>
            {/each}
          </div>
          <Button on:click={addProduct} color="light">Tambah</Button>
        {/if}
      </div>
    </div>
    <div class="border border-b-1 my-4"></div>
    <div class="md:flex">
      <P size="lg" weight="semibold" class="w-full md:w-1/4 pb-2">Deployed Platforms</P>
      <div class="w-full md:w-3/4 grid grid-cols-3 gap-4">
        {#each platforms as platform}
          <Checkbox checked={platform.selected} disabled={platform.disabled} custom>
            <div
              class={`font-normal p-5 w-full text-gray-500 bg-white rounded-lg border-2 border-gray-200 cursor-pointer dark:hover:text-gray-300 dark:border-gray-700 peer-checked:border-primary-600 hover:text-gray-600 dark:peer-checked:text-gray-300 peer-checked:text-gray-600 dark:text-gray-400 dark:bg-gray-800 hover:bg-primary-hover-light dark:hover:bg-primary-hover-dark
              ${platform.disabled ? "cursor-default bg-disabled-subtle-light dark:bg-disabled-subtle-dark bg-opacity-50 hover:bg-disabled-subtle-light dark:hover:bg-disabled-subtle-dark" : ""}`}
            >
              <div class="flex items-center gap-3">
                <div
                  class="w-10 h-10 p-2.5 bg-tertiary-light dark:bg-tertiary-dark rounded-full justify-center items-center flex"
                >
                  <!-- <WebIcon svgClass="w-5 h-5 text-violet-500" /> -->
                  <img src={platform.icon} alt="Web" class="w-5 h-5" />
                </div>
                <P
                  weight="medium"
                  class={`${platform.disabled ? "text-disabled-light dark:text-disabled-dark" : "text-secondary-light dark:text-secondary-dark"}`}
                  >{platform.name}</P
                >
              </div>
              <!-- <div class="w-full pt-2">
                <P size="sm" class="text-tertiary-light dark:text-tertiary-dark"
                  >{platform.description}</P
                >
              </div> -->
            </div>
          </Checkbox>
        {/each}
      </div>
    </div>
    <div class="border border-b-1 my-4"></div>
    <div class="md:flex">
      <P size="lg" weight="semibold" class="w-full md:w-1/4 pb-2">Avability</P>
      <div class="w-full md:w-3/4 flex flex-col justify-center gap-2">
        <div class="flex items-center gap-4">
          <P>Select day</P>
          {#each selectedDays as day}
            <Checkbox custom checked={day.selected}>
              <div
                class="flex items-center font-normal w-9 h-9 text-disabled-light dark:text-disabled-dark bg-quarterary-light dark:bg-quarterary-dark rounded-full cursor-pointer peer-checked:text-white dark:peer-checked:text-white peer-checked:bg-fg-brand-primary-light dark:peer-checked:bg-fg-brand-primary-dark"
              >
                <div class="w-full text-base font-semibold text-center">{day.label}</div>
              </div>
            </Checkbox>
          {/each}
        </div>
      </div>
    </div>
    <div class="border border-b-1 my-4"></div>
    <div class="md:flex">
      <P size="lg" weight="semibold" class="w-full md:w-1/4 pb-2">Additional Prompts</P>
      <div class="w-full md:w-3/4 flex flex-col justify-center gap-2">
        <P>Additional Prompts*</P>
        <Textarea class="resize-none" rows={10} />
      </div>
    </div>
  </div>
</div>
