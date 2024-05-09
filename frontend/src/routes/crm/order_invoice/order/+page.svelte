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
  import { _ } from "svelte-i18n";
  import CopyIcon from "$lib/icon/copy-03.svelte";
  import { onMount } from "svelte";
  import { CR3Call } from "hexacomponents-svelte";
  import Divider from "$lib/components/Divider.svelte";
  import { format } from "date-fns";

  const call = new CR3Call(globalThis.location.origin);

  const params = new URLSearchParams(window.location.search);
  const orderId = params.get("orderId");
  let onEdit = $state(false);

  let data = $state({
    orderId: "",
    customerId: "",
    createdAt: "",
    order_date: "",
    booking_date: "",
    total_amount: "",
    payment_status: "",
    order_status: "",
    booking: {
      time: "",
      booking_date: "",
      notes: "",
    },
    product: {
      id: "",
      name: "",
      resourceId: "",
    },
    shipping_address: "",
  });

  let orderStatusCopy = $state([]);
  let paymentStatusCopy = $state([]);

  async function getOrder() {
    let token = sessionStorage.getItem("token");

    if (token !== null) {
      try {
        let response = await call.call("client-dashboard", "get-order", {
          URLQuery: {},
          URLParams: {
            id: orderId,
          },
          Params: {},
          Headers: {
            Authorization: `Bearer ${token}`,
          },
        });

        if (response.status === 200) {
          if (response.result.data.length > 0) {
            console.log("order", response.result.data[0]);
            let order = response.result.data[0];
            data = {
              orderId: order.id,
              customerId: order.attributes.person.data.id,
              createdAt: order.attributes.created_at,
              order_date: format(new Date(order.attributes.order_date), "yyyy-MM-dd"),
              booking_date: format(new Date(order.attributes.booking_date), "yyyy-MM-dd"),
              total_amount: order.attributes.total_amount,
              payment_status:
                order.attributes.invoice.data.attributes.payment_status[
                  order.attributes.invoice.data.attributes.payment_status.length - 1
                ].status,
              order_status:
                order.attributes.order_status[order.attributes.order_status.length - 1].status,
              booking: {
                time: format(new Date(order.attributes.booking_date), "kk:mm"),
                booking_date: format(new Date(order.attributes.booking_date), "yyyy-MM-dd"),
                notes: order.attributes.note,
              },
              product: {
                id: order.attributes.product.data.id,
                name: order.attributes.product.data.attributes.product_name,
                resourceId: order.attributes.resources.data[0].id,
              },
              shipping_address: order.attributes.shipping_address,
            };

            orderStatusCopy = order.attributes.order_status;
          }
        }
      } catch (error) {
        console.log(error);
      }
    }
  }

  async function updateOrder() {
    let token = sessionStorage.getItem("token");
    try {
      let response = await call.call("client-dashboard", "update-order", {
        URLQuery: {},
        URLParams: {
          id: orderId,
        },
        Params: {
          order_date: new Date(data.order_date).toISOString().split("T")[0],
          booking_date: data.booking_date + "T" + data.booking.time + ":00.000Z",
          total_amount: data.total_amount,
          note: data.notes,
          order_status:
            data.order_status === orderStatusCopy[orderStatusCopy.length - 1].status
              ? JSON.stringify(orderStatusCopy)
              : JSON.stringify([
                  ...orderStatusCopy,
                  { status: data.order_status, time: new Date().toISOString() },
                ]),
          shipping_address: data.shipping_address,
        },
        Headers: {
          Authorization: "Bearer " + token,
        },
      });

      console.log("response", response);

      if (response.status === 200) {
        await call.goToPage("client-dashboard", "crm-order-invoice");
      }
    } catch (error) {
      console.log(error);
    }
  }

  async function deleteOrder() {
    let token = sessionStorage.getItem("token");
    try {
      let response = await call.call("client-dashboard", "delete-order", {
        URLQuery: {},
        URLParams: {
          id: orderId,
        },
        Params: {},
        Headers: {
          Authorization: "Bearer " + token,
        },
      });

      if (response.status === 200) {
        await call.goToPage("client-dashboard", "crm-order-invoice");
      }
    } catch (error) {
      console.log(error);
    }
  }

  onMount(() => {
    getOrder();
  });
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
              {`Order #${data.orderId}`}
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
            <Button onclick={updateOrder} class="px-3 focus-within:!ring-1">{$_("save")}</Button>
          {:else}
            <Button onclick={deleteOrder} color="light" class="px-3 focus-within:!ring-1"
              >{$_("delete")}</Button
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
      <Divider />
      <div class="flex flex-col gap-8 pt-8">
        <!-- general information -->
        <div class="flex gap-6">
          <div class="w-60">
            <P class="font-semibold text-secondary-light dark:text-secondary-dark">
              {$_("generalInformation")}
            </P>
          </div>
          <div class="flex flex-col gap-6">
            <div class="flex gap-5">
              <div class="flex flex-col w-max gap-2">
                <P
                  class="text-sm font-medium leading-tight text-secondary-light dark:text-secondary-dark"
                >
                  Order ID
                </P>
                <Input
                  readonly
                  bind:value={data.orderId}
                  class="bg-primary-light dark:bg-secondary-dark"
                />
              </div>
              <div class="flex flex-col w-max gap-2">
                <P
                  class="text-sm font-medium leading-tight text-secondary-light dark:text-secondary-dark"
                >
                  Customer ID
                </P>
                <Input
                  readonly
                  bind:value={data.customerId}
                  class="bg-primary-light dark:bg-secondary-dark"
                />
              </div>
              <div class="flex flex-col w-max gap-2">
                <P
                  class="text-sm font-medium leading-tight text-secondary-light dark:text-secondary-dark"
                >
                  Order Date
                </P>
                <Input
                  readonly={!onEdit}
                  bind:value={data.order_date}
                  type="date"
                  class="bg-primary-light dark:bg-secondary-dark"
                />
              </div>
              <div class="flex flex-col w-max gap-2">
                <P
                  class="text-sm font-medium leading-tight text-secondary-light dark:text-secondary-dark"
                >
                  Booking Date
                </P>
                <Input
                  readonly={!onEdit}
                  bind:value={data.booking_date}
                  type="date"
                  class="bg-primary-light dark:bg-secondary-dark"
                />
              </div>
            </div>
            <div class="flex gap-5">
              <div class="flex flex-col w-max gap-2">
                <P
                  class="text-sm font-medium leading-tight text-secondary-light dark:text-secondary-dark"
                >
                  Total Amounts
                </P>
                <Input
                  readonly={!onEdit}
                  bind:value={data.total_amount}
                  type="number"
                  class="bg-primary-light dark:bg-secondary-dark"
                />
              </div>
              <div class="flex flex-col w-max gap-2">
                <P
                  class="text-sm font-medium leading-tight text-secondary-light dark:text-secondary-dark"
                >
                  Payment Status
                </P>
                <Input
                  readonly
                  bind:value={data.payment_status}
                  class="bg-primary-light dark:bg-secondary-dark"
                />
              </div>
              <div class="flex flex-col w-max gap-2">
                <P
                  class="text-sm font-medium leading-tight text-secondary-light dark:text-secondary-dark"
                >
                  Order Status
                </P>
                <Input
                  readonly={!onEdit}
                  bind:value={data.order_status}
                  type="text"
                  class="bg-primary-light dark:bg-secondary-dark"
                />
              </div>
            </div>
          </div>
        </div>
        <Divider />
        <!-- general information -->
        <div class="flex gap-6">
          <div class="w-60">
            <P class="font-semibold text-secondary-light dark:text-secondary-dark">
              {$_("bookingInformation")}
            </P>
          </div>
          <div class="flex grow flex-col gap-6">
            <div class="flex justify-center items-center gap-5">
              <div class="flex flex-col grow gap-2">
                <P
                  class="text-sm font-medium leading-tight text-secondary-light dark:text-secondary-dark"
                >
                  Time
                </P>
                <Input
                  readonly={!onEdit}
                  bind:value={data.booking.time}
                  type="time"
                  class="bg-primary-light dark:bg-secondary-dark"
                />
              </div>
              <div class="mt-6">
                <P>-</P>
              </div>
              <div class="flex flex-col grow gap-2">
                <P
                  class="text-sm font-medium leading-tight text-secondary-light dark:text-secondary-dark"
                >
                  Booking Date
                </P>
                <Input
                  readonly={!onEdit}
                  bind:value={data.booking.booking_date}
                  type="date"
                  class="bg-primary-light dark:bg-secondary-dark"
                />
              </div>
            </div>
            <div class="flex flex-col grow gap-2">
              <P
                class="text-sm font-medium leading-tight text-secondary-light dark:text-secondary-dark"
              >
                Notes
              </P>
              <Textarea
                readonly={!onEdit}
                bind:value={data.booking.notes}
                rows="4"
                class="bg-primary-light dark:bg-secondary-dark"
              />
            </div>
          </div>
        </div>
        <Divider />
        <!-- general information -->
        <div class="flex gap-6">
          <div class="w-60">
            <P class="font-semibold text-secondary-light dark:text-secondary-dark">
              {$_("productInformation")}
            </P>
          </div>
          <div class="flex grow flex-col gap-6">
            <div class="flex justify-center items-center gap-5">
              <div class="flex flex-col grow gap-2">
                <P
                  class="text-sm font-medium leading-tight text-secondary-light dark:text-secondary-dark"
                >
                  Time
                </P>
                <Input
                  readonly={!onEdit}
                  bind:value={data.product.name}
                  class="bg-primary-light dark:bg-secondary-dark"
                />
              </div>
              <div class="flex flex-col w-max gap-2">
                <P
                  class="text-sm font-medium leading-tight text-secondary-light dark:text-secondary-dark"
                >
                  Product ID
                </P>
                <Input
                  readonly
                  bind:value={data.product.id}
                  class="bg-primary-light dark:bg-secondary-dark"
                />
              </div>
              <div class="flex flex-col w-max gap-2">
                <P
                  class="text-sm font-medium leading-tight text-secondary-light dark:text-secondary-dark"
                >
                  Resource ID
                </P>
                <Input
                  readonly
                  bind:value={data.product.resourceId}
                  class="bg-primary-light dark:bg-secondary-dark"
                />
              </div>
            </div>
          </div>
        </div>
        <Divider />
        <!-- general information -->
        <div class="flex gap-6">
          <div class="w-60">
            <P class="font-semibold text-secondary-light dark:text-secondary-dark">
              {$_("addressInformation")}
            </P>
          </div>
          <div class="flex grow flex-col gap-6">
            <div class="flex justify-center items-center gap-5">
              <div class="flex flex-col grow gap-2">
                <P
                  class="text-sm font-medium leading-tight text-secondary-light dark:text-secondary-dark"
                >
                  Shipping Address
                </P>
                <Input
                  readonly={!onEdit}
                  bind:value={data.shipping_address}
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
