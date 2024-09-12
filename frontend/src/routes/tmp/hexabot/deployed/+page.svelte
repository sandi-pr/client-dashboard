<script>
  import {
    P,
    Card,
    Button,
    Heading,
    Avatar,
    Toast,
  } from "flowbite-svelte";
  import { _ } from "svelte-i18n";
  import DotsVerticalIcon from "$lib/icon/dots-vertical.svelte";
  import ArrowUpRightIcon from "$lib/icon/arrow-up-right.svelte";
  import Copy01Icon from "$lib/icon/copy-01.svelte";
  import { Typography } from "$lib";

  let searchTerm = $state("");
  let bots = [
    {
      botsName: "Chat Qur’ani",
      botDescription: "Bot’s Description, Lorem ipsum dolor sit amet consectetur",
      botsID: "5904",
      clientID: "c2287dc7-0af7-4481-9ec3-dbaa910477e7",
      activationTime: "2023-12-13 05:00:56",
      urlBot: "https://webchat.aws.hexacode.org/25dfb75d-3985-40c5-8c7b-1ece3f50699a",
      timestamp: "2023-12-13 05:00:56",
      deployStatus: "Done",
      deployedPlatform: "Web",
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
          <div class="grid grid-cols-1 md:grid-cols-2 2xl:grid-cols-3 gap-4">
            {#each bots as bot (bot.botsID)}
              {#if bot.deployStatus === "Done"}
              <Card padding="none" class="shadow-none">
                <div
                  class="bg-active-light dark:bg-active-dark border-b border-secondary-light dark:border-secondary-dark px-4 py-6"
                >
                  <div class="flex justify-between items-center">
                    <div class="flex items-center gap-3">
                      <Avatar />
                      <div class="flex flex-col">
                        <Typography preset="Text lg/Semibold" color="primary" text={bot.botsName} />
                        <Typography preset="Text sm/Regular" color="secondary" text={bot.botDescription} />
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
                      <Typography preset="Text sm/Semibold" text={$_("view")} color="brand-secondary" />
                    </Button>
                    <Button
                      on:click={() => {}}
                      outline
                      class="!flex grow hover:bg-primary-light dark:hover:bg-primary-dark rounded-lg border-secondary-light dark:border-secondary-dark"
                    >
                      <Typography preset="Text sm/Semibold" text={$_("webchat")} color="secondary" />
                      <ArrowUpRightIcon svgClass="w-5 h-5 text-secondary-light dark:text-secondary-dark" />
                    </Button>
                  </div>
                  <button>
                    <Copy01Icon svgClass="w-5 h-5 text-quinary-light dark:text-quinary-dark" />
                  </button>
                </div>
              </Card>
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
