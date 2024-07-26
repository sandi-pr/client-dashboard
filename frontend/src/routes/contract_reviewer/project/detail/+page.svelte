<script>
  import { Typography, InputRw, UploadFile, triggerToast } from "$lib";
  import { Button, Modal, Input } from "flowbite-svelte";
  import { DatePicker } from "@svelte-plugins/datepicker";
  import { _ } from "svelte-i18n";
  import Divider from "$lib/components/Divider.svelte";
  import Stars03Icon from "$lib/icon/stars-03.svelte";
  import { ChevronDownOutline } from "flowbite-svelte-icons";
  import CalendarPlus02 from "$lib/icon/calendar-plus-02.svelte";
  import PlusCircleICon from "$lib/icon/plus-circle.svelte";
  import { formatDistanceToNow, format } from "date-fns";
  import PdfViewer from "svelte-pdf";
  import { onMount } from "svelte";
  import { CR3Call } from "hexacomponents-svelte";
  import { goto } from "$app/navigation";

  const params = new URLSearchParams(window.location.search);
  const projectId = params.get("projectId");

  const call = new CR3Call(globalThis.location.origin);

  let data = $state({});

  let documentViewModal = $state(false);
  let addNewDocumentModal = $state(false);
  let files = $state([]);
  let tmpFile = $state([]);
  let tmpDocumentName = $state("");
  let selectedFile = $state(null);

  let documentTypes = $state([
    {
      title: `${$_("documentGtc")}*`,
      type: "documentGtc",
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

      if (uploadedFile && role !== "other") {
        const base64Data = await readFileAsBase64(uploadedFile);
        let token = sessionStorage.getItem("token");
        if (token) {
          const response = await call.call("client-dashboard", "create-document-project", {
            URLQuery: {},
            URLParams: {},
            Params: {
              project_id: projectId,
              document: JSON.stringify({
                name: uploadedFile.name,
                size: formatBytes(uploadedFile.size),
                base64: base64Data,
                type: type,
              }),
            },
            Headers: {
              Authorization: "Bearer " + token,
            },
          });

          console.log(response);

          if (response.status !== 200) {
            console.log(response);
            return;
          }

          files.push({
            id: response.result.data.id,
            type: type,
            role: role,
            name: uploadedFile.name,
            size: formatBytes(uploadedFile.size),
            base64: base64Data,
          });
        }
      } else if (role === "other") {
        const base64Data = await readFileAsBase64(uploadedFile);
        tmpFile = {
          type: type,
          role: role,
          name: uploadedFile.name,
          size: formatBytes(uploadedFile.size),
          base64: base64Data,
        };
      } else {
        alert("Invalid file. Please select a valid file.");
      }
    });

    inputElement.click();
  };

  function getFileByType(type, role) {
    const index = files.findIndex((file) => file.type === type && file.role === role);
    return index !== -1 ? files[index] : [];
  }

  async function deleteFileByType(type, role) {
    const index = files.findIndex((file) => file.type === type && file.role === role);
    if (index !== -1) {
      let token = sessionStorage.getItem("token");
      if (token) {
        const response = await call.call("client-dashboard", "delete-document-project", {
          URLQuery: {},
          URLParams: {
            id: files[index].id.toString(),
          },
          Params: {},
          Headers: {
            Authorization: "Bearer " + token,
          },
        });
        console.log(response);

        if (response.status !== 200) {
          triggerToast("error");
          console.log(response);
          return;
        }

        triggerToast("deleted");
        files.splice(index, 1);
        if (role === "other") {
          otherDocType = otherDocType.filter((doc) => doc.type !== type);
        }
      }
    }
  }

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

  async function getProject() {
    let token = sessionStorage.getItem("token");
    if (token) {
      try {
        const response = await call.call("client-dashboard", "get-contract-project", {
          URLQuery: {},
          URLParams: {
            id: projectId,
          },
          Params: {},
          Headers: {
            Authorization: "Bearer " + token,
          },
        });

        console.log("response", response);
        if (response.status === 200) {
          if (response.result.data.length === 0) {
            triggerToast("error");
            setTimeout(async () => {
              goto(await call.getPage("client-dashboard", "contract-project-list"));
              return;
            }, 2000);
          }

          const response2 = await call.call("client-dashboard", "get-document-contract-project", {
            URLQuery: {},
            URLParams: {},
            Params: {
              project_id: projectId,
            },
            Headers: {
              Authorization: "Bearer " + token,
            },
          });

          if (response2.status === 200) {
            if (response2.result.files !== null) {
              response2.result.files.forEach((file) => {
                if (
                  file.type === "documentGtc" ||
                  file.type === "tradeConfirmation" ||
                  file.type === "dealDoneDocument"
                ) {
                  file.role = response.result.data[0].attributes.project_as;
                } else {
                  file.role = "other";
                }
              });
              files = response2.result.files;
              files.forEach((file) => {
                if (file.role === "other") {
                  otherDocType.push({
                    title: file.type,
                    type: file.type,
                  });
                }
              });
            }
          }

          data = {
            name: response.result.data[0].attributes.project_name,
            createdAt: formatDistanceToNow(new Date(response.result.data[0].attributes.createdAt), {
              addSuffix: true,
            }),
            industry: response.result.data[0].attributes.industry_type,
            as: response.result.data[0].attributes.project_as,
            start_date: format(
              new Date(response.result.data[0].attributes.project_start_date),
              "d MMM yyyy",
            ),
            end_date: format(
              new Date(response.result.data[0].attributes.project_end_date),
              "d MMM yyyy",
            ),
            description: response.result.data[0].attributes.description,
          };
        }
      } catch (error) {
        console.log(error);
      }
    }
  }

  onMount(async () => {
    if (projectId === null) {
      triggerToast("error");
      setTimeout(async () => {
        goto(await call.getPage("client-dashboard", "contract-project-list"));
      }, 2000);
    } else {
      console.log("projectId", projectId);
      getProject();
    }
  });
