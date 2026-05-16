<script lang="ts">
  import { auth } from "$lib/stores/auth.svelte";
  import { api } from "$lib/api";
  import { onMount } from "svelte";
  import { fade, fly } from "svelte/transition";

  interface Quiz {
    id: number;
    title: string;
    category: string;
    timeLimit: number;
    questions: { id: number }[];
  }

  let quizzes = $state<Quiz[]>([]);
  let loading = $state(true);
  let user = $state<{ username: string; role: string } | null>(null);

  onMount(async () => {
    try {
      quizzes = (await api.get<Quiz[]>("/api/quizzes")) || [];
    } catch {}

    if (auth.isLoggedIn) {
      try {
        const data = await api.get<{ username: string; role: string }>("/me");
        user = data;
        auth.login(auth.token!, data.username, data.role);
      } catch {
        auth.logout();
      }
    }

    loading = false;
  });
</script>

<svelte:head>
  <title>Balonggarut Quiz - Platform E-Learning Interaktif</title>
</svelte:head>

<div class="flex flex-col min-h-[calc(100vh-4rem)]">
  <!-- Hero Section -->
  <section class="relative flex-1 flex flex-col justify-center overflow-hidden pt-16 sm:pt-24 lg:pt-32 pb-16">
    <div class="absolute inset-0 bg-[radial-gradient(ellipse_at_top,rgba(99,102,241,0.15),transparent_50%)] [mask-image:linear-gradient(to_bottom,transparent,black_15%,black)] pointer-events-none"></div>
    <div class="mx-auto max-w-7xl px-6 lg:px-8 relative z-10">
      <div class="mx-auto max-w-3xl text-center">
        <div in:fly={{ y: -20, duration: 800, delay: 100 }}>
          <span class="inline-flex items-center rounded-full bg-indigo-500/10 px-3 py-1 text-sm font-medium text-indigo-400 ring-1 ring-inset ring-indigo-500/20 mb-8 shadow-[0_0_15px_rgba(99,102,241,0.2)]">
            Platform E-Learning Interaktif
          </span>
        </div>
        <h1 in:fly={{ y: 30, duration: 800, delay: 200 }} class="text-4xl font-extrabold tracking-tight text-white sm:text-6xl md:text-7xl bg-clip-text text-transparent bg-gradient-to-r from-indigo-400 via-purple-400 to-pink-400 drop-shadow-sm pb-2">
          Belajar Lebih Seru dengan Balonggarut Quiz
        </h1>
        <p in:fly={{ y: 30, duration: 800, delay: 300 }} class="mt-6 text-lg leading-8 text-slate-300 max-w-2xl mx-auto">
          Tingkatkan pengetahuanmu melalui kuis interaktif yang dirancang untuk membuat proses belajar menjadi menyenangkan, menantang, dan bermakna.
        </p>
        <div in:fly={{ y: 30, duration: 800, delay: 400 }} class="mt-10 flex flex-col sm:flex-row items-center justify-center gap-4">
          {#if !loading && auth.isLoggedIn}
            <a href="/quizzes" class="rounded-full bg-indigo-600 px-8 py-3.5 text-base font-semibold text-white shadow-lg shadow-indigo-500/30 hover:bg-indigo-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-400 transition-all hover:scale-105 hover:shadow-indigo-500/50">
              Lanjut Belajar, {user?.username}
            </a>
          {:else if !loading}
            <a href="/register" class="w-full sm:w-auto rounded-full bg-indigo-600 px-8 py-3.5 text-base font-semibold text-white shadow-lg shadow-indigo-500/30 hover:bg-indigo-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-400 transition-all hover:-translate-y-1 hover:shadow-indigo-500/50">
              Mulai Sekarang Gratis
            </a>
            <a href="/login" class="w-full sm:w-auto rounded-full bg-slate-800 border border-slate-700 px-8 py-3.5 text-base font-semibold text-white hover:bg-slate-700 hover:border-slate-600 transition-all">
              Masuk
            </a>
          {/if}
        </div>
      </div>
    </div>
  </section>

  <!-- Features Section -->
  <section class="py-16 sm:py-24 relative">
    <div class="mx-auto max-w-7xl px-6 lg:px-8">
      <div class="mx-auto max-w-2xl text-center mb-16">
        <h2 class="text-base font-semibold leading-7 text-indigo-400">Fitur Unggulan</h2>
        <p class="mt-2 text-3xl font-bold tracking-tight text-white sm:text-4xl">Mengapa Memilih Kami?</p>
      </div>
      <div class="grid grid-cols-1 md:grid-cols-3 gap-8">
        <!-- Feature 1 -->
        <div in:fly={{ y: 50, duration: 800, delay: 200 }} class="relative group bg-slate-800/40 p-8 rounded-3xl border border-slate-700/50 backdrop-blur-md transition-all duration-300 hover:-translate-y-2 hover:shadow-2xl hover:shadow-indigo-500/20 hover:border-indigo-500/50">
          <div class="absolute inset-0 bg-gradient-to-br from-indigo-500/5 to-purple-500/5 rounded-3xl opacity-0 group-hover:opacity-100 transition-opacity duration-500"></div>
          <div class="relative z-10">
            <div class="h-14 w-14 rounded-2xl bg-indigo-500/10 flex items-center justify-center mb-6 group-hover:scale-110 group-hover:-rotate-3 transition-transform duration-300 shadow-inner shadow-indigo-500/20">
              <svg class="h-7 w-7 text-indigo-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 10V3L4 14h7v7l9-11h-7z" />
              </svg>
            </div>
            <h3 class="text-xl font-bold text-white mb-3">Cepat & Interaktif</h3>
            <p class="text-slate-400 leading-relaxed">Antarmuka yang responsif dan menarik membuat pengalaman menjawab kuis menjadi lebih fokus dan menyenangkan.</p>
          </div>
        </div>
        
        <!-- Feature 2 -->
        <div in:fly={{ y: 50, duration: 800, delay: 350 }} class="relative group bg-slate-800/40 p-8 rounded-3xl border border-slate-700/50 backdrop-blur-md transition-all duration-300 hover:-translate-y-2 hover:shadow-2xl hover:shadow-purple-500/20 hover:border-purple-500/50">
          <div class="absolute inset-0 bg-gradient-to-br from-purple-500/5 to-pink-500/5 rounded-3xl opacity-0 group-hover:opacity-100 transition-opacity duration-500"></div>
          <div class="relative z-10">
            <div class="h-14 w-14 rounded-2xl bg-purple-500/10 flex items-center justify-center mb-6 group-hover:scale-110 group-hover:rotate-3 transition-transform duration-300 shadow-inner shadow-purple-500/20">
              <svg class="h-7 w-7 text-purple-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
            </div>
            <h3 class="text-xl font-bold text-white mb-3">Tantangan Waktu</h3>
            <p class="text-slate-400 leading-relaxed">Latih kecepatan berpikirmu dengan batas waktu di setiap pertanyaan. Semakin cepat, semakin baik!</p>
          </div>
        </div>
        
        <!-- Feature 3 -->
        <div in:fly={{ y: 50, duration: 800, delay: 500 }} class="relative group bg-slate-800/40 p-8 rounded-3xl border border-slate-700/50 backdrop-blur-md transition-all duration-300 hover:-translate-y-2 hover:shadow-2xl hover:shadow-pink-500/20 hover:border-pink-500/50">
          <div class="absolute inset-0 bg-gradient-to-br from-pink-500/5 to-rose-500/5 rounded-3xl opacity-0 group-hover:opacity-100 transition-opacity duration-500"></div>
          <div class="relative z-10">
            <div class="h-14 w-14 rounded-2xl bg-pink-500/10 flex items-center justify-center mb-6 group-hover:scale-110 group-hover:-rotate-3 transition-transform duration-300 shadow-inner shadow-pink-500/20">
              <svg class="h-7 w-7 text-pink-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
            </div>
            <h3 class="text-xl font-bold text-white mb-3">Beragam Kategori</h3>
            <p class="text-slate-400 leading-relaxed">Dari pengetahuan umum hingga mata pelajaran spesifik, temukan kuis yang sesuai dengan minat belajarmu.</p>
          </div>
        </div>
      </div>
    </div>
  </section>

  <!-- Popular Quizzes Section -->
  <section class="py-16 sm:py-24 relative overflow-hidden bg-slate-900/30 border-t border-slate-800">
    <div class="mx-auto max-w-7xl px-6 lg:px-8">
      <div class="flex flex-col sm:flex-row items-start sm:items-center justify-between mb-12 gap-4">
        <div>
          <h2 class="text-3xl font-bold tracking-tight text-white sm:text-4xl">Coba Kuis Sekarang</h2>
          <p class="mt-2 text-slate-400">Pilih dari berbagai kuis yang tersedia dan uji kemampuanmu.</p>
        </div>
        {#if quizzes.length > 3}
          <a href="/quizzes" class="inline-flex items-center gap-2 text-sm font-semibold text-indigo-400 hover:text-indigo-300 transition-colors bg-indigo-500/10 px-4 py-2 rounded-full border border-indigo-500/20 hover:bg-indigo-500/20">
            Lihat semua kuis
            <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 8l4 4m0 0l-4 4m4-4H3" />
            </svg>
          </a>
        {/if}
      </div>
      
      {#if loading}
        <div class="flex justify-center py-20">
          <div class="h-10 w-10 animate-spin rounded-full border-4 border-slate-700 border-t-indigo-500"></div>
        </div>
      {:else if quizzes.length === 0}
        <div class="text-center py-20 bg-slate-800/30 rounded-3xl border border-slate-700/50 border-dashed">
          <p class="text-lg text-slate-400">Belum ada kuis tersedia saat ini.</p>
          <p class="text-sm text-slate-500 mt-2">Silakan kembali lagi nanti!</p>
        </div>
      {:else}
        <div class="grid gap-6 sm:grid-cols-2 lg:grid-cols-3">
          {#each quizzes.slice(0, 3) as quiz, i}
            <a
              href="/quizzes/{quiz.id}"
              in:fly={{ y: 50, duration: 800, delay: 200 + i * 150 }}
              class="group relative flex flex-col justify-between rounded-3xl bg-slate-800/80 p-8 ring-1 ring-slate-700/50 transition-all duration-300 hover:ring-indigo-500/50 hover:bg-slate-800 hover:-translate-y-1 hover:shadow-2xl hover:shadow-indigo-500/10 overflow-hidden"
            >
              <div class="absolute inset-0 bg-gradient-to-br from-indigo-500/5 via-transparent to-transparent opacity-0 group-hover:opacity-100 transition-opacity duration-500"></div>
              
              <div class="relative z-10">
                <div class="flex items-start justify-between gap-x-4 mb-4">
                  <span class="inline-flex items-center rounded-lg bg-indigo-500/10 px-2.5 py-1 text-xs font-semibold text-indigo-400 border border-indigo-500/20">
                    {quiz.category}
                  </span>
                </div>
                <h3 class="text-xl font-bold text-white group-hover:text-indigo-400 transition-colors mb-2">
                  {quiz.title}
                </h3>
                <p class="text-sm text-slate-400 line-clamp-2">
                  Uji kemampuanmu dalam {quiz.category.toLowerCase()} dengan {quiz.questions?.length ?? 0} pertanyaan menantang.
                </p>
              </div>
              
              <div class="mt-8 flex items-center justify-between relative z-10 pt-6 border-t border-slate-700/50 group-hover:border-slate-700 transition-colors">
                <div class="flex items-center gap-4">
                  <div class="flex items-center gap-1.5 text-xs font-medium text-slate-300">
                    <svg class="w-4 h-4 text-slate-500" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8.228 9c.549-1.165 2.03-2 3.772-2 2.21 0 4 1.343 4 3 0 1.4-1.278 2.575-3.006 2.907-.542.104-.994.54-.994 1.093m0 3h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                    </svg>
                    {quiz.questions?.length ?? 0} Soal
                  </div>
                  <div class="flex items-center gap-1.5 text-xs font-medium text-slate-300">
                    <svg class="w-4 h-4 text-slate-500" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
                    </svg>
                    {quiz.timeLimit}s
                  </div>
                </div>
                <div class="w-8 h-8 rounded-full bg-slate-700/50 flex items-center justify-center group-hover:bg-indigo-500 group-hover:text-white transition-colors text-slate-400">
                  <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
                  </svg>
                </div>
              </div>
            </a>
          {/each}
        </div>
      {/if}
    </div>
  </section>
</div>
