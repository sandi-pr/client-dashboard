<script>
  import { onMount } from "svelte";
  import {
    DarkMode,
    Button,
    P,
    Modal,
    Dropzone,
    Avatar,
    SidebarGroup,
    Indicator,
    ListPlaceholder,
  } from "flowbite-svelte";
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
  import AlignRightIcon from "$lib/icon/align-right-01.svelte";
  import { _ } from "svelte-i18n";
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
  import { setSelectedConversation, call, triggerToast } from "$lib";
  import { goto } from "$app/navigation";

  let who = "Crossroads Bot";
  let theme = "light";
  let {
    newChat,
    hidden = $bindable(true),
    name = $bindable("User"),
    email = $bindable("user@gmail.com"),
    logout = () => {},
    getListConversations = () => {},
    sidebarData = $bindable({}),
    gettingDatas = $bindable(true),
  } = $props();

  let clickOutsideModal = $state(false);
  let dropzoneActive = false;
  let fileTextName = $state(null);
  let fileTextResult = $state(null);
  let base64Data = $state("");
  let dropdownUser = $state(false);
  let statusRequest = $state(null);
  let colorStatusRequest = $state("dark");

  let fileList = $state([]);
  let selectedFile = $state(null);

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
    inputElement.accept = ".csv";

    inputElement.addEventListener("change", async (event) => {
      const files = event.target.files;

      fileList = [];

      for (let i = 0; i < files.length; i++) {
        const file = files[i];

        if (file) {
          let base64 = "";
          base64 = await readFileAsBase64(file);

          const fileData = {
            name: file.name,
            size: file.size,
            base64: base64,
          };

          fileList = [...fileList, fileData];
          selectedFile = file;
        } else {
          alert("Invalid file. Please select a valid file.");
        }
      }

      importDatasetsWithCsv();
    });

    inputElement.click();
  };

  async function importDatasetsWithCsv() {
    let token = sessionStorage.getItem("token");
    if (token) {
      try {
        const response = await $call.call("client-dashboard", "import-datasets-with-csv", {
          URLQuery: {},
          URLParams: {},
          Params: {
            documents: JSON.stringify(fileList),
          },
          Headers: {
            Authorization: "Bearer " + token,
          },
        });

        if (response.status === 200) {
          triggerToast("created");
          getListConversations();
        }
      } catch (error) {
        console.log(error);
      }
    }
  }

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
        const response = await $call.call("client-dashboard", "request-fine-tuning", {
          URLQuery: {},
          URLParams: {},
          Params: {},
          Headers: {
            Authorization: "Bearer " + token,
          },
        });

        if (response.status === 200) {
          statusRequest = "REQUEST";
          triggerToast("created");
        }
      } catch (error) {
        console.log(error);
      }
    }
  }

  async function getStatusRequest() {
    let token = sessionStorage.getItem("token");
    if (token !== null) {
      try {
        const response = await $call.call("client-dashboard", "get-status-request-fine-tuning", {
          URLQuery: {},
          URLParams: {},
          Params: {},
          Headers: {
            Authorization: "Bearer " + token,
          },
        });
        if (response.status === 200) {
          if (response.result.data.length > 0) {
            statusRequest = response.result.data[0].attributes.status;
          }
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
    getStatusRequest();
  });

  $effect(() => {
    if (statusRequest !== null) {
      if (statusRequest === "REQUEST") colorStatusRequest = "blue";
      if (statusRequest === "PROCESS") colorStatusRequest = "yellow";
      if (statusRequest === "SUCCESS") colorStatusRequest = "green";
      if (statusRequest === "FAILED") colorStatusRequest = "red";
    }
  });
</script>

<div
  class={`h-full fixed z-10 w-full md:w-80 pt-9 pb-4 border-r border-secondary-light dark:border-secondary-dark bg-primary-light dark:bg-primary-dark flex flex-col justify-between items-start lg:inline-flex ${
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
        <Indicator color={colorStatusRequest} />
      </div>
    </div>
    <div class="flex flex-col gap-2 w-full">
      <Button
        on:click={newChat}
        class="!flex self-stretch h-12 !text-white dark:!text-white rounded-2xl"
      >
        <img src={MessagePlusIcon} alt="MessagePlusIcon" class="w-5 h-5" />
        &nbsp; &nbsp;
        <P class="!text-white dark:!text-white">New Dataset</P>
      </Button>
      <div class="grid grid-cols-2 gap-2">
        <Button
          disabled={statusRequest === "process"}
          on:click={() => {
            processDataset();
          }}
          outline
          class="!flex hover:bg-primary-light dark:hover:bg-primary-dark self-stretch h-12 rounded-2xl"
        >
          <DataFlow
            svgClass="w-5 h-5 text-utility-brand-700-light dark:text-utility-brand-700-dark"
          />
          &nbsp; &nbsp;
          <Typography preset="Text md/Semibold" text={$_("process")} color="brand-secondary" />
        </Button>
        <Button
          disabled={statusRequest === "process"}
          on:click={handleUploadWrapper}
          outline
          class="!flex hover:bg-primary-light dark:hover:bg-primary-dark self-stretch h-12 rounded-2xl"
        >
          <AlignRightIcon svgClass="w-5 h-5 text-utility-brand-700-light dark:text-utility-brand-700-dark" />
          &nbsp; &nbsp;
          <Typography preset="Text md/Semibold" text={$_("import")} color="brand-secondary" />
        </Button>
      </div>
    </div>
    {#if Object.keys(sidebarData).length === 0}
      {#if gettingDatas}
        <div class="flex grow">
          <ListPlaceholder />
        </div>
      {:else}
        <div
          class="flex w-full h-min py-4 justify-center mb-3 rounded-2xl items-center bg-secondary-light dark:bg-secondary-dark"
        >
          <Typography preset="Text md/Regular" text="No datasets yet" color="fg-disabled" />
        </div>
      {/if}
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
                class="w-full gap-1.5 bg-secondary-light dark:bg-secondary-dark border-0"
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
    class="absolute z-50 bottom-0 py-14 bg-gradient-to-t from-white dark:from-utility-bg-primary-dark from-80% w-full flex-col flex"
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
