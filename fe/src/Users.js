import React, { useEffect, useState } from "react";
import axios from "axios";

function Users() {
  const [users, setUsers] = useState([]);

  useEffect(() => {
    const token = localStorage.getItem("token");
    axios.get("/users/316", {
      headers: { Authorization: `Bearer ${token}` }
    }).then(res =>
      {
        setUsers([res.data]);
        // alert(JSON.stringify(res.data));
      });
  }, []);

  return (
    <div className="container mt-3">
      <h2>Utenti</h2>
      <table className="table table-striped">
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
