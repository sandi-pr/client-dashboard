const checkIsLoggedIn = () => {
  return JSON.parse(localStorage.getItem("user"));
};

export default checkIsLoggedIn;
