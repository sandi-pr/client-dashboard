<script>
  import { Input, Button, Spinner, Textarea, DarkMode } from "flowbite-svelte";
  import { _ } from "svelte-i18n";
  import { page } from "$app/stores";
  import Markdown from "svelte-exmarkdown";
  import { format } from "date-fns";
  import { onMount, onDestroy, tick } from "svelte";
  import { gfmPlugin } from "svelte-exmarkdown/gfm";

  import {
    call,
    check,
    Typography,
    triggerToast,
    triggerShowAI,
    triggerLoading,
    selectedConversation,
    setSelectedConversation,
    isBakrie,
  } from "$lib";
  import Chat from "$lib/components/Chat.svelte";

  import UserIcon from "$lib/icon/user-01.svg";
  import MenuIcon from "$lib/icon/menu-01.svg";
  import VerifiedIcon from "$lib/icon/verified.svg";

  import EditIcon from "$lib/icon/edit-05.svelte";
  import AlignRightIcon from "$lib/icon/align-right-01.svelte";
  import DotsVerticalIcon from "$lib/icon/dots-vertical.svelte";
  import ToastContainer from "$lib/ui/toast/ToastContainer.svelte";
  import LaodingLottie from "$lib/components/LaodingLottie.svelte";

  const plugins = [gfmPlugin()];
  const isFrontend = process.env.NODE_ENV === "development";

  let div;
  let textarea;

  let isLoading = $state(false);
  let editTitle = $state(false);
  let hideSidebar = $state(true);
  let deleteModal = $state(false);
  let gettingDatas = $state(true);
  let dropdownOpen = $state(false);
  let noAutoScroll = $state(false);
  let processGetData = $state(false);
  let showSendMessage = $state(false);
  // let fileDetailModal = $state(true);
  let disableInput = $state(!check("aiManagement.datasets:create"));

  let pageData = $state(1);
  let pageTotal = $state(1);
  let innerWidth = $state(0);

  let messageInput = $state("");
  let editedMessage = $state("");
  let selectedRole = $state("User");
  let selectedTypeDataset = $state($isBakrie ? "handwriting" : "conversation");
  // let typeDataset = $state(["handwriting"]);
  let disableButtonHandwriting = $state(true);
  let disableButtonConversation = $state(true);
  let disabledButtonProccess = $state(false);
  let currentChatTitle = $state("New Chat");

  $effect(() => {
    if (selectedTypeDataset === "conversation") {
      disabledButtonProccess = disableButtonConversation;
    }

    if (selectedTypeDataset === "handwriting") {
      disabledButtonProccess = disableButtonHandwriting;
    }
  });

  let fileList = $state([]);
  let datasetArray = $state([]);
  let selectedFile = $state(null);
  let conversationId = $state(null);
  let editingMessageId = $state(null);
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
  let datas = $state([]);

  let sidebarData = $state({});
  let userData = $state({ username: "", email: "", brand: "" });

  let winInnerWidth = $derived(innerWidth);

  const dummySendMessage = async () => {
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
  };

  const startEditingMessage = (id, message) => {
    editingMessageId = id;
    editedMessage = message;
  };

  const saveEditedMessage = async () => {
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
  };

  const cancelEditing = () => {
    editingMessageId = null;
    editedMessage = "";
  };

  const clearConversation = async () => {
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
  };

  const newChat = async () => {
    datasetArray = [];
    currentChatTitle = "New Chat";
    conversationId = null;
    selectedRole = "User";
    const newUrl =
      window.location.protocol + "//" + window.location.host + window.location.pathname;

    window.history.pushState({ path: newUrl }, "", newUrl);
    setSelectedConversation(null);

    dataConversationChats = [];
  };

  const handleKeyDown = (event) => {
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
  };

  const onclickSend = () => {
    if (messageInput.trim() !== "") {
      // sendMessage();
      dummySendMessage();
    }
  };

  const autosize = () => {
    textarea.addEventListener("input", (e) => {
      setTimeout(() => {
        textarea.style.height = "auto";
        let scHeight = e.target.scrollHeight + 2; // Tambahkan sedikit untuk padding
        textarea.style.height = `${scHeight}px`;
      }, 0);
    });
  };

  const autoscroll = async () => {
    await tick();
    if (div) {
      div.scrollTo({ top: div.scrollHeight, behavior: "smooth" });
    }
  };

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

  const getmessage = async (conversationId) => {
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
  };

  const getPrevMessage = async () => {
    pageData++;
    noAutoScroll = true;
    processGetData = true;
    getmessage(conversationId);
  };

  const updateTitle = async () => {
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
  };

  const getListConversations = async () => {
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

      if (response.result.data.length < 1) {
        disableButtonConversation = true;
      } else {
        disableButtonConversation = false;
      }

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
  };

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

  const readFileAsBase64 = async (file) => {
    return new Promise((resolve, reject) => {
      const reader = new FileReader();

      reader.onload = () => {
        const base64Data = reader.result.split(",")[1];
        resolve(base64Data);
      };
      reader.onerror = reject;

      reader.readAsDataURL(file);
    });
  };

  const importConversationWithCsv = async () => {
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
  };

  const handleDropdownFocusLoss = ({ relatedTarget, currentTarget }) => {
    if (relatedTarget instanceof HTMLElement && currentTarget.contains(relatedTarget)) return;
    dropdownOpen = false;
  };

  const autoResizeTextarea = (event) => {
    const textarea = event.target;
    textarea.style.height = "auto"; // Reset height
    textarea.style.height = textarea.scrollHeight + "px"; // Adjust height
  };

  let paginationHandwriting = $state({
    page: 1,
    pageCount: 1,
    pageSize: 15,
    total: 10,
  });

  const getDatasetsHandwriting = async (search) => {
    let token = localStorage.getItem("token");
    if (token) {
      if (search !== "") {
        paginationHandwriting.page = 1;
      }

      const response = await $call.call("client-dashboard", "get-datasets-handwriting", {
        URLQuery: {},
        URLParams: {},
        Params: {
          page: paginationHandwriting.page,
          page_size: paginationHandwriting.pageSize,
          search: search,
        },
        Headers: {
          Authorization: "Bearer " + token,
        },
      });

      if (response.status !== 200) {
        triggerToast("error");
        return;
      }

      if (response.result.data.length < 1) {
        disableButtonHandwriting = true;
      } else {
        disableButtonHandwriting = false;
      }

      let promise = response.result.data.map(async (data) => {
        let link = await $call.getPage("client-dashboard", "datasets-detail");

        return {
          id: data.id,
          filename: data.attributes.name,
          fileext: data.attributes.path.split(".").pop(),
          filesize: data.attributes.file_size,
          createdAt: format(new Date(data.attributes.createdAt), "MMM d, yyyy"),
          href: link + "?datasetId=" + data.id,
        };
      });

      datas = await Promise.all(promise);
    }
  };

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
    getDatasetsHandwriting();

    return () => {
      globalThis.removeEventListener("resize", resizeHandler);
    };
  });

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

  onDestroy(() => {
    window.removeEventListener("popstate", handleBeforePopState);
    window.removeEventListener("beforeunload", handleBeforeUnload);
  });
