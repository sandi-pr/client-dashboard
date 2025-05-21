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
    UploadFile,
    triggerToast,
    triggerLoading,
    Typography,
    InputRw,
    TextareaRw,
    Icon,
  } from "$lib";
  import LaodingLottie from "$lib/components/LaodingLottie.svelte";
  import { goto } from "$app/navigation";
  import { ExclamationCircleSolid } from "flowbite-svelte-icons";

  const isFrontend = process.env.NODE_ENV === "development";

  let { botId = null, modelId = null, editable = true } = $props();

  let collections = $state([]);
  let initCollections = $state([]);

  let modelName = $state("");
  let description = $state("");
  let fileImage = $state(null);
  let fileImageName = $state(null);
  let fileImageResult = $state(null);
  let platforms = $state([]);
  let screenWidth = $state(0);
  let ready = $state(false);
  let modalDelete = $state(false);
  let deleteTarget = $state("");
  let llmId = $state("");

  let selectedFile = $state(null);
  let fileSizeMB = $state("");
  let uploadProgress = $state(0);
  let isLoading = $state(false);
  let modalUpload = $state(false);

  const handleFileInput = (type = "image") => {
    const inputElement = document.createElement("input");
    inputElement.type = "file";

    if (type === "image") {
      inputElement.accept = "image/jpeg, image/png";
    } else if (type === "gguf") {
      inputElement.accept = ".gguf";
    } else {
      inputElement.accept = "image/jpeg, image/png, .gguf";
    }

    inputElement.addEventListener("change", async (event) => {
      const file = event.target.files[0];

      if (file) {
        if (["image/jpeg", "image/png"].includes(file.type)) {
          // Handle image files
          fileImageResult = await new Promise((resolve, reject) => {
            const reader = new FileReader();

            reader.onload = () => resolve(reader.result);
            reader.onerror = reject;

            reader.readAsDataURL(file);
          });

          fileImage = await readFileAsBase64(file);
          fileImageName = file.name;
        } else if (file.name.endsWith(".gguf")) {
          // Handle .gguf files
          selectedFile = file;
          modelName = modelName ? modelName : file.name;
          fileSizeMB = (file.size / (1024 * 1024)).toFixed(2);
          uploadProgress = 0;
        } else {
          if (type === "image") {
            alert("Invalid file. Please select a valid .jpg or .png file.");
          } else if (type === "gguf") {
            alert("Invalid file. Please select a valid .gguf file.");
          } else {
            alert("Invalid file. Please select a valid .jpg, .png, or .gguf file.");
          }
        }
      } else {
        alert("No file selected.");
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
        console.log("modelId", modelId);
        let response = await $call.call("hexa-ai-rag", "get-model-by-id", {
          URLQuery: {},
          URLParams: {
            id: modelId,
          },
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
          // let promise = response.result.models.map((model) => {
          //   return {
          //     id: i++,
          //     name: model.name,
          //     model: model.model,
          //     description: model.description,
          //     create_at: model.create_at,
          //   };
          // });

          // platforms = await Promise.all(promise);

          modelName = response.result.name;
          description = response.result.description;
        } else {
          triggerToast("error");
          console.log(response);
        }
      } catch (error) {
        console.log(error);
      }
    }
  }

  async function createModel() {
    if (!selectedFile) {
      triggerToast("please select file");
      return;
    } else if (fileImage === "" || fileImage === null || fileImage === undefined) {
      triggerToast("please select logo");
      return;
    }

    let token = localStorage.getItem("token");

    const logoData = JSON.stringify({ name: fileImageName, base64Data: fileImage });
    const extension = modelName.split(".").pop() || "";
    const chunkSize = 1024 * 1024; // 1MB chunks
    const totalChunks = Math.ceil(selectedFile.size / chunkSize);
    let offset = 0;
    let chunkIndex = 0;
    let limitError = 5;

    uploadProgress = 0;

    if (token !== null) {
      // triggerLoading(true);
      isLoading = true;

      while (offset < selectedFile.size) {
        chunkIndex++;
        const chunk = selectedFile.slice(offset, offset + 1024 * 1024);

        const formData = new FormData();
        formData.append("chunkIndex", chunkIndex.toString());
        formData.append("totalChunks", totalChunks.toString());
        formData.append("modelName", modelName);
        formData.append("extension", extension);
        formData.append("chunk", chunk, selectedFile.name);
        formData.append("description", description);
        formData.append("logo", logoData);

        // const payload = {
        //   chunkIndex: chunkIndex.toString(),
        //   totalChunks: totalChunks.toString(),
        //   modelName: modelName,
        //   extension: extension,
        //   chunk: await chunk.arrayBuffer(), // Convert the chunk to a buffer
        //   description: description,
        //   logo: logoData,
        // };

        // console.log("payload", payload);
        // triggerLoading(false);
        // return;

        try {
          const response = await $call.call("hexa-ai-rag", "upload-my-model", {
            URLQuery: {},
            URLParams: {},
            Params: formData,
            // Params: payload,
            Headers: {
              Authorization: `Bearer ${token}`,
            },
          });

          if (response.status !== 200) {
            console.log(response);
            triggerToast("Upload failed");
            isLoading = false;
            return;
          }

          offset += chunkSize;
          uploadProgress = (offset / selectedFile.size) * 100;
          console.log(`Upload progress: ${uploadProgress}%`);
        } catch (error) {
          console.log("Error during upload:", error);
          limitError--;
          if (limitError === 0) {
            triggerToast("Upload failed after multiple attempts");
            isLoading = false;
            return;
          }
        }
      }

      // triggerLoading(false);
      isLoading = false;
      triggerToast("created");
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
          goto(await $call.getPage("client-dashboard", "llms-my-models"));
          return;
        }

        const bot = response.result;
        const attributes = bot.attributes;

        // await getDocuments();

        // users = {
        //   modelName: attributes.bot_name,
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

        // modelName = attributes.bot_name;
        // description = attributes.description;

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
        goto(await $call.getPage("client-dashboard", "llms-my-models"));
      } catch (error) {
        triggerLoading(false);
        console.log(error);
      }
    }
  }

  async function deleteModel(id) {
    let token = localStorage.getItem("token");

    if (token) {
      triggerLoading(true);
      try {
        const response = await $call.call("client-dashboard", "delete-model", {
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
        goto(await $call.getPage("client-dashboard", "llms-my-models"));
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

  function getFileSize(fileSize) {
    if (!fileSize) {
      return "undefined";
    }
    if (fileSize >= 1024) {
      return fileSize / 1024 + " GB";
    } else {
      return fileSize + " MB";
    }
  }

  onMount(() => {
    getCollections();
    if (modelId !== null) {
      // console.log("modelId", modelId);
      getModels();
    }
    if (botId !== null) {
      getBot();
    }
    if (editable === null || editable === "false") {
      editable = false;
    }
  });
</script>

<svelte:window bind:innerWidth={screenWidth} />

<div class="m-4 pb-12">
  <div class="flex flex-col md:flex-row gap-4 justify-between">
    <Typography preset="Display sm/Semibold" text={$_("deployMyModel")} />
    <div class="flex text-sm gap-3">
      {#if botId !== null || modelId !== null}
        {#if editable}
          <Button
            onclick={() => {
              editable = false;
              // getBot();
              getModels();
            }}
            color="light"
            class="w-full px-3 focus-within:!ring-1"
          >
            {$_("cancel")}
          </Button>
          <Button on:click={() => updateBot()} class="w-full px-3 focus-within:!ring-1">
            {$_("save")}
          </Button>
        {:else}
          <Button
            color="light"
            class="w-full px-3 focus-within:!ring-1"
            on:click={async () => {
              if (isFrontend) {
                goto("/llms/my_models");
              } else {
                goto(await $call.getPage("client-dashboard", "llms-my-models"));
              }
            }}
          >
            {$_("back")}
          </Button>
          <Button
            onclick={() => (modalDelete = true)}
            outline
            color="red"
            class="w-full px-3 focus-within:!ring-1"
          >
            {$_("delete")}
          </Button>
          <Button
            onclick={() => {
              editable = true;
            }}
            class="w-full px-3 focus-within:!ring-1"
          >
            {$_("edit")}
          </Button>
        {/if}
      {:else}
        <Button
          on:click={async () => {
            if (isFrontend) {
              goto("/llms/my_models");
            } else {
              goto(await $call.getPage("client-dashboard", "llms-my-models"));
            }
          }}
          color="light"
          class="w-full"
        >
          {$_("cancel")}
        </Button>
        <Button
          disabled={modelName === "" || !fileImageResult || !selectedFile}
          on:click={async () => {
            // await createModel();
            modalUpload = true;
          }}
          class="w-full bg-hc-brand-600 hover:bg-hc-brand-700 dark:bg-hc-brand-600 dark:hover:bg-hc-brand-700"
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
          <InputRw bind:value={modelName} onEdit={editable} />
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
            handleFileUploadWrapper={() => handleFileInput()}
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
            handleFileUploadWrapper={() => handleFileInput("gguf")}
            disabled={!editable}
          />
          {#if modelName && fileSizeMB}
            <div
              class="flex justify-between mt-4 p-6 w-full gap-4 border border-primary-light dark:border-primary-dark rounded-lg bg-primary-light dark:bg-secondary-dark dark:hover:border-gray-600"
            >
              <div class="flex items-center gap-2">
                <Icon name="onedrive" size="w-8 h-8" />
                <Typography preset="Text sm/Medium" text={modelName} color="primary" />
              </div>
              <div class="flex items-center gap-2">
                <Typography
                  preset="Text sm/Medium"
                  text={getFileSize(fileSizeMB)}
                  color="secondary"
                />
                <div
                  class={editable ? "cursor-pointer" : ""}
                  onclick={() => {
                    if (editable && (botId !== null || modelId !== null)) {
                      deleteTarget = "llm";
                      llmId = 1;
                      modalDelete = true;
                    } else {
                      selectedFile = null;
                      modelName = "";
                      fileSizeMB = "";
                      uploadProgress = 0;
                    }
                  }}
                >
                  <Icon name="trash03" color="error-primary" />
                </div>
              </div>
            </div>
          {/if}
        </div>
      </div>
    </div>
  </div>
</div>

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

<Modal bind:open={modalUpload} size="xs" autoclose>
  <div class="text-center justify-center">
    <ExclamationCircleSolid
      class="mx-auto mb-4 text-warning-primary-light dark:text-warning-primary-dark w-12 h-12"
    />
    <Typography
      preset="Text lg/Semibold"
      text={$_("warningLargeFileUpload")}
      propClass="text-center"
    />
    <Typography
      preset="Text sm/Regular"
      text={$_(
        "PleaseEnsureYouHaveAStableInternetConnectionDoNotCloseThisTabOrNavigateAwayDuringTheUploadProcess",
      )}
      propClass="text-center mt-2 mb-8"
    />
    <div class="flex gap-2 w-full">
      <Button color="alternative" class="flex-grow">{$_("noCancel")}</Button>
      <Button
        onclick={() => {
          createModel();
        }}
        class="flex-grow">{$_("yesImSure")}</Button
      >
    </div>
  </div>
</Modal>

{#if isLoading}
  <div
    class="fixed top-0 left-0 w-full h-full flex flex-col gap-3 z-50 bg-fg-brand-primary-light dark:bg-fg-brand-primary-dark"
  >
    <div class="flex flex-col gap-2 items-center justify-center h-full">
      <LaodingLottie width={200} />
      <div class="flex flex-col gap-2 w-full items-center">
        <Typography
          preset="Text xl/Medium"
          color="primary-onbrand"
          text={`${$_("uploadingFilePleaseWait")}`}
          propClass="text-center"
        />
        <div class="w-1/2">
          <Progressbar progress={uploadProgress} size="h-4" labelInside />
        </div>
        <Typography
          preset="Text xl/Medium"
          color="primary-onbrand"
          text={`${$_("pleaseDontCloseTheSystemWhileUploading")}`}
          propClass="text-center"
        />
      </div>
    </div>
  </div>
{/if}
