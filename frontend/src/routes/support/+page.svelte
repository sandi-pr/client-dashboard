<script>
  import { Typography, Divider } from "$lib";
  import { Button } from "flowbite-svelte";
  import { _ } from "svelte-i18n";
  import LogoIcon from "$lib/icon/Logo.svelte";
  import ChevronUpIcon from "$lib/icon/chevron-up.svelte";
  import ChevronDownIcon from "$lib/icon/chevron-down.svelte";
  import SendIcon from "$lib/icon/send-01.svg";
  import { onMount } from "svelte";

  let messageInput = $state("");
  let textarea = $state(null);
  let innerWidth = $state(0);
  let winInnerWidth = $derived(innerWidth);
  let dropdownToggle = $state([true, false, false]);

  async function sendMessageInput() {
    messageInput = "";
    resetTextarea();
    console.log("messageInput", messageInput);
  }

  const handleKeyDown = (event) => {
    if (winInnerWidth > 768) {
      if (event.key === "Enter") {
        if (!event.shiftKey) {
          event.preventDefault();
          if (
            messageInput.trim() !== ""
          ) {
            onclickSend();
          }
        }
      }
    }
  };

  const onclickSend = () => {
    if (messageInput.trim() !== "") {
      // sendMessageInput();
    }
  };

  const autosize = () => {
    textarea.addEventListener("input", (e) => {
      setTimeout(() => {
        textarea.style.height = "auto";
        let scHeight = e.target.scrollHeight + 2; // Tambahkan sedikit untuk padding
        textarea.style.height = `${scHeight}px`;
      }, 0);
    });
  };

  const resetTextarea = () => {
    textarea.style.height = "auto";
  };

  onMount(() => {
    if (textarea) {
      autosize();
    }
    const resizeHandler = () => {
      innerWidth = globalThis.innerWidth;
    };
    globalThis.addEventListener("resize", resizeHandler);
    resizeHandler(); // call once initially to set the correct value
    return () => {
      globalThis.removeEventListener("resize", resizeHandler);
    };
  });

</script>

<div class="px-4 pt-6">
  <Typography preset="Display sm/Semibold" text={$_("support")} color="primary" />
  <div class="flex flex-col items-center gap-2">
    <LogoIcon />
    <Typography preset="Display xs/Semibold" text={$_("hexaProductSupport")} color="secondary" />
    <Typography preset="Text md/Regular" text={$_("supportDesc")} color="secondary" propClass="text-center" />
    <form
      onsubmit={messageInput !== "" ? sendMessageInput : () => {}}
      class="w-full p-6"
    >
      <div class="flex items-center">
        <div class="w-full flex gap-4">
          <!-- svelte-ignore event_directive_deprecated -->
          <textarea
            bind:this={textarea}
            bind:value={messageInput}
            onkeydown={handleKeyDown}
            id="chat"
            rows="1"
            class="block resize-none p-2.5 w-full fluid-sm text-brand-primary-light dark:text-primary-dark bg-white dark:bg-primary-dark rounded-xl border border-utility-brand-500-light dark:border-utility-brand-500-dark focus:border-utility-brand-500-light dark:focus:border-utility-brand-500-dark focus:ring-utility-brand-400-light dark:focus:ring-utility-brand-400-light dark:placeholder-gray-400"
            placeholder="Send a message..."
          ></textarea>
          <Button
            disabled={!messageInput}
            onclick={onclickSend}
            class="rounded-full h-11 w-11 !p-3"
          >
            <img src={SendIcon} alt="SendIcon" class="w-5 h-5" />
          </Button>
        </div>
      </div>
    </form>
    <div class="w-full">
      <Typography preset="Text lg/Semibold" text={$_("frequentlyAskedQuestion")} color="primary" propClass="pb-5" />
      <Divider />
      <Button
        on:click={() => (dropdownToggle[0] = !dropdownToggle[0])}
        class="flex w-full justify-between border-primary-light dark:border-primary-dark bg-secondary-light dark:bg-secondary-dark"
        color="light"
      >
        <Typography preset="Text lg/Semibold" text={$_("howDoIAccessTheChatbot")} color="primary" />
        {#if dropdownToggle[0]}
          <ChevronDownIcon svgClass="w-5 h-5 text-utility-gray-700-light dark:text-utility-gray-700-dark" />
        {:else}
          <ChevronUpIcon svgClass="w-5 h-5 text-utility-gray-700-light dark:text-utility-gray-700-dark" />
        {/if}
      </Button>
      {#if dropdownToggle[0]}
        <div class="pt-4 px-6">
          <Typography preset="Text md/Regular" text={$_("howDoIAccessTheChatbotDesc")} color="secondary" />
        </div>
      {/if}
    </div>
  </div>
</div>