import {
  fetchUtils,
  Admin,
  Resource,
  ListGuesser,
  EditGuesser,
  ShowGuesser,
} from "react-admin";
import simpleRestProvider from "ra-data-simple-rest";
import { authProvider } from "./authProvider";
import { BASE_URL, API_BASE } from "./common/constants.config";
import { KeysList } from "./keys";
import { KeysGroupList } from "./keys-group";

const httpClient = (url, options = {}) => {
  if (!options.headers) {
    options.headers = new Headers({ Accept: "application/json" });
  }

  const { token } = JSON.parse(localStorage.getItem("user"));
  options.headers.set("Authorization", `Bearer ${token}`);
  return fetchUtils.fetchJson(url, options);
};

let dataProvider = simpleRestProvider(`${BASE_URL}${API_BASE}`, httpClient);

export const App = () => (
  <Admin dataProvider={dataProvider} authProvider={authProvider}>
    <Resource
      name="keys"
      list={KeysList}
      edit={EditGuesser}
      show={ShowGuesser}
    />
    <Resource
      name="keys-group"
      list={KeysGroupList}
      edit={EditGuesser}
      show={ShowGuesser}
    ></Resource>
  </Admin>
);
