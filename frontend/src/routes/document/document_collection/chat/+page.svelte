<script>
  import { P, Input, Button } from "flowbite-svelte";
  import UsersIcon from "$lib/icon/users-01.svg";
  import UserIcon from "$lib/icon/user-03.svg";
  import VerifiedIcon from "$lib/icon/verified.svg";
  import { onMount } from "svelte";
  import { tick } from "svelte";
  import { CR3Call } from "hexacomponents-svelte";
  import { goto } from "$app/navigation";
  import { triggerToast } from "$lib";
  import LoadingSpinner from "$lib/components/LoadingSpinner.svelte";

  const call = new CR3Call(globalThis.location.origin);

  const params = new URLSearchParams(window.location.search);
  let collectionId = params.get("collectionId");

  let disableInput = $state(false);

  let div;

  let messageInput = $state("");

  let statusRAG = $state("");

  let isLoading = $state(true);

  let dataConversatonLogs = $state([
    {
      role: "assistant",
      message: "What would you like to know in this collection?",
    },
  ]);

  async function sendMessage() {
    disableInput = true;

    let token = sessionStorage.getItem("token");
    if (messageInput.trim() !== "" && token) {
      dataConversatonLogs = [
        ...dataConversatonLogs,
        { role: "user", message: messageInput.trim() },
      ];

      let newMessage = messageInput.trim();

      messageInput = "";

      let conversationIdCollection = sessionStorage.getItem(
        "conversationIdCollection",
      );

      try {
        let chat = await call.call("hexa-ai-rag", "chat", {
          URLQuery: {},
          URLParams: {
            collection_id: collectionId,
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

  async function getMessages() {
    let conversationIdCollection = sessionStorage.getItem(
      "conversationIdCollection",
    );

    let token = sessionStorage.getItem("token");

    if (conversationIdCollection && token) {
      try {
        let chat = await call.call("hexa-ai-rag", "get-messages", {
          URLQuery: {},
          URLParams: {
            collection_id: collectionId,
            conversation_id: conversationIdCollection,
          },
          Params: {},
          Headers: {
            Authorization: "Bearer " + token,
          },
        });

        console.log("chat", chat);

        if (chat.status === 200) {
          dataConversatonLogs = chat.result.messages;
          if (dataConversatonLogs.length > 0) {
            dataConversatonLogs = [
              {
                role: "assistant",
                message: "What would you like to know in this collection?",
              },
            ];
          }

          isLoading = false;
        }
      } catch (error) {
        triggerToast("error");
        console.log(error);
      }
    } else {
      try {
        let createConversationResponse = await call.call(
          "hexa-ai-rag",
          "create-conversation",
          {
            URLQuery: {},
            URLParams: {},
            Params: {
              collection_id: parseInt(collectionId),
            },
            Headers: {
              Authorization: "Bearer " + token,
            },
          },
        );

        console.log("createConversationResponse", createConversationResponse);

        if (createConversationResponse.status === 200) {
          conversationIdCollection =
            createConversationResponse.result.conversation_id;
          sessionStorage.setItem(
            "conversationIdCollection",
            conversationIdCollection,
          );

          isLoading = false;
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
      let rag = await call.call("hexa-ai-rag", "get-rag-status", {
        URLQuery: {},
        URLParams: {},
        Params: {},
        Headers: {
          Authorization: "Bearer " + token,
        },
      });

      let status = rag.result.rag_status;
      if (statusRAG.includes("!")) {
        let status2 = status.split("!")[0].trim();
        if (
          status2.includes("(NOVERB)") ||
          status2.includes("(INOVERB)") ||
          status2.includes("(!INOVERB)")
        ) {
          let status3 = status2
            .replace("(NOVERB)", "")
            .replace("(INOVERB)", "")
            .replace("(!INOVERB)", "")
            .replace("(!NOVERB)", "");
          statusRAG = status3;
        } else {
          statusRAG = status2;
        }
      } else {
        statusRAG = status;
      }
      console.log(rag);
      if (statusRAG !== "Ready") {
        triggerToast("error");
        goto(await call.getPage("client-dashboard", "document-collection"));
      }
    } catch (error) {
      console.error("Error:", error);
    }
  }

  onMount(async () => {
    getStatusRag();
    getMessages();
  });
</script>

{#if isLoading}
  <LoadingSpinner />
{:else}
  <div class="static h-full flex flex-col">
    <!-- <div class="mx-4 py-6 border-b border-secondary-light dark:border-secondary-dark">
    <P size="3xl" weight="semibold" class="text-primary-light dark:text-primary-dark leading-[38px]"
      >Conversation Logs 4c64986428240cn948929284
    </P>
  </div> -->
    <div class="flex h-screen flex-col px-4 pt-2">
      <div bind:this={div} class="flex flex-col gap-4">
        {#each dataConversatonLogs as qaPair}
          <div
            class="flex bg-secondary-light dark:bg-secondary-dark rounded-2xl px-8 py-6 mx-4 my-6 items-start gap-4"
          >
            <div
              class="w-12 h-12 relative flex justify-center items-center bg-tertiary-light dark:bg-tertiary-dark rounded-full border border-primary-light dark:border-primary-dark"
            >
              {#if qaPair.role === "user"}
                <img src={UsersIcon} alt="User" class="" />
              {/if}
              {#if qaPair.role === "assistant"}
                <img src={UserIcon} alt="User" class="" />
                <div class="absolute bottom-0 right-0">
                  <img
                    src={VerifiedIcon}
                    alt="User"
                    style={{ width: 16, height: 16 }}
                  />
                </div>
              {/if}
            </div>
            <div class="flex flex-col w-full gap-1">
              <div class="flex items-center space-x-2 rtl:space-x-reverse">
                <P
                  weight="bold"
                  class="text-primary-light dark:text-primary-dark"
                  >{qaPair.role === "user" ? "Customer" : "Assistant"}</P
                >
              </div>
              <P
                weight="medium"
                class="text-tertiary-light dark:text-tertiary-dark"
                >{qaPair.message}</P
              >
            </div>
          </div>
        {/each}
        {#if disableInput}
          <div
            class="flex bg-secondary-light dark:bg-secondary-dark rounded-2xl px-8 py-6 mx-4 my-6 items-start gap-4"
          >
            <div
              class="w-12 h-12 relative flex justify-center items-center bg-tertiary-light dark:bg-tertiary-dark rounded-full border border-primary-light dark:border-primary-dark"
            >
              <img src={UserIcon} alt="User" class="" />
              <div class="absolute bottom-0 right-0">
                <img
                  src={VerifiedIcon}
                  alt="User"
                  style={{ width: 16, height: 16 }}
                />
              </div>
            </div>
            <div class="flex flex-col w-full gap-1">
              <div class="flex items-center space-x-2 rtl:space-x-reverse">
                <P
                  weight="bold"
                  class="text-primary-light dark:text-primary-dark"
                  >{"Assistant"}</P
                >
              </div>
              <P
                weight="medium"
                class="text-tertiary-light dark:text-tertiary-dark"
                >{"Assistant is typing..."}</P
              >
            </div>
          </div>
        {/if}
      </div>
    </div>

    <div class="sticky bg-white dark:bg-primary-dark pb-4 bottom-0">
      <form
        on:submit|preventDefault={sendMessage}
        class="px-8 [&>div>div]:!pe-6"
      >
        <Input
          disabled={disableInput}
          bind:value={messageInput}
          placeholder="Send a message..."
          class=" text-wrap min-h-16 pe-24 rounded-2xl focus:border-utility-brand-500-light focus:ring-utility-brand-500-light dark:focus:border-utility-brand-500-dark dark:focus:ring-utility-brand-500-dark"
        >
          <Button
            disabled={disableInput}
            slot="right"
            size="sm"
            type="submit"
            class="my-3 bg-fg-brand-primary-light dark:bg-fg-brand-primary-dark"
            >Send</Button
          >
          <!-- <Button
          disabled={disableInput}
          color="red"
          size="sm"
          onclick={clear}
          class="my-3 bg-blue-800 dark:bg-blue-700">Clear</Button
        > -->
        </Input>
      </form>
      <!-- <div
      class="relative min-h-16 mx-8 py-3 border rounded-2xl border-brand-solid-light dark:border-brand-solid-dark"
    >
      <Button class="absolute right-4 lg:right-6 top-1/2 transform -translate-y-1/2 ">Send</Button>
    </div> -->
    </div>
  </div>
{/if}
