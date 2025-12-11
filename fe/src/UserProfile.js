import React, { useContext, useEffect, useState } from "react";
import { useParams, useNavigate } from "react-router-dom";
import { Container, Card, Form, Button, Alert } from "react-bootstrap";
import api from "./axios";
import { ThemeContext } from "./ThemeContext";
import { useTranslation } from "react-i18next";
import AssociationManager from "./AssociationManager";
import { getErrorMessage } from "./errorHandler";
import {GroupLinkView} from "./ContentWidgets";

function UserProfile() {
  const { t } = useTranslation();
  const { userId } = useParams();
  const navigate = useNavigate();
  const [user, setUser] = useState(null);
  const [groups, setGroups] = useState([]);
  const [confirmPwd, setConfirmPwd] = useState("");
  const [pwdError, setPwdError] = useState("");
  const [errorMessage, setErrorMessage] = useState("");
  const [successMessage, setSuccessMessage] = useState("");
  const { dark, themeClass } = useContext(ThemeContext);

  const currentUserId = localStorage.getItem("user_id");
  const userGroups = localStorage.getItem("groups") ? JSON.parse(localStorage.getItem("groups")) : [];
  const isAdmin = userGroups.includes("-2");
  const isOwnProfile = userId === currentUserId;

  useEffect(() => {
    // Check permissions: must be admin OR viewing own profile
    if (!isAdmin && !isOwnProfile) {
      navigate("/");
      return;
    }

    fetchUser();
    fetchGroups(); // Always fetch groups to display them
  }, [userId, isAdmin, isOwnProfile, navigate]);

  // Auto-dismiss messages after 5 seconds
  useEffect(() => {
    if (errorMessage) {
      const timer = setTimeout(() => setErrorMessage(""), 5000);
      return () => clearTimeout(timer);
    }
  }, [errorMessage]);

  useEffect(() => {
    if (successMessage) {
      const timer = setTimeout(() => setSuccessMessage(""), 5000);
      return () => clearTimeout(timer);
    }
  }, [successMessage]);

  const fetchUser = async () => {
    const token = localStorage.getItem("token");
    try {
      const res = await api.get(`/users/${userId}`, {
        headers: { Authorization: `Bearer ${token}` },
      });
      setUser({ ...res.data, pwd: "" });
    } catch (err) {
      setErrorMessage(getErrorMessage(err, t));
      console.error("Error loading user:", err);
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
      console.error("Error loading groups:", err);
    }
  };

  const handleChange = (e) => {
    const { name, value } = e.target;
    setUser((prev) => ({ ...prev, [name]: value }));
  };

  const handleSave = async () => {
    // Validation
    if (!user.login || user.login.trim() === "") {
      setPwdError(t("users.login_required"));
      return;
    }

    if (user.pwd && user.pwd !== confirmPwd) {
      setPwdError(t("users.password_mismatch"));
      return;
    }

    setPwdError("");
    setErrorMessage("");

    const token = localStorage.getItem("token");
    const payload = {
      login: user.login,
      fullname: user.fullname,
    };

    // Only include password if it was changed
    if (user.pwd && user.pwd.trim() !== "") {
      payload.pwd = user.pwd;
    }

    // Only admin can change groups
    if (isAdmin && user.group_ids) {
      payload.group_ids = user.group_ids;
    }

    try {
      await api.put(`/users/${userId}`, payload, {
        headers: { Authorization: `Bearer ${token}` },
      });
      setSuccessMessage(t("common.save") + " " + t("users.user") + " successful!");
      
      // If user changed their own username, update localStorage
      if (isOwnProfile && user.login !== localStorage.getItem("username")) {
        localStorage.setItem("username", user.login);
      }
      
      // Reload user data
      fetchUser();
      setConfirmPwd("");
    } catch (err) {
      setErrorMessage(getErrorMessage(err, t));
      console.error("Error saving user:", err);
    }
  };

  const handleCancel = () => {
    if (isAdmin) {
      navigate("/users");
    } else {
      navigate("/");
    }
  };

  const handleDelete = async () => {
    if (!window.confirm(t("users.delete_confirm"))) {
      return;
    }

    const token = localStorage.getItem("token");
    try {
      await api.delete(`/users/${userId}`, {
        headers: { Authorization: `Bearer ${token}` },
      });
      
      // If user deleted their own account, logout
      if (isOwnProfile) {
        localStorage.removeItem("token");
        localStorage.removeItem("username");
        localStorage.removeItem("user_id");
        localStorage.removeItem("groups");
        navigate("/login");
      } else {
        navigate("/users");
      }
    } catch (err) {
      setErrorMessage(getErrorMessage(err, t));
      console.error("Error deleting user:", err);
    }
  };

  if (!user) {
    return (
      <Container className={`mt-4 ${themeClass}`}>
        <p>Loading...</p>
      </Container>
    );
  }

  return (
    <Container className={`mt-4 ${themeClass}`}>
      <Card bg={dark ? "dark" : "light"} text={dark ? "light" : "dark"}>
        <Card.Header className={dark ? 'bg-secondary bg-opacity-25' : ''} style={dark ? { borderBottom: '1px solid rgba(255,255,255,0.1)' } : {}}>
          <h2>{isOwnProfile ? t("users.user_profile") : t("users.user") + ": " + user.login}</h2>
        </Card.Header>
        <Card.Body className={dark ? 'bg-secondary bg-opacity-25' : ''}>
          {errorMessage && <Alert variant="danger">{errorMessage}</Alert>}
          {successMessage && <Alert variant="success">{successMessage}</Alert>}

          <Form>
            <Form.Group className="mb-3">
              <Form.Label>{t("users.login")}</Form.Label>
              <Form.Control
                type="text"
                name="login"
                value={user.login || ""}
                onChange={handleChange}
                disabled={!isAdmin && !isOwnProfile}
              />
            </Form.Group>

            <Form.Group className="mb-3">
              <Form.Label>{t("users.fullname")}</Form.Label>
              <Form.Control
                type="text"
                name="fullname"
                value={user.fullname || ""}
                onChange={handleChange}
              />
            </Form.Group>

            <Form.Group className="mb-3">
              <Form.Label>{t("groups.group")}: </Form.Label>
              <GroupLinkView group_id={user.group_id} dark={dark} />
            </Form.Group>

            <Form.Group className="mb-3">
              <Form.Label>
                {t("users.password")} <small className="text-muted">({t("users.leave_blank")})</small>
              </Form.Label>
              <Form.Control
                type="password"
                name="pwd"
                value={user.pwd || ""}
                onChange={handleChange}
              />
            </Form.Group>

            <Form.Group className="mb-3">
              <Form.Label>{t("users.confirm_password")}</Form.Label>
              <Form.Control
                type="password"
                value={confirmPwd}
                onChange={(e) => setConfirmPwd(e.target.value)}
              />
              {pwdError && <small className="text-danger">{pwdError}</small>}
            </Form.Group>

            <Form.Group className="mb-3">
              <Form.Label>{t("users.groups")}</Form.Label>
              <AssociationManager
                title={t("users.groups")}
                available={groups}
                selected={user.group_ids || []}
                onChange={(newIds) =>
                  setUser((prev) => ({ ...prev, group_ids: newIds }))
                }
                labelKey="Name"
                valueKey="ID"
                disabled={!isAdmin}
                dark={dark}
              />
            </Form.Group>

            <div className="d-flex gap-2">
              <Button variant="primary" onClick={handleSave}>
                {t("common.save")}
              </Button>
              <Button variant="secondary" onClick={handleCancel}>
                {t("common.cancel")}
              </Button>
              {isAdmin && (
                <Button variant="danger" onClick={handleDelete}>
                  {t("common.delete")}
                </Button>
              )}
            </div>
          </Form>
        </Card.Body>
      </Card>
    </Container>
  );
}

export default UserProfile;
