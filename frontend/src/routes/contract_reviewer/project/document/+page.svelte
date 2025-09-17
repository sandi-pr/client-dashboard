<script>
  import { onMount, onDestroy } from "svelte";
  import { page } from "$app/stores";
  import { goto } from "$app/navigation";
  import { _ } from "svelte-i18n";
  import { format } from "date-fns";
  import { Button } from "flowbite-svelte";
  import { DatePicker } from "@svelte-plugins/datepicker";

  import { call, getIcons, triggerToast, triggerLoading } from "$lib";
  import {
    Divider,
    InputRw,
    Dropdown,
    UploadFile,
    TextareaRw,
    Typography,
    PdfjsViewer,
    FeaturedIcon,
  } from "$lib";

  const isFrontend = process.env.NODE_ENV === "development";

  const FilePlusIcon = getIcons("filePlus");
  const CalendarPlus02Icon = getIcons("calendarPlus02");

  let projectId = $state(null);
  let projectDocumentId = $state(null);

  let onEdit = $state(true);
  let onView = $state(false);
  let isOpen = $state(false);
  let docTypeOpen = $state(false);
  let startDate = $state(null);
  let endDate = $state(null);
  let file = $state(null);
  let dateFormat = $state("d MMM yyyy");
  let formattedStartDate = $state();
  let formattedEndDate = $state("");
  let tmpDocumentType = $state("Documents");
  let tmpSelected = $state("Documents");
  let selectedTypeCollection = $state("Documents");
  let acceptedFileTypes = $state(".txt, .pdf, .docx");
  let fileUrl = $state("");
  let collectionId = $state(null);
  let projectName = $state("");
  let documentTypes = $state([
    {
      name: "Documents (TXT, PDF)",
      value: "Documents",
    },
    {
      name: "Images (JPG, JPEG, PNG, WEBP, IMG)",
      value: "Images",
    },
    {
      name: "Handwriting (JPG, JPEG, PNG)",
      value: "Handwriting",
    },
  ]);

  let data = $state({
    document_name: "",
    content_date: "",
    description: "",
  });

  const getDocument = async (id) => {
    let token = localStorage.getItem("token");
    if (token) {
      try {
        const response = await $call.call(
          "client-dashboard",
          "get-contract-project-with-document-id",
          {
            URLQuery: {},
            URLParams: {},
            Params: {
              id: parseInt(id),
            },
            Headers: {
              Authorization: "Bearer " + token,
            },
          },
        );

        if (response.status !== 200) {
          triggerToast("error");
          console.log(response);
          return;
        }

        file = {
          base64: response.result.base64,
        };

        data = {
          document_name: response.result.document.data[0].attributes.document_name,
          content_date: response.result.document.data[0].attributes.content_date,
          description: response.result.document.data[0].attributes.description,
          ext: response.result.document.data[0].attributes.document_media_id.data.attributes.ext,
        };

        startDate = new Date(response.result.document.data[0].attributes.content_date);
      } catch (error) {
        triggerToast("error");
        console.log(error);
      }
    }
  };

  const getCollections = async (id) => {
    let token = localStorage.getItem("token");
    if (token !== null) {
      try {
        const response = await $call.call("client-dashboard", "get-contract-project", {
          URLQuery: {},
          URLParams: {
            id: id.toString(),
          },
          Params: {},
          Headers: {
            Authorization: "Bearer " + token,
          },
        });

        if (response.status !== 200) {
          triggerToast("error");
          console.log(response);
          return;
        }

        projectName = response.result.data[0].attributes.project_name;

        let response2 = await $call.call("client-dashboard", "get-collection-by-name", {
          URLQuery: {},
          URLParams: {},
          Params: {
            collection_name: id + ". " + projectName + " (Project Documents)",
          },
          Headers: {
            Authorization: "Bearer " + token,
          },
        });
        if (response2.status !== 200) {
          triggerToast("error");
          console.log(response2);
          return;
        }

        if (response2.result.data.length === 0) {
          let response3 = await $call.call("client-dashboard", "create-document-collection", {
            URLQuery: {},
            URLParams: {},
            Params: {
              collection_name: id + ". " + projectName + " (Project Documents)",
              description: "Collection for " + projectName,
            },
            Headers: {
              Authorization: "Bearer " + token,
            },
          });

          if (response3.status !== 200) {
            triggerToast("error");
            console.log(response3);
            return;
          }

          collectionId = response3.result.data.id;
          return;
        }

        collectionId = response2.result.data[0].id;
      } catch (error) {
        console.log(error);
      }
    }
  };

  const createDocument = async () => {
    triggerLoading(true);

    let token = localStorage.getItem("token");
    if (token) {
      const response = await $call.call("client-dashboard", "create-document-project", {
        URLQuery: {},
        URLParams: {},
        Params: {
          project_id: projectId.toString(),
          content_date: format(formattedStartDate, "yyyy-MM-dd"),
          document: JSON.stringify(file),
        },
        Headers: {
          Authorization: "Bearer " + token,
        },
      });

      if (response.status !== 200) {
        triggerLoading(false);
        triggerToast("error");
        console.log(response);
        return;
      }

      // TODO: save to collection
      file.id = response.result.data.id;
      file.document_title = data.document_name;
      file.description = "Document project " + projectName;
      file.contract = "PROJECT";

      let response3 = await $call.call("client-dashboard", "create-document", {
        URLQuery: {},
        URLParams: {},
        Params: {
          documents: JSON.stringify([file]),
          collection_id: collectionId.toString(),
          document_type: tmpDocumentType,
        },
        Headers: {
          Authorization: "Bearer " + token,
        },
      });

      if (response3.status !== 200) {
        triggerLoading(false);
        triggerToast("error");
      }

      triggerToast("created");
      data = {
        document_name: "",
        content_date: "",
        description: "",
      };

      selectedTypeCollection = "Documents";
      acceptedFileTypes = ".txt, .pdf";
      file = null;

      setTimeout(async () => {
        triggerLoading(false);
        const gotoUrl = isFrontend
          ? "/contract_reviewer/project/detail"
          : await $call.getPage("client-dashboard", "contract-project-detail");
        goto(gotoUrl + "?projectId=" + projectId);
      }, 2000);
    }
  };

  const readFileAsBase64 = async (file) => {
    return new Promise((resolve, reject) => {
      const reader = new FileReader();

      reader.onload = () => {
        const base64Data = reader.result.split(",")[1];
        resolve(base64Data);
      };
      reader.onerror = reject;

      reader.readAsDataURL(file);
    });
  };

  const selectTypeCollection = (type) => {
    tmpDocumentType = type;

    if (type === "Documents") {
      selectedTypeCollection = "Documents";
      acceptedFileTypes = ".txt, .pdf, .docx";
    } else if (type === "Images") {
      selectedTypeCollection = "Images";
      acceptedFileTypes = ".jpg, .jpeg, .png, .webp, .img";
    } else if (type === "Handwriting") {
      selectedTypeCollection = "Handwriting";
      acceptedFileTypes = ".jpg, .jpeg, .png";
    }

    onView = false;
    file = null;
    fileUrl = "";

    if (tmpSelected !== selectedTypeCollection) {
      file = null;
      tmpSelected = selectedTypeCollection;
    }
  };

  const formatBytes = (bytes, decimals = 2) => {
    if (bytes === 0) return "0 Bytes";

    const k = 1024;
    const dm = decimals < 0 ? 0 : decimals;
    const sizes = ["Bytes", "KB", "MB", "GB", "TB"];

    const i = Math.floor(Math.log(bytes) / Math.log(k));

    return parseFloat((bytes / Math.pow(k, i)).toFixed(dm)) + " " + sizes[i];
  };

  const handleFileUploadWrapper = () => {
    const inputElement = document.createElement("input");
    inputElement.type = "file";
    inputElement.accept = acceptedFileTypes;

    inputElement.addEventListener("change", async (event) => {
      const uploadedFile = event.target.files[0];

      if (uploadedFile) {
        const base64Data = await readFileAsBase64(uploadedFile);
        file = {
          type: projectDocumentId ? data.document_name : uploadedFile.type,
          role: tmpDocumentType,
          name: uploadedFile.name,
          size: formatBytes(uploadedFile.size),
          base64: base64Data,
          ext: "." + uploadedFile.name.split(".").pop(),
        };

        if (uploadedFile.type.startsWith("image/")) {
          fileUrl = URL.createObjectURL(uploadedFile);
          // console.log("fileUrl", fileUrl);
        }
        // console.log("file", file);
        onView = true;
      } else {
        alert("Invalid file. Please select a valid file.");
      }
    });

    inputElement.click();
  };

  const toggleDatePicker = () => {
    // if (onView) return;
    isOpen = !isOpen;
  };

  const formatDate = (dateString) => (dateString && format(new Date(dateString), dateFormat)) || "";

  onMount(() => {
    formattedStartDate = formatDate(new Date());
    const unsubscribe = page.subscribe(async ($page) => {
      projectId = $page.url.searchParams.get("projectId");
      projectDocumentId = $page.url.searchParams.get("projectDocumentId");
      if (projectId !== null) {
        if (projectDocumentId !== null) {
          onEdit = false;
          onView = true;
          getDocument(projectDocumentId);
        } else {
          onEdit = true;
          onView = false;
          getCollections(projectId);
        }
      } else {
        triggerToast("error");
        goto(
          isFrontend
            ? "/contract_reviewer/project"
            : await $call.getPage("client-dashboard", "contract-project-list"),
        );
      }
    });
    onDestroy(() => unsubscribe());
  });

  $effect(() => {
    formattedStartDate = formatDate(startDate);
    formattedEndDate = formatDate(endDate);
  });
