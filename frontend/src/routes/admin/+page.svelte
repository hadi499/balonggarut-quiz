<script lang="ts">
  import { page } from "$app/state";
  import { api } from "$lib/api";
  import { auth } from "$lib/stores/auth.svelte";
  import { goto } from "$app/navigation";
  import { onMount } from "svelte";

  import QuizzesTab from "./QuizzesTab.svelte";
  import QuestionsTab from "./QuestionsTab.svelte";
  import ScoresTab from "./ScoresTab.svelte";
  import UsersTab from "./UsersTab.svelte";
  import LogsTab from "./LogsTab.svelte";

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

  let quizzes = $state<Quiz[]>([]);
  let questions = $state<Question[]>([]);
  let users = $state<User[]>([]);
  let scores = $state<ScoreEntry[]>([]);
  let loading = $state(true);

  let activeTab = $state<"quizzes" | "questions" | "scores" | "users" | "logs">(
    "quizzes",
  );

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
    try {
      const results = await Promise.allSettled([
        api.get<Quiz[]>("/api/quizzes"),
        api.get<Question[]>("/api/questions"),
        api.get<ScoreEntry[]>("/api/admin/scores"),
        api.get<User[]>("/api/admin/users"),
      ]);
      if (results[0].status === "fulfilled") quizzes = results[0].value;
      if (results[1].status === "fulfilled") questions = results[1].value;
      if (results[2].status === "fulfilled") scores = results[2].value || [];
      if (results[3].status === "fulfilled") users = results[3].value || [];
    } finally {
      loading = false;
    }
  }
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
      class="text-lg border-0 py-1 px-2 font-bold text-slate-200 w-fit text-center bg-indigo-500 rounded-md te"
    >
      Admin Panel
    </h1>

    <div
      class="flex overflow-x-auto gap-2 border-b border-slate-700 [&::-webkit-scrollbar]:hidden [-ms-overflow-style:none] [scrollbar-width:none]"
    >
      <button
        onclick={() => switchTab("quizzes")}
        class="whitespace-nowrap px-4 py-2 text-md transition-colors {activeTab ===
        'quizzes'
          ? 'border-b-2 border-indigo-500 font-medium text-indigo-400'
          : 'text-slate-400 hover:text-slate-200'}"
      >
        Kuis
      </button>
      <button
        onclick={() => switchTab("questions")}
        class="whitespace-nowrap px-4 py-2 text-md transition-colors {activeTab ===
        'questions'
          ? 'border-b-2 border-indigo-500 font-medium text-indigo-400'
          : 'text-slate-400 hover:text-slate-200'}"
      >
        Soal
      </button>
      <button
        onclick={() => switchTab("scores")}
        class="whitespace-nowrap px-4 py-2 text-md transition-colors {activeTab ===
        'scores'
          ? 'border-b-2 border-indigo-500 font-medium text-indigo-400'
          : 'text-slate-400 hover:text-slate-200'}"
      >
        Nilai
      </button>
      <button
        onclick={() => switchTab("users")}
        class="whitespace-nowrap px-4 py-2 text-md transition-colors {activeTab ===
        'users'
          ? 'border-b-2 border-indigo-500 font-medium text-indigo-400'
          : 'text-slate-400 hover:text-slate-200'}"
      >
        Pengguna
      </button>
      <button
        onclick={() => switchTab("logs")}
        class="whitespace-nowrap px-4 py-2 text-md transition-colors {activeTab ===
        'logs'
          ? 'border-b-2 border-indigo-500 font-medium text-indigo-400'
          : 'text-slate-400 hover:text-slate-200'}"
      >
        Log Aktivitas
      </button>
    </div>

    {#if activeTab === "quizzes"}
      <QuizzesTab {quizzes} refresh={loadAll} />
    {:else if activeTab === "questions"}
      <QuestionsTab {questions} {quizzes} refresh={loadAll} />
    {:else if activeTab === "scores"}
      <ScoresTab {scores} refresh={loadAll} />
    {:else if activeTab === "users"}
      <UsersTab {users} refresh={loadAll} />
    {:else if activeTab === "logs"}
      <LogsTab />
    {/if}
  </div>
{/if}
