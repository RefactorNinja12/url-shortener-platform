"use client";

import { useRouter } from "next/navigation";

export default function LogoutButton() {
  const router = useRouter();

  async function handleLogout() {
    await fetch("/api/logout", { method: "POST" });
    router.refresh();
  }

  return (
    <button
      onClick={handleLogout}
      className="text-sm text-zinc-600 underline hover:text-black dark:text-zinc-400 dark:hover:text-zinc-50"
    >
      Logga ut
    </button>
  );
}
