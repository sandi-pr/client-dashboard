<script>
  import { Toast, Progressbar } from "flowbite-svelte";
  import InfoCircleIcon from "$lib/icon/info-circle.svelte";
  import { notifications, Typography, FeaturedIcon } from "$lib";
  import { _ } from "svelte-i18n";
  import StarsIcon from "$lib/icon/stars-01.svelte";

  // function truncateText(text, maxWords) {
  //   const words = text.split(" ");
  //   if (words.length > maxWords) {
  //     return words.slice(0, maxWords).join(" ") + " ...";
  //   }
  //   return text;
  // }

  function truncateText(text, maxWords, maxFirstWordLength = 40) {
    const words = text.split(" ");

    // Check if the first word is too long, or if the second word is long
    if (words[0] && words[0].length > maxFirstWordLength) {
      return words[0].slice(0, maxFirstWordLength) + "...";
    } else if (words[1] && words[1].length > maxFirstWordLength) {
      return words[0] + " " + words[1].slice(0, maxFirstWordLength) + "...";
    }

    // Limit to maxWords and append "..." if necessary
    if (words.length > maxWords) {
      return words.slice(1, maxWords).join(" ") + " ...";
    }

    return words.join(" ");
  }

  function truncateTextByChars(text, maxChars) {
    if (text.length > maxChars) {
      return text.slice(0, maxChars) + "...";
    }
    return text;
  }
</script>

<div class="static h-screen w-full bg-primary-light dark:bg-primary-dark">
  <div class="flex h-full flex-col px-4 pb-20 pt-4">
    <div class="flex flex-col gap-4 overflow-auto md:overflow-y-hidden md:hover:overflow-y-auto">
      {#if $notifications.length === 0}
        <Typography
          preset="Text sm/Medium"
          text={$_("noNotifications")}
          color="secondary"
          class="w-full text-center"
        />
        <!-- <Toast
          align={false}
          dismissable={false}
          divClass="w-full p-4 text-gray-500 bg-white shadow dark:text-gray-400 dark:bg-gray-800 gap-3"
        >
          <div class="flex w-full relative">
            <div class="absolute left-0">
              <FeaturedIcon classesSvg={StarsIcon} />
            </div>
            <div class="flex flex-col w-full pl-12">
              <span class="font-semibold text-gray-900 dark:text-white">{truncateTextByChars("Reading SOP Assistant", 25)}</span>
              <div class="mb-2 text-sm font-normal overflow-hidden text-ellipsis">
                SOP_Stream_Simpanan_ASRS_032771e73c_7fd66b044aASRS_032771e73c_7fd66b044a
              </div>
              <div class="mt-1 gap-2">
                <Progressbar
                  progress={80}
                  size="h-4"
                  labelInside
                />
              </div>
            </div>
          </div>
        </Toast> -->
      {:else}
        <!-- {#if $notifications.length < 1 || $notifications[0].attributes.type.includes("progress")}
          {#if !$notifications[0].attributes.content.includes("100")}
            <Toast
              align={false}
              dismissable={false}
              divClass="w-full p-4 text-gray-500 bg-white shadow dark:text-gray-400 dark:bg-gray-800 gap-3"
            >
              <svelte:fragment slot="icon">
                <InfoCircleIcon svgClass="w-6 h-6" />
              </svelte:fragment>

              <span class="font-semibold text-gray-900 dark:text-white"
                >{$notifications[0].attributes.title}</span
              >
              <div class="mt-1 gap-2">
                {#if $notifications[0].attributes.type.includes("progress")}
                  <Progressbar
                    progress={$notifications[0].attributes.content.replace("%", "")}
                    size="h-4"
                    labelInside
                  />
                {:else}
                  <div class="mb-2 text-sm font-normal">
                    {$notifications[0].attributes.content}
                  </div>
                {/if}
                <div class="flex gap-3 pt-2"></div>
              </div>
            </Toast>
          {/if}
        {/if} -->
        {#each $notifications as data}
          <!-- {#if !data.attributes.type.includes("progress")} -->
          <Toast
            align={false}
            dismissable={false}
            divClass="w-full p-4 text-gray-500 bg-white shadow dark:text-gray-400 dark:bg-gray-800 gap-3"
          >
            <svelte:fragment slot="icon">
              {#if data.attributes.type === "progress" || data.attributes.type === "progress-update"}
                <StarsIcon svgClass="w-5 h-5" />
              {:else}
                <InfoCircleIcon svgClass="w-5 h-5" />
              {/if}
            </svelte:fragment>

            <span class="font-semibold text-gray-900 dark:text-white"
              >{truncateTextByChars(data.attributes.title, 25)}</span
            >
            <div class="mb-2 text-sm font-normal">
              {truncateTextByChars(data.attributes.content, 100)}
            </div>
            <div class="mt-1 gap-2">
              {#if data.attributes.type.includes("progress")}
                {#if data.attributes.meta?.progress}
                  <Progressbar
                    progress={data.attributes.meta?.progress.replace("%", "")}
                    size="h-4"
                    labelInside
                  />
                {/if}
              {/if}
              <div class="flex gap-3 pt-2">
                <!-- {#if !data.attributes.read}
                  <button
                    onclick={async () => {
                      await markAsRead(data.id);
                    }}
                    class="text-sm font-semibold text-fg-button-tertiary-light dark:text-fg-button-tertiary-dark"
                    >Mark as read</button
                  >
                {/if}
                {#if data.attributes.type === "hexabot" && data.attributes.meta && "bot_id" in data.attributes.meta}
                  <button
                    onclick={async () => {
                      goto(
                        (await $call.getPage("client-dashboard", "hexabot-detail")) +
                          "?botId=" +
                          data.attributes.meta.bot_id,
                      );
                    }}
                    class="text-sm font-semibold text-fg-button-tertiary-color-light dark:text-fg-button-tertiary-color-dark"
                    >View detail</button
                  >
                {/if} -->
              </div>
            </div>
          </Toast>
          <!-- {/if} -->
        {/each}
      {/if}
    </div>
  </div>
</div>

<style>
  /* Hide scrollbar for Chrome, Safari and Opera */
  .no-scrollbar::-webkit-scrollbar {
    width: 1px;
  }

  /* Hide scrollbar for IE, Edge and Firefox */
  .no-scrollbar {
    -ms-overflow-style: none; /* IE and Edge */
    scrollbar-width: none; /* Firefox */
  }
</style>
