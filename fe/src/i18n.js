import i18n from "i18next";
import LanguageDetector from "i18next-browser-languagedetector";
import { initReactI18next } from "react-i18next";
import translationEN from "./locales/en/translation.json";
import translationFR from "./locales/fr/translation.json";
import translationDE from "./locales/de/translation.json";
import translationIT from "./locales/it/translation.json";
import errorsEN from "./locales/en/errors.json";
import errorsFR from "./locales/fr/errors.json";
import errorsDE from "./locales/de/errors.json";
import errorsIT from "./locales/it/errors.json";

const resources = {
  en: { 
    translation: translationEN,
    errors: errorsEN 
  },
  fr: { 
    translation: translationFR,
    errors: errorsFR 
  },
  de: { 
    translation: translationDE,
    errors: errorsDE 
  },
  it: { 
    translation: translationIT,
    errors: errorsIT 
  }
};

const savedLang = localStorage.getItem("lang") || "en";

i18n
  .use(LanguageDetector)
  .use(initReactI18next)
  .init({
    resources,
    lng: savedLang || "en", // lingua di default
    fallbackLng: "en",      // fallback se manca una traduzione
    interpolation: { escapeValue: false },
    detection: {
      order: ["localStorage", "navigator"],
      caches: ["localStorage"]
    }
  });

export default i18n;
