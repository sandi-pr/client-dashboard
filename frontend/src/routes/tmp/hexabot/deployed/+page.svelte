<script>
  import {
    P,
    Card,
    Button,
    Heading,
    A,
    Dropdown,
    DropdownItem,
    Search,
    Toast,
  } from "flowbite-svelte";
  import { _ } from "svelte-i18n";
  import { toCamelCase } from "$lib/utils/toCamelCase.js";
  import { DotsVerticalOutline } from "flowbite-svelte-icons";
  import CopyIcon from "$lib/icon/copy-03.svelte";
  import MessageChatSquare from "$lib/icon/message-chat-square.svelte";
  import PencilLine from "$lib/icon/pencil-line.svelte";
  import DotsVerticalIcon from "$lib/icon/dots-vertical.svelte";
  import FileAttachment_02 from "$lib/icon/file-attachment-02.svelte";

  let searchTerm = $state("");
  let bots = [
    {
      botsName: "Veronika",
      botsID: "5904",
      clientID: "c2287dc7-0af7-4481-9ec3-dbaa910477e7",
      activationTime: "2023-12-13 05:00:56",
      urlBot: "https://webchat.aws.hexacode.org/25dfb75d-3985-40c5-8c7b-1ece3f50699a",
      timestamp: "2023-12-13 05:00:56",
      deployStatus: "Done",
      deployedPlatform: "Web",
    },
    {
      botsName: "Woody Woodison",
      botsID: "5871",
      clientID: "96602bf7-28d4-4d56-a720-49dfc36b0e38",
      activationTime: "2024-01-02 09:56:53",
      urlBot: "https://webchat.aws.hexacode.org/9756f06b-a6ef-4bc6-8c9b-5fd6cddf64b6",
      timestamp: "2024-01-02 09:56:53",
      deployStatus: "Done",
      deployedPlatform: "Web",
    },
    {
      botsName: "JobsBots",
      botsID: "5919",
      clientID: "b608d45f-0f54-45f7-816e-8d6aa82e6eba",
      activationTime: "2024-02-06 10:30:13",
      urlBot: "https://webchat.aws.hexacode.org/1fff9299-6be6-4bb5-b5ec-0ff1a9695642",
      timestamp: "2024-02-06 10:30:13",
      deployStatus: "Done",
      deployedPlatform: "WhatsApp",
    },
    {
      botsName: "Gia",
      botsID: "5917",
      clientID: "b4257652-ae68-4e29-a527-0dcff052a2de",
      activationTime: "2024-01-22 04:20:50",
      urlBot: "https://webchat.aws.hexacode.org/d4ad076f-71b3-4297-a8f2-40039b3c00e0",
      timestamp: "2024-01-22 04:20:50",
      deployStatus: "On Progress",
    },
    {
      botsName: "JCS Plumbing and Heating Solutions",
      botsID: "5925",
      clientID: "6a159bd2-f3b9-4a46-b0f3-7f944095fc1f",
      activationTime: "2024-03-08 10:20:30",
      urlBot: "https://webchat.aws.hexacode.org/427b7d23-777c-43fe-af1f-f69532080285",
      timestamp: "2024-03-08 10:20:30",
      deployStatus: "On Progress",
    },
    {
      botsName: "Wealth Bank",
      botsID: "21411",
      clientID: "null",
      activationTime: "2024-02-06 10:30:13",
      urlBot: "https://webchat.aws.hexacode.org/e7b05ed9-66a2-4a82-ba31-93931f851923",
      timestamp: "2024-02-06 10:30:13",
      deployStatus: "Not Ready",
    },
  ];

  // $: filteredItems = bots.filter(
  //   (bot) => bot.botsName.toLowerCase().indexOf(searchTerm.toLowerCase()) !== -1,
  // );

  // let filteredItems = $state([]);

  // $effect(() => {
  //   filteredItems = bots.filter((bot) => {
  //     for (const key in bot) {
  //       if (
  //         typeof bot[key] === "string" &&
  //         bot[key].toLowerCase().indexOf(searchTerm.toLowerCase()) !== -1
  //       ) {
  //         return true;
  //       }
  //     }
  //     return false;
  //   });
  // });

  function copyToClipboard(text) {
    navigator.clipboard.writeText(text);
    trigger();
  }

  let open = false;
  let counter = 6;

  function trigger() {
    open = true;
    counter = 6;
    timeout();
  }

  function timeout() {
    if (--counter > 0) return setTimeout(timeout, 1000);
    open = false;
  }

  const classStatusType = (type) => {
    if (type === "Done") {
      return "bg-hc-success-50 border-hc-success-200 text-hc-success-700 dark:bg-hc-success-900 dark:border-hc-success-700 dark:text-hc-success-300";
    } else if (type === "Not Ready") {
      return "bg-hc-error-50 border-hc-error-200 text-hc-error-700 dark:bg-hc-error-900 dark:border-hc-error-700 dark:text-hc-error-300";
    } else if (type === "On Progress") {
      return "bg-hc-warning-50 border-hc-warning-200 text-hc-warning-700 dark:bg-hc-warning-900 dark:border-hc-warning-700 dark:text-hc-warning-300";
    } else if (type === "WhatsApp") {
      return "text-utility-success-700-light dark:text-utility-success-700-dark bg-utility-success-50-light dark:bg-utility-success-50-dark border-utility-success-200-light dark:border-utility-success-200-dark";
    } else if (type === "Web") {
      return "text-sky-700 dark:text-sky-300 bg-sky-50 dark:bg-sky-950 border-sky-200 dark:border-sky-800";
    }
  };
