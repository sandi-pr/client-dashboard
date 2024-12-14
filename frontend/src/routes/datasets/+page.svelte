<script>
  import {
    Button,
    Input,
    P,
    Heading,
    Span,
    Modal,
    Textarea,
    Spinner,
    DarkMode,
  } from "flowbite-svelte";
  import PdfViewer from "svelte-pdf";
  import { onMount, onDestroy } from "svelte";
  import { tick } from "svelte";
  import Markdown from "svelte-exmarkdown";
  import { fly } from "svelte/transition";
  import { goto } from "$app/navigation";
  import { gfmPlugin } from "svelte-exmarkdown/gfm";
  import { format, parse, sub } from "date-fns";
  import {
    Typography,
    triggerToast,
    setSelectedConversation,
    selectedConversation,
    call,
    triggerShowAI,
    triggerLoading,
    check,
  } from "$lib";
  import ToastContainer from "$lib/ui/toast/ToastContainer.svelte";
  import StateCard from "$lib/ui/cards/StateCard.svelte";
  import IllustrationEmpty from "$lib/icon/illustration_empty.svelte";
  import { CloseButton } from "flowbite-svelte";
  import Sidebar from "$lib/components/SidebarDataset.svelte";
  import { CR3Call } from "hexacomponents-svelte";
  import LaodingLottie from "$lib/components/LaodingLottie.svelte";
  import Chat from "$lib/components/Chat.svelte";
  import UserIcon from "$lib/icon/user-01.svg";
  import VerifiedIcon from "$lib/icon/verified.svg";
  import MenuIcon from "$lib/icon/menu-01.svg";
  // import ChevronLeftIcon from "$lib/icon/chevron-left.svg";
  import RefreshIcon from "$lib/icon/refresh-ccw-black.svg";
  // import FileIcon from "$lib/icon/file-plus-02.svg";
  // import FolderIcon from "$lib/icon/folder.svg";
  // import ImagePlusIcon from "$lib/icon/image-plus.svg";
  // import PenToolIcon from "$lib/icon/pen-tool-plus.svg";
  import PDFSolidIcon from "$lib/icon/PDFSolid.svg";
  import DOCSolidIcon from "$lib/icon/DOCSolid.svg";
  import DOCXSolidIcon from "$lib/icon/DOCXSolid.svg";
  import JPEGSolidIcon from "$lib/icon/JPEGSolid.svg";
  import JPGSolidIcon from "$lib/icon/JPGSolid.svg";
  import PNGSolidIcon from "$lib/icon/PNGSolid.svg";
  import EyeIcon from "$lib/icon/eye.svg";
  import ArrowLeftIcon from "$lib/icon/arrow-left.svg";
  import ArrowRightIcon from "$lib/icon/arrow-right.svg";
  import AlertCircleIcon from "$lib/icon/alert-circle.svg";
  import DotsVerticalIcon from "$lib/icon/dots-vertical.svelte";
  // import SendIcon from "$lib/icon/send-03.svg";
  import MessageChatCircleIcon from "$lib/icon/message-chat-circle.svg";
  // import MessagePlusIcon from "$lib/icon/message-plus-circle-white.svg";
  import MessagePlusIconGreen from "$lib/icon/message-plus-circle-green.svg";
  import EditIcon from "$lib/icon/edit-05.svelte";
  import AlignRightIcon from "$lib/icon/align-right-01.svelte";
  import FileIcon from "$lib/icon/file-04.svg";
  import { page } from "$app/stores";
  import { _ } from "svelte-i18n";
  import DownloadIcon from "$lib/icon/file-download-03.svg";
  import IllustrationProgress from "$lib/icon/illustration_progress.svelte";

  const plugins = [gfmPlugin()];
  const isFrontend = process.env.NODE_ENV === "development";

  let innerWidth = $state(0);
  let disableInput = $state(!check("aiManagement.datasets:create"));
  let hideSidebar = $state(true);
  let documentListModal = $state(false);
  let pageListDocument = $state(false);
  let previewModal = $state(false);
  let pagePreview = $state(false);
  let viewDocument = $state(false);
  let selectedFile = $state(null);
  let isLoading = $state(false);
  let allFileList = $state([]);
  let fileList = $state([]);
  let indexFileList = $state(0);
  let fileListReady = $state(false);
  let timer;
  let allData = $state([]);
  let files = $state([]);
  let conversations = $state([]);
  let initRAG = $state(false);
  let initDocList = $state(true);
  let showSendMessage = $state(false);
  let selectedRole = $state("User");
  let conversationId = $state(null);
  let editingMessageId = $state(null);
  let editedMessage = $state("");
  let deleteModal = $state(false);
  let dropdownOpen = $state(false);
  let gettingDatas = $state(true);

  let messageInput = $state("");

  let statusRAG = $state("Not Ready");

  let userData = $state({ username: "", email: "", brand: "" });

  let currentChatTitle = $state("New Chat");
  let editTitle = $state(false);

  let div;

  let ref;

  let textarea;

  let tmpResultTrocr = $state({
    text: "",
  });

  let sidebarData = $state({});

  let pageData = $state(1);

  let pageTotal = $state(1);

  let noAutoScroll = $state(false);

  let processGetData = $state(false);

  // let dataConversation = $state([
  //   {
  //     id: 0,
  //     chatId: [1, 2],
  //     date: new Date(),
  //   },
  // ]);

  let dataConversationChats = $state([
    // {
    //   id: 1,
    //   role: "assistant",
    //   message: `Hi! This is an AI assistant that is trained on Al-Quran. For example, you can ask for a Tafsir of a verse!
    // Now, how can I help you?`,
    //   date: "Chat " + format(new Date(), "dd.MM.yyyy"),
    // },
    // {
    //   id: 2,
    //   role: "user",
    //   message:
    //     "Hi! This is an AI assistant that is trained on Al-Quran. For example, you can ask for a Tafsir of a verse!",
    //   date: "Chat " + format(new Date(), "dd.MM.yyyy"),
    // },
  ]);

  async function dummySendMessage() {
    let init = false;
    let tmpConversationId;
    if (messageInput.trim() !== "") {
      console.log("messageInput", messageInput);
      const role = selectedRole === "Bot" ? "assistant" : "user";
      const message = messageInput.trim();

      if (selectedRole === "Bot") selectedRole = "User";
      else selectedRole = "Bot";

      messageInput = "";

      let token = localStorage.getItem("token");

      if (token === null) {
        triggerToast("error");
        return;
      }

      if (conversationId === null) {
        init = true;
        currentChatTitle = "New Chat " + format(new Date(), "dd.MM.yyyy");
        const response = await $call.call("client-dashboard", "create-conversation", {
          URLQuery: {},
          URLParams: {},
          Params: {
            title: currentChatTitle,
          },
          Headers: {
            Authorization: "Bearer " + token,
          },
        });

        if (response.status !== 200) {
          triggerToast("error");
          return;
        }

        conversationId = response.result.conversation_id;

        tmpConversationId = response.result.conversation_id;
      }

      noAutoScroll = false;

      let dataConversationChatsTmp = [...dataConversationChats];
      dataConversationChatsTmp.push({
        id: dataConversationChatsTmp.length + 1,
        role: role,
        message: message,
      });

      if (init) {
        const response2 = await $call.call("client-dashboard", "create-message", {
          URLQuery: {},
          URLParams: {},
          Params: {
            conversation_id: conversationId,
            message: JSON.stringify(dataConversationChatsTmp),
          },
          Headers: {
            Authorization: "Bearer " + token,
          },
        });

        if (response2.status !== 200) {
          triggerToast("error");
          return;
        }

        console.log("response2", response2);

        const newUrl =
          window.location.protocol +
          "//" +
          window.location.host +
          window.location.pathname +
          "?id=" +
          tmpConversationId;
        window.history.pushState({ path: newUrl }, "", newUrl);
        getListConversations();
      } else {
        const response2 = await $call.call("client-dashboard", "add-message", {
          URLQuery: {},
          URLParams: {},
          Params: {
            conversation_id: conversationId,
            message: JSON.stringify([
              {
                id: dataConversationChatsTmp.length,
                role: role,
                message: message,
              },
            ]),
          },
          Headers: {
            Authorization: "Bearer " + token,
          },
        });

        if (response2.status !== 200) {
          triggerToast("error");
          return;
        }
      }

      let currentDate = format(new Date(), "dd.MM.yyyy");
      let currentData = sidebarData[currentDate]?.find((data) => data.id === tmpConversationId);

      if (!currentData || currentData.title !== currentDate) {
        await getListConversations();

        currentData = sidebarData[currentDate]?.find((data) => data.id === tmpConversationId);
      }

      dataConversationChats = dataConversationChatsTmp;

      // addToDataset(role, message);
      messageInput = "";
    }
  }

  function startEditingMessage(id, message) {
    editingMessageId = id;
    editedMessage = message;
  }

  async function saveEditedMessage() {
    if (editingMessageId !== null) {
      let dataConversationChatsTmp = [...dataConversationChats];
      dataConversationChatsTmp = dataConversationChatsTmp.map((chat) => {
        if (chat.id === editingMessageId) {
          return { ...chat, message: editedMessage };
        }
        return chat;
      });

      let token = localStorage.getItem("token");

      const response = await $call.call("client-dashboard", "update-message", {
        URLQuery: {},
        URLParams: {},
        Params: {
          conversation_id: conversationId,
          message: JSON.stringify([
            {
              id: editingMessageId,
              role: "assistant",
              message: editedMessage,
            },
          ]),
        },
        Headers: {
          Authorization: "Bearer " + token,
        },
      });

      if (response.status !== 200) {
        triggerToast("error");
        return;
      }

      console.log("response", response);

      // datasetArray = [];
      // addToDatasetFromDB();

      dataConversationChats = dataConversationChatsTmp;
      editingMessageId = null;
      editedMessage = "";
    }
  }

  function cancelEditing() {
    editingMessageId = null;
    editedMessage = "";
  }

  function addToDataset(role, message) {
    const datasetEntry = {
      prompt: "",
      response: "",
    };

    if (role === "user") {
      datasetEntry.prompt = message;
    } else if (role === "assistant") {
      const lastUserMessage =
        dataConversationChats.length > 1 &&
        dataConversationChats[dataConversationChats.length - 2].role === "user"
          ? dataConversationChats[dataConversationChats.length - 2].message
          : "";

      datasetEntry.prompt = lastUserMessage;
      datasetEntry.response = message;

      saveDatasetEntry(datasetEntry);
    }
  }

  let datasetArray = $state([]);

  function saveDatasetEntry(entry) {
    console.log("Dataset Entry:", entry);
    datasetArray.push(entry);
  }

  async function clearConversation() {
    let token = localStorage.getItem("token");
    if (token) {
      const response = await $call.call("client-dashboard", "delete-conversation", {
        URLQuery: {},
        URLParams: {},
        Params: {
          conversation_id: conversationId,
        },
        Headers: {
          Authorization: "Bearer " + token,
        },
      });

      if (response.status !== 200) {
        triggerToast("error");
        return;
      }

      deleteModal = false;
      getListConversations();
      newChat();

      console.log("response", response);
    }
  }

  async function newChat() {
    datasetArray = [];
    currentChatTitle = "New Chat";
    conversationId = null;
    selectedRole = "User";
    const newUrl =
      window.location.protocol + "//" + window.location.host + window.location.pathname;

    window.history.pushState({ path: newUrl }, "", newUrl);
    setSelectedConversation(null);

    dataConversationChats = [];
  }

  const logout = async () => {
    sessionStorage.clear();
    localStorage.clear();

    try {
      isFrontend ? "" : await $call.goToPage("client-login", "home-page");
    } catch (error) {
      console.log(error);
    }
  };

  let winInnerWidth = $derived(innerWidth);

  function handleKeyDown(event) {
    if (winInnerWidth > 768) {
      if (event.key === "Enter") {
        if (!event.shiftKey) {
          event.preventDefault();
          if (messageInput.trim() !== "") {
            // sendMessage();
            dummySendMessage();
          }
        }
      }
    }
  }

  function onclickSend() {
    if (messageInput.trim() !== "") {
      // sendMessage();
      dummySendMessage();
    }
  }

  function autosize() {
    textarea.addEventListener("input", (e) => {
      setTimeout(() => {
        textarea.style.height = "auto";
        let scHeight = e.target.scrollHeight + 2; // Tambahkan sedikit untuk padding
        textarea.style.height = `${scHeight}px`;
      }, 0);
    });
  }

  async function autoscroll() {
    await tick();
    if (div) {
      div.scrollTo({ top: div.scrollHeight, behavior: "smooth" });
    }
  }

  const handleBeforePopState = (event) => {
    const message = "Are you sure you want to leave? Changes you made may not be saved.";
    event.preventDefault();
    event.returnValue = message; // Required for some browsers
    return message;
  };

  const handleBeforeUnload = (event) => {
    const message = "Are you sure you want to leave? Changes you made may not be saved.";
    event.preventDefault();
    event.returnValue = message; // Required for some browsers
    return message;
  };

  const handleChangeHome = () => {
    setHome(true);
  };

  async function getmessage(conversationId) {
    let token = localStorage.getItem("token");
    if (token) {
      const chatContainer = document.getElementById("chat-container");

      const scrollTop = chatContainer.scrollTop;
      const scrollHeight = chatContainer.scrollHeight;

      const response = await $call.call("client-dashboard", "get-messages", {
        URLQuery: {},
        URLParams: {},
        Params: {
          conversation_id: conversationId,
          limit: 100,
          page: pageData,
        },
        Headers: {
          Authorization: "Bearer " + token,
        },
      });

      if (response.status !== 200) {
        triggerToast("error");
        return;
      }

      currentChatTitle = response.result.attributes.title;
      dataConversationChats = [...response.result.attributes.message, ...dataConversationChats];
      pageTotal = response.result.page_total;

      await tick();
      const newScrollHeight = chatContainer.scrollHeight;
      const newScrollTop = scrollTop + (newScrollHeight - scrollHeight);

      chatContainer.scrollTop = newScrollTop;

      processGetData = false;

      let lastRole = dataConversationChats[dataConversationChats.length - 1].role;
      if (lastRole === "user") {
        selectedRole = "Bot";
      } else {
        selectedRole = "User";
      }
      // addToDatasetFromDB();
    }
  }

  function getPrevMessage() {
    pageData++;
    noAutoScroll = true;
    processGetData = true;
    getmessage(conversationId);
  }

  function addToDatasetFromDB() {
    datasetArray = [];
    for (let i = 0; i < dataConversationChats.length; i++) {
      const currentChat = dataConversationChats[i];
      const datasetEntry = {
        prompt: "",
        response: "",
      };

      if (
        currentChat.role === "assistant" &&
        i > 0 &&
        dataConversationChats[i - 1].role === "user"
      ) {
        datasetEntry.prompt = dataConversationChats[i - 1].message;
        datasetEntry.response = currentChat.message;

        saveDatasetEntry(datasetEntry);
      }
    }
  }

  async function updateTitle() {
    let token = localStorage.getItem("token");
    if (token) {
      const response = await $call.call("client-dashboard", "update-title", {
        URLQuery: {},
        URLParams: {},
        Params: {
          conversation_id: conversationId,
          title: currentChatTitle,
        },
        Headers: {
          Authorization: "Bearer " + token,
        },
      });

      if (response.status !== 200) {
        triggerToast("error");
        return;
      }
    }

    getListConversations();
    editTitle = false;

    triggerToast("updated");
  }

  async function getListConversations() {
    let token = localStorage.getItem("token");
    if (token) {
      const response = await $call.call("client-dashboard", "get-conversations", {
        URLQuery: {},
        URLParams: {},
        Params: {
          page: 1,
          page_size: 10,
        },
        Headers: {
          Authorization: "Bearer " + token,
        },
      });

      if (response.status !== 200) {
        triggerToast("error");
        return;
      }

      sidebarData = {};

      response.result.data.forEach((conversation) => {
        let updateAt = format(new Date(conversation.attributes.updatedAt), "yyyy-MM-dd");
        let title = conversation.attributes.title;
        let id = conversation.attributes.conversation_id;

        if (!sidebarData[updateAt]) {
          sidebarData[updateAt] = [];
        }

        sidebarData[updateAt].push({
          title: title,
          id: id,
        });
      });

      gettingDatas = false;
    }
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

        if (file && file.name.endsWith(".csv")) {
          let base64 = "";
          base64 = await readFileAsBase64(file);

          const fileData = {
            name: file.name,
            size: file.size,
            base64: base64,
          };

          fileList = [...fileList, fileData];
          selectedFile = file;
          importConversationWithCsv();
        } else {
          triggerToast("error");
        }
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

  async function importConversationWithCsv() {
    let token = localStorage.getItem("token");
    if (token) {
      console.log(fileList);
      try {
        if (conversationId === null) {
          currentChatTitle = "New Chat " + format(new Date(), "dd.MM.yyyy");
          const response = await $call.call("client-dashboard", "create-conversation", {
            URLQuery: {},
            URLParams: {},
            Params: {
              title: currentChatTitle,
            },
            Headers: {
              Authorization: "Bearer " + token,
            },
          });

          console.log("response", response);

          if (response.status !== 200) {
            triggerToast("error");
            return;
          }

          conversationId = response.result.conversation_id;

          let tmpConversationId = response.result.conversation_id;

          const newUrl =
            window.location.protocol +
            "//" +
            window.location.host +
            window.location.pathname +
            "?id=" +
            tmpConversationId;

          window.history.pushState({ path: newUrl }, "", newUrl);
          getListConversations();
        }

        const response = await $call.call("client-dashboard", "import-conversation-with-csv", {
          URLQuery: {},
          URLParams: {},
          Params: {
            conversation_id: conversationId,
            documents: JSON.stringify(fileList),
          },
          Headers: {
            Authorization: "Bearer " + token,
          },
        });

        if (response.status === 200) {
          getmessage(conversationId);
          triggerToast("uploaded");
        }
      } catch (error) {
        console.log(error);
      }
    }
  }

  let selectedTypeDataset = $state("conversation");

  let datas = $state([]);

  let configs = $state({});

  async function next() {
    // getCategories(pagination.page + 1);
  }

  async function previous() {
    // getCategories(pagination.page - 1);
  }

  onMount(() => {
    if (check("aiManagement.datasets:view") === false) {
      window.history.back();
      return;
    }
    triggerShowAI(false);
    userData = isFrontend
      ? { username: "John Doe", email: "johndoe@gmail.com" }
      : JSON.parse(localStorage.getItem("user"));

    if (textarea) {
      autosize();
    }

    const resizeHandler = () => {
      innerWidth = globalThis.innerWidth;
    };

    globalThis.addEventListener("resize", resizeHandler);
    resizeHandler(); // call once initially to set the correct value

    conversationId = $page.url.searchParams.get("id");

    if (conversationId) {
      getmessage(conversationId);
    }

    getListConversations();

    return () => {
      globalThis.removeEventListener("resize", resizeHandler);
    };
  });

  // $effect(async () => {
  //   const unsubscribe = page.subscribe(($page) => {
  //     conversationId = $page.url.searchParams.get("id");
  //     if (conversationId) {
  //       getmessage(conversationId);
  //     }
  //   });
  //   onDestroy(() => unsubscribe());
  // });

  $effect(() => {
    if ($selectedConversation !== null) {
      if (conversationId !== $selectedConversation) {
        dataConversationChats = [];
        getmessage($selectedConversation);
      }

      conversationId = $selectedConversation;
      const newUrl =
        window.location.protocol +
        "//" +
        window.location.host +
        window.location.pathname +
        "?id=" +
        $selectedConversation;

      window.history.pushState({ path: newUrl }, "", newUrl);

      pageData = 1;
      pageTotal = 1;
    }
  });

  $effect(() => {
    if (dataConversationChats.length) {
      if (!noAutoScroll) {
        autoscroll();
      }
    }
  });

  // use for handling sidebar data
  // $effect(() => {
  //   if (dataConversation.length) {
  //     let groupedConversations = {};

  //     dataConversation.forEach((log, index) => {
  //       // console.log(log);
  //       let date = log.date;

  //       if (!groupedConversations[date]) {
  //         groupedConversations[date] = [];
  //       }

  //       if (groupedConversations[date].find((conversation) => conversation.id === log.id)) {
  //         const conversationIndex = groupedConversations[date].findIndex(
  //           (conversation) => conversation.id === log.id
  //         );
  //         if (date > groupedConversations[date][conversationIndex].date) {
  //           groupedConversations[date][conversationIndex] = {
  //             id: log.id,
  //             title: "Chat " + format(date, "dd.MM.yyyy"),
  //           };
  //         }
  //       } else {
  //         groupedConversations[date].push({
  //           id: log.id,
  //           title: "Chat " + format(date, "dd.MM.yyyy"),
  //         });
  //       }
  //     });

  //     sidebarData = groupedConversations;
  //   }
  // });

  // used for handling conversation data
  // $effect(() => {
  //   if ($selectedConversation) {
  //     correlateConversations(dataConversation, dataConversationChats);
  //   }
  // });

  // $effect(() => {
  //   if ($isHome) {
  //     window.addEventListener("popstate", handleBeforePopState);
  //     window.addEventListener("beforeunload", handleBeforeUnload);
  //   } else {
  //     window.removeEventListener("popstate", handleBeforePopState);
  //     window.removeEventListener("beforeunload", handleBeforeUnload);
  //     window.addEventListener("popstate", handleChangeHome, { once: true });
  //   }
  // });

  // $effect(() => {
  //   if (pageListDocument || documentListModal) {
  //     initDocList = false;
  //   }
  // });

  // $effect(() => {
  //   if (!initDocList) {
  //     if (winInnerWidth < 768 && !pageListDocument) {
  //       if (!fileListReady) {
  //         triggerToast("getDocuments");
  //       }
  //     }
  //     if (winInnerWidth > 768 && !documentListModal) {
  //       if (!fileListReady) {
  //         triggerToast("getDocuments");
  //       }
  //     }
  //   }
  // });

  const handleDropdownFocusLoss = ({ relatedTarget, currentTarget }) => {
    if (relatedTarget instanceof HTMLElement && currentTarget.contains(relatedTarget)) return;
    dropdownOpen = false;
  };

  function autoResizeTextarea(event) {
    const textarea = event.target;
    textarea.style.height = "auto"; // Reset height
    textarea.style.height = textarea.scrollHeight + "px"; // Adjust height
  }

  onDestroy(() => {
    window.removeEventListener("popstate", handleBeforePopState);
    window.removeEventListener("beforeunload", handleBeforeUnload);
  });
</script>

<div
  bind:this={div}
  id="chat-container"
  class="static h-screen overflow-y-auto bg-primary-light dark:bg-primary-dark"
>
  <Sidebar
    bind:name={userData.username}
    bind:email={userData.email}
    {newChat}
    bind:hidden={hideSidebar}
    bind:sidebarData
    {logout}
    {getListConversations}
    bind:gettingDatas
    bind:selectedTypeDataset
  />
  {#if selectedTypeDataset === "conversation"}
    <div class={`flex flex-col lg:ml-80 pt-4 px-4 md:px-8`}>
      <div class="flex flex-col gap-3 md:gap-4">
        <div
          class="flex sticky top-0 z-20 py-4 bg-primary-light dark:bg-primary-dark w-full justify-between items-center"
        >
          <div class="lg:hidden flex w-max items-end justify-end pe-4">
            <Button
              outline
              on:click={() => (hideSidebar = !hideSidebar)}
              class="rounded-2xl !p-3 hover:bg-primary-light hover:dark:bg-primary-dark hover:text-primary-700 hover:dark:text-primary-300"
            >
              <img src={MenuIcon} alt="MenuIcon" class="w-5 h-5" />
            </Button>
          </div>
          <div class="flex justify-between items-center w-full">
            <div class="flex items-center gap-2">
              <!-- <img src={MessagePlusIconGreen} alt="MessagePlusIcon" class="w-8 h-8" /> -->
              {#if editTitle}
                <Input
                  bind:value={currentChatTitle}
                  class="w-48"
                  on:keydown={(e) => {
                    if (e.key === "Enter") {
                      updateTitle();
                      getListConversations();
                    }
                  }}
                  on:blur={() => {
                    editTitle = false;
                  }}
                />
              {:else}
                <Typography
                  preset="Text xl/Bold"
                  text={currentChatTitle}
                  color="primary"
                  propClass="fluid-lg line-clamp-1"
                />
              {/if}
              <!-- {#if conversationId !== null && dataConversationChats.length !== 0}
              <button
                on:click={() => {
                  editTitle = true;
                }}
                class="w-4 h-4 cursor-pointer"
              >
                <img src={EditIcon} alt="EditIcon" class="w-full h-full" />
              </button>
            {/if} -->
            </div>
            <div class="flex relative items-center justify-center gap-2">
              <!-- <button on:click={handleUploadWrapper}>import CSV</button>
            {#if selectedFile}
              <p>{selectedFile.name}</p>
            {/if}
            <button on:click={importConversationWithCsv}>Upload</button> -->
              <Button
                outline
                disabled={!check("aiManagement.datasets:edit")}
                on:click={() => {
                  editTitle = true;
                }}
                class="gap-1 h-max px-3 py-2 bg-secondary-light dark:bg-secondary-dark hover:bg-secondary-hover-light dark:hover:bg-secondary-hover-dark hover:text-brand-secondary-light dark:hover:text-brand-secondary-dark"
              >
                <EditIcon
                  svgClass="w-5 h-5 text-brand-secondary-light dark:text-brand-secondary-dark"
                />
                <!-- <img src={EditIcon} alt="Refresh Icon" class="w-5 h-5" /> -->
                <Typography
                  preset="Text md/Semibold"
                  text={$_("rename")}
                  color="brand-secondary"
                  propClass="hidden xl:block"
                />
              </Button>
              <Button
                outline
                disabled={!check("aiManagement.datasets:create")}
                on:click={handleUploadWrapper}
                class="gap-1 h-max px-3 py-2 bg-secondary-light dark:bg-secondary-dark hover:bg-secondary-hover-light dark:hover:bg-secondary-hover-dark hover:text-brand-secondary-light dark:hover:text-brand-secondary-dark"
              >
                <AlignRightIcon
                  svgClass="w-5 h-5 text-brand-secondary-light dark:text-brand-secondary-dark"
                />
                <Typography
                  preset="Text md/Semibold"
                  text={$_("importConversation")}
                  color="brand-secondary"
                  propClass="hidden xl:block"
                />
              </Button>
              <DarkMode
                btnClass="text-brand-secondary-light dark:text-brand-secondary-dark bg-secondary-light dark:bg-secondary-dark hover:bg-secondary-hover-light dark:hover:bg-secondary-hover-dark focus:outline-none rounded-lg px-3 py-2 border border-primary-light dark:border-primary-dark"
              />
              {#if check("aiManagement.datasets:delete")}
                <div on:focusout={handleDropdownFocusLoss} class="flex items-center relative">
                  <Button
                    outline
                    on:click={() => {
                      dropdownOpen = !dropdownOpen;
                    }}
                    class="gap-1 p-2 bg-secondary-light dark:bg-secondary-dark hover:bg-secondary-hover-light dark:hover:bg-secondary-hover-dark border-primary-light dark:border-primary-dark"
                  >
                    <DotsVerticalIcon
                      svgClass="w-5 h-5 text-utility-brand-700-light dark:text-utility-brand-700-dark"
                    />
                    <!-- <Typography preset="Text md/Semibold" text="Clear Conversation" color="secondary" /> -->
                  </Button>
                  <div
                    class={`z-50 ${dropdownOpen ? "absolute -bottom-16 right-0" : "hidden"} bg-white divide-y divide-gray-100 rounded-lg shadow w-max dark:bg-gray-700`}
                  >
                    <ul class="py-2 text-sm text-gray-700 dark:text-gray-200">
                      <button
                        class="flex p-2"
                        on:click={() => {
                          clearConversation();
                          dropdownOpen = false;
                        }}>{$_("delete")}</button
                      >
                    </ul>
                  </div>
                </div>
              {/if}
            </div>
          </div>
        </div>
        <div class="flex flex-col pb-40 gap-4">
          {#if pageData < pageTotal && dataConversationChats.length !== 0}
            <div class="flex w-full justify-center">
              <Button color="light" on:click={getPrevMessage} class="w-max"
                >Get Previous Messages &nbsp; {#if processGetData}<Spinner />{/if}</Button
              >
            </div>
          {/if}
          {#each dataConversationChats as qaPair}
            <div
              class="flex w-full bg-secondary-light dark:bg-secondary-dark rounded-2xl p-4 md:px-8 md:py-6 items-start gap-4"
            >
              <div
                class="w-8 h-7 md:w-12 md:h-11 relative hidden md:flex justify-center items-center bg-tertiary-light dark:bg-tertiary-dark rounded-full border border-primary-light dark:border-primary-dark"
              >
                <img src={UserIcon} alt="User" class="w-4 h-4 md:w-6 md:h-6" />
                {#if qaPair.role === "assistant"}
                  <div class="absolute bottom-0 right-0">
                    <img src={VerifiedIcon} alt="User" class="w-2.5 h-2.5 md:w-4 md:h-4" />
                  </div>
                {/if}
              </div>
              <div class="flex flex-col relative w-full gap-3 md:gap-1">
                <div class="flex items-center justify-between gap-2.5 rtl:space-x-reverse">
                  <div
                    class="w-9 h-9 md:w-12 md:h-11 relative flex md:hidden justify-center items-center bg-tertiary-light dark:bg-tertiary-dark rounded-full border border-primary-light dark:border-primary-dark"
                  >
                    <img src={UserIcon} alt="User" class="w-5 h-5 md:w-6 md:h-6" />
                    {#if qaPair.role === "assistant"}
                      <div class="absolute bottom-0 right-0">
                        <img src={VerifiedIcon} alt="User" class="w-2.5 h-2.5 md:w-4 md:h-4" />
                      </div>
                    {/if}
                  </div>
                  <Typography
                    preset="Text md/Bold"
                    text={qaPair.role === "user" ? "User" : "Bot"}
                    color="primary"
                    propClass="fluid-base md:fluid-base"
                  />

                  {#if editingMessageId !== qaPair.id}
                    <button on:click={() => startEditingMessage(qaPair.id, qaPair.message)}>
                      <EditIcon
                        svgClass="w-full h-full text-brand-secondary dark:text-brand-secondary"
                      />
                    </button>
                  {/if}
                </div>
                <!-- {#if qaPair?.documents && editingMessageId !== qaPair.id}
                <Button
                  color="green"
                  class="absolute right-6 -top-3 text-utility-success-500-light dark:text-white font-semibold bg-hc-success-50 hover:bg-hc-success-100 dark:bg-hc-success-500 hover:dark:bg-hc-success-400 rounded-2xl border border-utility-success-500-light dark:border-utility-success-500-dark"
                >
                  <img src={FileIcon} alt="File Icon" class="w-5 h-5" />
                  &nbsp; &nbsp; Document Retrieval</Button
                >
              {/if} -->
                {#if editingMessageId === qaPair.id}
                  <Textarea
                    bind:value={editedMessage}
                    on:input={autoResizeTextarea}
                    unWrappedClass="focus:border-fg-success-secondary-light dark:focus:border-fg-success-secondary-dark focus:ring-0"
                  ></Textarea>
                  <div class="flex gap-2 justify-end">
                    <Button on:click={saveEditedMessage} class="h-max">Save</Button>
                    <Button color="light" on:click={cancelEditing} class="h-max">Cancel</Button>
                  </div>
                {:else}
                  <div
                    class="prose max-w-none font-medium fluid-base md:fluid-base text-tertiary-light dark:text-primary-dark"
                  >
                    <Markdown md={qaPair.message} {plugins} />
                  </div>
                {/if}
              </div>
            </div>
          {/each}
          <!-- {#if disableInput}
          <div
            class="flex w-full bg-secondary-light dark:bg-secondary-dark rounded-2xl p-4 md:px-8 md:py-6 items-start gap-4"
          >
            <div
              class="w-8 h-7 md:w-12 md:h-11 relative hidden md:flex justify-center items-center bg-tertiary-light dark:bg-tertiary-dark rounded-full border border-primary-light dark:border-primary-dark"
            >
              <img src={UserIcon} alt="User" class="w-4 h-4 md:w-6 md:h-6" />
              <div class="absolute bottom-0 right-0">
                <img src={VerifiedIcon} alt="User" class="w-2.5 h-2.5 md:w-4 md:h-4" />
              </div>
            </div>
            <div class="flex flex-col w-full gap-3 md:gap-1">
              <div class="flex items-center gap-2.5 rtl:space-x-reverse">
                <div
                  class="w-9 h-9 md:w-12 md:h-11 relative flex md:hidden justify-center items-center bg-tertiary-light dark:bg-tertiary-dark rounded-full border border-primary-light dark:border-primary-dark"
                >
                  <img src={UserIcon} alt="User" class="w-5 h-5 md:w-6 md:h-6" />
                  <div class="absolute bottom-0 right-0">
                    <img src={VerifiedIcon} alt="User" class="w-2.5 h-2.5 md:w-4 md:h-4" />
                  </div>
                </div>
                <Typography
                  preset="Text md/Bold"
                  text={"Bot"}
                  color="primary"
                  propClass="fluid-base md:fluid-base"
                />
              </div>
              <div
                class="prose max-w-none font-medium fluid-base md:fluid-base text-tertiary-light dark:text-primary-dark"
              >
                <Markdown md={"Qur’ani is typing..."} {plugins} />
              </div>
            </div>
          </div>
        {/if} -->
        </div>
      </div>
    </div>

    <Chat
      {textarea}
      bind:messageInput
      {handleKeyDown}
      {onclickSend}
      {isLoading}
      {disableInput}
      {showSendMessage}
      bind:hideSidebar
      bind:selectedRole
    />
  {:else}
    <div class="flex flex-col lg:ml-80 pt-4 px-4 md:px-8 dark:bg-primary-dark">
      <div class="w-full pb-10">
        <div class="flex flex-col m-4 pt-2">
          <!-- <Navbar /> -->

          <div class="flex flex-col">
            <!-- User Information -->
            <div class="flex items-center justify-between">
              <Typography
                preset="Display sm/Semibold"
                color="primary"
                text={$_("allDatasetsHandwriting")}
              />
              <Button onclick={() => {}} icon={DownloadIcon}
                ><img class="w-4 h-4 me-1" src={DownloadIcon} alt="plus-icon" />
                {$_("addNew")}</Button
              >
            </div>

            <!-- Person Table Data -->
            {#if datas.length > 0}
              <DataTable bind:datas bind:configs bind:next bind:previous />
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
          </div>
        </div>
      </div>
    </div>
  {/if}

  {#if pagePreview}
    <div
      class={`fixed inset-0 z-40 flex overflow-auto flex-col items-center justify-between ${viewDocument ? "bg-black bg-opacity-80" : "bg-white"} `}
    >
      <div class="flex flex-col items-center pb-6">
        {#if viewDocument}
          <div class="flex w-full justify-end items-end pt-12 pr-4">
            <CloseButton
              on:click={() => {
                pagePreview = false;
                viewDocument = false;
              }}
              class="text-white"
            />
          </div>
        {:else}
          <div
            class="flex sticky top-0 bg-primary-light dark:bg-primary-dark justify-start items-center gap-4 py-4 px-6 border-b border-secondary-light dark:border-secondary-dark"
          >
            <div class="w-8 h-7 p-2 bg-[#feefc6] rounded-full flex justify-center items-center">
              <img src={AlertCircleIcon} alt="AlertCircleIcon" class="w-5 h-5" />
            </div>
            <div class="flex flex-col">
              <Typography preset="Text lg/Semibold" text="Double-check" color="primary" />
              <Typography
                preset="Text sm/Regular"
                text="Are you sure this is the correct document that you want to upload?"
                color="tertiary"
              />
            </div>
          </div>
        {/if}
        <div
          class={`mt-4 flex flex-col m-4 py-12 px-6 ${viewDocument ? "bg-transparent" : "bg-secondary-light dark:bg-secondary-dark"} rounded-lg items-center justify-center`}
        >
          {#if selectedFile.ext.includes(".png") || selectedFile.ext.includes(".jpg") || selectedFile.ext.includes(".jpeg")}
            <img
              src={`data:${selectedFile.type};base64,${selectedFile.base64}`}
              style="max-width: 100%; height: auto;"
              alt="PNG Preview"
            />
          {:else if selectedFile.ext.includes(".pdf")}
            <PdfViewer
              showBorder={false}
              showButtons={["navigation", "zoom", "rotate", "pageInfo"]}
              scale={1}
              data={atob(selectedFile.base64)}
            />
          {/if}
        </div>
        <!-- detected handwriting -->
        {#if selectedFile.type === "Handwriting"}
          <div class="flex w-full flex-col gap-3 px-4 pb-1">
            <Typography
              preset="Text md/Semibold"
              text="Handwriting Information"
              color="secondary"
            />
            <div class="flex flex-col gap-1.5">
              <Typography preset="Text sm/Medium" text="Detected Text" color="secondary" />
              <textarea
                bind:value={tmpResultTrocr.text}
                id="chat"
                rows="4"
                class="block p-2.5 w-full text-base text-brand-primary-light dark:text-primary-dark bg-white dark:bg-primary-dark rounded-lg border border-utility-brand-500-light dark:border-utility-brand-500-dark focus:border-utility-brand-500-light dark:focus:utility-brand-500-dark focus:ring-utility-brand-400-light dark:focus:ring-utility-brand-400-light dark:placeholder-gray-400"
              ></textarea>
            </div>
          </div>
        {/if}
      </div>
      {#if !viewDocument}
        <div
          class="flex sticky bottom-0 w-full bg-primary-light gap-2 justify-center items-center py-3 px-4"
        >
          <Button
            color="light"
            on:click={() => {
              pagePreview = false;
              files = [];
            }}
            class="w-full"
          >
            Cancel
          </Button>
          <Button color="primary" on:click={handleConfirm} class="rounded-lg w-full">
            I’m sure
          </Button>
        </div>
      {/if}
    </div>
  {/if}

  {#if isLoading}
    <div class="fixed inset-0 z-50 flex flex-col items-center justify-center bg-[#8340c4]">
      <LaodingLottie />
      <div class="flex flex-col gap-2 items-center pt-4">
        <Typography
          preset="Text xl/Bold"
          color="primary-onbrand"
          text="Please wait, AI is on training!"
        />
        <Typography
          preset="Text sm/Medium"
          color="primary-onbrand"
          text="This loading takes a few second to few minutes at most..."
          propClass="text-center"
        />
      </div>
    </div>
  {/if}

  <ToastContainer />
</div>

<Modal bind:open={documentListModal} bodyClass="p-4 md:p-5 space-y-4 flex-1 overflow-x-hidden">
  <svelte:fragment slot="header">
    <div class="flex flex-col">
      <Typography
        preset="Text lg/Semibold"
        text="Document List"
        color="primary"
        propClass="fluid-base md:fluid-lg"
      />
    </div>
  </svelte:fragment>
  <div class="flex flex-col gap-2">
    {#if fileList[indexFileList] && fileList[indexFileList].length > 0}
      <P
        class="hidden md:block text-sm font-medium leading-tight text-secondary-light dark:text-secondary-dark"
        >File name</P
      >
      <div class="flex flex-col gap-2">
        {#each fileList[indexFileList] as file, index (file.name)}
          <div
            class="flex justify-between items-center px-3 py-3 rounded-lg border-b border-primary-light dark:border-primary-dark"
          >
            <div class="flex justify-start items-center gap-3">
              {#if file.ext.includes(".pdf") || file.ext.includes(".PDF")}
                <img src={PDFSolidIcon} alt="" />
              {:else if file.ext.includes(".docx") || file.ext.includes(".DOCX")}
                <img src={DOCXSolidIcon} alt="" />
              {:else if file.ext.includes(".doc") || file.ext.includes(".DOC")}
                <img src={DOCSolidIcon} alt="" />
              {:else if file.ext.includes(".jpg") || file.ext.includes(".JPG")}
                <img src={JPGSolidIcon} alt="" />
              {:else if file.ext.includes(".jpeg") || file.ext.includes(".JPEG")}
                <img src={JPEGSolidIcon} alt="" />
              {:else if file.ext.includes(".png") || file.ext.includes(".PNG")}
                <img src={PNGSolidIcon} alt="" />
              {/if}
              <div
                class="flex flex-col justify-start items-start max-w-48 md:max-w-max text-ellipsis"
              >
                <p
                  class="text-hc-graylight-600 text-sm text-ellipsis font-semibold dark:text-gray-400"
                >
                  {file.name}
                </p>
              </div>
            </div>
            <div class="flex justify-center items-center">
              <!-- <CheckCircleOutline class="w-5 h-5 ml-3 text-hc-brand-700" />
              <TrashBinOutline
                class="w-5 h-5 ml-3 text-hc-error-700 cursor-pointer"
                onclick={() => removeFile(index)}
              /> -->
              <button
                class="flex items-center gap-1.5"
                on:click={async () => {
                  goto(
                    (isFrontend ? "" : await $call.getPage("lajnah", "document-detail")) +
                      "?collectionId=" +
                      sessionStorage.getItem("collectionIdFinetune") +
                      "&documentId=" +
                      file.id,
                  );
                  return;
                }}
              >
                <img src={EyeIcon} alt="EyeIcon" class="w-5 h-5" />
                <Typography
                  preset="Text sm/Semibold"
                  text="View"
                  color="brand-primary"
                  propClass="hidden sm:block fluid-base md:fluid-md"
                />
              </button>
            </div>
          </div>
        {/each}
      </div>
    {:else}
      <StateCard
        illustration={IllustrationEmpty}
        description={fileListReady ? "You don't have any document yet" : "Getting document list..."}
      />
    {/if}
  </div>
  <svelte:fragment slot="footer">
    {#if fileList[indexFileList] && fileList[indexFileList].length > 0}
      <div class="flex w-full justify-between items-center">
        <Button
          outline
          color="light"
          disabled={indexFileList === 0}
          on:click={() => {
            if (indexFileList > 0) {
              indexFileList = indexFileList - 1;
            }
          }}
          class="rounded-2xl !p-3 hover:bg-primary-light hover:dark:bg-primary-dark"
        >
          <img src={ArrowLeftIcon} alt="ArrowLeftIcon" class="w-5 h-5" />
        </Button>
        <Typography
          preset="Text sm/Regular"
          text={indexFileList + 1 + " of " + fileList.length}
          color="brand-primary"
          propClass="fluid-base md:fluid-md"
        />
        <Button
          outline
          color="light"
          disabled={indexFileList === fileList.length - 1}
          on:click={() => {
            if (indexFileList < fileList.length - 1) {
              indexFileList = indexFileList + 1;
            }
          }}
          class="rounded-2xl !p-3 hover:bg-primary-light hover:dark:bg-primary-dark"
        >
          <img src={ArrowRightIcon} alt="ArrowRightIcon" class="w-5 h-5" />
        </Button>
      </div>
    {/if}
  </svelte:fragment>
</Modal>

<Modal bind:open={previewModal} size="lg">
  <svelte:fragment slot="header">
    <div class="flex justify-start items-center gap-4 py-4 px-6">
      <div class="w-8 h-7 p-2 bg-[#feefc6] rounded-full flex justify-center items-center">
        <img src={AlertCircleIcon} alt="AlertCircleIcon" class="w-5 h-5" />
      </div>
      <div class="flex flex-col">
        <Typography preset="Text lg/Semibold" text="Double-check" color="primary" />
        <Typography
          preset="Text sm/Regular"
          text="Are you sure this is the correct document that you want to upload?"
          color="tertiary"
        />
      </div>
    </div>
  </svelte:fragment>
  <div class="flex">
    <div
      class="flex w-1/2 justify-center items-center py-6 rounded-lg bg-secondary-light dark:bg-secondary-dark"
    >
      {#if selectedFile.ext.includes(".png") || selectedFile.ext.includes(".jpg") || selectedFile.ext.includes(".jpeg")}
        <img
          src={`data:${selectedFile.type};base64,${selectedFile.base64}`}
          style="max-width: 80%; height: auto;"
          alt="PNG Preview"
        />
      {:else if selectedFile.ext.includes(".pdf")}
        <PdfViewer
          showBorder={false}
          showButtons={false}
          scale={1}
          data={atob(selectedFile.base64)}
        />
      {/if}
    </div>
    <!-- detected handwriting -->
    {#if selectedFile.type === "Handwriting"}
      <div class="flex w-full flex-col gap-3 px-4 pb-1">
        <Typography preset="Text md/Semibold" text="Handwriting Information" color="secondary" />
        <div class="flex flex-col gap-1.5">
          <Typography preset="Text sm/Medium" text="Detected Text" color="secondary" />
          <textarea
            bind:value={tmpResultTrocr.text}
            id="chat"
            rows="4"
            class="block p-2.5 w-full text-sm text-brand-primary-light dark:text-primary-dark bg-white dark:bg-primary-dark rounded-lg border border-utility-brand-500-light dark:border-utility-brand-500-dark focus:border-utility-brand-500-light dark:focus:utility-brand-500-dark focus:ring-utility-brand-400-light dark:focus:ring-utility-brand-400-light dark:placeholder-gray-400"
          ></textarea>
        </div>
      </div>
    {/if}
  </div>
  <svelte:fragment slot="footer">
    <div class="flex w-full gap-2 justify-center items-center py-3 px-4">
      <Button
        color="light"
        on:click={() => {
          previewModal = false;
          files = [];
        }}
        class="w-full">Cancel</Button
      >
      <Button color="primary" on:click={handleConfirm} class="rounded-lg w-full">I’m sure</Button>
    </div>
  </svelte:fragment>
</Modal>
