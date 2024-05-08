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
  import { P, Heading, Input, Button, Select, Textarea, Toggle } from "flowbite-svelte";
  import upload from "$lib/icon/upload-cloud-02.svg";
  import { _ } from "svelte-i18n";
  import CopyIcon from "$lib/icon/copy-03.svelte";
  import TabBar from "$lib/ui/tab_bar/TabBar.svelte";
  import TabDetail from "$lib/components/TabDetail.svelte";
  import Divider from "$lib/components/Divider.svelte";
  import { onMount } from "svelte";
  import { CR3Call, Typography } from "hexacomponents-svelte";

  const call = new CR3Call(globalThis.location.origin);

  let onEdit = $state(true);
  let uploadedImage = $state(null);
  let uploadedImageName = $state(null);
  let fileImage = $state(null);
  let categories = $state([
    {
      value: "1",
      name: "General",
    },
    {
      value: "2",
      name: "Contract",
    }
  ]);

  let personClient = $state([
    {
      value: "1",
      name: "Azizi",
    },
    {
      value: "2",
      name: "Habibi",
    }
  ]);

  let ischecked = $state(false);

  let data = $state({
    document_title: "",
    document_upload: "",
    description: "",
    document_type: "",
    person_client: "",
  });

  let resources = $state([]);

  let tabDatas = $state({
    selected: "Information",
    data: [
      {
        name: "Information",
      },
      {
        name: "Resource",
      },
    ],
  });

  let index = $state(resources.length);

  function addResource() {
    const newResource = {
      // dummy id
      id: Date.now(),
      resource_name: "",
      image: "",
      description: "",
      resource_type: "",
      resource_type_id: "",
      is_active: true,
    };
    resources = [...resources, newResource];
  }

  function deleteResource(id) {
    resources = resources.filter((resource) => resource.id !== id);
  }

  const handleImageUploadWrapper = (index) => {
    const inputElement = document.createElement("input");
    inputElement.type = "file";
    inputElement.accept = "image/*";

    inputElement.addEventListener("change", (event) => {
      const file = event.target.files[0];

      if (file && isValidImage(file)) {
        const reader = new FileReader();

        reader.onload = () => {
          fileImage = file;
          uploadedImageName = file.name;
          resources[index].image = reader.result;
        };

        reader.readAsDataURL(file);
      } else {
        alert("Invalid file. Please select a valid image file (max. 800x400px).");
      }
    });
    inputElement.click();
  };

  async function createProduct() {
    if (
      data.product_name === "" ||
      data.description === "" ||
      data.price === "" ||
      data.category === null
    ) {
      return;
    }
    let token = sessionStorage.getItem("token");

    if (token !== null) {
      try {
        let response = await call.call("client-dashboard", "create-product", {
          URLQuery: {},
          URLParams: {},
          Params: {
            product_name: data.product_name,
            description: data.description,
            price: data.price,
            is_active: data.is_active,
            category: data.category,
          },
          Headers: {
            Authorization: "Bearer " + token,
          },
        });

        console.log("kelinci", response);

        if (response.status === 200) {
          if (resources.length > 0) {
            for (let i = 0; i < resources.length; i++) {
              let response2 = await call.call("client-dashboard", "create-resource", {
                URLQuery: {},
                URLParams: {},
                Params: {
                  resource_name: resources[i].resource_name,
                  // image: resources[i].image !== null ? resources[i].image : null,
                  description: resources[i].description,
                  // resource_type: resources[i].resource_type,
                  is_active: resources[i].is_active === true ? "Yes" : "No",
                  product: response.result.data.id,
                },
                Headers: {
                  Authorization: "Bearer " + token,
                },
              });

              if (response2.status !== 200) {
                console.log(response2);
                return;
              }
            }
          }
          await call.goToPage("client-dashboard", "crm-product-service");
          return;
        } else {
          console.log(response);
          return;
        }
      } catch (error) {
        console.log(error);
      }
    }
  }

  async function getCategories() {
    let token = sessionStorage.getItem("token");
    if (token !== null) {
      let response = await call.call("client-dashboard", "get-categories", {
        URLQuery: {},
        URLParams: {},
        Params: {
          page: 1,
        },
        Headers: {
          Authorization: "Bearer " + token,
        },
      });

      if (response.status === 200) {
        let promise = response.result.data.map((category) => {
          return {
            value: category.id,
            name: category.attributes.category_name,
          };
        });

        categories = await Promise.all(promise);
      } else {
        console.log(response);
        return;
      }
    }
  }

  let buttonDisable = $state(true);

  $effect(() => {
    if (
      data.product_name !== "" &&
      data.description !== "" &&
      data.price !== "" &&
      data.is_active !== ""
    ) {
      buttonDisable = false;
    } else {
      buttonDisable = true;
    }
  });

  onMount(async () => {
    getCategories();
  });

  // let dropdownOpen = false;
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
              Add new document
            </Heading>
          </div>
        </div>
        <div class="flex gap-2">
          <Button
            onclick={async function () {
              await call.goToPage("client-dashboard", "crm-product-service");
              return;
            }}
            color="light"
            class="px-3 focus-within:!ring-1">{$_("cancel")}</Button
          >
          <Button onclick={createProduct} class="px-3 focus-within:!ring-1">{$_("create")}</Button>
        </div>
      </div>

      <Divider />

      <div class="flex flex-col gap-4 pt-8">
        <!-- general information -->
        <div class="flex gap-6">
          <div class="w-60">
            <P class="font-semibold text-secondary-light dark:text-secondary-dark">
              {$_("generalInformation")}
            </P>
          </div>
          <div class="flex flex-col grow gap-6">
            <div class="flex gap-5">
              <div class="flex flex-col grow gap-2">
                <P
                  class="text-sm font-medium leading-tight text-secondary-light dark:text-secondary-dark"
                >
                  {$_("Document Title")}*
                </P>
                <Input
                  readonly={!onEdit}
                  bind:value={data.document_title}
                  class="bg-primary-light dark:bg-secondary-dark"
                />
              </div>
            </div>
            <div class="flex flex-col col-span-5 gap-2">
              <P
                class="text-sm font-medium leading-tight text-secondary-light dark:text-secondary-dark"
                >{$_("Document Upload")}*</P
              >
              <Button
                onclick={handleImageUploadWrapper(index)}
                color="light"
                class="w-full rounded-xl border border-primary-light dark:border-primary-dark bg-primary-light dark:bg-primary-dark flex-col justify-start items-center gap-1 flex py-8"
              >
                {#if uploadedImage}
                  <img
                    src={uploadedImage}
                    alt="Uploaded Logo"
                    class="rounded-xl max-h-48 mb-4"
                  />
                  <P>{uploadedImageName}</P>
                {:else}
                  <div
                    class="rounded-xl shadow border border-primary-light dark:border-primary-dark justify-center items-center inline-flex mb-4"
                  >
                    <img src={upload} alt="mail.svg" class="p-4 dark:invert" />
                  </div>
                  <div class="flex space-x-1">
                    <P class="text-hc-brand-700">Click to upload</P><P>or drag and drop</P>
                  </div>
                  <P>SVG, PNG, JPG, or GIF (max. 800x400px)</P>
                {/if}
              </Button>
            </div>
            <div class="grid grid-cols-1">
              <div class="flex flex-col gap-2">
                <P
                  class="text-sm font-medium leading-tight text-secondary-light dark:text-secondary-dark"
                  >{$_("Description")}</P
                >
                <Textarea
                  readonly={!onEdit}
                  bind:value={data.description}
                  rows="3"
                  class="bg-primary-light dark:bg-secondary-dark"
                />
              </div>
            </div>
          </div>
        </div>
        <Divider />
        <!-- product detail -->
        <div class="flex gap-6">
          <div class="w-60">
            <P class="font-semibold text-secondary-light dark:text-secondary-dark">
              {$_("Document Type")}
            </P>
          </div>
          <div class="flex flex-col grow gap-6">
            <div class="flex flex-col grow gap-2">
              <P
                class="text-sm font-medium leading-tight text-secondary-light dark:text-secondary-dark"
                >{$_("Document Type")}*</P
              >
              <Select
                items={categories}
                bind:value={data.document_type}
                placeholder="Select document type"
              />
              <!-- <Input
                readonly={!onEdit}
                bind:value={data.category}
                class="bg-primary-light dark:bg-secondary-dark"
              />   -->
            </div>
          </div>
        </div>
        <Divider />
        <!-- product detail -->
        <div class="flex gap-6">
          <div class="w-60">
            <P class="font-semibold text-secondary-light dark:text-secondary-dark">
              {$_("Contract Information")}
            </P>
          </div>
          <div class="flex flex-col grow gap-6">
            <div class="flex flex-col grow gap-2">
              <P
                class="text-sm font-medium leading-tight text-secondary-light dark:text-secondary-dark"
                >{$_("Persons or Client")}*</P
              >
              <Select
                items={personClient}
                bind:value={data.person_client}
                placeholder="Search in persons or client database"
              />
              <!-- <Input
                readonly={!onEdit}
                bind:value={data.category}
                class="bg-primary-light dark:bg-secondary-dark"
              />   -->
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</div>
