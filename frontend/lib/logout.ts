export async function logout() {
  try {
    const res = await fetch("/api/backend/auth/logout", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
    });

    if (res.ok) {
      window.location.href = "/login";
    } else {
      console.error("Logout failed on server side");
    }
  } catch (error) {
    console.error("Network error during logout:", error);
  }
}
