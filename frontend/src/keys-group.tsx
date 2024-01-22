import { List, Datagrid, TextField, DateField, EditButton } from "react-admin";
import { TextInput } from "react-admin";

const postFilters = [<TextInput label="Search" source="key" alwaysOn />];

export const KeysGroupList = () => {
  return (
    <List filters={postFilters}>
      <Datagrid>
        <TextField source="id" />
        <TextField source="value" />
        <EditButton />
      </Datagrid>
    </List>
  );
};
