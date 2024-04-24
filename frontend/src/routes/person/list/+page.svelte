<script>
  // import { onMount } from "svelte";
  import { Heading } from "flowbite-svelte";
  import DataTable from "$lib/components/DataTable.svelte";

  // export let data;
  let personsLatest = [];
  let pagination = { page: 1, pageSize: 25, pageCount: 1, total: 25 };

  personsLatest = [
    {
      attributes: {
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
    },
    {
      attributes: {
        customerId: "34678291",
        fullName: "John Doe",
        dateOfBirth: "03/15/1985",
        email: "johndoe@example.com",
        phone: "+1 123456789",
        country: "USA",
        address: "123 Main St, Anytown USA",
        city: "Anytown",
        state: "AnyState",
        zipCode: "12345",
        createdAt: "01-12-2024 03:02",
        active: true,
      },
    },
    {
      attributes: {
        customerId: "56473829",
        fullName: "Alice Johnson",
        dateOfBirth: "06/22/1978",
        email: "alice@example.com",
        phone: "+61 987654321",
        country: "Australia",
        address: "123 Main St, Anytown AU",
        city: "Anytown",
        state: "AnyState",
        zipCode: "12345",
        createdAt: "01-12-2024 03:02",
        active: false,
      },
    },
    {
      attributes: {
        customerId: "97261835",
        fullName: "Mohammed Ahmed",
        dateOfBirth: "12/04/1990",
        email: "mohammed@example.com",
        phone: "+91 9876543210",
        country: "India",
        address: "123 Main St, Anytown IN",
        city: "Anytown",
        state: "AnyState",
        zipCode: "12345",
        createdAt: "01-12-2024 03:02",
        active: true,
      },
    },
  ];

  let datas = personsLatest.map((person, index) => {
    return {
      id: index + 1,
      customerId: person.attributes.customerId,
      fullName: person.attributes.fullName,
      email: person.attributes.email,
      phone: person.attributes.phone,
      country: person.attributes.country,
      actions: {
        objects: [
          {
            href: "/person/" + person.attributes.customerId,
            name: "Edit",
          },
        ],
      },
    };
  });

  let configs = {
    search: true,
    pagination: pagination,
  };
</script>

<div class="flex dark:bg-hc-graydark-950">
  <div class="w-full pb-10">
    <div class="flex flex-col m-4">
      <!-- <Navbar /> -->

      <div class="flex flex-col">
        <!-- User Information -->
        <Heading
          tag="h3"
          class="text-3xl leading-[38px] text-hc-graylight-900 dark:text-hc-graydark-50"
        >
          Persons List
        </Heading>

        <!-- Document Table Data -->
        <DataTable {datas} {configs} />
      </div>
    </div>
  </div>
</div>
