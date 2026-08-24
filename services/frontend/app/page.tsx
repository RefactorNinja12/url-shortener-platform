import Link from "next/link";
import { cookies } from "next/headers";
import { SESSION_COOKIE } from "@/app/lib/session";
import ShortenForm from "@/app/components/ShortenForm";
import LogoutButton from "@/app/components/LogoutButton";

export default async function Home() {
  const cookieStore = await cookies();
  const isLoggedIn = Boolean(cookieStore.get(SESSION_COOKIE)?.value);

  return (
    <div className="flex min-h-screen flex-col items-center bg-zinc-50 font-sans dark:bg-black">
      <main className="flex w-full max-w-xl flex-1 flex-col items-center gap-6 px-6 py-32">
        <div className="flex w-full items-center justify-between">
          <h1 className="text-3xl font-semibold tracking-tight text-black dark:text-zinc-50">
            URL Shortener
          </h1>
          {isLoggedIn && <LogoutButton />}
        </div>

        {isLoggedIn ? (
          <ShortenForm />
        ) : (
          <div className="flex w-full flex-col items-center gap-3 rounded border border-zinc-300 bg-white p-6 text-center dark:border-zinc-700 dark:bg-zinc-900">
            <p className="text-zinc-700 dark:text-zinc-300">
              Du måste logga in för att skapa korta länkar.
            </p>
            <div className="flex gap-4">
              <Link href="/login" className="font-medium underline">
                Logga in
              </Link>
              <Link href="/register" className="font-medium underline">
                Skapa konto
              </Link>
            </div>
          </div>
        )}
      </main>
    </div>
  );
}
