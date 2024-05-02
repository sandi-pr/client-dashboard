<script>
  import { onMount } from "svelte";
  import { Heading } from "flowbite-svelte";
  import DataTable from "$lib/components/DataTable.svelte";
  import { CR3Call } from "hexacomponents-svelte";

  const call = new CR3Call(globalThis.location.origin);

  let pagination = { page: 1, pageSize: 25, pageCount: 1, total: 25 };

  let datas = $state([]);

  let configs = $state({});

  async function getPerson() {
    let token = sessionStorage.getItem("token");
    if (token !== null) {
      let response = await call.call("client-dashboard", "get-persons", {
        URLQuery: {},
        URLParams: {},
        Params: {
          page: pagination.page,
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
        configs = {
          search: true,
          pagination: response.result.meta.pagination,
        };
        let promises2 = response.result.data.map(async (person, index) => {
          return {
            id: person.id,
            fullname: person.attributes.fullname,
            email: person.attributes.email,
            phone: person.attributes.phone,
            // address: person.attributes.address,
            // city: person.attributes.city,
            // state: person.attributes.state,
            // zip_code: person.attributes.zip_code,
            country: person.attributes.country,
            // date_of_birth: person.attributes.date_of_birth,
            is_active: person.attributes.is_active,
            actions: {
              objects: [
                {
                  href:
                    `${await call.getPage("client-dashboard", "crm-person-detail")}` +
                    "?personId=" +
                    person.id,
                  name: "View",
                },
              ],
            },
          };
        });

        datas = await Promise.all(promises2);
      }
    }
  }

  async function next() {
    let token = sessionStorage.getItem("token");
    if (token !== null) {
      try {
        let response = await call.call("client-dashboard", "get-persons", {
          URLQuery: {},
          URLParams: {},
          Params: {
            page: pagination.page + 1,
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
          configs = {
            search: true,
            pagination: response.result.meta.pagination,
          };
          let promises2 = response.result.data.map(async (person, index) => {
            return {
              id: person.id,
              fullname: person.attributes.fullname,
              email: person.attributes.email,
              phone: person.attributes.phone,
              // address: person.attributes.address,
              // city: person.attributes.city,
              // state: person.attributes.state,
              // zip_code: person.attributes.zip_code,
              country: person.attributes.country,
              // date_of_birth: person.attributes.date_of_birth,
              is_active: person.attributes.is_active,
              actions: {
                objects: [
                  {
                    href:
                      `${await call.getPage("client-dashboard", "crm-person-detail")}` +
                      "?personId=" +
                      person.id,
                    name: "View",
                  },
                ],
              },
            };
          });
          datas = await Promise.all(promises2);
        }
      } catch (error) {
        console.log(error);
      }
    }
  }

  async function previous() {
    let token = sessionStorage.getItem("token");
    if (token !== null) {
      try {
        let response = await call.call("client-dashboard", "get-persons", {
          URLQuery: {},
          URLParams: {},
          Params: {
            page: pagination.page - 1,
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
          configs = {
            search: true,
            pagination: response.result.meta.pagination,
          };
          let promises2 = response.result.data.map(async (person, index) => {
            return {
              id: person.id,
              fullname: person.attributes.fullname,
              email: person.attributes.email,
              phone: person.attributes.phone,
              // address: person.attributes.address,
              // city: person.attributes.city,
              // state: person.attributes.state,
              // zip_code: person.attributes.zip_code,
              country: person.attributes.country,
              // date_of_birth: person.attributes.date_of_birth,
              is_active: person.attributes.is_active,
              actions: {
                objects: [
                  {
                    href:
                      `${await call.getPage("client-dashboard", "crm-person-detail")}` +
                      "?personId=" +
                      person.id,
                    name: "View",
                  },
                ],
              },
            };
          });
          datas = await Promise.all(promises2);
        }
      } catch (error) {
        console.log(error);
      }
    }
  }

  onMount(async () => {
    getPerson();
  });
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
        <DataTable bind:datas bind:configs bind:next bind:previous />
      </div>
    </div>
  </div>
</div>
