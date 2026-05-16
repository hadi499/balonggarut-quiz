<script lang="ts">
  import { page } from "$app/state";
  import { api } from "$lib/api";
  import { auth } from "$lib/stores/auth.svelte";
  import { goto } from "$app/navigation";
  import { onMount } from "svelte";

  $effect(() => {
    if (!auth.isLoggedIn) {
      goto("/login");
    } else if (!auth.isTeacher) {
      goto("/dashboard");
    }
  });

  interface Quiz {
    id: number;
    title: string;
    category: string;
    timeLimit: number;
    questions: { id: number }[];
  }

  interface Question {
    id: number;
    quiz_id: number;
    question: string;
    options: string[];
    answer: number;
  }

  interface ScoreEntry {
    id: number;
    username: string;
    score: number;
    quiz_title: string;
    created_at: string;
  }

  interface User {
    id: number;
    username: string;
    role: string;
  }

  interface ActivityLog {
    id: number;
    username: string;
    action: string;
    timestamp: string;
  }

  let quizzes = $state<Quiz[]>([]);
  let questions = $state<Question[]>([]);
  let users = $state<User[]>([]);
  let logs = $state<ActivityLog[]>([]);
  let searchLog = $state("");
  let filteredLogs = $derived(
    logs.filter(
      (log) =>
        log.username.toLowerCase().includes(searchLog.toLowerCase()) ||
        log.action.toLowerCase().includes(searchLog.toLowerCase()),
    ),
  );
  let searchQuestion = $state("");
  let filteredQuestions = $derived(
    questions.filter(
      (q) =>
        q.question.toLowerCase().includes(searchQuestion.toLowerCase()) ||
        getQuizTitle(q.quiz_id)
          .toLowerCase()
          .includes(searchQuestion.toLowerCase()),
    ),
  );
  let scores = $state<ScoreEntry[]>([]);
  let loading = $state(true);
  let refreshing = $state(false);

  let activeTab = $state<"quizzes" | "questions" | "scores" | "users" | "logs">(
    "quizzes",
  );

  let quizForm = $state({ title: "", category: "", timeLimit: 60 });
  let editingQuizId = $state<number | null>(null);
  let quizFormError = $state("");

  let questionForm = $state({
    quiz_id: 0,
    question: "",
    optionA: "",
    optionB: "",
    optionC: "",
    optionD: "",
    answer: 0,
  });
  let editingQuestionId = $state<number | null>(null);
  let questionFormError = $state("");

  let resetConfirm = $state(false);
  let deleteQuizId = $state<number | null>(null);
  let deleteUserId = $state<number | null>(null);

  async function deleteUser() {
    if (!deleteUserId) return;
    try {
      await api.delete(`/api/admin/users/${deleteUserId}`);
      deleteUserId = null;
      await loadAll();
    } catch {}
  }

  async function changeUserRole(id: number, newRole: string) {
    try {
      await api.put(`/api/admin/users/${id}/role`, { role: newRole });
      await loadAll();
    } catch {}
  }

  onMount(async () => {
    const tab = page.url.searchParams.get("tab");
    if (
      tab === "questions" ||
      tab === "scores" ||
      tab === "users" ||
      tab === "logs"
    ) {
      activeTab = tab as any;
    }
    await loadAll();
  });

  function switchTab(
    tab: "quizzes" | "questions" | "scores" | "users" | "logs",
  ) {
    activeTab = tab;
    const url = new URL(page.url);
    url.searchParams.set("tab", tab);
    goto(url, { replaceState: true, noScroll: true, keepFocus: true });
  }

  async function loadAll() {
    refreshing = true;
    try {
      const results = await Promise.allSettled([
        api.get<Quiz[]>("/api/quizzes"),
        api.get<Question[]>("/api/questions"),
        api.get<ScoreEntry[]>("/api/admin/scores"),
        api.get<User[]>("/api/admin/users"),
        api.get<ActivityLog[]>("/api/admin/logs"),
      ]);
      if (results[0].status === "fulfilled") quizzes = results[0].value;
      if (results[1].status === "fulfilled") questions = results[1].value;
      if (results[2].status === "fulfilled") scores = results[2].value || [];
      if (results[3].status === "fulfilled") users = results[3].value || [];
      if (results[4].status === "fulfilled") logs = results[4].value || [];
    } finally {
      loading = false;
      refreshing = false;
    }
  }

  async function saveQuiz() {
    quizFormError = "";
    const body = {
      title: quizForm.title,
      category: quizForm.category,
      timeLimit: quizForm.timeLimit,
    };
    try {
      if (editingQuizId) {
        await api.put(`/api/quizzes/${editingQuizId}`, body);
      } else {
        await api.post("/api/quizzes", body);
      }
      resetQuizForm();
      await loadAll();
    } catch (err) {
      quizFormError = (err as Error).message;
    }
  }

  function editQuiz(quiz: Quiz) {
    editingQuizId = quiz.id;
    quizForm = {
      title: quiz.title,
      category: quiz.category,
      timeLimit: quiz.timeLimit,
    };
    activeTab = "quizzes";
  }

  function resetQuizForm() {
    editingQuizId = null;
    quizForm = { title: "", category: "", timeLimit: 60 };
    quizFormError = "";
  }

  async function deleteQuiz() {
    if (!deleteQuizId) return;
    try {
      await api.delete(`/api/quizzes/${deleteQuizId}`);
      deleteQuizId = null;
      await loadAll();
    } catch {}
  }

  async function saveQuestion() {
    questionFormError = "";
    const body = {
      quiz_id: questionForm.quiz_id,
      question: questionForm.question,
      options: [
        questionForm.optionA,
        questionForm.optionB,
        questionForm.optionC,
        questionForm.optionD,
      ],
      answer: questionForm.answer,
    };
    try {
      if (editingQuestionId) {
        await api.put(`/api/questions/${editingQuestionId}`, body);
      } else {
        await api.post("/api/questions", body);
      }
      resetQuestionForm();
      await loadAll();
    } catch (err) {
      questionFormError = (err as Error).message;
    }
  }

  function editQuestion(q: Question) {
    editingQuestionId = q.id;
    questionForm = {
      quiz_id: q.quiz_id,
      question: q.question,
      optionA: q.options[0] || "",
      optionB: q.options[1] || "",
      optionC: q.options[2] || "",
      optionD: q.options[3] || "",
      answer: q.answer,
    };
    activeTab = "questions";
  }

  function resetQuestionForm() {
    editingQuestionId = null;
    questionForm = {
      quiz_id: 0,
      question: "",
      optionA: "",
      optionB: "",
      optionC: "",
      optionD: "",
      answer: 0,
    };
    questionFormError = "";
  }

  async function deleteQuestion(id: number) {
    try {
      await api.delete(`/api/questions/${id}`);
      await loadAll();
    } catch {}
  }

  async function resetAllScores() {
    try {
      await api.delete("/api/admin/scores/reset");
      resetConfirm = false;
      await loadAll();
    } catch {}
  }

  function getQuizTitle(id: number): string {
    return quizzes.find((q) => q.id === id)?.title ?? `Quiz #${id}`;
  }

  const inputClass =
    "w-full rounded-lg border border-slate-600 bg-slate-800 px-3 py-2 text-sm text-white placeholder-slate-500 focus:border-indigo-500 focus:outline-none focus:ring-1 focus:ring-indigo-500";
