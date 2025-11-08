import React, { useContext, useEffect, useState } from "react";
import axios from "axios";
import { ThemeContext } from "./ThemeContext";

function Users() {
  const [users, setUsers] = useState([]);
  const [query, setQuery] = useState("");
  const { dark, themeClass } = useContext(ThemeContext);

    // Carica tutti gli utenti all'inizio
  useEffect(() => {
    fetchUsers();
  }, []);

  const fetchUsers = async (search = "") => {
    const token = localStorage.getItem("token");
    try {
      const res = await axios.get(
        search
          ? `/users?search=${encodeURIComponent(search)}`
          : "/users",
        {
          headers: { Authorization: `Bearer ${token}` },
        }
      );
      setUsers(res.data); // supponendo che l'API restituisca un array
    } catch (err) {
      alert("Errore caricamento utenti");
    }
  };

  const handleSearch = (e) => {
    e.preventDefault();
    fetchUsers(query);
  };

  // useEffect(() => {
  //   const token = localStorage.getItem("token");
  //   axios.get("/users/316", {
  //     headers: { Authorization: `Bearer ${token}` }
  //   }).then(res =>
  //     {
  //       setUsers([res.data]);
  //       // alert(JSON.stringify(res.data));
  //     });
  // }, []);

  return (
    <div className={`container mt-3 ${themeClass}`}>
      <h2 className={dark ? "text-light" : "text-dark"}>Utenti</h2>

      {/* Form di ricerca */}
      <form className="d-flex mb-3" onSubmit={handleSearch}>
        <input
          type="text"
          className="form-control me-2"
          placeholder="Cerca utente..."
          value={query}
          onChange={(e) => setQuery(e.target.value)}
        />
        <button className="btn btn-primary">Cerca</button>
      </form>

      <table 
       className={`table ${dark ? "table-dark" : "table-striped"} table-hover`}
     >
        <thead>
          <tr><th>ID</th><th>Login</th><th>Fullname</th><th>Group</th></tr>
        </thead>
        <tbody>
          {users.map(u => (
            <tr key={u.ID}>
              <td>{u.ID}</td>
              <td>{u.Login}</td>
              <td>{u.Fullname}</td>
              <td>{u.GroupID}</td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}

export default Users;
