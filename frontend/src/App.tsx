import {
  Admin,
  Resource,
  ListGuesser,
  EditGuesser,
  ShowGuesser,
} from "react-admin";
import simpleRestProvider from "ra-data-simple-rest";
//import { dataProvider } from "./dataProvider";
import { authProvider } from "./authProvider";

let dataProvider = simpleRestProvider("http://127.0.0.1:8080/api/v1");
// dataProvider.getOne = (resource, params) => {
//   return fetch(`http://127.0.0.1:8080/api/v1/keys/TEST_1::newkey`, {
//     method: "GET",
//     headers: new Headers({
//       "Content-Type": "application/json",
//       Authorization: "Bearer " + localStorage.getItem("token"),
//     }),
//   })
//     .then((response) => response.json())
//     .then((json) => {
//       return Promise.resolve({
//         data: json,
//       });
//     })
//     .catch((error) => {
//       console.error(error);
//     });
// };

export const App = () => (
  <Admin dataProvider={dataProvider} authProvider={authProvider}>
    {/* <Resource
      name="posts"
      list={ListGuesser}
      edit={EditGuesser}
      show={ShowGuesser}
    />
    <Resource
      name="comments"
      list={ListGuesser}
      edit={EditGuesser}
      show={ShowGuesser}
    /> */}
    <Resource
      name="keys"
      list={ListGuesser}
      edit={EditGuesser}
      show={ShowGuesser}
    />
    <Resource
      name="keys-group"
      list={ListGuesser}
      edit={EditGuesser}
      show={ShowGuesser}
    ></Resource>
  </Admin>
);
