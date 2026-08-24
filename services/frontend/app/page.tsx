"use client";

import { useState } from "react";

export default function Home() {
  const [url, setUrl] = useState("");
  const [shortUrl, setShortUrl] = useState<string | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault();
    setError(null);
    setShortUrl(null);
    setLoading(true);

    try {
      const res = await fetch("/api/shorten", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ url }),
      });
      const data = await res.json();

      if (!res.ok) {
        setError(data.error ?? "Något gick fel");
        return;
      }

      setShortUrl(data.shortUrl);
    } catch {
      setError("Kunde inte nå servern");
    } finally {
      setLoading(false);
    }
  }

  return (
    <div className="flex min-h-screen flex-col items-center bg-zinc-50 font-sans dark:bg-black">
      <main className="flex w-full max-w-xl flex-1 flex-col items-center gap-6 px-6 py-32">
        <h1 className="text-3xl font-semibold tracking-tight text-black dark:text-zinc-50">
          URL Shortener
        </h1>

        <form onSubmit={handleSubmit} className="flex w-full gap-2">
          <input
            type="url"
            required
            value={url}
            onChange={(e) => setUrl(e.target.value)}
            placeholder="https://din-langa-url.se/..."
            className="flex-1 rounded border border-zinc-300 bg-white px-3 py-2 text-black focus:border-black focus:outline-none dark:border-zinc-700 dark:bg-zinc-900 dark:text-zinc-50"
          />
          <button
            type="submit"
            disabled={loading}
            className="rounded bg-black px-5 py-2 font-medium text-white transition-colors hover:bg-zinc-800 disabled:opacity-50 dark:bg-zinc-50 dark:text-black dark:hover:bg-zinc-200"
          >
            {loading ? "Skapar…" : "Förkorta"}
          </button>
        </form>

        {shortUrl && (
          <div className="w-full rounded border border-green-300 bg-green-50 p-3 dark:border-green-800 dark:bg-green-950">
            <a
              href={shortUrl}
              target="_blank"
              rel="noreferrer"
              className="break-all font-mono text-sm text-green-800 underline dark:text-green-300"
            >
              {shortUrl}
            </a>
          </div>
        )}

        {error && (
          <div className="w-full rounded border border-red-300 bg-red-50 p-3 text-sm text-red-700 dark:border-red-800 dark:bg-red-950 dark:text-red-300">
            {error}
          </div>
        )}
      </main>
    </div>
  );
}
