const checkIsLoggedIn = () => {
  return JSON.parse(sessionStorage.getItem("user"));
};

export default checkIsLoggedIn;
