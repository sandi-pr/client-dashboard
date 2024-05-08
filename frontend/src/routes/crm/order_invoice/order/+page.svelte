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

  const call = new CR3Call(globalThis.location.origin);

  const params = new URLSearchParams(window.location.search);
  const catergoryId = params.get("");
  let onEdit = $state(false);

  let data = $state({
    orderId: "19846829",
    customerId: "19846829",
    createdAt: "2023-07-02 08:51 PM",
    order_date: "10/29/2023 11:08",
    booking_date: "10/29/2023 11:08",
    total_amount: "$9836",
    payment_status: "Unpaid",
    order_status: "Calendar Booked",
    booking: {
      time: "24 Jan 2023, 08:00 PM",
      booking_date: "24 Jan 2023, 09:00 PM",
      notes: "-",
    },
    product: {
      id: "19846829",
      name: "Solid Hardwood Flooring",
      resourceId: "19846829",
    },
    shipping_address: "John Smith, 42 Primrose Avenue, London, EC2M 7EB, United Kingdom",
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
            <Button onclick={() => {}} class="px-3 focus-within:!ring-1">{$_("save")}</Button>
          {:else}
            <Button onclick={() => {}} color="light" class="px-3 focus-within:!ring-1"
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
                  readonly={!onEdit}
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
                  readonly={!onEdit}
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
                  readonly={!onEdit}
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
                  bind:value={data.booking.time}
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
                  readonly={!onEdit}
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
                  readonly={!onEdit}
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
