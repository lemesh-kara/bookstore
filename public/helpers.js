function decodeJWT(token) {
  const base64Url = token.split(".")[1];
  const base64 = base64Url.replace(/-/g, "+").replace(/_/g, "/");
  const jsonPayload = decodeURIComponent(
    atob(base64)
      .split("")
      .map((c) => {
        return "%" + ("00" + c.charCodeAt(0).toString(16)).slice(-2);
      })
      .join("")
  );
  return JSON.parse(jsonPayload);
}

function getUserIdFromJWT(token) {
  return decodeJWT(token).user_id;
}

function getCurrentUserId() {
  return getUserIdFromJWT(Cookies.get("access_token"));
}

export default {
  decodeJWT,
  getUserIdFromJWT,
  getCurrentUserId
};
