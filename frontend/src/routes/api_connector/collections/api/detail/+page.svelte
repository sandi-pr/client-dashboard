<script>
  import { onMount } from "svelte";
  import { _ } from "svelte-i18n";
  import { format } from "date-fns";
  import { page } from "$app/stores";

  import { call, triggerToast, errorHandling } from "$lib";
  import { goto } from "$app/navigation";

  import FormApiConnector from "$lib/pages/FormApiConnector.svelte";
  import FormApiParams from "$lib/pages/FormApiParams.svelte";

  let apiID = $state(0);
  let isBody = $state(false);
  let onEdit = $state(false);
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

        // console.log("response", response);

        if (response?.status !== 200) {
          const error = new Error("Failed to fetch API details");
          error.response = response;
          error.status = response?.status;
          throw error;
        }

        data.apiName = response.result.attributes.name;
        data.description = response.result.attributes.description || "";
        data.urlLink = response.result.attributes.url || "";
        data.method = response.result.attributes.method || "GET";
        data.authorization = response.result.attributes.auth_type || "No Auth";
        data.username = response.result.attributes.username || "";
        data.password = response.result.attributes.password || "";
        data.token = response.result.attributes.token || "";
        data.createOn = format(new Date(response.result.attributes.createdAt), "dd-MM-yyyy");
        data.body_type = response.result.attributes.body_type
          ? response.result.attributes.body_type
          : "raw";
        paramsData = response.result.attributes.params ? response.result.attributes.params : [];
        bodyData = response.result.attributes.body ? response.result.attributes.body : [];
        // console.log("data", data);
        // console.log("paramsData", paramsData);
        // console.log("bodyData", bodyData);
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

        // window.location.href = link + "?" + queryString;
        goto(link + "?" + queryString);

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

        const keyOrder = ["id", "key", "type", "description", "actions"];
        const actionKeyOrder = ["name", "icon"];
        const sortedParamsData = paramsData.map((obj) => {
          const sortedObj = {};
          keyOrder.forEach((key) => {
            if (obj.hasOwnProperty(key)) {
              // If the key is 'actions' and it's an object, sort its keys too
              if (
                key === "actions" &&
                obj[key] &&
                typeof obj[key] === "object" &&
                !Array.isArray(obj[key])
              ) {
                const actionsObj = obj[key];
                // If actions.objects is an array, sort each object's keys
                if (Array.isArray(actionsObj.objects)) {
                  sortedObj[key] = {
                    objects: actionsObj.objects.map((actionItem) => {
                      const sortedActionItem = {};
                      actionKeyOrder.forEach((aKey) => {
                        if (actionItem.hasOwnProperty(aKey)) {
                          sortedActionItem[aKey] = actionItem[aKey];
                        }
                      });
                      // Add any extra keys not in actionKeyOrder at the end
                      Object.keys(actionItem).forEach((aKey) => {
                        if (!actionKeyOrder.includes(aKey)) {
                          sortedActionItem[aKey] = actionItem[aKey];
                        }
                      });
                      return sortedActionItem;
                    }),
                  };
                  // Add any extra keys in actionsObj (besides 'objects') at the end
                  Object.keys(actionsObj).forEach((aKey) => {
                    if (aKey !== "objects") {
                      sortedObj[key][aKey] = actionsObj[aKey];
                    }
                  });
                } else {
                  // If no objects array, just copy as is
                  sortedObj[key] = actionsObj;
                }
              } else {
                sortedObj[key] = obj[key];
              }
            }
          });
          // Add any extra keys not in keyOrder at the end
          Object.keys(obj).forEach((key) => {
            if (!keyOrder.includes(key)) {
              sortedObj[key] = obj[key];
            }
          });
          return sortedObj;
        });

        if (sortedParamsData.length > 0) {
          newData.params = JSON.stringify(sortedParamsData);
        } else {
          newData.params = "";
        }
        if (bodyData.length > 0) {
          newData.body = JSON.stringify(bodyData);
        } else {
          newData.body = "";
        }

        // console.log("paramsData", paramsData, typeof newData.params);
        // console.log("bodyData", bodyData);
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
    bind:onEdit
    {saveApi}
    {deleteApi}
  />
{:else if stateCreate === "formParams"}
  <FormApiParams bind:stateCreate bind:paramsData />
{:else if stateCreate === "formBody"}
  <FormApiParams bind:stateCreate bind:paramsData={bodyData} bind:data bind:isBody />
{/if}
