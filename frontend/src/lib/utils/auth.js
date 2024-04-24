export function setUserDataSession(data) {
  const user = sessionStorage.setItem("user", JSON.stringify(data));
}

export function userDataSession() {
  const user = sessionStorage.getItem("user");
  return JSON.parse(user);
}

export function setOrganizationSession(data) {
  const user = sessionStorage.setItem("organization", JSON.stringify(data));
}

export function organizationSession() {
  const user = sessionStorage.getItem("organization");
  return JSON.parse(user);
}