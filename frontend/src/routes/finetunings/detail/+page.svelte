<script>
  import { P, Button, Modal } from "flowbite-svelte";
  import { _ } from "svelte-i18n";
  import { page } from "$app/stores";
  import { goto } from "$app/navigation";
  import { format } from "date-fns";
  import { onMount, onDestroy, tick } from "svelte";
  import { ChevronLeftOutline, ExclamationCircleSolid } from "flowbite-svelte-icons";

  import {
    call,
    check,
    Typography,
    TextareaRw,
    triggerToast,
    triggerShowAI,
    triggerLoading,
    selectedConversation,
    setSelectedConversation,
  } from "$lib";
  import Divider from "$lib/components/Divider.svelte";
  import StateCard from "$lib/ui/cards/StateCard.svelte";
  import PdfjsViewer from "$lib/components/PdfjsViewer.svelte";

  import ToastContainer from "$lib/ui/toast/ToastContainer.svelte";
  import LaodingLottie from "$lib/components/LaodingLottie.svelte";
  import IllustrationEmpty from "$lib/icon/illustration_empty.svelte";
  import IllustrationProgress from "$lib/icon/illustration_progress.svelte";

  const isFrontend = process.env.NODE_ENV === "development";

  let div;
  let textarea;

  let pdfReady = $state(false);
  let isLoading = $state(false);
  let hideSidebar = $state(true);
  let gettingDatas = $state(true);
  let noAutoScroll = $state(false);
  let processGetData = $state(false);
  let addFinetuneModal = $state(false);

  let pageData = $state(1);
  let pageTotal = $state(1);
  let innerWidth = $state(0);

  let base64 = $state(null);
  let contentFile = $state(null);
  let selectedRole = $state("User");
  let currentChatTitle = $state("New Chat");

  let finetuneArray = $state([]);
  let conversationId = $state(null);
  let dataConversationChats = $state([]);

  let fileUrl = $state("null");
  let data = $state({});
  let tmpData = $state({});
  let popupModal = $state(false);
  let onEdit = $state(false);

  let sidebarData = $state({});
  let userData = $state({ username: "", email: "", brand: "" });
  let fileData = $state({});

  const newChat = async () => {
    finetuneArray = [];
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
    triggerLoading(true);
    let token = localStorage.getItem("token");
    if (token !== null) {
      try {
        let response = await $call.call("client-dashboard", "get-dataset-handwriting", {
          URLQuery: {},
          URLParams: {
            id: conversationId.toString(),
          },
          Params: {},
          Headers: {
            Authorization: "Bearer " + token,
          },
        });

        // console.log("document", response);

        if (response.status !== 200) {
          console.log(response);
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
          goto(await $call.getPage("client-dashboard", "finetunings"));
          return;
        }

        let documentData = response.result.data;

        // collectionId = response.result.data.parent_id;

        base64 = response.result.base64;
        pdfReady = true;

        contentFile = atob(base64);

        let fileType = "";
        if (response.result.ext === ".pdf") {
          fileType = "application/pdf";
        } else if (response.result.ext === ".docx") {
          fileType = "application/vnd.openxmlformats-officedocument.wordprocessingml.document";
        } else if (response.result.ext === ".doc") {
          fileType = "application/msword";
        } else {
          fileType = "text/plain";
        }

        const blob = base64toBlob(base64, fileType);
        fileUrl = URL.createObjectURL(blob);
        console.log(fileUrl);

        tmpData = {
          id: response.result.id,
          createdAt: documentData.createdAt,
          name: documentData.name,
          document_type: "Handwriting",
          ext: response.result.ext,
          // document_type: "Document Audio",
          // ext: ".mp3",
          description: documentData.detected_text ? documentData.detected_text : "",
          slices: documentData.result_trocr ? documentData.result_trocr.slices : "",
          result_trocr: documentData.result_trocr,
          href: fileUrl,
        };

        fileData = {
          id: response.result.id,
          createdAt: documentData.createdAt,
          name: documentData.name,
          document_type: "Handwriting",
          ext: response.result.ext,
          // document_type: "Document Audio",
          // ext: ".mp3",
          description: documentData.detected_text ? documentData.detected_text : "",
          slices: documentData.result_trocr ? documentData.result_trocr.slices : "",
          result_trocr: documentData.result_trocr,
          href: fileUrl,
        };

        triggerLoading(false);
      } catch (error) {
        triggerLoading(false);
        console.error("Error during document get:", error);
        triggerToast($_("failedGetDocument"), true, "error");
      }
    } else {
      triggerToast($_("unauthorized"), true, "error");
      return;
    }
  };

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

  // maybe this?
  const getDocument = async () => {
    triggerLoading(true);
    let token = localStorage.getItem("token");
    if (token !== null) {
      try {
        let response = await $call.call("client-dashboard", "get-document", {
          URLQuery: {},
          URLParams: {
            id: documentId.toString(),
          },
          Params: {},
          Headers: {
            Authorization: "Bearer " + token,
          },
        });

        if (response.status !== 200) {
          triggerLoading(false);
          triggerToast("error");
          console.log(response);
          if (collectionId !== null) {
            goto(
              (await $call.getPage("client-dashboard", "document-collection")) +
                "?collectionId=" +
                collectionId,
            );
            return;
          } else {
            goto(await $call.getPage("client-dashboard", "document-collection"));
            return;
          }
        }

        if (response.result.document.data.length > 0) {
          let documentData = response.result.document.data[0];

          base64 = response.result.base64;
          pdfReady = true;

          contentFile = atob(base64);

          let fileType = "";
          if (data.ext === ".pdf") {
            fileType = "application/pdf";
          } else if (data.ext === ".docx") {
            fileType = "application/vnd.openxmlformats-officedocument.wordprocessingml.document";
          } else if (data.ext === ".doc") {
            fileType = "application/msword";
          } else {
            fileType = "text/plain";
          }

          const blob = base64toBlob(base64, fileType);
          fileUrl = URL.createObjectURL(blob);

          // categories = await getCategories(collectionId);

          fileData = {
            id: documentData.id,
            createdAt: documentData.attributes.createdAt,
            document_title: documentData.attributes.document_title,
            description: documentData.attributes.description,
            category: documentData.attributes.category,
            document_type: "Handwriting",
            ext: documentData.attributes.document.data.attributes.ext,
            detectedText: documentData.attributes.detected_text
              ? documentData.attributes.detected_text
              : "",
            slices: documentData.attributes.result_trocr
              ? documentData.attributes.result_trocr.slices
              : "",
            result_trocr: documentData.attributes.result_trocr,
          };

          tmpData = {
            id: documentData.id,
            createdAt: documentData.attributes.createdAt,
            document_title: documentData.attributes.document_title,
            description: documentData.attributes.description,
            category: documentData.attributes.category,
            document_type: "Handwriting",
            ext: documentData.attributes.document.data.attributes.ext,
            detectedText: documentData.attributes.detected_text
              ? documentData.attributes.detected_text
              : "",
            slices: documentData.attributes.result_trocr
              ? documentData.attributes.result_trocr.slices
              : "",
            result_trocr: documentData.attributes.result_trocr,
          };

          tags = [documentData.attributes.category];

          console.log("tags", documentData.attributes.category);

          fileReady = true;
          triggerLoading(false);
        }
      } catch (error) {
        triggerLoading(false);
        triggerToast("error");
        console.log(error);
      }
    } else {
      // triggerToast("error");
    }
  };

  function handleDeleteFinetune(id) {
    popupModal = true;
  }

  async function deleteFinetune() {
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
          id: fileData.id.toString(),
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

      goto(await $call.call("client-dashboard", "finetunings"));

      triggerLoading(false);
      triggerToast($_("successDeleteDocument"), true);
    } catch (error) {
      console.error("Error during document deletion:", error);
      triggerToast($_("failedDeleteDocument"), true, "error");
    }
  }

  async function updateFinetune() {
    let token = localStorage.getItem("token");
    try {
      const response = await $call.call("client-dashboard", "update-dataset-handwriting", {
        URLQuery: {},
        URLParams: {
          id: fileData.id.toString(),
        },
        Params: {
          detected_text: fileData.description,
        },
        Headers: {
          Authorization: "Bearer " + token,
        },
      });

      if (response.status !== 200) {
        console.log(response);
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

      onEdit = false;
      triggerToast("updated");
    } catch (e) {
      console.log(e);
      return;
    }
  }

  onMount(async () => {
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

    conversationId = $page.url.searchParams.get("datasetId");

    if (conversationId) {
      getmessage(conversationId);
    } else {
      goto(await $call.getPage("client-dashboard", "finetunings"));
    }

    getListConversations();

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

<div
  bind:this={div}
  id="chat-container"
  class="static h-screen overflow-y-auto bg-primary-light dark:bg-primary-dark"
>
  <div class="flex flex-col pt-4 px-4 md:px-8 dark:bg-primary-dark">
    <div class="w-full pb-10">
      <div class="flex flex-col m-4 pt-2">
        <!-- <Navbar /> -->

        <div class="flex flex-col">
          <button
            class="ml-4 flex items-center gap-1 mb-2"
            onclick={async () => {
              if (isFrontend) {
                goto("/finetunings/handwriting");
              } else {
                goto(await $call.getPage("client-dashboard", "finetunings-handwriting"));
              }
            }}
          >
            <ChevronLeftOutline class="w-7 h-7 text-fg-tertiary-light dark:text-fg-tertiary-dark" />
            <P class="font-semibold text-fg-tertiary-light dark:text-fg-tertiary-dark">
              {$_("back")}
            </P>
          </button>

          <div
            class={`flex justify-between mb-4 ${innerWidth < 768 ? "flex-col" : "items-center"}`}
          >
            <Typography
              preset="Display sm/Semibold"
              color="primary"
              text={fileData.name + fileData.ext}
            />
            <div class={`flex gap-2 ${innerWidth < 768 ? "mt-4" : ""}`}>
              {#if !onEdit}
                <Button
                  onclick={handleDeleteFinetune}
                  outline
                  color="red"
                  class={innerWidth < 768 ? "w-full" : ""}
                >
                  {$_("delete")}
                </Button>
                <Button
                  onclick={() => {
                    onEdit = true;
                  }}
                  outline
                  class={innerWidth < 768 ? "w-full" : ""}
                >
                  {$_("edit")}
                </Button>
              {:else}
                <Button
                  onclick={() => {
                    fileData = tmpData;
                    onEdit = false;
                  }}
                  outline
                  color="alternative"
                  class={innerWidth < 768 ? "w-full" : ""}
                >
                  {$_("cancel")}
                </Button>
                <Button
                  onclick={() => {
                    updateFinetune();
                  }}
                  outline
                  class={innerWidth < 768 ? "w-full" : ""}
                >
                  {$_("save")}
                </Button>
              {/if}
            </div>
          </div>

          <Divider />

          <!-- Person Table Data -->
          {#if Object.keys(fileData).length > 0}
            <div class="flex flex-col md:flex-row gap-4 md:gap-6 my-8">
              <div class="w-60">
                <P class="font-semibold text-secondary-light dark:text-secondary-dark">
                  {$_("handwritingInformation")}
                </P>
                <!-- <button
                  class="font-semibold text-[#6941C6] mt-2 text-left"
                  on:click={() => {
                    onEdit = true;
                  }}
                >
                  {$_("correctHandwriting")}
                </button> -->
              </div>
              <div class="flex flex-col grow gap-4 md:gap-6">
                <div class="grid grid-cols-1">
                  <div class="flex flex-col gap-2">
                    <P
                      class="text-sm font-medium leading-tight text-secondary-light dark:text-secondary-dark"
                    >
                      {$_("detectedText")}
                    </P>
                    <TextareaRw {onEdit} bind:value={fileData.description} />
                  </div>
                </div>
              </div>
            </div>

            <Divider />

            <div class="flex flex-col md:flex-row gap-6 mt-8">
              <div class="w-60">
                <P class="font-semibold text-secondary-light dark:text-secondary-dark">
                  {$_("viewer")}
                </P>
              </div>
              <div
                class="flex flex-col relative items-center grow py-4 md:py-12 rounded-lg bg-secondary-light dark:bg-secondary-dark"
              >
                {#if fileData.ext === ".pdf"}
                  {#if pdfReady}
                    <PdfjsViewer pdfFile={base64} />
                  {/if}
                {:else if fileData.ext === ".jpg" || fileData.ext === ".png" || fileData.ext === ".jpeg" || fileData.ext === ".webp"}
                  <div class="flex justify-center">
                    <img class="w-full h-screen" src={`${fileData.href}`} alt="imageViewer" />
                  </div>
                {:else}
                  {#each contentFile.split("\n") as paragraphs, i}
                    <p class="dark:text-white">{paragraphs}</p>
                  {/each}
                {/if}
              </div>
            </div>
          {:else if gettingDatas}
            <StateCard
              illustration={IllustrationProgress}
              description={$_("fetchingDataPleaseWait")}
            />
          {:else}
            <StateCard illustration={IllustrationEmpty} description={$_("youDontHaveAnyDataYet")} />
          {/if}
        </div>
      </div>
    </div>
  </div>

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

<Modal bind:open={popupModal} size="xs" autoclose>
  <div class="text-center">
    <ExclamationCircleSolid class="mx-auto mb-4 text-red-600 w-12 h-12" />
    <h3 class="mb-5 text-lg font-normal text-gray-500 dark:text-gray-400">
      {$_("areYouSureYouWantToDelete")}
      {$_("thisDocument")}
    </h3>
    <Button onclick={deleteFinetune} color="red" class="me-2">{$_("yesImSure")}</Button>
    <Button color="alternative">{$_("noCancel")}</Button>
  </div>
</Modal>
