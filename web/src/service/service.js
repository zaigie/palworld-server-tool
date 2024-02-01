import { useFetch } from "@vueuse/core";
import router from "@/router";

class Service {
  /**
   * Fetches data from a specified URL.
   *
   * @param {string} url - The URL to fetch data from.
   * @return {Promise<Response>} A Promise that resolves to the response from the server.
   */
  fetch(url) {
    return useFetch(`${url}`, {
      updateDataOnError: true,
      beforeFetch({ options }) {
        const token = localStorage.getItem("palworld_token");
        options.headers = {
          Authorization: `Bearer ${token}`,
          "Content-Type": "application/json",
          ...options.headers,
          "Remote-Ip-Address": localStorage.getItem("ip") || "127.0.0.1",
        };
        return {
          options,
        };
      },
      onFetchError(context) {
        if (context.response.status === 401) {
          localStorage.removeItem("palworld_token");
          return context;
        }
        return context;
      },
    });
  }

  /**
   * Generates a query string from a given credential object.
   *
   * @param {Object} credential - The credential object.
   * @return {string} - The generated query string.
   */
  generateQuery(credential) {
    const entries = Object.entries(credential);
    return entries
      .reduce((accumulation, [key, value]) => {
        if (value) {
          accumulation.push(`${key}=${value}`);
        }
        return accumulation;
      }, [])
      .join("&");
  }
}

export default Service;