</script>

<div class="flex flex-col px-4 pt-6 gap-6 pb-12">
  <!-- section header -->
  <div class="flex justify-between">
    <div class="flex items-center gap-3">
      <FeaturedIcon classesSvg={FilePlusIcon} />
      <div class="flex flex-col">
        {#if data?.document_name !== ""}
          <Typography preset="Display sm/Semibold" text={data.document_name} color="primary" />
        {:else}
          <Typography preset="Text lg/Semibold" text={$_("addNewDocument")} color="primary" />
          <Typography preset="Text sm/Regular" text={$_("descAddNewDocument")} color="tertiary" />
        {/if}
      </div>
    </div>
    <div class="flex gap-3">
      <Button
        outline
        color="light"
        on:click={async () => {
          const gotoUrl = isFrontend
            ? "/contract_reviewer/project/detail"
            : await $call.getPage("client-dashboard", "contract-project-detail");
          goto(gotoUrl + "?projectId=" + projectId);
        }}
      >
        <Typography preset="Text sm/Medium" text={$_("cancel")} />
      </Button>
      <Button
        disabled={file === null}
        on:click={() => {
          createDocument();
        }}
      >
        {#if projectDocumentId !== null}
          {$_("save")}
        {:else}
          {$_("add")}
        {/if}
      </Button>
    </div>
  </div>
  <Divider />
  <!-- section form -->
  <div class="flex gap-3 w-full">
    <div class="w-60">
      <Typography preset="Text md/Semibold" text={$_("documentInformation")} color="secondary" />
    </div>
    <div class="flex grow flex-col gap-3">
      <div class="flex gap-4">
        <div class="flex flex-col grow gap-1.5">
          <Typography preset="Text sm/Medium" text={$_("documentName") + "*"} color="secondary" />
          <InputRw {onEdit} bind:value={data.document_name} />
        </div>
        <div class="flex grow flex-col gap-1.5">
          <Typography preset="Text sm/medium" color="secondary" text={`${$_("contentDate")}*`} />
          <DatePicker bind:isOpen bind:startDate enableFutureDates>
            <button
              onclick={toggleDatePicker}
              class:open={isOpen}
              class="flex w-full p-2.5 border rounded-lg border-gray-300 dark:border-gray-500 gap-2 justify-center items-center"
            >
              <CalendarPlus02Icon
                svgClass="w-4 h-4 ms-2 text-secondary-light dark:text-secondary-dark"
              />
              <div class="date text-secondary-light dark:text-secondary-dark">
                {#if startDate}
                  {formattedStartDate}
                {:else}
                  {$_("selectDate")}
                {/if}
              </div>
            </button>
          </DatePicker>
        </div>
      </div>
      <div class="flex flex-col grow gap-1.5">
        <Typography preset="Text sm/Medium" text={$_("description")} color="secondary" />
        <TextareaRw {onEdit} bind:value={data.description} />
      </div>
      {#if !onView}
        <UploadFile
          title={`${$_("documentUpload")}*`}
          handleFileUploadWrapper={() => handleFileUploadWrapper()}
          {file}
          actionButton={true}
          actionDelete={() => (file = null)}
          isFlex
        />
      {/if}
    </div>
  </div>
  <Divider />
  <!-- section form -->
  <div class="flex gap-3">
    <div class="w-60">
      <Typography preset="Text md/Semibold" text={`${$_("documentType")}`} color="secondary" />
    </div>
    <div class="flex grow flex-col gap-3">
      <div class="flex flex-col grow gap-1.5">
        <Typography preset="Text sm/Medium" text={$_("documentType")} color="secondary" />
        <Dropdown
          dropdownOpen={docTypeOpen}
          bind:selectedItem={selectedTypeCollection}
          items={documentTypes}
          placeholder={$_("selectDocumentType")}
          isOverlay={true}
          onClick={selectTypeCollection}
        />
      </div>
    </div>
  </div>
  <Divider />
  <!-- {#if !onView}
    <div class="flex justify-end gap-3">
      <Button
        outline
        on:click={async () => {
          goto(
            (await $call.getPage("client-dashboard", "contract-project-detail")) +
              "?projectId=" +
              projectId,
          );
        }}
        color="light">{$_("cancel")}</Button
      >
      <Button
        disabled={file === null}
        on:click={() => {
          createDocument();
        }}>{$_("add")}</Button
      >
    </div>
  {/if} -->
  {#if onView && (file?.base64 !== undefined || fileUrl !== "")}
    <div class="flex gap-6">
      <div class="w-60">
        <Typography preset="Text md/Medium" text={$_("viewer")} color="secondary" />
      </div>
      <div class="flex flex-col grow gap-6 rounded-lg bg-secondary-light dark:bg-secondary-dark">
        {#if data.ext === ".pdf" || file?.ext === ".pdf"}
          <div class="flex justify-center">
            <PdfjsViewer pdfFile={file.base64} />
          </div>
        {:else if data.ext === ".jpg" || data.ext === ".png" || data.ext === ".jpeg" || data.ext === ".webp" || file?.ext === ".jpg" || file?.ext === ".png" || file?.ext === ".jpeg" || file?.ext === ".webp" || file?.ext === ".img"}
          <div class="flex justify-center">
            <img class="w-full h-screen" src={`${fileUrl}`} alt="imageViewer" />
          </div>
          <!-- {:else}
          {#each contentFile.split("\n") as paragraphs, i}
            <p class="dark:text-white">{paragraphs}</p>
          {/each} -->
        {/if}
      </div>
    </div>
  {/if}
</div>
