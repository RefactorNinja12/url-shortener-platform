import { NextRequest, NextResponse } from "next/server";
import { SESSION_COOKIE } from "@/app/lib/session";

const GO_API_URL = process.env.GO_API_URL ?? "http://localhost:8080";

export async function POST(request: NextRequest) {
  const token = request.cookies.get(SESSION_COOKIE)?.value;
  if (!token) {
    return NextResponse.json({ error: "du måste vara inloggad" }, { status: 401 });
  }

  const body = await request.json().catch(() => null);
  const url = body?.url;

  if (typeof url !== "string" || url.trim() === "") {
    return NextResponse.json({ error: "url is required" }, { status: 400 });
  }

  let apiRes: Response;
  try {
    apiRes = await fetch(`${GO_API_URL}/shorten`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${token}`,
      },
      body: JSON.stringify({ url }),
    });
  } catch {
    return NextResponse.json({ error: "kunde inte nå go-api" }, { status: 502 });
  }

  if (!apiRes.ok) {
    const text = await apiRes.text();
    return NextResponse.json({ error: text || "kunde inte förkorta url:en" }, { status: apiRes.status });
  }

  const data = await apiRes.json();
  return NextResponse.json(data, { status: apiRes.status });
}
