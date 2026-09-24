import { useState } from "react"
import Login from "./components/Login"
import Register from "./components/Register"
import Dashboard from "./components/Dashboard"
import BulkTransfer from "./components/BulkTransfer"

function App() {
  const [token, setToken] = useState(localStorage.getItem("token"))
  const [view, setView] = useState("login")
  const [page, setPage] = useState("dashboard") // trang sau khi đăng nhập: "dashboard" | "bulk-transfer"

  function handleLoginSuccess(newToken) {
    localStorage.setItem("token", newToken)
    setToken(newToken)
  }

  function handleLogout() {
    localStorage.removeItem("token")
    setToken(null)
  }

  if (token) {
    if (page === "bulk-transfer") {
      return <BulkTransfer token={token} onBack={() => setPage("dashboard")} />
    }
    return (
      <Dashboard token={token} onLogout={handleLogout} onGoToBulkTransfer={() => setPage("bulk-transfer")} />
    )
  }

  if (view === "register") {
    return <Register onSwitchToLogin={() => setView("login")} />
  }

  return (
    <Login
      onLoginSuccess={handleLoginSuccess}
      onSwitchToRegister={() => setView("register")}
    />
  )
}

export default App
