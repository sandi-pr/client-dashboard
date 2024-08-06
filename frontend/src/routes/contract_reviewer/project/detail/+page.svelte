<script>
  import {
    Typography,
    InputRw,
    UploadFile,
    triggerToast,
    Dropdown,
    DocumentProject,
    triggerLoading,
  } from "$lib";
  import { Button, Modal, Input } from "flowbite-svelte";
  import { DatePicker } from "@svelte-plugins/datepicker";
  import { _ } from "svelte-i18n";
  import Divider from "$lib/components/Divider.svelte";
  import Stars03Icon from "$lib/icon/stars-03.svelte";
  import * as Icon from "flowbite-svelte-icons";
  import CalendarPlus02 from "$lib/icon/calendar-plus-02.svelte";
  import PlusCircleICon from "$lib/icon/plus-circle.svelte";
  import { formatDistanceToNow, format } from "date-fns";
  import PdfViewer from "svelte-pdf";
  import { onMount, onDestroy } from "svelte";
  import { CR3Call } from "hexacomponents-svelte";
  import { goto } from "$app/navigation";
  import { page } from "$app/stores";

  const call = new CR3Call(globalThis.location.origin);

  let projectId = $state(null);

  let data = $state({});

  let documentViewModal = $state(false);
  let addNewDocumentModal = $state(false);
  let deleteModal = $state(false);
  let ready = $state(false);
  let files = $state([]);
  let tmpFile = $state([]);
  let typeTmp = $state("");
  let asTmp = $state("");
  let rulesTmp = $state([]);
  let tmpDocumentName = $state("");
  let selectedFile = $state(null);
  let init = $state(true);

  let industries = $state([]);
  let role = $state([
    { name: "Buyer", value: "Buyer" },
    { name: "Seller", value: "Seller" },
  ]);

  let selectedIndustry = $state("");

  let documentUploads = $state([
    // {
    //   title: "Trade Confirmation",
    //   type: "tradeConfirmation",
    // },
  ]);

  let rules = $state([
    // {
    //   name: "Patra Niaga POC",
    // },
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
        tmpFile = {
          type: type,
          role: role,
          name: uploadedFile.name,
          size: formatBytes(uploadedFile.size),
          base64: base64Data,
        };
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

        if (response.status !== 200) {
          triggerToast("error");
          console.log(response);
          return;
        }

        documentUploads = documentUploads.filter((doc) => doc.type !== type);
        triggerToast("deleted");
        files.splice(index, 1);
        if (role === "other") {
          otherDocType = otherDocType.filter((doc) => doc.type !== type);
        }
      }
    }
  }

  async function addOtherDocument() {
    addNewDocumentModal = false;
    triggerLoading(true);
    if (tmpDocumentName) {
      documentUploads.push({
        title: tmpDocumentName,
        type: tmpDocumentName,
      });

      let token = sessionStorage.getItem("token");
      if (token) {
        const response = await call.call("client-dashboard", "create-document-project", {
          URLQuery: {},
          URLParams: {},
          Params: {
            project_id: projectId.toString(),
            document: JSON.stringify(tmpFile),
          },
          Headers: {
            Authorization: "Bearer " + token,
          },
        });

        if (response.status !== 200) {
          triggerLoading(false);
          triggerToast("error");
          addNewDocumentModal = true;
          console.log(response);
          return;
        }

        const response2 = await call.call("hexa-ai-contract-review", "storing-to-rag", {
          URLQuery: {},
          URLParams: {},
          Params: {
            document_base64: tmpFile.base64,
            document_id: response.result.data.id.toString(),
            document_name: tmpFile.name,
            document_type: "project",
          },
          Headers: {
            Authorization: "Bearer " + token,
          },
        });

        if (response2.status !== 200) {
          triggerLoading(false);
          triggerToast("error");
          addNewDocumentModal = true;
          console.log(response2);
          return;
        }

        const response3 = await call.call("client-dashboard", "update-project-document", {
          URLQuery: {},
          URLParams: {
            id: response.result.data.id.toString(),
          },
          Params: {
            rag_id: response2.result.rag_id,
          },
          Headers: {
            Authorization: "Bearer " + token,
          },
        });

        if (response3.status !== 200) {
          triggerLoading(false);
          triggerToast("error");
          addNewDocumentModal = true;
          console.log(response3);
          return;
        }

        triggerLoading(false);
        files.push(tmpFile);
        getProject();
        triggerToast("created");
      }

      tmpFile = [];
      tmpDocumentName = "";
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

  async function getRules() {
    let token = sessionStorage.getItem("token");
    if (token) {
      rulesTmp = [];
      let page = 1;
      try {
        while (true) {
          const response = await call.call("client-dashboard", "get-contract-rules", {
            URLQuery: {},
            URLParams: {},
            Params: {
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

          if (response.result.data.length === 0) {
            break;
          }
          rulesTmp = rulesTmp.concat(response.result.data);
          page++;
        }

        if (rulesTmp.length > 0) {
          let promise = rulesTmp.map((rule) => {
            return {
              value: rule.id,
              name: rule.attributes.rule_name,
            };
          });
          rules = await Promise.all(promise);
        }
      } catch (error) {
        console.log(error);
      }
    }
  }

  async function getProject() {
    if (init) triggerLoading(true);

    let token = sessionStorage.getItem("token");
    if (token) {
      try {
        const response = await call.call("client-dashboard", "get-contract-project", {
          URLQuery: {},
          URLParams: {
            id: projectId.toString(),
          },
          Params: {},
          Headers: {
            Authorization: "Bearer " + token,
          },
        });

        if (response.status !== 200) {
          if (init) triggerLoading(false);
          triggerToast("error");
          console.log(response);
          return;
        }

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
              project_id: projectId.toString(),
            },
            Headers: {
              Authorization: "Bearer " + token,
            },
          });

          if (response2.status !== 200) {
            if (init) triggerLoading(false);
            triggerToast("error");
            console.log(response2);
            return;
          }

          if (response2.status === 200) {
            files = [];
            documentUploads = [];
            if (response2.result.files !== null) {
              response2.result.files.forEach((file) => {
                file.role = response.result.data[0].attributes.project_as;
              });
              files = response2.result.files;
              files.forEach((file) => {
                documentUploads.push({
                  title: file.type,
                  type: file.type,
                });
              });
            }
          }

          selectedIndustry = response.result.data[0].attributes.industry_id;

          await getIndustries();
          await getRules();

          data = {
            name: response.result.data[0].attributes.project_name,
            createdAt: formatDistanceToNow(new Date(response.result.data[0].attributes.createdAt), {
              addSuffix: true,
            }),
            industry: response.result.data[0].attributes.industry_id,
            contract_rule_id:
              response.result.data[0].attributes.contract_rule_id !== null
                ? response.result.data[0].attributes.contract_rule_id
                : rules[0].value,
            as: response.result.data[0].attributes.project_as,
            start_date: format(
              new Date(response.result.data[0].attributes.duration_start),
              "d MMM yyyy",
            ),
            end_date: format(
              new Date(response.result.data[0].attributes.duration_end),
              "d MMM yyyy",
            ),
            description: response.result.data[0].attributes.description,
          };

          ready = true;

          if (init) triggerLoading(false);
          getStatusGenerated();
        }
      } catch (error) {
        triggerLoading(false);
        triggerToast("error");
        console.log(error);
      }
    }
  }

  async function generatedAI(file) {
    triggerLoading(true);
    let token = sessionStorage.getItem("token");
    if (token) {
      try {
        const response2 = await call.call("client-dashboard", "get-document-contract-project", {
          URLQuery: {},
          URLParams: {},
          Params: {
            project_id: projectId.toString(),
          },
          Headers: {
            Authorization: "Bearer " + token,
          },
        });

        if (response2.status !== 200) {
          triggerLoading(false);
          triggerToast("error");
          console.log(response2);
          return;
        }

        if (response2.status === 200) {
          files = [];
          documentUploads = [];
          let generated = [];
          if (response2.result.files !== null) {
            response2.result.files.forEach((file) => {
              file.role = response.result.data[0].attributes.project_as;
            });
            files = response2.result.files;
            files.forEach((file) => {
              documentUploads.push({
                title: file.type,
                type: file.type,
              });
            });

            for (let i = 0; i < files.length; i++) {
              if (files[i].status === "generated") {
                generated.push(files[i].id);
                continue;
              } else if (files[i].status === "processing") {
                break;
              } else {
                generated.push(files[i].id);
                continue;
              }
            }

            if (generated.length < files.length) {
              setTimeout(getStatusGenerated, 2000);
            } else {
              console.log("finish");
            }
          }
        }
      } catch (error) {
        console.log(error);
      }
    }

    return;
  }

  async function getStatusGenerated() {
    let token = sessionStorage.getItem("token");
    if (token) {
      try {
        let generated = [];
        for (let i = 0; i < files.length; i++) {
          if (files[i].status === "generated") {
            generated.push(files[i].id);
            continue;
          } else {
            init = false;
            getProject();
            break;
          }
        }

        if (generated.length < files.length) {
          setTimeout(getStatusGenerated, 2000);
        } else {
          console.log("finish");
        }
      } catch (error) {
        console.log(error);
      }
    }
  }

  async function goResultPage(id) {
    let token = sessionStorage.getItem("token");
    if (token) {
      try {
        const response = await call.call(
          "client-dashboard",
          "get-contract-project-result-with-document-id",
          {
            URLQuery: {},
            URLParams: {},
            Params: {
              project_document_id: id.toString(),
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

        goto(
          (await call.getPage("client-dashboard", "contract-project-review")) +
            "?projectDocumentId=" +
            response.result.data[0].id,
        );
      } catch (error) {
        console.log(error);
      }
    }
  }

  onMount(() => {
    const unsubscribe = page.subscribe(async ($page) => {
      projectId = $page.url.searchParams.get("projectId");
      if (projectId !== null) {
        getProject();
      } else {
        triggerToast("error");
        goto(await call.getPage("client-dashboard", "contract-project-list"));
      }
    });
    return () => unsubscribe();
  });

  $effect(async () => {
    const unsubscribe = page.subscribe(async ($page) => {
      projectId = $page.url.searchParams.get("projectId");
      if (projectId !== null) {
        getProject();
      } else {
        goto(await call.getPage("client-dashboard", "contract-project-list"));
      }
    });
    onDestroy(() => unsubscribe());
  });
</script>

{#if ready}
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
    </div>
    <Divider />
    <!-- section form -->
    <div class="flex gap-3">
      <div class="w-60">
        <Typography preset="Text md/Semibold" text={`${$_("industry")}`} color="secondary" />
      </div>
      <div class="flex flex-col grow gap-3">
        <!-- form dropdown -->
        <div class="flex flex-col gap-1.5">
          <Typography preset="Text sm/medium" color="secondary" text={$_("projectTitle")} />
          <InputRw onEdit={false} bind:value={data.name} />
        </div>
        <div class="flex gap-3">
          <div class="flex grow flex-col gap-1.5">
            <Typography preset="Text sm/Medium" text={`${$_("industryType")}`} color="secondary" />
            <Dropdown
              dropdownOpen={false}
              bind:selectedItem={data.industry}
              items={industries}
              isOverlay={true}
              disabled={true}
            />
          </div>
          <!-- <div class="flex w-max flex-col gap-1.5">
            <Typography preset="Text sm/medium" color="secondary" text={$_("as")} />
            <Dropdown
              dropdownOpen={false}
              bind:selectedItem={data.as}
              items={role}
              isOverlay={true}
              disabled={true}
            />
          </div> -->
          <div class="flex grow flex-col gap-1.5">
            <Typography preset="Text sm/medium" color="secondary" text={$_("rules")} />
            <Dropdown
              dropdownOpen={false}
              bind:selectedItem={data.contract_rule_id}
              items={rules}
              isOverlay={true}
              disabled={true}
            />
          </div>
        </div>
        <div class="flex grow flex-col gap-1.5">
          <Typography preset="Text sm/medium" color="secondary" text={`${$_("projectDuration")}`} />
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
    <Divider />
    <div class="flex gap-3">
      <div class="w-60">
        <Typography preset="Text md/Semibold" text={`${$_("documentUploads")}`} color="secondary" />
      </div>
      <div class="flex flex-col grow gap-6">
        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
          <!-- form upload -->
          {#each documentUploads as docType, index}
            <DocumentProject
              title={docType.title}
              handleFileUploadWrapper={() => handleFileUploadWrapper(docType.type, data.as)}
              file={getFileByType(docType.type, data.as)}
              placeholderDropdown={$_("setRules")}
              actionButton={true}
              actionView={() => (
                (documentViewModal = true), (selectedFile = getFileByType(docType.type, data.as))
              )}
              actionDelete={() => {
                typeTmp = docType.type;
                asTmp = data.as;
                deleteModal = true;
              }}
              actionGenerate={() => {
                const index = files.findIndex((file) => file.type === docType.type);
                if (files[index].status === "generated") {
                  goResultPage(files[index].id);
                } else if (files[index].status === "pending") {
                  console.log("pending");
                  alert("Please wait for the file to be generated");
                } else {
                  generatedAI(files[index]);
                }
              }}
            />
          {/each}
          <!-- {#each otherDocType as docType}
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
        {/each} -->
          <Button
            class="w-full min-h-32 flex gap-2.5 rounded-xl border border-primary-light dark:border-primary-dark bg-primary-light dark:bg-primary-dark"
            on:click={() => {
              addNewDocumentModal = true;
            }}
            color="light"
          >
            <PlusCircleICon svgClass="w-6 h-6 text-fg-tertiary-light dark:text-fg-tertiary-dark" />
            {$_("addComplianceDocs")}
          </Button>
        </div>
      </div>
    </div>
  </div>
{/if}

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
      <Input
        bind:value={tmpDocumentName}
        placeholder={$_("name")}
        class="bg-primary-light dark:bg-primary-dark"
      />
    </div>
    <UploadFile
      title={`${$_("uploadDocuments")}*`}
      handleFileUploadWrapper={() => handleFileUploadWrapper(tmpDocumentName, data.as)}
      file={tmpFile}
      actionButton={true}
      actionDelete={() => (tmpFile = [])}
    />
  </div>
  <svelte:fragment slot="footer">
    <Button on:click={addOtherDocument} disabled={tmpFile.length === 0 || !tmpDocumentName}
      >{$_("add")}</Button
    >
  </svelte:fragment>
</Modal>

<Modal bind:open={deleteModal} size="xs" autoclose>
  <div class="text-center">
    <Icon.ExclamationCircleSolid class="mx-auto mb-4 text-red-600 w-12 h-12" />
    <h3 class="mb-5 text-lg font-normal text-gray-500 dark:text-gray-400">
      {$_("areYouSureYouWantToDelete")}
    </h3>
    <Button
      onclick={() => {
        deleteFileByType(typeTmp, asTmp);
      }}
      color="red"
      class="me-2">{$_("yesImSure")}</Button
    >
    <Button
      color="alternative"
      onclick={() => {
        deleteModal = false;
        tmpSectionId = null;
      }}>{$_("noCancel")}</Button
    >
  </div>
</Modal>
