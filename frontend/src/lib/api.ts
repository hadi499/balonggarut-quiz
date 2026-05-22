// const BASE_URL = 'http://localhost:8080';
const BASE_URL = import.meta.env.VITE_API_URL;

async function request<T>(url: string, options: RequestInit = {}): Promise<T> {
  const headers: HeadersInit = {
    "Content-Type": "application/json",
    ...(options.headers || {}),
  };

  const res = await fetch(`${BASE_URL}${url}`, {
    ...options,
    headers,
    credentials: "include", // Wajib agar browser mengirim HttpOnly cookie ke backend
  });

  const data = await res.json();
  if (!res.ok) {
    // Auto-logout hanya jika init() sudah selesai — mencegah race condition di mana
    // /me dari init() yang gagal (belum login) menimpa state yang baru di-set login()
    if (res.status === 401) {
      const { auth } = await import("$lib/stores/auth.svelte");
      if (auth.initialized) {
        auth.logout();
      }
    }
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
