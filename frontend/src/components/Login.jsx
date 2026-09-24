import { useState } from "react"
import { login } from "../api"

function Login({ onLoginSuccess, onSwitchToRegister }) {
  const [username, setUsername] = useState("")
  const [password, setPassword] = useState("")
  const [error, setError] = useState("")

  async function handleSubmit(e) {
    e.preventDefault()
    try {
      const token = await login(username, password)
      onLoginSuccess(token)
    } catch (err) {
      setError(err.message)
    }
  }

  return (
    <div className="card">
      <h2>Đăng nhập</h2>
      <form onSubmit={handleSubmit}>
        <input placeholder="Username" value={username} onChange={(e) => setUsername(e.target.value)} />
        <input placeholder="Password" type="password" value={password} onChange={(e) => setPassword(e.target.value)} />
        <button type="submit">Đăng nhập</button>
      </form>
      {error && <p>{error}</p>}
      <button onClick={onSwitchToRegister}>Chưa có tài khoản? Đăng ký</button>
    </div>
  )
}

export default Login
