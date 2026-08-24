import { NextRequest, NextResponse } from "next/server";

const AUTH_API_URL = process.env.AUTH_API_URL ?? "http://localhost:5080";

export async function POST(request: NextRequest) {
  const body = await request.json().catch(() => null);
  const { email, password } = body ?? {};

  if (typeof email !== "string" || typeof password !== "string") {
    return NextResponse.json({ error: "email and password are required" }, { status: 400 });
  }

  let authRes: Response;
  try {
    authRes = await fetch(`${AUTH_API_URL}/register`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ email, password }),
    });
  } catch {
    return NextResponse.json({ error: "kunde inte nå auth-admin" }, { status: 502 });
  }

  if (authRes.status === 409) {
    return NextResponse.json({ error: "e-postadressen är redan registrerad" }, { status: 409 });
  }
  if (!authRes.ok) {
    return NextResponse.json({ error: "kunde inte registrera användaren" }, { status: authRes.status });
  }

  return NextResponse.json({ ok: true }, { status: 201 });
}
