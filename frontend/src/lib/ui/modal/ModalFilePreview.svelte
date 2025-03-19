<script>
  import { Button, Modal } from "flowbite-svelte";
  import { Typography } from "$lib";
  import AudioPlayer from "$lib/components/AudioPlayer.svelte";

  import AlertCircleIcon from "$lib/icon/alert-circle.svg";

  let {
    previewModal = $bindable(false),
    dataSend = $bindable({ detected_text: "" }),
    cancelOnClick,
    okOnClick,
    selectedFile,
    audioUrl,
  } = $props();

  let textarea = $state();

  // Function to adjust the height of the textarea
    function autosize() {
    textarea.style.height = "auto"; // Reset height
    textarea.style.height = `${textarea.scrollHeight}px`; // Set height to scroll height

    textarea.addEventListener("input", (e) => {
      e.target.style.height = "auto";
      e.target.style.height = `${e.target.scrollHeight}px`;
    });
  }

  $effect(() => {
    if (textarea) {
      autosize();
    }
  });
</script>

<Modal bind:open={previewModal} size="lg">
  <svelte:fragment slot="header">
    <div class="flex justify-start items-center gap-4 py-4 px-6">
      <div class="w-8 h-7 p-2 bg-[#feefc6] rounded-full flex justify-center items-center">
        <img src={AlertCircleIcon} alt="AlertCircleIcon" class="w-5 h-5" />
      </div>
      <div class="flex flex-col">
        <Typography preset="Text lg/Semibold" text="Double-check" color="primary" />
        <Typography
          preset="Text sm/Regular"
          text="Are you sure this is the correct document that you want to upload?"
          color="tertiary"
        />
      </div>
    </div>
  </svelte:fragment>
  <div class="flex">
    <div
      class="flex w-full justify-center items-center py-6 rounded-lg bg-secondary-light dark:bg-secondary-dark"
    >
      {#if selectedFile.ext.includes(".png") || selectedFile.ext.includes(".jpg") || selectedFile.ext.includes(".jpeg")}
        <img
          src={`data:${selectedFile.type};base64,${selectedFile.base64}`}
          style="max-width: 80%; height: auto;"
          alt="PNG Preview"
        />
      {/if}
      {#if audioUrl}
        <AudioPlayer src={audioUrl} />
      {/if}
    </div>
    <!-- detected handwriting -->
    {#if selectedFile.type === "Handwriting" || selectedFile.type === "Document Images"}
      <div class="flex w-full flex-col gap-3 px-4 pb-1">
        <Typography
          preset="Text md/Semibold"
          text={selectedFile.type === "Handwriting"
            ? "Handwriting Information"
            : "Document Information"}
          color="secondary"
        />
        <div class="flex flex-col gap-1.5">
          <Typography preset="Text sm/Medium" text="Detected Text" color="secondary" />
          <textarea
            bind:this={textarea}
            bind:value={dataSend.detected_text}
            id="chat"
            rows="1"
            class="block p-2.5 w-full text-sm text-brand-primary-light dark:text-primary-dark bg-white dark:bg-primary-dark rounded-lg border border-utility-brand-500-light dark:border-utility-brand-500-dark focus:border-utility-brand-500-light dark:focus:utility-brand-500-dark focus:ring-utility-brand-400-light dark:focus:ring-utility-brand-400-light dark:placeholder-gray-400"
          ></textarea>
        </div>
      </div>
    {/if}
    {#if selectedFile.type === "Document Audio"}
      <div class="flex w-full flex-col gap-3 px-4 pb-1">
        <Typography preset="Text md/Semibold" text="Audio Information" color="secondary" />
        <div class="flex flex-col gap-1.5">
          <Typography preset="Text sm/Medium" text="Detected Text" color="secondary" />
          <textarea
            bind:this={textarea}
            bind:value={dataSend.detected_text}
            id="chat"
            rows="1"
            class="block p-2.5 w-full text-sm text-brand-primary-light dark:text-primary-dark bg-white dark:bg-primary-dark rounded-lg border border-utility-brand-500-light dark:border-utility-brand-500-dark focus:border-utility-brand-500-light dark:focus:utility-brand-500-dark focus:ring-utility-brand-400-light dark:focus:ring-utility-brand-400-light dark:placeholder-gray-400"
          ></textarea>
        </div>
      </div>
    {/if}
  </div>
  <svelte:fragment slot="footer">
    <div class="flex w-full gap-2 justify-center items-center py-3 px-4">
      <Button
        color="light"
        on:click={() => {
          previewModal = false;
          cancelOnClick();
        }}
        class="w-full">Cancel</Button
      >
      <Button color="primary" on:click={okOnClick(true)} class="rounded-lg w-full">I’m sure</Button>
    </div>
  </svelte:fragment>
</Modal>
