<!--
CR3 Client Dashboard
Package: gitlab.hexacode.org/cr3/client-dashboard
Maintainer:
  - Dean Ramhan <dean@hexacode.org>
  - Sandi Permana <sandy@hexacode.org>
  - Taqie Fadlillah <taqief@hexacode.org>
Copyright (C) 2024 Hexacode Teknologi Indonesia
All Rights Reserved
-->
<script>
  import { onMount } from "svelte";
  import { Modal, Button, Input, Heading } from "flowbite-svelte";
  import { page } from "$app/stores";
  import { goto } from "$app/navigation";
  import { ExclamationCircleSolid } from "flowbite-svelte-icons";
  import { _ } from "svelte-i18n";
  import { format } from "date-fns";
  import { call, check, getIcons, triggerToast, triggerLoading } from "$lib";
  import {
    InputRw,
    Divider,
    TabDetail,
    StateCard,
    Typography,
    UploadFile,
    TextareaRw,
    PdfjsViewer,
  } from "$lib";

  const IllustrationEmpty = getIcons("illustrationEmpty");
  const PlusCircleIcon = getIcons("plusCircle");
  const PlusIcon = getIcons("plus");

  const isFrontend = process.env.NODE_ENV === "development";

  let init = $state(false);
  let edit = $state(false);
  let caseDelete = $state(false);
  let deleteModal = $state(false);
  let addRefDocModal = $state(false);
  let documentViewModal = $state(false);
  let addRulesSectionsModal = $state(false);

  let formData = $state({ title: "", prompt: "", sample: "" });
  let tabDatas = $state({
    selected: "Referenced Documents",
    data: [
      check("aiManagement.contractreviewer.document:create_delete") && {
        name: "Referenced Documents",
      },
      check("aiManagement.contractreviewer.rule:create_delete") && {
        name: "Rules",
      },
    ].filter(Boolean),
  });
  let datas = $state({
    title: "",
    createdAt: "",
    data: [],
  });
  let datasDataDuplicate = $state([]);

  let editRules = $state([]);
  let files = $state([
    // {
    //   id: 1,
    //   type: "Testing",
    //   role: "buyer",
    //   name: "test.pdf",
    //   size: "1.2 MB",
    //   base64: "",
    // }
  ]);

  let tmpId = $state(null);
  let ruleId = $state(null);
  let tmpFile = $state(null);
  let selectedFile = $state(null);
  let tmpSectionId = $state(null);
  let collectionId = $state(null);
  let tmpDocumentName = $state(null);

  // let ifDataMultiple = $derived(() => {
  //   return datas.data.length > 1;
  // });

  const getCollections = async (name, id) => {
    let token = localStorage.getItem("token");
    if (token !== null) {
      try {
        let response = await $call.call("client-dashboard", "get-collection-by-name", {
          URLQuery: {},
          URLParams: {},
          Params: {
            collection_name: id + ". " + name + " (Referenced Documents)",
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
          let response2 = await $call.call("client-dashboard", "create-document-collection", {
            URLQuery: {},
            URLParams: {},
            Params: {
              collection_name: id + ". " + name + " (Referenced Documents)",
              description: "Collection for " + name,
            },
            Headers: {
              Authorization: "Bearer " + token,
            },
          });

          if (response2.status !== 200) {
            console.log(response2);
            return;
          }

          collectionId = response2.result.data.id;
          return;
        }

        collectionId = response.result.data[0].id;
      } catch (error) {
        console.log(error);
      }
    }
  };

  const getRulesSections = async () => {
    triggerLoading(true);
    let token = localStorage.getItem("token");
    if (token) {
      try {
        const response = await $call.call("client-dashboard", "get-contract-rule", {
          URLQuery: {},
          URLParams: {
            id: ruleId.toString(),
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

        await getCollections(
          response.result.data[0].attributes.rule_name,
          response.result.data[0].id,
        );

        let sectionsTmp = [];
        let page = 1;

        while (true) {
          const response2 = await $call.call("client-dashboard", "get-contract-rule-section", {
            URLQuery: {},
            URLParams: {},
            Params: {
              contract_rule_id: ruleId.toString(),
              page: page,
              page_size: 25,
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

          if (response2.result.data.length < 1) {
            break;
          }

          page += 1;
          sectionsTmp = sectionsTmp.concat(response2.result.data);
          datasDataDuplicate = sectionsTmp.map((item) => {
            return {
              id: item.id,
              title: item.attributes.section_title,
              prompt: item.attributes.section_prompt,
              reference_prompt: item.attributes.reference_prompt,
              sample: item.attributes.section_sample,
              createdAt: format(item.attributes.createdAt, "yyyy-MM-dd"),
            };
          });
        }

        editRules = [];
        if (sectionsTmp.length < 1) {
          init = true;
          editRules.push(true);
        }

        for (let i = 0; i < sectionsTmp.length; i++) {
          editRules.push(false);
        }

        datas = {
          title: response.result.data[0].attributes.rule_name,
          createdAt: format(response.result.data[0].attributes.createdAt, "yyyy-MM-dd HH:mm:ss"),
          data:
            sectionsTmp.length < 1
              ? datas.data
              : sectionsTmp.map((item) => {
                  return {
                    id: item.id,
                    title: item.attributes.section_title,
                    prompt: item.attributes.section_prompt,
                    reference_prompt: item.attributes.reference_prompt,
                    sample: item.attributes.section_sample,
                    createdAt: format(item.attributes.createdAt, "yyyy-MM-dd"),
                  };
                }),
        };

        triggerLoading(false);
      } catch (error) {
        triggerLoading(false);
        triggerToast("error");
        console.log(error);
      }
    }
  };

  const getReferences = async () => {
    let token = localStorage.getItem("token");
    let industry_id = sessionStorage.getItem("industry_id");

    if (industry_id === null) {
      triggerLoading(false);

      return;
    }

    try {
      let page = 1;
      files = [];

      while (true) {
        const response = await $call.call("client-dashboard", "get-references-document", {
          URLQuery: {},
          URLParams: {},
          Params: {
            industry_id: parseInt(industry_id),
            contract_rule_id: parseInt(ruleId),
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

        const filesWithExt = response.result.files.map((file) => ({
          ...file,
          ext: ".pdf",
        }));

        files = files.concat(filesWithExt);
        page++;
      }
      // console.log("files", files);
    } catch (error) {
      console.log(error);
    }
  };

  const createRulesSections = async (data) => {
    triggerLoading(true);
    let token = localStorage.getItem("token");
    if (token) {
      try {
        const response = await $call.call("client-dashboard", "create-contract-rule-section", {
          URLQuery: {},
          URLParams: {},
          Params: {
            contract_rule_id: ruleId.toString(),
            section_title: data.title,
            section_prompt: data.prompt,
            reference_prompt: data.reference_prompt,
            section_sample: data.sample,
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

        triggerLoading(false);
        triggerToast("created");
      } catch (error) {
        triggerLoading(false);
        triggerToast("error");
        console.log(error);
      }
    }
  };

  const addOtherDocument = async () => {
    addRefDocModal = false;
    triggerLoading(true);
    let token = localStorage.getItem("token");
    let industry_id = sessionStorage.getItem("industry_id");
    console.log("industry_id", industry_id);
    if (industry_id === null) {
      triggerLoading(false);

      return;
    }
    try {
      tmpFile.type = tmpDocumentName;

      const response = await $call.call("client-dashboard", "create-references-document", {
        URLQuery: {},
        URLParams: {},
        Params: {
          references: JSON.stringify([tmpFile]),
          contract_rule_id: parseInt(ruleId),
          industry_id: parseInt(industry_id),
        },
        Headers: {
          Authorization: "Bearer " + token,
        },
      });
      // console.log(response);

      if (response.status !== 200) {
        triggerLoading(false);
        addRefDocModal = true;
        triggerToast("error");
      }

      tmpFile.id = response.result.data[0].id;

      tmpFile.document_title = tmpDocumentName;
      tmpFile.description = "Document reference " + datas.title;
      tmpFile.contract = "REFERENCE";

      let response2 = await $call.call("client-dashboard", "create-document", {
        URLQuery: {},
        URLParams: {},
        Params: {
          documents: JSON.stringify([tmpFile]),
          collection_id: collectionId.toString(),
          document_type: "Documents",
        },
        Headers: {
          Authorization: "Bearer " + token,
        },
      });

      if (response2.status !== 200) {
        triggerLoading(false);
        addRefDocModal = true;
        triggerToast("error");
      }

      files.push(tmpFile);

      tmpFile = null;
      tmpDocumentName = "";
      triggerLoading(false);
      triggerToast("updated");
    } catch (error) {
      console.log(error);
    }
  };

  const editSection = async () => {
    triggerLoading(true);
    init = false;
    let token = localStorage.getItem("token");
    let tmpSection = datas.data.find((item) => item.id === tmpSectionId);
    if (token) {
      try {
        const response = await $call.call("client-dashboard", "update-contract-rule-section", {
          URLQuery: {},
          URLParams: {
            id: tmpSection.id.toString(),
          },
          Params: {
            section_title: tmpSection.title,
            section_prompt: tmpSection.prompt,
            reference_prompt: tmpSection.reference_prompt,
            section_sample: tmpSection.sample,
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

        triggerLoading(false);
        triggerToast("updated");
        await getRulesSections();
      } catch (error) {
        console.log(error);
      }
    }
  };

  const deleteSection = async () => {
    deleteModal = false;
    triggerLoading(true);
    let token = localStorage.getItem("token");
    if (token) {
      try {
        const response = await $call.call("client-dashboard", "delete-contract-rule-section", {
          URLQuery: {},
          URLParams: {
            id: tmpSectionId.toString(),
          },
          Params: {},
          Headers: {
            Authorization: "Bearer " + token,
          },
        });

        if (response.status !== 200) {
          triggerLoading(false);
          triggerToast("error");
          deleteModal = true;
          console.log(response);
          return;
        }

        triggerLoading(false);
        triggerToast("deleted");
        await getRulesSections();
      } catch (error) {
        console.log(error);
      }
    }
  };

  const deleteDocument = async () => {
    deleteModal = false;
    triggerLoading(true);

    let token = localStorage.getItem("token");

    if (tmpId !== 0) {
      const response = await $call.call("client-dashboard", "delete-references-document", {
        URLQuery: {},
        URLParams: {
          id: tmpId.toString(),
        },
        Params: {},
        Headers: {
          Authorization: "Bearer " + token,
        },
      });

      if (response.status !== 200) {
        triggerLoading(false);
        triggerToast("error");
        deleteModal = true;
      }

      const response2 = await $call.call(
        "client-dashboard",
        "delete-document-collection-by-contract",
        {
          URLQuery: {},
          URLParams: {},
          Params: {
            contract: "REFERENCE",
            contract_doc_id: tmpId.toString(),
          },
          Headers: {
            Authorization: "Bearer " + token,
          },
        },
      );

      if (response2.status !== 200) {
        triggerLoading(false);
        // triggerToast("error");
        // deleteModal = true;
      }

      files = files.filter((file) => file.id !== tmpId);
      triggerLoading(false);
      triggerToast("deleted");
    }
  };

  const deleteCase = async () => {
    let token = localStorage.getItem("token");
    if (token) {
      try {
        const response = await $call.call("client-dashboard", "delete-rules", {
          URLQuery: {},
          URLParams: {
            id: ruleId.toString(),
          },
          Params: {},
          Headers: {
            Authorization: "Bearer " + token,
          },
        });

        if (response.status === 200) {
          deleteModal = false;
        }

        triggerToast("deleted");
        goto(
          isFrontend
            ? "/contract_reviewer/case"
            : await $call.getPage("client-dashboard", "contract-case"),
        );
      } catch (error) {
        console.log(error);
      }
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

  const handleFileUploadWrapper = (type) => {
    const inputElement = document.createElement("input");
    inputElement.type = "file";
    inputElement.accept = ".pdf";

    inputElement.addEventListener("change", async (event) => {
      const uploadedFile = event.target.files[0];

      const base64Data = await readFileAsBase64(uploadedFile);
      tmpFile = {
        id: 0,
        type: type,
        role: "Buyer",
        name: uploadedFile.name,
        size: formatBytes(uploadedFile.size),
        base64: base64Data,
        ext: ".pdf",
      };
      console.log(tmpFile);
    });

    inputElement.click();
  };

  const formatBytes = (bytes, decimals = 2) => {
    if (bytes === 0) return "0 Bytes";

    const k = 1024;
    const dm = decimals < 0 ? 0 : decimals;
    const sizes = ["Bytes", "KB", "MB", "GB", "TB"];

    const i = Math.floor(Math.log(bytes) / Math.log(k));

    return parseFloat((bytes / Math.pow(k, i)).toFixed(dm)) + " " + sizes[i];
  };

  const checkSectionInputs = (data) => {
    return (
      data.title === "" || data.prompt === "" || data.sample === "" || data.reference_prompt === ""
    );
  };

  onMount(() => {
    if (
      check("aiManagement.contractreviewer.rule:create_delete") === false &&
      check("aiManagement.contractreviewer.document:create_delete") === false
    ) {
      console.log("no access");
      window.history.back();
      return;
    }
    const unsubscribe = page.subscribe(async ($page) => {
      ruleId = $page.url.searchParams.get("caseId");
      triggerLoading(true);
      if (ruleId !== null) {
        getReferences();
        getRulesSections();
      } else {
        triggerToast("error");
        setTimeout(async () => {
          triggerLoading(false);
          goto(
            isFrontend
              ? "/contract_reviewer/case"
              : await $call.getPage("client-dashboard", "contract-case"),
          );
        }, 2000);
      }
    });
    return () => unsubscribe();
  });
</script>

<div class="flex dark:bg-primary-dark">
  <div class="w-full pb-10">
    <div class="flex flex-col m-4">
      <!-- <Navbar /> -->

      <div class="flex flex-col">
        <!-- User Information -->
        <div class="flex justify-between w-full">
          <!-- Header -->
          <div class="grow gap-2">
            <Heading
              tag="h3"
              class="text-3xl leading-[38px] text-hc-graylight-900 dark:text-hc-graydark-50"
            >
              {datas.title}
            </Heading>
            <Typography
              preset="Text md/Medium"
              color="tertiary"
              text={`${$_("createOn")} ${datas.createdAt}`}
            />
          </div>
          {#if check("aiManagement.contractreviewer.case:delete")}
            <div class="flex flex-col gap-2">
              <div class="flex flex-col">
                <Button
                  onclick={() => {
                    caseDelete = true;
                    tmpId = null;
                    deleteModal = true;
                  }}
                  color="red"
                  outline
                >
                  {$_("delete")}
                </Button>
              </div>
            </div>
          {/if}
        </div>

        <!-- Tab Bar -->
        <TabDetail {tabDatas} on:tabChange={(e) => (tabDatas.selected = e.detail)} />

        {#if tabDatas.selected === "Referenced Documents" && check("aiManagement.contractreviewer.document:create_delete")}
          <div>
            <!-- header -->
            <div class="flex justify-between pt-3 pb-6">
              <div class="flex flex-col">
                <Typography
                  preset="Text lg/Semibold"
                  color="primary"
                  text={$_("setYourReferencedDocument")}
                />
                <Typography
                  preset="Text sm/Regular"
                  color="tertiary"
                  text={$_("descSetYourReferencedDocument")}
                />
              </div>
            </div>
            <!-- reference document -->
            <div class="flex gap-3">
              <div class="w-60">
                <Typography
                  preset="Text md/Semibold"
                  text={`${$_("documents")}`}
                  color="secondary"
                />
              </div>
              <div class="flex flex-col grow gap-6">
                <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                  {#each files as doc}
                    <UploadFile
                      title={doc.type}
                      handleFileUploadWrapper={() => {
                        addRefDocModal = true;
                      }}
                      file={doc}
                      actionButton={true}
                      actionView={() => {
                        selectedFile = doc;
                        documentViewModal = true;
                      }}
                      actionDelete={() => {
                        tmpId = doc.id;
                        caseDelete = false;
                        deleteModal = true;
                      }}
                    />
                  {/each}
                  <Button
                    class="w-full min-h-32 md:mt-6 flex gap-2.5 rounded-xl text-tertiary-light dark:text-tertiary-dark border border-primary-light dark:border-primary-dark bg-primary-light dark:bg-primary-dark"
                    on:click={() => (addRefDocModal = true)}
                    color="light"
                  >
                    <PlusCircleIcon
                      svgClass="w-6 h-6 text-fg-tertiary-light dark:text-fg-tertiary-dark"
                    />
                    {$_("addReferenceDocs")}
                  </Button>
                </div>
              </div>
            </div>
          </div>
        {:else if tabDatas.selected === "Rules" && check("aiManagement.contractreviewer.rule:create_delete")}
          <div>
            <!-- header -->
            <div class="flex justify-between pt-3 pb-6">
              <div class="flex flex-col">
                <Typography preset="Text lg/Semibold" color="primary" text={$_("setRules")} />
                <Typography preset="Text sm/Regular" color="tertiary" text={$_("descSetRules")} />
              </div>
              <Button
                class="self-center gap-1"
                color="primary"
                on:click={() => {
                  datas.data = [
                    ...datas.data,
                    {
                      id: datas.data.length,
                      title: "",
                      prompt: "",
                      reference_prompt: "",
                      sample: "",
                      createdAt: "",
                    },
                  ];
                  editRules = [...editRules, true];
                }}
              >
                <PlusIcon svgClass="w-4 h-4" />
                {$_("addNewSection")}
              </Button>
            </div>
            <!-- table latest -->
            {#if datas.data.length > 0}
              <div class="flex flex-col gap-6">
                {#each datas.data as data, index}
                  <div class="flex gap-4">
                    <div class="flex flex-col w-60 gap-2">
                      <Typography
                        preset="Text md/Semibold"
                        text={$_("section") + " " + (data.id === 0 ? "1" : index + 1)}
                        color="secondary"
                      />
                      <div class="flex flex-col gap-2">
                        {#if !editRules[index]}
                          {#if datas.data.length > 1 && !init}
                            <button
                              class="text-sm text-utility-gray-600-light dark:text-utility-gray-600-dark"
                              onclick={() => {
                                tmpSectionId = data.id;
                                caseDelete = false;
                                deleteModal = true;
                              }}
                            >
                              <Typography
                                preset="Text sm/Semibold"
                                text={$_("deleteSection")}
                                color="error-primary"
                              />
                            </button>
                          {/if}
                          <button
                            class="text-sm text-brand-secondary-light dark:text-brand-secondary-dark"
                            onclick={() => {
                              editRules[index] = true;
                            }}
                          >
                            <Typography
                              preset="Text sm/Semibold"
                              text={$_("edit")}
                              color="fg-brand-primary"
                            />
                          </button>
                        {:else}
                          {#if !init}
                            <button
                              onclick={() => {
                                if (datasDataDuplicate[index] !== undefined) {
                                  datas.data[index] = datasDataDuplicate[index];
                                  editRules[index] = false;
                                } else {
                                  datas.data.splice(index, 1);
                                  editRules.splice(index, 1);
                                }
                              }}
                            >
                              <Typography
                                preset="Text sm/Semibold"
                                text={$_("cancel")}
                                color="tertiary"
                              />
                            </button>
                          {/if}
                          <button
                            class={checkSectionInputs(data) ? "cursor-not-allowed" : ""}
                            onclick={() => {
                              editRules[index] = false;
                              if (datasDataDuplicate[index] === undefined) {
                                createRulesSections(data);
                              } else {
                                tmpSectionId = data.id;
                                editSection();
                              }
                            }}
                            disabled={checkSectionInputs(data)}
                          >
                            <Typography
                              preset="Text sm/Semibold"
                              text={$_("save")}
                              color={checkSectionInputs(data) ? "tertiary" : "fg-brand-primary"}
                            />
                          </button>
                        {/if}
                      </div>
                    </div>
                    <div class="flex flex-col grow gap-3">
                      <div class="flex gap-3">
                        <!-- <div class="flex flex-col gap-1.5">
                          <Typography
                            preset="Text sm/Medium"
                            text={$_("order")}
                            color="secondary"
                          />
                          <InputRw value={"#" + (index + 1)} onEdit={ifDataMultiple()} />
                        </div> -->
                        <div class="flex flex-col grow gap-1.5">
                          <Typography preset="Text sm/Medium" text={$_("name")} color="secondary" />
                          <InputRw bind:value={data.title} onEdit={editRules[index]} />
                        </div>
                      </div>
                      <div class="flex flex-col grow gap-1.5">
                        <Typography
                          preset="Text sm/Medium"
                          text={$_("extractDataPrompt")}
                          color="secondary"
                        />
                        <TextareaRw bind:value={data.prompt} onEdit={editRules[index]} />
                      </div>
                      <div class="flex flex-col grow gap-1.5">
                        <Typography
                          preset="Text sm/Medium"
                          text={$_("referencePrompt")}
                          color="secondary"
                        />
                        <TextareaRw bind:value={data.reference_prompt} onEdit={editRules[index]} />
                      </div>

                      <div class="flex flex-col grow gap-1.5">
                        <Typography
                          preset="Text sm/Medium"
                          text={$_("sampleOutputPrompt")}
                          color="secondary"
                        />
                        <TextareaRw bind:value={data.sample} onEdit={editRules[index]} />
                      </div>
                    </div>
                  </div>
                  <Divider />
                {/each}
                {#if datas.data.length > 1}
                  <Button
                    class="w-auto self-center"
                    onclick={() => {
                      datas.data = [
                        ...datas.data,
                        {
                          id: datas.data.length,
                          title: "",
                          prompt: "",
                          reference_prompt: "",
                          sample: "",
                          createdAt: "",
                        },
                      ];
                      editRules = [...editRules, true];
                    }}
                  >
                    <Typography preset="Text sm/Semibold" text={$_("addNewSection")} />
                  </Button>
                {/if}
              </div>
            {:else}
              <div class="pt-8">
                <StateCard
                  illustration={IllustrationEmpty}
                  actionText={$_("addNewSection")}
                  funcButton={() => {
                    datas.data = [
                      ...datas.data,
                      {
                        id: datas.data.length,
                        title: "",
                        prompt: "",
                        reference_prompt: "",
                        sample: "",
                        createdAt: "",
                      },
                    ];
                    editRules = [...editRules, true];
                  }}
                  description={$_("youDontHaveAnyDataYet")}
                />
              </div>
            {/if}
          </div>
        {/if}
      </div>
    </div>
  </div>
</div>

<Modal title={$_("addNewReferenceDocument")} bind:open={addRefDocModal} size="lg" outsideclose>
  <div class="flex flex-col gap-4">
    <div class="flex flex-col gap-1.5">
      <Typography preset="Text sm/medium" color="secondary" text={`${$_("documentName")}*`} />
      <Input
        bind:value={tmpDocumentName}
        placeholder={$_("name")}
        class="bg-primary-light dark:bg-gray-800"
      />
    </div>
    <UploadFile
      title={`${$_("uploadDocuments")}*`}
      handleFileUploadWrapper={() => handleFileUploadWrapper(tmpDocumentName)}
      file={tmpFile}
      actionButton={true}
      actionDelete={() => (tmpFile = null)}
      actionView={() => ((documentViewModal = true), (selectedFile = tmpFile))}
      propClass="bg-primary-light dark:bg-gray-800 border-gray-300 dark:border-gray-500"
    />
    {#if tmpFile?.base64 !== undefined}
      <div class="flex gap-6">
        <div class="flex flex-col grow gap-6 rounded-lg bg-secondary-light dark:bg-secondary-dark">
          {#if tmpFile?.ext === ".pdf"}
            <div class="flex justify-center">
              <PdfjsViewer pdfFile={tmpFile.base64} />
            </div>
          {/if}
        </div>
      </div>
    {/if}
  </div>
  <svelte:fragment slot="footer">
    <Button
      class="w-full"
      on:click={() => {
        addOtherDocument();
      }}
      disabled={tmpFile === null || !tmpDocumentName}
    >
      {$_("addNewCase")}
    </Button>
  </svelte:fragment>
</Modal>

<Modal title={$_("documentViewer")} bind:open={documentViewModal} autoclose outsideclose>
  {#if selectedFile && (selectedFile.name.includes(".pdf") || selectedFile.name.includes(".PDF"))}
    <div class="flex justify-center">
      <PdfjsViewer pdfFile={selectedFile.base64} />
    </div>
  {/if}
</Modal>

<Modal bind:open={deleteModal} size="xs" autoclose>
  <div class="text-center">
    <ExclamationCircleSolid class="mx-auto mb-4 text-red-600 w-12 h-12" />
    <h3 class="mb-5 text-lg font-normal text-gray-500 dark:text-gray-400">
      {$_("areYouSureYouWantToDelete")}
    </h3>
    <Button
      onclick={() => {
        if (caseDelete) {
          deleteCase();
        } else {
          if (tabDatas.selected === "Rules") {
            deleteSection();
          } else {
            deleteDocument();
          }
        }
      }}
      color="red"
      class="me-2">{$_("yesImSure")}</Button
    >
    <Button
      color="alternative"
      onclick={() => {
        deleteModal = false;
        caseDelete = false;
      }}>{$_("noCancel")}</Button
    >
  </div>
</Modal>
