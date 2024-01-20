import { AuthProvider, HttpError } from "react-admin";
import data from "./users.json";
import { stringify } from "query-string";
import { fetchUtils } from "react-admin";
import { API_BASE, BASE_URL } from "./common/constants.config";

/**
 * This authProvider is only for test purposes. Don't use it in production.
 */
export const authProvider: AuthProvider = {
  login: ({ username, password }) => {
    const httpClient = (url, options = {}) => {
      return fetchUtils.fetchJson(url, options);
    };

    return httpClient(`${BASE_URL}/auth`, {
      method: "POST",
      body: JSON.stringify({ username: username, password: password }),
    }).then(({ json }) => {
      localStorage.setItem("user", JSON.stringify(json));
    });
  },
  logout: () => {
    localStorage.removeItem("user");
    return Promise.resolve();
  },
  checkError: (error) => {
    const status = error.status;
    if (status === 401 || status === 403) {
      localStorage.removeItem("user");
      return Promise.reject();
    }
    return Promise.resolve();
  },
  checkAuth: () => {
    return localStorage.getItem("user") ? Promise.resolve() : Promise.reject();
  },
  getPermissions: () => {
    return Promise.resolve(undefined);
  },
  getIdentity: () => {
    const persistedUser = localStorage.getItem("user");
    const user = persistedUser ? JSON.parse(persistedUser) : null;

    return Promise.resolve(user);
  },
};

export default authProvider;
