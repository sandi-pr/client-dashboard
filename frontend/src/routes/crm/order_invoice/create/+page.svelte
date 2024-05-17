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
  import { P, Heading, Input, Button, Textarea, Select } from "flowbite-svelte";
  import { _ } from "svelte-i18n";
  import { CR3Call } from "hexacomponents-svelte";
  import { onMount } from "svelte";
  import { goto } from "$app/navigation";

  const call = new CR3Call(globalThis.location.origin);

  const params = new URLSearchParams(window.location.search);
  let onEdit = $state(true);

  let data = $state({
    customerId: "",
    order_date: "",
    order_time: "",
    productId: "",
    resourceId: "",
    shipping_address: "",
    total_amount: 0,
  });

  async function createOrder() {
    let token = sessionStorage.getItem("token");
    if (token !== null) {
      try {
        let response = await call.call("client-dashboard", "create-order", {
          URLQuery: {},
          URLParams: {},
          Params: {
            person: data.customerId,
            order_date: new Date(data.order_date + " " + data.order_time + ":00").toISOString(),
            product: data.productId,
            resource: data.resourceId,
            total_amount: data.total_amount,
            shipping_address: data.shipping_address,
            order_status: JSON.stringify([
              {
                status: "waiting for payment",
                time: new Date().toISOString(),
              },
            ]),
          },
          Headers: {
            Authorization: "Bearer " + token,
          },
        });
        console.log("create-order", response);
        if (response.status === 200) {
          let orderId = response.result.data.id;
          if (orderId !== null) {
            let currentDateTime = new Date();
            const oneHourLater = new Date(currentDateTime.getTime() + 60 * 60 * 1000);

            const invoice_date = currentDateTime.toISOString();
            const due_date = oneHourLater.toISOString();
            console.log("svelte", data.total_amount);
            let response2 = await call.call("client-dashboard", "create-invoice", {
              URLQuery: {},
              URLParams: {},
              Params: {
                order_booking: orderId,
                total_amount: data.total_amount,
                person: data.customerId,
                invoice_date: invoice_date,
                due_date: due_date,
                payment_status: JSON.stringify([
                  {
                    status: "Unpaid",
                    time: new Date().toISOString(),
                  },
                ]),
              },
              Headers: {
                Authorization: "Bearer " + token,
              },
            });

            console.log("create-invoice", response2);

            if (response2.status !== 200) {
              console.log(response2);
              return;
            }

            if (response2.status === 200) {
              goto(await call.getPage("client-dashboard", "crm-order-invoice"));
            }
          }
        } else {
          console.log(response);
        }
      } catch (error) {
        console.log(error);
      }
    }
  }

  let allProducts = $state([]);

  let products = $state([]);
  let resources = $state([]);

  let persons = $state([]);

  async function getProducts() {
    let token = sessionStorage.getItem("token");
    if (token !== null) {
      let page = 1;
      while (true) {
        let response = await call.call("client-dashboard", "get-products", {
          URLQuery: {},
          URLParams: {},
          Params: {
            page: page,
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
          allProducts = allProducts.concat(response.result.data);
          let promise = response.result.data.map((product) => {
            return {
              value: product.id,
              name: product.attributes.product_name,
            };
          });

          let data = await Promise.all(promise);
          products = products.concat(data);
          page++;
        } else {
          break;
        }
      }
    }
  }

  async function getPersons() {
    let token = sessionStorage.getItem("token");
    if (token !== null) {
      let page = 1;
      while (true) {
        let response = await call.call("client-dashboard", "get-persons", {
          URLQuery: {},
          URLParams: {},
          Params: {
            page: page,
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
          let promise = response.result.data.map((person) => {
            return {
              value: person.id,
              name:
                person.attributes.fullname !== null
                  ? person.attributes.fullname
                  : `<${person.attributes.email}>`,
            };
          });

          let data = await Promise.all(promise);
          persons = persons.concat(data);
          page++;
        } else {
          break;
        }
      }
    }
  }

  onMount(async () => {
    getProducts();
    getPersons();
  });

  $effect(async () => {
    if (!data.productId) return;

    if (data.productId !== "") {
      for (let i = 0; i < allProducts.length; i++) {
        if (allProducts[i].id === data.productId) {
          let promise = allProducts[i].attributes.resources.data.map((resource) => {
            return {
              value: resource.id,
              name: resource.attributes.resource_name,
            };
          });
          resources = await Promise.all(promise);

          console.log("resources", resources);
          break;
        }
      }
    }
  });
</script>

<!-- Form layout -->
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
              Create New Order
            </Heading>
          </div>
        </div>
        <div class="flex gap-2">
          {#if onEdit}
            <Button
              onclick={async function () {
                goto(await call.getPage("client-dashboard", "crm-categories"));
                return;
              }}
              color="light"
              class="px-3 focus-within:!ring-1">{$_("cancel")}</Button
            >
            <Button onclick={createOrder} class="px-3 focus-within:!ring-1">{$_("create")}</Button>
          {/if}
        </div>
      </div>
      <div class="flex flex-col gap-4 pt-8">
        <!-- Order Information -->
        <div class="flex gap-6">
          <div class="w-60">
            <P class="font-semibold text-secondary-light dark:text-secondary-dark">
              Order Information
            </P>
          </div>
          <div class="flex flex-col grow gap-6">
            <!-- Customer ID -->
            <div class="flex gap-5">
              <div class="flex flex-col grow gap-2">
                <P
                  class="text-sm font-medium leading-tight text-secondary-light dark:text-secondary-dark"
                  >Customer</P
                >
                <Select
                  bind:items={persons}
                  bind:value={data.customerId}
                  placeholder="Select Customer"
                />
              </div>
            </div>
            <!-- Order Date -->
            <div class="flex gap-5">
              <div class="flex flex-col grow gap-2">
                <P
                  class="text-sm font-medium leading-tight text-secondary-light dark:text-secondary-dark"
                >
                  Order Date
                </P>
                <Input
                  bind:value={data.order_date}
                  type="date"
                  class="bg-primary-light dark:bg-secondary-dark"
                />
              </div>
            </div>
            <div class="flex gap-5">
              <div class="flex flex-col grow gap-2">
                <P
                  class="text-sm font-medium leading-tight text-secondary-light dark:text-secondary-dark"
                >
                  Order Time
                </P>
                <Input
                  bind:value={data.order_time}
                  type="time"
                  class="bg-primary-light dark:bg-secondary-dark"
                />
              </div>
            </div>
            <!-- Total Amount -->
            <div class="flex gap-5">
              <div class="flex flex-col grow gap-2">
                <P
                  class="text-sm font-medium leading-tight text-secondary-light dark:text-secondary-dark"
                >
                  Total Amount
                </P>
                <Input
                  bind:value={data.total_amount}
                  type="number"
                  class="bg-primary-light dark:bg-secondary-dark"
                />
              </div>
            </div>
            <!-- Product Information -->
            <div class="flex gap-5">
              <div class="flex flex-col grow gap-2">
                <!-- Product Name -->
                <div class="flex gap-5">
                  <div class="flex flex-col grow gap-2">
                    <P
                      class="text-sm font-medium leading-tight text-secondary-light dark:text-secondary-dark"
                      >Product</P
                    >
                    <Select
                      bind:items={products}
                      bind:value={data.productId}
                      placeholder="Select Product"
                    />
                  </div>
                </div>
                <!-- Product Resource ID -->
                <div class="flex gap-5">
                  <div class="flex flex-col grow gap-2">
                    <P
                      class="text-sm font-medium leading-tight text-secondary-light dark:text-secondary-dark"
                      >Resource</P
                    >
                    <Select
                      bind:items={resources}
                      bind:value={data.resourceId}
                      placeholder="Select Resource"
                    />
                  </div>
                </div>
              </div>
            </div>
            <!-- Shipping Address -->
            <div class="flex gap-5">
              <div class="flex flex-col grow gap-2">
                <P
                  class="text-sm font-medium leading-tight text-secondary-light dark:text-secondary-dark"
                >
                  Shipping Address
                </P>
                <Textarea
                  bind:value={data.shipping_address}
                  rows="4"
                  class="bg-primary-light dark:bg-secondary-dark"
                />
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</div>
