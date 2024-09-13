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
  import { P, Avatar, Button, Toggle, Indicator, Modal } from "flowbite-svelte";
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
  import { Typography, Badge, InputRw, TextareaRw, SelectedCard, UploadFile, Dropdown, triggerToast } from "$lib";
  // import Highlight from "svelte-highlight";
  // import typescript from "svelte-highlight/languages/typescript";
  // import github from "svelte-highlight/styles/github";

  // const API_URL = process.env.NODE_ENV === "development" ? PUBLIC_API_URL_DEV : PUBLIC_API_URL_PROD;
  // const STRAPI_URL =
  //   process.env.NODE_ENV === "development"
  //     ? PUBLIC_STRAPI_MOTHERSHIP_DEV
  //     : PUBLIC_STRAPI_MOTHERSHIP_PROD;

  let isLoading = false;
  let screenWidth;
  let botFile = [];
  let userFile = [];
  let logoFile = [];
  let backgroundFile = [];
  let modalThemeColor = false;
  let modalOverlayColor = false;
  let dropdownBgAlignOpen = false;
  let selectedBgAlign = "center";
  let bgAlign = [
    { name: "Center", value: "center" },
    { name: "Left", value: "left" },
    { name: "Right", value: "right" },
  ]
  // let placement = "bottom-start";

  let users = {
    botName: "Chat Qur’ani",
    logoName: "Hexabot",
    logoOnly: false,
    deployedPlatform: "Web",
    description: "A chatbot who has vast knowledge of Al Quran",
    additionalPrompts:
      "1. Policy Information and Management: The assistant can provide policyholders with information about their insurance policies, such as coverage details, premium due dates, and policy terms and conditions. It can also assist in managing their policies, such as updating contact information or beneficiary details. \n \n 2. Claims Processing: The assistant can guide policyholders through the claims process, helping them understand what information is needed, how to submit a claim, and tracking the status of their claim.",
    welcomeMessage: "Hi! This is an AI assistant that is trained on Al-Quran. For example, you can ask for a Tafsir of a verse! Now, how can I help you?",    
    username: "You",
    primaryBrandColor: "purple",
    selectedBrandColor: 1,
    overlayColor: "purple",
    selectedOverlayColor: 1,
    selectedPlatform: 1,
    selectedLayout: 1,
    userId: "1454523",
    status: "Active",
    embedCode: `<iframe style=”width: 100%; height: 100%; border: 0px; border-radius: 10px; display: block;" src="http://webchat.hexacode.org/19ca7b92-4bfd-4c20-b6ab-048207e66258"></iframe>`,
  };

  let tabDatas = {
    selected: "Customization",
    data: [
      {
        name: "Information",
      },
      {
        name: "Customization",
      },
      // {
      //   name: "Data",
      // },
      {
        name: "Widget Embed",
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

  let platforms = [
    {
      id: 1,
      name: "LLama 3.1",
    },
  ];

  let layout = [
    {
      id: 1,
      name: "Modern Minimalism",
    },
    {
      id: 2,
      name: "Modern with Background",
    }
  ]

  let themeColors = [
    {
      id: 1,
      name: "Purple",
      value: "purple",
    },
    {
      id: 2,
      name: "Blue",
      value: "blue",
    },
    {
      id: 3,
      name: "Green",
      value: "green",
    },
  ];

  function getBrandColorName(selectedId, value) {
    const selectedColor = themeColors.find((color) => color.id === selectedId);
    return selectedColor ? value === "name" ? selectedColor.name : selectedColor.value : "Unknown Color";
  }

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
      return "success";
    } else if (type === "Not Ready") {
      return "error";
    } else if (type === "On Progress") {
      return "warning";
    } else if (type === "WhatsApp") {
      return "success";
    } else if (type === "Web") {
      return "sky";
    }
  };

  const handleInputFile = (as) => {
    const inputElement = document.createElement("input");
    inputElement.type = "file";
    inputElement.accept = "image/jpeg, image/png";

    inputElement.addEventListener("change", async (event) => {
      const file = event.target.files[0];

      if (file && (file.type === "image/jpeg" || file.type === "image/png")) {
        if (as === "bot") {
          botFile.botImageResult = await new Promise((resolve, reject) => {
            const reader = new FileReader();

            reader.onload = () => resolve(reader.result);
            reader.onerror = reject;

            reader.readAsDataURL(file);
          });

          botFile.botImage = await readFileAsBase64(file);
          botFile.botImageName = file.name;
        } else if (as === "user") {
          userFile.userImageResult = await new Promise((resolve, reject) => {
            const reader = new FileReader();

            reader.onload = () => resolve(reader.result);
            reader.onerror = reject;

            reader.readAsDataURL(file);
          });

          userFile.userImage = await readFileAsBase64(file);
          userFile.userImageName = file.name;
        } else if (as === "logo") {
          logoFile.logoImageResult = await new Promise((resolve, reject) => {
            const reader = new FileReader();

            reader.onload = () => resolve(reader.result);
            reader.onerror = reject;

            reader.readAsDataURL(file);
          });

          logoFile.logoImage = await readFileAsBase64(file);
          logoFile.logoImageName = file.name;
        } else if (as === "background") {
          backgroundFile.backgroundImageResult = await new Promise((resolve, reject) => {
            const reader = new FileReader();

            reader.onload = () => resolve(reader.result);
            reader.onerror = reject;

            reader.readAsDataURL(file);
          });

          backgroundFile.backgroundImage = await readFileAsBase64(file);
          backgroundFile.backgroundImageName = file.name;
        }
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

  function copyToClipboard(text) {
    navigator.clipboard.writeText(text)
      .then(() => {
        console.log('Text copied to clipboard');
        // Opsional: Tambahkan notifikasi atau toast untuk menunjukkan teks telah disalin
        triggerToast('copyToClipboard');
      })
      .catch(err => {
        console.error('Failed to copy: ', err);
      });
  }

  let formattedDate = "Created on 2021-08-12 14:00 PM";

  let code = `<iframe style=”width: 100%; height: 100%; border: 0px; border-radius: 10px; display: block;" src="http://webchat.hexacode.org/19ca7b92-4bfd-4c20-b6ab-048207e66258"></iframe>`;
</script>

<!-- <svelte:head>
  {@html github}
</svelte:head> -->

<!-- <svelte:window bind:innerWidth={screenWidth} /> -->

<!-- <LoadingSpinner {isLoading} /> -->

<div class="m-4 pb-12">
  <div class="sticky top-0 bg-primary-light dark:bg-primary-dark z-10">
    <div class="flex justify-start items-center mx-4">
      <Avatar class="size-14 border-[1px] border-hc-graylight-200" stacked />
      <div class="ml-4">
        <div class="flex items-center gap-2">
          <Typography preset="Display sm/Semibold" text={users.botName} color="primary" />
          <Badge type={classStatusType(users.deployedPlatform)} text={users.deployedPlatform} />
        </div>
        <Typography preset="Text md/Medium" text={formattedDate} color="quarterary" />
      </div>
    </div>

    <!-- Submenu -->
    <TabDetail {tabDatas} on:tabChange={(e) => (tabDatas.selected = e.detail)} />
  </div>

  <!-- Content -->
  {#if tabDatas.selected === "Information"}
    <div class="mt-6 space-y-6">
      <div class="md:flex">
        <Typography preset="Text md/Semibold" text="Bots Information" color="secondary" propClass="w-full md:w-1/4" />

        <div class="w-full md:w-3/4 grid grid-cols-1 md:grid-cols-4 gap-3">
          <div
            class="flex flex-col gap-1.5 col-span-4"
          >
            <Typography preset="Text sm/Medium" text="Bots Name" color="secondary" />
            <InputRw bind:value={users.botName} onEdit />
          </div>
          <div
            class="flex flex-col gap-1.5 col-span-4"
          >
            <Typography preset="Text sm/Medium" text="Description (Optional)" color="secondary" />
            <TextareaRw bind:value={users.description} onEdit />
          </div>
        </div>
      </div>
      <Divider />
      <div class="md:flex">
        <Typography preset="Text md/Semibold" text="AI Engines" color="secondary" propClass="w-full md:w-1/4" />
        <div class="w-full md:w-3/4">
          <SelectedCard
            selectionData={platforms}
            bind:selectedData={users.selectedPlatform}
            isFlex
            isRadio
            widthFlex="w-max"
          />
        </div>
      </div>
      <Divider />
      <div class="md:flex">
        <Typography preset="Text md/Semibold" text="Additional Prompts" color="secondary" propClass="w-full md:w-1/4" />
        <div
          class="w-full md:w-3/4 flex flex-col justify-center gap-2 [&>*]:bg-primary-light [&>*]:dark:bg-primary-dark"
        >
          <Typography preset="Text sm/Medium" text="Bot’s Role" color="secondary" />
          <TextareaRw bind:value={users.additionalPrompts} onEdit />
        </div>
      </div>
      <Divider />
      <div class="flex justify-end gap-3">
        <Button
          onclick={() => {}}
          color="light"
          class="px-3 focus-within:!ring-1">{"Cancel"}</Button
        >
        <Button
          onclick={() => {}}
          class="px-3 focus-within:!ring-1">{"Save"}</Button
        >
      </div>
    </div>
  {:else if tabDatas.selected === "Customization"}
    <div class="mt-6 space-y-6">
      <div class="md:flex">
        <Typography preset="Text md/Semibold" text="General" color="secondary" propClass="w-full md:w-1/4" />

        <div class="w-full md:w-3/4 grid grid-cols-1 md:grid-cols-4 gap-4">
          <div class="flex flex-col gap-1.5 col-span-4">
            <Typography preset="Text sm/Medium" text="Welcome Message" color="secondary" />
            <TextareaRw bind:value={users.welcomeMessage} onEdit />
          </div>
          <div class="flex gap-4 col-span-4">
            <div class="flex flex-col gap-1.5 grow">
              <Typography preset="Text sm/Medium" text="Bots Name" color="secondary" />
              <InputRw bind:value={users.botName} onEdit />
            </div>
            <div class="flex flex-col gap-1.5 grow">
              <Typography preset="Text sm/Medium" text="User Name" color="secondary" />
              <InputRw bind:value={users.username} onEdit />
            </div>
          </div>
          <div class="col-span-2">
            <UploadFile
              title="Bots Profile Picture"
              handleFileUploadWrapper={() => handleInputFile("bot")}
              isPreview
              srcAvatar={botFile.botImageName ? botFile.botImageResult : ""}
            />
          </div>
          <div class="col-span-2">
            <UploadFile
              title="User Profile Picture"
              handleFileUploadWrapper={() => handleInputFile("user")}
              isPreview
              srcAvatar={userFile.userImageName ? userFile.userImageResult : ""}
            />
          </div>
        </div>
      </div>
      <Divider />
      <div class="md:flex">
        <Typography preset="Text md/Semibold" text="Logo" color="secondary" propClass="w-full md:w-1/4" />
        <div class="w-full md:w-3/4 grid grid-cols-1 md:grid-cols-4 gap-4">
          <div class="col-span-4">
            <Toggle checked={users.logoOnly} class="w-max">Logo only</Toggle>
          </div>
          <div class="col-span-2">
            <UploadFile
              title="Logo"
              handleFileUploadWrapper={() => handleInputFile("logo")}
              isPreview
              srcAvatar={logoFile.logoImageName ? logoFile.logoImageResult : ""}
            />
          </div>
          <div class="flex flex-col gap-1.5 col-span-2">
            <Typography preset="Text sm/Medium" text="Logo Name" color="secondary" />
            <InputRw bind:value={users.logoName} onEdit />
          </div>
        </div>
      </div>
      <Divider />
      <div class="md:flex">
        <Typography preset="Text md/Semibold" text="Theme Colors" color="secondary" propClass="w-full md:w-1/4" />
        <div class="w-full md:w-3/4 grid grid-cols-1 md:grid-cols-4 gap-4">
          <div class="flex flex-col gap-1.5 col-span-2">
            <Typography preset="Text sm/Medium" text="Primary Brand Color" color="secondary" />
            <Button
              on:click={() => (modalThemeColor = true)}
              class="flex justify-between border-primary-light dark:border-primary-dark dark:bg-gray-800"
              color="light"
            >
              <div class="flex items-center gap-2">
                <Indicator color={users.primaryBrandColor} />
                {getBrandColorName(users.selectedBrandColor, "name")}
              </div>
              <ChevronDownOutline
                class="w-3 h-3 ms-2 text-fg-quaternary-light dark:text-fg-quaternary-dark"
              />
            </Button>
          </div>
        </div>
      </div>
      <Divider />
      <div class="md:flex">
        <Typography preset="Text md/Semibold" text="Layout" color="secondary" propClass="w-full md:w-1/4" />
        <div class="w-full md:w-3/4">
          <SelectedCard
            selectionData={layout}
            bind:selectedData={users.selectedLayout}
            isFlex
            isRadio
            widthFlex="w-full"
            spacing="!p-4"
          />
        </div>
      </div>
      <Divider />
      <div class="md:flex">
        <Typography preset="Text md/Semibold" text="Background" color="secondary" propClass="w-full md:w-1/4" />
        <div class="w-full md:w-3/4 grid grid-cols-1 md:grid-cols-4 gap-4">
          <div class="col-span-4">
            <UploadFile
              title="Background Picture"
              handleFileUploadWrapper={() => handleInputFile("background")}
              isPreview
              largePreview
              srcAvatar={backgroundFile.backgroundImageName ? backgroundFile.backgroundImageResult : ""}
            />
          </div>
          <div class="flex flex-col gap-1.5 col-span-2">
            <Typography preset="Text sm/Medium" text="Overlay Color" color="secondary" />
            <Button
              on:click={() => (modalOverlayColor = true)}
              class="flex justify-between border-primary-light dark:border-primary-dark dark:bg-gray-800"
              color="light"
            >
              <div class="flex items-center gap-2">
                <Indicator color={users.overlayColor} />
                {getBrandColorName(users.selectedOverlayColor, "name")}
              </div>
              <ChevronDownOutline
                class="w-3 h-3 ms-2 text-fg-quaternary-light dark:text-fg-quaternary-dark"
              />
            </Button>
          </div>
          <div class="flex flex-col gap-1.5 col-span-2">
            <Typography preset="Text sm/Medium" text="Background Align" color="secondary" />
            <Dropdown
              dropdownOpen={dropdownBgAlignOpen}
              bind:selectedItem={selectedBgAlign}
              items={bgAlign}
              placeholder={"Choose Collection"}
              isOverlay={true}
            />
          </div>
        </div>
      </div>
      <Divider />
      <div class="flex justify-end gap-3">
        <Button
          onclick={() => {}}
          color="light"
          class="px-3 focus-within:!ring-1">{"Cancel"}</Button
        >
        <Button
          onclick={() => {}}
          class="px-3 focus-within:!ring-1">{"Save"}</Button
        >
      </div>
    </div>
  <!-- {:else if tabDatas.selected === "Conversation Logs"}
    <div>
      <DataTable datas={conversationDatas} configs={conversationConfigs} />
    </div> -->
  {:else if tabDatas.selected === "Widget Embed"}
    <div class="mt-6 space-y-6">
      <div class="md:flex">
        <Typography preset="Text md/Semibold" text="Embed Codes" color="secondary" propClass="w-full md:w-1/4" />

        <div class="w-full md:w-3/4 grid grid-cols-1 md:grid-cols-4 gap-4">
          <Button
            outline
            on:click={() => copyToClipboard(users.embedCode)}
            color="purple"
            class=" w-max col-span-4 hover:!bg-transparent !text-sm !font-semibold !text-utility-brand-600-light dark:!text-utility-brand-600-dark !border-brand-light dark:!border-brand-dark focus-within:!ring-1 px-3 py-[10px] gap-1"
          >
            <CopyIcon
              svgClass="w-5 h-5 text-fg-brand-primary-light dark:text-fg-brand-primary-dark"
            />
            Copy to clipboard
          </Button>
          <div class="col-span-4">
            <TextareaRw bind:value={users.embedCode} onEdit={false} />
          </div>
        </div>
      </div>
    </div>
  {/if}
</div>

<Modal title="Terms of Service" bind:open={modalThemeColor} autoclose outsideclose>
  <SelectedCard
    selectionData={themeColors}
    bind:selectedData={users.selectedBrandColor}
    numGrid="2"
    isRadio
    widthFlex="w-full"
    spacing="!p-4"
    widgetIndicator={Indicator}
  />
  <svelte:fragment slot="footer">
    <Button 
      onclick={() => (users.primaryBrandColor = getBrandColorName(users.selectedBrandColor))}
      class="px-3 w-full focus-within:!ring-1">{"Save"}</Button
    >
  </svelte:fragment>
</Modal>

<Modal title="Terms of Service" bind:open={modalOverlayColor} autoclose outsideclose>
  <SelectedCard
    selectionData={themeColors}
    bind:selectedData={users.selectedOverlayColor}
    numGrid="2"
    isRadio
    widthFlex="w-full"
    spacing="!p-4"
    widgetIndicator={Indicator}
  />
  <svelte:fragment slot="footer">
    <Button 
      onclick={() => (users.overlayColor = getBrandColorName(users.selectedOverlayColor))}
      class="px-3 w-full focus-within:!ring-1">{"Save"}</Button
    >
  </svelte:fragment>
</Modal>
