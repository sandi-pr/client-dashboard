<script>
  import { Toast, Progressbar } from "flowbite-svelte";
  import InfoCircleIcon from "$lib/icon/info-circle.svelte";
  import { notifications, Typography } from "$lib";

  // function truncateText(text, maxWords) {
  //   const words = text.split(" ");
  //   if (words.length > maxWords) {
  //     return words.slice(0, maxWords).join(" ") + " ...";
  //   }
  //   return text;
  // }

  function truncateText(text, maxWords, maxFirstWordLength = 40) {
    const words = text.split(" ");

    // Check if the first word is too long
    if (words[1].length > maxFirstWordLength) {
      words[1] = words[1].slice(1, maxFirstWordLength) + "...";
    }

    // Limit to maxWords and append "..." if necessary
    if (words.length > maxWords) {
      return words.slice(1, maxWords).join(" ") + " ...";
    }

    return words.join(" ");
  }
</script>

<div class="static h-screen w-full bg-primary-light dark:bg-primary-dark">
  <div class="flex h-full flex-col px-4 pb-20 pt-4">
    <div class="flex flex-col gap-4 overflow-auto md:overflow-y-hidden md:hover:overflow-y-auto">
      {#if $notifications.length === 0}
        <Typography
          preset="Text sm/Medium"
          text="No Notification"
          color="secondary"
          class="w-full text-center"
        />
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
              <InfoCircleIcon svgClass="w-6 h-6" />
            </svelte:fragment>

            <span class="font-semibold text-gray-900 dark:text-white">{truncateText(data.attributes.title, 2)}</span>
            <div class="mt-1 gap-2">
              {#if data.attributes.type.includes("progress")}
                <Progressbar
                  progress={data.attributes.content.replace("%", "")}
                  size="h-4"
                  labelInside
                />
              {:else}
                <div class="mb-2 text-sm font-normal">
                  {data.attributes.content}
                </div>
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
