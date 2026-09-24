const API_BASE = "http://localhost:8080"

export async function register(username, password, fullName) {
  const res = await fetch(`${API_BASE}/register`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ username, password, full_name: fullName }),
  })
  const data = await res.json()
  if (!res.ok) throw new Error(data.error)
  return data
}

export async function login(username, password) {
  const res = await fetch(`${API_BASE}/login`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ username, password }),
  })
  const data = await res.json()
  if (!res.ok) throw new Error(data.error)
  return data.token
}

export function decodeToken(token) {
  const payloadBase64 = token.split(".")[1]
  const payloadJson = atob(payloadBase64.replace(/-/g, "+").replace(/_/g, "/"))
  return JSON.parse(payloadJson)
}

export async function getAccount(token) {
  const res = await fetch(`${API_BASE}/accounts/me`, {
    headers: { Authorization: `Bearer ${token}` },
  })
  const data = await res.json()
  if (!res.ok) throw new Error(data.error)
  return data
}

export async function deposit(token, amount) {
  const res = await fetch(`${API_BASE}/deposit`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      Authorization: `Bearer ${token}`,
    },
    body: JSON.stringify({ amount: Number(amount) }),
  })
  const data = await res.json()
  if (!res.ok) throw new Error(data.error)
  return data
}

export async function transfer(token, toAccountNo, amount) {
  const res = await fetch(`${API_BASE}/transfer`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      Authorization: `Bearer ${token}`,
    },
    body: JSON.stringify({ to_account_no: toAccountNo, amount: Number(amount) }),
  })
  const data = await res.json()
  if (!res.ok) throw new Error(data.error)
  return data
}

export async function bulkTransfer(token, transfers) {
  const res = await fetch(`${API_BASE}/transfer/bulk`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      Authorization: `Bearer ${token}`,
    },
    body: JSON.stringify({
      transfers: transfers.map((t) => ({
        to_account_no: t.toAccountNo,
        amount: Number(t.amount),
      })),
    }),
  })
  const data = await res.json()
  if (!res.ok) throw new Error(data.error)
  return data
}

export async function getTransactions(token) {
  const res = await fetch(`${API_BASE}/transactions`, {
    headers: { Authorization: `Bearer ${token}` },
  })
  const data = await res.json()
  if (!res.ok) throw new Error(data.error)
  return data
}

export async function searchTransaction(token, code) {
  const res = await fetch(`${API_BASE}/transactions/${code}`, {
    headers: { Authorization: `Bearer ${token}` },
  })
  const data = await res.json()
  if (!res.ok) throw new Error(data.error)
  return data
}