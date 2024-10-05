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
  import { P, Heading, Input, Button } from "flowbite-svelte";
  import { ToastContainer, triggerToast, InputRw } from "$lib";
  import { _ } from "svelte-i18n";
  import CopyIcon from "$lib/icon/copy-03.svelte";
  import { onMount } from "svelte";
  import { CR3Call } from "hexacomponents-svelte";
  import Divider from "$lib/components/Divider.svelte";
  import { format, formatDistanceToNow } from "date-fns";
  import { goto } from "$app/navigation";
  import { generateUniqueId } from "$lib";

  import { call } from "$lib";

  const params = new URLSearchParams(window.location.search);
  const orderId = params.get("orderId");
  let invoiceId = $state("");
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

  async function getOrder() {
    if (orderId === null) {
      triggerToast("error");
      // setTimeout(async () => {
      goto(await $call.getPage("client-dashboard", "crm-order-invoice"));
      // }, 3000);
      return;
    }
    let token = localStorage.getItem("token");

    if (token !== null) {
      try {
        let response = await $call.call("client-dashboard", "get-order", {
          URLQuery: {},
          URLParams: {
            id: orderId,
          },
          Params: {},
          Headers: {
            Authorization: `Bearer ${token}`,
          },
        });

        if (response.status !== 200) {
          console.log(response);
        }

        if (response.status === 200) {
          if (response.result.data.length > 0) {
            console.log("order", response.result.data[0]);
            let order = response.result.data[0];
            invoiceId = order.attributes.invoice.data?.id;
            data = {
              orderId: generateUniqueId(order.id, "order"),
              customerId: generateUniqueId(order.attributes.person.data?.id, "person"),
              createdAt: order.attributes.created_at,
              order_date: formatDistanceToNow(new Date(order.attributes.order_date), {
                addSuffix: true,
              }),
              total_amount: order.attributes.total_amount,
              payment_status:
                order.attributes.invoice.data?.attributes.payment_status[
                  order.attributes.invoice.data?.attributes.payment_status.length - 1
                ]?.status,
              order_status:
                order.attributes.order_status[order.attributes.order_status.length - 1].status,
              // booking: {
              //   time: format(new Date(order.attributes.booking_date), "yyyy-MM-dd, kk:mm"),
              //   booking_date: format(new Date(order.attributes.booking_date), "yyyy-MM-dd, kk:mm"),
              //   notes: order.attributes.note,
              // },
              product: {
                id: generateUniqueId(order.attributes.product.data?.id, "product"),
                name: order.attributes.product.data?.attributes.product_name,
                resourceId: order.attributes.resources.data
                  .map((resource) => {
                    return generateUniqueId(resource.id, "resource");
                  })
                  .join(","),
                // resourceId: order.attributes.resources.data.map((resource) => generateUniqueId(resource.id)).join(","),
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

  // async function updateOrder() {
  //   let token = localStorage.getItem("token");
  //   try {
  //     let response = await $call.call("client-dashboard", "update-order", {
  //       URLQuery: {},
  //       URLParams: {
  //         id: orderId,
  //       },
  //       Params: {
  //         order_date: new Date(data.order_date).toISOString().split("T")[0],
  //         booking_date: data.booking_date + "T" + data.booking.time + ":00.000Z",
  //         total_amount: data.total_amount,
  //         note: data.notes,
  //         order_status:
  //           data.order_status === orderStatusCopy[orderStatusCopy.length - 1].status
  //             ? JSON.stringify(orderStatusCopy)
  //             : JSON.stringify([
  //                 ...orderStatusCopy,
  //                 { status: data.order_status, time: new Date().toISOString() },
  //               ]),
  //         shipping_address: data.shipping_address,
  //       },
  //       Headers: {
  //         Authorization: "Bearer " + token,
  //       },
  //     });

  //     console.log("response", response);

  //     if (response.status === 200) {
  //       goto(await $call.getPage("client-dashboard", "crm-order-invoice"));
  //     }
  //   } catch (error) {
  //     console.log(error);
  //   }
  // }

  async function deleteOrder() {
    let token = localStorage.getItem("token");
    if (invoiceId !== "" && invoiceId !== null && invoiceId !== undefined) {
      try {
        let response = await $call.call("client-dashboard", "delete-invoice", {
          URLQuery: {},
          URLParams: {
            id: invoiceId.toString(),
          },
          Headers: {
            Authorization: `Bearer ${token}`,
          },
        });

        if (response.status !== 200) {
          triggerToast("error");
          console.log(response);
          return;
        }

        if (response.status === 200) {
          let response = await $call.call("client-dashboard", "delete-order", {
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
            triggerToast("deleted");
            // setTimeout(async () => {
            goto(await $call.getPage("client-dashboard", "crm-order-invoice"));
            // }, 3000);
          }
        }
      } catch (error) {
        triggerToast("error");
        console.log(error);
      }
    } else {
      let response = await $call.call("client-dashboard", "delete-order", {
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
        triggerToast("deleted");
        // setTimeout(async () => {
        goto(await $call.getPage("client-dashboard", "crm-order-invoice"));
        // }, 3000);
      }
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
              onclick={async function () {
                goto(await $call.getPage("client-dashboard", "crm-order-invoice"));
              }}
              color="light"
              class="px-3 focus-within:!ring-1">{$_("cancel")}</Button
            >
            <Button onclick={updateOrder} class="px-3 focus-within:!ring-1">{$_("save")}</Button>
          {:else}
            <Button onclick={deleteOrder} color="light" class="px-3 focus-within:!ring-1"
              >{$_("delete")}</Button
            >
            <!-- <Button
              onclick={() => {
                onEdit = true;
              }}
              class="px-3 focus-within:!ring-1">{$_("edit")}</Button
            > -->
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
                  {$_("orderId")}
                </P>
                <InputRw bind:value={data.orderId} />
              </div>
              <div class="flex flex-col w-max gap-2">
                <P
                  class="text-sm font-medium leading-tight text-secondary-light dark:text-secondary-dark"
                >
                  {$_("customerId")}
                </P>
                <InputRw bind:value={data.customerId} />
              </div>
              <div class="flex flex-col w-max gap-2">
                <P
                  class="text-sm font-medium leading-tight text-secondary-light dark:text-secondary-dark"
                >
                  {$_("orderDate")}
                </P>
                <InputRw bind:value={data.order_date} />
              </div>
              <!-- <div class="flex flex-col w-max gap-2">
                <P
                  class="text-sm font-medium leading-tight text-secondary-light dark:text-secondary-dark"
                >
                  Booking Date
                </P>
                <InputRw bind:value={data.booking_date} />
              </div> -->
            </div>
            <div class="flex gap-5">
              <div class="flex flex-col w-max gap-2">
                <P
                  class="text-sm font-medium leading-tight text-secondary-light dark:text-secondary-dark"
                >
                  {$_("totalAmount")}
                </P>
                <InputRw bind:value={data.total_amount} />
              </div>
              <div class="flex flex-col w-max gap-2">
                <P
                  class="text-sm font-medium leading-tight text-secondary-light dark:text-secondary-dark"
                >
                  {$_("paymentStatus")}
                </P>
                <InputRw bind:value={data.payment_status} />
              </div>
              <div class="flex flex-col w-max gap-2">
                <P
                  class="text-sm font-medium leading-tight text-secondary-light dark:text-secondary-dark"
                >
                  {$_("orderStatus")}
                </P>
                <InputRw bind:value={data.order_status} />
              </div>
            </div>
          </div>
        </div>
        <Divider />
        <!-- general information -->
        <!-- <div class="flex gap-6">
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
        <Divider /> -->
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
                  {$_("productName")}
                </P>
                <InputRw bind:value={data.product.name} />
              </div>
              <div class="flex flex-col w-max gap-2">
                <P
                  class="text-sm font-medium leading-tight text-secondary-light dark:text-secondary-dark"
                >
                  {$_("productId")}
                </P>
                <InputRw bind:value={data.product.id} />
              </div>
              <div class="flex flex-col w-max gap-2">
                <P
                  class="text-sm font-medium leading-tight text-secondary-light dark:text-secondary-dark"
                >
                  {$_("resourceId")}
                </P>
                <InputRw bind:value={data.product.resourceId} />
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
                  {$_("shippingAddress")}
                </P>
                <InputRw bind:value={data.shipping_address} />
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
  <ToastContainer />
</div>
