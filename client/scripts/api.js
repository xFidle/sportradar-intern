export async function apiFetch(path, options = {}) {
  const res = await fetch(`${path}`, {
    headers: { "Content-Type": "application/json", ...(options.headers || {}) },
    ...options
  })

  let body = null
  try {
    body = await res.json()
  } catch (_) {
    body = null
  }

  if (!res.ok) {
    const details = body && body.details ? `: ${JSON.stringify(body.details)}` : ""
    const message = body && body.error ? `${body.error}${details}` : `Request failed (${res.status})`
    throw new Error(message)
  }

  return body
}
