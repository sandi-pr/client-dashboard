<script>
  import { P, Heading, Input, Button } from "flowbite-svelte";
  import { ChevronDownOutline } from "flowbite-svelte-icons";
  import { _ } from "svelte-i18n";
  // import CheckWhiteIcon from "$lib/icon/check-white.svg";
  // import XIcon from "$lib/icon/x.svg";
  // import ShareIcon from "$lib/icon/share-06.svg";
  import CopyIcon from "../../lib/icon/copy-03.svelte";
  import Divider from "../../lib/components/Divider.svelte";
  // import { toCamelCase } from "$lib/utils/toCamelCase.js";

  let data = {
    person: {
      customerId: "19846829",
      fullName: "Azzisi Bassem",
      dateOfBirth: "01/12/1994",
      email: "azzis144635@gmail.com",
      phone: "+44 33771021",
      country: "UK",
      address: "John Smith 42, Primrose Avenue, London EC2M 7EB, United Kingdom",
      city: "London",
      state: "London",
      zipCode: "EC2M 7EB",
      createdAt: "01-12-2024 03:02",
      active: true,
    },
  };

  // let dropdownOpen = false;
</script>

<div class="flex dark:bg-hc-graydark-950">
  <div class="w-full pb-10">
    <div class="flex flex-col m-4">
      <div class="flex items-center pb-8">
        <div class="flex flex-col gap-2 grow">
          <!-- Header -->
          <div class="gap-2">
            <Heading
              tag="h3"
              class="text-3xl leading-[38px] text-hc-graylight-900 dark:text-hc-graydark-50"
            >
              {data.person.fullName}
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
            {$_("createOn")}{` ${data.person.createdAt}`}
          </P>
        </div>
        <div class="flex gap-2">
          <Button color="light" class="px-3 focus-within:!ring-1">{$_("delete")}</Button>
          <Button class="px-3 focus-within:!ring-1">{$_("save")}</Button>
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
                  {$_("fullname")}
                </P>
                <Input
                  readonly
                  value={data.person.fullName}
                  class="bg-primary-light dark:bg-secondary-dark"
                />
              </div>
              <div class="flex flex-col gap-2">
                <P
                  class="text-sm font-medium leading-tight text-secondary-light dark:text-secondary-dark"
                >
                  {$_("customerId")}
                </P>
                <Input
                  readonly
                  value={data.person.customerId}
                  class="bg-primary-light dark:bg-secondary-dark"
                />
              </div>
            </div>
            <div class="grid grid-cols-2 gap-5">
              <div class="flex flex-col gap-2">
                <P
                  class="text-sm font-medium leading-tight text-secondary-light dark:text-secondary-dark"
                  >{$_("email")}</P
                >
                <Input
                  readonly
                  value={data.person.email}
                  class="bg-primary-light dark:bg-secondary-dark"
                />
              </div>
              <div class="flex flex-col gap-2">
                <P
                  class="text-sm font-medium leading-tight text-secondary-light dark:text-secondary-dark"
                  >{$_("phone")}</P
                >
                <Input
                  readonly
                  value={data.person.phone}
                  class="bg-primary-light dark:bg-secondary-dark"
                />
              </div>
            </div>
            <div class="flex flex-col grow gap-2">
              <P
                class="text-sm font-medium leading-tight text-secondary-light dark:text-secondary-dark"
              >
                {$_("address")}
              </P>
              <Input
                readonly
                value={data.person.address}
                class="bg-primary-light dark:bg-secondary-dark"
              />
            </div>
            <div class="grid grid-cols-2 gap-5">
              <div class="flex flex-col gap-2">
                <P
                  class="text-sm font-medium leading-tight text-secondary-light dark:text-secondary-dark"
                  >{$_("city")}</P
                >
                <Input
                  readonly
                  value={data.person.city}
                  class="bg-primary-light dark:bg-secondary-dark"
                />
              </div>
              <div class="flex flex-col gap-2">
                <P
                  class="text-sm font-medium leading-tight text-secondary-light dark:text-secondary-dark"
                  >{$_("state")}</P
                >
                <Input
                  readonly
                  value={data.person.state}
                  class="bg-primary-light dark:bg-secondary-dark"
                />
              </div>
            </div>
            <div class="grid grid-cols-2 gap-5">
              <div class="flex flex-col gap-2">
                <P
                  class="text-sm font-medium leading-tight text-secondary-light dark:text-secondary-dark"
                  >{$_("zipCode")}</P
                >
                <Input
                  readonly
                  value={data.person.zipCode}
                  class="bg-primary-light dark:bg-secondary-dark"
                />
              </div>
              <div class="flex flex-col gap-2">
                <P
                  class="text-sm font-medium leading-tight text-secondary-light dark:text-secondary-dark"
                  >{$_("country")}</P
                >
                <Input
                  readonly
                  value={data.person.country}
                  class="bg-primary-light dark:bg-secondary-dark"
                />
              </div>
            </div>
            <div class="grid grid-cols-2 gap-5">
              <div class="flex flex-col gap-2">
                <P
                  class="text-sm font-medium leading-tight text-secondary-light dark:text-secondary-dark"
                  >{$_("dateOfBirth")}</P
                >
                <Input
                  readonly
                  value={data.person.dateOfBirth}
                  class="bg-primary-light dark:bg-secondary-dark"
                />
              </div>
              <div class="flex flex-col gap-2">
                <P
                  class="text-sm font-medium leading-tight text-secondary-light dark:text-secondary-dark"
                  >{$_("active")}</P
                >
                <!-- <Input readonly value={data.person.active} /> -->
                <div class="flex grow gap-2">
                  <Button
                    color="light"
                    class="w-full flex justify-between focus-within:!ring-1 bg-primary-light dark:bg-secondary-dark hover:bg-transparent dark:hover:bg-secondary-dark"
                  >
                    {data.person.active}
                    <ChevronDownOutline class="w-3 h-3 ms-2 dark:text-white" />
                  </Button>
                  <!-- <Dropdown
                    placement="bottom-start"
                    bind:open={dropdownOpen}
                    class="w-48 overflow-y-auto py-1 dark:bg-secondary-dark"
                  >
                    <DropdownItem on:click={() => (dropdownOpen = false)}>true</DropdownItem>
                    <DropdownItem on:click={() => (dropdownOpen = false)}>false</DropdownItem>
                  </Dropdown> -->
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</div>
