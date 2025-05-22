import triggerToast from "./toast";

export const getUser = async (call, userData) => {
  let token = localStorage.getItem("token");
  if (token !== null) {
    try {
      let response = await call.call("client-dashboard", "get-user", {
        URLQuery: {},
        URLParams: { id: userData.upload_by.toString() },
        Params: {},
        Headers: {
          Authorization: "Bearer " + token,
        },
      });

      if (response.status !== 200) {
        console.log(response);
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
        throw new Error("Failed to fetch user data");
      }

      return response?.result?.user?.username || "Unknown User";
    } catch (error) {
      console.error("Error during get user:", error);
      triggerToast($_("failedGetUser"), true, "error");
    }
  } else {
    triggerToast($_("unauthorized"), true, "error");
  }
  return {};
};
