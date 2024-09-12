<script>
  import { onMount } from "svelte";
  import {
    P,
    A,
    Card,
    Modal,
    Button,
    Avatar,
    Heading,
    // Dropdown,
    // DropdownItem,
  } from "flowbite-svelte";
  import { _ } from "svelte-i18n";
  import { DotsVerticalOutline } from "flowbite-svelte-icons";
  import MessagePlus from "$lib/icon/message-plus-square.svelte";
  // import MessageChatSquare from "$lib/icon/message-chat-square.svelte";
  // import FileAttachment_02 from "$lib/icon/file-attachment-02.svelte";
  // import PencilLine from "$lib/icon/pencil-line.svelte";
  import CopyIcon from "$lib/icon/copy-03.svelte";
  import Eye from "$lib/icon/eye.svelte";
  import { toCamelCase } from "$lib/utils/toCamelCase.js";
  import { Typography, call } from "$lib";
  import DotsVerticalIcon from "$lib/icon/dots-vertical.svelte";
  import ArrowUpRightIcon from "$lib/icon/arrow-up-right.svelte";
  import Copy01Icon from "$lib/icon/copy-01.svelte";
  // import { PUBLIC_STRAPI_URL_DEV, PUBLIC_STRAPI_URL_PROD } from "$env/static/public";

  // const isDev = process.env.NODE_ENV === "development";

  // const strapiUrl = isDev ? PUBLIC_STRAPI_URL_DEV : PUBLIC_STRAPI_URL_PROD;

  // export let data;
  let data = {
    admin: false,
  };

  let moduleMonitoringData = [
    {
      href: "/hexabot/create",
      icon: MessagePlus,
      moduleName: "Create New Bots",
      desc: "desc Create New Bots",
    },
    // {
    //   href: "#",
    //   action: () => {
    //     deployStatusModal = true;
    //   },
    //   icon: Eye,
    //   moduleName: "View bot deployment status",
    //   desc: "desc View Bot Deployment Status",
    // },
  ];

  let Header;

  let bots = [
    {
      botsName: "Bot 1",
      botDescription: "Bot 1 Description",
      botsID: "711",
      clientID: "b608d45f-0f54-45f7-816e-8d6aa82e6eba",
      activationTime: "2024-02-06 10:30:13",
      urlBot: "https://webchat.hexacode.org/1fff",
      timestamp: "2024-02-06 10:30:13",
      deployStatus: "Done",
      deployedPlatform: "Web",
    },
  ];

  let deployStatusModal = false;

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

  async function getBots() {
    let token = sessionStorage.getItem("token");

    if (token) {
      try {
        const response = await $call.call("client-dashboard", "get-bots", {
          URLQuery: {},
          URLParams: {},
          Params: {
            page: 1,
            page_size: 8,
          },
          Headers: {
            Authorization: "Bearer " + token,
          },
        });

        console.log(response);

        if (response.status !== 200) {
          console.log(response);
          return;
        }

        let tmp = response.result.data.map((bot) => {
          let fileImageResult = "";
          if (bot.attributes.bot_profile.data !== null) {
            let ext = bot.attributes.bot_profile.data.attributes.ext;

            let fileType;
            if (ext === ".png") {
              fileType = "image/png";
            } else if (ext === ".jpg") {
              fileType = "image/jpg";
            } else if (ext === ".jpeg") {
              fileType = "image/jpeg";
            } else if (ext === ".svg") {
              fileType = "image/svg+xml";
            } else {
              fileType = "image/png";
            }

            if (bot.bot_profile !== "") {
              let fileImage = bot.bot_profile;
              let blob = base64toBlob(fileImage, fileType);
              fileImageResult = URL.createObjectURL(blob);
            }
          }
          return {
            botsName: bot.attributes.bot_name,
            botDescription: bot.attributes.description,
            botsID: bot.id,
            // clientID: bot.attributes.client_id,
            activationTime: bot.attributes.createdAt,
            urlBot: "",
            timestamp: bot.attributes.updatedAt,
            // deployStatus: bot.attributes.deploy_status,
            deployedPlatform: bot.attributes.ai_engine,
            bot_profile: fileImageResult,
          };
        });

        bots = await Promise.all(tmp);
      } catch (error) {
        console.log(error);
      }
    } else {
      console.log("token not found");
    }
  }

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

  onMount(() => {
    getBots();
  });
