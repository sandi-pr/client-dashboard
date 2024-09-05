<script>
  import { onMount } from "svelte";
  import { DarkMode, Button, P, Modal, Dropzone, Avatar, SidebarGroup } from "flowbite-svelte";
  import { Typography } from "$lib";
  import { format, isToday, isYesterday, set } from "date-fns";
  import SunIcon from "$lib/icon/sun.svg";
  // import Lajnah from "$lib/icon/lajnah.svg";
  import MenuIcon from "$lib/icon/menu-01.svg";
  import UploadIcon from "$lib/icon/upload-cloud-02.svg";
  import MessagePlusIcon from "$lib/icon/message-plus-circle-white.svg";
  import MessageChatCircleIcon from "$lib/icon/message-chat-circle.svg";
  import DataIcon from "$lib/icon/data.svelte";
  import DataFlow from "$lib/icon/dataflow-03.svelte";
  // import HexaBotLogo from "$lib/icon/hexabot-logo.svg";
  // import FileIcon from "$lib/icon/file-plus-02.svg";
  // import HomeIcon from "$lib/icon/home-01.svg";
  // import {
  //   CheckCircleOutline,
  //   TrashBinOutline,
  //   ArrowRightToBracketOutline,
  //   ChevronUpOutline,
  //   ChevronDownOutline,
  // } from "flowbite-svelte-icons";
  // import PDFSolidIcon from "$lib/icon/PDFSolid.svg";
  // import DOCSolidIcon from "$lib/icon/DOCSolid.svg";
  // import DOCXSolidIcon from "$lib/icon/DOCXSolid.svg";
  // import JPEGSolidIcon from "$lib/icon/JPEGSolid.svg";
  // import JPGSolidIcon from "$lib/icon/JPGSolid.svg";
  // import PNGSolidIcon from "$lib/icon/PNGSolid.svg";
  import "@fontsource/plus-jakarta-sans";
  import { setSelectedConversation, call } from "$lib";
  import { goto } from "$app/navigation";

  let who = "Crossroads Bot";
  let theme = "light";
  let {
    newChat,
    hidden = $bindable(true),
    name = $bindable("User"),
    email = $bindable("user@gmail.com"),
    logout = () => {},
    sidebarData = $bindable({}),
  } = $props();

  let clickOutsideModal = $state(false);
  let dropzoneActive = false;
  let fileTextName = $state(null);
  let fileTextResult = $state(null);
  let base64Data = $state("");
  let dropdownUser = $state(false);

  let fileList = $state([]);

  async function readFileAsText(file) {
    return new Promise((resolve, reject) => {
      const reader = new FileReader();

      reader.onload = () => {
        const text = reader.result.replace(/\r\n/g, "\n");
        const base64 = btoa(text);
        resolve({ text, base64 });
      };
      reader.onerror = reject;

      reader.readAsText(file);
    });
  }

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

  const handleUploadWrapper = () => {
    const inputElement = document.createElement("input");
    inputElement.type = "file";
    inputElement.accept = ".jpg, .jpeg, .png, .pdf, .doc, .docx";
    inputElement.multiple = true;

    inputElement.addEventListener("change", async (event) => {
      const files = event.target.files;

      for (let i = 0; i < files.length; i++) {
        const file = files[i];

        if (file) {
          let base64 = "";
          if (file.type === "text/plain") {
            const { text, base64: fileBase64 } = await readFileAsText(file);
            fileTextResult = text;
            base64 = fileBase64;
          } else {
            base64 = await readFileAsBase64(file);
          }

          const fileData = {
            name: file.name,
            size: file.size,
            base64Data: base64,
          };

          fileList = [...fileList, fileData];
        } else {
          alert("Invalid file. Please select a valid file.");
        }
      }
    });

    inputElement.click();
  };

  async function dropHandle(event) {
    event.preventDefault();

    const files = event.dataTransfer.files;
    const file = files[0];

    if (file) {
      if (file.type === "text/plain") {
        const { text, base64 } = await readFileAsText(file);
        fileTextResult = text;
        base64Data = base64;
      } else {
        base64Data = await readFileAsBase64(file);
      }

      const fileData = {
        name: file.name,
        size: file.size,
        base64Data: base64,
      };

      fileList = [...fileList, fileData];
    } else {
      alert("Invalid file. Please select a valid file.");
    }
  }

  function formatBytes(bytes, decimals = 2) {
    if (bytes === 0) return "0 Bytes";

    const k = 1024;
    const dm = decimals < 0 ? 0 : decimals;
    const sizes = ["Bytes", "KB", "MB", "GB", "TB", "PB", "EB", "ZB", "YB"];

    const i = Math.floor(Math.log(bytes) / Math.log(k));

    return parseFloat((bytes / Math.pow(k, i)).toFixed(dm)) + " " + sizes[i];
  }

  function removeFile(index) {
    fileList = fileList.filter((_, i) => i !== index);
  }

  function formatDate(dateString) {
    // console.log(dateString);
    const date = new Date(dateString);

    if (isToday(date)) {
      return "Today";
    } else if (isYesterday(date)) {
      return "Yesterday";
    } else {
      return format(date, "dd MMMM yyyy");
    }
  }

  const handleLogout = () => {
    logout();
  };

  async function processDataset() {
    let token = sessionStorage.getItem("token");
    if (token !== null) {
      try {
        const response = await $call.call("lajnah", "request-fine-tuning", {
          URLQuery: {},
          URLParams: {},
          Params: {},
          Headers: {
            Authorization: "Bearer " + token,
          },
        });

        if (response.status === 200) {
          console.log(response);
        }
      } catch (error) {
        console.log(error);
      }
    }
  }

  onMount(() => {
    // Initialize theme based on system preferences
    // const prefersDarkMode = window.matchMedia("(prefers-color-scheme: dark)").matches;
    // if (prefersDarkMode) {
    //   changeTheme();
    // }
  });
