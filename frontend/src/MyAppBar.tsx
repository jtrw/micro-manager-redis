import { AppBar, Toolbar, Box, Select, MenuItem } from "@mui/material";
import { TitlePortal, UserMenu, Logout, RefreshIconButton } from "react-admin";
import SettingsIcon from "@mui/icons-material/Settings";
import { IconButton } from "@mui/material";
import React from "react";

const SettingsButton = () => (
  <IconButton color="inherit">
    <SettingsIcon />
  </IconButton>
);

const SelectDbButton = () => {
  const [db, setDb] = React.useState("db1");

  const handleChange = (event) => {
    setDb(event.target.value);
    // Тут можна додати логіку для обробки зміни вибору БД
  };

  return (
    <Box display="flex" alignItems="center">
      <Select value={db} onChange={handleChange} variant="standard">
        <MenuItem value="db1">db1</MenuItem>
        <MenuItem value="db2">db2</MenuItem>
        <MenuItem value="db3">db3</MenuItem>
      </Select>
    </Box>
  );
};

const MyUserMenu = () => (
  <Box display="flex" alignItems="center">
    <UserMenu>
      <Logout />
    </UserMenu>
  </Box>
);

export const MyAppBar = () => (
  <AppBar position="static">
    <Toolbar>
      <TitlePortal />
      <Box sx={{ flex: "1" }} />
      <SelectDbButton />
      <RefreshIconButton />
      <MyUserMenu>
        <Logout />
      </MyUserMenu>
    </Toolbar>
  </AppBar>
);
