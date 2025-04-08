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
  import { Button, Modal, Progressbar } from "flowbite-svelte";
  import { onMount } from "svelte";
  import { _ } from "svelte-i18n";
  import Divider from "$lib/components/Divider.svelte";
  import {
    call,
    // SelectedCard,
    UploadFile,
    triggerToast,
    triggerLoading,
    Typography,
    InputRw,
    TextareaRw,
    Icon,
    // StateCard,
  } from "$lib";
  // import DataTableListCollection from "$lib/components/DataTableListCollection.svelte";
  import { goto } from "$app/navigation";
  // import IllustrationEmpty from "$lib/icon/illustration_empty.svelte";
  // import IllustrationProgress from "$lib/icon/illustration_progress.svelte";
  import { ExclamationCircleSolid } from "flowbite-svelte-icons";

  const isFrontend = process.env.NODE_ENV === "development";

  let { botId = null, editable = true } = $props();

  let collections = $state([]);
  let initCollections = $state([]);

  let botName = $state("");
  let description = $state("");
  let fileImage = $state(null);
  let fileImageName = $state(null);
  let fileImageResult = $state(null);
  let selectedplatforms = $state(null);
  // let selectedCollection = $state("");
  // let oldSelectedCollection = $state("");
  let additionalPrompt = $state("");
  let needAuth = $state(false);
  // let modalOpen = $state(false);
  // let isSelectedLibrary = $state(false);
  // let allData = $state([]);
  // let searchTerm = $state("");
  // let search = $state(false);
  let checkAll = $state(false);
  let datasChecked = $state([]);
  let tmpDatasChecked = $state({});
  // let configs = $state({});
  // let gettingDatas = $state(false);
  let pagination = $state({ page: 1, pageSize: 10, pageCount: 0, total: 0 });
  // let paginationCollection = $state({ page: 1, pageSize: 20, pageCount: 0, total: 0 });
  let platforms = $state([]);
  let screenWidth = $state(0);
  let ready = $state(false);
  let modalDelete = $state(false);
  let llmExists = $state(true);
  let llmStatus = $state("");
  let deleteTarget = $state("");
  let llmId = $state("");

  // async function getCollectionSub(id) {
  //   let token = localStorage.getItem("token");
  //   if (token !== null) {
  //     if (!search) {
  //       gettingDatas = true;
  //     }
  //     let response = await $call.call("client-dashboard", "get-library-collections", {
  //       URLQuery: {},
  //       URLParams: {},
  //       Params: {
  //         page: pagination.page,
  //         page_size: pagination.pageSize,
  //         parent_id: id,
  //         search: searchTerm,
  //       },
  //       Headers: {
  //         Authorization: "Bearer " + token,
  //       },
  //     });

  //     if (response.status !== 200) {
  //       console.log(response);
  //       gettingDatas = false;
  //       return;
  //     }

  //     console.log("collections", response);

  //     if (response.result.data.length > 0) {
  //       pagination = response.result.meta.pagination;
  //       configs = {
  //         search: true,
  //         searchHandler: handleSearch,
  //         pagination: response.result.meta.pagination,
  //         space: "p-4",
  //         actionInLeft: true,
  //         actions: [
  //           {
  //             id: "deselect_all",
  //             name: "Deselect All",
  //             onclick: () => {
  //               handleCheckAllCollection(false);
  //             },
  //           },
  //           {
  //             id: "select_all",
  //             name: "Select All",
  //             type: "brand-inverse",
  //             onclick: () => {
  //               handleCheckAllCollection(true);
  //             },
  //           },
  //         ],
  //       };
  //       let promises = response.result.data.map(async (document) => {
  //         const sub_collections = await getSubCollections(document.id);
  //         return {
  //           id: document.id,
  //           collection_name: document.attributes.name,
  //           sub_collections,
  //         };
  //       });
  //       allData = await Promise.all(promises);
  //       // console.log("all data with subcollections", allData);

  //       if (checkAll) {
  //         datasChecked = allData.map((data) => {
  //           return {
  //             id: data.id,
  //             checked: true,
  //             sub_collection: data.sub_collections.map((sub) => ({
  //               id: sub.id,
  //               checked: true,
  //             })),
  //           };
  //         });
  //         // console.log("datasChecked all", datasChecked);
  //       } else {
  //         if (tmpDatasChecked[pagination.page] === undefined) {
  //           datasChecked = allData.map((data) => {
  //             return {
  //               id: data.id,
  //               checked: false,
  //               sub_collection: data.sub_collections.map((sub) => ({
  //                 id: sub.id,
  //                 checked: false,
  //               })),
  //             };
  //           });
  //         } else {
  //           datasChecked = tmpDatasChecked[pagination.page];
  //         }
  //         // console.log("datasChecked", datasChecked);
  //       }

  //       // if (checkAll) {
  //       //   datasChecked = datasChecked.map((data) => {
  //       //     return { id: data.id, checked: true };
  //       //   });
  //       // }
  //     } else {
  //       pagination = response.result.meta.pagination;
  //       allData = [];
  //       datasChecked = [];
  //       configs = {
  //         search: true,
  //         searchHandler: handleSearch,
  //         pagination: response.result.meta.pagination,
  //         space: "p-4",
  //         actionInLeft: true,
  //         actions: [
  //           {
  //             id: "deselect_all",
  //             name: "Deselect All",
  //             onclick: () => {
  //               handleCheckAllCollection(false);
  //             },
  //           },
  //           {
  //             id: "select_all",
  //             name: "Select All",
  //             type: "brand-inverse",
  //             onclick: () => {
  //               handleCheckAllCollection(true);
  //             },
  //           },
  //         ],
  //         // actions: [
  //         //   {
  //         //     id: "delete",
  //         //     name: "Delete",
  //         //     icon: "trash",
  //         //     type: "error",
  //         //     disabled: datasChecked.find((data) => data.checked) ? false : true,
  //         //     onclick: () => {
  //         //       if (datasChecked.find((data) => data.checked)) {
  //         //         popupModal = true;
  //         //       }
  //         //     },
  //         //   },
  //         //   {
  //         //     id: "add_new",
  //         //     name: "Add New",
  //         //     icon: "plus",
  //         //     type: "brand-inverse",
  //         //     onclick: () => {
  //         //       selectTypeCollection("Documents");
  //         //       addDocumentsModal = true;
  //         //       addOtherCategory = false;
  //         //     },
  //         //   },
  //         // ],
  //       };
  //     }

  //     gettingDatas = false;
  //     search = false;
  //   }
  // }

  // async function getSubCollections(id) {
  //   let token = localStorage.getItem("token");

  //   if (token !== null) {
  //     gettingDatas = true;
  //     const aggregatedData = [];

  //     paginationCollection.page = 1;

  //     do {
  //       let response = await $call.call("client-dashboard", "get-library-collections", {
  //         URLQuery: {},
  //         URLParams: {},
  //         Params: {
  //           page: paginationCollection.page,
  //           page_size: paginationCollection.pageSize,
  //           parent_id: parseInt(id),
  //         },
  //         Headers: {
  //           Authorization: "Bearer " + token,
  //         },
  //       });

  //       if (response.status !== 200) {
  //         console.log(response);
  //         gettingDatas = false;
  //         if (response.status === 400) {
  //           console.error("Invalid parameters:", response);
  //           triggerToast($_("invalidParameters"), true, "error");
  //         } else if (response.status === 401) {
  //           console.error("Unauthorized access:", response);
  //           triggerToast($_("unauthorized"), true, "error");
  //         } else if (response.status === 500) {
  //           console.error("Server error:", response);
  //           triggerToast($_("serverError"), true, "error");
  //         } else {
  //           console.error("Unexpected status code:", response);
  //           triggerToast(`${$_("unexpectedError")}: ${response.status}`, true, "error");
  //         }
  //         return;
  //       }

  //       aggregatedData.push(...response.result.data);

  //       paginationCollection = response.result.meta.pagination;

  //       paginationCollection.page++;
  //     } while (paginationCollection.page <= paginationCollection.pageCount);

  //     const promise = aggregatedData.map(async (collection) => ({
  //       id: collection.id,
  //       collection_name: collection.attributes.name,
  //     }));

  //     const subCollectionList = await Promise.all(promise);
  //     console.log("subCollectionList", subCollectionList);
  //     gettingDatas = false;
  //     return subCollectionList;
  //   } else {
  //     triggerToast($_("unauthorized"), true, "error");
  //     gettingDatas = false;
  //     return;
  //   }
  // }

  // const handleSearch = async (value) => {
  //   search = true;
  //   pagination.page = 1;
  //   searchTerm = value;

  //   await getCollectionSub(selectedCollection);
  // };

  const handleInputFile = (type = "image") => {
    const inputElement = document.createElement("input");
    inputElement.type = "file";
    if (type === "image") {
      inputElement.accept = "image/jpeg, image/png";
    } else {
      inputElement.accept = "image/jpeg, image/png";
    }

    inputElement.addEventListener("change", async (event) => {
      const file = event.target.files[0];

      if (file && ["image/jpeg", "image/png"].includes(file.type)) {
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
          } else {
            triggerToast("error");
          }
          return;
        }

        if (response.status === 200) {
          let i = 0;
          let promise = response.result.models.map((model) => {
            return {
              id: i++,
              name: model.name,
              model: model.model,
              description: model.description,
              create_at: model.create_at,
            };
          });

          platforms = await Promise.all(promise);
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
        library_id: 404, // TODO : Change this to the actual library ID
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
      console.log("no token");
      return;
    }
  }

  async function getBot() {
    let token = localStorage.getItem("token");
    if (token) {
      try {
        let response = await $call.call("client-dashboard", "get-bot", {
          URLQuery: {},
          URLParams: {
            id: botId,
          },
          Params: {},
          Headers: {
            Authorization: "Bearer " + token,
          },
        });

        if (response.status !== 200) {
          console.log(response);
          triggerLoading(false);
          triggerToast("error");
          goto(await $call.getPage("client-dashboard", "bots-list"));
          return;
        }

        const bot = response.result;
        const attributes = bot.attributes;

        // await getDocuments();

        // users = {
        //   botName: attributes.bot_name,
        //   logoName: attributes.logo_name,
        //   logoOnly: attributes.logo_only,
        //   deployedPlatform: attributes?.deployed_platform || "Web",
        //   description: attributes.description,
        //   additionalPrompts: attributes.additional_prompt,
        //   welcomeMessage: attributes.greeting_text,
        //   username: attributes?.username || "User",
        //   primaryBrandColor:
        //     attributes.color_theme !== null
        //       ? themeColors.find((color) => color.value === attributes.color_theme).id
        //       : 1,
        //   overlayColor:
        //     attributes.overlay_color !== null
        //       ? themeColors.find((color) => color.value === attributes.overlay_color).id
        //       : 1,
        //   platform: platforms.find((platform) => platform.model === attributes.ai_engine).id,
        //   layout:
        //     attributes.layout !== null
        //       ? layouts.find((layout) => layout.name === attributes.layout).id
        //       : 1,
        //   background_align:
        //     attributes.background_align !== null ? attributes.background_align : "center",
        //   embedCode: `<iframe style=”width: 100%; height: 100%; border: 0px; border-radius: 10px; display: block;" src="${urlBot}"></iframe>`,
        //   needAuth: attributes.need_auth,
        //   documentOptions: attributes.manage_document,
        // };

        botName = attributes.bot_name;
        description = attributes.description;

        // if (attributes.background.data !== null) {
        //   backgroundFile = getFile(
        //     bot.background,
        //     attributes.background.data.attributes.ext,
        //     attributes.background.data.attributes.name,
        //   );
        // }

        // if (attributes.bot_logo.data !== null) {
        //   logoFile = getFile(
        //     bot.bot_logo,
        //     attributes.bot_logo.data.attributes.ext,
        //     attributes.bot_logo.data.attributes.name,
        //   );
        // }

        // if (attributes.user_profile.data !== null) {
        //   userFile = getFile(
        //     bot.user_profile,
        //     attributes.user_profile.data.attributes.ext,
        //     attributes.user_profile.data.attributes.name,
        //   );
        // }

        if (attributes.bot_profile.data !== null) {
          let botFile = getFile(
            bot.bot_profile,
            attributes.bot_profile.data.attributes.ext,
            attributes.bot_profile.data.attributes.name,
          );
          fileImageName = botFile.name;
          fileImageResult = botFile.imageResult;
        }

        ready = true;
        triggerLoading(false);
      } catch (error) {
        console.log(error);
      }
    }
  }

  async function deleteBot(id) {
    let token = localStorage.getItem("token");

    if (token) {
      triggerLoading(true);
      try {
        const response = await $call.call("client-dashboard", "delete-bot", {
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
          triggerLoading(false);
          triggerToast("error");
          console.log(response);
          return;
        }

        // popupModal = false;
        triggerLoading(false);
        triggerToast("deleted");
        goto(await $call.getPage("client-dashboard", "bots-list"));
      } catch (error) {
        triggerLoading(false);
        console.log(error);
      }
    }
  }

  function getFile(base64, ext, name) {
    let fileType;
    if (ext === ".png") {
      fileType = "image/png";
    } else if (ext === ".jpg") {
      fileType = "image/jpg";
    } else if (ext === ".jpeg") {
      fileType = "image/jpeg";
    } else if (ext === ".svg") {
      fileType = "image/svg+xml";
    } else {
      fileType = "image/png";
    }

    let blob = base64toBlob(base64, fileType);
    let fileImageResult = URL.createObjectURL(blob);
    return {
      imageResult: fileImageResult,
      name: name,
      base64: base64,
    };
  }

  function base64toBlob(base64Data, fileType) {
    const byteCharacters = atob(base64Data);
    const byteArrays = [];

    for (let offset = 0; offset < byteCharacters.length; offset += 512) {
      const slice = byteCharacters.slice(offset, offset + 512);

      const byteNumbers = new Array(slice.length);
      for (let i = 0; i < slice.length; i++) {
        byteNumbers[i] = slice.charCodeAt(i);
      }

      const byteArray = new Uint8Array(byteNumbers);
      byteArrays.push(byteArray);
    }

    const blob = new Blob(byteArrays, { type: fileType });
    return blob;
  }

  // async function handleCheckAll(allChecked) {
  //   datasChecked = datasChecked.map((data) => ({
  //     id: data.id,
  //     checked: allChecked ? false : true,
  //     sub_collection: data.sub_collection.map((sub) => ({
  //       id: sub.id,
  //       checked: allChecked ? false : true,
  //     })),
  //   }));

  //   // console.log("Updated datasChecked (handleCheckAll):", datasChecked);
  // }

  // async function handleCheckAllCollection(allChecked) {
  //   checkAll = allChecked;

  //   datasChecked = datasChecked.map((data) => ({
  //     id: data.id,
  //     checked: allChecked,
  //     sub_collection: data.sub_collection.map((sub) => ({
  //       id: sub.id,
  //       checked: allChecked,
  //     })),
  //   }));

  //   if (!allChecked) {
  //     // Jika allChecked adalah false, update tmpDatasChecked
  //     Object.keys(tmpDatasChecked).forEach((page) => {
  //       tmpDatasChecked[page] = tmpDatasChecked[page].map((data) => ({
  //         id: data.id,
  //         checked: false,
  //         sub_collection: data.sub_collection.map((sub) => ({
  //           id: sub.id,
  //           checked: false,
  //         })),
  //       }));
  //     });
  //   }

  //   // console.log("Updated datasChecked (handleCheckAllCollection):", datasChecked);
  // }

  // async function handleCheck(id, isSubCollection = false, parentId = null) {
  //   if (!Array.isArray(datasChecked)) {
  //     datasChecked = [];
  //   }

  //   datasChecked = datasChecked.map((data) => {
  //     // Jika parentId tidak ada, kita memproses level utama
  //     if (!isSubCollection && data.id === id) {
  //       if (data.checked) {
  //         checkAll = false;
  //       }
  //       return {
  //         ...data,
  //         checked: !data.checked,
  //       };
  //     }

  //     // console.log("checkAll", checkAll);

  //     // Jika parentId cocok, kita memproses sub_collection
  //     if (isSubCollection && data.id === parentId) {
  //       const updatedSubCollection = data.sub_collection.map((sub) => {
  //         if (sub.id === id) {
  //           return { ...sub, checked: !sub.checked };
  //         }
  //         return sub;
  //       });

  //       return {
  //         ...data,
  //         sub_collection: updatedSubCollection,
  //       };
  //     }

  //     return data;
  //   });

  //   let pageNow = pagination.page;
  //   if (tmpDatasChecked[pageNow] === undefined) {
  //     tmpDatasChecked = {
  //       ...tmpDatasChecked,
  //       [pageNow]: datasChecked,
  //     };
  //   } else {
  //     tmpDatasChecked[pageNow] = datasChecked;
  //   }
  //   // console.log("datasChecked setelah pembaruan:", datasChecked);
  // }

  // function areAllUnchecked() {
  //   // Gabungkan data dari halaman sekarang dengan data yang tersimpan di tmpDatasChecked
  //   const allData = [...datasChecked, ...Object.values(tmpDatasChecked).flat()];

  //   // Cek apakah semua data tidak ada yang checked
  //   return allData.every((item) => item.checked === false);
  // }

  // async function next() {
  //   if (pagination.page >= pagination.pageCount) return;

  //   let pageNow = pagination.page;
  //   if (tmpDatasChecked[pageNow] === undefined) {
  //     tmpDatasChecked = {
  //       ...tmpDatasChecked,
  //       [pageNow]: datasChecked,
  //     };
  //   } else {
  //     tmpDatasChecked[pageNow] = datasChecked;
  //   }

  //   pagination.page = pagination.page + 1;
  //   // console.log("cek tmpDatasChecked", tmpDatasChecked);
  //   // getDocuments();
  //   getCollectionSub(selectedCollection);
  // }

  // async function previous() {
  //   if (pagination.page <= 1) return;

  //   let pageNow = pagination.page;
  //   if (tmpDatasChecked[pageNow] === undefined) {
  //     tmpDatasChecked = {
  //       ...tmpDatasChecked,
  //       [pageNow]: datasChecked,
  //     };
  //   } else {
  //     tmpDatasChecked[pageNow] = datasChecked;
  //   }

  //   pagination.page = pagination.page - 1;
  //   // getDocuments();
  //   getCollectionSub(selectedCollection);
  // }

  onMount(() => {
    getCollections();
    getModels();
    if (botId !== null) {
      getBot();
    }
    if (editable === null || editable === "false") {
      editable = false;
    }
  });

  // $effect(() => {
  //   if (selectedCollection) {
  //     if (typeof selectedCollection === "number" && selectedCollection !== oldSelectedCollection) {
  //       oldSelectedCollection = selectedCollection;
  //       getCollectionSub(selectedCollection);
  //     }
  //   }
  // });
</script>

<svelte:window bind:innerWidth={screenWidth} />

<div class="m-4 pb-12">
  <div class="flex justify-between">
    <Typography preset="Display sm/Semibold" text={$_("deployMyModel")} />
    <div class="flex text-sm gap-3">
      {#if botId !== null}
        {#if editable}
          <Button
            onclick={() => {
              editable = false;
              getBot();
            }}
            color="light"
            class="px-3 focus-within:!ring-1"
          >
            {$_("cancel")}
          </Button>
          <Button on:click={() => updateBot()} class="px-3 focus-within:!ring-1">
            {$_("save")}
          </Button>
        {:else}
          <Button
            onclick={async function () {
              goto(await $call.getPage("client-dashboard", "my_models"));
              return;
            }}
            color="light"
            class="px-3 focus-within:!ring-1"
          >
            {$_("back")}
          </Button>
          <Button
            onclick={() => (modalDelete = true)}
            outline
            color="red"
            class="px-3 focus-within:!ring-1"
          >
            {$_("delete")}
          </Button>
          <Button
            onclick={() => {
              editable = true;
            }}
            class="px-3 focus-within:!ring-1"
          >
            {$_("edit")}
          </Button>
        {/if}
      {:else}
        <Button
          on:click={async () => {
            goto(await $call.getPage("client-dashboard", "my_models"));
          }}
          color="light"
        >
          {$_("cancel")}
        </Button>
        <Button
          disabled={botName === "" || selectedplatforms === ""}
          on:click={async () => {
            await createBot();
            goto(await $call.getPage("client-dashboard", "my_models"));
          }}
          class="bg-hc-brand-600 hover:bg-hc-brand-700 dark:bg-hc-brand-600 dark:hover:bg-hc-brand-700"
        >
          {$_("deploy")}
        </Button>
      {/if}
    </div>
  </div>

  <!-- Content -->
  <div class="mt-6 space-y-6">
    <div class="md:flex">
      <Typography
        preset="Text md/Semibold"
        text={$_("generalInformation")}
        propClass="w-full md:w-1/4"
      />
      <div class="w-full md:w-3/4 grid grid-cols-1 md:grid-cols-4 gap-3">
        <div class="flex flex-col gap-1.5 col-span-4">
          <Typography preset="Text sm/Medium" text={$_("modelName")} color="secondary" />
          <InputRw bind:value={botName} onEdit={editable} />
        </div>
        <div class="flex flex-col gap-1.5 col-span-4">
          <Typography preset="Text sm/Medium" text={$_("modelDescription")} color="secondary" />
          <TextareaRw bind:value={description} onEdit={editable} />
        </div>
      </div>
    </div>
    <Divider />
    <div class="md:flex">
      <Typography preset="Text md/Semibold" text={$_("logo")} propClass="w-full md:w-1/4" />
      <div class="w-full md:w-3/4 grid grid-cols-1 md:grid-cols-4 gap-3">
        <div class="col-span-4">
          <UploadFile
            title={$_("logo")}
            handleFileUploadWrapper={() => handleInputFile()}
            isPreview
            srcAvatar={fileImageName ? fileImageResult : ""}
            disabled={!editable}
          />
        </div>
      </div>
    </div>
    <Divider />
    <div class="md:flex">
      <Typography preset="Text md/Semibold" text={$_("sourceFile")} propClass="w-full md:w-1/4" />
      <div class="w-full md:w-3/4 grid grid-cols-1 md:grid-cols-4 gap-3">
        <div class="col-span-4">
          <UploadFile
            title={$_("modelSourceFile")}
            handleFileUploadWrapper={() => handleInputFile()}
            disabled={!editable}
          />
          {#if llmExists}
            <div
              class="flex justify-between mt-4 p-6 w-full gap-4 border border-primary-light dark:border-primary-dark rounded-lg bg-primary-light dark:bg-secondary-dark dark:hover:border-gray-600"
            >
              {#if llmStatus === "loading"}
                <div class="flex flex-col justify-center gap-2 w-full">
                  <Typography preset="Text sm/Medium" text={"Uploading File"} color="primary" />
                  <Progressbar progress={60} size="h-4" labelInside />
                </div>
              {:else}
                <div class="flex items-center gap-2">
                  <Icon name="onedrive" size="w-8 h-8" />
                  <Typography
                    preset="Text sm/Medium"
                    text={"AI Document Archive"}
                    color="primary"
                  />
                </div>
                <div class="flex items-center gap-2">
                  <Typography preset="Text sm/Medium" text={"2.13 GB"} color="secondary" />
                  <div
                    class={editable ? "cursor-pointer" : ""}
                    onclick={() => {
                      if (editable) {
                        deleteTarget = "llm";
                        llmId = 1;
                        modalDelete = true;
                      }
                    }}
                  >
                    <Icon name="trash03" color="error-primary" />
                  </div>
                </div>
              {/if}
            </div>
          {/if}
        </div>
      </div>
    </div>
  </div>
</div>

<!-- <Modal bind:open={modalOpen} outsideclose>
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
</Modal> -->

<Modal bind:open={modalDelete} size="xs" autoclose>
  <div class="text-center">
    <ExclamationCircleSolid class="mx-auto mb-4 text-red-600 w-12 h-12" />
    <h3 class="mb-5 text-lg font-normal text-gray-500 dark:text-gray-400">
      {deleteTarget === "llm"
        ? $_("areYouSureYouWantToDeleteLlm")
        : $_("areYouSureYouWantToDeleteBot")}
    </h3>
    <Button
      onclick={() => {
        if (deleteTarget === "llm") {
          // deleteAiDocument(llmId)
          console.log("hapus llm");
        } else {
          deleteBot(botId);
        }
      }}
      color="red"
      class="me-2">{$_("yesImSure")}</Button
    >
    <Button color="alternative">{$_("noCancel")}</Button>
  </div>
</Modal>