</script>

<div class="flex dark:bg-hc-graydark-950">
  <div class="w-full pb-10">
    <div class="flex flex-col m-4">
      <!-- <Navbar /> -->

      <div class="flex flex-col gap-8">
        <!-- User Information -->
        <Heading tag="h2" class="font-jakarta text-hc-graylight-900 dark:text-hc-graydark-50"
          >{$_("deployedBotsList")}</Heading
        >

        <!-- Bots List Data -->
        <div>
          <!-- <Search
            class="bg-white dark:bg-hc-graydark-950 mb-6 rounded-lg"
            placeholder={$_("searchBotName")}
            bind:value={searchTerm}
          /> -->
          <div class="grid grid-cols-1 md:grid-cols-2 2xl:grid-cols-3 gap-4">
            {#each bots as bot (bot.botsID)}
              {#if bot.deployStatus === "Done"}
                <div class="span-col-1">
                  <Card padding="none" size="md">
                    <div
                      class="bg-active-light dark:bg-active-dark border-b border-secondary-light dark:border-secondary-dark px-4 py-3"
                    >
                      <div class="flex justify-between items-center">
                        <div class="flex flex-col">
                          <div class="flex items-center gap-2">
                            <Heading
                              tag="h6"
                              class="!font-semibold text-primary-light dark:text-primary-dark"
                              >{`${bot.botsName}`}</Heading
                            >
                            <div
                              class={`px-2 py-0.5 rounded-full border justify-start items-center inline-flex text-center text-xs font-medium leading-[18px] 
                              ${classStatusType(bot.deployedPlatform)}`}
                            >
                              {bot.deployedPlatform}
                            </div>
                          </div>
                          <P size="sm" class="text-secondary-light dark:text-secondary-dark"
                            >{`ID: ${bot.botsID}`}</P
                          >
                        </div>
                        <div class="flex items-center pl-2 gap-2">
                          <A
                            href="/hexabot/detail"
                            class="text-fg-brand-primary-light dark:text-fg-brand-primary-dark text-sm font-semibold hover:no-underline"
                            >{$_("view")}</A
                          >
                          <div class=" flex items-center">
                            <Button class="focus-within:ring-0 !p-0" color="none" size="xs"
                              ><DotsVerticalIcon
                                class="text-fg-quinary-light dark:text-fg-quinary-dark"
                              /></Button
                            >
                            <!-- <Dropdown>
                              <DropdownItem>
                                <div class="flex">
                                  <MessageChatSquare
                                    svgClass="text-fg-quaternary-light dark:text-fg-quaternary-dark mr-2"
                                  />{$_("conversationLogs")}
                                </div>
                              </DropdownItem>
                              <DropdownItem>
                                <div class="flex">
                                  <PencilLine
                                    svgClass="text-fg-quaternary-light dark:text-fg-quaternary-dark mr-2"
                                  />{$_("editBotName")}
                                </div>
                              </DropdownItem>
                              <DropdownItem>
                                <div class="flex">
                                  <FileAttachment_02
                                    svgClass="text-fg-quaternary-light dark:text-fg-quaternary-dark mr-2"
                                  />{$_("viewDocuments")}
                                </div>
                              </DropdownItem>
                            </Dropdown> -->
                          </div>
                        </div>
                      </div>
                    </div>
                    <div
                      class="pb-4 bg-primary-light dark:bg-primary-dark border-b border-secondary-light dark:border-secondary-dark p-4"
                    >
                      <div class="flex justify-between items-center">
                        <div>
                          <P
                            size="sm"
                            weight="semibold"
                            class="text-primary-light dark:text-primary-dark">{$_("clientId")}</P
                          >
                          <P
                            size="sm"
                            weight="normal"
                            class="text-secondary-light dark:text-secondary-dark line-clamp-1"
                            >{bot.clientID}</P
                          >
                        </div>
                        <Button
                          on:click={() => copyToClipboard(bot.clientID)}
                          class="p-0 bg-transparent dark:bg-transparent hover:bg-transparent hover:dark:bg-transparent focus-within:ring-0"
                        >
                          <div>
                            <CopyIcon
                              svgClass="text-fg-brand-primary-light dark:text-fg-brand-primary-dark"
                            />
                          </div>
                        </Button>
                      </div>
                    </div>
                    <div
                      class="pb-4 bg-primary-light dark:bg-primary-dark border-b border-secondary-light dark:border-secondary-dark p-4"
                    >
                      <div class="flex justify-between items-center">
                        <div>
                          <P
                            size="sm"
                            weight="semibold"
                            class="text-primary-light dark:text-primary-dark">{$_("activation")}</P
                          >
                          <P
                            size="sm"
                            weight="normal"
                            class="text-secondary-light dark:text-secondary-dark line-clamp-1"
                            >{bot.activationTime}</P
                          >
                        </div>
                        <Button
                          on:click={() => copyToClipboard(bot.activationTime)}
                          class="p-0 bg-transparent dark:bg-transparent hover:bg-transparent hover:dark:bg-transparent focus-within:ring-0"
                        >
                          <div>
                            <CopyIcon
                              svgClass="text-fg-brand-primary-light dark:text-fg-brand-primary-dark"
                            />
                          </div>
                        </Button>
                      </div>
                    </div>
                    <div
                      class="pb-4 bg-primary-light dark:bg-primary-dark border-b border-secondary-light dark:border-secondary-dark p-4"
                    >
                      <div class="flex justify-between items-center">
                        <div>
                          <P
                            size="sm"
                            weight="semibold"
                            class="text-primary-light dark:text-primary-dark">{$_("urlBot")}</P
                          >
                          <P
                            size="sm"
                            weight="normal"
                            class="text-secondary-light dark:text-secondary-dark line-clamp-2"
                            >{bot.urlBot}</P
                          >
                        </div>
                        <Button
                          on:click={() => copyToClipboard(bot.urlBot)}
                          class="p-0 bg-transparent dark:bg-transparent hover:bg-transparent hover:dark:bg-transparent focus-within:ring-0"
                        >
                          <div>
                            <CopyIcon
                              svgClass="text-fg-brand-primary-light dark:text-fg-brand-primary-dark"
                            />
                          </div>
                        </Button>
                      </div>
                    </div>
                  </Card>
                </div>
              {/if}
            {/each}
          </div>
        </div>
      </div>
    </div>
  </div>
</div>

<div class="flex justify-center">
  <Toast bind:open>{$_("copiedToClipboard")}</Toast>
</div>
