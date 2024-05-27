<script>
  import { Button, Input, P } from "flowbite-svelte";
  import UserIcon from "$lib/icon/user-01.svg";
  import UsersIcon from "$lib/icon/users-03.svg";
  import VerifiedIcon from "$lib/icon/verified.svg";
  import Sidebar from "$lib/components/Sidebar.svelte";
  import { onMount } from "svelte";
  import { CR3Call } from "hexacomponents-svelte";
  import { Pane } from "svelte-splitpanes";
  import { tick } from "svelte";
  import Markdown from "svelte-exmarkdown";
  import { gfmPlugin } from "svelte-exmarkdown/gfm";

  const plugins = [gfmPlugin()];

  const call = new CR3Call(globalThis.location.origin);

  let div;

  let chat_id = $state("");
  let content = $state("");

  let chatId = $state("");
  let messageInput = $state("");

  let disableInput = $state(false);

  let dataConversatonLogs = $state([
    {
      role: "assistant",
      message: "Hello, I'm your Hexa AI assistant. How can I help you today?",
    },
    {
      role: "assistant",
      message: "Hello, I'm your Hexa AI assistant. How can I help you today?",
    },
    {
      role: "assistant",
      message: "Hello, I'm your Hexa AI assistant. How can I help you today?",
    },
    {
      role: "assistant",
      message: "Hello, I'm your Hexa AI assistant. How can I help you today?",
    }, {
      role: "assistant",
      message: "Hello, I'm your Hexa AI assistant. How can I help you today?",
    }, {
      role: "assistant",
      message: "Hello, I'm your Hexa AI assistant. How can I help you today?",
    }, {
      role: "assistant",
      message: "Hello, I'm your Hexa AI assistant. How can I help you today?",
    }, {
      role: "assistant",
      message: "Hello, I'm your Hexa AI assistant. How can I help you today?",
    }, {
      role: "assistant",
      message: "Hello, I'm your Hexa AI assistant. How can I help you today?",
    },
  ]);

  $effect.pre(() => {
    dataConversatonLogs;
    const autoscroll =
      div && div.offsetHeight + div.scrollTop > div.scrollHeight - 50;

    if (autoscroll) {
      tick().then(() => {
        div.scrollTo(0, div.scrollHeight);
      });
    }
  });

  async function sendMessage() {
    disableInput = true;
    let token = sessionStorage.getItem("token");

    if (messageInput.trim() !== "") {
      dataConversatonLogs = [
        ...dataConversatonLogs,
        { role: "user", message: messageInput.trim() },
      ];

      let newMessage = messageInput.trim();

      messageInput = "";

      let conversationId = localStorage.getItem("conversationId");

      if (!conversationId) {
        try {
          let createConversationResponse = await call.call(
            "hexa-ai",
            "create-conversation",
            {
              URLQuery: {},
              URLParams: {},
              Params: {},
              Headers: {
                Authorization: "Bearer " + token,
              },
            },
          );

          if (createConversationResponse.status === 200) {
            conversationId = createConversationResponse.result.conversation_id;
            localStorage.setItem("conversationId", conversationId);
          } else {
            console.log(createConversationResponse);
          }
        } catch (error) {
          console.log(error);
        }
      }

      try {
        let chat = await call.call("hexa-ai", "chat", {
          URLQuery: {},
          URLParams: {
            conversation_id: conversationId,
          },
          Params: {
            message: newMessage,
          },
          Headers: {
            Authorization: "Bearer " + token,
          },
        });

        console.log("merpati", chat);

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
        console.log(error);
      }

      disableInput = false;
    }
  }

  async function getMessage() {
    let token = sessionStorage.getItem("token");

    let conversationId = localStorage.getItem("conversationId");

    if (conversationId) {
      try {
        let chat = await call.call("hexa-ai", "get-messages", {
          URLQuery: {},
          URLParams: {
            conversation_id: conversationId,
          },
          Params: {},
          Headers: {
            Authorization: "Bearer " + token,
          },
        });

        console.log(chat);

        if (chat.status === 200) {
          dataConversatonLogs = [
            ...dataConversatonLogs,
            ...chat.result.messages,
          ];
        }
      } catch (error) {
        console.log(error);
      }
    }
  }

  onMount(() => {
    getMessage();
  });
</script>

<div
  bind:this={div}
  class="static h-screen w-full bg-primary-light dark:bg-primary-dark"
>
  <div class="flex h-full flex-col px-4 pb-[178px] pt-2">
    <div class="flex flex-col gap-4 overflow-y-auto">
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
              class="prose font-medium text-tertiary-light dark:text-primary-dark"
            >
              <Markdown md={qaPair.message} {plugins} />
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
              <P weight="bold" class="text-primary-light dark:text-primary-dark"
                >{"HexaAI"}</P
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
  <div class="sticky bottom-0 right-0 bg-white dark:bg-primary-dark px-4 py-2 flex items-center">
    <form on:submit|preventDefault={sendMessage} class="flex-grow">
      <Input
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
    </Input>
    </form>
  </div>
  
</div>
