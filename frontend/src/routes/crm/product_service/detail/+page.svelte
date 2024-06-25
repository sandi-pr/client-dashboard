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
  import { P, Heading, Input, Button, Select, Textarea, Toggle, Modal } from "flowbite-svelte";
  import { ToastContainer, triggerToast } from "$lib";
  import { ExclamationCircleSolid } from "flowbite-svelte-icons";
  import upload from "$lib/icon/upload-cloud-02.svg";
  import { _ } from "svelte-i18n";
  import CopyIcon from "$lib/icon/copy-03.svelte";
  import TabBar from "$lib/ui/tab_bar/TabBar.svelte";
  import TabDetail from "$lib/components/TabDetail.svelte";
  import Divider from "$lib/components/Divider.svelte";
  import { onMount } from "svelte";
  import { CR3Call, Typography } from "hexacomponents-svelte";
  import { goto } from "$app/navigation";

  const call = new CR3Call(globalThis.location.origin);

  const params = new URLSearchParams(window.location.search);
  const productId = params.get("productId");
  const edit = params.get("edit");
  let onEdit = $state(false);
  let uploadedImage = $state(null);
  let uploadedImageName = $state(null);
  let fileImage = $state(null);
  let resourceReady = $state([]);

  let ischecked = $state(false);
  let categories = $state([]);

  let data = $state({});

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

  let toastCreated = $state(false);
  let toastUpdated = $state(false);
  let toastError = $state(false);
  let toastDeleted = $state(false);
  let counter = $state(6);

  let popupModal = $state(false);

  let index = $state(resources.length);

  function addResource() {
    if (onEdit) {
      const newResource = {
        // dummy id
        id: Date.now(),
        resource_name: "",
        image: "",
        description: "",
        resource_type_name: "",
        resource_type_id: "",
        resource_type_active: false,
      };
      resources = [...resources, newResource];
    }
  }

  async function deleteResource(id) {
    if (!onEdit) return;
    if (onEdit) {
      console.log(id);
      try {
        let token = sessionStorage.getItem("token");
        let response = await call.call("client-dashboard", "delete-resource", {
          URLQuery: {},
          URLParams: { id: id.toString() },
          Params: {},
          Headers: {
            Authorization: "Bearer " + token,
          },
        });

        if (response.status !== 200) {
          console.log(response);
          return;
        }

        console.log(response);
      } catch (error) {
        console.log(error);
      } finally {
        resources = resources.filter((resource) => {
          if (resource.id === id) {
            return false;
          }
          return true;
        });
      }
    }
  }

  const handleImageUploadWrapper = () => {
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
          uploadedImage = reader.result;
        };

        reader.readAsDataURL(file);
      } else {
        alert("Invalid file. Please select a valid image file (max. 800x400px).");
      }
    });
    inputElement.click();
  };

  function formatDate(inputDate) {
    const date = new Date(inputDate);

    const formattedDate = `${(date.getMonth() + 1).toString().padStart(2, "0")}/${date
      .getDate()
      .toString()
      .padStart(2, "0")}/${date.getFullYear()}, ${date
      .getHours()
      .toString()
      .padStart(2, "0")}:${date
      .getMinutes()
      .toString()
      .padStart(2, "0")} ${date.getHours() >= 12 ? "PM" : "AM"}`;

    return formattedDate;
  }

  async function getProduct() {
    if (productId === null) {
      triggerToast("error");
      // setTimeout(async () => {
        goto(await call.getPage("client-dashboard", "crm-product-service"));
      // }, 3000);
      return;
    };
    let token = sessionStorage.getItem("token");
    try {
      let response = await call.call("client-dashboard", "get-product", {
        URLQuery: {},
        URLParams: { id: productId },
        Params: {},
        Headers: {
          Authorization: "Bearer " + token,
        },
      });

      if (response.status !== 200) {
        console.log(response);
        return;
      }

      console.log("ayam", response);

      if (response.status === 200 && response.result.data.length > 0) {
        data = {
          productId: response.result.data[0].id,
          product_name: response.result.data[0].attributes.product_name,
          createdAt: formatDate(response.result.data[0].attributes.createdAt),
          description: response.result.data[0].attributes.description,
          price: response.result.data[0].attributes.price,
          category: response.result.data[0].attributes.category_product.data.id || null,
          is_active: response.result.data[0].attributes.is_active === "Yes" ? true : false,
        };

        if (response.result.data[0].attributes.resources.data.length > 0) {
          let promise = response.result.data[0].attributes.resources.data.map(
            async (resource, index) => {
              return {
                id: resource.id,
                resource_name: resource.attributes.resource_name,
                image: "",
                description: resource.attributes.description,
                resource_type_name: "",
                resource_type_id: "",
                is_active: resource.attributes.is_active === "Yes" ? true : false,
              };
            },
          );
          resources = await Promise.all(promise);
          for (let i = 0; i < resources.length; i++) {
            resourceReady = [...resourceReady, resources[i].id];
          }
          console.log("ayam", resourceReady);
        }
      }
    } catch (error) {
      console.log(error);
    }
  }

  async function updateProduct() {
    let token = sessionStorage.getItem("token");
    try {
      let response;
      if (tabDatas.selected === "Information") {
        response = await call.call("client-dashboard", "update-product", {
          URLQuery: {},
          URLParams: { id: productId },
          Params: {
            product_name: data.product_name,
            description: data.description,
            price: data.price,
            is_active: data.is_active ? "Yes" : "No",
            category: data.category,
          },
          Headers: {
            Authorization: "Bearer " + token,
          },
        });

        if (response.status !== 200) {
          console.log(response);
          return;
        }

        if (response.status === 200) {
          triggerToast("updated");
          // setTimeout(async () => {
            goto(await call.getPage("client-dashboard", "crm-product-service"));
          // }, 3000);
          return;
        }
      } else if (tabDatas.selected === "Resource") {
        for (let i = 0; i < resources.length; i++) {
          if (resourceReady.includes(resources[i].id)) {
            response = await call.call("client-dashboard", "update-resource", {
              URLQuery: {},
              URLParams: { id: resources[i].id.toString() },
              Params: {
                resource_name: resources[i].resource_name,
                description: resources[i].description,
                is_active: resources[i].is_active ? "Yes" : "No",
              },
              Headers: {
                Authorization: "Bearer " + token,
              },
            });
          } else {
            response = await call.call("client-dashboard", "create-resource", {
              URLQuery: {},
              URLParams: {},
              Params: {
                resource_name: resources[i].resource_name,
                description: resources[i].description,
                is_active: resources[i].is_active ? "Yes" : "No",
                product: parseInt(productId),
              },
              Headers: {
                Authorization: "Bearer " + token,
              },
            });
          }
        }

        if (response.status !== 200) {
          console.log(response);
          return;
        }

        if (response.status === 200) {
          triggerToast("updated");
          // setTimeout(async () => {
            goto(await call.getPage("client-dashboard", "crm-product-service"));
          // }, 3000);
          return;
        }
      }
    } catch (error) {
      console.log(error);
    }
  }

  async function deleteProduct() {
    let token = sessionStorage.getItem("token");
    try {
      let response = await call.call("client-dashboard", "delete-product", {
        URLQuery: {},
        URLParams: { id: productId },
        Params: {},
        Headers: {
          Authorization: "Bearer " + token,
        },
      });

      if (response.status !== 200) {
        console.log(response);
        return;
      }

      if (response.status === 200) {
        triggerToast("deleted");
        // setTimeout(async () => {
          goto(await call.getPage("client-dashboard", "crm-product-service"));
        // }, 3000);
        return;
      }
    } catch (error) {
      console.log(error);
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
          page_size: 25,
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

  async function backPage() {
    goto(await call.getPage("client-dashboard", "crm-product-service"));
    return;
  }

  onMount(async () => {
    if (edit === "true") onEdit = true;
    getProduct();
    getCategories();
  });

  // let dropdownOpen = false;
</script>

<div class="flex dark:bg-hc-graydark-950">
  <div class="w-full pb-10">
    <div class="flex flex-col m-4">
      <!-- Header -->
      <div class="flex items-center pb-6">
        <div class="flex flex-col gap-2 grow">
          <!-- Header -->
          <div class="gap-2">
            <Heading
              tag="h3"
              class="text-3xl leading-[38px] text-hc-graylight-900 dark:text-hc-graydark-50"
            >
              {data.product_name}
              <Button
                outline
                color="red"
                class="hover:!bg-transparent !border-hc-graylight-200 dark:!border-hc-graydark-600 focus-within:!ring-1 !p-2"
              >
                <CopyIcon svgClass="w-5 h-5 text-fg-secondary-light dark:text-fg-secondary-dark" />
              </Button>
            </Heading>
          </div>
          <P
            class="text-base font-medium leading-normal text-hc-graylight-500 dark:text-hc-graydark-400"
          >
            {$_("createOn")}{` ${data.createdAt}`}
          </P>
        </div>
        <div class="flex gap-2">
          {#if onEdit}
            <Button
              onclick={() => {
                onEdit = false;
              }}
              color="light"
              class="px-3 focus-within:!ring-1">{$_("cancel")}</Button
            >
            <Button onclick={updateProduct} class="px-3 focus-within:!ring-1">{$_("save")}</Button>
          {:else}
            <Button color="light" class="px-3 focus-within:!ring-1" onclick={backPage}
              >{$_("back")}</Button
            >
            <Button
              onclick={() => (popupModal = true)}
              outline color="red"
              class="px-3 focus-within:!ring-1">{$_("delete")}</Button
            >
            <Button
              onclick={() => {
                onEdit = true;
              }}
              class="px-3 focus-within:!ring-1">{$_("edit")}</Button
            >
          {/if}
        </div>
      </div>

      <!-- Submenu -->
      <TabDetail {tabDatas} on:tabChange={(e) => (tabDatas.selected = e.detail)} />

      {#if tabDatas.selected === "Information"}
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
                    {$_("productName")}
                  </P>
                  <Input
                    readonly={!onEdit}
                    bind:value={data.product_name}
                    class="bg-primary-light dark:bg-secondary-dark"
                  />
                </div>
                <div class="flex flex-col gap-2">
                  <P
                    class="text-sm font-medium leading-tight text-secondary-light dark:text-secondary-dark"
                  >
                    {$_("productId")}
                  </P>
                  <Input
                    type="number"
                    readonly
                    bind:value={data.productId}
                    class="bg-primary-light dark:bg-secondary-dark"
                  />
                </div>
              </div>
              <div class="grid grid-cols-1">
                <div class="flex flex-col gap-2">
                  <P
                    class="text-sm font-medium leading-tight text-secondary-light dark:text-secondary-dark"
                    >{$_("description")}</P
                  >
                  <Textarea
                    readonly={!onEdit}
                    bind:value={data.description}
                    rows="4"
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
                {$_("productDetail")}
              </P>
            </div>
            <div class="flex flex-col grow gap-6">
              <div class="grid grid-cols-5 gap-5">
                <div class="flex flex-col col-span-2 gap-2">
                  <P
                    class="text-sm font-medium leading-tight text-secondary-light dark:text-secondary-dark"
                    >{$_("pricing")}</P
                  >
                  <Input
                    type="number"
                    readonly={!onEdit}
                    bind:value={data.price}
                    class="bg-primary-light dark:bg-secondary-dark"
                  />
                </div>
                <div class="flex flex-col col-span-2 gap-2">
                  <P
                    class="text-sm font-medium leading-tight text-secondary-light dark:text-secondary-dark"
                    >{$_("category")}</P
                  >
                  <Select
                    disabled={!onEdit}
                    items={categories}
                    bind:value={data.category}
                    placeholder={$_("selectCategory")}
                    class="bg-primary-light dark:bg-secondary-dark"
                  />
                </div>
                <div class="flex flex-col gap-2 justify-center w-max">
                  <P
                    class="text-sm font-medium leading-tight text-secondary-light dark:text-secondary-dark"
                    >{$_("active")}</P
                  >
                  <Toggle disabled={!onEdit} bind:checked={data.is_active}></Toggle>
                </div>
              </div>
            </div>
          </div>
        </div>
      {:else if tabDatas.selected === "Resource"}
        <div class="flex gap-6 pt-8">
          <div class="w-60 space-y-2">
            <P class="font-semibold text-secondary-light dark:text-secondary-dark">
              {$_("resources")}
            </P>
            {#if onEdit}
              <Button color="light" on:click={addResource}>Add</Button>
            {/if}
          </div>
          <div class="flex flex-col grow gap-6">
            {#each resources as resource (resource.id)}
              <div class="grid grid-cols-5 gap-5 border rounded-lg p-4">
                <div class="flex items-center col-span-5 gap-2">
                  <div class="flex items-center w-full justify-between gap-2">
                    <div></div>
                    <Button on:click={() => deleteResource(resource.id)} color="light" class="w-max"
                      >{$_("delete")}</Button
                    >
                  </div>
                </div>
                <div class="flex flex-col col-span-5 gap-2">
                  <!-- <P
                    class="text-sm font-medium leading-tight text-secondary-light dark:text-secondary-dark"
                    >{$_("image")}</P
                  >
                  <Button
                    on:click={handleImageUploadWrapper}
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
                  </Button> -->
                </div>
                <div class="flex flex-col col-span-5 gap-2">
                  <P
                    class="text-sm font-medium leading-tight text-secondary-light dark:text-secondary-dark"
                    >{$_("name")}</P
                  >
                  <Input
                    readonly={!onEdit}
                    bind:value={resource.resource_name}
                    class="bg-primary-light dark:bg-secondary-dark"
                  />
                </div>
                <div class="flex flex-col col-span-5 gap-2">
                  <P
                    class="text-sm font-medium leading-tight text-secondary-light dark:text-secondary-dark"
                    >{$_("description")}</P
                  >
                  <Input
                    readonly={!onEdit}
                    bind:value={resource.description}
                    class="bg-primary-light dark:bg-secondary-dark"
                  />
                </div>
                <!-- <div class="flex flex-col col-span-2 gap-2">
                  <P
                    class="text-sm font-medium leading-tight text-secondary-light dark:text-secondary-dark"
                    >{$_("typeName")}</P
                  >
                  <Input
                    readonly={!onEdit}
                    bind:value={resource.resource_type_name}
                    class="bg-primary-light dark:bg-secondary-dark"
                  />
                </div>
                <div class="flex flex-col col-span-2 gap-2">
                  <P
                    class="text-sm font-medium leading-tight text-secondary-light dark:text-secondary-dark"
                    >{$_("resourceTypeId")}</P
                  >
                  <Input
                    readonly={!onEdit}
                    bind:value={resource.resource_type_id}
                    class="bg-primary-light dark:bg-secondary-dark"
                  />
                </div> -->
                <div class="flex flex-col gap-2 justify-center w-max">
                  <P
                    class="text-sm font-medium leading-tight text-secondary-light dark:text-secondary-dark"
                    >{$_("active")}</P
                  >
                  <Toggle disabled={!onEdit} bind:checked={resource.is_active}></Toggle>
                </div>
              </div>
            {/each}
          </div>
        </div>
      {/if}
    </div>
  </div>
  <ToastContainer />
</div>

<Modal bind:open={popupModal} size="xs" autoclose>
  <div class="text-center">
    <ExclamationCircleSolid class="mx-auto mb-4 text-red-600 w-12 h-12" />
    <h3 class="mb-5 text-lg font-normal text-gray-500 dark:text-gray-400">
      {$_("areYouSureYouWantToDelete")}
    </h3>
    <Button onclick={deleteProduct} color="red" class="me-2">{$_("yesImSure")}</Button>
    <Button color="alternative">{$_("noCancel")}</Button>
  </div>
</Modal>