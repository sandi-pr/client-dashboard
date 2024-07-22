<script>
  import { _ } from "svelte-i18n";
  import { Button } from "flowbite-svelte";
  import { Typography } from "$lib";
  import UploadIcon from "$lib/icon/upload-cloud-02.svelte";
  import PDFSolidIcon from "$lib/icon/PDFSolid.svg";
  import DOCSolidIcon from "$lib/icon/DOCSolid.svg";
  import DOCXSolidIcon from "$lib/icon/DOCXSolid.svg";

  let { 
    title, 
    handleFileUploadWrapper, 
    multiple = false, 
    type = "document", 
    file = [], 
    actionButton = false, 
    actionView = () => {},
    actionDelete = () => {},
  } = $props();

  const textFileType = (type) => {
    if (type === "document") {
      return "PDF, DOC, DOCX (Max. 15MB)";
    } else {
      return "";
    }
  };

</script>

<div class="relative flex flex-col gap-1.5">
  <Typography preset="Text sm/Medium" text={title} color="secondary" />
  {#if file.length === 0}
  <Button
    on:click={handleFileUploadWrapper}
    color="light"
    class={`w-full rounded-xl border border-primary-light dark:border-primary-dark bg-primary-light dark:bg-primary-dark gap-3 py-4 flex flex-col justify-start items-center`}
  >
    <div class="rounded-xl p-2.5 shadow border border-primary-light dark:border-primary-dark justify-center items-center">
      <UploadIcon class="w-5 h-5 text-fg-secondary-light dark:text-fg-secondary-dark" />
    </div>
    <div class="flex flex-col items-center gap-0.5">
      <Typography preset="Text sm/Semibold" text={$_("clickToUpload")} color="brand-secondary" />
      <Typography preset="Text xs/Regular" text={textFileType(type)} color="tertiary" />
    </div>
  </Button>
  {:else}
    <div class="flex flex-col gap-3 px-6 py-4 rounded-xl border border-primary-light dark:border-primary-dark">
      <div class="flex items-center justify-start gap-3">
        {#if file.name.includes(".pdf") || file.name.includes(".PDF")}
          <img src={PDFSolidIcon} alt="" />
        {:else if file.name.includes(".docx")}
          <img src={DOCXSolidIcon} alt="" />
        {:else if file.name.includes(".doc")}
          <img src={DOCSolidIcon} alt="" />
        {/if}
        <div class="flex flex-col justify-center">
          <Typography preset="Text sm/Semibold" text={file.name} color="secondary" />
          <Typography preset="Text sm/Regular" text={file.size} color="tertiary" />
        </div>
      </div>
      {#if actionButton}
        <div class="flex gap-3">
          <Button class="w-full" color="light" on:click={actionView}>
            <Typography preset="Text sm/Semibold" text={$_("view")} color="secondary" />
          </Button>
          <Button class="w-full border-0" color="light" on:click={actionDelete}>
            <Typography preset="Text sm/Semibold" text={$_("delete")} color="secondary" />
          </Button>
        </div>
      {/if}
    </div>
  {/if}
</div>