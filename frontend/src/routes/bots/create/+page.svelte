<!--
  Copyright (C) 2024 - Hexacode TI
  All Rights Reserved.

  Author(s):
  * Dean Ramhan <dean@hexacode.org>
  * Taqie Fadlillah <taqief@hexacode.org>
  * Sandi Permana <sandy@hexacode.org>

  This file is part of the "Crossroads" project.
-->

<script>
  import { Button, Modal } from "flowbite-svelte";
  import { onMount } from "svelte";
  import { _ } from "svelte-i18n";
  import Divider from "$lib/components/Divider.svelte";
  import {
    call,
    SelectedCard,
    UploadFile,
    triggerToast,
    triggerLoading,
    Typography,
    InputRw,
    TextareaRw,
    StateCard,
    Dropdown,
  } from "$lib";
  import { findIconFromModelName } from "$lib";
  import DataTableListCollection from "$lib/components/DataTableListCollection.svelte";
  import { goto } from "$app/navigation";
  import { ChevronDownOutline } from "flowbite-svelte-icons";
  import IllustrationEmpty from "$lib/icon/illustration_empty.svelte";
  import IllustrationProgress from "$lib/icon/illustration_progress.svelte";

  const isFrontend = process.env.NODE_ENV === "development";

  let dropdownCollectionOpen = false;
  let collections = $state([]);
  let initCollections = $state([]);
  let documents = $state([]);

  let botName = $state("");
  let description = $state("");
  let fileImage = $state(null);
  let fileImageName = $state(null);
  let fileImageResult = $state(null);
  let selectedplatforms = $state(null);
  let selectedCollection = $state("");
  let oldSelectedCollection = $state("");
  let additionalPrompt = $state("");
  let needAuth = $state(false);
  let modalOpen = $state(false);
  let isSelectedLibrary = $state(false);
  let allData = $state([]);
  let searchTerm = $state("");
  let search = $state(false);
  // let page = $state(1);
  let checkAll = $state(false);
  let stateCheckAll = $state("");
  let datasChecked = $state([]);
  let tmpDatasChecked = $state({});
  let configs = $state({});
  let gettingDatas = $state(false);
  let pagination = $state({ page: 1, pageSize: 10, pageCount: 0, total: 0 });
  let paginationCollection = $state({ page: 1, pageSize: 20, pageCount: 0, total: 0 });
  let dropdownUsecaseOpen = $state(false);
  let usecase = $state([
    { name: "My Bots (Default)", value: "myBots" },
    { name: "Hexassistant", value: "hexassistant" },
    { name: "Coding Assistant", value: "codingAssistant" },
  ]);
  let usecaseSelected = $state("myBots");

  let platforms = $state([]);

  let selectedDays = [
    {
      name: "Monday",
      label: "M",
      selected: false,
    },
    {
      name: "Tuesday",
      label: "T",
      selected: false,
    },
    {
      name: "Wednesday",
      label: "W",
      selected: true,
    },
    {
      name: "Thursday",
      label: "T",
      selected: true,
    },
    {
      name: "Friday",
      label: "F",
      selected: true,
    },
    {
      name: "Saturday",
      label: "S",
      selected: false,
    },
    {
      name: "Sunday",
      label: "S",
      selected: false,
    },
  ];

  let screenWidth = $state(0);
  let params = [];
  let products_and_services = [];

  let templates = [
    {
      id: 1,
      attributes: {
        name: "Template 1",
        parameters: [
          {
            children: [
              {
                label: "Parameter 1",
                required: true,
                type: "text",
                placeholder: "Enter Parameter 1",
                enum: ["Option 1", "Option 2", "Option 3"],
                description: "Description for Parameter 1",
              },
            ],
          },
          {
            label: "Label 1",
            children: [
              {
                list_dict_format: [
                  {
                    label: "Field 1",
                    required: true,
                    type: "text",
                    placeholder: "Enter Field 1",
                    name: "field1",
                  },
                  {
                    label: "Field 2",
                    required: false,
                    type: "number",
                    placeholder: "Enter Field 2",
                    min: 0,
                    max: 100,
                    name: "field2",
                  },
                ],
              },
            ],
          },
        ],
      },
    },
    {
      id: 2,
      attributes: {
        name: "Template 2",
        parameters: [
          {
            children: [
              {
                label: "Parameter 2",
                required: false,
                type: "number",
                placeholder: "Enter Parameter 2",
                enum: ["Option A", "Option B", "Option C"],
                description: "Description for Parameter 2",
              },
            ],
          },
        ],
      },
    },
    {
      id: 3,
      attributes: {
        name: "Template 3",
        parameters: [
          {
            children: [
              {
                label: "Parameter 3",
                required: true,
                type: "text",
                placeholder: "Enter Parameter 3",
                enum: ["Option X", "Option Y", "Option Z"],
                description: "Description for Parameter 3",
              },
            ],
          },
        ],
      },
    },
  ];

  async function getCollectionSub(id) {
    let token = localStorage.getItem("token");
    if (token !== null) {
      if (!search) {
        gettingDatas = true;
      }
      let response = await $call.call("client-dashboard", "get-library-collections", {
        URLQuery: {},
        URLParams: {},
        Params: {
          page: pagination.page,
          page_size: pagination.pageSize,
          parent_id: id,
          search: searchTerm,
        },
        Headers: {
          Authorization: "Bearer " + token,
        },
      });

      if (response.status !== 200) {
        console.log(response);
        gettingDatas = false;
        return;
      }

      console.log("collections", response);

      if (response.result.data.length > 0) {
        pagination = response.result.meta.pagination;
        configs = {
          search: true,
          searchHandler: handleSearch,
          pagination: response.result.meta.pagination,
          space: "p-4",
          actionInLeft: true,
          actions: [
            {
              id: "deselect_all",
              name: "Deselect All",
              onclick: () => {
                handleCheckAllCollection(false);
              },
            },
            {
              id: "select_all",
              name: "Select All",
              type: "brand-inverse",
              onclick: () => {
                handleCheckAllCollection(true);
              },
            },
          ],
        };
        let promises = response.result.data.map(async (document) => {
          const sub_collections = await getSubCollections(document.id);
          return {
            id: document.id,
            collection_name: document.attributes.name,
            sub_collections,
          };
        });
        allData = await Promise.all(promises);
        // console.log("all data with subcollections", allData);

        if (checkAll) {
          datasChecked = allData.map((data) => {
            return {
              id: data.id,
              checked: true,
              sub_collection: data.sub_collections.map((sub) => ({
                id: sub.id,
                checked: true,
              })),
            };
          });
          // console.log("datasChecked all", datasChecked);
        } else {
          if (tmpDatasChecked[pagination.page] === undefined) {
            datasChecked = allData.map((data) => {
              return {
                id: data.id,
                checked: false,
                sub_collection: data.sub_collections.map((sub) => ({
                  id: sub.id,
                  checked: false,
                })),
              };
            });
          } else {
            datasChecked = tmpDatasChecked[pagination.page];
          }
          // console.log("datasChecked", datasChecked);
        }

        // if (checkAll) {
        //   datasChecked = datasChecked.map((data) => {
        //     return { id: data.id, checked: true };
        //   });
        // }
      } else {
        pagination = response.result.meta.pagination;
        allData = [];
        datasChecked = [];
        configs = {
          search: true,
          searchHandler: handleSearch,
          pagination: response.result.meta.pagination,
          space: "p-4",
          actionInLeft: true,
          actions: [
            {
              id: "deselect_all",
              name: "Deselect All",
              onclick: () => {
                handleCheckAllCollection(false);
              },
            },
            {
              id: "select_all",
              name: "Select All",
              type: "brand-inverse",
              onclick: () => {
                handleCheckAllCollection(true);
              },
            },
          ],
          // actions: [
          //   {
          //     id: "delete",
          //     name: "Delete",
          //     icon: "trash",
          //     type: "error",
          //     disabled: datasChecked.find((data) => data.checked) ? false : true,
          //     onclick: () => {
          //       if (datasChecked.find((data) => data.checked)) {
          //         popupModal = true;
          //       }
          //     },
          //   },
          //   {
          //     id: "add_new",
          //     name: "Add New",
          //     icon: "plus",
          //     type: "brand-inverse",
          //     onclick: () => {
          //       selectTypeCollection("Documents");
          //       addDocumentsModal = true;
          //       addOtherCategory = false;
          //     },
          //   },
          // ],
        };
      }

      gettingDatas = false;
      search = false;
    }
  }

  async function getSubCollections(id) {
    let token = localStorage.getItem("token");

    if (token !== null) {
      gettingDatas = true;
      const aggregatedData = [];

      paginationCollection.page = 1;

      do {
        let response = await $call.call("client-dashboard", "get-library-collections", {
          URLQuery: {},
          URLParams: {},
          Params: {
            page: paginationCollection.page,
            page_size: paginationCollection.pageSize,
            parent_id: parseInt(id),
          },
          Headers: {
            Authorization: "Bearer " + token,
          },
        });

        if (response.status !== 200) {
          console.log(response);
          gettingDatas = false;
          if (response.status === 400) {
            console.error("Invalid parameters:", response);
            triggerToast($_("invalidParameters"), true, "error");
          } else if (response.status === 401) {
            console.error("Unauthorized access:", response);
            triggerToast($_("unauthorized"), true, "error");
          } else if (response.status === 500) {
            console.error("Server error:", response);
            triggerToast($_("serverError"), true, "error");
          } else {
            console.error("Unexpected status code:", response);
            triggerToast(`${$_("unexpectedError")}: ${response.status}`, true, "error");
          }
          return;
        }

        aggregatedData.push(...response.result.data);

        paginationCollection = response.result.meta.pagination;

        paginationCollection.page++;
      } while (paginationCollection.page <= paginationCollection.pageCount);

      const promise = aggregatedData.map(async (collection) => ({
        id: collection.id,
        collection_name: collection.attributes.name,
      }));

      const subCollectionList = await Promise.all(promise);
      console.log("subCollectionList", subCollectionList);
      gettingDatas = false;
      return subCollectionList;
    } else {
      triggerToast($_("unauthorized"), true, "error");
      gettingDatas = false;
      return;
    }
  }

  const handleSearch = async (value) => {
    search = true;
    pagination.page = 1;
    searchTerm = value;

    await getCollectionSub(selectedCollection);
  };

  const handleInputFile = () => {
    const inputElement = document.createElement("input");
    inputElement.type = "file";
    inputElement.accept = "image/jpeg, image/png";

    inputElement.addEventListener("change", async (event) => {
      const file = event.target.files[0];

      if (file && (file.type === "image/jpeg" || file.type === "image/png")) {
        fileImageResult = await new Promise((resolve, reject) => {
          const reader = new FileReader();

          reader.onload = () => resolve(reader.result);
          reader.onerror = reject;

          reader.readAsDataURL(file);
        });

        fileImage = await readFileAsBase64(file);
        fileImageName = file.name;
      } else {
        alert("Invalid file. Please select a valid .jpg or .png file.");
      }
    });

    inputElement.click();
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

  let platformsReady = $state(false);

  async function getCollections() {
    let token = localStorage.getItem("token");

    if (token !== null) {
      let page = 1;
      initCollections = [];
      while (true) {
        let response = await $call.call("client-dashboard", "get-librarys", {
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

        if (response.result.data.length > 0) {
          initCollections = initCollections.concat(response.result.data);
          page++;
        } else {
          break;
        }
      }

      console.log("libraries", initCollections);

      collections = initCollections.map((collection) => {
        return {
          value: collection.id,
          id: collection.id,
          name: collection.attributes.name,
          description: collection.attributes.meta?.description,
        };
      });
    }
  }

  async function getModels() {
    let token = localStorage.getItem("token");
    if (token) {
      try {
        let response = await $call.call("hexa-ai-rag", "get-models", {
          URLQuery: {},
          URLParams: {},
          Params: {},
          Headers: {
            Authorization: "Bearer " + token,
          },
        });

        console.log("models", response);

        if (response === null) {
          if (isFrontend) {
            platforms = [
              {
                id: 1,
                name: "Llama3.1",
                model: "llama3.1:8b",
                description: "Llama3.1 Develop by Meta",
                create_at: "2024-07-23",
              },
            ];
            platformsReady = true;
          } else {
            triggerToast("error");
          }
          return;
        }

        if (response.status === 200) {
          // BRI ADJUSTMENT
          response.result.models = response.result.models.filter(
            (model) => !model.name.includes("Qurani"),
          );

          let i = 0;
          let promise = response.result.models.map((model) => {
            return {
              id: i++,
              name: model.name,
              model: model.model,
              description: model.description,
              create_at: model.create_at,
              icon: findIconFromModelName(model.name),
              recommended: model.name === "Llama3.1" ? true : false,
            };
          });

          platforms = await Promise.all(promise);

          platformsReady = true;
        } else {
          triggerToast("error");
          console.log(response);
        }
      } catch (error) {
        console.log(error);
      }
    }
  }

  async function createBot() {
    let token = localStorage.getItem("token");

    if (token !== null) {
      triggerLoading(true);

      if (platforms?.find((platform) => platform.id === selectedplatforms)?.model !== "openai") {
      let allCollection = "";

      // console.log("tmpDatasChecked", tmpDatasChecked);

      let pageNow = pagination.page;
      if (tmpDatasChecked[pageNow] === undefined) {
        tmpDatasChecked = {
          ...tmpDatasChecked,
          [pageNow]: datasChecked,
        };
      } else {
        tmpDatasChecked[pageNow] = datasChecked;
      }

      // Iterasi key dalam tmpDatasChecked
      Object.values(tmpDatasChecked).forEach((pageData) => {
        if (!pageData) {
          triggerLoading(false);
          return;
        } // Skip jika pageData undefined/null

        console.log("pageData", pageData);

        const collections = pageData
          .flatMap((item) => {
            // Ambil ID dari sub_collection jika checked = true
            const subCollectionIds =
              item.sub_collection?.filter((sub) => sub.checked).map((sub) => sub.id) || [];

            // Tambahkan item.id jika checked = true
            return item.checked ? [item.id, ...subCollectionIds] : subCollectionIds;
          })
          .join(","); // Gabungkan ID dengan koma

        console.log("collections check", collections);

        if (collections) {
          allCollection += collections + ",";
        }
      });

      // Hapus koma terakhir jika ada
      allCollection = allCollection.replace(/,$/, "");

      // console.log("All Collection IDs (checked = true):", allCollection);
      // return;

      // let collections2 = collections.split(",").map(Number);

      let dataSend = {
        bot_name: botName,
        description: description,
        ai_engine: platforms.find((platform) => platform.id === selectedplatforms).model,
        library_id: selectedCollection,
        collection_id: checkAll ? "all" : allCollection,
        additional_prompt: additionalPrompt,
        need_auth: needAuth,
      };
      if (fileImage !== "" && fileImage !== null && fileImage !== undefined) {
        dataSend.bot_logo = JSON.stringify({ name: fileImageName, base64Data: fileImage });
      }

      // console.log("dataSend", dataSend);

      try {
        const response = await $call.call("client-dashboard", "create-bot", {
          URLQuery: {},
          URLParams: {},
          Params: dataSend,
          Headers: {
            Authorization: "Bearer " + token,
          },
        });

        console.log(response);

        if (response.status !== 200) {
          console.log(response);
          triggerLoading(false);
          return;
        }

        goto(await $call.getPage("client-dashboard", "bots"));
        triggerLoading(false);
        triggerToast("created");
      } catch (error) {
        triggerLoading(false);
        console.log(error);
      }
    } else {

      const response = await $call.call("client-dashboard", "find-create-library-llm", {
          URLQuery: {},
          URLParams: {},
          Params: {},
          Headers: {
            Authorization: "Bearer " + token,
          },
        });

        if (response.status !== 200) {
          // console.log(response);
          return;
        }

        let collectionId = response.result.collection_id;
        let libraryId = response.result.library_id;

      let dataSend = {
        bot_name: botName,
        description: description,
        ai_engine: platforms.find((platform) => platform.id === selectedplatforms).model,
        library_id: parseInt(libraryId),
        collection_id: collectionId.toString(),
        additional_prompt: additionalPrompt,
        need_auth: needAuth,
      };

      if (fileImage !== "" && fileImage !== null && fileImage !== undefined) {
        dataSend.bot_logo = JSON.stringify({ name: fileImageName, base64Data: fileImage });
      }

      // console.log("dataSend", dataSend);

      try {
        const response = await $call.call("client-dashboard", "create-bot", {
          URLQuery: {},
          URLParams: {},
          Params: dataSend,
          Headers: {
            Authorization: "Bearer " + token,
          },
        });

        console.log(response);

        if (response.status !== 200) {
          console.log(response);
          triggerLoading(false);
          return;
        }

        goto(await $call.getPage("client-dashboard", "bots"));
        triggerLoading(false);
        triggerToast("created");
      } catch (error) {
        triggerLoading(false);
        console.log(error);
      }
    }
    } else {
      console.log("no token");
      return;
    }
  }

  async function handleCheckAll(allChecked) {
    datasChecked = datasChecked.map((data) => ({
      id: data.id,
      checked: allChecked ? false : true,
      sub_collection: data.sub_collection.map((sub) => ({
        id: sub.id,
        checked: allChecked ? false : true,
      })),
    }));

    // console.log("Updated datasChecked (handleCheckAll):", datasChecked);
  }

  async function handleCheckAllCollection(allChecked) {
    checkAll = allChecked;

    datasChecked = datasChecked.map((data) => ({
      id: data.id,
      checked: allChecked,
      sub_collection: data.sub_collection.map((sub) => ({
        id: sub.id,
        checked: allChecked,
      })),
    }));

    if (!allChecked) {
      // Jika allChecked adalah false, update tmpDatasChecked
      Object.keys(tmpDatasChecked).forEach((page) => {
        tmpDatasChecked[page] = tmpDatasChecked[page].map((data) => ({
          id: data.id,
          checked: false,
          sub_collection: data.sub_collection.map((sub) => ({
            id: sub.id,
            checked: false,
          })),
        }));
      });
    }

    // console.log("Updated datasChecked (handleCheckAllCollection):", datasChecked);
  }

  async function handleCheck(id, isSubCollection = false, parentId = null) {
    if (!Array.isArray(datasChecked)) {
      datasChecked = [];
    }

    datasChecked = datasChecked.map((data) => {
      // Jika parentId tidak ada, kita memproses level utama
      if (!isSubCollection && data.id === id) {
        if (data.checked) {
          checkAll = false;
        }
        return {
          ...data,
          checked: !data.checked,
        };
      }

      // console.log("checkAll", checkAll);

      // Jika parentId cocok, kita memproses sub_collection
      if (isSubCollection && data.id === parentId) {
        const updatedSubCollection = data.sub_collection.map((sub) => {
          if (sub.id === id) {
            return { ...sub, checked: !sub.checked };
          }
          return sub;
        });

        return {
          ...data,
          sub_collection: updatedSubCollection,
        };
      }

      return data;
    });

    let pageNow = pagination.page;
    if (tmpDatasChecked[pageNow] === undefined) {
      tmpDatasChecked = {
        ...tmpDatasChecked,
        [pageNow]: datasChecked,
      };
    } else {
      tmpDatasChecked[pageNow] = datasChecked;
    }
    // console.log("datasChecked setelah pembaruan:", datasChecked);
  }

  async function handleCheckSubCollection(parentId, subId) {
    datasChecked = datasChecked.map((data) => {
      if (data.id === parentId) {
        return {
          ...data,
          sub_collection: data.sub_collection.map((sub) => {
            if (sub.id === subId) {
              return { ...sub, checked: !sub.checked };
            }
            return sub;
          }),
        };
      }
      return data;
    });

    // console.log("Updated datasChecked (handleCheckSubCollection):", datasChecked);
  }

  function areAllUnchecked() {
    // Gabungkan data dari halaman sekarang dengan data yang tersimpan di tmpDatasChecked
    const allData = [...datasChecked, ...Object.values(tmpDatasChecked).flat()];

    // Cek apakah semua data tidak ada yang checked
    return allData.every((item) => item.checked === false);
  }

  async function next() {
    if (pagination.page >= pagination.pageCount) return;

    let pageNow = pagination.page;
    if (tmpDatasChecked[pageNow] === undefined) {
      tmpDatasChecked = {
        ...tmpDatasChecked,
        [pageNow]: datasChecked,
      };
    } else {
      tmpDatasChecked[pageNow] = datasChecked;
    }

    pagination.page = pagination.page + 1;
    // console.log("cek tmpDatasChecked", tmpDatasChecked);
    // getDocuments();
    getCollectionSub(selectedCollection);
  }

  async function previous() {
    if (pagination.page <= 1) return;

    let pageNow = pagination.page;
    if (tmpDatasChecked[pageNow] === undefined) {
      tmpDatasChecked = {
        ...tmpDatasChecked,
        [pageNow]: datasChecked,
      };
    } else {
      tmpDatasChecked[pageNow] = datasChecked;
    }

    pagination.page = pagination.page - 1;
    // getDocuments();
    getCollectionSub(selectedCollection);
  }

  onMount(() => {
    getCollections();
    getModels();
    // getCollectionSub();
  });

  $effect(() => {
    if (selectedCollection) {
      if (typeof selectedCollection === "number" && selectedCollection !== oldSelectedCollection) {
        oldSelectedCollection = selectedCollection;
        getCollectionSub(selectedCollection);
      }
    }
  });
</script>

<svelte:window bind:innerWidth={screenWidth} />

<div class="m-4 pb-12">
  <div class="flex justify-between">
    <Typography preset="Display sm/Semibold" text={$_("createNewBots")} />
    <div class="flex text-sm gap-3">
      <Button
        on:click={async () => {
          if (isFrontend) {
            goto("/bots");
          } else {
            goto(await $call.getPage("client-dashboard", "bots"));
          }
        }}
        color="light"
      >
        {$_("cancel")}
      </Button>
      <Button
        disabled={botName === "" || selectedplatforms === "" || platforms?.find((platform) => platform.id === selectedplatforms)?.model !== "openai" ? selectedCollection === "" : false}
        on:click={() => {
          createBot();
        }}
        class="bg-hc-brand-600 hover:bg-hc-brand-700 dark:bg-hc-brand-600 dark:hover:bg-hc-brand-700"
      >
        {$_("create")}
      </Button>
    </div>
  </div>

  <!-- Content -->
  <div class="mt-6 space-y-6">
    <div class="md:flex">
      <Typography
        preset="Text md/Semibold"
        text={$_("botsInformation")}
        propClass="w-full md:w-1/4"
      />
      <div class="w-full md:w-3/4 grid grid-cols-1 md:grid-cols-4 gap-3">
        <div class="flex flex-col gap-1.5 col-span-4">
          <Typography preset="Text sm/Medium" text={$_("botsName")} color="secondary" />
          <InputRw bind:value={botName} onEdit />
        </div>
        <div class="flex flex-col gap-1.5 col-span-4">
          <Typography preset="Text sm/Medium" text={$_("descriptionOptional")} color="secondary" />
          <TextareaRw bind:value={description} onEdit />
        </div>
        <div class="col-span-4">
          <UploadFile
            title={$_("botsProfilePicture")}
            handleFileUploadWrapper={() => handleInputFile()}
            isPreview
            srcAvatar={fileImageName ? fileImageResult : ""}
          />
        </div>
      </div>
    </div>
    <Divider />
    <div class="md:flex">
      <Typography
        preset="Text md/Semibold"
        text={$_("aiEngines")}
        color="secondary"
        propClass="w-full md:w-1/4"
      />
      <div class="w-full md:w-3/4">
        <SelectedCard
          selectionData={platforms}
          bind:selectedData={selectedplatforms}
          numGrid={3}
          isRadio
        />
      </div>
    </div>
    <!-- <Divider />
    <div class="md:flex">
      <Typography
        preset="Text md/Semibold"
        text={$_("useCase")}
        color="secondary"
        propClass="w-full md:w-1/4"
      />
      <div class="w-full md:w-auto">
        <Dropdown
          dropdownOpen={dropdownUsecaseOpen}
          bind:selectedItem={usecaseSelected}
          items={usecase}
          isOverlay={true}
        />
      </div>
    </div> -->
    {#if platformsReady}
    {#if platforms?.find((platform) => platform.id === selectedplatforms)?.model !== "openai"}
    <Divider />
    <div class="md:flex">
      <Typography
        preset="Text md/Semibold"
        text={$_("knowledgeSource")}
        color="secondary"
        propClass="w-full md:w-1/4"
      />
      <div class="w-full md:w-3/4">
        <div class="flex items-center gap-4">
          <Button
            on:click={() => (modalOpen = true)}
            color="light"
            class={`w-max rounded-lg border border-primary-light dark:border-primary-dark bg-primary-light dark:bg-primary-dark gap-1.5 flex justify-between items-center py-2.5 px-3`}
          >
            <Typography
              preset="Text md/Semibold"
              text={selectedCollection
                ? collections.find((item) => item.value === selectedCollection)?.name
                : $_("configure")}
              color="secondary"
            />
            <ChevronDownOutline
              class="w-3 h-3 ms-2 text-fg-quaternary-light dark:text-fg-quaternary-dark"
            />
          </Button>
        </div>
      </div>
    </div>
    {/if}
    {/if}
    
    <Divider />
    <div class="md:flex">
      <Typography
        preset="Text md/Semibold"
        text={$_("additionalPrompts")}
        color="secondary"
        propClass="w-full md:w-1/4"
      />
      <div class="w-full md:w-3/4 flex flex-col gap-1.5">
        <Typography preset="Text sm/Medium" text={$_("botsRole")} color="secondary" />
        <TextareaRw bind:value={additionalPrompt} onEdit rows={8} />
      </div>
    </div>
    <Divider />
    <div class="md:flex">
      <Typography
        preset="Text md/Semibold"
        text={$_("additionalSettings")}
        color="secondary"
        propClass="w-full md:w-1/4"
      />
      <div class="w-full md:w-3/4 grid grid-cols-1 md:grid-cols-4 gap-4">
        <div class="col-span-4">
          <label class="inline-flex items-center cursor-pointer">
            <input type="checkbox" value="" class="hidden peer" bind:checked={needAuth} />
            <div
              class="relative w-11 h-6 bg-tertiary-light dark:bg-tertiary-dark peer-focus:outline-none peer-focus:ring-4 peer-focus:ring-utility-brand-200-light dark:peer-focus:ring-utility-brand-200-dark rounded-full peer peer-checked:after:translate-x-full rtl:peer-checked:after:-translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:start-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all dark:border-gray-600 peer-checked:bg-brand-solid-light"
            ></div>
            <span class="ms-3 text-sm font-medium text-gray-900 dark:text-gray-300"
              >{$_("requiresLogin")}</span
            >
          </label>
        </div>
      </div>
    </div>
  </div>
</div>

<Modal bind:open={modalOpen} outsideclose>
  <svelte:fragment slot="header">
    <div class="flex gap-4 me-6 grow w-full">
      <div
        class="grow pt-2 border-t-4 border-fg-brand-primary-alt-light dark:border-fg-brand-primary-alt-dark"
      >
        <Typography
          preset="Text md/Bold"
          text={$_("chooseDocumentsLibrary")}
          color={!isSelectedLibrary ? "brand-secondary" : "secondary"}
        />
        <Typography
          preset="Text sm/Regular"
          text={$_("documentLibraryAsThePrimarySource")}
          color={!isSelectedLibrary ? "brand-tertiary" : "tertiary"}
        />
      </div>

      <div
        class={`grow pt-2 border-t-4 ${isSelectedLibrary ? "border-fg-brand-primary-alt-light dark:border-fg-brand-primary-alt-dark" : "border-utility-gray-200-light dark:border-utility-gray-200-dark"}`}
      >
        <Typography
          preset="Text md/Bold"
          text={$_("chooseWhichCollections")}
          color={isSelectedLibrary ? "brand-secondary" : "secondary"}
        />
        <Typography
          preset="Text sm/Regular"
          text={$_("chooseTheOneYouWantToInclude")}
          color={isSelectedLibrary ? "brand-tertiary" : "tertiary"}
        />
      </div>
    </div>
  </svelte:fragment>
  {#if !isSelectedLibrary}
    <div class="flex flex-col justify-center gap-4">
      <Typography preset="Text sm/Medium" text={$_("latestDocumentLibrary")} color="quaterary" />

      <SelectedCard
        selectionData={collections}
        bind:selectedData={selectedCollection}
        numGrid={2}
        isRadio
      />
    </div>
  {:else if allData.length > 0 || searchTerm !== ""}
    {#if !gettingDatas}
      <DataTableListCollection
        bind:datas={allData}
        bind:datasChecked
        bind:configs
        bind:gettingDatas
        bind:next
        bind:previous
        bind:handleCheckAll
        bind:handleCheck
      />
    {:else}
      <StateCard illustration={IllustrationProgress} description={$_("fetchingDataPleaseWait")} />
    {/if}
  {:else if gettingDatas}
    <StateCard illustration={IllustrationProgress} description={$_("fetchingDataPleaseWait")} />
  {:else}
    <StateCard illustration={IllustrationEmpty} description={$_("youDontHaveAnyDataYet")} />
  {/if}
  <svelte:fragment slot="footer">
    {#if !isSelectedLibrary}
      <Button
        on:click={() => {
          isSelectedLibrary = true;
        }}
        class="w-full"
        disabled={!selectedCollection}>{$_("next")}</Button
      >
    {:else}
      <div class="flex gap-4 w-full">
        <button
          onclick={() => {
            isSelectedLibrary = false;
          }}
          class="px-5">{$_("goBack")}</button
        >
        <Button
          on:click={() => {
            modalOpen = false;
          }}
          class="grow"
          disabled={areAllUnchecked()}>{$_("finish")}</Button
        >
      </div>
    {/if}
  </svelte:fragment>
</Modal>