</script>

<div class="flex dark:bg-hc-graydark-950">
  <div class="w-full pb-10">
    <div class="flex flex-col m-4">
      <!-- <Navbar /> -->

      <div class="flex flex-col gap-8">
        <!-- User Information -->
        {#if !data.admin}
          <Heading tag="h2" class="font-jakarta text-hc-graylight-900 dark:text-hc-graydark-50"
            >{$_("hexaBots")}</Heading
          >

          <!-- Module Monitoring Data -->
          <div>
            <div
              class="flex border-b border-hc-graylight-200 dark:border-hc-graydark-800 pb-4 mb-4"
            >
              <Heading tag="h6" class="font-jakarta text-hc-graylight-900 dark:text-hc-graydark-50"
                >{$_("action")}</Heading
              >
            </div>
            <div class="grid grid-cols-1 lg:grid-cols-3 gap-4 w-full items-center justify-center">
              {#each moduleMonitoringData as { href, action, icon, moduleName, desc }}
                <Card
                  {href}
                  on:click={action}
                  horizontal="true"
                  class="items-center w-full !p-4 min-h-28 bg-primary-light dark:bg-primary-dark hover:bg-primary-hover-light dark:hover:bg-primary-hover-dark border-brand-light dark:border-brand-dark shadow-sm hover:no-underline"
                  size="none"
                >
                  <Avatar class="!bg-hc-brand-100 dark:!bg-hc-brand-900 !p-2 !h-10"
                    ><svelte:component
                      this={icon}
                      svgClass="w-5 h-5 text-fg-brand-primary-light dark:text-fg-brand-primary-dark"
                    /></Avatar
                  >
                  <div class="ml-4 gap[2px]">
                    <P
                      class="font-jakarta text-base font-semibold leading-normal text-hc-graylight-900 dark:text-hc-graydark-50"
                      >{$_(toCamelCase(moduleName))}</P
                    >
                    <P
                      class="font-jakarta text-sm font-normal leading-tight text-hc-graylight-700 dark:text-hc-graydark-300"
                      >{$_(toCamelCase(desc))}</P
                    >
                  </div>
                </Card>
              {/each}
            </div>
          </div>

          <!-- Document Monitoring Data -->
          <!-- <div>
          <div class="flex border-b border-secondary-light dark:border-secondary-dark pb-4 mb-4">
            <Heading tag="h6" class="font-jakarta text-primary-light dark:text-primary-dark">
              Reports
            </Heading>
          </div>
          <div
            class="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-3 gap-4 items-center justify-center"
          >
            {#each docData as { href, dataName, value, change, changeType }, index ({ href, dataName, value, change, changeType })}
              <Card
                horizontal="true"
                class="grid grid-cols-12 items-center !px-4 !py-0 min-h-28 justify-between bg-primary-light dark:bg-primary-dark border-secondary-light dark:border-secondary-dark shadow-sm"
                size="none"
              >
                <div class="space-y-2 col-span-6 sm:col-span-7 md:col-span-5">
                  <div class="flex items-center gap-2">
                    <P
                      size="base"
                      weight="semibold"
                      color="text-secondary-light dark:text-secondary-dark"
                      class="font-jakarta">{dataName}</P
                    >
                    <P
                      size="sm"
                      weight="bold"
                      color={`${
                        changeType === "decrease"
                          ? "text-error-primary-light dark:text-error-primary-dark"
                          : "text-success-primary-light dark:text-success-primary-dark"
                      }`}
                      class="font-jakarta leading-tight"
                      >{`${changeType === "decrease" ? "-" : "+"}${change}`}</P
                    >
                  </div>
                  <P
                    size="xl"
                    weight="semibold"
                    color="text-gray-900"
                    class="font-jakarta leading-[30px] text-primary-light dark:text-primary-dark"
                    >{value}</P
                  >
                  <A
                    {href}
                    class="text-sm font-semibold text-fg-brand-primary-light dark:text-fg-brand-primary-dark font-jakarta leading-tight"
                    >See All</A
                  >
                </div>
                <div class="col-span-6 sm:col-span-5 md:col-span-7">
                  {#if chartReady}
                    {#if index === 0}
                      <Chart dataChart={dataChartAll} />
                    {:else if index === 1}
                      <Chart dataChart={dataChartNeed} />
                    {:else}
                      <Chart dataChart={dataChartComplete} />
                    {/if}
                  {:else}
                    <div></div>
                  {/if}
                </div>
              </Card>
            {/each}
          </div>
        </div> -->

          <!-- Document Table Data -->
          <div>
            <div
              class={`flex items-center justify-between pb-4 mb-4 border-b border-hc-graylight-200 dark:border-hc-graydark-800`}
            >
              <Heading
                tag="h6"
                class="font-jakarta w-fit text-hc-graylight-900 dark:text-hc-graydark-50"
                >{$_("recentDeployedBots")}</Heading
              >
              <Button
                href="/hexabot/deployed"
                color="light"
                class="bg-hc-graylight-50 dark:bg-hc-graydark-900 border-hc-graylight-200 dark:border-hc-graydark-800 hover:bg-hc-graylight-100 hover:dark:bg-hc-graydark-800 focus-within:ring-0"
                >{$_("viewAll")}</Button
              >
            </div>
            <div class="grid grid-cols-1 md:grid-cols-2 2xl:grid-cols-4 gap-4">
              {#each bots as bot (bot.botsID)}
                <!-- {#if bot.deployStatus === "Done"} -->
                <Card padding="none" class="shadow-none grow w-full">
                  <div
                    class="bg-active-light dark:bg-active-dark border-b border-secondary-light dark:border-secondary-dark px-4 py-6"
                  >
                    <div class="flex justify-between items-center">
                      <div class="flex items-center gap-3">
                        <Avatar src={bot.bot_profile ? bot.bot_profile : ""} />
                        <div class="flex flex-col">
                          <Typography
                            preset="Text lg/Semibold"
                            color="primary"
                            text={bot.botsName}
                          />
                          <Typography
                            preset="Text sm/Regular"
                            color="secondary"
                            text={bot.botDescription}
                          />
                        </div>
                      </div>
                      <DotsVerticalIcon svgClass="w-5 h-5 text-fg-quinary" />
                    </div>
                  </div>
                  <div class="flex gap-3 px-4 pt-4 pb-6">
                    <div class="grid grow grid-cols-2 gap-3">
                      <Button
                        on:click={() => {}}
                        outline
                        class="!flex grow hover:bg-primary-light dark:hover:bg-primary-dark rounded-lg"
                      >
                        <Typography
                          preset="Text sm/Semibold"
                          text={$_("view")}
                          color="brand-secondary"
                        />
                      </Button>
                      <Button
                        on:click={() => {}}
                        outline
                        class="!flex grow hover:bg-primary-light dark:hover:bg-primary-dark rounded-lg border-secondary-light dark:border-secondary-dark"
                      >
                        <Typography
                          preset="Text sm/Semibold"
                          text={$_("webchat")}
                          color="secondary"
                        />
                        <ArrowUpRightIcon
                          svgClass="w-5 h-5 text-secondary-light dark:text-secondary-dark"
                        />
                      </Button>
                    </div>
                    <button
                      on:click={() => {
                        navigator.clipboard.writeText(bot.urlBot);
                      }}
                    >
                      <Copy01Icon svgClass="w-5 h-5 text-quinary-light dark:text-quinary-dark" />
                    </button>
                  </div>
                </Card>
                <!-- <div class="span-col-1">
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
                              href="hexabot/detail"
                              class="text-fg-brand-primary-light dark:text-fg-brand-primary-dark text-sm font-semibold"
                              >{$_("view")}</A
                            >
                            <div class="cursor-pointer flex items-center">
                              <Button class="focus-within:ring-0 !p-0" color="none" size="xs"
                                ><DotsVerticalOutline /></Button
                              >
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
                          <div class="flex pl-2">
                            <div class="cursor-pointer">
                              <CopyIcon
                                svgClass="text-fg-brand-primary-light dark:text-fg-brand-primary-dark"
                              />
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
                              class="text-primary-light dark:text-primary-dark"
                              >{$_("activation")}</P
                            >
                            <P
                              size="sm"
                              weight="normal"
                              class="text-secondary-light dark:text-secondary-dark line-clamp-1"
                              >{bot.activationTime}</P
                            >
                          </div>
                          <div class="flex pl-2">
                            <div class="cursor-pointer">
                              <CopyIcon
                                svgClass="text-fg-brand-primary-light dark:text-fg-brand-primary-dark"
                              />
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
                              class="text-primary-light dark:text-primary-dark">{$_("urlBot")}</P
                            >
                            <P
                              size="sm"
                              weight="normal"
                              class="text-secondary-light dark:text-secondary-dark line-clamp-1"
                              >{bot.urlBot}</P
                            >
                          </div>
                          <div class="flex pl-2">
                            <div class="cursor-pointer">
                              <CopyIcon
                                svgClass="text-fg-brand-primary-light dark:text-fg-brand-primary-dark"
                              />
                            </div>
                          </div>
                        </div>
                      </div>
                    </Card>
                  </div> -->
                <!-- {/if} -->
              {/each}
            </div>
          </div>
        {:else}
          <Heading tag="h2" class="font-jakarta text-hc-graylight-900 dark:text-hc-graydark-50"
            >{$_("Hexa AI")}</Heading
          >

          <!-- Document Table Data -->
          <div>
            <div
              class={`flex items-center justify-between pb-4 mb-4 border-b border-hc-graylight-200 dark:border-hc-graydark-800`}
            >
              <Heading
                tag="h6"
                class="font-jakarta w-fit text-hc-graylight-900 dark:text-hc-graydark-50"
                >Bots</Heading
              >
            </div>
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
                              href="#"
                              class="text-fg-brand-primary-light dark:text-fg-brand-primary-dark text-sm font-semibold"
                              >{$_("view")}</A
                            >
                            <div class="cursor-pointer flex items-center">
                              <Button class="focus-within:ring-0 !p-0" color="none" size="xs"
                                ><DotsVerticalOutline /></Button
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
                              class="text-primary-light dark:text-primary-dark">{$_("clientID")}</P
                            >
                            <P
                              size="sm"
                              weight="normal"
                              class="text-secondary-light dark:text-secondary-dark line-clamp-1"
                              >{bot.clientID}</P
                            >
                          </div>
                          <div class="flex pl-2">
                            <div class="cursor-pointer">
                              <CopyIcon
                                svgClass="text-fg-brand-primary-light dark:text-fg-brand-primary-dark"
                              />
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
                              class="text-primary-light dark:text-primary-dark"
                              >{$_("activation")}</P
                            >
                            <P
                              size="sm"
                              weight="normal"
                              class="text-secondary-light dark:text-secondary-dark line-clamp-1"
                              >{bot.activationTime}</P
                            >
                          </div>
                          <div class="flex pl-2">
                            <div class="cursor-pointer">
                              <CopyIcon
                                svgClass="text-fg-brand-primary-light dark:text-fg-brand-primary-dark"
                              />
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
                              class="text-primary-light dark:text-primary-dark">{$_("urlBot")}</P
                            >
                            <P
                              size="sm"
                              weight="normal"
                              class="text-secondary-light dark:text-secondary-dark line-clamp-1"
                              >{bot.urlBot}</P
                            >
                          </div>
                          <div class="flex pl-2">
                            <div class="cursor-pointer">
                              <CopyIcon
                                svgClass="text-fg-brand-primary-light dark:text-fg-brand-primary-dark"
                              />
                            </div>
                          </div>
                        </div>
                      </div>
                    </Card>
                  </div>
                {/if}
              {/each}
            </div>
          </div>
        {/if}
      </div>
    </div>
  </div>
</div>

<Modal bind:open={deployStatusModal} autoclose>
  <div class="space-y-2">
    <Heading tag="h5">{$_("botDeploymentStatus")}</Heading>
    <P>{$_("descBotDeploymentStatus")}</P>
  </div>
  {#each bots as bot}
    {#if bot.deployStatus !== "Done"}
      <Card size="lg">
        <div class="flex justify-between items-center">
          <div>
            <Heading tag="h6">{`${bot.botsName}`}</Heading>
            <P>{`${bot.timestamp}`}</P>
          </div>
          <Button
            class={`p-[2px] border-[1px] h-1/2 rounded-xl text-center ${classStatusType(
              bot.deployStatus,
            )}`}>{bot.deployStatus}</Button
          >
        </div>
      </Card>
    {/if}
  {/each}
  <Button
    on:click={() => {
      deployStatusModal = false;
    }}
    class="w-full"
    color="light">{$_("close")}</Button
  >
</Modal>
