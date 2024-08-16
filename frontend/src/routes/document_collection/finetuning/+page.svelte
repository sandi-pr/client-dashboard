<!--
CR3 Client Dashboard
Package: gitlab.hexacode.org/cr3/client-dashboard
Maintainer:
  - Azzis Arswendo <azzis@hexacode.org>
  - Dean Ramhan <dean@hexacode.org>
  - Taqie Fadlillah <taqief@hexacode.org>
  - Sandi Permana <sandy@hexacode.org>
Copyright (C) 2024 Hexacode Teknologi Indonesia
All Rights Reserved
-->
<script>
  import { P, Heading, Input, Button, CloseButton, Textarea } from "flowbite-svelte";
  import { ToastContainer, triggerToast } from "$lib";
  import upload from "$lib/icon/upload-cloud-02.svg";
  import { _ } from "svelte-i18n";
  import Divider from "$lib/components/Divider.svelte";
  import { onMount } from "svelte";
  import { CR3Call, Typography } from "hexacomponents-svelte";
  import { format } from "date-fns";
  import { goto } from "$app/navigation";
  import PDFSolidIcon from "$lib/icon/PDFSolid.svg";
  import DOCSolidIcon from "$lib/icon/DOCSolid.svg";
  import DOCXSolidIcon from "$lib/icon/DOCXSolid.svg";

  const params = new URLSearchParams(window.location.search);
  let collectionId = params.get("collectionId");

  import { call } from "$lib";

  let fileTextName = $state(null);

  let images = $state([]);

  const handleFileUploadWrapper = () => {
    const inputElement = document.createElement("input");
    inputElement.type = "file";
    inputElement.accept = ".jpg, .jpeg, .png";
    inputElement.multiple = true;

    inputElement.addEventListener("change", async (event) => {
      const files = event.target.files;

      if (files.length > 0) {
        for (let i = 0; i < files.length; i++) {
          const file = files[i];
          const imageUrl = URL.createObjectURL(file);
          const base64 = await readFileAsBase64(file);

          images = [
            ...images,
            {
              id: Date.now() + i,
              file: file,
              imageUrl: imageUrl,
              base64: base64,
              name: file.name,
              sourceText: "",
            },
          ];
        }
      } else {
        alert("Invalid file. Please select valid files.");
      }
    });

    inputElement.click();
  };

  const removeImage = (id) => {
    images = images.filter((image) => image.id !== id);
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

  async function addDataFineTuning() {
    let token = sessionStorage.getItem("token");

    if (token !== null) {
      if (images.length > 0) {
        try {
          let response = await $call.call("client-dashboard", "add-data-handwriting", {
            URLQuery: {},
            URLParams: {},
            Params: {
              images: JSON.stringify(images),
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

          if (response.status === 200) {
            triggerToast("created");
            // setTimeout(async () => {
            window.location.reload();
            // }, 3000);
          }
        } catch (error) {
          triggerToast("error");
          console.log(error);
        }
      }
      return;
    } else {
      triggerToast("error");
    }
  }
</script>

<div class="flex dark:bg-hc-graydark-950">
  <div class="w-full pb-10">
    <div class="flex flex-col m-4">
      <!-- Header -->
      <div class="flex items-center pb-8">
        <div class="flex flex-col gap-2 grow">
          <!-- Header -->
          <div class="gap-2">
            <Heading
              tag="h3"
              class="text-3xl leading-[38px] text-hc-graylight-900 dark:text-hc-graydark-50"
            >
              {$_("addNewFinetuningData")}
            </Heading>
          </div>
        </div>
        <div class="flex gap-2">
          <Button
            onclick={() => {
              window.history.back();
            }}
            color="light"
            class="px-3 focus-within:!ring-1">{$_("cancel")}</Button
          >
          <Button onclick={addDataFineTuning} class="px-3 focus-within:!ring-1"
            >{$_("startTraining")}</Button
          >
        </div>
      </div>

      <Divider />

      <div class="flex flex-col gap-4 pt-8">
        <!-- general information -->
        <div class="flex gap-6">
          <div class="w-60">
            <P class="font-semibold text-secondary-light dark:text-secondary-dark">
              {$_("trainingData")}
            </P>
          </div>
          <div class="flex flex-col grow gap-6">
            <div class="flex flex-col col-span-5 gap-2">
              <P
                class="text-sm font-medium leading-tight text-secondary-light dark:text-secondary-dark"
                >{$_("imageUpload")}*</P
              >
              <div class="flex flex-col col-span-5 gap-2">
                <Button
                  onclick={() => {
                    handleFileUploadWrapper();
                  }}
                  color="light"
                  class={`w-full rounded-xl border border-primary-light dark:border-primary-dark bg-primary-light dark:bg-primary-dark gap-1 py-8 flex flex-col justify-start ${fileTextName ? "items-start cursor-default" : "items-center"}`}
                >
                  <div
                    class="rounded-xl shadow border border-primary-light dark:border-primary-dark justify-center items-center inline-flex mb-4 pointer-events-none"
                  >
                    <img src={upload} alt="mail.svg" class="p-4 dark:invert pointer-events-none" />
                  </div>
                  <div class="flex space-x-1">
                    <P class="text-hc-brand-700 pointer-events-none">{$_("clickToUpload")}</P>
                    <P class="pointer-events-none">
                      <!-- {$_("orDragAndDrop")} -->
                    </P>
                  </div>
                  <P class="pointer-events-none">PNG, JPG, JPEG (Max. 15MB)</P>
                </Button>
              </div>
              <!-- Display uploaded file names -->
              {#if images.length > 0}
                <div class="flex flex-col gap-3">
                  {#each images as image}
                    <div
                      class="flex justify-between items-start p-4 rounded-xl border border-secondary-light dark:border-secondary-dark"
                    >
                      <div class="flex flex-col gap-3">
                        <P class="text-secondary-light dark:text-secondary-dark">{image.name}</P>
                        <img src={image.imageUrl} alt="" />
                        <div class="flex flex-col gap-1 justify-center">
                          <P class="text-secondary-light dark:text-secondary-dark"
                            >{$_("sourceText")}</P
                          >
                          <Textarea
                            bind:value={image.sourceText}
                            class="bg-primary-light dark:bg-secondary-dark"
                          />
                        </div>
                      </div>
                      <CloseButton
                        on:click={() => removeImage(image.id)}
                        class="h-min text-secondary-light dark:text-secondary-dark"
                      />
                    </div>
                  {/each}
                </div>
              {/if}
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
  <ToastContainer />
</div>
