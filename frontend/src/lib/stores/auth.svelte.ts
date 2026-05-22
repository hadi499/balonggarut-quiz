// Token tidak lagi disimpan di client — dikelola sepenuhnya oleh HttpOnly cookie di browser.
// State login ditentukan dari data /me, bukan keberadaan token di localStorage.
//
// Menggunakan class dengan $state fields untuk reaktivitas Svelte 5 yang lebih eksplisit dan reliable.

class AuthStore {
  username = $state<string | null>(null);
  role = $state<string | null>(null);
  /**
   * True setelah init() selesai pertama kali.
   * Dipakai oleh api.ts untuk mencegah auto-logout saat init sedang berjalan.
   */
  initialized = $state(false);

  get isLoggedIn() {
    return !!this.username;
  }

  get isTeacher() {
    return this.role === 'teacher';
  }

  login(username: string, role: string) {
    this.username = username;
    this.role = role;
  }

  logout() {
    this.username = null;
    this.role = null;
  }

  /**
   * Dipanggil sekali di onMount layout. Mencoba restore sesi dari cookie yang masih valid.
   * Guard `initialized` mencegah race condition: jika init() belum selesai saat user login,
   * response /me yang gagal tidak akan menimpa state yang baru di-set oleh login().
   */
  async init() {
    if (typeof window === 'undefined') return; // SSR guard
    if (this.initialized) return; // Jangan jalankan dua kali

    try {
      const { api } = await import('$lib/api');
      // Panggil /api/auth/session yang selalu return 200 — tidak pernah throw 401
      const session = await api.get<{ loggedIn: boolean; username?: string; role?: string }>('/api/auth/session');
      if (session.loggedIn && session.username) {
        this.username = session.username;
        this.role = session.role ?? null;
      }
    } catch {
      // Catch network error dll — state tetap null (tidak login)
    } finally {
      this.initialized = true;
    }
  }
}

export const auth = new AuthStore();
