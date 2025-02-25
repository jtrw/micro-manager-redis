import { Layout } from "react-admin";
import { MyAppBar } from "./MyAppBar";

export const MyLayout = ({ children }) => (
  <Layout appBar={MyAppBar}>{children}</Layout>
);
