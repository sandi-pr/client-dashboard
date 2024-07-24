<script>
  import { P, Button, ToolbarButton } from "flowbite-svelte";

  import UserIcon from "$lib/icon/user-01.svg";
  import UsersIcon from "$lib/icon/users-03.svg";
  import VerifiedIcon from "$lib/icon/verified.svg";
  import Sidebar from "$lib/components/Sidebar.svelte";
  import { onMount, onDestroy } from "svelte";
  import { CR3Call } from "hexacomponents-svelte";
  import { Pane } from "svelte-splitpanes";
  import { tick } from "svelte";
  import { goto } from "$app/navigation";
  import Markdown from "svelte-exmarkdown";
  import { gfmPlugin } from "svelte-exmarkdown/gfm";
  import { PaperPlaneOutline, TrashBinOutline, RefreshOutline } from "flowbite-svelte-icons";
  import FilePlus02Icon from "$lib/icon/file-plus-02.svelte";
  import { _ } from "svelte-i18n";
  import { triggerToast, selectedCollection } from "$lib";

  const plugins = [gfmPlugin()];

  const call = new CR3Call(globalThis.location.origin);

  let selectedCollectionOld = $state({});

  let statusRAG = $state("");

  let imageFileName = $state("");
  let imageUrl = $state(null);
  let imageFile = $state(null);

  let div;
  let textarea;

  let chat_id = $state("");
  let content = $state("");

  let chatId = $state("");
  let messageInput = $state("");
  let timer;

  let disableInput = $state(false);
  let disableButton = $state(false);

  let init = $state(false);

  let dataConversatonLogs = $state([
    {
      role: "assistant",
      message:
        "Hi! I’m your virtual assistant. This is Collections AI. Please select source document by clicking “Configure” above.",
    },
  ]);

  // $effect.pre(() => {
  //   if (!div) return;
  //   dataConversatonLogs;
  //   const autoscroll =
  //     div.offsetHeight + div.scrollTop > div.scrollHeight - 50;

  //   if (autoscroll) {
  //     tick().then(() => {
  //       div.scrollTo(0, div.scrollHeight);
  //     });
  //   }
  // });

  $effect(() => {
    // if (!div) return; // not yet mounted

    // reference `messages` so that this code re-runs whenever it changes
    dataConversatonLogs;

    // autoscroll when new messages are added
    if (div) {
      tick().then(() => {
        div.scrollTo(0, div.scrollHeight);
      });
    }
  });

  $effect(() => {
    if (imageUrl) {
      console.log("image url:", imageUrl);

      dataConversatonLogs = [
        ...dataConversatonLogs,
        {
          role: "user",
          message: `![image](${imageUrl})`,
        },
      ];

      imageUrl = "";
    }
  });

  // Function to adjust the height of the textarea
  function autosize() {
    textarea.addEventListener("keyup", (e) => {
      textarea.style.height = "auto";
      let scHeight = e.target.scrollHeight;
      textarea.style.height = `${scHeight}px`;
    });
  }

  $effect(() => {
    if (textarea) {
      autosize();
    }
  });

  $effect(() => {
    disableButton = statusRAG !== "Ready" || messageInput.length < 1;
  });

  // // Add event listener on mount
  // onMount(() => {
  //   if (textarea) {
  //     textarea.addEventListener('keydown', autosize);
  //   }
  //   return () => {
  //     if (textarea) {
  //       textarea.removeEventListener('keydown', autosize);
  //     }
  //   };
  // });

  async function sendMessage() {
    disableInput = true;

    let token = sessionStorage.getItem("token");
    if (messageInput.trim() !== "" && token) {
      dataConversatonLogs.push({
        role: "user",
        message: messageInput.trim(),
      });

      let newMessage = messageInput.trim();

      messageInput = "";

      let conversationIdCollection = sessionStorage.getItem(
        "conversationIdCollection-" + $selectedCollection.id.toString(),
      );

      try {
        let chat = await call.call("hexa-ai-rag", "chat", {
          URLQuery: {},
          URLParams: {
            collection_id: $selectedCollection.id.toString(),
            conversation_id: conversationIdCollection,
          },
          Params: {
            message: newMessage,
          },
          Headers: {
            Authorization: "Bearer " + token,
          },
        });

        console.log("chat", chat);

        if (chat.status === 200) {
          disableInput = false;
          dataConversatonLogs = [
            ...dataConversatonLogs,
            {
              role: chat.result.message.role,
              message: chat.result.message.message,
            },
          ];
        }
      } catch (error) {
        triggerToast("error");
        console.log(error);
      }

      disableInput = false;
    }
    disableInput = false;
  }

  async function getMessage() {
    if (dataConversatonLogs.length < 2) {
      dataConversatonLogs.push({
        role: "assistant",
        message: "Collection is selected. Please wait while we prepare your assistant.",
      });
    }

    let conversationIdCollection = sessionStorage.getItem(
      "conversationIdCollection-" + $selectedCollection.id.toString(),
    );

    let token = sessionStorage.getItem("token");

    let response2 = await call.call("hexa-ai-rag", "add-document-by-collection-id", {
      URLQuery: {},
      URLParams: {},
      Params: {
        collection_id: $selectedCollection.id.toString(),
      },
      Headers: {
        Authorization: "Bearer " + token,
      },
    });

    if (response2.status !== 200) {
      console.log(response2);
      return;
    }

    console.log(response2);

    if (conversationIdCollection) {
      try {
        let chat = await call.call("hexa-ai-rag", "get-messages", {
          URLQuery: {},
          URLParams: {
            collection_id: $selectedCollection.id.toString(),
            conversation_id: conversationIdCollection,
          },
          Params: {},
          Headers: {
            Authorization: "Bearer " + token,
          },
        });

        console.log("chat", chat);

        if (chat.status === 200) {
          if (chat.result.messages.length > 0) {
            dataConversatonLogs = chat.result.messages;
          }

          getStatusRag();
        }
      } catch (error) {
        triggerToast("error");
        console.log(error);
      }
    } else {
      try {
        let createConversationResponse = await call.call("hexa-ai-rag", "create-conversation", {
          URLQuery: {},
          URLParams: {},
          Params: {
            collection_id: parseInt($selectedCollection.id.toString()),
          },
          Headers: {
            Authorization: "Bearer " + token,
          },
        });

        console.log("createConversationResponse", createConversationResponse);

        if (createConversationResponse.status === 200) {
          conversationIdCollection = createConversationResponse.result.conversation_id;
          sessionStorage.setItem(
            "conversationIdCollection-" + $selectedCollection.id.toString(),
            conversationIdCollection,
          );

          getStatusRag();
        } else {
          triggerToast("error");
          console.log(createConversationResponse);
        }
      } catch (error) {
        triggerToast("error");
        console.log(error);
      }
    }
  }

  async function getStatusRag() {
    let token = sessionStorage.getItem("token");
    try {
      // let rag = await call.call("hexa-ai-rag", "get-rag-status", {
      //   URLQuery: {},
      //   URLParams: {},
      //   Params: {},
      //   Headers: {
      //     Authorization: "Bearer " + token,
      //   },
      // });

      // let status = rag.result.rag_status;
      // if (statusRAG.includes("!")) {
      //   let status2 = status.split("!")[0].trim();
      //   if (
      //     status2.includes("(NOVERB)") ||
      //     status2.includes("(INOVERB)") ||
      //     status2.includes("(!INOVERB)")
      //   ) {
      //     let status3 = status2
      //       .replace("(NOVERB)", "")
      //       .replace("(INOVERB)", "")
      //       .replace("(!INOVERB)", "")
      //       .replace("(!NOVERB)", "");
      //     statusRAG = status3;
      //   } else {
      //     statusRAG = status2;
      //   }
      // } else {
      //   statusRAG = status;
      // }
      // console.log(rag);
      let allData = [];
      let page = 1;
      while (true) {
        let response = await call.call("client-dashboard", "get-documents-by-collectionid", {
          URLQuery: {},
          URLParams: {},
          Params: {
            page: page,
            page_size: 25,
            collection_id: selectedCollectionOld.id.toString(),
          },
          Headers: {
            Authorization: "Bearer " + token,
          },
        });

        if (response.status !== 200) {
          console.log(response);
          return;
        }

        if (response.result.data.length > 0) {
          allData = allData.concat(response.result.data);
          page++;
        } else {
          break;
        }
      }

      if (allData.length === 0) {
        dataConversatonLogs.push({
          role: "assistant",
          message: "Your collection is empty",
        });
        return;
      }

      let allReady = true;

      for (let i = 0; i < allData.length; i++) {
        if (allData[i].attributes.status !== "Ready") {
          allReady = false;
          break;
        }
      }

      if (allReady) {
        statusRAG = "Ready";
        if (dataConversatonLogs[dataConversatonLogs.length - 1].message.includes("Processing...")) {
          dataConversatonLogs[dataConversatonLogs.length - 1].message = "Your document is ready";
        } else {
          if (
            dataConversatonLogs[dataConversatonLogs.length - 1].message.includes(
              "Your document is ready",
            )
          ) {
            console.log("Already");
            clearTimeout(timer);
            return;
          } else {
            dataConversatonLogs.push({
              role: "assistant",
              message: "Your document is ready",
            });
          }
        }
      } else {
        let totalDocument = allData.length;
        let totalProcessed = 0;
        for (let i = 0; i < allData.length; i++) {
          if (allData[i].attributes.status === "Ready") {
            totalProcessed++;
          }
        }
        let percentage = (totalProcessed / totalDocument) * 100;
        statusRAG = percentage.toFixed(1) + "%";

        if (dataConversatonLogs[dataConversatonLogs.length - 1].message.includes("Processing...")) {
          dataConversatonLogs[dataConversatonLogs.length - 1].message =
            "Processing... " + statusRAG;
        } else if (
          dataConversatonLogs[dataConversatonLogs.length - 1].message.includes(
            "Your document is ready",
          )
        ) {
          dataConversatonLogs[dataConversatonLogs.length - 1].message =
            "Processing... " + statusRAG;
        } else {
          dataConversatonLogs.push({
            role: "assistant",
            message: "Processing... " + statusRAG,
          });
        }
      }
    } catch (error) {
      console.error("Error:", error);
    }
    timer = setTimeout(getStatusRag, 2000);
  }

  const handleFileUploadWrapper = () => {
    const inputElement = document.createElement("input");
    inputElement.type = "file";
    inputElement.accept = ".jpg, .jpeg, .png";

    inputElement.addEventListener("change", async (event) => {
      const file = event.target.files[0];

      if (file) {
        imageUrl = URL.createObjectURL(file);
        imageFile = file;

        imageFileName = file.name;
      } else {
        alert("Invalid file. Please select a valid file.");
      }
    });

    inputElement.click();
  };

  function clearChat() {
    sessionStorage.removeItem("conversationId");
    window.location.reload();
  }

  // $effect(() => {
  //   if (selectedCollection !== null && statusRAG !== "Ready") {
  //   }
  // });

  $effect(() => {
    if ($selectedCollection) {
      if ($selectedCollection !== selectedCollectionOld) {
        if ($selectedCollection !== null) {
          selectedCollectionOld = $selectedCollection;
          console.log("$selectedCollection");
          // getMessage();
        }
      }
    }
  });

  onMount(() => {
    // getMessage();
  });

  onDestroy(() => {
    if (timer) {
      clearTimeout(timer);
    }
  });
