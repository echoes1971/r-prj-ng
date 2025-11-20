import React, { useContext, useEffect, useState } from "react";
import api from "./axios";
import { ThemeContext } from "./ThemeContext";
import { useTranslation } from "react-i18next";
import AssociationManager from "./AssociationManager";
import { getErrorMessage } from "./errorHandler";

function Users() {
  const { t } = useTranslation();
  const [users, setUsers] = useState([]);
  const [groups, setGroups] = useState([]);
  const [query, setQuery] = useState("");
  const [editingUser, setEditingUser] = useState(null); // utente in modifica
  const [confirmPwd, setConfirmPwd] = useState("");
  const [pwdError, setPwdError] = useState("");
  const [errorMessage, setErrorMessage] = useState("");
  const { dark, themeClass } = useContext(ThemeContext);

    // Carica tutti gli utenti all'inizio
  useEffect(() => {
    fetchUsers();
    fetchGroups();
  }, []);

  // Auto-dismiss error message after 5 seconds
  useEffect(() => {
    if (errorMessage) {
      const timer = setTimeout(() => {
        setErrorMessage("");
      }, 5000);
      return () => clearTimeout(timer);
    }
  }, [errorMessage]);

  const fetchUsers = async (search = "") => {
    const token = localStorage.getItem("token");
    try {
      const res = await api.get(
        search
          ? `/users?search=${encodeURIComponent(search)}&order_by=login`
          : "/users?order_by=login",
        { headers: { Authorization: `Bearer ${token}` }, }
      );
      setUsers(res.data || []); // supponendo che l'API restituisca un array
    } catch (err) {
      console.log("Error loading users.");
    }
  };

  const fetchGroups = async () => {
    const token = localStorage.getItem("token");
    try {
      const res = await api.get("/groups", {
        headers: { Authorization: `Bearer ${token}` },
      });
      setGroups(res.data || []);
    } catch (err) {
      console.log("Error loading groups.");
    }
  };

  const handleSearch = (e) => {
    e.preventDefault();
    fetchUsers(query);
  };

  // Helper per ottenere il nome del gruppo dall'ID
  const getGroupName = (groupId) => {
    const group = groups.find(g => g.ID === groupId);
    return group ? group.Name : groupId;
  };

  const handleEditClick = async (user) => {
    // Carica i dettagli completi dell'utente inclusi i group_ids
    const token = localStorage.getItem("token");
    try {
      const res = await api.get(`/users/${user.ID}`, {
        headers: { Authorization: `Bearer ${token}` },
      });
      setEditingUser({ ...user, group_ids: res.data.group_ids || [], Pwd: "" });
      setConfirmPwd("");
      setPwdError("");
      setErrorMessage("");
    } catch (err) {
      console.log("Error loading user details.");
      // Fallback ai dati base se la chiamata fallisce
      setEditingUser({ ...user, group_ids: [], Pwd: "" });
      setConfirmPwd("");
      setPwdError("");
      setErrorMessage("");
    }
  };

  const handleEditChange = (e) => {
    const { name, value } = e.target;
    // alert(name+"="+value);
    setEditingUser((prev) => ({ ...prev, [name]: value }));
  };

  const handleSave = async () => {
    // Validazione
    if (!editingUser.Login || editingUser.Login.trim() === "") {
      setPwdError(t("users.login_required") || "Login is required");
      return;
    }
    
    if (!editingUser.ID && (!editingUser.Pwd || editingUser.Pwd.trim() === "")) {
      setPwdError(t("users.password_required") || "Password is required for new users");
      return;
    }
    
    if (editingUser.Pwd && editingUser.Pwd !== confirmPwd) {
      setPwdError(t("users.password_mismatch") || "Passwords do not match");
      return;
    }
    
    setPwdError("");
    const token = localStorage.getItem("token");
    try {
      if (!editingUser.ID) {
        // Nuovo utente
        await api.post(`/users`, {
          login: editingUser.Login,
          pwd: editingUser.Pwd,
          fullname: editingUser.Fullname,
          group_id: editingUser.GroupID,
          group_ids: editingUser.group_ids || []
        }, { headers: { Authorization: `Bearer ${token}` } } );
        setEditingUser(null);
        setConfirmPwd("");
        fetchUsers();
        fetchGroups();
        return;
      }
      // Utente esistente - invia password solo se modificata
      const updateData = {
        login: editingUser.Login,
        fullname: editingUser.Fullname,
        group_id: editingUser.GroupID,
        group_ids: editingUser.group_ids || []
      };
      if (editingUser.Pwd && editingUser.Pwd.trim() !== "") {
        updateData.pwd = editingUser.Pwd;
      }
      await api.put(`/users/${editingUser.ID}`, updateData, { headers: { Authorization: `Bearer ${token}` } } );
      setEditingUser(null);
      setConfirmPwd("");
      fetchUsers();
    } catch (err) {
      // Extract and translate error message from response
      const errorMsg = getErrorMessage(err, t("users.error_saving") || "Error saving user");
      setErrorMessage(errorMsg);
    }
  };

  const handleDelete = async () => {
    if (window.confirm("Are you sure to delete this user?")) {
      const token = localStorage.getItem("token");
      await api.delete(`/users/${editingUser.ID}`, {
        headers: { Authorization: `Bearer ${token}` },
      });
      setEditingUser(null);
      fetchUsers();
    }
  };

  return (
    <div className={`container mt-3 ${themeClass}`}>
      <h2 className={dark ? "text-light" : "text-dark"}>{t("users.users")}</h2>

      {/* Form di ricerca */}
      {!editingUser && (
        <form className="d-flex mb-3" onSubmit={handleSearch}>
          <input
            type="text"
            className={`form-control me-2 ${dark ? "bg-secondary text-light" : ""}`}
            placeholder={t("common.search")}
            value={query}
            onChange={(e) => setQuery(e.target.value)}
          />
          <button className="btn btn-primary">{t("common.search")}</button>
        </form>
      )}

      {!editingUser && (
        <button
          className="btn btn-success mb-3"
          onClick={() => {
            setEditingUser({ ID: "", Login: "", Fullname: "", GroupID: "", group_ids: [], Pwd: "" });
            setConfirmPwd("");
            setPwdError("");
            setErrorMessage("");
          }}
        >
          {t("users.newUser")}
        </button>
      )}

      {!editingUser && users.length > 0 && (
        <table 
        className={`table ${dark ? "table-dark" : "table-striped"} table-hover`}
        >
          <thead>
            <tr>
              <th className="d-none d-md-table-cell">#</th>
              <th>Login</th>
              <th>Fullname</th>
              <th className="d-none d-md-table-cell">Group</th>
              <th>Actions</th>
            </tr>
          </thead>
          <tbody>
            {users.map((u, index) => (
              <tr key={u.ID}>
                <td className="d-none d-md-table-cell">{index+1}</td>
                <td>{u.Login}</td>
                <td>{u.Fullname}</td>
                <td className="d-none d-md-table-cell">{getGroupName(u.GroupID)}</td>
                <td>
                  <button
                    className="btn btn-sm btn-warning"
                    onClick={() => handleEditClick(u)}
                  >
                    {t("common.edit")}
                  </button>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      )}

      {/* Form di editing */}
      {editingUser && (
        <div className={`card p-3 mt-3 ${dark ? "bg-dark text-light" : "bg-light text-dark" }`}>
          <h4>{editingUser.ID>'' ? t("common.edit") : t("common.create")} {t("users.user")}</h4>
          
          {/* Error message at the top */}
          {errorMessage && (
            <div className="alert alert-danger alert-dismissible fade show" role="alert">
              {errorMessage}
              <button type="button" className="btn-close" onClick={() => setErrorMessage("")} aria-label="Close"></button>
            </div>
          )}
          
          <input
            className={`form-control mb-2 ${dark ? "bg-secondary text-light" : ""}`}
            name="ID"
            title="ID"
            value={editingUser.ID}
            readOnly
            onChange={handleEditChange}
          />
          <label className="form-label">{t("users.login") || "Login"} *</label>
          <input
            className={`form-control mb-2 ${dark ? "bg-secondary text-light" : ""}`}
            name="Login"
            placeholder={t("users.login") || "Login"}
            value={editingUser.Login}
            disabled={editingUser.ID !== ""}
            onChange={handleEditChange}
            required
          />
          <label className="form-label">{t("users.fullname") || "Fullname"}</label>
          <input
            className={`form-control mb-2 ${dark ? "bg-secondary text-light" : ""}`}
            name="Fullname"
            placeholder={t("users.fullname") || "Fullname"}
            value={editingUser.Fullname}
            onChange={handleEditChange}
          />
          <label className="form-label">
            {t("users.password") || "Password"}
            {!editingUser.ID && " *"}
            {editingUser.ID && " (" + (t("users.leave_blank") || "leave blank to keep current") + ")"}
          </label>
          <input
            className={`form-control mb-2 ${dark ? "bg-secondary text-light" : ""}`}
            name="Pwd"
            type="password"
            placeholder={t("users.password") || "Password"}
            value={editingUser.Pwd || ""}
            onChange={handleEditChange}
            required={!editingUser.ID}
          />
          <label className="form-label">
            {t("users.confirm_password") || "Confirm Password"}
            {!editingUser.ID && " *"}
          </label>
          <input
            className={`form-control mb-2 ${dark ? "bg-secondary text-light" : ""}`}
            type="password"
            placeholder={t("users.confirm_password") || "Confirm Password"}
            value={confirmPwd}
            onChange={(e) => setConfirmPwd(e.target.value)}
            required={!editingUser.ID}
          />
          {pwdError && (
            <div className="alert alert-danger" role="alert">
              {pwdError}
            </div>
          )}
          <input
            className={`form-control mb-2 ${dark ? "bg-secondary text-light" : ""}`}
            name="GroupID"
            title="GroupID"
            value={getGroupName(editingUser.GroupID)}
            readOnly
            onChange={handleEditChange}
          />
          
          {/* Association Manager per i gruppi */}
          <AssociationManager
            title={t("users.groups") || "Groups"}
            available={groups}
            selected={editingUser.group_ids || []}
            onChange={(newGroupIds) => setEditingUser(prev => ({ ...prev, group_ids: newGroupIds }))}
            labelKey="Name"
            valueKey="ID"
          />

          <div>
            <button className="btn btn-success me-2" onClick={handleSave}>
              { t("common.save") }
            </button>
            <button
              className="btn btn-secondary me-4"
              onClick={() => setEditingUser(null)}
            >
              { t("common.cancel") }
            </button>
            {editingUser.ID>"" ?
                      <button
                        className="btn btn-danger"
                        onClick={handleDelete}
                      >
                        { t("common.delete") }
                      </button>
          : null}
          </div>
        </div>
      )}
    </div>
  );
}

export default Users;
