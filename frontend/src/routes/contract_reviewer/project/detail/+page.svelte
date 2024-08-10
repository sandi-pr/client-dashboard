<script>
  import {
    Typography,
    InputRw,
    TextareaRw,
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
  import { breadcrumbLastName } from "$lib/stores/breadcrumbLastName.js";

  const call = new CR3Call(globalThis.location.origin);

  let projectId = $state(null);

  let data = $state({});

  let deleteModal = $state(false);
  let ready = $state(false);
  let files = $state([]);
  let tmpFile = $state([]);
  let typeTmp = $state("");
  let asTmp = $state("");
  let rulesTmp = $state([]);
  let init = $state(true);
  let contractReady = $state(false);
  let timer;
  let timer2;
  let timer3;
  let interval = $state(false);
  let ragStatus = $state(false);
  let countSection = $state(0);
  let idResults = $state([]);

  let rules = $state([]);

  function formatBytes(bytes, decimals = 2) {
    if (bytes === 0) return "0 Bytes";

    const k = 1024;
    const dm = decimals < 0 ? 0 : decimals;
    const sizes = ["Bytes", "KB", "MB", "GB", "TB"];

    const i = Math.floor(Math.log(bytes) / Math.log(k));

    return parseFloat((bytes / Math.pow(k, i)).toFixed(dm)) + " " + sizes[i];
  }

  function getFileByType(type) {
    const index = files.findIndex((file) => file.type === type);
    return index !== -1 ? files[index] : [];
  }

  async function deleteFileByType(type) {
    deleteModal = false;
    triggerLoading(true);
    const index = files.findIndex((file) => file.type === type);
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
          triggerLoading(false);
          triggerToast("error");
          console.log(response);
          return;
        }

        triggerToast("deleted");
        triggerLoading(false);
        files.splice(index, 1);
      }
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

          if (response.result.data.length < 1) {
            break;
          }

          rulesTmp = rulesTmp.concat(response.result.data);
          page += 1;
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

          await getRules();

          data = {
            id: response.result.data[0].id,
            name: response.result.data[0].attributes.project_name,
            createdAt: formatDistanceToNow(new Date(response.result.data[0].attributes.createdAt), {
              addSuffix: true,
            }),
            industry: response.result.data[0].attributes.industry_id,
            contract_rule_id: response.result.data[0].attributes.contract_rule_id,
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

          await getCountSection();
          await getDocumentContractProject();

          getStatusReference();
          getStatusGenerated();

          // breadcrumbLastName.set(data.name);

          ready = true;

          if (init) triggerLoading(false);
        }
      } catch (error) {
        triggerLoading(false);
        triggerToast("error");
        console.log(error);
      }
    }
    // breadcrumbLastName.set("Bruh");
  }

  async function getCountSection() {
    let token = sessionStorage.getItem("token");
    if (token) {
      try {
        let tmpCount = 0;
        let page = 1;
        while (true) {
          const response = await call.call("client-dashboard", "get-contract-rule-section", {
            URLQuery: {},
            URLParams: {},
            Params: {
              contract_rule_id: data.contract_rule_id.toString(),
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

          page += 1;
          tmpCount += response.result.data.length;
        }

        countSection = tmpCount;
      } catch (error) {
        triggerToast("error");
        console.log(error);
      }
    }
  }

  async function generatedAI(file) {
    let token = sessionStorage.getItem("token");
    if (token) {
      try {
        const response = await call.call("hexa-ai-contract-review", "generate-review", {
          URLQuery: {},
          URLParams: {},
          Params: {
            project_id: projectId.toString(),
            project_document_id: file.id,
          },
          Headers: {
            Authorization: "Bearer " + token,
          },
        });

        if (response.status !== 200) {
          triggerToast("error");
          console.log(response);
          return;
        }

        const response2 = await call.call(
          "client-dashboard",
          "update-project-document-status-generated",
          {
            URLQuery: {},
            URLParams: {
              id: file.id.toString(),
            },
            Params: {
              status_generated: "START",
            },
            Headers: {
              Authorization: "Bearer " + token,
            },
          },
        );

        if (response2.status !== 200) {
          triggerToast("error");
          console.log(response2);
          return;
        }

        setIntervalGetDocument();
      } catch (error) {
        triggerToast("error");
        console.log(error);
      }
    }
  }

  async function getDocumentContractProject() {
    let token = sessionStorage.getItem("token");
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
        if (init) triggerLoading(false);
        triggerToast("error");
        console.log(response2);
        return;
      }

      if (response2.status === 200) {
        if (response2.result.files.length > 0) {
          let promise = response2.result.files.map((file) => {
            return {
              type: file.type,
              name: file.name,
              base64: file.base64,
              id: file.id,
              status: ragStatus ? file.status : "PROCESSING",
              status_generated: file.status_generated,
              role: "Buyer",
              rag_id: file.rag_id,
            };
          });

          files = await Promise.all(promise);

          if (files.length < 1) {
            interval = false;
            clearInterval(timer3);
            return;
          }
        }

        let allReady = true;

        for (let i = 0; i < files.length; i++) {
          if (files[i].status !== "FINISH") {
            allReady = false;
            break;
          }
        }

        if (allReady) {
          // let prosess = false;
          // for (let i = 0; i < files.length; i++) {
          //   if (files[i].status_generated === "START") {
          //     prosess = true;
          //     break;
          //   }
          // }

          // if (prosess) {
          //   return;
          // } else {
          interval = false;
          clearInterval(timer3);
          // }
        }
      }
    } catch (error) {
      triggerToast("error");
      console.log(error);
    }
  }

  function setIntervalGetDocument() {
    if (interval) {
      return;
    } else {
      interval = true;
      timer3 = setInterval(() => getDocumentContractProject(), 3000);
    }
  }

  async function getStatusReference() {
    let token = sessionStorage.getItem("token");
    if (token) {
      try {
        const response = await call.call("client-dashboard", "check-reference-rag", {
          URLQuery: {},
          URLParams: {},
          Params: {
            contract_rule_id: data.contract_rule_id,
          },
          Headers: {
            Authorization: "Bearer " + token,
          },
        });

        if (response.status !== 200) {
          triggerToast("error");
          console.log("status reference", response);
          return;
        }

        ragStatus = response.result.ready;

        if (ragStatus) {
          setIntervalGetDocument();
          clearTimeout(timer2);
        } else {
          timer2 = setTimeout(() => getStatusReference(), 3000);
        }
      } catch (error) {
        console.log(error);
      }
    }
  }

  async function getStatusGenerated() {
    let token = sessionStorage.getItem("token");
    if (token) {
      try {
        for (let i = 0; i < files.length; i++) {
          if (files[i].status_generated !== "START") continue;
          else {
            const response = await call.call(
              "client-dashboard",
              "get-contract-project-result-with-document-id",
              {
                URLQuery: {},
                URLParams: {},
                Params: {
                  project_document_id: files[i].id.toString(),
                },
                Headers: {
                  Authorization: "Bearer " + token,
                },
              },
            );

            if (response.status !== 200) {
              triggerLoading(false);
              triggerToast("error");
              console.log(response);
              break;
            }

            if (response.result.data.length > 0) {
              console.log(response.result.data[0].attributes.result.length, countSection);
              if (response.result.data[0].attributes.result.length === countSection) {
                const response2 = await call.call(
                  "client-dashboard",
                  "update-project-document-status-generated",
                  {
                    URLQuery: {},
                    URLParams: {
                      id: files[i].id.toString(),
                    },
                    Params: {
                      status_generated: "FINISH",
                    },
                    Headers: {
                      Authorization: "Bearer " + token,
                    },
                  },
                );

                if (response2.status !== 200) {
                  triggerLoading(false);
                  triggerToast("error");
                  console.log(response2);
                  break;
                }

                files[i].status_generated = "FINISH";
                files[i].result_id = response.result.data[0].id;

                let id = response.result.data[0].id;

                idResults.push({
                  idDocument: files[i].id,
                  idResult: id,
                });
              }
            }
          }
        }

        let allGenerated = true;

        for (let i = 0; i < files.length; i++) {
          if (files[i].status_generated !== "FINISH") {
            allGenerated = false;
            break;
          }
        }

        if (allGenerated) {
          clearTimeout(timer);
        } else {
          timer = setTimeout(() => {
            getStatusGenerated();
          }, 3000);
        }
      } catch (error) {
        triggerLoading(false);
        triggerToast("error");
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

  onDestroy(() => {
    clearTimeout(timer);
    clearTimeout(timer2);
    clearInterval(timer3);
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
        <Typography
          preset="Text md/Semibold"
          text={`${$_("generalInformation")}`}
          color="secondary"
        />
      </div>
      <div class="flex flex-col grow gap-3">
        <!-- form dropdown -->
        <div class="flex flex-col gap-1.5">
          <Typography preset="Text sm/medium" color="secondary" text={$_("projectTitle")} />
          <InputRw onEdit={false} bind:value={data.name} />
        </div>
        <div class="flex gap-3">
          <div class="flex grow flex-col gap-1.5">
            <Typography preset="Text sm/medium" color="secondary" text={$_("case")} />
            <Dropdown
              dropdownOpen={false}
              bind:selectedItem={data.contract_rule_id}
              items={rules}
              isOverlay={true}
              disabled={true}
            />
          </div>
          <div class="flex grow flex-col gap-1.5">
            <Typography
              preset="Text sm/medium"
              color="secondary"
              text={`${$_("projectDuration")}`}
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
                class="flex w-full p-2 border rounded-lg bg-primary-light dark:bg-gray-800 border-primary-light dark:border-primary-dark gap-2 justify-center items-center hover:cursor-default"
              >
                <svelte:component
                  this={CalendarPlus02}
                  svgClass="w-4 h-4 ms-2 text-secondary-light dark:text-secondary-dark"
                />
                <div class="date text-secondary-light dark:text-secondary-dark">
                  {data.start_date} - {data.end_date}
                </div>
              </button>
            </DatePicker>
          </div>
        </div>
        <div class="flex flex-col gap-1.5">
          <Typography preset="Text sm/medium" color="secondary" text={$_("description")} />
          <TextareaRw onEdit={false} bind:value={data.description} />
        </div>
      </div>
    </div>
    <Divider />
    <div class="flex gap-3">
      <div class="w-60">
        <Typography preset="Text md/Semibold" text={`${$_("documentUploads")}`} color="secondary" />
        <!-- <button
          onclick={() => {
            getStatusGenerated();
          }}
          class="text-brand-secondary-light dark:text-brand-secondary-dark"
          >{$_("checkStatus")}</button
        > -->
      </div>
      <div class="flex flex-col grow gap-6">
        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
          <!-- form upload -->
          {#each files as file}
            <DocumentProject
              title={file.type}
              handleFileUploadWrapper={() => handleFileUploadWrapper(file.type)}
              file={getFileByType(file.type)}
              bind:status={file.status_generated}
              actionButton={true}
              isLoading={file.status !== "FINISH"}
              actionView={async () => {
                goto(
                  (await call.getPage("client-dashboard", "contract-project-detail-document")) +
                    "?projectDocumentId=" +
                    file.id,
                );
                // (documentViewModal = true), (selectedFile = getFileByType(file.type));
              }}
              actionDelete={() => {
                typeTmp = file.type;
                deleteModal = true;
              }}
              actionGenerate={async () => {
                if (file.status_generated !== "FINISH") {
                  generatedAI(file);
                } else {
                  if (file.result_id !== undefined) {
                    goto(
                      (await call.getPage("client-dashboard", "contract-project-review")) +
                        "?projectDocumentId=" +
                        file.result_id,
                    );
                  } else {
                    console.log("gamasuk");
                    let index = idResults.findIndex((result) => result.idDocument === file.id);
                    goto(
                      (await call.getPage("client-dashboard", "contract-project-review")) +
                        "?projectDocumentId=" +
                        idResults[index].idResult,
                    );
                  }
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
            on:click={async () => {
              goto(
                (await call.getPage("client-dashboard", "contract-project-detail-document")) +
                  "?projectId=" +
                  projectId,
              );
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

<Modal bind:open={deleteModal} size="xs" autoclose>
  <div class="text-center">
    <Icon.ExclamationCircleSolid class="mx-auto mb-4 text-red-600 w-12 h-12" />
    <h3 class="mb-5 text-lg font-normal text-gray-500 dark:text-gray-400">
      {$_("areYouSureYouWantToDelete")}
    </h3>
    <Button
      onclick={() => {
        deleteFileByType(typeTmp);
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
