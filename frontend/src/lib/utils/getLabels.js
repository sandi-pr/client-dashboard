import triggerToast from "./toast";

export const getLabels = async (call) => {
  let token = localStorage.getItem("token");
  if (token !== null) {
    try {
      let response = await call.call("client-dashboard", "get-label", {
        URLQuery: {},
        URLParams: {},
        Params: {},
        Headers: {
          Authorization: "Bearer " + token,
        },
      });

      if (response.status !== 200) {
        if (response.status === 400) {
          console.error("Invalid parameters:", response);
          triggerToast($_("invalidParameters"), true, "error");
        } else if (response.status === 401) {
          console.error("Unauthorized access:", response);
          triggerToast($_("unauthorized"), true, "error");
        } else if (response.status === 500) {
          console.error("Server error:", response);
          triggerToast($_("serverError"), true, "error");
        } else {
          console.error("Unexpected status code:", response);
          triggerToast(`${$_("unexpectedError")}: ${response.status}`, true, "error");
        }
        throw new Error("Failed to fetch labels");
      }

      return response.result.labels;
    } catch (error) {
      console.error("Error during get labels:", error);
      triggerToast($_("failedGetLabels"), true, "error");
    }
  } else {
    triggerToast($_("unauthorized"), true, "error");
  }
  return [];
};
