import { useState } from "react"
import { register } from "../api"

function Register({ onSwitchToLogin }) {
  const [username, setUsername] = useState("")
  const [password, setPassword] = useState("")
  const [fullName, setFullName] = useState("")
  const [message, setMessage] = useState("")

  async function handleSubmit(e) {
    e.preventDefault()
    try {
      const data = await register(username, password, fullName)
      setMessage(`Đăng ký thành công! account_id: ${data.account_id}`)
    } catch (err) {
      setMessage(`Lỗi: ${err.message}`)
    }
  }

  return (
    <div className="card">
      <h2>Đăng ký</h2>
      <form onSubmit={handleSubmit}>
        <input placeholder="Username" value={username} onChange={(e) => setUsername(e.target.value)} />
        <input placeholder="Password" type="password" value={password} onChange={(e) => setPassword(e.target.value)} />
        <input placeholder="Họ tên" value={fullName} onChange={(e) => setFullName(e.target.value)} />
        <button type="submit">Đăng ký</button>
      </form>
      {message && <p>{message}</p>}
      <button onClick={onSwitchToLogin}>Đã có tài khoản? Đăng nhập</button>
    </div>
  )
}

export default Register
