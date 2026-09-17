"use client";

import { useEffect, useState } from "react";

interface UserProfile {
  id: string;
  email: string;
  name: string;
  avatar_url: string;
}

const BACKEND_URL = process.env.NEXT_PUBLIC_BACKEND_URL || "http://localhost:8080";

export default function Home() {
  const [user, setUser] = useState<UserProfile | null>(null);
  const [loading, setLoading] = useState(true);
  const [errorMessage] = useState<string | null>(() => {
    if (typeof window !== "undefined") {
      const err = new URLSearchParams(window.location.search).get("error");
      return err ? decodeURIComponent(err) : null;
    }
    return null;
  });

  useEffect(() => {
    // Kiểm tra phiên đăng nhập hiện tại qua HttpOnly Cookie
    fetch(`${BACKEND_URL}/api/auth/me`, {
      method: "GET",
      credentials: "include", // Rất quan trọng để gửi kèm HttpOnly Cookie
    })
      .then((res) => {
        if (res.ok) return res.json();
        throw new Error("Not authenticated");
      })
      .then((res) => {
        if (res.success && res.data) {
          setUser(res.data);
        }
      })
      .catch(() => {
        setUser(null);
      })
      .finally(() => {
        setLoading(false);
      });
  }, []);

  const handleLogout = async () => {
    try {
      const res = await fetch(`${BACKEND_URL}/api/auth/logout`, {
        method: "POST",
        credentials: "include",
      });
      if (res.ok) {
        setUser(null);
      }
    } catch (e) {
      console.error("Logout error:", e);
    }
  };

  return (
    <div className="flex min-h-screen flex-col items-center justify-center bg-zinc-50 p-6 font-sans dark:bg-zinc-950">
      <main className="w-full max-w-md rounded-2xl border border-zinc-200 bg-white p-8 shadow-sm dark:border-zinc-800 dark:bg-zinc-900">
        {/* Header */}
        <div className="mb-6 flex flex-col items-center text-center">
          <div className="mb-3 flex h-12 w-12 items-center justify-center rounded-full bg-emerald-100 text-emerald-600 dark:bg-emerald-950/50 dark:text-emerald-400">
            <svg
              className="h-6 w-6"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z"
              />
            </svg>
          </div>
          <h1 className="text-2xl font-bold tracking-tight text-zinc-900 dark:text-zinc-50">
            Field Booking
          </h1>
          <p className="mt-1 text-sm text-zinc-500 dark:text-zinc-400">
            Hệ thống xác thực Google OAuth 2.0
          </p>
        </div>

        {/* Thông báo lỗi nếu có */}
        {errorMessage && (
          <div className="mb-6 rounded-lg border border-red-200 bg-red-50 p-3 text-sm text-red-600 dark:border-red-900/50 dark:bg-red-950/50 dark:text-red-400">
            <div className="flex items-center gap-2">
              <svg className="h-4 w-4 shrink-0" fill="currentColor" viewBox="0 0 20 20">
                <path
                  fillRule="evenodd"
                  d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7 4a1 1 0 11-2 0 1 1 0 012 0zm-1-9a1 1 0 00-1 1v4a1 1 0 102 0V6a1 1 0 00-1-1z"
                  clipRule="evenodd"
                />
              </svg>
              <span>{errorMessage}</span>
            </div>
          </div>
        )}

        {/* Nội dung trạng thái */}
        {loading ? (
          <div className="flex justify-center py-8">
            <div className="h-8 w-8 animate-spin rounded-full border-2 border-zinc-300 border-t-emerald-600" />
          </div>
        ) : user ? (
          /* Đã đăng nhập */
          <div className="flex flex-col items-center gap-4 text-center">
            {user.avatar_url ? (
              /* eslint-disable-next-line @next/next/no-img-element */
              <img
                src={user.avatar_url}
                alt={user.name}
                className="h-20 w-20 rounded-full border-2 border-emerald-500 object-cover shadow-sm"
              />
            ) : (
              <div className="flex h-20 w-20 items-center justify-center rounded-full bg-zinc-200 text-2xl font-semibold text-zinc-700 dark:bg-zinc-800 dark:text-zinc-300">
                {user.name ? user.name.charAt(0).toUpperCase() : "U"}
              </div>
            )}
            <div>
              <h2 className="text-lg font-semibold text-zinc-900 dark:text-zinc-100">
                {user.name}
              </h2>
              <p className="text-sm text-zinc-500 dark:text-zinc-400">{user.email}</p>
            </div>

            <div className="w-full rounded-lg bg-zinc-100 p-3 text-left text-xs font-mono text-zinc-600 dark:bg-zinc-800 dark:text-zinc-300">
              <span className="font-semibold">User ID:</span> {user.id}
            </div>

            <button
              onClick={handleLogout}
              className="mt-2 flex w-full items-center justify-center gap-2 rounded-xl bg-red-600 px-4 py-2.5 text-sm font-medium text-white transition-colors hover:bg-red-700"
            >
              Đăng xuất
            </button>
          </div>
        ) : (
          /* Chưa đăng nhập */
          <div className="flex flex-col gap-4">
            <p className="text-center text-sm text-zinc-600 dark:text-zinc-400">
              Đăng nhập an toàn qua tài khoản Google với HttpOnly Cookie.
            </p>

            <a
              href={`${BACKEND_URL}/api/auth/google`}
              className="flex w-full items-center justify-center gap-3 rounded-xl border border-zinc-300 bg-white px-4 py-3 text-sm font-medium text-zinc-700 shadow-sm transition-all hover:bg-zinc-50 dark:border-zinc-700 dark:bg-zinc-800 dark:text-zinc-200 dark:hover:bg-zinc-750"
            >
              <svg className="h-5 w-5" viewBox="0 0 24 24">
                <path
                  fill="#4285F4"
                  d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92c-.26 1.37-1.04 2.53-2.21 3.31v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.09z"
                />
                <path
                  fill="#34A853"
                  d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z"
                />
                <path
                  fill="#FBBC05"
                  d="M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.06H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.94l2.85-2.22.81-.63z"
                />
                <path
                  fill="#EA4335"
                  d="M12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.06l3.66 2.84c.87-2.6 3.3-4.52 6.16-4.52z"
                />
              </svg>
              <span>Continue with Google</span>
            </a>
          </div>
        )}
      </main>
    </div>
  );
}
