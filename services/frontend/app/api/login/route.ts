import { NextRequest, NextResponse } from "next/server";
import { SESSION_COOKIE } from "@/app/lib/session";

const AUTH_API_URL = process.env.AUTH_API_URL ?? "http://localhost:5080";
const SESSION_MAX_AGE_SECONDS = 60 * 60; // matchar auth-admins JWT-livstid (1h)

export async function POST(request: NextRequest) {
  const body = await request.json().catch(() => null);
  const { email, password } = body ?? {};

  if (typeof email !== "string" || typeof password !== "string") {
    return NextResponse.json({ error: "email and password are required" }, { status: 400 });
  }

  let authRes: Response;
  try {
    authRes = await fetch(`${AUTH_API_URL}/login`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ email, password }),
    });
  } catch {
    return NextResponse.json({ error: "kunde inte nå auth-admin" }, { status: 502 });
  }

  if (!authRes.ok) {
    return NextResponse.json({ error: "fel e-post eller lösenord" }, { status: 401 });
  }

  const { token } = await authRes.json();
  const response = NextResponse.json({ ok: true });
  response.cookies.set(SESSION_COOKIE, token, {
    httpOnly: true,
    sameSite: "lax",
    secure: process.env.NODE_ENV === "production",
    path: "/",
    maxAge: SESSION_MAX_AGE_SECONDS,
  });
  return response;
}
