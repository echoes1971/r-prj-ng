import i18n from './i18n';

/**
 * Extract and translate error message from API response
 * @param {Error} err - The error object from axios
 * @param {string} fallback - Fallback message if no error found
 * @returns {string} - Translated error message
 */
export const getErrorMessage = (err, fallback = "An error occurred") => {
  // Check if response has structured error (APIError format)
  if (err.response?.data?.code) {
    const { code, message, params } = err.response.data;
    
    // Try to get translation from errors namespace
    const translationKey = `errors:${code}`;
    const translated = i18n.t(translationKey, params || {});
    
    // If translation not found (returns key), use fallback message from backend
    if (translated === translationKey) {
      return message || fallback;
    }
    
    return translated;
  }
  
  // Legacy format: {error: "message"}
  if (err.response?.data?.error) {
    return err.response.data.error;
  }
  
  // Axios error message
  if (err.message) {
    return err.message;
  }
  
  return fallback;
};