</script>

<svelte:head>
  <title>Admin Panel - Balonggarut Quiz</title>
</svelte:head>

{#if !auth.isLoggedIn}
  <p class="text-slate-400">Redirecting...</p>
{:else if !auth.isTeacher}
  <p class="text-slate-400">Redirecting...</p>
{:else if loading}
  <div class="flex justify-center py-12">
    <div
      class="h-8 w-8 animate-spin rounded-full border-2 border-indigo-500 border-t-transparent"
    ></div>
  </div>
{:else}
  <div class="space-y-6">
    <h1
      class="text-xl border-0 py-1 px-2 font-bold text-slate-200 w-fit text-center bg-indigo-500 rounded-md te"
    >
      Admin Panel
    </h1>

    <div class="flex gap-2 border-b border-slate-700">
      <button
        onclick={() => switchTab("quizzes")}
        class="px-4 py-2 text-md transition-colors {activeTab === 'quizzes'
          ? 'border-b-2 border-indigo-500 font-medium text-indigo-400'
          : 'text-slate-400 hover:text-slate-200'}"
      >
        Kuis
      </button>
      <button
        onclick={() => switchTab("questions")}
        class="px-4 py-2 text-md transition-colors {activeTab === 'questions'
          ? 'border-b-2 border-indigo-500 font-medium text-indigo-400'
          : 'text-slate-400 hover:text-slate-200'}"
      >
        Soal
      </button>
      <button
        onclick={() => switchTab("scores")}
        class="px-4 py-2 text-md transition-colors {activeTab === 'scores'
          ? 'border-b-2 border-indigo-500 font-medium text-indigo-400'
          : 'text-slate-400 hover:text-slate-200'}"
      >
        Nilai
      </button>
      <button
        onclick={() => switchTab("users")}
        class="px-4 py-2 text-md transition-colors {activeTab === 'users'
          ? 'border-b-2 border-indigo-500 font-medium text-indigo-400'
          : 'text-slate-400 hover:text-slate-200'}"
      >
        Pengguna
      </button>
      <button
        onclick={() => switchTab("logs")}
        class="px-4 py-2 text-md transition-colors {activeTab === 'logs'
          ? 'border-b-2 border-indigo-500 font-medium text-indigo-400'
          : 'text-slate-400 hover:text-slate-200'}"
      >
        Log Aktivitas
      </button>
    </div>

    {#if activeTab === "quizzes"}
      <div class="grid gap-6 lg:grid-cols-2 items-start">
        <div
          class="h-fit rounded-xl border border-slate-700 bg-slate-800/50 backdrop-blur-sm p-5 lg:sticky lg:top-20"
        >
          <h2 class="mb-4 text-lg font-semibold text-white">
            {editingQuizId ? "Edit Kuis" : "Buat Kuis Baru"}
          </h2>
          {#if quizFormError}
            <p class="mb-3 text-sm text-red-400">{quizFormError}</p>
          {/if}
          <form
            onsubmit={(e) => {
              e.preventDefault();
              saveQuiz();
            }}
            class="space-y-3"
            autocomplete="off"
          >
            <input
              placeholder="Judul"
              bind:value={quizForm.title}
              required
              class={inputClass}
            />
            <input
              placeholder="Kategori"
              bind:value={quizForm.category}
              required
              class={inputClass}
            />
            <input
              type="number"
              placeholder="Waktu per soal (detik)"
              bind:value={quizForm.timeLimit}
              required
              min="1"
              class={inputClass}
            />
            <div class="flex gap-2">
              <button
                type="submit"
                class="rounded-lg bg-indigo-500 px-4 py-2 text-sm text-white hover:bg-indigo-400 transition"
              >
                {editingQuizId ? "Update" : "Buat"}
              </button>
              {#if editingQuizId}
                <button
                  type="button"
                  onclick={resetQuizForm}
                  class="rounded-lg border border-slate-600 px-4 py-2 text-sm text-slate-300 hover:bg-slate-700 transition"
                  >Batal</button
                >
              {/if}
            </div>
          </form>
        </div>

        <div
          class="flex flex-col h-[600px] rounded-xl border border-slate-700 bg-slate-800/50 backdrop-blur-sm p-5"
        >
          <h2 class="flex-none mb-4 text-lg font-semibold text-white">
            Daftar Kuis
          </h2>
          {#if quizzes.length === 0}
            <p class="text-sm text-slate-500">Belum ada kuis.</p>
          {:else}
            <div class="space-y-2 overflow-y-auto pr-2 flex-1">
              {#each quizzes as quiz}
                <div
                  class="flex items-center justify-between rounded-lg border border-slate-700 px-3 py-2"
                >
                  <div>
                    <span class="text-md font-medium text-slate-200"
                      >{quiz.title}</span
                    >
                    <span class="ml-2 text-sm text-slate-500"
                      >{quiz.questions?.length ?? 0} soal</span
                    >
                  </div>
                  <div class="flex gap-1">
                    <button
                      onclick={() => editQuiz(quiz)}
                      class="rounded px-2 py-1 text-sm text-indigo-400 hover:bg-indigo-500/10 transition"
                      >Edit</button
                    >
                    <button
                      onclick={() => (deleteQuizId = quiz.id)}
                      class="rounded px-2 py-1 text-sm text-red-400 hover:bg-red-500/10 transition"
                      >Hapus</button
                    >
                  </div>
                </div>
              {/each}
            </div>
          {/if}
        </div>
      </div>
    {:else if activeTab === "questions"}
      <div class="grid gap-6 lg:grid-cols-2 items-start">
        <div
          class="h-fit rounded-xl border border-slate-700 bg-slate-800/50 backdrop-blur-sm p-5 lg:sticky lg:top-20"
        >
          <h2 class="mb-4 text-lg font-semibold text-white">
            {editingQuestionId ? "Edit Soal" : "Buat Soal Baru"}
          </h2>
          {#if questionFormError}
            <p class="mb-3 text-sm text-red-400">{questionFormError}</p>
          {/if}
          <form
            onsubmit={(e) => {
              e.preventDefault();
              saveQuestion();
            }}
            class="space-y-3"
            autocomplete="off"
          >
            <select bind:value={questionForm.quiz_id} class={inputClass}>
              <option value={0} disabled>Pilih Kuis</option>
              {#each quizzes as quiz}
                <option value={quiz.id}>{quiz.title}</option>
              {/each}
            </select>
            <input
              placeholder="Pertanyaan"
              bind:value={questionForm.question}
              required
              class={inputClass}
            />
            <input
              placeholder="A. Opsi jawaban"
              bind:value={questionForm.optionA}
              required
              class={inputClass}
            />
            <input
              placeholder="B. Opsi jawaban"
              bind:value={questionForm.optionB}
              required
              class={inputClass}
            />
            <input
              placeholder="C. Opsi jawaban"
              bind:value={questionForm.optionC}
              required
              class={inputClass}
            />
            <input
              placeholder="D. Opsi jawaban"
              bind:value={questionForm.optionD}
              required
              class={inputClass}
            />
            <select bind:value={questionForm.answer} class={inputClass}>
              <option value={0}>Jawaban: A</option>
              <option value={1}>Jawaban: B</option>
              <option value={2}>Jawaban: C</option>
              <option value={3}>Jawaban: D</option>
            </select>
            <div class="flex gap-2">
              <button
                type="submit"
                class="rounded-lg bg-indigo-500 px-4 py-2 text-sm text-white hover:bg-indigo-400 transition"
              >
                {editingQuestionId ? "Update" : "Buat"}
              </button>
              {#if editingQuestionId}
                <button
                  type="button"
                  onclick={resetQuestionForm}
                  class="rounded-lg border border-slate-600 px-4 py-2 text-sm text-slate-300 hover:bg-slate-700 transition"
                  >Batal</button
                >
              {/if}
            </div>
          </form>
        </div>

        <div
          class="flex flex-col h-[600px] rounded-xl border border-slate-700 bg-slate-800/50 backdrop-blur-sm p-5"
        >
          <div
            class="flex-none flex flex-col sm:flex-row sm:items-center justify-between gap-4 mb-4"
          >
            <h2 class="text-lg font-semibold text-white">Daftar Soal</h2>
            <div class="relative w-full sm:w-64 flex items-center">
              <svg
                class="absolute left-2.5 h-4 w-4 text-slate-500"
                fill="none"
                viewBox="0 0 24 24"
                stroke="currentColor"
              >
                <path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  stroke-width="2"
                  d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"
                />
              </svg>
              <input
                type="text"
                placeholder="Cari soal atau kuis..."
                bind:value={searchQuestion}
                class="w-full rounded-lg border border-slate-600 bg-slate-800 py-1.5 pl-9 pr-8 text-sm text-white placeholder-slate-500 focus:border-indigo-500 focus:outline-none focus:ring-1 focus:ring-indigo-500"
              />
              {#if searchQuestion}
                <button
                  type="button"
                  onclick={() => (searchQuestion = "")}
                  class="absolute right-2.5 text-slate-400 hover:text-white"
                >
                  <svg
                    class="h-4 w-4"
                    fill="none"
                    viewBox="0 0 24 24"
                    stroke="currentColor"
                  >
                    <path
                      stroke-linecap="round"
                      stroke-linejoin="round"
                      stroke-width="2"
                      d="M6 18L18 6M6 6l12 12"
                    />
                  </svg>
                </button>
              {/if}
            </div>
          </div>
          {#if questions.length === 0}
            <p class="text-sm text-slate-500">Belum ada soal.</p>
          {:else if filteredQuestions.length === 0}
            <p class="text-sm text-slate-500">Soal tidak ditemukan.</p>
          {:else}
            <div class="space-y-2 overflow-y-auto pr-2 flex-1">
              {#each filteredQuestions as question}
                <div
                  class="flex items-center justify-between rounded-lg border border-slate-700 px-3 py-2"
                >
                  <div class="min-w-0 flex-1">
                    <p class="truncate text-sm font-medium text-slate-200">
                      {question.question}
                    </p>
                    <span class="text-xs text-slate-500"
                      >{getQuizTitle(question.quiz_id)}</span
                    >
                  </div>
                  <div class="ml-2 flex shrink-0 gap-1">
                    <button
                      onclick={() => editQuestion(question)}
                      class="rounded px-2 py-1 text-xs text-indigo-400 hover:bg-indigo-500/10 transition"
                      >Edit</button
                    >
                    <button
                      onclick={() => deleteQuestion(question.id)}
                      class="rounded px-2 py-1 text-xs text-red-400 hover:bg-red-500/10 transition"
                      >Hapus</button
                    >
                  </div>
                </div>
              {/each}
            </div>
          {/if}
        </div>
      </div>
    {:else if activeTab === "scores"}
      <div class="space-y-4">
        <div class="flex items-center justify-between">
          <h2 class="text-lg font-semibold text-white">Semua Nilai</h2>
          {#if !resetConfirm}
            <button
              onclick={() => (resetConfirm = true)}
              class="rounded-lg bg-red-500/80 px-3 py-1.5 text-sm text-white hover:bg-red-500 transition"
              >Reset Semua Nilai</button
            >
          {:else}
            <div class="flex items-center gap-2">
              <span class="text-sm text-red-400">Yakin?</span>
              <button
                onclick={resetAllScores}
                class="rounded-lg bg-red-600 px-3 py-1.5 text-sm text-white hover:bg-red-500 transition"
                >Ya, Reset</button
              >
              <button
                onclick={() => (resetConfirm = false)}
                class="rounded-lg border border-slate-600 px-3 py-1.5 text-sm text-slate-300 hover:bg-slate-700 transition"
                >Batal</button
              >
            </div>
          {/if}
        </div>
        {#if scores.length === 0}
          <p class="text-sm text-slate-500">Belum ada nilai.</p>
        {:else}
          <div
            class="overflow-x-auto rounded-xl border border-slate-700 bg-slate-800/50 backdrop-blur-sm"
          >
            <table class="w-full text-left text-sm">
              <thead class="border-b border-slate-700">
                <tr>
                  <th class="px-4 py-3 font-medium text-slate-400">ID</th>
                  <th class="px-4 py-3 font-medium text-slate-400">Username</th>
                  <th class="px-4 py-3 font-medium text-slate-400">Kuis</th>
                  <th class="px-4 py-3 font-medium text-slate-400">Skor</th>
                  <th class="px-4 py-3 font-medium text-slate-400">Tanggal</th>
                </tr>
              </thead>
              <tbody>
                {#each scores as s}
                  <tr class="border-b border-slate-700/50 last:border-0">
                    <td class="px-4 py-3 text-slate-500">{s.id}</td>
                    <td class="px-4 py-3 text-slate-300">{s.username}</td>
                    <td class="px-4 py-3 text-slate-300">{s.quiz_title}</td>
                    <td class="px-4 py-3 font-semibold text-indigo-400"
                      >{s.score}</td
                    >
                    <td class="px-4 py-3 text-slate-500"
                      >{new Date(s.created_at).toLocaleDateString("id-ID")}</td
                    >
                  </tr>
                {/each}
              </tbody>
            </table>
          </div>
        {/if}
      </div>
    {:else if activeTab === "users"}
      <div class="space-y-4">
        <h2 class="text-lg font-semibold text-white">Semua Pengguna</h2>
        {#if users.length === 0}
          <p class="text-sm text-slate-500">Belum ada pengguna.</p>
        {:else}
          <div
            class="overflow-x-auto rounded-xl border border-slate-700 bg-slate-800/50 backdrop-blur-sm"
          >
            <table class="w-full text-left text-sm">
              <thead class="border-b border-slate-700">
                <tr>
                  <th class="px-4 py-3 font-medium text-slate-400">ID</th>
                  <th class="px-4 py-3 font-medium text-slate-400">Username</th>
                  <th class="px-4 py-3 font-medium text-slate-400">Role</th>
                  <th class="px-4 py-3 font-medium text-slate-400">Aksi</th>
                </tr>
              </thead>
              <tbody>
                {#each users as u}
                  <tr class="border-b border-slate-700/50 last:border-0">
                    <td class="px-4 py-3 text-slate-500">{u.id}</td>
                    <td class="px-4 py-3 text-slate-300 font-medium"
                      >{u.username}</td
                    >
                    <td class="px-4 py-3">
                      <select
                        class="rounded-md border border-slate-600 bg-slate-800 px-2 py-1 text-xs text-white focus:border-indigo-500 focus:outline-none"
                        value={u.role}
                        onchange={(e) =>
                          changeUserRole(u.id, e.currentTarget.value)}
                        disabled={u.username === auth.username}
                      >
                        <option value="student">Student</option>
                        <option value="teacher">Teacher</option>
                      </select>
                    </td>
                    <td class="px-4 py-3">
                      {#if u.username !== auth.username}
                        <button
                          onclick={() => (deleteUserId = u.id)}
                          class="rounded px-2 py-1 text-xs text-red-400 hover:bg-red-500/10 transition"
                          >Hapus</button
                        >
                      {:else}
                        <span class="text-xs text-slate-500 italic">Anda</span>
                      {/if}
                    </td>
                  </tr>
                {/each}
              </tbody>
            </table>
          </div>
        {/if}
      </div>
    {:else if activeTab === "logs"}
      <div class="space-y-4">
        <div class="flex items-center justify-between">
          <h2 class="text-lg font-semibold text-white">Log Aktivitas</h2>

          <div class="relative w-64">
            <input
              type="text"
              bind:value={searchLog}
              placeholder="Cari log..."
              class="w-full rounded-lg border border-slate-600 bg-slate-800 py-1.5 pl-3 pr-8 text-sm text-white placeholder-slate-400 focus:border-indigo-500 focus:outline-none focus:ring-1 focus:ring-indigo-500"
            />
            {#if searchLog}
              <button
                onclick={() => (searchLog = "")}
                class="absolute right-2 top-1/2 -translate-y-1/2 text-slate-400 hover:text-white"
                aria-label="Clear search"
              >
                <svg
                  class="h-4 w-4"
                  fill="none"
                  viewBox="0 0 24 24"
                  stroke="currentColor"
                >
                  <path
                    stroke-linecap="round"
                    stroke-linejoin="round"
                    stroke-width="2"
                    d="M6 18L18 6M6 6l12 12"
                  />
                </svg>
              </button>
            {/if}
          </div>
        </div>

        {#if logs.length === 0}
          <p class="text-sm text-slate-500">Belum ada aktivitas.</p>
        {:else if filteredLogs.length === 0}
          <p class="text-sm text-slate-500">Log tidak ditemukan.</p>
        {:else}
          <div
            class="overflow-x-auto rounded-xl border border-slate-700 bg-slate-800/50 backdrop-blur-sm max-h-[600px] overflow-y-auto custom-scrollbar"
          >
            <table class="w-full text-left text-sm relative">
              <thead
                class="border-b border-slate-700 bg-slate-800/90 backdrop-blur-sm sticky top-0"
              >
                <tr>
                  <th class="px-4 py-3 font-medium text-slate-400">Waktu</th>
                  <th class="px-4 py-3 font-medium text-slate-400">Username</th>
                  <th class="px-4 py-3 font-medium text-slate-400">Aktivitas</th
                  >
                </tr>
              </thead>
              <tbody>
                {#each filteredLogs as log}
                  <tr class="border-b border-slate-700/50 last:border-0">
                    <td class="px-4 py-3 text-slate-400 whitespace-nowrap">
                      {new Date(log.timestamp).toLocaleDateString("id-ID", {
                        year: "numeric",
                        month: "short",
                        day: "numeric",
                        hour: "2-digit",
                        minute: "2-digit",
                      })}
                    </td>
                    <td class="px-4 py-3 text-slate-300 font-medium"
                      >{log.username}</td
                    >
                    <td class="px-4 py-3">
                      {#if log.action === "REGISTER"}
                        <span
                          class="inline-flex items-center rounded-md bg-green-500/10 px-2 py-1 text-xs font-medium text-green-400 ring-1 ring-inset ring-green-500/20"
                          >Registrasi Akun Baru</span
                        >
                      {:else if log.action === "DELETE_ACCOUNT"}
                        <span
                          class="inline-flex items-center rounded-md bg-red-500/10 px-2 py-1 text-xs font-medium text-red-400 ring-1 ring-inset ring-red-500/20"
                          >Hapus Akun Sendiri</span
                        >
                      {:else if log.action === "DELETED_BY_ADMIN"}
                        <span
                          class="inline-flex items-center rounded-md bg-orange-500/10 px-2 py-1 text-xs font-medium text-orange-400 ring-1 ring-inset ring-orange-500/20"
                          >Dihapus Admin</span
                        >
                      {:else}
                        <span
                          class="inline-flex items-center rounded-md bg-slate-500/10 px-2 py-1 text-xs font-medium text-slate-400 ring-1 ring-inset ring-slate-500/20"
                          >{log.action}</span
                        >
                      {/if}
                    </td>
                  </tr>
                {/each}
              </tbody>
            </table>
          </div>
        {/if}
      </div>
    {/if}
  </div>

  {#if deleteQuizId}
    <!-- svelte-ignore a11y_click_events_have_key_events a11y_no_static_element_interactions -->
    <div
      class="fixed inset-0 z-50 flex items-center justify-center bg-black/50 backdrop-blur-sm"
      onclick={() => (deleteQuizId = null)}
    >
      <!-- svelte-ignore a11y_click_events_have_key_events a11y_no_static_element_interactions -->
      <div
        class="rounded-xl border border-slate-600 bg-slate-800 p-6 shadow-2xl"
        onclick={(e) => e.stopPropagation()}
      >
        <h3 class="text-lg font-semibold text-white">Konfirmasi Hapus</h3>
        <p class="mt-2 text-sm text-slate-400">
          Yakin menghapus <span class="font-medium text-white"
            >{getQuizTitle(deleteQuizId)}</span
          >? Semua soal di dalamnya juga akan ikut terhapus.
        </p>
        <div class="mt-4 flex justify-end gap-3">
          <button
            onclick={() => (deleteQuizId = null)}
            class="rounded-lg border border-slate-600 px-4 py-2 text-sm text-slate-300 hover:bg-slate-700 transition"
            >Batal</button
          >
          <button
            onclick={deleteQuiz}
            class="rounded-lg bg-red-500 px-4 py-2 text-sm text-white hover:bg-red-400 transition"
            >Hapus</button
          >
        </div>
      </div>
    </div>
  {/if}

  {#if deleteUserId}
    <!-- svelte-ignore a11y_click_events_have_key_events a11y_no_static_element_interactions -->
    <div
      class="fixed inset-0 z-50 flex items-center justify-center bg-black/50 backdrop-blur-sm"
      onclick={() => (deleteUserId = null)}
    >
      <!-- svelte-ignore a11y_click_events_have_key_events a11y_no_static_element_interactions -->
      <div
        class="rounded-xl border border-slate-600 bg-slate-800 p-6 shadow-2xl"
        onclick={(e) => e.stopPropagation()}
      >
        <h3 class="text-lg font-semibold text-white">
          Konfirmasi Hapus Pengguna
        </h3>
        <p class="mt-2 text-sm text-slate-400 max-w-md">
          Yakin menghapus pengguna ini? <strong class="text-red-400"
            >Semua nilai skor dari pengguna ini juga akan ikut terhapus
            selamanya.</strong
          >
        </p>
        <div class="mt-6 flex justify-end gap-3">
          <button
            onclick={() => (deleteUserId = null)}
            class="rounded-lg border border-slate-600 px-4 py-2 text-sm text-slate-300 hover:bg-slate-700 transition"
            >Batal</button
          >
          <button
            onclick={deleteUser}
            class="rounded-lg bg-red-500 px-4 py-2 text-sm text-white hover:bg-red-400 transition"
            >Hapus</button
          >
        </div>
      </div>
    </div>
  {/if}
{/if}
