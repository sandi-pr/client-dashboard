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
  // import {
  //   PUBLIC_API_URL_DEV,
  //   PUBLIC_API_URL_PROD,
  //   PUBLIC_STRAPI_MOTHERSHIP_DEV,
  //   PUBLIC_STRAPI_MOTHERSHIP_PROD,
  // } from "$env/static/public";
  import {
    P,
    Textarea,
    Avatar,
    Input,
    Card,
    Dropdown,
    Button,
    DropdownItem,
  } from "flowbite-svelte";
  // import SvelteMarkdown from "svelte-markdown";
  import ID from "$lib/icon/ID.svg";
  import LoadingSpinner from "$lib/components/LoadingSpinner.svelte";
  import DataTable from "$lib/components/DataTable.svelte";
  import Divider from "$lib/components/Divider.svelte";
  import CopyIcon from "$lib/icon/copy-03.svelte";
  // import card from "$lib/icon/credit-card-refresh.svg";
  import crm from "$lib/icon/users-right.svg";
  import box from "$lib/icon/box.svg";
  import { ChevronDownOutline, DotsVerticalOutline } from "flowbite-svelte-icons";
  import TabDetail from "$lib/components/TabDetail.svelte";
  import Highlight from "svelte-highlight";
  import typescript from "svelte-highlight/languages/typescript";
  import github from "svelte-highlight/styles/github";

  // const API_URL = process.env.NODE_ENV === "development" ? PUBLIC_API_URL_DEV : PUBLIC_API_URL_PROD;
  // const STRAPI_URL =
  //   process.env.NODE_ENV === "development"
  //     ? PUBLIC_STRAPI_MOTHERSHIP_DEV
  //     : PUBLIC_STRAPI_MOTHERSHIP_PROD;

  let isLoading = false;
  let screenWidth;
  // let placement = "bottom-start";

  let users = {
    botName: "Vira",
    deployedPlatform: "Web",
    description: "A personal assistant for Insurance company named Vira",
    additionalPrompts:
      "1. Policy Information and Management: The assistant can provide policyholders with information about their insurance policies, such as coverage details, premium due dates, and policy terms and conditions. It can also assist in managing their policies, such as updating contact information or beneficiary details. \n \n 2. Claims Processing: The assistant can guide policyholders through the claims process, helping them understand what information is needed, how to submit a claim, and tracking the status of their claim.",
    username: "ABC",
    userId: "1454523",
    status: "Active",
    fullname: "ABCABC",
    email: "****@gmail.com",
    role: "Owner Company",
    companyName: "Bank Jasa Jakarta",
  };

  let tabDatas = {
    selected: "Information",
    data: [
      {
        name: "Information",
      },
      {
        name: "Conversation Logs",
      },
      {
        name: "Widget Embed",
      },
      {
        name: "Usage",
      },
    ],
  };

  let datas = [
    [
      {
        id: 1,
        feature: "Profile",
        usage: "100 of 1000 users",
      },
      {
        id: 2,
        feature: "Admins",
        usage: "4 of 5 admins",
      },
    ],
    [
      {
        id: 1,
        feature: "Inventory",
        usage: "100 of 1000 users",
      },
      {
        id: 2,
        feature: "Admins",
        usage: "4 of 5 admins",
      },
    ],
  ];

  let configs = [
    {
      doubleHeader: "CRM",
      iconDoubleHeader: crm,
    },
    {
      doubleHeader: "Inventory",
      iconDoubleHeader: box,
    },
  ];

  let timeoutId;
  const handleSearch = (evt) => {
    clearTimeout(timeoutId);
    if (evt.keyCode !== 13) {
      timeoutId = setTimeout(() => {
        console.log(evt.target.value);
      }, 500);
    } else {
      console.log(evt.target.value);
    }
  };

  let conversationDatas = [
    {
      id: 1,
      timestamp: "03/19/2024 22:57",
      customerName: "Elias Hoppe IV",
      conversationId: "4c64986428240cn948929284",
      actions: {
        objects: [
          {
            name: "View",
            href: "/conversation_logs",
          },
        ],
      },
    },
    {
      id: 2,
      timestamp: "03/20/2024 09:15",
      customerName: "Amelia Jones",
      conversationId: "8d74683756ds33l90928f",
      actions: {
        objects: [
          {
            name: "View",
            href: "/conversation_logs",
          },
        ],
      },
    },
    {
      id: 3,
      timestamp: "03/21/2024 13:45",
      customerName: "Mikhail Petrov",
      conversationId: "5f739dks76402dk23kj4",
      actions: {
        objects: [
          {
            name: "View",
            href: "/conversation_logs",
          },
        ],
      },
    },
  ];

  let conversationConfigs = {
    search: handleSearch,
    pagination: {
      page: 1,
      pageSize: 5,
      pageCount: 1,
      total: 5,
    },
  };

  const classStatusType = (type) => {
    if (type === "Done") {
      return "bg-hc-success-50 border-hc-success-200 text-hc-success-700 dark:bg-hc-success-900 dark:border-hc-success-700 dark:text-hc-success-300";
    } else if (type === "Not Ready") {
      return "bg-hc-error-50 border-hc-error-200 text-hc-error-700 dark:bg-hc-error-900 dark:border-hc-error-700 dark:text-hc-error-300";
    } else if (type === "On Progress") {
      return "bg-hc-warning-50 border-hc-warning-200 text-hc-warning-700 dark:bg-hc-warning-900 dark:border-hc-warning-700 dark:text-hc-warning-300";
    } else if (type === "WhatsApp") {
      return "text-utility-success-700-light dark:text-utility-success-700-dark bg-utility-success-50-light dark:bg-utility-success-50-dark border-utility-success-200-light dark:border-utility-success-200-dark";
    } else if (type === "Web") {
      return "text-sky-700 dark:text-sky-300 bg-sky-50 dark:bg-sky-950 border-sky-200 dark:border-sky-800";
    }
  };

  let formattedDate = "Created on 2021-08-12 14:00 PM";

  let code = `<iframe style=”width: 100%; height: 100%; border: 0px; border-radius: 10px; display: block;" src="http://webchat.hexacode.org/19ca7b92-4bfd-4c20-b6ab-048207e66258"></iframe>`;
