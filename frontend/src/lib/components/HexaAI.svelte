<script>
  import { Button, Input, P } from "flowbite-svelte";
  import UserIcon from "$lib/icon/user-01.svg";
  import UsersIcon from "$lib/icon/users-03.svg";
  import VerifiedIcon from "$lib/icon/verified.svg";
  import Sidebar from "$lib/components/Sidebar.svelte";
  import { onMount } from "svelte";
  import { CR3Call } from "hexacomponents-svelte";

  const call = new CR3Call(globalThis.location.origin);

  let chat_id = $state("");
  let content = $state("");

  let chatId = $state("");
  let messageInput = $state("");

  let disableInput = $state(false);

  let dataConversatonLogs = $state([
    {
      role: "user",
      message:
        "Hi there, I'm currently in [Location] and I'm in need of some services. Specifically, I'm looking for a reliable wood flooring company. Do you have any recommendations?",
    },
    {
      role: "assistant",
      message:
        'Of course! London has several reputable wood flooring companies. One notable option is "Zway Flooring Solutions" located in the heart of the city. They are well-regarded for their premium materials and skilled craftsmanship. Would you like more information about them or any other suggestions?',
    },
  ]);

  // const ws = new WebSocket(
  //   "ws://localhost:3000/api-hexabot-webchat/c/9756f06b-a6ef-4bc6-8c9b-5fd6cddf64b6",
  // );

  const ws = new WebSocket("wss://cr3-api.hexacode.org/api-hexabot-webchat/c/123");

  function SendMessage(message) {
    ws.send(JSON.stringify(message));
  }

  function Joining(session_id) {
    SendMessage({
      type: "joining",
      session_id: session_id,
      joining: {
        session_id: session_id,
      },
    });
  }

  function Leaving(session_id) {
    SendMessage({
      type: "leaving",
      session_id: session_id,
    });
  }

  function Message(session_id, message) {
    SendMessage({
      type: "content",
      session_id: session_id,
      content: {
        role: "user",
        message: message,
        timestamp: new Date().valueOf(),
        time_zone: timezonse(),
      },
    });
  }

  ws.onopen = function (event) {
    console.log("WebSocket connected");
  };

  ws.onerror = function (error) {
    console.error("WebSocket error:", error);
  };

  ws.onmessage = function (event) {
    console.log(JSON.parse(event.data));
    data = JSON.parse(event.data);
    if (data.type == "content") {
      // content.append(CreateMessage(data.content));
      dataConversatonLogs = [...dataConversatonLogs, { role: "user", message: data.content }];
    }

    if (data.type == "accept") {
      // chat_id.val(data.accept.session_id);
      chatId = data.accept.session_id;
      for (var i = 0; i < data.accept.contents.length; i++) {
        // content.append(CreateMessage(data.accept.contents[i]));
        dataConversatonLogs = [
          ...dataConversatonLogs,
          { role: "assistant", message: data.accept.contents[i] },
        ];
      }
    }
  };

  async function sendMessage() {
    disableInput = true;
    let token = sessionStorage.getItem("token");
    if (messageInput.trim() !== "") {
      dataConversatonLogs = [
        ...dataConversatonLogs,
        { role: "user", message: messageInput.trim() },
      ];
      try {
        let conversationId = await call.call("hexa-ai", "create-conversation", {
          URLQuery: {},
          URLParams: {},
          Params: {},
          Headers: {
            Authorization: "Bearer " + token,
          },
        });

        if (conversationId.status === 200) {
          let chat = await call.call("hexa-ai", "chat", {
            URLQuery: {},
            URLParams: {
              conversation_id: conversationId.result.conversation_id,
            },
            Params: {
              message: messageInput.trim(),
            },
            Headers: {
              Authorization: "Bearer " + token,
            },
          });

          console.log("merpati", chat);

          if (chat.status === 0) {
            dataConversatonLogs = [
              ...dataConversatonLogs,
              { role: chat.result.message.role, message: chat.result.message.message },
            ];

            disableInput = false;
          } else {
            disableInput = false;
          }
        } else {
          console.log(conversationId);

          disableInput = false;
        }
        // Message(chatId, messageInput.trim());
        messageInput = "";
        disableInput = false;
      } catch (error) {
        console.log(error);

        disableInput = false;
      }
    }
  }

  function newChat() {
    dataConversatonLogs = [];
    Joining("");
  }

  onMount(() => {
    // Joining("");
  });
</script>

<div class="relative h-screen overflow-auto bg-primary-light dark:bg-primary-dark">
  <div class="flex flex-col px-4 pt-6 pb-24">
    <div class="flex flex-col gap-4">
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
          <div
            class="w-12 h-11 relative flex justify-center items-center bg-tertiary-light dark:bg-tertiary-dark rounded-full border border-primary-light dark:border-primary-dark"
          >
            <img src={UserIcon} alt="User" class="" />
            {#if qaPair.role === "assistant"}
              <div class="absolute bottom-0 right-0">
                <img src={VerifiedIcon} alt="User" class="w-4 h-4" />
              </div>
            {/if}
          </div>
          <div class="flex flex-col w-full gap-1">
            <div class="flex items-center space-x-2 rtl:space-x-reverse">
              <P weight="bold" class="text-primary-light dark:text-primary-dark"
                >{qaPair.role === "user" ? "You" : "Assistant"}</P
              >
            </div>
            <P weight="medium" class="text-tertiary-light dark:text-tertiary-dark"
              >{qaPair.message}</P
            >
          </div>
        </div>
      {/each}
    </div>
  </div>
  <div class="fixed bottom-4">
    <form on:submit|preventDefault={sendMessage} class="px-4 [&>div>div]:!pe-4">
      <Input
        bind:disable={disableInput}
        bind:value={messageInput}
        placeholder="Send a message..."
        class=" text-wrap min-h-16 pe-24 rounded-2xl"
      >
        <Button slot="right" size="sm" type="submit" class="my-3">Send</Button>
      </Input>
    </form>
  </div>
</div>