</script>

<div class="flex flex-col px-4 pt-6 gap-6 pb-12">
  <!-- section header -->
  <div class="flex justify-between">
    <div class="flex flex-col gap-2">
      <div class="flex flex-col">
        <Typography preset="Display sm/Semibold" text={data.name} color="primary" />
      </div>
      <Typography
        preset="Text md/Medium"
        text={`${$_("createOn")} - ${data.createdAt}`}
        color="quarterary"
      />
    </div>
    <Button class="gap-1 h-min" on:click={() => {}}>
      <Stars03Icon class="w-5 h-5 text-white" />
      {$_("AIReview")}
    </Button>
  </div>
  <Divider />
  <!-- section form -->
  <div class="flex gap-3">
    <div class="w-60">
      <Typography preset="Text md/Semibold" text={`${$_("industry")}*`} color="secondary" />
    </div>
    <div class="flex flex-col grow gap-3">
      <!-- form dropdown -->
      <div class="flex flex-col gap-1.5">
        <Typography preset="Text sm/medium" color="secondary" text={$_("projectTitle")} />
        <InputRw onEdit={false} bind:value={data.name} />
      </div>
      <div class="flex gap-3">
        <div class="flex grow flex-col gap-1.5">
          <Typography preset="Text sm/Medium" text={`${$_("industryType")}*`} color="secondary" />
          <div class="flex flex-col gap-1">
            <Button
              on:click={() => {}}
              class="flex justify-between border-primary-light dark:border-primary-dark"
              color="light"
            >
              {data.industry}
              <ChevronDownOutline
                class="w-3 h-3 ms-2 text-fg-quaternary-light dark:text-fg-quaternary-dark"
              />
            </Button>
          </div>
        </div>
        <div class="flex w-min flex-col gap-1.5">
          <Typography preset="Text sm/medium" color="secondary" text={$_("as")} />
          <Button
            on:click={() => {}}
            class="flex justify-between border-primary-light dark:border-primary-dark"
            color="light"
          >
            {data.as}
            <ChevronDownOutline
              class="w-3 h-3 ms-2 text-fg-quaternary-light dark:text-fg-quaternary-dark"
            />
          </Button>
        </div>
        <div class="flex grow flex-col gap-1.5">
          <Typography
            preset="Text sm/medium"
            color="secondary"
            text={`${$_("projectDuration")}*`}
          />
          <DatePicker
            isOpen={false}
            startDate={data.start_date}
            endDate={data.end_date}
            isRange
            isMultipane
            enableFutureDates
          >
            <button
              onclick={() => {}}
              class="flex w-full p-2 border rounded-lg border-gray-300 dark:border-gray-500 gap-2 justify-center items-center"
            >
              <svelte:component
                this={CalendarPlus02}
                class="w-2 h-2 ms-2 text-gray-700 dark:text-gray-400"
              />
              <div class="date">
                {data.start_date} - {data.end_date}
              </div>
            </button>
          </DatePicker>
        </div>
      </div>
    </div>
  </div>
  <Divider />
  <div class="flex gap-3">
    <div class="w-60">
      <Typography preset="Text md/Semibold" text={`${$_("documentUploads")}`} color="secondary" />
    </div>
    <div class="flex flex-col grow gap-6">
      <div class="grid grid-cols-2 gap-4">
        <!-- form upload -->
        {#each documentTypes as docType}
          <UploadFile
            title={docType.title}
            handleFileUploadWrapper={() => handleFileUploadWrapper(docType.type, data.as)}
            file={getFileByType(docType.type, data.as)}
            actionButton={true}
            actionView={() => (
              (documentViewModal = true), (selectedFile = getFileByType(docType.type, data.as))
            )}
            actionDelete={() => deleteFileByType(docType.type, data.as)}
          />
        {/each}
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
  </div>
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