</script>

<svelte:head>
  {@html github}
</svelte:head>

<!-- <svelte:window bind:innerWidth={screenWidth} /> -->

<LoadingSpinner {isLoading} />

<div class="m-4">
  <div class="flex justify-start items-center mx-4">
    <Avatar class="size-16 border-[1px] border-hc-graylight-200 ml-4" stacked />
    <div class="ml-4">
      <div class="flex items-center gap-2">
        <P size="3xl" weight="semibold">{users.botName}</P>
        <div
          class={`px-2 py-0.5 rounded-full border justify-start items-center inline-flex text-center text-xs font-medium leading-[18px] 
          ${classStatusType(users.deployedPlatform)}`}
        >
          {users.deployedPlatform}
        </div>
      </div>
      <P class="text-hc-graylight-500 dark:text-hc-graydark-500">{formattedDate}</P>
    </div>
  </div>

  <!-- Submenu -->
  <TabDetail {tabDatas} on:tabChange={(e) => (tabDatas.selected = e.detail)} />

  <!-- Content -->
  {#if tabDatas.selected === "Information"}
    <div class="mx-4 mt-6 space-y-6">
      <div class="md:flex">
        <P
          size="base"
          weight="semibold"
          class="w-full md:w-1/4 text-secondary-light dark:text-secondary-dark">Bots Information</P
        >

        <div class="w-full md:w-3/4 grid grid-cols-1 md:grid-cols-4 gap-4">
          <div
            class="flex flex-col gap-1 col-span-4 [&>*]:bg-primary-light [&>*]:dark:bg-primary-dark [&>textarea]:text-primary-light dark:[&>textarea]:text-primary-dark"
          >
            <P size="sm" weight="medium" class="text-secondary-light dark:text-secondary-dark"
              >Bots Name*</P
            >
            <Textarea
              value={users.botName}
              class="resize-none"
              rows={screenWidth > 1024 ? "1" : "2"}
              unWrappedClass="p-2.5 text-sm border border-primary-light dark:border-primary-dark focus:ring-primary-500 focus:border-primary-500 dark:focus:ring-primary-500 dark:focus:border-primary-500"
            />
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
              >Isurance<ChevronDownOutline
                class="w-3 h-3 ms-2 text-fg-quaternary-light dark:text-fg-quaternary-dark"
              /></Button
            >
            <!-- <Dropdown containerClass="w-10/12 sm:w-1/4">
              <DropdownItem>Isurance</DropdownItem>
              <DropdownItem>Information Technology (IT)</DropdownItem>
              <DropdownItem>Healthcare</DropdownItem>
              <DropdownItem>Food and Beverage</DropdownItem>
              <DropdownItem>Automotive</DropdownItem>
              <DropdownItem>Energy</DropdownItem>
              <DropdownItem>Finance</DropdownItem>
              <DropdownItem>Tourism and Hospitality</DropdownItem>
              <DropdownItem>Manufacturing</DropdownItem>
              <DropdownItem>Retail</DropdownItem>
              <DropdownItem>Construction and Real Estate</DropdownItem>
              <DropdownItem>Entertainment and Media</DropdownItem>
              <DropdownItem>Education</DropdownItem>
            </Dropdown> -->
            <P size="sm" class="!mt-1" color="text-tertiary-light dark:text-tertiary-dark"
              >This will set bots base configuration</P
            >
          </div>
          <div
            class="space-y-1 col-span-4 sm:col-span-2 [&>*]:bg-primary-light [&>*]:dark:bg-primary-dark [&>textarea]:text-primary-light dark:[&>textarea]:text-primary-dark"
          >
            <P size="sm" weight="medium" class="text-secondary-light dark:text-secondary-dark"
              >Use Cases*</P
            >
            <Button
              class="w-full flex justify-between border-primary-light dark:border-primary-dark"
              color="light"
              >Personal Agent Assistant<ChevronDownOutline
                class="w-3 h-3 ms-2 text-fg-quaternary-light dark:text-fg-quaternary-dark"
              /></Button
            >
            <!-- <Dropdown containerClass="w-10/12 sm:w-1/4">
              <DropdownItem>Personal Agent Assistant</DropdownItem>
              <DropdownItem>Customer Support Chatbot</DropdownItem>
              <DropdownItem>Automated Data Analysis Tool</DropdownItem>
              <DropdownItem>Virtual Personal Trainer</DropdownItem>
              <DropdownItem>Remote Health Monitoring System</DropdownItem>
              <DropdownItem>Supply Chain Optimization Software</DropdownItem>
              <DropdownItem>Smart Home Automation System</DropdownItem>
              <DropdownItem>Online Learning Platform</DropdownItem>
              <DropdownItem>Travel Planning Application</DropdownItem>
              <DropdownItem>Financial Portfolio Management Tool</DropdownItem>
              <DropdownItem>Social Media Analytics Platform</DropdownItem>
              <DropdownItem>Remote Collaboration Tool</DropdownItem>
            </Dropdown> -->
            <P size="sm" class="!mt-1" color="text-tertiary-light dark:text-tertiary-dark"
              >This will impact how your bots tailored to specific use cases</P
            >
          </div>
          <div
            class="space-y-1 col-span-4 [&>*]:bg-primary-light [&>*]:dark:bg-primary-dark [&>textarea]:text-primary-light dark:[&>textarea]:text-primary-dark"
          >
            <P size="sm" weight="medium" class="text-secondary-light dark:text-secondary-dark"
              >Description (Opsional)</P
            >
            <Textarea value={users.description} class="resize-none" />
          </div>
        </div>
      </div>
      <Divider />
      <div class="md:flex">
        <P size="lg" weight="semibold" class="w-full md:w-1/4 pb-2">Additional Prompts</P>
        <div
          class="w-full md:w-3/4 flex flex-col justify-center gap-2 [&>*]:bg-primary-light [&>*]:dark:bg-primary-dark"
        >
          <P>Additional Prompts*</P>
          <Textarea value={users.additionalPrompts} class="" rows={10} />
        </div>
      </div>
    </div>
  {:else if tabDatas.selected === "Conversation Logs"}
    <!-- Data Tables -->
    <div>
      <DataTable datas={conversationDatas} configs={conversationConfigs} />
    </div>
  {:else if tabDatas.selected === "Widget Embed"}
    <div
      class="grid grid-cols-1 md:grid-cols-4 gap-4 mx-4 mt-8 [&>*]:text-hc-graylight-700 [&>*]:bg-white [&>*]:dark:bg-hc-graydark-950 [&>textarea]:text-hc-graylight-500"
    >
      <P size="lg" weight="semibold" class="col-span-1">Embed Codes</P>
      <div class="col-span-3 space-y-2">
        <Button
          outline
          color="purple"
          class=" w-max hover:!bg-transparent !text-sm !font-semibold !text-utility-brand-600-light dark:!text-utility-brand-600-dark !border-brand-light dark:!border-brand-dark focus-within:!ring-1 px-3 py-[10px] gap-1"
        >
          <CopyIcon
            svgClass="w-5 h-5 text-fg-brand-primary-light dark:text-fg-brand-primary-dark"
          />
          Copy to clipboard
        </Button>
        <Highlight language={typescript}   {code} />
        <!-- <SvelteMarkdown></SvelteMarkdown> -->
      </div>
    </div>
  {:else if tabDatas.selected === "Usage"}
    <div
      class="grid grid-cols-1 md:grid-cols-4 gap-4 mx-4 mt-8 [&>*]:text-hc-graylight-700 [&>*]:bg-white [&>*]:dark:bg-hc-graydark-950 [&>textarea]:text-hc-graylight-500"
    >
      <P size="lg" weight="semibold" class="col-span-1">Usages</P>
      <div class="col-span-3">
        <div>
          <DataTable datas={datas[0]} configs={configs[0]} />
        </div>
        <div>
          <DataTable datas={datas[1]} configs={configs[1]} />
        </div>
      </div>
    </div>
  {/if}
</div>
