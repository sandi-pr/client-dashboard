<script>
  import { P, Card, Button, Heading, Avatar, Toast } from "flowbite-svelte";
  import { _ } from "svelte-i18n";
  import DotsVerticalIcon from "$lib/icon/dots-vertical.svelte";
  import ArrowUpRightIcon from "$lib/icon/arrow-up-right.svelte";
  import Copy01Icon from "$lib/icon/copy-01.svelte";
  import { Typography, call, triggerToast } from "$lib";
  import { onMount } from "svelte";
  import { goto } from "$app/navigation";

  let searchTerm = $state("");
  let bots = $state([]);
  let page = $state(1);

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

  async function getBots() {
    let token = sessionStorage.getItem("token");

    if (token) {
      try {
        const response = await $call.call("client-dashboard", "get-bots", {
          URLQuery: {},
          URLParams: {},
          Params: {
            page: page,
            page_size: 18,
          },
          Headers: {
            Authorization: "Bearer " + token,
          },
        });

        if (response === null) {
          console.log(response);
          return;
        }

        if (response.status !== 200) {
          console.log(response);
          return;
        }

        let tmp = response.result.data.map(async (bot) => {
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

          let urlBot = (await $call.getPage("hexabot", "home-page")) + "?id=" + bot.bot_id;

          return {
            botsName: bot.attributes.bot_name,
            botDescription: bot.attributes.description,
            botsID: bot.id,
            // clientID: bot.attributes.client_id,
            activationTime: bot.attributes.createdAt,
            urlBot: urlBot,
            timestamp: bot.attributes.updatedAt,
            // deployStatus: bot.attributes.deploy_status,
            deployedPlatform: bot.attributes.ai_engine,
            bot_profile: fileImageResult,
          };
        });

        bots = await Promise.all(tmp);
        console.log(bots);
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
        <Typography preset="Display sm/Semibold" text={$_("deployedBots")} color="secondary" />

        <!-- Bots List Data -->
        <div>
          <!-- <Search
            class="bg-white dark:bg-hc-graydark-950 mb-6 rounded-lg"
            placeholder={$_("searchBotName")}
            bind:value={searchTerm}
          /> -->
          <div class="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-3 2xl:grid-cols-4 gap-4">
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
                        <Typography preset="Text lg/Semibold" color="primary" text={bot.botsName} />
                        <Typography
                          preset="Text sm/Regular"
                          color="secondary"
                          text={bot.botDescription}
                        />
                      </div>
                    </div>
                    <!-- <DotsVerticalIcon svgClass="w-5 h-5 text-fg-quinary" /> -->
                  </div>
                </div>
                <div class="flex gap-3 px-4 pt-4 pb-6">
                  <div class="grid grow grid-cols-2 gap-3">
                    <Button
                      on:click={async () => {
                        goto(
                          (await $call.getPage("client-dashboard", "hexabot-detail")) +
                            "?botId=" +
                            bot.botsID,
                        );
                      }}
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
                      on:click={() => {
                        goto(bot.urlBot);
                      }}
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
                      navigator.clipboard
                        .writeText(bot.urlBot)
                        .then(() => {
                          console.log("Text copied to clipboard");
                          triggerToast("copyToClipboard");
                        })
                        .catch((err) => {
                          console.error("Failed to copy: ", err);
                        });
                    }}
                  >
                    <Copy01Icon svgClass="w-5 h-5 text-quinary-light dark:text-quinary-dark" />
                  </button>
                </div>
              </Card>
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