</script>

<div
  class={`h-full fixed z-10 w-full md:w-80 pt-9 pb-4 border-r border-secondary-light dark:border-secondary-dark bg-primary-light dark:bg-primary-dark flex flex-col justify-between items-start md:inline-flex ${
    hidden ? "hidden" : ""
  } `}
>
  <div
    class="flex-col w-full h-full px-4 md:px-8 justify-start items-center gap-4 md:gap-6 pb-4 md:pb-8 flex"
  >
    <div class="flex flex-col w-full justify-center items-center gap-2 pb-1">
      <!-- <img src={HexaBotLogo} alt="HexaBot Logo" class="" /> -->
      <!-- <img src={Lajnah} alt="Lajnah Logo" class="" /> -->
      <div class="flex gap-2 items-center">
        <DataIcon svgClass="text-utility-brand-700-light dark:text-utility-brand-700-dark" />
        <Typography preset="Display sm/Bold" text="Datasets" />
      </div>
    </div>
    <div class="flex flex-col gap-2 w-full">
      <Button
        on:click={newChat}
        class="!hidden lg:!flex self-stretch h-12 !text-white dark:!text-white rounded-2xl"
      >
        <img src={MessagePlusIcon} alt="MessagePlusIcon" class="w-5 h-5" />
        &nbsp; &nbsp;
        <P class="!text-white dark:!text-white">New Dataset</P>
      </Button>
      <Button
        on:click={() => {
          processDataset();
        }}
        outline
        class="!hidden lg:!flex hover:bg-primary-light dark:hover:bg-primary-dark self-stretch h-12 rounded-2xl"
      >
        <DataFlow
          svgClass="w-5 h-5 text-utility-brand-700-light dark:text-utility-brand-700-dark"
        />
        &nbsp; &nbsp;
        <Typography preset="Text md/Semibold" text="Process Dataset" color="brand-secondary" />
      </Button>
    </div>
    {#if Object.keys(sidebarData).length === 0}
      <div
        class="flex w-full h-full justify-center mb-3 rounded-2xl items-center bg-secondary-light dark:bg-secondary-dark"
      >
        <Typography preset="Text md/Regular" text="No sidebarData yet" color="fg-disabled" />
      </div>
    {:else}
      <div class="flex flex-col overflow-hidden hover:overflow-auto pb-20 w-full">
        {#each Object.keys(sidebarData) as date}
          <div
            class="flex flex-col w-full h-min py-4 justify-start mb-3 rounded-2xl items-center bg-secondary-light dark:bg-secondary-dark"
          >
            <Typography
              preset="Text md/Regular"
              text={formatDate(date)}
              color="disabled"
              propClass="mb-2"
            />
            {#each sidebarData[date] as conversation}
              <Button
                color="light"
                on:click={() => {
                  setSelectedConversation(conversation.id);
                  hidden = true;
                }}
                class="w-full gap-1.5 bg-transparent border-0"
              >
                <img src={MessageChatCircleIcon} alt="MessageChatCircleIcon" class="w-5 h-5" />
                <Typography preset="Text md/Semibold" text={conversation.title} color="secondary" />
              </Button>
            {/each}
          </div>
        {/each}
      </div>
    {/if}
  </div>

  <div
    class="absolute z-50 bottom-0 py-14 bg-gradient-to-t from-white dark:from-black from-80% w-full flex-col flex"
  >
    <div class="flex w-full justify-center px-8">
      <Button
        on:click={async () => {
          goto(await $call.getPage("client-dashboard", "datasets-list"));
        }}
        class="!bg-secondary-light dark:!bg-secondary-dark grow h-12 !text-slate-700 dark:!text-gray-300 w-full"
      >
        <P>See All</P>
      </Button>
    </div>
  </div>
</div>

<Modal bind:open={clickOutsideModal}>
  <svelte:fragment slot="header">
    <div class="flex flex-col">
      <Typography
        preset="Text lg/Semibold"
        text="Document Assistant"
        color="primary"
        propClass="fluid-base md:fluid-lg"
      />
      <Typography
        preset="Text sm/Regular"
        text="Register new document to the AI database"
        color="primary"
        propClass="fluid-sm md:fluid-base"
      />
    </div>
  </svelte:fragment>
  <div class="flex flex-col gap-2">
    <P
      class="hidden md:block text-sm font-medium leading-tight text-secondary-light dark:text-secondary-dark"
      >Document Upload</P
    >
    <Button
      ondrop={dropHandle}
      ondragover={(event) => {
        event.preventDefault();
      }}
      onclick={() => {
        handleUploadWrapper();
      }}
      color="light"
      class={`w-full rounded-xl border border-primary-light dark:border-primary-dark bg-primary-light dark:bg-primary-dark gap-1 py-4 md:py-8 flex md:flex-col items-center justify-start ${fileTextName ? "items-center cursor-default" : "items-center"}`}
    >
      <div
        class="rounded-xl shadow border border-primary-light dark:border-primary-dark justify-center items-center hidden md:inline-flex mb-4 pointer-events-none"
      >
        <img src={UploadIcon} alt="UploadIcon" class="p-4 dark:invert pointer-events-none" />
      </div>
      <div class="flex flex-col items-start md:items-center gap-1">
        <Typography
          preset="Text sm/Semibold"
          text="Click to upload"
          color="fg-brand-primary"
          propClass="fluid-base md:fluid-md"
        />
        <Typography
          preset="Text xs/Semibold"
          text="JPG, JPEG, PNG, PDF, DOC, DOCX  (Max. 15MB)"
          color="tertiary"
          propClass="fluid-xs md:fluid-base"
        />
      </div>
    </Button>
    <!-- {#if fileList && fileList.length > 0}
      <div class="flex flex-col gap-2">
        {#each fileList as file, index (file.name)}
          <div
            class="flex justify-between items-center bg-white px-3 py-3 rounded-lg border-2 border-hc-graylight-200"
          >
            <div class="flex justify-start items-center gap-3">
              {#if file.name.includes(".pdf") || file.name.includes(".PDF")}
                <img src={PDFSolidIcon} alt="" />
              {:else if file.name.includes(".docx") || file.name.includes(".DOCX")}
                <img src={DOCXSolidIcon} alt="" />
              {:else if file.name.includes(".doc") || file.name.includes(".DOC")}
                <img src={DOCSolidIcon} alt="" />
              {:else if file.name.includes(".jpg") || file.name.includes(".JPG")}
                <img src={JPGSolidIcon} alt="" />
              {:else if file.name.includes(".jpeg") || file.name.includes(".JPEG")}
                <img src={JPEGSolidIcon} alt="" />
              {:else if file.name.includes(".png") || file.name.includes(".PNG")}
                <img src={PNGSolidIcon} alt="" />
              {/if}
              <div class="flex flex-col justify-start items-start truncate">
                <p class="text-hc-graylight-600 text-sm font-semibold dark:text-gray-400">
                  {file.name}
                </p>
                <P
                  size="sm"
                  weight="light"
                  class="text-hc-graylight-600 dark:text-gray-400 text-sm"
                  >{formatBytes(file.size)}</P
                >
              </div>
            </div>
            <div class="flex justify-center items-center">
              <CheckCircleOutline class="w-5 h-5 ml-3 text-hc-brand-700" />
              <TrashBinOutline
                class="w-5 h-5 ml-3 text-hc-error-700 cursor-pointer"
                onclick={() => removeFile(index)}
              />
            </div>
          </div>
        {/each}
      </div>
    {/if} -->
  </div>
  <svelte:fragment slot="footer">
    <Button onclick={() => alert('Handle "success"')} class="w-full">Finish</Button>
  </svelte:fragment>
</Modal>
