<script>
  import { Button, Modal, Input } from "flowbite-svelte";
  import { FeaturedIcon, Typography, UploadFile, triggerToast } from "$lib";
  import PdfViewer from "svelte-pdf";
  import { _ } from "svelte-i18n";
  import Divider from "$lib/components/Divider.svelte";
  import FilePlusIcon from "$lib/icon/file-plus-02.svelte";
  import { ChevronDownOutline } from "flowbite-svelte-icons";
  import PlusCircleICon from "$lib/icon/plus-circle.svelte";
  import { CR3Call } from "hexacomponents-svelte";

  const call = new CR3Call(globalThis.location.origin);

  let dropdownIndustryOpen = $state(false);
  let documentViewModal = $state(false);
  let addNewDocumentModal = $state(false);
  let tmpDocumentName = $state("");
  let tmpFile = $state([]);
  let selectedIndustry = $state("");
  let files = $state([]);
  let selectedFile = $state(null);
  let requiredFileTypes = $state(["documentGtc", "tradeConfirmation", "dealDoneDocument"]);
  let roles = $state(["buyer", "seller"]);

  let industries = $state([
    { name: "Oil & Gas", value: "oil-gas" },
    { name: "Other", value: "other" },
  ]);

  let documentTypes = $state([
    {
      title: `${$_("documentGtc")}*`,
      type: "documentGtc",
    },
    {
      title: $_("gtcAmendment"),
      type: "gtcAmendment",
    },
    {
      title: `${$_("tradeConfirmation")}*`,
      type: "tradeConfirmation",
    },
    {
      title: `${$_("dealDoneDocument")}*`,
      type: "dealDoneDocument",
    },
  ]);

  let otherDocType = $state([]);

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

  const handleFileUploadWrapper = (type, role) => {
    const inputElement = document.createElement("input");
    inputElement.type = "file";
    inputElement.accept = ".pdf, .doc, .docx";

    inputElement.addEventListener("change", async (event) => {
      const uploadedFile = event.target.files[0];

      if (uploadedFile) {
        const base64Data = await readFileAsBase64(uploadedFile);
        files.push({
          type: type,
          role: role,
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

  function addOtherDocument() {
    if (tmpDocumentName) {
      otherDocType.push({
        title: tmpDocumentName,
        type: tmpDocumentName,
      });

      files.push(tmpFile);

      tmpFile = [];
      tmpDocumentName = "";
      addNewDocumentModal = false;
    }
  }

  function getFileByType(type, role) {
    const index = files.findIndex((file) => file.type === type && file.role === role);
    return index !== -1 ? files[index] : [];
  }

  function deleteFileByType(type, role) {
    const index = files.findIndex((file) => file.type === type && file.role === role);
    if (index !== -1) {
      files.splice(index, 1);
      if (role === "other") {
        otherDocType = otherDocType.filter((doc) => doc.type !== type);
      }
    }
  }

  function areRequiredFilesPresent(files) {
    return roles.every((role) =>
      requiredFileTypes.every((type) =>
        files.some((file) => file.type === type && file.role === role),
      ),
    );
  }

  let allRequiredFilesPresent = $state(false);

  $effect(() => {
    allRequiredFilesPresent = areRequiredFilesPresent(files);
  });

  async function saveReferences() {
    if (allRequiredFilesPresent) {
      let token = sessionStorage.getItem("token");
      try {
        const response = await call.call("client-dashboard", "set-references", {
          URLQuery: {},
          URLParams: {},
          Params: {
            references: JSON.stringify(files),
            industry_type: selectedIndustry,
          },
          Headers: {
            Authorization: "Bearer " + token,
          },
        });
        console.log(response);

        if (response.status === 200) {
          triggerToast("updated");
        }
      } catch (error) {
        console.log(error);
      }
    }
  }

  async function getReferences() {
    let token = sessionStorage.getItem("token");
    try {
      const response = await call.call("client-dashboard", "get-references", {
        URLQuery: {},
        URLParams: {},
        Params: {
          industry_type: selectedIndustry,
        },
        Headers: {
          Authorization: "Bearer " + token,
        },
      });

      console.log(response);
      if (response.status === 200) {
        if (response.result.files !== null) {
          files = response.result.files;
        }
      }
    } catch (error) {
      console.log(error);
    }
  }

  $effect(() => {
    if (selectedIndustry) {
      if (selectedIndustry !== "") {
        getReferences();
      }
    }
  });
</script>

<div class="flex flex-col px-4 pt-6 gap-6 pb-12">
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
      <Button
        disabled={!allRequiredFilesPresent}
        on:click={() => {
          saveReferences();
        }}>{$_("save")}</Button
      >
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
            <ChevronDownOutline
              class="w-3 h-3 ms-2 text-fg-quaternary-light dark:text-fg-quaternary-dark"
            />
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
                  onclick={() => {
                    dropdownIndustryOpen = false;
                    selectedIndustry = industry.name;
                  }}
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
    <!-- section form -->
    <div class="flex gap-3">
      <div class="w-60">
        <Typography preset="Text md/Semibold" text={`${$_("documentAsBuyer")}`} color="secondary" />
      </div>
      <div class="flex flex-col grow gap-6">
        <div class="grid grid-cols-2 gap-4">
          <!-- form upload -->
          {#each documentTypes as docType}
            <UploadFile
              title={docType.title}
              handleFileUploadWrapper={() => handleFileUploadWrapper(docType.type, "buyer")}
              file={getFileByType(docType.type, "buyer")}
              actionButton={true}
              actionView={() => (
                (documentViewModal = true), (selectedFile = getFileByType(docType.type, "buyer"))
              )}
              actionDelete={() => deleteFileByType(docType.type, "buyer")}
            />
          {/each}
        </div>
      </div>
    </div>
    <Divider />
    <!-- section form -->
    <div class="flex gap-3">
      <div class="w-60">
        <Typography
          preset="Text md/Semibold"
          text={`${$_("documentAsSeller")}`}
          color="secondary"
        />
      </div>
      <div class="flex flex-col grow gap-6">
        <div class="grid grid-cols-2 gap-4">
          <!-- form upload -->
          {#each documentTypes as docType}
            <UploadFile
              title={docType.title}
              handleFileUploadWrapper={() => handleFileUploadWrapper(docType.type, "seller")}
              file={getFileByType(docType.type, "seller")}
              actionButton={true}
              actionView={() => (
                (documentViewModal = true), (selectedFile = getFileByType(docType.type, "seller"))
              )}
              actionDelete={() => deleteFileByType(docType.type, "seller")}
            />
          {/each}
        </div>
      </div>
    </div>
    <Divider />
    <!-- section form -->
    <div class="flex gap-3">
      <div class="w-60">
        <Typography
          preset="Text md/Semibold"
          text={`${$_("otherComplianceDocument")}`}
          color="secondary"
        />
      </div>
      <div class="flex flex-col grow gap-6">
        <div class="grid grid-cols-2 gap-4">
          <!-- form upload -->
          {#each otherDocType as docType}
            <UploadFile
              title={docType.title}
              handleFileUploadWrapper={() => handleFileUploadWrapper(docType.type, "other")}
              file={getFileByType(docType.type, "other")}
              actionButton={true}
              actionView={() => (
                (documentViewModal = true), (selectedFile = getFileByType(docType.type, "other"))
              )}
              actionDelete={() => deleteFileByType(docType.type, "other")}
            />
          {/each}
          <Button
            class="w-full min-h-32 flex gap-2.5 rounded-xl border border-primary-light dark:border-primary-dark bg-primary-light dark:bg-primary-dark"
            on:click={() => {addNewDocumentModal = true}}
            color="light"
          >
            <PlusCircleICon svgClass="w-6 h-6 text-fg-tertiary-light dark:text-fg-tertiary-dark" />
            {$_("addOtherComplianceDocs")}
          </Button>
        </div>
      </div>
    </div>
    <Divider />
    <div class="flex justify-end gap-3">
      <Button outline on:click={() => {}} color="light">{$_("cancel")}</Button>
      <Button disabled={files.length === 0} on:click={() => {}}>{$_("save")}</Button>
    </div>
  {/if}
</div>

<Modal title={$_("documentViewer")} bind:open={documentViewModal} autoclose outsideclose>
  {#if selectedFile && (selectedFile.name.includes(".pdf") || selectedFile.name.includes(".PDF"))}
    <div class="flex justify-center">
      <PdfViewer
        showBorder={false}
        showButtons={["navigation", "zoom", "rotate", "autoflip", "pageInfo"]}
        scale={1}
        data={atob(selectedFile.base64)}
      />
    </div>
  {/if}
</Modal>

<Modal title={$_("addNewDocument")} bind:open={addNewDocumentModal}>
  <div class="flex flex-col gap-4">
    <div class="flex flex-col gap-1.5">
      <Typography preset="Text sm/medium" color="secondary" text={`${$_("documentName")}*`} />
      <Input bind:value={tmpDocumentName} placeholder={$_("name")} />
    </div>
    <UploadFile
      title={`${$_("uploadDocuments")}*`}
      handleFileUploadWrapper={() => handleFileUploadWrapper(tmpDocumentName, "other")}
      file={tmpFile}
      actionButton={true}
      actionDelete={() => (tmpFile = [])}
    />
  </div>
  <svelte:fragment slot="footer">
    <Button on:click={addOtherDocument}>{$_("add")}</Button>
  </svelte:fragment>
</Modal>