</script>

<div bind:this={div} id="chat-container" class="static bg-primary-light dark:bg-primary-dark">
  <div class="flex flex-col pt-4 px-4 md:px-8">
    <div class="flex flex-col gap-3 md:gap-4">
      <div
        class="flex sticky top-0 z-[1] py-4 bg-primary-light dark:bg-primary-dark w-full justify-between items-center"
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
          </div>
          <div class="flex relative items-center justify-center gap-2">
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
            {#if check("aiManagement.datasets:delete")}
              <div onfocusout={handleDropdownFocusLoss} class="flex items-center relative">
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
                </Button>
                <div
                  class={`z-50 ${dropdownOpen ? "absolute -bottom-16 right-0" : "hidden"} bg-white divide-y divide-gray-100 rounded-lg shadow w-max dark:bg-gray-700`}
                >
                  <ul class="py-2 text-sm text-gray-700 dark:text-gray-200">
                    <button
                      class="flex p-2"
                      onclick={() => {
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
                  <button onclick={() => startEditingMessage(qaPair.id, qaPair.message)}>
                    <EditIcon
                      svgClass="w-full h-full text-brand-secondary-light dark:text-brand-secondary-dark"
                    />
                  </button>
                {/if}
              </div>
              {#if editingMessageId === qaPair.id}
                <Textarea
                  bind:value={editedMessage}
                  on:input={autoResizeTextarea}
                  unWrappedClass="focus:border-fg-success-secondary-light dark:focus:border-fg-success-secondary-dark focus:ring-0"
                ></Textarea>
                <div class="flex gap-2 justify-end">
                  <Button on:click={saveEditedMessage} class="h-max">{$_("save")}</Button>
                  <Button color="light" on:click={cancelEditing} class="h-max"
                    >{$_("cancel")}</Button
                  >
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
    spaceSidebar={false}
  />

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
