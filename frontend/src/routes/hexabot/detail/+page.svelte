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
  import {
    Typography,
    Badge,
    InputRw,
    TextareaRw,
    SelectedCard,
    UploadFile,
    Dropdown,
    triggerToast,
    triggerLoading,
    call,
  } from "$lib";
  import { onMount } from "svelte";
  import { page } from "$app/stores";
  import { goto } from "$app/navigation";
  // import Highlight from "svelte-highlight";
  // import typescript from "svelte-highlight/languages/typescript";
  // import github from "svelte-highlight/styles/github";

  // const API_URL = process.env.NODE_ENV === "development" ? PUBLIC_API_URL_DEV : PUBLIC_API_URL_PROD;
  // const STRAPI_URL =
  //   process.env.NODE_ENV === "development"
  //     ? PUBLIC_STRAPI_MOTHERSHIP_DEV
  //     : PUBLIC_STRAPI_MOTHERSHIP_PROD;

  const isFrontend = process.env.NODE_ENV === "development";

  let ready = false;
  let botId;
  let botFile = { name: "" };
  let userFile = { name: "" };
  let logoFile = { name: "" };
  let backgroundFile = {
    name: "",
  };
  let modalThemeColor = false;
  let modalOverlayColor = false;
  let dropdownBgAlignOpen = false;
  let background_align = "center";
  let bgAlign = [
    { name: "Center", value: "center" },
    { name: "Left", value: "left" },
    { name: "Right", value: "right" },
  ];

  let users = {};

  let tmpThemeColor = "";
  let tmpOverlayColor = "";

  let tabDatas = {
    selected: "Information",
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

  let platforms = [];

  let layouts = [
    {
      id: 1,
      name: "Modern Minimalism",
    },
    {
      id: 2,
      name: "Modern with Background",
    },
  ];

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
    return selectedColor
      ? value === "name"
        ? selectedColor.name
        : selectedColor.value
      : "Unknown Color";
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
          botFile.imageResult = await new Promise((resolve, reject) => {
            const reader = new FileReader();

            reader.onload = () => resolve(reader.result);
            reader.onerror = reject;

            reader.readAsDataURL(file);
          });

          botFile.base64Data = await readFileAsBase64(file);
          botFile.name = file.name;
        } else if (as === "user") {
          userFile.imageResult = await new Promise((resolve, reject) => {
            const reader = new FileReader();

            reader.onload = () => resolve(reader.result);
            reader.onerror = reject;

            reader.readAsDataURL(file);
          });

          userFile.base64Data = await readFileAsBase64(file);
          userFile.name = file.name;
        } else if (as === "logo") {
          logoFile.imageResult = await new Promise((resolve, reject) => {
            const reader = new FileReader();

            reader.onload = () => resolve(reader.result);
            reader.onerror = reject;

            reader.readAsDataURL(file);
          });

          logoFile.base64Data = await readFileAsBase64(file);
          logoFile.name = file.name;
        } else if (as === "background") {
          backgroundFile.imageResult = await new Promise((resolve, reject) => {
            const reader = new FileReader();

            reader.onload = () => resolve(reader.result);
            reader.onerror = reject;

            reader.readAsDataURL(file);
          });

          backgroundFile.base64Data = await readFileAsBase64(file);
          backgroundFile.name = file.name;
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
    navigator.clipboard
      .writeText(text)
      .then(() => {
        console.log("Text copied to clipboard");
        // Opsional: Tambahkan notifikasi atau toast untuk menunjukkan teks telah disalin
        triggerToast("copyToClipboard");
      })
      .catch((err) => {
        console.error("Failed to copy: ", err);
      });
  }

  let formattedDate = "Created on 2021-08-12 14:00 PM";

  let code = `<iframe style=”width: 100%; height: 100%; border: 0px; border-radius: 10px; display: block;" src="http://webchat.hexacode.org/19ca7b92-4bfd-4c20-b6ab-048207e66258"></iframe>`;

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

  async function getBot() {
    let token = sessionStorage.getItem("token");
    if (token) {
      try {
        let response = await $call.call("client-dashboard", "get-bot", {
          URLQuery: {},
          URLParams: {
            id: botId,
          },
          Params: {},
          Headers: {
            Authorization: "Bearer " + token,
          },
        });

        if (response.status !== 200) {
          console.log(response);
          triggerLoading(false);
          return;
        }

        const bot = response.result;
        const attributes = bot.attributes;

        let urlBot = (await $call.getPage("hexabot", "home-page")) + "?id=" + attributes.bot_id;

        users = {
          botName: attributes.bot_name,
          logoName: attributes.logo_name,
          logoOnly: attributes.logo_only,
          deployedPlatform: attributes?.deployed_platform || "Web",
          description: attributes.description,
          additionalPrompts: attributes.additional_prompt,
          welcomeMessage: attributes.greeting_text,
          username: attributes?.username || "User",
          primaryBrandColor:
            attributes.color_theme !== null
              ? themeColors.find((color) => color.value === attributes.color_theme).id
              : 1,
          overlayColor:
            attributes.overlay_color !== null
              ? themeColors.find((color) => color.value === attributes.overlay_color).id
              : 1,
          platform: platforms.find(
            (platform) => platform.model === attributes.ai_engine || "llama3.1:8b",
          ).id,
          layout:
            attributes.layout !== null
              ? layouts.find((layout) => layout.name === attributes.layout).id
              : 1,
          background_align:
            attributes.background_align !== null ? attributes.background_align : "center",
          embedCode: `<iframe style=”width: 100%; height: 100%; border: 0px; border-radius: 10px; display: block;" src="${urlBot}"></iframe>`,
        };

        if (attributes.background.data !== null) {
          backgroundFile = getFile(
            bot.background,
            attributes.background.data.attributes.ext,
            attributes.background.data.attributes.name,
          );
        }

        if (attributes.bot_logo.data !== null) {
          logoFile = getFile(
            bot.bot_logo,
            attributes.bot_logo.data.attributes.ext,
            attributes.bot_logo.data.attributes.name,
          );
        }

        if (attributes.user_profile.data !== null) {
          userFile = getFile(
            bot.user_profile,
            attributes.user_profile.data.attributes.ext,
            attributes.user_profile.data.attributes.name,
          );
        }

        if (attributes.bot_profile.data !== null) {
          botFile = getFile(
            bot.bot_profile,
            attributes.bot_profile.data.attributes.ext,
            attributes.bot_profile.data.attributes.name,
          );
        }

        ready = true;
        triggerLoading(false);
      } catch (error) {
        console.log(error);
      }
    }
  }

  function getFile(base64, ext, name) {
    let fileType;
    if (ext === ".png") {
      fileType = "image/png";
    } else if (ext === ".jpg") {
      fileType = "image/jpg";
    } else if (ext === ".jpeg") {
      fileType = "image/jpeg";
    } else if (ext === ".svg") {
      fileType = "image/svg+xml";
    } else {
      fileType = "image/png";
    }

    let blob = base64toBlob(base64, fileType);
    let fileImageResult = URL.createObjectURL(blob);
    return {
      imageResult: fileImageResult,
      name: name,
      base64: base64,
    };
  }

  function base64toBlob(base64Data, fileType) {
    const byteCharacters = atob(base64Data);
    const byteArrays = [];

    for (let offset = 0; offset < byteCharacters.length; offset += 512) {
      const slice = byteCharacters.slice(offset, offset + 512);

      const byteNumbers = new Array(slice.length);
      for (let i = 0; i < slice.length; i++) {
        byteNumbers[i] = slice.charCodeAt(i);
      }

      const byteArray = new Uint8Array(byteNumbers);
      byteArrays.push(byteArray);
    }

    const blob = new Blob(byteArrays, { type: fileType });
    return blob;
  }

  async function updateBot() {
    let token = sessionStorage.getItem("token");
    if (token) {
      try {
        let response = await $call.call("client-dashboard", "update-bot", {
          URLQuery: {},
          URLParams: {
            id: botId,
          },
          Params: {
            additional_prompt: users.additionalPrompts,
            ai_engine: platforms.find((platform) => platform.id === users.platform).model,
            background: JSON.stringify(backgroundFile),
            background_align: users.background_align,
            bot_logo: JSON.stringify(logoFile),
            bot_name: users.botName,
            bot_profile: JSON.stringify(botFile),
            color_theme: themeColors.find((color) => color.id === users.primaryBrandColor).value,
            description: users.description,
            greeting_text: users.welcomeMessage,
            layout: layouts.find((layout) => layout.id === users.layout).name,
            logo_name: users.logoName,
            logo_only: users.logoOnly,
            overlay_color: themeColors.find((color) => color.id === users.overlayColor).value,
            user_profile: JSON.stringify(userFile),
            deployed_platform: users.deployedPlatform,
            username: users.username,
          },
          Headers: {
            Authorization: "Bearer " + token,
          },
        });

        console.log("update bot", response);
        triggerToast("updated");
      } catch (error) {
        console.log(error);
      }
    }
  }

  onMount(() => {
    triggerLoading(true);
    getModels();
    const unsubscribe = page.subscribe(async ($page) => {
      botId = $page.url.searchParams.get("botId");
      if (!botId) {
        triggerToast("error");
        triggerLoading(false);
        goto(await $call.getPage("client-dashboard", "hexabot-list"));
        return;
      }
      getBot();
    });
    return () => unsubscribe();
  });

  $: disabledButton = users.botName === "" || users.username === "";
