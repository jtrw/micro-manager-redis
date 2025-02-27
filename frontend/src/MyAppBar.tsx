import { AppBar, Toolbar, Box, Select, MenuItem } from "@mui/material";
import { TitlePortal, UserMenu, Logout, RefreshIconButton } from "react-admin";
import SettingsIcon from "@mui/icons-material/Settings";
import { IconButton } from "@mui/material";
import React from "react";
import { BASE_URL, API_BASE } from "./common/constants.config";


const SettingsButton = () => (
  <IconButton color="inherit">
    <SettingsIcon />
  </IconButton>
);

const SelectDbButton = () => {
  const [db, setDb] = React.useState(() => localStorage.getItem('selectedDb') || "");
  const [dbValues, setDbValues] = React.useState([]);

  const { token } = JSON.parse(localStorage.getItem("user"));

  React.useEffect(() => {
    fetch(`${BASE_URL}${API_BASE}/databases`, {
      headers: {
        Authorization: `Bearer ${token}`,
      },
    })
      .then((response) => response.json())
      .then((data) => {
        const availableDbs = Object.entries(data)
          .filter(([_, value]) => value === 1)
          .map(([key]) => `db${key}`);
        
        setDbValues(availableDbs);
        // Встановлюємо значення тільки якщо немає збереженого або збережене значення відсутнє в доступних
        if (!db || !availableDbs.includes(db)) {
          const newDb = availableDbs[0] || "";
          setDb(newDb);
          localStorage.setItem('selectedDb', newDb);
        }
      })
      .catch((error) => console.error("Error fetching databases:", error));
  }, []);

  const handleChange = (event) => {
    const newDb = event.target.value;
    setDb(newDb);
    localStorage.setItem('selectedDb', newDb);
    window.location.reload();
    // try {
    //   const dbIndex = newDb.replace('db', '');
      
    //   fetch(`${BASE_URL}${API_BASE}/set-database`, {
    //     method: 'POST',
    //     headers: {
    //       'Authorization': `Bearer ${token}`,
    //       'Content-Type': 'application/json',
    //     },
    //     body: JSON.stringify({ database: dbIndex })
    //   }).then(response => {
    //     if (response.ok) {
    //       window.location.reload();
    //     } else {
    //       console.error('Failed to set database');
    //     }
    //   });
    // } catch (error) {
    //   console.error('Error setting database:', error);
    // }
  };

  return (
    <Box display="flex" alignItems="center">
      <Select value={db} onChange={handleChange} variant="standard">
        {dbValues.map((dbName) => (
          <MenuItem key={dbName} value={dbName}>
            {dbName}
          </MenuItem>
        ))}
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
