<script>
  import { Button, P } from "flowbite-svelte";
  import { FeaturedIcon, Typography, UploadFile } from "$lib";
  import { _ } from "svelte-i18n";
  import Divider from "$lib/components/Divider.svelte";
  import FilePlusIcon from "$lib/icon/file-plus-02.svelte";
  import UploadIcon from "$lib/icon/upload-cloud-02.svelte";
  import { ChevronDownOutline } from "flowbite-svelte-icons";

  let dropdownIndustryOpen = $state(false);
  let selectedIndustry = $state("");
  let files = $state([
    // {
    //   type: "DocumentGtc",
    //   name: "Document GTC_PatraNiaga2024.PDF",
    //   size: "200 KB",
    //   base64: "",
    // },
  ]);

  let industries = $state([
    { name: "Oil & Gas", value: "oil-gas" },
    { name: "Other", value: "other" },
  ]);

  const handleDropdownFocusLoss = ({ relatedTarget, currentTarget }) => {
    if (relatedTarget instanceof HTMLElement && currentTarget.contains(relatedTarget)) return;
    dropdownIndustryOpen = false;
  };

  async function readFileAsBase64(file) {
    return new Promise((resolve, reject) => {
      const reader = new FileReader();

      reader.onload = () => {
        const base64Data = reader.result.split(",")[1];
        resolve(base64Data);
      };
      reader.onerror = reject;

      reader.readAsDataURL(file);
    });
  }

  function formatBytes(bytes, decimals = 2) {
    if (bytes === 0) return "0 Bytes";

    const k = 1024;
    const dm = decimals < 0 ? 0 : decimals;
    const sizes = ["Bytes", "KB", "MB", "GB", "TB"];

    const i = Math.floor(Math.log(bytes) / Math.log(k));

    return parseFloat((bytes / Math.pow(k, i)).toFixed(dm)) + " " + sizes[i];
  }

  const handleFileUploadWrapper = (type) => {
    const inputElement = document.createElement("input");
    inputElement.type = "file";
    inputElement.accept = ".pdf, .doc, .docx";

    inputElement.addEventListener("change", async (event) => {
      const uploadedFile = event.target.files[0];

      if (uploadedFile) {
        const base64Data = await readFileAsBase64(uploadedFile);
        files.push({
          type: type,
          name: uploadedFile.name,
          size: formatBytes(uploadedFile.size),
          base64: base64Data,
        });
      } else {
        alert("Invalid file. Please select a valid file.");
      }
    });

    inputElement.click();
  };

  function getFileByType(type) {
    const index = files.findIndex((file) => file.type === type);
    return index !== -1 ? files[index] : [];
  }
</script>

<div class="flex flex-col px-4 pt-6 gap-6">
  <!-- section header -->
  <div class="flex justify-between">
    <div class="flex items-center gap-3">
      <FeaturedIcon classesSvg={FilePlusIcon} />
      <div class="flex flex-col">
        <Typography
          preset="Text lg/Semibold"
          text={$_("setYourReferencedDocument")}
          color="primary"
        />
        <Typography
          preset="Text sm/Regular"
          text={$_("descSetYourReferencedDocument")}
          color="tertiary"
        />
      </div>
    </div>
    <div class="flex gap-3">
      <Button outline on:click={() => {}} color="light">{$_("cancel")}</Button>
      <Button disabled={true} on:click={() => {}}>{$_("save")}</Button>
    </div>
  </div>
  <Divider />
  <!-- section form -->
  <div class="flex gap-3">
    <div class="w-60">
      <Typography preset="Text md/Semibold" text={`${$_("industry")}*`} color="secondary" />
    </div>
    <div class="flex flex-col grow gap-6">
      <!-- form dropdown -->
      <div class="relative flex flex-col gap-1.5">
        <Typography preset="Text sm/Medium" text={`${$_("industryType")}*`} color="secondary" />
        <div onfocusout={handleDropdownFocusLoss} class="flex flex-col gap-1 w-72">
          <Button
            on:click={() => (dropdownIndustryOpen = !dropdownIndustryOpen)}
            class="flex justify-between border-primary-light dark:border-primary-dark"
            color="light"
          >
            {selectedIndustry ? selectedIndustry : "select industry"}
            <ChevronDownOutline class="w-3 h-3 ms-2 text-fg-quaternary-light dark:text-fg-quaternary-dark"/>
          </Button>
          {#if dropdownIndustryOpen}
          <ul
            style:visibility={dropdownIndustryOpen ? "visible" : "hidden"}
            class="py-2 flex flex-col items-center border border-secondary-light dark:border-secondary-dark rounded-lg text-sm text-primary-light dark:text-primary-dark"
            aria-labelledby="dropdownDefaultButton"
          >
            {#each industries as industry}
              <button
                class="flex w-full py-2 px-4"
                onclick={() => ((dropdownIndustryOpen = false), (selectedIndustry = industry.name))}
              >
                {industry.name}
              </button>
            {/each}
          </ul>
          {/if}
        </div>
      </div>
    </div>
  </div>
  {#if selectedIndustry === "Oil & Gas"}
    <Divider />
    <div class="flex justify-end gap-3">
      <Button outline on:click={() => {}} color="light">{$_("cancel")}</Button>
      <Button disabled on:click={() => {}}>{$_("save")}</Button>
    </div>
    <!-- section form -->
    <div class="flex gap-3">
      <div class="w-60">
        <Typography preset="Text md/Semibold" text={`${$_("documentAsBuyer")}`} color="secondary" />
      </div>
      <div class="flex flex-col grow gap-6">
        <div class="grid grid-cols-2 gap-4">
          <!-- form upload -->
          <UploadFile
            title={`${$_("documentGtc")}*`}
            handleFileUploadWrapper={() => handleFileUploadWrapper("DocumentGtc")}
            file={getFileByType("DocumentGtc")}
            actionButton={true}
          />
          <UploadFile
            title={$_("gtcAmendment")}
            handleFileUploadWrapper={() => handleFileUploadWrapper("GtcAmendment")}
            file={getFileByType("GtcAmendment")}
            actionButton={true}
          />
        </div>
      </div>
    </div>
  {/if}
</div>
