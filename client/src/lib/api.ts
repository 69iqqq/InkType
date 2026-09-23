export async function getAuthToken(): Promise<string> {
  if (typeof window === "undefined") return "";

  // 1. If Clerk session is already active, return token immediately
  if (window.Clerk?.session) {
    const token = await window.Clerk.session.getToken();
    if (token) return token;
  }

  // 2. If Clerk is still hydrating on page refresh, wait up to 3 seconds
  for (let i = 0; i < 30; i++) {
    await new Promise((resolve) => setTimeout(resolve, 100));
    if (window.Clerk?.session) {
      const token = await window.Clerk.session.getToken();
      if (token) return token;
    }
  }

  return "";
}

export async function uploadDocument(file: File): Promise<Record<string, unknown>> {
  const token = await getAuthToken();
  const apiBase = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080/api/v1";

  const res = await fetch(`${apiBase}/documents`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      "Authorization": `Bearer ${token}`
    },
    body: JSON.stringify({ filename: file.name })
  });

  if (!res.ok) {
    throw new Error(`Failed to get presigned URL: ${res.statusText}`);
  }

  const data = await res.json();
  const { upload_url, document } = data;

  const uploadRes = await fetch(upload_url, {
    method: "PUT",
    body: file,
    headers: {
      "Content-Type": file.type || "application/pdf",
    }
  });

  if (!uploadRes.ok) {
    throw new Error(`Failed to upload to S3: ${uploadRes.statusText}`);
  }

  // Trigger document conversion pipeline
  const convertRes = await fetch(`${apiBase}/documents/${document.id}/convert`, {
    method: "POST",
    headers: {
      "Authorization": `Bearer ${token}`
    }
  });

  if (!convertRes.ok) {
    throw new Error(`Failed to trigger conversion: ${convertRes.statusText}`);
  }

  return document as Record<string, unknown>;
}

export async function deleteDocument(id: string): Promise<boolean> {
  const token = await getAuthToken();
  const apiBase = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080/api/v1";

  const res = await fetch(`${apiBase}/documents/${id}`, {
    method: "DELETE",
    headers: {
      "Authorization": `Bearer ${token}`
    }
  });

  return res.ok;
}
