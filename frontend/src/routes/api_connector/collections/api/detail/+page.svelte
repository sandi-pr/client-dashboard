<script>
  import { onMount } from "svelte";
  import { _ } from "svelte-i18n";
  import { format } from "date-fns";
  import { page } from "$app/stores";

  import { call, triggerToast, errorHandling } from "$lib";

  import FormApiConnector from "$lib/pages/FormApiConnector.svelte";
  import FormApiParams from "$lib/pages/FormApiParams.svelte";

  let apiID = $state(0);
  let isBody = $state(false);
  let stateCreate = $state("formApi");
  let paramsData = $state([]);
  let bodyData = $state([]);
  let data = $state({
    apiName: "",
    description: "",
    urlLink: "",
    method: "GET",
    authorization: "No Auth",
    username: "",
    password: "",
    token: "",
    body_type: "raw",
  });

  const getApi = async () => {
    try {
      let token = localStorage.getItem("token");

      if (token !== null) {
        let response = await $call.call("client-dashboard", "get-api", {
          URLQuery: {},
          URLParams: {
            id: apiID,
          },
          Params: {},
          Headers: {
            Authorization: "Bearer " + token,
          },
        });

        // console.log(response);

        if (response?.status !== 200) {
          const error = new Error("Failed to fetch API details");
          error.response = response;
          error.status = response?.status;
          throw error;
        }

        data.apiName = response.result.attributes.name;
        data.description = response.result.attributes.description;
        data.urlLink = response.result.attributes.url;
        data.method = response.result.attributes.method;
        data.authorization = response.result.attributes.auth_type;
        data.username = response.result.attributes.username;
        data.password = response.result.attributes.password;
        data.token = response.result.attributes.token;
        data.createOn = format(new Date(response.result.attributes.createdAt), "dd-MM-yyyy");
        data.body_type = response.result.attributes.body_type
          ? response.result.attributes.body_type
          : "raw";
        paramsData = response.result.attributes.params ? response.result.attributes.params : [];
        bodyData = response.result.attributes.body ? response.result.attributes.body : [];
      }
    } catch (error) {
      errorHandling($_, error, "get");
    }
  };

  const deleteApi = async () => {
    try {
      let token = localStorage.getItem("token");

      if (token !== null) {
        let response = await $call.call("client-dashboard", "delete-api", {
          URLQuery: {},
          URLParams: {
            id: apiID,
          },
          Params: {},
          Headers: {
            Authorization: "Bearer " + token,
          },
        });

        if (response?.status !== 200) {
          const error = new Error("Failed to delete API");
          error.response = response;
          error.status = response?.status;
          throw error;
        }

        let link = await $call.getPage("client-dashboard", "api-connector-collections-api");

        const params = new URLSearchParams(window.location.search);

        params.delete("apiId");

        const queryParams = {};
        for (const [key, value] of params.entries()) {
          queryParams[key] = value;
        }
        const queryString = new URLSearchParams(queryParams).toString();

        window.location.href = link + "?" + queryString;

        triggerToast($_("successfullyDeleted"));
      }
    } catch (error) {
      errorHandling($_, error, "delete");
    }
  };

  const saveApi = async () => {
    try {
      let token = localStorage.getItem("token");

      if (token !== null) {
        let newData = {
          name: data.apiName,
          description: data.description,
          url: data.urlLink,
          method: data.method,
          auth_type: data.authorization,
          username: data.username,
          password: data.password,
          token: data.token,
          body_type: data.body_type,
        };

        let collectionId = $page.url.searchParams.get("collectionId");
        newData.parent_id = parseInt(collectionId);
        if (paramsData.length > 0) {
          newData.params = JSON.stringify(paramsData);
        }
        if (bodyData.length > 0) {
          newData.body = JSON.stringify(bodyData);
        }
        // console.log("data", newData);

        let response = await $call.call("client-dashboard", "update-api", {
          URLQuery: {},
          URLParams: {
            id: apiID,
          },
          Params: newData,
          Headers: {
            Authorization: "Bearer " + token,
          },
        });

        if (response.status !== 200) {
          const error = new Error("Failed to update API");
          error.response = response;
          error.status = response?.status;
          throw error;
        }

        // console.log(response);
        triggerToast($_("successfullyUpdated"));
      }
    } catch (error) {
      errorHandling($_, error, "update");
    }
  };

  onMount(() => {
    apiID = $page.url.searchParams.get("apiId");
    getApi();
  });
</script>

{#if stateCreate === "formApi"}
  <FormApiConnector
    bind:data
    bind:stateCreate
    bind:paramsData
    bind:bodyData
    bind:apiID
    bind:isBody
    {saveApi}
    {deleteApi}
  />
{:else if stateCreate === "formParams"}
  <FormApiParams bind:stateCreate bind:paramsData />
{:else if stateCreate === "formBody"}
  <FormApiParams bind:stateCreate bind:paramsData={bodyData} bind:isBody />
{/if}
