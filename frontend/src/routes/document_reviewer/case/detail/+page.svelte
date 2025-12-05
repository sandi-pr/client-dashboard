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
  import { Modal, Button, Input, MultiSelect, Checkbox } from "flowbite-svelte";
  import { page } from "$app/stores";
  import { goto } from "$app/navigation";
  import { ExclamationCircleSolid } from "flowbite-svelte-icons";
  import { _ } from "svelte-i18n";
  import { format } from "date-fns";
  import { call, check, getIcons, triggerToast, triggerLoading, errorHandling } from "$lib";
  import {
    InputRw,
    Divider,
    TabDetail,
    StateCard,
    Typography,
    UploadFile,
    TextareaRw,
    PdfjsViewer,
    Dropdown,
  } from "$lib";
  import { DataPresetCase } from "$lib/utils/DataPresetCase";

  const PlusIcon = getIcons("plus");
  const CopyIcon = getIcons("copy03");
  const Edit05Icon = getIcons("edit05");
  const Check2Icon = getIcons("check2");
  const PlusCircleIcon = getIcons("plusCircle");
  const IllustrationEmpty = getIcons("illustrationEmpty");

  const isFrontend = process.env.NODE_ENV === "development";

  let init = $state(false);
  let caseDelete = $state(false);
  let deleteModal = $state(false);
  let isEditTitle = $state(false);
  let addRefDocModal = $state(false);
  let documentViewModal = $state(false);
  let workflowPresetOpen = $state(false);
  let workflowAIEngineOpen = $state(false);

  let titleInput = $state("");
  let workflowPresetSelected = $state([""]);

  let tabDatas = $state({
    selected: "Referenced Documents",
    data: [
      check("aiManagement.contractreviewer.document:create_delete") && {
        name: "Referenced Documents",
      },
      check("aiManagement.contractreviewer.rule:create_delete") && {
        name: "Rules",
      },
      {
        name: "Workflow",
      },
    ].filter(Boolean),
  });
  let datas = $state({
    title: "",
    createdAt: "",
    data: [],
  });

  let editRules = $state([]);
  let editWorkflow = $state([]);
  let datasDataDuplicate = $state([]);
  let datasDataWorkflowDuplicate = $state([]);
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
  let workflowPresetOptions = [
    // {
    //   value: "",
    //   name: "No Preset",
    // },
    {
      value: "StartAgent",
      name: "Start Agent",
    },
    {
      value: "IteratorAgent",
      name: "Iterator Agent",
    },
    {
      value: "CompletenessAgent",
      name: "Completeness Agent",
    },
    {
      value: "ComparisonAgent",
      name: "Comparison Agent",
    },
    {
      value: "Kondisi00RouterAgent",
      name: "Kondisi 00 Router Agent",
    },
    {
      value: "Kolektibilitas2CheckerAgent",
      name: "Kolektibilitas 2 Checker Agent",
    },
    {
      value: "Kolektibilitas3CheckerAgent",
      name: "Kolektibilitas 3 Checker Agent",
    },
    {
      value: "Kolektibilitas4CheckerAgent",
      name: "Kolektibilitas 4 Checker Agent",
    },
    {
      value: "Kolektibilitas5CheckerAgent",
      name: "Kolektibilitas 5 Checker Agent",
    },
    {
      value: "Kode02CheckerAgent",
      name: "Kode 02 Checker Agent",
    },
    {
      value: "Kode03CheckerAgent",
      name: "Kode 03 Checker Agent",
    },
    {
      value: "Kode04CheckerAgent",
      name: "Kode 04 Checker Agent",
    },
    {
      value: "Kode05PlusCheckerAgent",
      name: "Kode 05 Plus Checker Agent",
    },
  ];
  let workflowAIEngineOptions = [
    { value: "gpt-oss", name: "GPT OSS" },
    { value: "lamma-3.1.8b", name: "Lamma 3.1 8B" },
  ];
  let workflowAIToolsOptions = [
    {
      value: "fetch_next_record",
      name: "Fetch Next Record",
    },
    {
      value: "run_completeness_check",
      name: "Run Completeness Check",
    },
    {
      value: "run_comparison_check",
      name: "Run Comparison Check",
    },
    {
      value: "run_checker",
      name: "Run Checker",
    },
  ];
  let workflowHandoffToOptions = [
    {
      value: "IteratorAgent",
      name: "Iterator Agent",
    },
    {
      value: "CompletenessAgent",
      name: "Completeness Agent",
    },
    {
      value: "FinalAgent",
      name: "Final Agent",
    },
    {
      value: "ComparisonAgent",
      name: "Comparison Agent",
    },
    {
      value: "AccuracyRouterAgent",
      name: "Accuracy Router Agent",
    },
  ];

  let tmpId = $state(null);
  let ruleId = $state(null);
  let tmpFile = $state(null);
  let selectedFile = $state(null);
  let tmpSectionId = $state(null);
  let collectionId = $state(null);
  let tmpDocumentName = $state(null);

  async function generateAllWorkflowAgents() {
    datas.dataWorkflow = workflowPresetOptions.map((preset, index) => ({
      id: index,
      preset: preset.value,
      name: "",
      description: "",
      aiEngine: "",
      aiTools: [],
      handoffTo: [],
      systemPrompt: "",
      setAsRootAgent: false,
    }));

    // Buat flag edit untuk semua
    editWorkflow = workflowPresetOptions.map(() => true);

    // Pilihan preset sesuai preset yang di-generate
    workflowPresetSelected = workflowPresetOptions.map((p) => p.value);
  }

  const getCollections = async (name, id) => {
    try {
      let token = localStorage.getItem("token");
      let responseCollection = await $call.call("client-dashboard", "get-collection-by-name", {
        URLQuery: {},
        URLParams: {},
        Params: {
          collection_name: id + ". " + name + " (Referenced Documents)",
        },
        Headers: {
          Authorization: "Bearer " + token,
        },
      });

      if (responseCollection.status !== 200) {
        const error = new Error("Failed to get collections");
        error.responseCollection = responseCollection;
        error.status = responseCollection.status;
        throw error;
      }

      if (responseCollection.result.data.length === 0) {
        let responseCreateCollection = await $call.call(
          "client-dashboard",
          "create-document-collection",
          {
            URLQuery: {},
            URLParams: {},
            Params: {
              collection_name: id + ". " + name + " (Referenced Documents)",
              description: "Collection for " + name,
            },
            Headers: {
              Authorization: "Bearer " + token,
            },
          },
        );

        if (responseCreateCollection.status !== 200) {
          const error = new Error("Failed to create collection");
          error.response = responseCreateCollection;
          error.status = responseCreateCollection.status;
          throw error;
        }

        collectionId = responseCreateCollection.result.data.id;
        return;
      }

      collectionId = responseCollection.result.data[0].id;
    } catch (error) {
      errorHandling($_, error, "get");
    }
  };

  const getRulesSections = async () => {
    try {
      let token = localStorage.getItem("token");
      if (token) {
        const responseRule = await $call.call("client-dashboard", "get-contract-rule", {
          URLQuery: {},
          URLParams: {
            id: ruleId.toString(),
          },
          Params: {},
          Headers: {
            Authorization: "Bearer " + token,
          },
        });

        if (responseRule.status !== 200) {
          const error = new Error("Failed to get rules");
          error.response = responseRule;
          error.status = responseRule.status;
          throw error;
        }

        await getCollections(
          responseRule.result.data[0].attributes.rule_name,
          responseRule.result.data[0].id,
        );

        let sectionsTmp = [];
        let page = 1;

        while (true) {
          const responseRuleSection = await $call.call(
            "client-dashboard",
            "get-contract-rule-section",
            {
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
            },
          );

          if (responseRuleSection.status !== 200) {
            const error = new Error("Failed to get rule sections");
            error.response = responseRuleSection;
            error.status = responseRuleSection.status;
            throw error;
          }

          if (responseRuleSection.result.data.length < 1) {
            break;
          }

          page += 1;
          sectionsTmp = sectionsTmp.concat(responseRuleSection.result.data);
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
          title: responseRule.result.data[0].attributes.rule_name,
          createdAt: format(
            responseRule.result.data[0].attributes.createdAt,
            "yyyy-MM-dd HH:mm:ss",
          ),
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
          dataWorkflow: [],
        };

        titleInput = datas.title;
      }
    } catch (error) {
      errorHandling($_, error, "get", "case");
      triggerToast("error");
    } finally {
      triggerLoading(false);
    }
  };

  const updateRulesSection = async () => {
    try {
      let token = localStorage.getItem("token");

      if (!token) {
        throw new Error("No token found");
      }

      const response = await $call.call("client-dashboard", "update-contract-rules", {
        URLQuery: {},
        URLParams: {
          id: ruleId,
        },
        Params: {
          rule_name: titleInput,
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

      await getReferences();
      await getRulesSections();
      triggerToast("updated");
    } catch (error) {
      console.log("error", error);
      triggerToast("error");
    }
  };

  const getReferences = async () => {
    try {
      let token = localStorage.getItem("token");
      let industry_id = sessionStorage.getItem("industry_id");

      if (industry_id === null) {
        console.log("no industry id");
        triggerLoading(false);
        return;
      }

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
          const error = new Error("Failed to get references document");
          error.response = response;
          error.status = response.status;
          throw error;
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
    } catch (error) {
      errorHandling($_, error, "get", "reference");
    }
  };

  const createRulesSections = async (data) => {
    try {
      triggerLoading(true);
      let token = localStorage.getItem("token");
      if (token) {
        const response = await $call.call("client-dashboard", "create-contract-rule-section", {
          URLQuery: {},
          URLParams: {},
          Params: {
            contract_rule_id: ruleId.toString(),
            section_title: data.title,
            section_prompt: data.prompt,
            reference_prompt: data.reference_prompt,
            // reference_prompt: ".",
            section_sample: data.sample,
          },
          Headers: {
            Authorization: "Bearer " + token,
          },
        });

        if (response.status !== 200) {
          const error = new Error("Failed to create section");
          error.response = response;
          error.status = response.status;
          throw error;
        }

        triggerToast("created");
      }
    } catch (error) {
      errorHandling($_, error, "create", "case");
    } finally {
      triggerLoading(false);
    }
  };

  const createWorkflowSections = async (data) => {
    try {
      console.log("not yet implemented");
      console.log("dataWorkflow", data);
      // triggerLoading(true);
      // let token = localStorage.getItem("token");
      // if (token) {
      //   const response = await $call.call("client-dashboard", "create-contract-rule-section", {
      //     URLQuery: {},
      //     URLParams: {},
      //     Params: {
      //       contract_rule_id: ruleId.toString(),
      //       section_title: data.title,
      //       section_prompt: data.prompt,
      //       reference_prompt: data.reference_prompt,
      //       // reference_prompt: ".",
      //       section_sample: data.sample,
      //     },
      //     Headers: {
      //       Authorization: "Bearer " + token,
      //     },
      //   });

      //   if (response.status !== 200) {
      //     const error = new Error("Failed to create section");
      //     error.response = response;
      //     error.status = response.status;
      //     throw error;
      //   }

      //   triggerToast("created");
      // }
    } catch (error) {
      errorHandling($_, error, "create", "case");
    } finally {
      triggerLoading(false);
    }
  };

  const addOtherDocument = async () => {
    try {
      triggerLoading(true);
      let token = localStorage.getItem("token");
      let industry_id = sessionStorage.getItem("industry_id");

      addRefDocModal = false;

      if (industry_id === null) {
        triggerLoading(false);
        return;
      }

      tmpFile.type = tmpDocumentName;

      const responseReference = await $call.call("client-dashboard", "create-references-document", {
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

      if (responseReference.status !== 200) {
        const error = new Error("Failed to create reference document");
        error.response = responseReference;
        error.status = responseReference.status;
        throw error;
      }

      tmpFile.id = responseReference.result.data[0].id;

      tmpFile.document_title = tmpDocumentName;
      tmpFile.description = "Document reference " + datas.title;
      tmpFile.contract = "REFERENCE";

      const responseDocument = await $call.call("client-dashboard", "create-document", {
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

      if (responseDocument.status !== 200) {
        const error = new Error("Failed to create document");
        error.response = responseDocument;
        error.status = responseDocument.status;
        throw error;
      }

      files.push(tmpFile);

      tmpFile = null;
      tmpDocumentName = "";
      triggerToast("updated");
    } catch (error) {
      errorHandling($_, error, "create", "reference");
    } finally {
      triggerLoading(false);
    }
  };

  const editSectionRules = async () => {
    try {
      triggerLoading(true);
      init = false;
      let token = localStorage.getItem("token");
      let tmpSection = datas.data.find((item) => item.id === tmpSectionId);
      if (token) {
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
          const error = new Error("Failed to update section");
          error.response = response;
          error.status = response.status;
          throw error;
        }

        triggerToast("updated");
        await getRulesSections();
      }
    } catch (error) {
      triggerLoading(false);
      errorHandling($_, error, "update", "case");
    }
  };

  const editSectionWorkflow = async () => {
    try {
      console.log("not yet implemented");
      // triggerLoading(true);
      // init = false;
      // let token = localStorage.getItem("token");
      // let tmpSection = datas.data.find((item) => item.id === tmpSectionId);
      // if (token) {
      //   const response = await $call.call("client-dashboard", "update-contract-rule-section", {
      //     URLQuery: {},
      //     URLParams: {
      //       id: tmpSection.id.toString(),
      //     },
      //     Params: {
      //       section_title: tmpSection.title,
      //       section_prompt: tmpSection.prompt,
      //       reference_prompt: tmpSection.reference_prompt,
      //       section_sample: tmpSection.sample,
      //     },
      //     Headers: {
      //       Authorization: "Bearer " + token,
      //     },
      //   });

      //   if (response.status !== 200) {
      //     const error = new Error("Failed to update section");
      //     error.response = response;
      //     error.status = response.status;
      //     throw error;
      //   }

      //   triggerToast("updated");
      //   await getRulesSections();
      // }
    } catch (error) {
      triggerLoading(false);
      errorHandling($_, error, "update", "case");
    }
  };

  const deleteSection = async () => {
    try {
      triggerLoading(true);
      deleteModal = false;
      let token = localStorage.getItem("token");
      if (token) {
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
          const error = new Error("Failed to delete section");
          error.response = response;
          error.status = response.status;
          throw error;
        }

        triggerToast("deleted");
        await getRulesSections();
      }
    } catch (error) {
      triggerLoading(false);
      errorHandling($_, error, "delete", "case");
    }
  };

  const deleteDocument = async () => {
    try {
      triggerLoading(true);
      deleteModal = false;

      let token = localStorage.getItem("token");

      const responseDocument = await $call.call("client-dashboard", "delete-references-document", {
        URLQuery: {},
        URLParams: {
          id: tmpId.toString(),
        },
        Params: {},
        Headers: {
          Authorization: "Bearer " + token,
        },
      });

      if (responseDocument.status !== 200) {
        const error = new Error("Failed to delete reference document");
        error.response = responseDocument;
        error.status = responseDocument.status;
        throw error;
      }

      const responseDelete = await $call.call(
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

      if (responseDelete.status !== 200) {
        const error = new Error("Failed to delete document");
        error.response = responseDelete;
        error.status = responseDelete.status;
        throw error;
      }

      files = files.filter((file) => file.id !== tmpId);
      triggerToast("deleted");
    } catch (error) {
      errorHandling($_, error, "delete", "reference");
    } finally {
      triggerLoading(false);
    }
  };

  const deleteCase = async () => {
    try {
      let token = localStorage.getItem("token");
      if (token) {
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
          const error = new Error("Failed to delete case");
          error.response = response;
          error.status = response.status;
          throw error;
        }

        triggerToast("deleted");
        goto(
          isFrontend
            ? "/document_reviewer/case"
            : await $call.getPage("client-dashboard", "contract-case"),
        );
      }
    } catch (error) {
      triggerLoading(false);
      errorHandling($_, error, "delete", "case");
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

  const checkSectionInputs = (data, type) => {
    if (type === "workflow") {
      return (
        data.name === "" ||
        data.description === "" ||
        data.aiEngine === "" ||
        data.handoffTo.length === 0 ||
        data.systemPrompt === ""
      );
    } else {
      return (
        data.title === "" ||
        data.prompt === "" ||
        data.sample === "" ||
        data.reference_prompt === ""
      );
    }
  };

  function stop(e) {
    e.preventDefault?.();
    e.stopPropagation?.();
  }

  onMount(async () => {
    if (
      (await check("aiManagement.contractreviewer.rule:create_delete")) === false &&
      (await check("aiManagement.contractreviewer.document:create_delete")) === false
    ) {
      console.log("no access");
      window.history.back();
      return;
    }

    ruleId = $page.url.searchParams.get("caseId");
    triggerLoading(true);
    if (ruleId !== null) {
      await getReferences();
      await getRulesSections();
      await generateAllWorkflowAgents();
    } else {
      triggerToast("error");
      triggerLoading(false);
      goto(
        isFrontend
          ? "/document_reviewer/case"
          : await $call.getPage("client-dashboard", "contract-case"),
      );
    }
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
            <div class="flex gap-2 items-center">
              {#if !isEditTitle}
                <Typography preset="Display sm/Semibold" text={datas.title} />
                <button
                  class="p-2 border border-primary-light dark:border-primary-dark rounded-lg bg-primary-light dark:bg-primary-dark hover:bg-secondary-light dark:hover:bg-secondary-dark"
                  onclick={() => {
                    navigator.clipboard.writeText(datas.title);
                    triggerToast("copyToClipboard");
                  }}
                >
                  <CopyIcon svgClass="w-4 h-4 text-secondary-light dark:text-secondary-dark" />
                </button>
                <button
                  class="p-2 border border-primary-light dark:border-primary-dark rounded-lg bg-primary-light dark:bg-primary-dark hover:bg-secondary-light dark:hover:bg-secondary-dark"
                  onclick={() => {
                    isEditTitle = true;
                  }}
                >
                  <Edit05Icon svgClass="w-4 h-4 text-secondary-light dark:text-secondary-dark" />
                </button>
              {:else}
                <div class="flex">
                  <InputRw bind:value={titleInput} />
                </div>
                <button
                  class="flex items-center justify-center p-2 border border-primary-light dark:border-primary-dark rounded-lg bg-primary-light dark:bg-primary-dark hover:bg-secondary-light dark:hover:bg-secondary-dark"
                  onclick={() => {
                    isEditTitle = false;
                    titleInput = datas.title;
                  }}
                >
                  <Typography
                    preset="Text sm/Regular"
                    text="X"
                    color="secondary"
                    propClass="w-4 h-4 text-center"
                  />
                </button>
                <button
                  class="p-2 border border-primary-light dark:border-primary-dark rounded-lg bg-primary-light dark:bg-primary-dark hover:bg-secondary-light dark:hover:bg-secondary-dark"
                  onclick={() => {
                    isEditTitle = false;
                    updateRulesSection();
                  }}
                >
                  <Check2Icon svgClass="w-4 h-4 text-secondary-light dark:text-secondary-dark" />
                </button>
              {/if}
            </div>
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
                outline
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
                            class={checkSectionInputs(data, "rules") ? "cursor-not-allowed" : ""}
                            onclick={() => {
                              editRules[index] = false;
                              if (datasDataDuplicate[index] === undefined) {
                                createRulesSections(data);
                              } else {
                                tmpSectionId = data.id;
                                editSectionRules();
                              }
                            }}
                            disabled={checkSectionInputs(data, "rules")}
                          >
                            <Typography
                              preset="Text sm/Semibold"
                              text={$_("save")}
                              color={checkSectionInputs(data, "rules")
                                ? "tertiary"
                                : "fg-brand-primary"}
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
                    class="w-auto self-center gap-1"
                    color="primary"
                    outline
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
        {:else if tabDatas.selected === "Workflow"}
          <div>
            <!-- header -->
            <div class="flex justify-between pt-3 pb-6">
              <div class="flex flex-col">
                <Typography
                  preset="Text lg/Semibold"
                  color="primary"
                  text={$_("setWorkflowForAIAgents")}
                />
                <Typography
                  preset="Text sm/Regular"
                  color="tertiary"
                  text={$_("descSetWorkflowForAIAgents")}
                />
              </div>
              <Button
                class="self-center gap-1"
                color="primary"
                outline
                on:click={() => {
                  datas.dataWorkflow = [
                    ...datas.dataWorkflow,
                    {
                      id: datas.dataWorkflow.length,
                      preset: "",
                      name: "",
                      description: "",
                      aiEngine: "",
                      aiTools: [],
                      handoffTo: [],
                      systemPrompt: "",
                      setAsRootAgent: false,
                    },
                  ];
                  editWorkflow = [...editWorkflow, true];
                  workflowPresetSelected = [...workflowPresetSelected, ""];
                }}
              >
                <PlusIcon svgClass="w-4 h-4" />
                {$_("addNewAIAgent")}
              </Button>
            </div>
            <!-- table latest -->
            {#if datas?.dataWorkflow?.length > 0 && workflowPresetSelected.length > 0}
              <div class="flex flex-col gap-6">
                {#each datas.dataWorkflow as data, index}
                  <div class="flex gap-4">
                    <div class="flex flex-col w-60 gap-2">
                      <Typography
                        preset="Text md/Semibold"
                        text={$_("section") + " " + (data.id === 0 ? "1" : index + 1)}
                        color="secondary"
                      />
                      <div class="flex flex-col gap-2">
                        {#if !editWorkflow[index]}
                          {#if datas.dataWorkflow.length > 1 && !init}
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
                              editWorkflow[index] = true;
                            }}
                          >
                            <Typography
                              preset="Text sm/Semibold"
                              text={$_("edit")}
                              color="fg-brand-primary"
                            />
                          </button>
                        {:else}
                          <button
                            onclick={() => {
                              if (datasDataWorkflowDuplicate[index] !== undefined) {
                                datas.dataWorkflow[index] = datasDataWorkflowDuplicate[index];
                                editWorkflow[index] = false;
                              } else {
                                datas.dataWorkflow.splice(index, 1);
                                editWorkflow.splice(index, 1);
                              }
                            }}
                          >
                            <Typography
                              preset="Text sm/Semibold"
                              text={$_("cancel")}
                              color="tertiary"
                            />
                          </button>
                          <button
                            class={checkSectionInputs(data, "workflow") ? "cursor-not-allowed" : ""}
                            onclick={() => {
                              editWorkflow[index] = false;
                              if (datasDataWorkflowDuplicate[index] === undefined) {
                                createWorkflowSections(data);
                              } else {
                                tmpSectionId = data.id;
                                editSectionWorkflow();
                              }
                            }}
                            disabled={checkSectionInputs(data, "workflow")}
                          >
                            <Typography
                              preset="Text sm/Semibold"
                              text={$_("save")}
                              color={checkSectionInputs(data, "workflow")
                                ? "tertiary"
                                : "fg-brand-primary"}
                            />
                          </button>
                        {/if}
                      </div>
                    </div>
                    <div class="flex flex-col grow gap-3">
                      <div class="flex flex-col gap-1.5">
                        <Typography preset="Text sm/Medium" text={$_("preset")} color="secondary" />
                        <Dropdown
                          dropdownOpen={workflowPresetOpen}
                          bind:selectedItem={workflowPresetSelected[index]}
                          items={workflowPresetOptions}
                          isOverlay={true}
                          placeholder={$_("selectPreset")}
                          onClick={() => {
                            if (workflowPresetSelected[index] !== "") {
                              const dataPreset = DataPresetCase(workflowPresetSelected[index]);
                              data.preset = workflowPresetSelected[index];
                              data.name = dataPreset.name;
                              data.description = dataPreset.description;
                              data.aiEngine = "gpt-oss";
                              data.aiTools = [...dataPreset.tools];
                              data.handoffTo = [...dataPreset.handoff_to];
                              data.systemPrompt = dataPreset.system_prompt;
                              if (dataPreset.name === "StartAgent") {
                                console.log("trigger");
                                data.setAsRootAgent = true;
                              } else {
                                data.setAsRootAgent = false;
                              }
                              editWorkflow[index] = false;
                            } else {
                              data.preset = "";
                              data.name = "";
                              data.description = "";
                              data.aiEngine = "";
                              data.aiTools = [];
                              data.handoffTo = [];
                              data.systemPrompt = "";
                              data.setAsRootAgent = false;
                              editWorkflow[index] = true;
                            }
                          }}
                        />
                      </div>
                      <div class="flex flex-col grow gap-1.5">
                        <Typography
                          preset="Text sm/Medium"
                          text={$_("name") + "*"}
                          color="secondary"
                        />
                        <InputRw bind:value={data.name} onEdit={editWorkflow[index]} />
                      </div>
                      <div class="flex flex-col grow gap-1.5">
                        <Typography
                          preset="Text sm/Medium"
                          text={$_("description") + "*"}
                          color="secondary"
                        />
                        <TextareaRw bind:value={data.description} onEdit={editWorkflow[index]} />
                      </div>
                      <div class="flex flex-col gap-1.5">
                        <Typography
                          preset="Text sm/Medium"
                          text={$_("AIEngine") + "*"}
                          color="secondary"
                        />
                        <Dropdown
                          dropdownOpen={workflowAIEngineOpen}
                          bind:selectedItem={data.aiEngine}
                          items={workflowAIEngineOptions}
                          isOverlay={true}
                          placeholder={$_("selectAIEngine")}
                          disabled={!editWorkflow[index]}
                        />
                      </div>
                      <div class="flex flex-col gap-1.5">
                        <Typography
                          preset="Text sm/Medium"
                          text={$_("AITools")}
                          color="secondary"
                        />
                        <div class="relative">
                          <MultiSelect
                            items={workflowAIToolsOptions}
                            bind:value={data.aiTools}
                            size="lg"
                            class="text-primary-light dark:text-primary-dark"
                          />
                          {#if !editWorkflow[index]}
                            <div
                              class="absolute inset-0 z-20 bg-transparent cursor-not-allowed"
                              onclick={stop}
                              onmousedown={stop}
                              onmouseup={stop}
                              onkeydown={stop}
                              aria-hidden="true"
                            />
                          {/if}
                        </div>
                      </div>
                      <div class="flex flex-col gap-1.5">
                        <Typography
                          preset="Text sm/Medium"
                          text={$_("AIEngine") + "*"}
                          color="secondary"
                        />
                        <div class="relative">
                          <MultiSelect
                            items={workflowHandoffToOptions}
                            bind:value={data.handoffTo}
                            disabled={!editWorkflow[index]}
                            size="lg"
                            class="text-primary-light dark:text-primary-dark"
                          />
                          {#if !editWorkflow[index]}
                            <div
                              class="absolute inset-0 z-20 bg-transparent cursor-not-allowed"
                              onclick={stop}
                              onmousedown={stop}
                              onmouseup={stop}
                              onkeydown={stop}
                              aria-hidden="true"
                            />
                          {/if}
                        </div>
                      </div>
                      <div class="flex flex-col grow gap-1.5">
                        <Typography
                          preset="Text sm/Medium"
                          text={$_("systemPrompt") + "*"}
                          color="secondary"
                        />
                        <TextareaRw bind:value={data.systemPrompt} onEdit={editWorkflow[index]} />
                        <Typography
                          preset="Text sm/Regular"
                          text={$_("descSystemPrompt")}
                          color="tertiary"
                        />
                      </div>
                      <div class="flex flex-col grow gap-1.5">
                        <div class="flex gap-2">
                          <Checkbox
                            bind:value={data.setAsRootAgent}
                            checked={data.setAsRootAgent}
                            disabled={!editWorkflow[index]}
                            class="disabled:hover:cursor-not-allowed"
                          />
                          <Typography
                            preset="Text sm/Regular"
                            text={$_("rootAgent")}
                            color="primary"
                          />
                        </div>
                        <Typography
                          preset="Text sm/Regular"
                          text={$_("descRootAgent")}
                          color="tertiary"
                        />
                      </div>
                    </div>
                  </div>
                  <Divider />
                {/each}
                {#if datas.data.length > 1}
                  <Button
                    class="w-auto self-center gap-1"
                    color="primary"
                    outline
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
                      editWorkflow = [...editWorkflow, true];
                    }}
                  >
                    <PlusIcon svgClass="w-4 h-4" />
                    {$_("addNewSection")}
                  </Button>
                {/if}
              </div>
            {:else}
              <div class="pt-8">
                <StateCard
                  illustration={IllustrationEmpty}
                  actionText={$_("addNewSection")}
                  funcButton={() => {
                    datas.dataWorkflow = [
                      ...datas.dataWorkflow,
                      {
                        id: datas.dataWorkflow.length,
                        preset: "",
                        name: "",
                        description: "",
                        aiEngine: "",
                        aiTools: [],
                        handoffTo: [],
                        systemPrompt: "",
                        setAsRootAgent: false,
                      },
                    ];
                    editWorkflow = [...editWorkflow, true];
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
