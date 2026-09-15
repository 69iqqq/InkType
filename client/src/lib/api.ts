export async function getAuthToken(): Promise<string> {
  return "placeholder-jwt-token";
}

export async function uploadDocument(file: File): Promise<Record<string, unknown>> {
  const token = await getAuthToken();
  const apiBase = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080/v1";

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

  return document as Record<string, unknown>;
}