</script>

<div class="static h-screen w-full bg-primary-light dark:bg-primary-dark">
  <div class="flex h-full flex-col px-4 pb-[210px] pt-4">
    <div bind:this={div} class="flex flex-col gap-4 overflow-y-auto">
      <!-- <P
        size="xl"
        weight="bold"
        class="text-fg-brand-primary-light dark:text-fg-brand-primary-dark leading-loose"
        >Chat 24.12.2023
      </P> -->
      {#each dataConversatonLogs as qaPair}
        <div
          class=" bg-secondary-light flex dark:bg-secondary-dark rounded-2xl p-4 items-start gap-4"
        >
          <div class="flex flex-col w-full gap-2">
            <div class="flex items-center space-x-2 rtl:space-x-reverse">
              <div
                class="w-8 h-8 relative flex justify-center items-center bg-tertiary-light dark:bg-tertiary-dark rounded-full border border-primary-light dark:border-primary-dark"
              >
                <img src={UserIcon} alt="User" class="w-6 h-6" />
                {#if qaPair.role === "assistant"}
                  <div class="absolute bottom-0 right-0">
                    <img src={VerifiedIcon} alt="User" class="w-3 h-3" />
                  </div>
                {/if}
              </div>
              <P weight="bold" class="text-primary-light dark:text-primary-dark"
                >{qaPair.role === "user" ? "You" : "HexaAI"}</P
              >
            </div>
            <div
              class="prose font-medium text-tertiary-light dark:text-primary-dark [&>p>a]:text-primary-light dark:[&>p>a]:text-primary-dark"
            >
              <Markdown md={qaPair.message} />
            </div>
          </div>
        </div>
      {/each}
      {#if disableInput}
        <div
          class=" bg-secondary-light flex dark:bg-secondary-dark rounded-2xl p-4 items-start gap-4"
        >
          <div class="flex flex-col w-full gap-2">
            <div class="flex items-center space-x-2 rtl:space-x-reverse">
              <div
                class="w-8 h-8 relative flex justify-center items-center bg-tertiary-light dark:bg-tertiary-dark rounded-full border border-primary-light dark:border-primary-dark"
              >
                <img src={UserIcon} alt="User" class="w-6 h-6" />
                <div class="absolute bottom-0 right-0">
                  <img src={VerifiedIcon} alt="User" class="w-3 h-3" />
                </div>
              </div>
              <P weight="bold" class="text-primary-light dark:text-primary-dark">{"HexaAI"}</P>
            </div>
            <P weight="medium" class="text-tertiary-light dark:text-tertiary-dark"
              >{"Assistant is typing..."}</P
            >
          </div>
        </div>
      {/if}
    </div>
  </div>
  <div class="sticky bottom-0 right-0 bg-white dark:bg-primary-dark px-2 py-2 flex items-center">
    <form on:submit|preventDefault={sendMessage} class="flex-grow">
      <!-- <Input
      bind:disabled={disableInput}
      bind:value={messageInput}
      placeholder="Send a message..."
      class=" text-wrap min-h-12 pe-24 rounded-2xl"
    >
      <Button
        slot="right"
        size="sm"
        type="submit"
        class="my-3"
        disabled={disableInput}>Send</Button
      >
    </Input> -->
      <label for="chat" class="sr-only">Your message</label>
      <div class="flex items-center py-2 rounded-lg">
        <!-- <Button
          on:click={() => handleFileUploadWrapper()}
          color="light"
          class="rounded-full !px-2.5"
        >
          <FilePlus02Icon
            svgClass="w-5 h-5 text-fg-brand-secondary-light dark:text-fg-brand-secondary-dark"
          />
        </Button> -->
        <!-- <Textarea
          id="chat"
          bind:disabled={disableInput}
          bind:value={messageInput}
          class="mx-4"
          rows="1"
          placeholder="Send a message..."
        /> -->
        <textarea
          bind:this={textarea}
          bind:value={messageInput}
          readonly={disableInput}
          id="chat"
          rows="1"
          class="block mx-2 p-2.5 w-full text-sm text-brand-primary-light dark:text-primary-dark bg-white dark:bg-primary-dark rounded-lg border border-utility-brand-500-light dark:border-utility-brand-500-dark focus:border-utility-brand-500-light dark:focus:utility-brand-500-dark focus:ring-utility-brand-400-light dark:focus:ring-utility-brand-400-light dark:placeholder-gray-400"
          placeholder="Send a message..."
        ></textarea>
        <Button type="submit" disabled={disableButton} color="primary" class="rounded-full !px-2.5">
          <PaperPlaneOutline class="w-5 h-5 rotate-45" />
          <!-- <span class="sr-only">Send</span> -->
        </Button>
      </div>
    </form>
  </div>
</div>

<style>
  textarea {
    /* overflow-y: hidden; */
    max-height: 220px;
  }
  textarea::-webkit-scrollbar {
    display: none;
  }
</style>
