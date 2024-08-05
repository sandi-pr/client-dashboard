<script>
  import { Button, Modal, Input } from "flowbite-svelte";
  import { FeaturedIcon, Typography, UploadFile, triggerToast, SelectedCard, Dropdown } from "$lib";
  import PdfViewer from "svelte-pdf";
  import { _ } from "svelte-i18n";
  import Divider from "$lib/components/Divider.svelte";
  import FilePlusIcon from "$lib/icon/file-plus-02.svelte";
  import { ChevronDownOutline } from "flowbite-svelte-icons";
  import PlusCircleICon from "$lib/icon/plus-circle.svelte";
  import { CR3Call } from "hexacomponents-svelte";
  import { onMount } from "svelte";

  const call = new CR3Call(globalThis.location.origin);

  let dropdownIndustryOpen = $state(false);
  let documentViewModal = $state(false);
  let addNewDocumentModal = $state(false);
  let addSellerDocumentModal = $state(false);
  let addBuyerDocumentModal = $state(false);
  let tmpDocumentName = $state("");
  let tmpFile = $state([]);
  let selectedIndustry = $state("");
  let files = $state([]);
  let selectedFile = $state(null);
  let requiredFileTypes = $state(["documentGtc"]);
  let roles = $state(["Buyer", "Seller"]);
  let idUploaded = $state([]);
  let tmpIdDeleted = $state([]);

  let industries = $state([]);

  let documentTypes = $state([
    {
      title: `${$_("documentGtc")}*`,
      type: "documentGtc",
    },
    // {
    //   title: $_("gtcAmendment"),
    //   type: "gtcAmendment",
    // },
    // {
    //   title: `${$_("tradeConfirmation")}*`,
    //   type: "tradeConfirmation",
    // },
    // {
    //   title: `${$_("dealDoneDocument")}*`,
    //   type: "dealDoneDocument",
    // },
  ]);

  let sellerDocuments = $state([
    {
      title: `${$_("documentGtc")}*`,
      type: "documentGtc",
    },
  ]);

  let buyerDocuments = $state([
    {
      title: `${$_("documentGtc")}*`,
      type: "documentGtc",
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

      if (uploadedFile && type === "documentGtc") {
        const base64Data = await readFileAsBase64(uploadedFile);
        files.push({
          id: 0,
          type: type,
          role: role,
          name: uploadedFile.name,
          size: formatBytes(uploadedFile.size),
          base64: base64Data,
        });
      } else if (uploadedFile && type !== "documentGtc") {
        const base64Data = await readFileAsBase64(uploadedFile);
        tmpFile = {
          id: 0,
          type: type,
          role: role,
          name: uploadedFile.name,
          size: formatBytes(uploadedFile.size),
          base64: base64Data,
        };
        console.log(tmpFile);
      } else {
        alert("Invalid file. Please select a valid file.");
      }
    });

    inputElement.click();
  };

  function addOtherDocument(role) {
    if (role === "Seller") {
      sellerDocuments.push({
        title: tmpDocumentName,
        type: tmpDocumentName,
      });

      tmpFile.type = tmpDocumentName;

      files.push(tmpFile);

      tmpFile = [];
      tmpDocumentName = "";
      addSellerDocumentModal = false;
    } else if (role === "Buyer") {
      buyerDocuments.push({
        title: tmpDocumentName,
        type: tmpDocumentName,
      });

      tmpFile.type = tmpDocumentName;

      files.push(tmpFile);

      tmpFile = [];
      tmpDocumentName = "";
      addBuyerDocumentModal = false;
    }
  }

  function getFileByType(type, role) {
    const index = files.findIndex((file) => file.type === type && file.role === role);
    return index !== -1 ? files[index] : [];
  }

  function deleteFileByType(type, role) {
    const index = files.findIndex((file) => file.type === type && file.role === role);

    if (index !== -1) {
      if (files[index].id !== 0) tmpIdDeleted.push(files[index].id);
      files.splice(index, 1);
      if (role === "Seller" && type !== "documentGtc") {
        sellerDocuments = sellerDocuments.filter((doc) => doc.type !== type);
      } else if (role === "Buyer" && type !== "documentGtc") {
        buyerDocuments = buyerDocuments.filter((doc) => doc.type !== type);
      }
    }
  }

  function areRequiredFilesPresent(files) {
    // return roles.every((role) =>
    //   requiredFileTypes.every((type) =>
    //     files.some((file) => file.type === type && file.role === role),
    //   ),
    // );

    return requiredFileTypes.every((type) => files.some((file) => file.type === type));
  }

  let allRequiredFilesPresent = $state(false);

  $effect(() => {
    allRequiredFilesPresent = areRequiredFilesPresent(files);
  });

  async function saveReferences() {
    if (allRequiredFilesPresent) {
      let token = sessionStorage.getItem("token");
      try {
        if (tmpIdDeleted.length > 0) {
          for (let i = 0; i < tmpIdDeleted.length; i++) {
            if (tmpIdDeleted[i] !== 0) {
              const response = await call.call("client-dashboard", "delete-references-document", {
                URLQuery: {},
                URLParams: {
                  id: tmpIdDeleted[i].toString(),
                },
                Params: {},
                Headers: {
                  Authorization: "Bearer " + token,
                },
              });

              if (response.status !== 200) {
                triggerToast("error");
                break;
              }

              tmpIdDeleted.filter((id) => id !== tmpIdDeleted[i]);
              files = files.filter((file) => file.id !== tmpIdDeleted[i]);
            }
          }
          tmpIdDeleted = [];
        }

        for (let i = 0; i < files.length; i++) {
          if (files[i].id === 0) {
            const response2 = await call.call("client-dashboard", "create-references-document", {
              URLQuery: {},
              URLParams: {},
              Params: {
                references: JSON.stringify([files[i]]),
                industry_id: selectedIndustry,
              },
              Headers: {
                Authorization: "Bearer " + token,
              },
            });
            console.log(response2);

            if (response2.status !== 200) {
              triggerToast("error");
              break;
            }

            files[i].id = response2.result.data[0].id;
          }
        }
        triggerToast("updated");
      } catch (error) {
        console.log(error);
      }
    }
  }

  async function getReferences() {
    let token = sessionStorage.getItem("token");
    try {
      let page = 1;
      idUploaded = [];
      files = [];

      while (true) {
        const response = await call.call("client-dashboard", "get-references-document", {
          URLQuery: {},
          URLParams: {},
          Params: {
            industry_id: selectedIndustry,
            page: page,
            page_size: 25,
          },
          Headers: {
            Authorization: "Bearer " + token,
          },
        });

        if (response.status !== 200) {
          console.log(response);
          return;
        }

        if (response.result.files.length === 0) {
          break;
        }

        files = files.concat(response.result.files);
        page++;
      }

      if (files.length > 0) {
        files.forEach((file) => {
          idUploaded.push(file.id);
          if (file.role === "Seller" && file.type !== "documentGtc") {
            sellerDocuments.push({
              title: file.type,
              type: file.type,
            });
          } else if (file.role === "Buyer" && file.type !== "documentGtc") {
            buyerDocuments.push({
              title: file.type,
              type: file.type,
            });
          }
        });
      }
    } catch (error) {
      console.log(error);
    }
  }

  async function getIndustries() {
    let token = sessionStorage.getItem("token");

    if (token !== null) {
      let page = 1;
      let tmpIndustries = [];
      while (true) {
        const response = await call.call("client-dashboard", "get-industries", {
          URLQuery: {},
          URLParams: {},
          Params: {
            industry_id: selectedIndustry,
            page: page,
            page_size: 25,
          },
          Headers: {
            Authorization: "Bearer " + token,
          },
        });

        if (response.status === 200) {
          if (response.result.data.length === 0) {
            break;
          } else {
            tmpIndustries = tmpIndustries.concat(response.result.data);
            page++;
          }
        }
      }

      let promise = tmpIndustries.map((industry) => {
        return {
          value: industry.id,
          name: industry.attributes.industry_name,
          description: industry.attributes.industry_description,
        };
      });

      industries = await Promise.all(promise);
    }
  }

  onMount(() => {
    getIndustries();
  });

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
  <div class="flex gap-3 w-full">
    <div class="w-60">
      <Typography preset="Text md/Semibold" text={`${$_("industry")}*`} color="secondary" />
    </div>
    <div class="flex grow gap-6">
      <!-- form dropdown -->
      <div class="relative flex flex-col gap-1.5 w-full">
        <Typography preset="Text sm/Medium" text={`${$_("industryType")}*`} color="secondary" />
        <!-- <SelectedCard selectionData={industries} numGrid={3} /> -->
        <Dropdown
          dropdownOpen={dropdownIndustryOpen}
          bind:selectedItem={selectedIndustry}
          items={industries}
          placeholder={$_("selectIndustry")}
          isOverlay={true}
        />
      </div>
    </div>
  </div>
  {#if selectedIndustry !== ""}
    <Divider />
    <!-- section form -->
    <div class="flex gap-3">
      <div class="w-60">
        <Typography preset="Text md/Semibold" text={`${$_("documentAsBuyer")}`} color="secondary" />
      </div>
      <div class="flex flex-col grow gap-6">
        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
          <!-- form upload -->
          {#each buyerDocuments as docType}
            <UploadFile
              title={docType.title}
              handleFileUploadWrapper={() => handleFileUploadWrapper(docType.type, "Buyer")}
              file={getFileByType(docType.type, "Buyer")}
              actionButton={true}
              actionView={() => (
                (documentViewModal = true), (selectedFile = getFileByType(docType.type, "Buyer"))
              )}
              actionDelete={() => deleteFileByType(docType.type, "Buyer")}
            />
          {/each}
          <Button
            class="w-full min-h-32 md:mt-6 flex gap-2.5 rounded-xl border border-primary-light dark:border-primary-dark bg-primary-light dark:bg-primary-dark"
            on:click={() => {
              addBuyerDocumentModal = true;
            }}
            color="light"
          >
            <PlusCircleICon svgClass="w-6 h-6 text-fg-tertiary-light dark:text-fg-tertiary-dark" />
            {$_("addOtherComplianceDocs")}
          </Button>
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
        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
          <!-- form upload -->
          {#each sellerDocuments as docType}
            <UploadFile
              title={docType.title}
              handleFileUploadWrapper={() => handleFileUploadWrapper(docType.type, "Seller")}
              file={getFileByType(docType.type, "Seller")}
              actionButton={true}
              actionView={() => (
                (documentViewModal = true), (selectedFile = getFileByType(docType.type, "Seller"))
              )}
              actionDelete={() => deleteFileByType(docType.type, "Seller")}
            />
          {/each}
          <Button
            class="w-full min-h-32 md:mt-6 flex gap-2.5 rounded-xl border border-primary-light dark:border-primary-dark bg-primary-light dark:bg-primary-dark"
            on:click={() => {
              addSellerDocumentModal = true;
            }}
            color="light"
          >
            <PlusCircleICon svgClass="w-6 h-6 text-fg-tertiary-light dark:text-fg-tertiary-dark" />
            {$_("addOtherComplianceDocs")}
          </Button>
        </div>
      </div>
    </div>
    <Divider />
    <!-- section form -->
    <!-- <div class="flex gap-3">
      <div class="w-60">
        <Typography
          preset="Text md/Semibold"
          text={`${$_("otherComplianceDocument")}`}
          color="secondary"
        />
      </div>
      <div class="flex flex-col grow gap-6">
        <div class="grid grid-cols-2 gap-4">
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
            on:click={() => {
              addNewDocumentModal = true;
            }}
            color="light"
          >
            <PlusCircleICon svgClass="w-6 h-6 text-fg-tertiary-light dark:text-fg-tertiary-dark" />
            {$_("addOtherComplianceDocs")}
          </Button>
        </div>
      </div>
    </div> -->
    <!-- <Divider />
    <div class="flex justify-end gap-3">
      <Button outline on:click={() => {}} color="light">{$_("cancel")}</Button>
      <Button disabled={files.length === 0} on:click={() => {}}>{$_("save")}</Button>
    </div> -->
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

<Modal title={$_("addNewDocumentAsBuyer")} bind:open={addBuyerDocumentModal}>
  <div class="flex flex-col gap-4">
    <div class="flex flex-col gap-1.5">
      <Typography preset="Text sm/medium" color="secondary" text={`${$_("documentName")}*`} />
      <Input bind:value={tmpDocumentName} placeholder={$_("name")} />
    </div>
    <UploadFile
      title={`${$_("uploadDocuments")}*`}
      handleFileUploadWrapper={() => handleFileUploadWrapper(tmpDocumentName, "Buyer")}
      file={tmpFile}
      actionButton={true}
      actionDelete={() => (tmpFile = [])}
    />
  </div>
  <svelte:fragment slot="footer">
    <Button
      on:click={() => {
        addOtherDocument("Buyer");
      }}
      disabled={tmpFile.length === 0 || !tmpDocumentName}>{$_("add")}</Button
    >
  </svelte:fragment>
</Modal>

<Modal title={$_("addNewDocumentAsSeller")} bind:open={addSellerDocumentModal}>
  <div class="flex flex-col gap-4">
    <div class="flex flex-col gap-1.5">
      <Typography preset="Text sm/medium" color="secondary" text={`${$_("documentName")}*`} />
      <Input bind:value={tmpDocumentName} placeholder={$_("name")} />
    </div>
    <UploadFile
      title={`${$_("uploadDocuments")}*`}
      handleFileUploadWrapper={() => handleFileUploadWrapper(tmpDocumentName, "Seller")}
      file={tmpFile}
      actionButton={true}
      actionDelete={() => (tmpFile = [])}
    />
  </div>
  <svelte:fragment slot="footer">
    <Button
      on:click={() => {
        addOtherDocument("Seller");
      }}
      disabled={!tmpFile && !tmpDocumentName}>{$_("add")}</Button
    >
  </svelte:fragment>
</Modal>
