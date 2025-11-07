<script>
  import { onMount, onDestroy } from "svelte";
  import { goto } from "$app/navigation";
  import { page } from "$app/stores";
  import { _ } from "svelte-i18n";
  import { Button, Modal } from "flowbite-svelte";
  import { ExclamationCircleSolid } from "flowbite-svelte-icons";
  import { formatDistanceToNow, format } from "date-fns";
  import { call, check, getIcons, triggerToast, triggerLoading } from "$lib";
  import { Divider, InputRw, Dropdown, TextareaRw, Typography, DocumentProject } from "$lib";

  const isFrontend = process.env.NODE_ENV === "development";

  const PlusCircleIcon = getIcons("plusCircle");

  let timer;
  let timer2;
  let timer3;

  let init = $state(true);
  let ready = $state(false);
  let onEdit = $state(false);
  let initRAG = $state(true);
  let interval = $state(false);
  let ragStatus = $state(false);
  let deleteModal = $state(false);
  let updateModal = $state(false);
  let projectDelete = $state(false);

  let countSection = $state(0);

  let name = $state("");
  let typeTmp = $state("");
  let description = $state("");
  let projectCase = $state("");

  let data = $state({});

  let rules = $state([]);
  let files = $state([]);
  let rulesTmp = $state([]);
  let idResults = $state([]);

  let compareId = $state(null);

  const getRules = async () => {
    let token = localStorage.getItem("token");
    if (token) {
      rulesTmp = [];
      let page = 1;
      try {
        while (true) {
          const response = await $call.call("client-dashboard", "get-contract-rules", {
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
  };

  const getProject = async () => {
    if (init) triggerLoading(true);

    let token = localStorage.getItem("token");
    if (token) {
      try {
        const response = await $call.call("client-dashboard", "get-contract-project", {
          URLQuery: {},
          URLParams: {
            id: compareId.toString(),
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
              goto(
                isFrontend
                  ? "document_reviewer/compare"
                  : await $call.getPage("client-dashboard", "contract-compare-list"),
              );
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
            start_date: response.result.data[0].attributes.duration_start,
            end_date: response.result.data[0].attributes.duration_end,
            description: response.result.data[0].attributes.description,
          };

          name = data.name;
          description = data.description;
          projectCase = data.contract_rule_id;

          await getCountSection();
          await getDocumentContractProject();

          getStatusReference();
          getStatusGenerated();

          ready = true;

          if (init) triggerLoading(false);
        }
      } catch (error) {
        triggerLoading(false);
        triggerToast("error");
        console.log(error);
      }
    }
  };

  const getCountSection = async () => {
    let token = localStorage.getItem("token");
    if (token) {
      try {
        let tmpCount = 0;
        let page = 1;
        while (true) {
          const response = await $call.call("client-dashboard", "get-contract-rule-section", {
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
  };

  const getDocumentContractProject = async () => {
    let token = localStorage.getItem("token");
    try {
      const response2 = await $call.call("client-dashboard", "get-document-contract-project", {
        URLQuery: {},
        URLParams: {},
        Params: {
          project_id: compareId.toString(),
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

          // if (initRAG) {
          //   const response3 = await $call.call("hexa-ai-contract-review", "process-to-rag", {
          //     URLQuery: {},
          //     URLParams: {},
          //     Params: {
          //       project_id: compareId.toString(),
          //     },
          //     Headers: {
          //       Authorization: "Bearer " + token,
          //     },
          //   });

          //   if (response3.status !== 200) {
          //     triggerLoading(false);
          //     triggerToast("error");
          //     console.log(response3);
          //     return;
          //   }
          // }
        }

        initRAG = false;

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
  };

  const getStatusReference = async () => {
    let token = localStorage.getItem("token");
    if (token) {
      try {
        const response = await $call.call("client-dashboard", "check-reference-rag", {
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
  };

  const getStatusGenerated = async () => {
    let token = localStorage.getItem("token");
    if (token) {
      try {
        for (let i = 0; i < files.length; i++) {
          if (
            files[i].status_generated !== "START" &&
            files[i].status_generated !== "FINISH" &&
            files[i].status_generated !== "GENERATING"
          )
            continue;
          else {
            const response = await $call.call(
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
              for (let j = 0; j < response.result.data.length; j++) {
                if (response.result.data[j].attributes.status === "GENERATED") {
                  if (files[i].status_generated === "GENERATING") {
                    const response2 = await $call.call(
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
                  }

                  let id = response.result.data[j].id;

                  console.log("id", id);

                  // check if already push to idResults
                  let index = idResults.findIndex((result) => result.id_doc === files[i].id);
                  if (index !== -1) {
                    idResults[index].id = id;
                  } else {
                    idResults.push({
                      id_doc: files[i].id,
                      id: id,
                    });
                  }
                }

                if (response.result.data[j].attributes.status === "GENERATING") {
                  if (files[i].status_generated === "START") {
                    const response2 = await $call.call(
                      "client-dashboard",
                      "update-project-document-status-generated",
                      {
                        URLQuery: {},
                        URLParams: {
                          id: files[i].id.toString(),
                        },
                        Params: {
                          status_generated: "GENERATING",
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

                    files[i].status_generated = "GENERATING";
                  }

                  let id = response.result.data[j].id;

                  console.log("id", id);

                  let index = idResults.findIndex((result) => result.id_doc === files[i].id);
                  if (index !== -1) {
                    idResults[index].id = id;
                  } else {
                    idResults.push({
                      id_doc: files[i].id,
                      id: id,
                    });
                  }
                }

                if (response.result.data[j].attributes.status === "FAILED") {
                  files[i].status_generated = "ERROR";
                  triggerToast("error");
                  break;
                }
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
  };

  const generatedAI = async (file) => {
    let token = localStorage.getItem("token");
    if (token) {
      try {
        const response = await $call.call("hexa-ai-contract-review", "generate-review", {
          URLQuery: {},
          URLParams: {},
          Params: {
            project_id: compareId.toString(),
            project_document_id: file.id.toString(),
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

        const response2 = await $call.call(
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
        getStatusGenerated();
      } catch (error) {
        triggerToast("error");
        console.log(error);
      }
    }
  };

  const deleteFileByType = async (type) => {
    deleteModal = false;
    triggerLoading(true);
    const index = files.findIndex((file) => file.type === type);
    if (index !== -1) {
      let token = localStorage.getItem("token");
      if (token) {
        const response = await $call.call("client-dashboard", "delete-document-project", {
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

        const response2 = await $call.call(
          "client-dashboard",
          "delete-document-collection-by-contract",
          {
            URLQuery: {},
            URLParams: {},
            Params: {
              contract: "PROJECT",
              contract_doc_id: files[index].id.toString(),
            },
            Headers: {
              Authorization: "Bearer " + token,
            },
          },
        );

        if (response2.status !== 200) {
          triggerLoading(false);
          // triggerToast("error");
          // console.log(response2);
        }

        triggerToast("deleted");
        triggerLoading(false);
        files.splice(index, 1);
      }
    }
  };

  const deleteProject = async () => {
    let token = localStorage.getItem("token");

    if (token !== null) {
      try {
        const response = await $call.call("client-dashboard", "delete-contract-project", {
          URLQuery: {},
          URLParams: {
            id: compareId,
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

        triggerToast("deleted");
        goto(
          isFrontend
            ? "document_reviewer/compare"
            : await $call.getPage("client-dashboard", "contract-compare-list"),
        );
      } catch (error) {
        console.log(error);
      }
    }
  };

  const updateProject = async () => {
    try {
      let token = localStorage.getItem("token");

      if (!token) {
        throw new Error("No token found");
      }

      const response = await $call.call("client-dashboard", "update-contract-project", {
        URLQuery: {},
        URLParams: {
          id: compareId,
        },
        Params: {
          project_name: name,
          duration_start: format(new Date(data.start_date), "yyyy-MM-dd"),
          duration_end: format(new Date(data.end_date), "yyyy-MM-dd"),
          description: description,
        },
        Headers: {
          Authorization: "Bearer " + token,
        },
      });

      if (response.status !== 200) {
        const error = new Error("Failed to update project");
        error.response = response;
        throw error;
      }

      await getProject();
      triggerToast("updated");
    } catch (error) {
      console.log("error", error);
      triggerToast("error");
    }
  };

  const getFileByType = (type, id) => {
    return files.find((file) => file.type === type && file.id === id) || null;
  };

  const setIntervalGetDocument = () => {
    if (interval) {
      return;
    } else {
      interval = true;
      timer3 = setInterval(() => getDocumentContractProject(), 3000);
    }
  };

  onMount(() => {
    if (check("aiManagement.contractreviewer.project:view") === false) {
      console.log("no access");
      window.history.back();
      return;
    }
    const unsubscribe = page.subscribe(async ($page) => {
      compareId = $page.url.searchParams.get("compareId");
      if (compareId !== null) {
        getProject();
      } else {
        triggerToast("error");
        goto(
          isFrontend
            ? "/document_reviewer/compare"
            : await $call.getPage("client-dashboard", "contract-compare-list"),
        );
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
          <Typography preset="Display sm/Semibold" text={name} color="primary" />
        </div>
        <Typography
          preset="Text md/Medium"
          text={`${$_("createOn")} - ${data.createdAt}`}
          color="quarterary"
        />
      </div>
      {#if check("aiManagement.contractreviewer.project:delete")}
        <div class="flex flex-col gap-2">
          <div class="flex gap-2">
            {#if !onEdit}
              <Button
                on:click={() => {
                  onEdit = !onEdit;
                }}
                outline
              >
                {$_("edit")}
              </Button>
              <Button
                on:click={() => {
                  projectDelete = true;
                  typeTmp = null;
                  deleteModal = true;
                }}
                color="red"
                outline
              >
                {$_("delete")}
              </Button>
            {:else}
              <Button
                on:click={() => {
                  onEdit = !onEdit;
                  name = data.name;
                  description = data.description;
                  projectCase = data.contract_rule_id;
                }}
                color="alternative"
                outline
              >
                {$_("cancel")}
              </Button>
              <Button
                on:click={() => {
                  updateModal = true;
                }}
                outline
              >
                {$_("update")}
              </Button>
            {/if}
          </div>
        </div>
      {/if}
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
          <InputRw {onEdit} bind:value={name} />
        </div>
        <div class="flex gap-3">
          <div class="flex grow flex-col gap-1.5">
            <Typography preset="Text sm/medium" color="secondary" text={$_("case")} />
            <Dropdown
              dropdownOpen={false}
              bind:selectedItem={projectCase}
              items={rules}
              isOverlay={true}
              disabled={true}
            />
          </div>
        </div>
        <div class="flex flex-col gap-1.5">
          <Typography preset="Text sm/medium" color="secondary" text={$_("description")} />
          <TextareaRw {onEdit} bind:value={description} />
        </div>
      </div>
    </div>
    {#if check("aiManagement.contractreviewer.document:create")}
      <Divider />
      <div class="flex gap-3">
        <div class="w-60">
          <Typography
            preset="Text md/Semibold"
            text={`${$_("documentUploads")}`}
            color="secondary"
          />
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
                file={getFileByType(file.type, file.id)}
                bind:status={file.status_generated}
                actionButton={true}
                isLoading={file.status !== "FINISH"}
                actionView={async () => {
                  let gotoUrl = isFrontend
                    ? "/document_reviewer/compare/document"
                    : await $call.getPage("client-dashboard", "contract-compare-detail-document");
                  goto(gotoUrl + "?compareDocumentId=" + file.id + "&compareId=" + compareId);
                  // (documentViewModal = true), (selectedFile = getFileByType(file.type));
                }}
                actionDelete={() => {
                  typeTmp = file.type;

                  projectDelete = false;
                  deleteModal = true;
                }}
                actionGenerate={async () => {
                  if (
                    file.status_generated !== "FINISH" &&
                    file.status_generated !== "GENERATING"
                  ) {
                    generatedAI(file);
                  } else {
                    if (file.result_id !== undefined) {
                      let gotoUrl = isFrontend
                        ? "/document_reviewer/compare/review"
                        : await $call.getPage("client-dashboard", "contract-compare-review");
                      goto(
                        gotoUrl +
                          "?compareDocumentId=" +
                          file.result_id +
                          "&compareId=" +
                          compareId,
                      );
                    } else {
                      let index = idResults.findIndex((result) => result.id_doc === file.id);
                      let gotoUrl = isFrontend
                        ? "/document_reviewer/compare/review"
                        : await $call.getPage("client-dashboard", "contract-compare-review");
                      goto(
                        gotoUrl +
                          "?compareDocumentId=" +
                          idResults[index].id +
                          "&compareId=" +
                          compareId,
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
                let gotoUrl = isFrontend
                  ? "/document_reviewer/compare/document"
                  : await $call.getPage("client-dashboard", "contract-compare-detail-document");
                goto(gotoUrl + "?compareId=" + compareId);
              }}
              color="light"
            >
              <PlusCircleIcon
                svgClass="w-6 h-6 text-fg-tertiary-light dark:text-fg-tertiary-dark"
              />
              {$_("addComplianceDocs")}
            </Button>
          </div>
        </div>
      </div>
    {/if}
  </div>
{/if}

<Modal bind:open={deleteModal} size="xs" autoclose>
  <div class="text-center">
    <ExclamationCircleSolid class="mx-auto mb-4 text-red-600 w-12 h-12" />
    <h3 class="mb-5 text-lg font-normal text-gray-500 dark:text-gray-400">
      {$_("areYouSureYouWantToDelete")}
    </h3>
    <Button
      onclick={() => {
        if (projectDelete === false) {
          deleteFileByType(typeTmp);
        } else {
          deleteProject();
        }
      }}
      color="red"
      class="me-2"
    >
      {$_("yesImSure")}
    </Button>
    <Button
      color="alternative"
      onclick={() => {
        deleteModal = false;
        typeTmp = null;
        projectDelete = false;
      }}
    >
      {$_("noCancel")}
    </Button>
  </div>
</Modal>

<Modal bind:open={updateModal} size="xs" autoclose>
  <div class="text-center">
    <ExclamationCircleSolid class="mx-auto mb-4 text-yellow-400 w-12 h-12" />
    <h3 class="mb-5 text-lg font-normal text-gray-500 dark:text-gray-400">
      {$_("areYouSureYouWantToUpdate")}
    </h3>
    <Button
      onclick={() => {
        updateProject();
        onEdit = false;
      }}
      color="primary"
      class="me-2"
    >
      {$_("yesImSure")}
    </Button>
    <Button
      color="alternative"
      onclick={() => {
        updateModal = false;
        onEdit = false;
      }}
    >
      {$_("noCancel")}
    </Button>
  </div>
</Modal>
