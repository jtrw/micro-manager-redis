import { List, Datagrid, TextField, EditButton } from "react-admin";
import { TextInput } from "react-admin";

const postFilters = [<TextInput label="ID" source="key" alwaysOn />];

export const KeysList = () => {
  return (
    <List filters={postFilters}>
      <Datagrid>
        <TextField source="id" />
        <TextField
          component="pre"
          source="value"
          sx={{
            maxWidth: "30em",
            textOverflow: "ellipsis",
            whiteSpace: "nowrap",
            overflow: "hidden",
          }}
        />
        <TextField source="expire" />
        <EditButton />
      </Datagrid>
    </List>
  );
};