</script>

<!-- <svelte:head>
  {@html github}
</svelte:head> -->

<!-- <svelte:window bind:innerWidth={screenWidth} /> -->

<!-- <LoadingSpinner {isLoading} /> -->
{#if ready}
  <div class="m-4 pb-12">
    <div class="sticky top-0 bg-primary-light dark:bg-primary-dark z-10">
      <div class="flex justify-start items-center mx-4">
        <Avatar
          class="size-14 border-[1px] border-hc-graylight-200"
          src={botFile.name ? botFile.imageResult : ""}
          stacked
        />
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
          <Typography
            preset="Text md/Semibold"
            text="Bots Information"
            color="secondary"
            propClass="w-full md:w-1/4"
          />

          <div class="w-full md:w-3/4 grid grid-cols-1 md:grid-cols-4 gap-3">
            <div class="flex flex-col gap-1.5 col-span-4">
              <Typography preset="Text sm/Medium" text="Bots Name" color="secondary" />
              <InputRw bind:value={users.botName} onEdit />
            </div>
            <div class="flex flex-col gap-1.5 col-span-4">
              <Typography preset="Text sm/Medium" text="Description (Optional)" color="secondary" />
              <TextareaRw bind:value={users.description} onEdit />
            </div>
          </div>
        </div>
        <Divider />
        <div class="md:flex">
          <Typography
            preset="Text md/Semibold"
            text="AI Engines"
            color="secondary"
            propClass="w-full md:w-1/4"
          />
          <div class="w-full md:w-3/4">
            <SelectedCard
              selectionData={platforms}
              bind:selectedData={users.platform}
              isFlex
              isRadio
              widthFlex="w-max"
            />
          </div>
        </div>
        <Divider />
        <div class="md:flex">
          <Typography
            preset="Text md/Semibold"
            text="Additional Prompts"
            color="secondary"
            propClass="w-full md:w-1/4"
          />
          <div
            class="w-full md:w-3/4 flex flex-col justify-center gap-2 [&>*]:bg-primary-light [&>*]:dark:bg-primary-dark"
          >
            <Typography preset="Text sm/Medium" text="Bot’s Role" color="secondary" />
            <TextareaRw bind:value={users.additionalPrompts} onEdit />
          </div>
        </div>
        <Divider />
        <div class="flex justify-end gap-3">
          <Button onclick={() => {}} color="light" class="px-3 focus-within:!ring-1"
            >{"Cancel"}</Button
          >
          <Button
            disabled={disabledButton}
            onclick={() => {
              updateBot();
            }}
            class="px-3 focus-within:!ring-1">{"Save"}</Button
          >
        </div>
      </div>
    {:else if tabDatas.selected === "Customization"}
      <div class="mt-6 space-y-6">
        <div class="md:flex">
          <Typography
            preset="Text md/Semibold"
            text="General"
            color="secondary"
            propClass="w-full md:w-1/4"
          />

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
                srcAvatar={botFile.name ? botFile.imageResult : ""}
              />
            </div>
            <div class="col-span-2">
              <UploadFile
                title="User Profile Picture"
                handleFileUploadWrapper={() => handleInputFile("user")}
                isPreview
                srcAvatar={userFile.name ? userFile.imageResult : ""}
              />
            </div>
          </div>
        </div>
        <Divider />
        <div class="md:flex">
          <Typography
            preset="Text md/Semibold"
            text="Logo"
            color="secondary"
            propClass="w-full md:w-1/4"
          />
          <div class="w-full md:w-3/4 grid grid-cols-1 md:grid-cols-4 gap-4">
            <div class="col-span-4">
              <Toggle bind:checked={users.logoOnly} class="w-max">Logo only</Toggle>
            </div>
            <div class="col-span-2">
              <UploadFile
                title="Logo"
                handleFileUploadWrapper={() => handleInputFile("logo")}
                isPreview
                srcAvatar={logoFile.name ? logoFile.imageResult : ""}
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
          <Typography
            preset="Text md/Semibold"
            text="Theme Colors"
            color="secondary"
            propClass="w-full md:w-1/4"
          />
          <div class="w-full md:w-3/4 grid grid-cols-1 md:grid-cols-4 gap-4">
            <div class="flex flex-col gap-1.5 col-span-2">
              <Typography preset="Text sm/Medium" text="Primary Brand Color" color="secondary" />
              <Button
                on:click={() => {
                  tmpThemeColor = users.primaryBrandColor;
                  modalThemeColor = true;
                }}
                class="flex justify-between border-primary-light dark:border-primary-dark dark:bg-gray-800"
                color="light"
              >
                <div class="flex items-center gap-2">
                  <Indicator color={getBrandColorName(users.primaryBrandColor)} />
                  {getBrandColorName(users.primaryBrandColor, "name")}
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
          <Typography
            preset="Text md/Semibold"
            text="Layout"
            color="secondary"
            propClass="w-full md:w-1/4"
          />
          <div class="w-full md:w-3/4">
            <SelectedCard
              selectionData={layouts}
              bind:selectedData={users.layout}
              isFlex
              isRadio
              widthFlex="w-full"
              spacing="!p-4"
            />
          </div>
        </div>
        <Divider />
        <div class="md:flex">
          <Typography
            preset="Text md/Semibold"
            text="Background"
            color="secondary"
            propClass="w-full md:w-1/4"
          />
          <div class="w-full md:w-3/4 grid grid-cols-1 md:grid-cols-4 gap-4">
            <div class="col-span-4">
              <UploadFile
                title="Background Picture"
                handleFileUploadWrapper={() => handleInputFile("background")}
                isPreview
                largePreview
                srcAvatar={backgroundFile.name ? backgroundFile.imageResult : ""}
              />
            </div>
            <div class="flex flex-col gap-1.5 col-span-2">
              <Typography preset="Text sm/Medium" text="Overlay Color" color="secondary" />
              <Button
                on:click={() => {
                  tmpOverlayColor = users.overlayColor;
                  modalOverlayColor = true;
                }}
                class="flex justify-between border-primary-light dark:border-primary-dark dark:bg-gray-800"
                color="light"
              >
                <div class="flex items-center gap-2">
                  <Indicator color={getBrandColorName(users.overlayColor)} />
                  {getBrandColorName(users.overlayColor, "name")}
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
                bind:selectedItem={users.background_align}
                items={bgAlign}
                placeholder={"Choose Collection"}
                isOverlay={true}
              />
            </div>
          </div>
        </div>
        <Divider />
        <div class="flex justify-end gap-3">
          <Button onclick={() => {}} color="light" class="px-3 focus-within:!ring-1"
            >{"Cancel"}</Button
          >
          <Button
            disabled={disabledButton}
            onclick={() => {
              updateBot();
            }}
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
          <Typography
            preset="Text md/Semibold"
            text="Embed Codes"
            color="secondary"
            propClass="w-full md:w-1/4"
          />

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
      bind:selectedData={tmpThemeColor}
      numGrid="2"
      isRadio
      widthFlex="w-full"
      spacing="!p-4"
      widgetIndicator={Indicator}
    />
    <svelte:fragment slot="footer">
      <Button
        onclick={() => (users.primaryBrandColor = tmpThemeColor)}
        class="px-3 w-full focus-within:!ring-1">{"Save"}</Button
      >
    </svelte:fragment>
  </Modal>

  <Modal title="Terms of Service" bind:open={modalOverlayColor} autoclose outsideclose>
    <SelectedCard
      selectionData={themeColors}
      bind:selectedData={tmpOverlayColor}
      numGrid="2"
      isRadio
      widthFlex="w-full"
      spacing="!p-4"
      widgetIndicator={Indicator}
    />
    <svelte:fragment slot="footer">
      <Button
        onclick={() => (users.overlayColor = tmpOverlayColor)}
        class="px-3 w-full focus-within:!ring-1">{"Save"}</Button
      >
    </svelte:fragment>
  </Modal>
{/if}
