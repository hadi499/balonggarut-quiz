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

  onMount(async () => {
    try {
      quizzes = (await api.get<Quiz[]>("/api/quizzes")) || [];
    } catch {}

    loading = false;
  });
</script>

<svelte:head>
  <title>Balonggarut Quiz - Platform E-Learning Interaktif</title>
</svelte:head>

<div class="flex flex-col min-h-[calc(100vh-4rem)]">
  <!-- Hero Section -->
  <section
    class="relative flex-1 flex flex-col justify-center overflow-hidden pt-10 sm:pt-14 pb-16"
  >
    <div
      class="absolute inset-0 bg-[radial-gradient(ellipse_at_top,rgba(99,102,241,0.15),transparent_50%)] [mask-image:linear-gradient(to_bottom,transparent,black_15%,black)] pointer-events-none"
    ></div>
    <div class="mx-auto max-w-7xl px-6 lg:px-8 relative z-10">
      <div class="mx-auto max-w-3xl text-center">
        <div in:fly={{ y: -20, duration: 800, delay: 100 }}>
          <span
            class="inline-flex items-center rounded-full bg-indigo-500/10 px-3 py-1 text-sm font-medium text-indigo-400 ring-1 ring-inset ring-indigo-500/20 mb-8 shadow-[0_0_15px_rgba(99,102,241,0.2)]"
          >
            Platform E-Learning Interaktif
          </span>
        </div>
        <h1
          in:fly={{ y: 30, duration: 800, delay: 200 }}
          class="text-4xl font-extrabold tracking-tight text-white sm:text-6xl md:text-7xl bg-clip-text text-transparent bg-gradient-to-r from-indigo-400 via-purple-400 to-pink-400 drop-shadow-sm pb-2"
        >
          Belajar Lebih Seru dengan Balonggarut Quiz
        </h1>
        <p
          in:fly={{ y: 30, duration: 800, delay: 300 }}
          class="mt-6 text-lg leading-8 text-slate-300 max-w-2xl mx-auto"
        >
          Tingkatkan pengetahuanmu melalui kuis interaktif yang dirancang untuk
          membuat proses belajar menjadi menyenangkan, menantang, dan bermakna.
        </p>
        <div
          in:fly={{ y: 30, duration: 800, delay: 400 }}
          class="mt-10 flex flex-col sm:flex-row items-center justify-center gap-4"
        >
          {#if !loading && auth.isLoggedIn}
            <a
              href="/quizzes"
              class="rounded-full bg-indigo-600 px-8 py-3.5 text-base font-semibold text-white shadow-lg shadow-indigo-500/30 hover:bg-indigo-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-400 transition-all hover:scale-105 hover:shadow-indigo-500/50"
            >
              Lanjut Belajar, {auth.username}
            </a>
          {:else if !loading}
            <a
              href="/register"
              class="w-full sm:w-auto rounded-full bg-indigo-600 px-8 py-3.5 text-base font-semibold text-white shadow-lg shadow-indigo-500/30 hover:bg-indigo-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-400 transition-all hover:-translate-y-1 hover:shadow-indigo-500/50"
            >
              Mulai Sekarang Gratis
            </a>
            <a
              href="/login"
              class="w-full sm:w-auto rounded-full bg-slate-800 border border-slate-700 px-8 py-3.5 text-base font-semibold text-white hover:bg-slate-700 hover:border-slate-600 transition-all"
            >
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
        <h2 class="text-base font-semibold leading-7 text-indigo-400">
          Fitur Unggulan
        </h2>
        <p
          class="mt-2 text-3xl font-bold tracking-tight text-white sm:text-4xl"
        >
          Mengapa Memilih Kami?
        </p>
      </div>
      <div class="grid grid-cols-1 md:grid-cols-3 gap-8">
        <!-- Feature 1 -->
        <div
          in:fly={{ y: 50, duration: 800, delay: 200 }}
          class="relative group bg-slate-800/40 p-8 rounded-3xl border border-slate-700/50 backdrop-blur-md transition-all duration-300 hover:-translate-y-2 hover:shadow-2xl hover:shadow-indigo-500/20 hover:border-indigo-500/50"
        >
          <div
            class="absolute inset-0 bg-gradient-to-br from-indigo-500/5 to-purple-500/5 rounded-3xl opacity-0 group-hover:opacity-100 transition-opacity duration-500"
          ></div>
          <div class="relative z-10">
            <div
              class="h-14 w-14 rounded-2xl bg-indigo-500/10 flex items-center justify-center mb-6 group-hover:scale-110 group-hover:-rotate-3 transition-transform duration-300 shadow-inner shadow-indigo-500/20"
            >
              <svg
                class="h-7 w-7 text-indigo-400"
                fill="none"
                viewBox="0 0 24 24"
                stroke="currentColor"
              >
                <path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  stroke-width="2"
                  d="M13 10V3L4 14h7v7l9-11h-7z"
                />
              </svg>
            </div>
            <h3 class="text-xl font-bold text-white mb-3">
              Cepat & Interaktif
            </h3>
            <p class="text-slate-400 leading-relaxed">
              Antarmuka yang responsif dan menarik membuat pengalaman menjawab
              kuis menjadi lebih fokus dan menyenangkan.
            </p>
          </div>
        </div>

        <!-- Feature 2 -->
        <div
          in:fly={{ y: 50, duration: 800, delay: 350 }}
          class="relative group bg-slate-800/40 p-8 rounded-3xl border border-slate-700/50 backdrop-blur-md transition-all duration-300 hover:-translate-y-2 hover:shadow-2xl hover:shadow-purple-500/20 hover:border-purple-500/50"
        >
          <div
            class="absolute inset-0 bg-gradient-to-br from-purple-500/5 to-pink-500/5 rounded-3xl opacity-0 group-hover:opacity-100 transition-opacity duration-500"
          ></div>
          <div class="relative z-10">
            <div
              class="h-14 w-14 rounded-2xl bg-purple-500/10 flex items-center justify-center mb-6 group-hover:scale-110 group-hover:rotate-3 transition-transform duration-300 shadow-inner shadow-purple-500/20"
            >
              <svg
                class="h-7 w-7 text-purple-400"
                fill="none"
                viewBox="0 0 24 24"
                stroke="currentColor"
              >
                <path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  stroke-width="2"
                  d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"
                />
              </svg>
            </div>
            <h3 class="text-xl font-bold text-white mb-3">Tantangan Waktu</h3>
            <p class="text-slate-400 leading-relaxed">
              Latih kecepatan berpikirmu dengan batas waktu di setiap
              pertanyaan. Semakin cepat, semakin baik!
            </p>
          </div>
        </div>

        <!-- Feature 3 -->
        <div
          in:fly={{ y: 50, duration: 800, delay: 500 }}
          class="relative group bg-slate-800/40 p-8 rounded-3xl border border-slate-700/50 backdrop-blur-md transition-all duration-300 hover:-translate-y-2 hover:shadow-2xl hover:shadow-pink-500/20 hover:border-pink-500/50"
        >
          <div
            class="absolute inset-0 bg-gradient-to-br from-pink-500/5 to-rose-500/5 rounded-3xl opacity-0 group-hover:opacity-100 transition-opacity duration-500"
          ></div>
          <div class="relative z-10">
            <div
              class="h-14 w-14 rounded-2xl bg-pink-500/10 flex items-center justify-center mb-6 group-hover:scale-110 group-hover:-rotate-3 transition-transform duration-300 shadow-inner shadow-pink-500/20"
            >
              <svg
                class="h-7 w-7 text-pink-400"
                fill="none"
                viewBox="0 0 24 24"
                stroke="currentColor"
              >
                <path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  stroke-width="2"
                  d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"
                />
              </svg>
            </div>
            <h3 class="text-xl font-bold text-white mb-3">Beragam Kategori</h3>
            <p class="text-slate-400 leading-relaxed">
              Dari pengetahuan umum hingga mata pelajaran spesifik, temukan kuis
              yang sesuai dengan minat belajarmu.
            </p>
          </div>
        </div>
      </div>
    </div>
  </section>

  <!-- Popular Quizzes Section -->
  <section
    class="py-16 sm:py-24 relative overflow-hidden bg-slate-900/30 border-t border-slate-800"
  >
    <div class="mx-auto max-w-7xl px-6 lg:px-8">
      <div
        class="flex flex-col sm:flex-row items-start sm:items-center justify-between mb-12 gap-4"
      >
        <div>
          <h2 class="text-3xl font-bold tracking-tight text-white sm:text-4xl">
            Coba Kuis Sekarang
          </h2>
          <p class="mt-2 text-slate-400">
            Pilih dari berbagai kuis yang tersedia dan uji kemampuanmu.
          </p>
        </div>
        {#if quizzes.length > 3}
          <a
            href="/quizzes"
            class="inline-flex items-center gap-2 text-sm font-semibold text-indigo-400 hover:text-indigo-300 transition-colors bg-indigo-500/10 px-4 py-2 rounded-full border border-indigo-500/20 hover:bg-indigo-500/20"
          >
            Lihat semua kuis
            <svg
              class="w-4 h-4"
              fill="none"
              viewBox="0 0 24 24"
              stroke="currentColor"
            >
              <path
                stroke-linecap="round"
                stroke-linejoin="round"
                stroke-width="2"
                d="M17 8l4 4m0 0l-4 4m4-4H3"
              />
            </svg>
          </a>
        {/if}
      </div>

      {#if loading}
        <div class="flex justify-center py-20">
          <div
            class="h-10 w-10 animate-spin rounded-full border-4 border-slate-700 border-t-indigo-500"
          ></div>
        </div>
      {:else if quizzes.length === 0}
        <div
          class="text-center py-20 bg-slate-800/30 rounded-3xl border border-slate-700/50 border-dashed"
        >
          <p class="text-lg text-slate-400">
            Belum ada kuis tersedia saat ini.
          </p>
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
              <div
                class="absolute inset-0 bg-gradient-to-br from-indigo-500/5 via-transparent to-transparent opacity-0 group-hover:opacity-100 transition-opacity duration-500"
              ></div>

              <div class="relative z-10">
                <div class="flex items-start justify-between gap-x-4 mb-4">
                  <span
                    class="inline-flex items-center rounded-lg bg-indigo-500/10 px-2.5 py-1 text-xs font-semibold text-indigo-400 border border-indigo-500/20"
                  >
                    {quiz.category}
                  </span>
                </div>
                <h3
                  class="text-xl font-bold text-white group-hover:text-indigo-400 transition-colors"
                >
                  {quiz.title}
                </h3>
              </div>

              <div
                class="mt-8 flex items-center justify-between relative z-10 pt-6 border-t border-slate-700/50 group-hover:border-slate-700 transition-colors"
              >
                <div class="flex items-center gap-4">
                  <div
                    class="flex items-center gap-1.5 text-xs font-medium text-slate-300"
                  >
                    <svg
                      class="w-4 h-4 text-slate-500"
                      fill="none"
                      viewBox="0 0 24 24"
                      stroke="currentColor"
                    >
                      <path
                        stroke-linecap="round"
                        stroke-linejoin="round"
                        stroke-width="2"
                        d="M8.228 9c.549-1.165 2.03-2 3.772-2 2.21 0 4 1.343 4 3 0 1.4-1.278 2.575-3.006 2.907-.542.104-.994.54-.994 1.093m0 3h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
                      />
                    </svg>
                    {quiz.questions?.length ?? 0} Soal
                  </div>
                  <div
                    class="flex items-center gap-1.5 text-xs font-medium text-slate-300"
                  >
                    <svg
                      class="w-4 h-4 text-slate-500"
                      fill="none"
                      viewBox="0 0 24 24"
                      stroke="currentColor"
                    >
                      <path
                        stroke-linecap="round"
                        stroke-linejoin="round"
                        stroke-width="2"
                        d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"
                      />
                    </svg>
                    {quiz.timeLimit}s
                  </div>
                </div>
                <div
                  class="w-8 h-8 rounded-full bg-slate-700/50 flex items-center justify-center group-hover:bg-indigo-500 group-hover:text-white transition-colors text-slate-400"
                >
                  <svg
                    class="w-4 h-4"
                    fill="none"
                    viewBox="0 0 24 24"
                    stroke="currentColor"
                  >
                    <path
                      stroke-linecap="round"
                      stroke-linejoin="round"
                      stroke-width="2"
                      d="M9 5l7 7-7 7"
                    />
                  </svg>
                </div>
              </div>
            </a>
          {/each}
        </div>
      {/if}
    </div>
  </section>

  <!-- Footer -->
  <footer class="bg-slate-900 border-t border-slate-800 pt-16 pb-8">
    <div class="mx-auto max-w-7xl px-6 lg:px-8">
      <div class="grid grid-cols-1 md:grid-cols-3 gap-12 md:gap-8 mb-12">
        <div class="space-y-4">
          <h3
            class="text-xl font-bold bg-clip-text text-transparent bg-gradient-to-r from-indigo-400 to-purple-400"
          >
            Balonggarut Quiz
          </h3>
          <p class="text-sm text-slate-400 leading-relaxed max-w-xs">
            Platform e-learning interaktif modern untuk meningkatkan kualitas
            belajar melalui kuis yang menantang dan menyenangkan.
          </p>
        </div>
        <div>
          <h4 class="text-white font-semibold mb-4">Tautan Cepat</h4>
          <ul class="space-y-2 text-sm text-slate-400">
            <li>
              <a href="/" class="hover:text-indigo-400 transition-colors"
                >Beranda</a
              >
            </li>
            <li>
              <a href="/quizzes" class="hover:text-indigo-400 transition-colors"
                >Daftar Kuis</a
              >
            </li>
            <li>
              <a href="/login" class="hover:text-indigo-400 transition-colors"
                >Masuk</a
              >
            </li>
            <li>
              <a
                href="/register"
                class="hover:text-indigo-400 transition-colors">Daftar Akun</a
              >
            </li>
          </ul>
        </div>
        <div>
          <h4 class="text-white font-semibold mb-4">Hubungi Kami</h4>
          <ul class="space-y-2 text-sm text-slate-400">
            <li class="flex items-center gap-2">
              <svg
                class="w-4 h-4 text-slate-500"
                fill="none"
                viewBox="0 0 24 24"
                stroke="currentColor"
                ><path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  stroke-width="2"
                  d="M3 8l7.89 5.26a2 2 0 002.22 0L21 8M5 19h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z"
                /></svg
              >
              info@balonggarutquiz.com
            </li>
            <li class="flex items-center gap-2">
              <svg
                class="w-4 h-4 text-slate-500"
                fill="none"
                viewBox="0 0 24 24"
                stroke="currentColor"
                ><path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  stroke-width="2"
                  d="M17.657 16.657L13.414 20.9a1.998 1.998 0 01-2.827 0l-4.242-4.243a8 8 0 1111.314 0z"
                /><path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  stroke-width="2"
                  d="M15 11a3 3 0 11-6 0 3 3 0 016 0z"
                /></svg
              >
              Sidoarjo, Jawa Timur
            </li>
          </ul>
        </div>
      </div>
      <div
        class="pt-8 border-t border-slate-800 flex flex-col md:flex-row items-center justify-between gap-4"
      >
        <p class="text-xs text-slate-500">
          &copy; {new Date().getFullYear()} Balonggarut Quiz. Hak Cipta Dilindungi.
        </p>
        <div class="flex gap-4">
          <a
            href="#"
            class="text-slate-500 hover:text-indigo-400 transition-colors"
          >
            <span class="sr-only">Facebook</span>
            <svg
              class="h-5 w-5"
              fill="currentColor"
              viewBox="0 0 24 24"
              aria-hidden="true"
              ><path
                fill-rule="evenodd"
                d="M22 12c0-5.523-4.477-10-10-10S2 6.477 2 12c0 4.991 3.657 9.128 8.438 9.878v-6.987h-2.54V12h2.54V9.797c0-2.506 1.492-3.89 3.777-3.89 1.094 0 2.238.195 2.238.195v2.46h-1.26c-1.243 0-1.63.771-1.63 1.562V12h2.773l-.443 2.89h-2.33v6.988C18.343 21.128 22 16.991 22 12z"
                clip-rule="evenodd"
              /></svg
            >
          </a>
          <a
            href="#"
            class="text-slate-500 hover:text-indigo-400 transition-colors"
          >
            <span class="sr-only">Instagram</span>
            <svg
              class="h-5 w-5"
              fill="currentColor"
              viewBox="0 0 24 24"
              aria-hidden="true"
              ><path
                fill-rule="evenodd"
                d="M12.315 2c2.43 0 2.784.013 3.808.06 1.064.049 1.791.218 2.427.465a4.902 4.902 0 011.772 1.153 4.902 4.902 0 011.153 1.772c.247.636.416 1.363.465 2.427.048 1.067.06 1.407.06 4.123v.08c0 2.643-.012 2.987-.06 4.043-.049 1.064-.218 1.791-.465 2.427a4.902 4.902 0 01-1.153 1.772 4.902 4.902 0 01-1.772 1.153c-.636.247-1.363.416-2.427.465-1.067.048-1.407.06-4.123.06h-.08c-2.643 0-2.987-.012-4.043-.06-1.064-.049-1.791-.218-2.427-.465a4.902 4.902 0 01-1.772-1.153 4.902 4.902 0 01-1.153-1.772c-.247-.636-.416-1.363-.465-2.427-.047-1.024-.06-1.379-.06-3.808v-.63c0-2.43.013-2.784.06-3.808.049-1.064.218-1.791.465-2.427a4.902 4.902 0 011.153-1.772A4.902 4.902 0 015.45 2.525c.636-.247 1.363-.416 2.427-.465C8.901 2.013 9.256 2 11.685 2h.63zm-.081 1.802h-.468c-2.456 0-2.784.011-3.807.058-.975.045-1.504.207-1.857.344-.467.182-.8.398-1.15.748-.35.35-.566.683-.748 1.15-.137.353-.3.882-.344 1.857-.047 1.023-.058 1.351-.058 3.807v.468c0 2.456.011 2.784.058 3.807.045.975.207 1.504.344 1.857.182.466.399.8.748 1.15.35.35.683.566 1.15.748.353.137.882.3 1.857.344 1.054.048 1.37.058 4.041.058h.08c2.597 0 2.917-.01 3.96-.058.976-.045 1.505-.207 1.858-.344.466-.182.8-.398 1.15-.748.35-.35.566-.683.748-1.15.137-.353.3-.882.344-1.857.048-1.055.058-1.37.058-4.041v-.08c0-2.597-.01-2.917-.058-3.96-.045-.976-.207-1.505-.344-1.858a3.097 3.097 0 00-.748-1.15 3.098 3.098 0 00-1.15-.748c-.353-.137-.882-.3-1.857-.344-1.023-.047-1.351-.058-3.807-.058zM12 6.865a5.135 5.135 0 110 10.27 5.135 5.135 0 010-10.27zm0 1.802a3.333 3.333 0 100 6.666 3.333 3.333 0 000-6.666zm5.338-3.205a1.2 1.2 0 110 2.4 1.2 1.2 0 010-2.4z"
                clip-rule="evenodd"
              /></svg
            >
          </a>
        </div>
      </div>
    </div>
  </footer>
</div>
