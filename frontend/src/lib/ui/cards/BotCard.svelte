<script>
  import { Card, Button, Avatar } from "flowbite-svelte";
  import { _ } from "svelte-i18n";
  import { Typography, call, triggerToast, check, Badge, Icon } from "$lib";
  import DotsVerticalIcon from "$lib/icon/dots-vertical.svelte";
  import ArrowUpRightIcon from "$lib/icon/arrow-up-right.svelte";
  import Settings03Icon from "$lib/icon/settings-03.svelte";
  import Star04Icon from "$lib/icon/star-04.svelte";
  import Copy01Icon from "$lib/icon/copy-01.svelte";
  import EyeIcon from "$lib/icon/eye.svelte";
  import TrashIcon from "$lib/icon/trash-03.svelte";
  import { goto } from "$app/navigation";

  let {
    bot,
    handleDropdownFocusLoss,
    index,
    dropdownOpen,
    config = null,
    startReading = () => {},
    popupModal = $bindable(false),
    tmpId = $bindable(null),
  } = $props();
</script>

<Card
  padding="none"
  class="relative shadow-none grow w-full {config?.cardFullLength ? 'max-w-none' : ''}"
>
  {#if bot.deployStatus === "Reading"}
    <div class="absolute top-2 right-2">
      <Badge type="warning" text={bot.deployStatus} />
    </div>
  {:else if bot.deployStatus === "Ready"}
    <div class="absolute top-2 right-2">
      <Badge type="success" text={bot.deployStatus} />
    </div>
  {/if}
  <div
    class="bg-active-light dark:bg-active-dark border-b border-secondary-light dark:border-secondary-dark px-4 py-6"
  >
    <div class="flex justify-between items-center">
      <div class="flex items-center gap-3">
        <Avatar src={bot.bot_profile ? bot.bot_profile : ""} />
        <div class="flex flex-col">
          <Typography preset="Text lg/Semibold" color="primary" text={bot.botsName} />
          <Typography preset="Text sm/Regular" color="secondary" text={bot.botDescription} />
        </div>
      </div>
      {#if (check("aiManagement.hexabots:delete") && config?.detail === true) || config?.startReading === true}
        <div onfocusout={handleDropdownFocusLoss} class=" flex items-center relative">
          <Button
            on:click={() => (dropdownOpen[index] = !dropdownOpen[index])}
            class="focus-within:ring-0 !p-0"
            color="none"
            size="xs"
          >
            <DotsVerticalIcon class="text-fg-quinary-light dark:text-fg-quinary-dark" />
          </Button>
          <div
            class={`z-10 ${dropdownOpen[index] ? "absolute top-9 right-0" : "hidden"} bg-white divide-y divide-gray-100 rounded-lg shadow w-max dark:bg-gray-700`}
          >
            <ul class="py-2 grow text-sm">
              <!-- <button
              class="flex py-2 px-4 w-full text-secondary-light dark:text-secondary-dark hover:text-secondary-hover-light dark:hover:text-secondary-hover-dark hover:bg-primary-hover-light dark:hover:bg-primary-hover-dark"
              onclick={goToEditPage(bot.botsID)}>{$_("edit")}</button
            > -->
              {#if config?.startReading === true}
                <button
                  class="flex py-2 px-4 w-full text-secondary-light dark:text-secondary-dark hover:text-secondary-hover-light dark:hover:text-secondary-hover-dark hover:bg-primary-hover-light dark:hover:bg-primary-hover-dark"
                  onclick={() => {
                    tmpId = bot.botsID;
                    dropdownOpen[index] = !dropdownOpen[index];
                    startReading(bot.botsID);
                  }}
                >
                  <EyeIcon svgClass="mr-2 text-fg-quaternary-light dark:text-fg-quaternary-dark" />
                  {$_("startReading")}
                </button>
              {/if}
              {#if config?.delete === true}
                <button
                  class="flex py-2 px-4 w-full text-secondary-light dark:text-secondary-dark hover:text-secondary-hover-light dark:hover:text-secondary-hover-dark hover:bg-primary-hover-light dark:hover:bg-primary-hover-dark"
                  onclick={() => {
                    popupModal = true;
                    tmpId = bot.botsID;
                    dropdownOpen[index] = !dropdownOpen[index];
                  }}
                >
                  <TrashIcon
                    svgClass="mr-2 text-fg-quaternary-light dark:text-fg-quaternary-dark"
                  />
                  {$_("delete")}
                </button>
              {/if}
            </ul>
          </div>
        </div>
      {/if}
    </div>
  </div>
  <div class="flex gap-3 px-4 pt-4 pb-6">
    <div class="flex w-full gap-3">
      {#if config?.detail === true}
        {#if config?.detailView === true}
          <Button
            disabled={!check("aiManagement.hexabots:view")}
            on:click={async () => {
              if (config?.isLlm) {
                goto((await $call.getPage("client-dashboard", "llms-list")) + "?id=" + bot.botsID);
              } else {
                goto(
                  (await $call.getPage("client-dashboard", "bots-detail")) + "?botId=" + bot.botsID,
                );
              }
            }}
            outline
            class="!flex grow rounded-lg *:hover:text-white dark:*:hover:text-white"
          >
            <Typography
              preset="Text sm/Semibold"
              text={config?.detailView ? $_("view") : $_("detail")}
              color="brand-secondary"
            />
          </Button>
        {/if}
        {#if config?.detailTry === true}
          <Button
            on:click={async () => {
              goto(
                (await $call.getPage("client-dashboard", "bots-detail")) +
                  "?botId=" +
                  bot.botsID +
                  "&try=true",
              );
            }}
            outline
            class="!flex grow rounded-lg *:hover:text-white dark:*:hover:text-white"
          >
            <Typography
              preset="Text sm/Semibold"
              text={$_("tryThisModel")}
              color="brand-secondary"
            />
            <Icon
              name="arrowUpRight"
              propClass="text-brand-secondary-light dark:text-brand-secondary-dark"
            />
          </Button>
        {/if}
      {/if}
      {#if config?.configure === true}
        <Button
          disabled={!check("aiManagement.hexabots:view")}
          on:click={async () => {
            goto(
              (await $call.getPage("client-dashboard", "bots-detail")) +
                "?botId=" +
                bot.botsID +
                "&edit=true",
            );
          }}
          outline
          class="!flex grow gap-1 hover:bg-primary-light dark:hover:bg-primary-dark rounded-lg border-secondary-light dark:border-secondary-dark"
        >
          <Typography preset="Text sm/Semibold" text={$_("configure")} color="secondary" />
          <Settings03Icon svgClass="w-5 h-5 text-secondary-light dark:text-secondary-dark" />
        </Button>
      {/if}
      {#if config?.openWebchat === true}
        <Button
          disabled={bot.deployStatus !== "Ready"}
          on:click={() => {
            console.log("open webchat", bot.urlBot);
            window.open(bot.urlBot, "_blank");
          }}
          outline
          class={`!flex grow hover:bg-primary-light dark:hover:bg-primary-dark rounded-lg ${config?.detail === true ? "border-secondary-light dark:border-secondary-dark" : "border-button-secondary-color-light dark:border-button-secondary-color-dark"}`}
        >
          <Typography
            preset="Text sm/Semibold"
            text={$_("webchat")}
            color={config?.detail === true ? "secondary" : "brand-secondary"}
          />
          <ArrowUpRightIcon
            svgClass={`w-5 h-5 ${config?.detail === true ? "text-secondary-light dark:text-secondary-dark" : "text-brand-secondary-light dark:text-brand-secondary-dark"}`}
          />
        </Button>
      {/if}
      {#if config?.insights === true}
        <Button
          on:click={() => {
            // window.open(bot.urlBot);
            goto(bot.urlBot);
          }}
          outline
          class={`!flex grow gap-1.5 hover:bg-primary-light dark:hover:bg-primary-dark rounded-lg ${config?.detail === true ? "border-secondary-light dark:border-secondary-dark" : "border-button-secondary-color-light dark:border-button-secondary-color-dark"}`}
        >
          <Typography
            preset="Text sm/Semibold"
            text={$_("askForInsight")}
            color={config?.detail === true ? "secondary" : "brand-secondary"}
          />
          <Star04Icon
            svgClass={`w-5 h-5 ${config?.detail === true ? "text-secondary-light dark:text-secondary-dark" : "text-brand-secondary-light dark:text-brand-secondary-dark"}`}
          />
        </Button>
      {/if}
    </div>
    <button
      onclick={() => {
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
      <Copy01Icon svgClass="w-5 h-5 text-fg-quinary-light dark:text-fg-quinary-dark" />
    </button>
  </div>
</Card>
