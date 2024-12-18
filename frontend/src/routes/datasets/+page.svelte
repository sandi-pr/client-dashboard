<script>
  import {
    P,
    Span,
    Modal,
    Input,
    Button,
    Toggle,
    Spinner,
    Heading,
    Textarea,
    DarkMode,
    Search,
  } from "flowbite-svelte";
  import { _ } from "svelte-i18n";
  import { page } from "$app/stores";
  import PdfViewer from "svelte-pdf";
  import { goto } from "$app/navigation";
  import { fly } from "svelte/transition";
  import Markdown from "svelte-exmarkdown";
  import { format, parse, sub } from "date-fns";
  import { CloseButton } from "flowbite-svelte";
  import { onMount, onDestroy, tick } from "svelte";
  import { gfmPlugin } from "svelte-exmarkdown/gfm";
  import { SearchOutline, ExclamationCircleSolid } from "flowbite-svelte-icons";

  import { CR3Call } from "hexacomponents-svelte";
  import {
    call,
    check,
    Typography,
    triggerToast,
    triggerShowAI,
    triggerLoading,
    selectedConversation,
    setSelectedConversation,
  } from "$lib";
  import Chat from "$lib/components/Chat.svelte";
  import StateCard from "$lib/ui/cards/StateCard.svelte";
  import Sidebar from "$lib/components/SidebarDataset.svelte";
  import DataTableLibraryCollectionBakrie from "$lib/components/DataTableLibraryCollectionBakrie.svelte";

  import EyeIcon from "$lib/icon/eye.svg";
  import UserIcon from "$lib/icon/user-01.svg";
  import MenuIcon from "$lib/icon/menu-01.svg";
  import VerifiedIcon from "$lib/icon/verified.svg";
  import PDFSolidIcon from "$lib/icon/PDFSolid.svg";
  import DOCSolidIcon from "$lib/icon/DOCSolid.svg";
  import JPGSolidIcon from "$lib/icon/JPGSolid.svg";
  import PNGSolidIcon from "$lib/icon/PNGSolid.svg";
  import DOCXSolidIcon from "$lib/icon/DOCXSolid.svg";
  import JPEGSolidIcon from "$lib/icon/JPEGSolid.svg";
  import ArrowLeftIcon from "$lib/icon/arrow-left.svg";
  import ArrowRightIcon from "$lib/icon/arrow-right.svg";
  import UploadIcon from "$lib/icon/upload-cloud-02.svg";
  import AlertCircleIcon from "$lib/icon/alert-circle.svg";
  import DownloadIcon from "$lib/icon/file-download-03.svg";

  import EditIcon from "$lib/icon/edit-05.svelte";
  import FileTypeIcon from "$lib/icon/file-type-icon.svelte";
  import FilePlus02Icon from "$lib/icon/file-plus-02.svelte";
  import AlignRightIcon from "$lib/icon/align-right-01.svelte";
  import DotsVerticalIcon from "$lib/icon/dots-vertical.svelte";
  import ToastContainer from "$lib/ui/toast/ToastContainer.svelte";
  import LaodingLottie from "$lib/components/LaodingLottie.svelte";
  import IllustrationEmpty from "$lib/icon/illustration_empty.svelte";
  import IllustrationProgress from "$lib/icon/illustration_progress.svelte";

  const plugins = [gfmPlugin()];
  const isFrontend = process.env.NODE_ENV === "development";

  let div;
  let ref;
  let timer;
  let textarea;

  let initRAG = $state(false);
  let isLoading = $state(false);
  let editTitle = $state(false);
  let hideSidebar = $state(true);
  let initDocList = $state(true);
  let deleteModal = $state(false);
  let gettingDatas = $state(true);
  let pagePreview = $state(false);
  let previewModal = $state(false);
  let viewDocument = $state(false);
  let dropdownOpen = $state(false);
  let noAutoScroll = $state(false);
  let fileListReady = $state(false);
  let processGetData = $state(false);
  let showSendMessage = $state(false);
  let addDatasetModal = $state(false);
  // let fileDetailModal = $state(true);
  let pageListDocument = $state(false);
  let documentListModal = $state(false);
  let disableInput = $state(!check("aiManagement.datasets:create"));

  let pageData = $state(1);
  let pageTotal = $state(1);
  let innerWidth = $state(0);
  let indexFileList = $state(0);

  let base64Data = $state("");
  let fileTextName = $state("");
  let messageInput = $state("");
  let editedMessage = $state("");
  let selectedRole = $state("User");
  let statusRAG = $state("Not Ready");
  let selectedTypeDataset = $state("handwriting");
  let typeDataset = $state(["handwriting"]);
  // let typeDataset = $state(["handwriting", "conversation"]);
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

  let files = $state([]);
  let allData = $state([]);
  let fileList = $state([]);
  let tmpFile = $state(null);
  let allFileList = $state([]);
  let datasetArray = $state([]);
  let conversations = $state([]);
  let selectedFile = $state(null);
  let conversationId = $state(null);
  let tmpDocumentName = $state(null);
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

  let pagination = $state({
    page: 1,
    pageCount: 1,
    pageSize: 10,
    total: 10,
  });
  let sidebarData = $state({});
  let tmpResultTrocr = $state({ text: "" });
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

  const addToDataset = (role, message) => {
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
  };

  const saveDatasetEntry = (entry) => {
    console.log("Dataset Entry:", entry);
    datasetArray.push(entry);
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

  const logout = async () => {
    sessionStorage.clear();
    localStorage.clear();

    try {
      isFrontend ? "" : await $call.goToPage("client-login", "home-page");
    } catch (error) {
      console.log(error);
    }
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

  const handleChangeHome = () => {
    setHome(true);
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

  const addToDatasetFromDB = () => {
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

  const readFileAsText = async (file) => {
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

  let next = async () => {
    // getCategories(pagination.page + 1);
  };

  let previous = async () => {
    // getCategories(pagination.page - 1);
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

  const handleDatasetModal = () => {
    if (!addDatasetModal) filesHandwriting = [];
    addDatasetModal = !addDatasetModal;
  };

  const dropHandle = async (event) => {
    event.preventDefault();

    const files = event.dataTransfer.files;

    for (const file of files) {
      if (file) {
        let base64Data;
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
          base64Data: base64Data,
        };

        filesHandwriting = [...filesHandwriting, fileData];
      } else {
        alert("Invalid file. Please select a valid file.");
      }
    }
  };

  let filesHandwriting = $state([]);
  let multipleUpload = $state(true);
  let dataSend = $state({});
  let paginationHandwriting = $state({
    page: 1,
    pageCount: 1,
    pageSize: 15,
    total: 10,
  });
  let tmpId = $state(null);
  let popupModal = $state(false);

  const handleUploadWrapper2 = () => {
    const inputElement = document.createElement("input");
    inputElement.type = "file";
    inputElement.accept = ".png,.jpg,.jpeg";
    inputElement.multiple = multipleUpload;

    inputElement.addEventListener("change", async (event) => {
      let filesSelect = event.target.files;

      if (!multipleUpload) {
        filesHandwriting = [];
        filesSelect = [filesSelect[0]];
      }

      const fileTypeRegex = new RegExp(
        inputElement.accept
          .split(",")
          .map((type) => type.trim().replace(".", "\\."))
          .join("|") + "$",
        "i",
      );

      for (const file of filesSelect) {
        if (file && fileTypeRegex.test(file.name)) {
          const base64Data = await readFileAsBase64(file);

          filesHandwriting.push({
            id: filesHandwriting.length < 1 ? 1 : filesHandwriting.length + 1,
            name: file.name,
            size: formatBytes(file.size),
            base64: base64Data,
            ext: "." + file.name.split(".").pop(),
            type: "Handwriting",
          });
        }
      }
    });
    inputElement.click();
  };

  const dropHandle2 = async (event) => {
    event.preventDefault();

    const accept = ".png,.jpg,.jpeg";

    const files = event.dataTransfer.files;

    const fileTypeRegex = new RegExp(
      accept
        .split(",")
        .map((type) => type.trim().replace(".", "\\."))
        .join("|") + "$",
      "i",
    );

    for (const file of files) {
      if (file && fileTypeRegex.test(file.name)) {
        const base64Data = await readFileAsBase64(file);

        filesHandwriting.push({
          id: filesHandwriting.length < 1 ? 1 : filesHandwriting.length + 1,
          name: file.name,
          size: formatBytes(file.size),
          base64: base64Data,
          ext: "." + file.name.split(".").pop(),
          type: "Handwriting",
        });
      }
    }
  };

  function formatBytes(bytes, decimals = 2) {
    if (bytes === 0) return "0 Bytes";

    const k = 1024;
    const dm = decimals < 0 ? 0 : decimals;
    const sizes = ["Bytes", "KB", "MB", "GB", "TB", "PB", "EB", "ZB", "YB"];

    const i = Math.floor(Math.log(bytes) / Math.log(k));

    return parseFloat((bytes / Math.pow(k, i)).toFixed(dm)) + " " + sizes[i];
  }

  async function addDatasetHandwriting() {
    let token = localStorage.getItem("token");

    if (token !== null) {
      try {
        addDatasetModal = false;
        triggerLoading(true);

        for (let i = 0; i < filesHandwriting.length; i++) {
          let response3 = await $call.call("hexa-ai-rag", "trocr-predict", {
            URLQuery: {},
            URLParams: {},
            Params: {
              image: filesHandwriting[i].base64,
              model: sessionStorage.getItem("model_trocr"),
            },
            Headers: {
              Authorization: "Bearer " + token,
            },
          });

          if (response3.status !== 200) {
            addDatasetModal = true;
            triggerLoading(false);
            if (response3.status === 400) {
              console.error("Invalid parameters:", response3);
              triggerToast($_("invalidParameters"), true, "error");
            } else if (response3.status === 401) {
              console.error("Unauthorized access:", response3);
              triggerToast($_("unauthorized"), true, "error");
            } else if (response3.status === 500) {
              console.error("Server error:", response3);
              triggerToast($_("serverError"), true, "error");
            } else {
              console.error("Unexpected status code:", response3);
              triggerToast(`${$_("unexpectedError")}: ${response3.status}`, true, "error");
            }
            return;
          }

          if (response3?.result?.text !== null) {
            let originalName = filesHandwriting[i].name.split(".").slice(0, -1).join(".");
            let newName = originalName;

            filesHandwriting[i].name = newName + filesHandwriting[i].ext;

            dataSend = {
              name: originalName,
              documents: JSON.stringify([filesHandwriting[i]]),
              detected_text: response3.result.text,
              result_trocr: JSON.stringify(response3.result),
            };

            if (filesHandwriting.length < 2) {
              triggerLoading(false);
            }
          }

          if (filesHandwriting.length > 1) {
            await sendDocument();
          } else {
            previewModal = true;
            selectedFile = filesHandwriting[0];
          }
        }

        addDatasetModal = false;

        if (filesHandwriting.length > 1) {
          getDatasetsHandwriting();
          filesHandwriting = [];
          dataSend = null;
          previewModal = false;
          triggerLoading(false);
          triggerToast($_("succesfullyUploaded"), true);
        }
      } catch (error) {
        addDatasetModal = true;
        console.error("Error during document upload:", error);
        triggerToast($_("failedUploadDocument"), true, "error");
        triggerLoading(false);
        return;
      }
    } else {
      triggerToast($_("unauthorized"), true, "error");
      return;
    }
  }

  async function sendDocument() {
    let token = localStorage.getItem("token");

    if (token === null) {
      triggerToast($_("unauthorized"), true, "error");
      return;
    }

    triggerLoading(true);
    let response = await $call.call("client-dashboard", "create-dataset-handwriting", {
      URLQuery: {},
      URLParams: {},
      Params: dataSend,
      Headers: {
        Authorization: "Bearer " + token,
      },
    });

    if (response.status !== 200) {
      addDatasetModal = true;
      triggerLoading(false);
      if (response.status === 400) {
        console.error("Invalid parameters:", response);
        triggerToast($_("invalidParameters"), true, "error");
      } else if (response.status === 401) {
        console.error("Unauthorized access:", response);
        triggerToast($_("unauthorized"), true, "error");
      } else if (response.status === 500) {
        console.error("Server error:", response);
        triggerToast($_("serverError"), true, "error");
      } else {
        console.error("Unexpected status code:", response);
        triggerToast(`${$_("unexpectedError")}: ${response.status}`, true, "error");
      }

      return;
    }

    if (filesHandwriting.length === 1) {
      getDatasetsHandwriting();
      triggerLoading(false);
      filesHandwriting = [];
      dataSend = null;
      previewModal = false;
      triggerLoading(false);
      triggerToast($_("succesfullyUploaded"), true);
    }

    return;
  }

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

  function handleDeleteDataset(id) {
    tmpId = id;
    popupModal = true;
  }

  async function deleteDataset() {
    popupModal = false;
    triggerLoading(true);
    let token = localStorage.getItem("token");

    if (token === null) {
      triggerToast($_("unauthorized"), true, "error");
      return;
    }

    try {
      const response = await $call.call("client-dashboard", "delete-dataset-handwriting", {
        URLQuery: {},
        URLParams: {
          id: tmpId.toString(),
        },
        Params: {},
        Headers: {
          Authorization: "Bearer " + token,
        },
      });

      if (response.status !== 200) {
        console.log(response);
        triggerLoading(false);
        popupModal = true;
        if (response.status === 400) {
          console.error("Invalid parameters:", response);
          triggerToast($_("invalidParameters"), true, "error");
        } else if (response.status === 401) {
          console.error("Unauthorized access:", response);
          triggerToast($_("unauthorized"), true, "error");
        } else if (response.status === 500) {
          console.error("Server error:", response);
          triggerToast($_("serverError"), true, "error");
        } else {
          console.error("Unexpected status code:", response);
          triggerToast(`${$_("unexpectedError")}: ${response.status}`, true, "error");
        }
        return;
      }

      tmpId = null;

      triggerLoading(false);
      triggerToast("successDeleteDocument", true);
      getDatasetsHandwriting();
    } catch (error) {
      console.error("Error during document deletion:", error);
      triggerToast($_("failedDeleteDocument"), true, "error");
    }
  }

  function handleSearchDataset(event) {
    getDatasetsHandwriting(event.target.value);
  }

  async function handleDownloadDataset() {
    let token = localStorage.getItem("token");
    if (token) {
      const response = await $call.call("client-dashboard", "download-dataset-handwriting-zip", {
        URLQuery: {},
        URLParams: {},
        Params: {},
        Headers: {
          Authorization: "Bearer " + token,
        },
      });

      if (response.status !== 200 || response.result.base64 === null) {
        triggerToast("error");
        return;
      }

      const now = new Date();
      const formattedDate = format(now, "ddMMyyyy");

      const filename = `datasets-handwriting-${formattedDate}.zip`;

      const binaryString = atob(response.result.base64);

      const byteArray = new Uint8Array(binaryString.length);
      for (let i = 0; i < binaryString.length; i++) {
        byteArray[i] = binaryString.charCodeAt(i);
      }

      const blob = new Blob([byteArray], { type: "application/zip" });

      const url = URL.createObjectURL(blob);

      const link = document.createElement("a");
      link.href = url;
      link.download = filename;

      document.body.appendChild(link);
      link.click();

      document.body.removeChild(link);

      URL.revokeObjectURL(url);
    }
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
    getDatasetsHandwriting();

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
    {typeDataset}
    bind:disabledButtonProccess
    bind:selectedTypeDataset
    {handleDatasetModal}
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
            <div class={`flex justify-between ${innerWidth < 768 ? "flex-col" : "items-center"}`}>
              <Typography
                preset="Display sm/Semibold"
                color="primary"
                text={$_("allDatasetsHandwriting")}
              />
              <Button
                disabled={disableButtonHandwriting}
                onclick={handleDownloadDataset}
                icon={DownloadIcon}
                class={innerWidth < 768 ? "mt-4" : ""}
              >
                <img class="w-6 h-6 me-1" src={DownloadIcon} alt="plus-icon" />
                {$_("downloadAllAsZip")}
              </Button>
            </div>

            <div class="mt-8 mb-4">
              <Input
                id="search"
                type="search"
                class="bg-white dark:bg-hc-graydark-950 items-center"
                placeholder={$_("search")}
                on:input={handleSearchDataset}
              >
                <SearchOutline slot="left" class="w-5 h-5 text-gray-500 dark:text-gray-400" />
              </Input>
            </div>

            <!-- Person Table Data -->
            {#if datas.length > 0}
              <DataTableLibraryCollectionBakrie
                functDeleteCategories={handleDeleteDataset}
                collectionListDatas={datas}
                otherFunc={() => {
                  hideSidebar = true;
                }}
                {pagination}
                {next}
                {previous}
              />
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
      {#if filesHandwriting[0].ext.includes(".png") || filesHandwriting[0].ext.includes(".jpg") || filesHandwriting[0].ext.includes(".jpeg")}
        <img
          src={`data:${"Handwriting"};base64,${filesHandwriting[0].base64}`}
          style="max-width: 80%; height: auto;"
          alt="PNG Preview"
        />
      {:else if filesHandwriting[0].ext.includes(".pdf")}
        <PdfViewer
          showBorder={false}
          showButtons={false}
          scale={1}
          data={atob(filesHandwriting[0].base64)}
        />
      {/if}
    </div>
    <!-- detected handwriting -->
    {#if filesHandwriting[0].type === "Handwriting"}
      <div class="flex w-full flex-col gap-3 px-4 pb-1">
        <Typography preset="Text md/Semibold" text="Handwriting Information" color="secondary" />
        <div class="flex flex-col gap-1.5">
          <Typography preset="Text sm/Medium" text="Detected Text" color="secondary" />
          <textarea
            bind:value={dataSend.detected_text}
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
          filesHandwriting = [];
        }}
        class="w-full">Cancel</Button
      >
      <Button color="primary" on:click={sendDocument} class="rounded-lg w-full">I’m sure</Button>
    </div>
  </svelte:fragment>
</Modal>

<Modal bind:open={addDatasetModal} size="lg">
  <svelte:fragment slot="header">
    <div class="flex justify-start items-center gap-4 py-4 px-6">
      <div class="p-2 bg-[#F4EBFF] rounded-full flex justify-center items-center">
        <FilePlus02Icon svgClass="w-8 h-8 text-[#7F56D9]" />
      </div>
      <div class="flex flex-col">
        <Typography preset="Text lg/Semibold" text="Add Datasets" color="primary" />
        <Typography
          preset="Text sm/Regular"
          text="It can be anything. PDF, Images, etc!"
          color="tertiary"
        />
      </div>
    </div>
  </svelte:fragment>
  <div class="flex flex-col w-full justify-center">
    <Button
      ondrop={dropHandle}
      ondragover={(event) => {
        event.preventDefault();
      }}
      onclick={() => {
        handleUploadWrapper2();
      }}
      color="light"
      class={`w-full rounded-xl border border-primary-light dark:border-primary-dark bg-primary-light dark:bg-primary-dark gap-1 py-4 md:py-8 flex md:flex-col items-center justify-start ${fileTextName ? "items-center cursor-default" : "items-center"}`}
    >
      {#if filesHandwriting.length > 0}
        {#if !multipleUpload}
          <div class="flex flex-col items-center">
            <FileTypeIcon
              svgClass="w-20 h-20 text-utility-brand-600-light dark:text-utility-brand-600-dark"
            />
            <Typography
              preset="Text sm/Semibold"
              text={filesHandwriting[0].name}
              color="fg-brand-primary"
              propClass="fluid-base md:fluid-md"
            />
          </div>
        {:else}
          <div class="flex items-center">
            <FileTypeIcon
              svgClass="w-20 h-20 text-utility-brand-600-light dark:text-utility-brand-600-dark"
            />
            <div
              class="flex items-center justify-center bg-[#7F56D9] rounded-full w-8 h-8 ml-[-1.25rem]"
            >
              <P class="text-white font-semibold">{filesHandwriting.length}</P>
            </div>
          </div>
        {/if}
      {:else}
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
            text="JPG, JPEG, PNG (Max. 15MB)"
            color="tertiary"
            propClass="fluid-xs md:fluid-base"
          />
        </div>
      {/if}
    </Button>

    <!-- <Toggle
      class="mt-4"
      on:change={(e) => {
        console.log("toggled");
        multipleUpload = !multipleUpload;
        if (filesHandwriting.length > 1) {
          filesHandwriting = [filesHandwriting[0]];
        }
        // if (e.detail) {
        //   bulkUpload = true;
        // } else {
        //   bulkUpload = false;
        // }
      }}
    >
      <P>Bulk Upload</P>
    </Toggle> -->
  </div>
  <svelte:fragment slot="footer">
    <div class="flex w-full gap-2 justify-center items-center py-3 px-4">
      <Button
        color="primary"
        on:click={addDatasetHandwriting}
        class="rounded-lg w-full"
        disabled={filesHandwriting.length === 0}>Add</Button
      >
    </div>
  </svelte:fragment>
</Modal>

<Modal bind:open={popupModal} size="xs" autoclose>
  <div class="text-center">
    <ExclamationCircleSolid class="mx-auto mb-4 text-red-600 w-12 h-12" />
    <h3 class="mb-5 text-lg font-normal text-gray-500 dark:text-gray-400">
      {$_("areYouSureYouWantToDelete")}
      {$_("thisDocument")}
    </h3>
    <Button onclick={deleteDataset} color="red" class="me-2">{$_("yesImSure")}</Button>
    <Button color="alternative">{$_("noCancel")}</Button>
  </div>
</Modal>
