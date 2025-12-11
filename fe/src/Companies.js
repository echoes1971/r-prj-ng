import React, { useContext, useEffect, useState } from "react";
import { ThemeContext } from "./ThemeContext";
import { useTranslation } from "react-i18next";
import { ObjectSearch } from "./DBObject";


export function Companies() {
  const { t } = useTranslation();
  const { dark, themeClass } = useContext(ThemeContext);
  // const [query, setQuery] = useState("");
  // const [editingFolder, setEditingFolder] = useState(null); // folder in editing

  const searchClassname = "DBCompany";

  const searchColumns = [
    { name: t("dbobjects.name") || "Name", attribute: "name", type: "string" },
    { name: t("dbobjects.description") || "Description", attribute: "description", type: "string" },
    { name: t("dbobjects.parent") || "Parent", attribute: "father_id", type: "objectLink" },

    // { name: t("dbobjects.name") || "Name", attribute: "name2", type: "string" },
    // { name: t("dbobjects.name") || "Name", attribute: "name3", type: "string" },
    // { name: t("dbobjects.name") || "Name", attribute: "name4", type: "string" },
  ];

  const resultsColumns = [
    { name: t("dbobjects.created") || "Created", attribute: "creator", type: "userLink", hideOnSmall: true },
    { name: t("dbobjects.group") || "Group", attribute: "group_id", type: "groupLink", hideOnSmall: true },
    { name: t("dbobjects.parent") || "Parent", attribute: "father_id", type: "objectLink", hideOnSmall: true },
    { name: t("dbobjects.name") || "Name", attribute: "name", type: "string", hideOnSmall: false },
    { name: t("dbobjects.description") || "Description", attribute: "description", type: "string", hideOnSmall: true },
  ]
  return (
    <ObjectSearch searchClassname={searchClassname} searchColumns={searchColumns} resultsColumns={resultsColumns} dark={dark} themeClass={themeClass} />
    );
}
