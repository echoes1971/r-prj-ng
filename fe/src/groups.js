import React, { useContext, useEffect, useState } from "react";
import api from "./axios";
import { ThemeContext } from "./ThemeContext";
import { useTranslation } from "react-i18next";
import AssociationManager from "./AssociationManager";
import { getErrorMessage } from "./errorHandler";

function Groups() {
  const { t } = useTranslation();
  const [groups, setGroups] = useState([]);
  const [users, setUsers] = useState([]);
  const [query, setQuery] = useState("");
  const [editingGroup, setEditingGroup] = useState(null); // gruppo in modifica
  const [errorMessage, setErrorMessage] = useState("");
  const { dark, themeClass } = useContext(ThemeContext);

    // Carica tutti i gruppi all'inizio
  useEffect(() => {
    fetchGroups();
    fetchUsers();
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

  const fetchGroups = async (search = "") => {
    const token = localStorage.getItem("token");
    try {
      const res = await api.get(
        search
          ? `/groups?search=${encodeURIComponent(search)}`
          : "/groups",
        { headers: { Authorization: `Bearer ${token}` }, }
      );
      setGroups(res.data || []); // supponendo che l'API restituisca un array
    } catch (err) {
      console.log("Error loading groups.");
    }
  };

  const fetchUsers = async () => {
    const token = localStorage.getItem("token");
    try {
      const res = await api.get("/users", {
        headers: { Authorization: `Bearer ${token}` },
      });
      setUsers(res.data || []);
    } catch (err) {
      console.log("Error loading users.");
    }
  };

  const handleSearch = (e) => {
    e.preventDefault();
    fetchGroups(query);
  };

  const handleEditClick = async (group) => {
    // Carica i dettagli completi del gruppo inclusi i user_ids
    const token = localStorage.getItem("token");
    try {
      const res = await api.get(`/groups/${group.ID}`, {
        headers: { Authorization: `Bearer ${token}` },
      });
      setEditingGroup({ ...group, user_ids: res.data.user_ids || [] });
      setErrorMessage("");
    } catch (err) {
      console.log("Error loading group details.");
      // Fallback ai dati base se la chiamata fallisce
      setEditingGroup({ ...group, user_ids: [] });
      setErrorMessage("");
    }
  };

  const handleEditChange = (e) => {
    const { name, value } = e.target;
    // alert(name+"="+value);
    setEditingGroup((prev) => ({ ...prev, [name]: value }));
  };

  const handleSave = async () => {
    // Validazione
    if (!editingGroup.Name || editingGroup.Name.trim() === "") {
      setErrorMessage(t("groups.name_required") || "Group name is required");
      return;
    }
    
    setErrorMessage("");
    const token = localStorage.getItem("token");
    try {
      if (!editingGroup.ID) {
        // Nuovo gruppo
        await api.post(`/groups`, {
          name: editingGroup.Name,
          description: editingGroup.Description,
          user_ids: editingGroup.user_ids || []
        }, { headers: { Authorization: `Bearer ${token}` } } );
        setEditingGroup(null);
        fetchGroups();
        return;
      }
      // Gruppo esistente
      await api.put(`/groups/${editingGroup.ID}`, {
        name: editingGroup.Name,
        description: editingGroup.Description,
        user_ids: editingGroup.user_ids || []
      }, { headers: { Authorization: `Bearer ${token}` } } );
      setEditingGroup(null);
      fetchGroups();
    } catch (err) {
      // Extract and translate error message from response
      const errorMsg = getErrorMessage(err, t("groups.error_saving") || "Error saving group");
      setErrorMessage(errorMsg);
    }
  };

    const handleDelete = async () => {
    if (window.confirm("Are you sure to delete this group?")) {
        try {
            const token = localStorage.getItem("token");
            await api.delete(`/groups/${editingGroup.ID}`, {
                headers: { Authorization: `Bearer ${token}` },
            });
            setEditingGroup(null);
            fetchGroups();
        } catch (err) {
            // Extract and translate error message from response
            const errorMsg = getErrorMessage(err, t("groups.error_deleting") || "Error deleting group");
            setErrorMessage(errorMsg);
        }
    }
  };

  return (
    <div className={`container ${themeClass}`}>
      <h2 className={dark ? "text-light" : "text-dark"}>{t("groups.groups")}</h2>

      {/* Form di ricerca */}
      {!editingGroup && (
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

      {!editingGroup && (
        <button
          className="btn btn-success mb-3"
          onClick={() => {
            setEditingGroup({ ID: "", Name: "", Description: "", user_ids: []});
            setErrorMessage("");
          }}
        >
          {t("groups.newGroup")}
        </button>
      )}

      {!editingGroup && groups.length > 0 && (
        <table 
        className={`table ${dark ? "table-dark" : "table-striped"} table-hover`}
        >
          <thead>
            <tr>
              <th className="d-none d-md-table-cell">#</th>
              <th>Name</th>
              <th>Description</th>
              <th>Actions</th>
            </tr>
          </thead>
          <tbody>
            {groups.map((g, index) => (
              <tr key={g.ID}>
                <td className="d-none d-md-table-cell">{index+1}</td>
                <td>{g.Name}</td>
                <td>{g.Description}</td>
                <td>
                  <button
                    className="btn btn-sm btn-warning"
                    onClick={() => handleEditClick(g)}
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
      {editingGroup && (
        <div className={`card p-3 mt-3 ${dark ? "bg-dark text-light" : "bg-light text-dark" }`}>
          <h4>{editingGroup.ID>'' ? t("common.edit") : t("common.create")} {t("groups.group")}</h4>
          
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
            value={editingGroup.ID}
            readOnly
            onChange={handleEditChange}
          />
          <input
            className={`form-control mb-2 ${dark ? "bg-secondary text-light" : ""}`}
            name="Name"
            title="Name"
            value={editingGroup.Name}
            {...editingGroup.ID ? "disabled" : null}
            onChange={handleEditChange}
          />
          <input
            className={`form-control mb-2 ${dark ? "bg-secondary text-light" : ""}`}
            name="Description"
            title="Description"
            value={editingGroup.Description}
            onChange={handleEditChange}
          />
          
          {/* Association Manager per gli utenti */}
          {editingGroup.ID==='' ? '' :
          <AssociationManager
            title={t("groups.users") || "Users"}
            available={users}
            selected={editingGroup.user_ids || []}
            onChange={(newUserIds) => setEditingGroup(prev => ({ ...prev, user_ids: newUserIds }))}
            labelKey="Fullname"
            valueKey="ID"
          />
          }

          <div>
            <button className="btn btn-success me-2" onClick={handleSave}>
              { t("common.save") }
            </button>
            <button
              className="btn btn-secondary me-4"
              onClick={() => setEditingGroup(null)}
            >
              { t("common.cancel") }
            </button>
            {editingGroup.ID>"" ?
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

export default Groups;