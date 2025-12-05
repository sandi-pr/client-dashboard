export const DataPresetCase = (name) => {
  const preset = [
    {
      name: "StartAgent",
      description: "Agent pertama yang memulai workflow dan meneruskan eksekusi ke IteratorAgent.",
      tools: [],
      handoff_to: ["IteratorAgent"],
      system_prompt:
        'Kamu StartAgent.\n\nTUGAS:\n1. SEGERA melakukan handoff ke IteratorAgent menggunakan tool internal handoff.\n\nFORMAT:\nGunakan tool:\n- to_agent = "IteratorAgent"\n- reason = "Mulai workflow"\n\nDILARANG:\n- Mengoutput JSON manual (gunakan tool handoff saja).\n- Mengakhiri workflow sendiri.\n- Memberikan penjelasan panjang.',
    },
    {
      name: "IteratorAgent",
      description: "Mengambil record berikutnya. Jika data habis, lanjut ke FinalAgent.",
      tools: ["fetch_next_record"],
      handoff_to: ["CompletenessAgent", "FinalAgent"],
      system_prompt:
        'Kamu IteratorAgent.\n\nSETIAP KALI AKTIF:\n1. WAJIB panggil tool fetch_next_record().\n\nHASIL SETELAH TOOL:\n\n- Jika hasil adalah None:\n    → WAJIB panggil tool handoff ke "FinalAgent"\n      dengan reason "Tidak ada data tersisa".\n\n- Jika hasil berisi data valid:\n    → WAJIB panggil tool handoff ke "CompletenessAgent"\n      dengan reason "Mulai pemeriksaan completeness untuk record ini".\n\nDILARANG:\n- Mengakhiri workflow sendiri.\n- Tidak memanggil tool.',
    },
    {
      name: "CompletenessAgent",
      description: "Evaluasi completeness setiap record.",
      tools: ["run_completeness_check"],
      handoff_to: ["ComparisonAgent"],
      system_prompt:
        'Kamu CompletenessAgent.\n\nTUGAS:\n1. Panggil tool run_completeness_check() satu kali.\n\nSETELAH TOOL:\n→ WAJIB panggil tool handoff ke ComparisonAgent\n   dengan reason "Pemeriksaan completeness selesai".\n\nDILARANG:\n- Melewati pemanggilan tool.\n- Mengakhiri workflow sendiri.',
    },
    {
      name: "ComparisonAgent",
      description: "Membandingkan hasil extraction dengan DWH.",
      tools: ["run_comparison_check"],
      handoff_to: ["AccuracyRouterAgent"],
      system_prompt:
        'Kamu ComparisonAgent.\n\nTUGAS:\n1. Panggil tool run_comparison_check().\n\nSETELAH TOOL:\n→ WAJIB panggil tool handoff ke AccuracyRouterAgent\n   dengan reason "Perbandingan selesai".\n\nDILARANG:\n- Mengakhiri workflow sendiri.',
    },
    {
      name: "Kondisi00CheckerAgent",
      description: "Agent validator untuk rule Kondisi00Checker.",
      tools: ["run_checker"],
      handoff_to: ["AccuracyRouterAgent"],
      system_prompt:
        'Kamu Kondisi00CheckerAgent.\n\nTUGAS:\n1. Jalankan tool run_checker() sesuai rule Kondisi00Checker.\n\nSETELAH TOOL:\n→ WAJIB panggil tool handoff ke AccuracyRouterAgent\n   dengan reason "Rule checker selesai".\n\nDILARANG:\n- Mengakhiri workflow sendiri.',
    },
    {
      name: "Kolektabilitas2CheckerAgent",
      description: "Agent validator untuk rule Kolektabilitas2Checker.",
      tools: ["run_checker"],
      handoff_to: ["AccuracyRouterAgent"],
      system_prompt: "...",
    },
    {
      name: "Kolektabilitas3CheckerAgent",
      description: "Agent validator untuk rule Kolektabilitas3Checker.",
      tools: ["run_checker"],
      handoff_to: ["AccuracyRouterAgent"],
      system_prompt: "...",
    },
    {
      name: "Kolektabilitas4CheckerAgent",
      description: "Agent validator untuk rule Kolektabilitas4Checker.",
      tools: ["run_checker"],
      handoff_to: ["AccuracyRouterAgent"],
      system_prompt: "...",
    },
    {
      name: "Kolektabilitas5CheckerAgent",
      description: "Agent validator untuk rule Kolektabilitas5Checker.",
      tools: ["run_checker"],
      handoff_to: ["AccuracyRouterAgent"],
      system_prompt: "...",
    },
    {
      name: "Kode02CheckerAgent",
      description: "Agent validator untuk rule Kode02Checker.",
      tools: ["run_checker"],
      handoff_to: ["AccuracyRouterAgent"],
      system_prompt: "...",
    },
    {
      name: "Kode03CheckerAgent",
      description: "Agent validator untuk rule Kode03Checker.",
      tools: ["run_checker"],
      handoff_to: ["AccuracyRouterAgent"],
      system_prompt: "...",
    },
    {
      name: "Kode04CheckerAgent",
      description: "Agent validator untuk rule Kode04Checker.",
      tools: ["run_checker"],
      handoff_to: ["AccuracyRouterAgent"],
      system_prompt: "...",
    },
    {
      name: "Kode05PlusCheckerAgent",
      description: "Agent validator untuk rule Kode05PlusChecker.",
      tools: ["run_checker"],
      handoff_to: ["AccuracyRouterAgent"],
      system_prompt: "...",
    },
  ];
  return preset.find((p) => p.name === name) || null;
};
