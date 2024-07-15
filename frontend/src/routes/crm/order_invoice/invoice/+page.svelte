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
  import { ToastContainer, triggerToast, InputRw, TextareaRw } from "$lib";
  import { _ } from "svelte-i18n";
  import CopyIcon from "$lib/icon/copy-03.svelte";
  import { onMount } from "svelte";
  import { CR3Call } from "hexacomponents-svelte";
  import Divider from "$lib/components/Divider.svelte";
  import { format, formatDistanceToNow } from "date-fns";
  import { goto } from "$app/navigation";
  import { generateUniqueId } from "$lib";

  const call = new CR3Call(globalThis.location.origin);

  const params = new URLSearchParams(window.location.search);
  const invoiceId = params.get("invoiceId");
  let onEdit = $state(false);

  let data = $state({
    invoiceId: "",
    customerId: "",
    createdAt: "",
    invoice_date: "",
    due_date: "",
    total_amount: "",
    payment_status: "",
    notes: "",
    orderId: "",
  });

  let toastCreated = $state(false);
  let toastUpdated = $state(false);
  let toastError = $state(false);
  let toastDeleted = $state(false);
  let counter = $state(6);

  async function getInvoice() {
    if (invoiceId === null) {
      triggerToast("error");
      // setTimeout(async () => {
      goto(await call.getPage("client-dashboard", "crm-order-invoice"));
      // }, 3000);
      return;
    }
    let token = sessionStorage.getItem("token");

    if (token !== null) {
      try {
        let response = await call.call("client-dashboard", "get-invoice", {
          URLQuery: {},
          URLParams: {
            id: invoiceId,
          },
          Headers: {
            Authorization: `Bearer ${token}`,
          },
        });

        console.log(response);

        if (response.status !== 200) {
          console.log(response);
        }

        if (response.status === 200) {
          if (response.result.data.length > 0) {
            let invoice = response.result.data[0];
            data = {
              invoiceId: generateUniqueId(invoice.id, "invoice"),
              customerId: generateUniqueId(invoice.attributes.person?.data.id, "person"),
              createdAt: invoice.attributes.created_at,
              invoice_date: formatDistanceToNow(new Date(invoice.attributes.invoice_date), {
                addSuffix: true,
              }),
              due_date: formatDistanceToNow(new Date(invoice.attributes.due_date), {
                addSuffix: true,
              }),
              total_amount: invoice.attributes.order_booking?.data.attributes.total_amount,
              payment_status:
                invoice.attributes.payment_status[invoice.attributes.payment_status.length - 1]
                  .status,
              notes: invoice.attributes.note,
              orderId: generateUniqueId(invoice.attributes.order_booking?.data.id),
            };
          }
        }
      } catch (error) {
        console.log(error);
      }
    }
  }

  async function deleteInvoice() {
    let token = sessionStorage.getItem("token");

    if (token !== null) {
      try {
        let response = await call.call("client-dashboard", "delete-invoice", {
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
          triggerToast("deleted");
          // setTimeout(async () => {
          goto(await call.getPage("client-dashboard", "crm-order-invoice"));
          // }, 3000);
        }
      } catch (error) {
        triggerToast("error");
        console.log(error);
      }
    }
  }

  onMount(() => {
    getInvoice();
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
              {`Invoice #${data.invoiceId}`}
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
                goto(await call.getPage("client-dashboard", "crm-order-invoice"));
              }}
              color="light"
              class="px-3 focus-within:!ring-1">{$_("cancel")}</Button
            >
            <Button onclick={() => {}} class="px-3 focus-within:!ring-1">{$_("save")}</Button>
          {:else}
            <Button onclick={deleteInvoice} color="light" class="px-3 focus-within:!ring-1"
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
          <div class="flex grow flex-col gap-6">
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
            </div>
            <div class="flex gap-5">
              <div class="flex flex-col w-max gap-2">
                <P
                  class="text-sm font-medium leading-tight text-secondary-light dark:text-secondary-dark"
                >
                  {$_("invoiceId")}
                </P>
                <InputRw bind:value={data.invoiceId} />
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
                  {$_("orderId")}
                </P>
                <InputRw bind:value={data.orderId} />
              </div>
            </div>
            <div class="flex justify-center items-center gap-5">
              <div class="flex flex-col grow gap-2">
                <P
                  class="text-sm font-medium leading-tight text-secondary-light dark:text-secondary-dark"
                >
                  {$_("invoiceDate")}
                </P>
                <InputRw bind:value={data.invoice_date} />
              </div>
              <div class="mt-6">
                <P>-</P>
              </div>
              <div class="flex flex-col grow gap-2">
                <P
                  class="text-sm font-medium leading-tight text-secondary-light dark:text-secondary-dark"
                >
                  {$_("dueDate")}
                </P>
                <InputRw bind:value={data.due_date} />
              </div>
            </div>
            <div class="flex flex-col grow gap-2">
              <P
                class="text-sm font-medium leading-tight text-secondary-light dark:text-secondary-dark"
              >
                {$_("notes")}
              </P>
              <TextareaRw bind:value={data.notes} />
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
  <ToastContainer />
</div>
