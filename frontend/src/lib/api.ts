// const BASE_URL = 'http://localhost:8080';
const BASE_URL = import.meta.env.VITE_API_URL;

async function request<T>(url: string, options: RequestInit = {}): Promise<T> {
  const token = localStorage.getItem("token");
  const headers: HeadersInit = {
    "Content-Type": "application/json",
    ...(options.headers || {}),
  };
  if (token) {
    (headers as Record<string, string>)["Authorization"] = `Bearer ${token}`;
  }

  const res = await fetch(`${BASE_URL}${url}`, {
    ...options,
    headers,
  });

  const data = await res.json();
  if (!res.ok) {
    throw new Error(data.error || "Something went wrong");
  }
  return data as T;
}

export const api = {
  get<T>(url: string) {
    return request<T>(url);
  },
  post<T>(url: string, body: unknown) {
    return request<T>(url, { method: "POST", body: JSON.stringify(body) });
  },
  put<T>(url: string, body: unknown) {
    return request<T>(url, { method: "PUT", body: JSON.stringify(body) });
  },
  delete<T>(url: string) {
    return request<T>(url, { method: "DELETE" });
  },
};
