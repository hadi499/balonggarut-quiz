let _token = $state<string | null>(null);
let _username = $state<string | null>(null);
let _role = $state<string | null>(null);

export const auth = {
  get token() {
    return _token;
  },
  get isLoggedIn() {
    return !!_token;
  },
  get username() {
    return _username;
  },
  get role() {
    return _role;
  },
  get isTeacher() {
    return _role === 'teacher';
  },
  login(token: string, username: string, role: string) {
    _token = token;
    _username = username;
    _role = role;
    localStorage.setItem('token', token);
    localStorage.setItem('username', username);
    localStorage.setItem('role', role);
  },
  logout() {
    _token = null;
    _username = null;
    _role = null;
    localStorage.removeItem('token');
    localStorage.removeItem('username');
    localStorage.removeItem('role');
  },
  init() {
    const t =
      typeof localStorage !== 'undefined'
        ? localStorage.getItem('token')
        : null;
    if (t) {
      _token = t;
      _username = localStorage.getItem('username');
      _role = localStorage.getItem('role');
    }
  },
};
