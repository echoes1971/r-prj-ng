import React, { useContext, useEffect, useState } from "react";
import api from "./axios";
import { ThemeContext } from "./ThemeContext";
import { useTranslation } from "react-i18next";

function Groups() {
  const { t } = useTranslation();
  const [groups, setGroups] = useState([]);
  const [query, setQuery] = useState("");
  const [editingGroup, setEditingGroup] = useState(null); // gruppo in modifica
  const { dark, themeClass } = useContext(ThemeContext);

    // Carica tutti i gruppi all'inizio
  useEffect(() => {
    fetchGroups();
  }, []);

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

  const handleSearch = (e) => {
    e.preventDefault();
    fetchGroups(query);
  };

  const handleEditClick = (group) => {
    // alert("Editing group="+group.Name);
    setEditingGroup({ ...group }); // copia per non modificare direttamente
  };

  const handleEditChange = (e) => {
    const { name, value } = e.target;
    // alert(name+"="+value);
    setEditingGroup((prev) => ({ ...prev, [name]: value }));
  };

  const handleSave = async () => {
    const token = localStorage.getItem("token");
    try {
      if (!editingGroup.ID) {
        // Nuovo gruppo
        await api.post(`/groups`, editingGroup, { headers: { Authorization: `Bearer ${token}` } } );
        setEditingGroup(null);
        fetchGroups();
        return;
      }
      // Gruppo esistente
      await api.put(`/groups/${editingGroup.ID}`, editingGroup, { headers: { Authorization: `Bearer ${token}` } } );
      setEditingGroup(null);
      fetchGroups();
    } catch (err) {
      alert("Error saving group.");
    }
  };

    const handleDelete = async () => {
    if (window.confirm("Are you sure to delete this group?")) {
        try {
            const token = localStorage.getItem("token");
            await api.delete(`/groups/${editingGroup.ID}`, {
                headers: { Authorization: `Bearer ${token}` },
            });
        } catch (err) {
            alert("Error deleting group.");
        }
        setEditingGroup(null);
        fetchGroups();
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
          onClick={() => setEditingGroup({ ID: "", Name: "", Description: ""})}
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
              <th className="d-none d-md-table-cell">ID</th>
              <th>Name</th>
              <th>Description</th>
              <th>Actions</th>
            </tr>
          </thead>
          <tbody>
            {groups.map(g => (
              <tr key={g.ID}>
                <td className="d-none d-md-table-cell">{g.ID}</td>
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