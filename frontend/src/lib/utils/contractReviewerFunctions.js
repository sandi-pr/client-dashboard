export const getRules = async (call, pagination = { page: 1, pageSize: 10 }, searchTerm = "") => {
  let token = localStorage.getItem("token");
  if (token) {
    try {
      const response = await call.call("client-dashboard", "get-contract-rules", {
        URLQuery: {},
        URLParams: {},
        Params: {
          search: searchTerm,
          page: pagination.page,
          page_size: pagination.pageSize,
        },
        Headers: {
          Authorization: "Bearer " + token,
        },
      });

      if (response.status !== 200) {
        const error = new Error("Failed to fetch rules");
        error.response = response;
        throw error;
      }

      return response.result.data;
    } catch (error) {
      console.log("error", error);
      return [];
    }
  }
};

export const getAllRules = async (call) => {
  let pagination = { page: 1, pageSize: 1000 };
  let allRules = [];
  while (true) {
    const rules = await getRules(call, pagination);
    if (rules.length === 0) break;
    allRules = [...allRules, ...rules];
    pagination.page += 1;
  }

  return allRules;
};

export const getProjects = async (call, pagination = { page: 1, pageSize: 10 }) => {
  let token = localStorage.getItem("token");
  if (token) {
    try {
      const response = await call.call("client-dashboard", "get-contract-projects", {
        URLQuery: {},
        URLParams: {},
        Params: {
          page: pagination.page,
          page_size: pagination.pageSize,
        },
        Headers: {
          Authorization: "Bearer " + token,
        },
      });

      if (response.status !== 200) {
        const error = new Error("Failed to fetch projects");
        error.response = response; // Attach the response object
        throw error;
      }

      return response.result.data;
    } catch (error) {
      console.log("error", error);
      return [];
    }
  }
};
