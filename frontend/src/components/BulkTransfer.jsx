import { useState, useEffect } from "react"
import { bulkTransfer, getAccount } from "../api"

function BulkTransfer({ token, onBack }) {
  const [account, setAccount] = useState(null)
  const [rows, setRows] = useState([
    { toAccountNo: "", amount: "" },
    { toAccountNo: "", amount: "" },
  ])
  const [result, setResult] = useState(null)
  const [error, setError] = useState("")

  useEffect(() => {
    getAccount(token).then(setAccount)
  }, [])

  function updateRow(index, field, value) {
    const newRows = [...rows]
    newRows[index] = { ...newRows[index], [field]: value }
    setRows(newRows)
  }

  function addRow() {
    setRows([...rows, { toAccountNo: "", amount: "" }])
  }

  function removeRow(index) {
    setRows(rows.filter((_, i) => i !== index))
  }

  async function handleSubmit(e) {
    e.preventDefault()
    setError("")
    setResult(null)
    try {
      const data = await bulkTransfer(token, rows)
      setResult(data)
      getAccount(token).then(setAccount)
    } catch (err) {
      setError(err.message)
    }
  }

  return (
    <div>
      <button onClick={onBack}>← Quay lại Dashboard</button>
      <h2>Chuyển tiền hàng loạt</h2>
      {account && <p>Số dư hiện tại: {account.balance.toLocaleString()}đ</p>}

      <form onSubmit={handleSubmit}>
        <table border="1">
          <thead>
            <tr>
              <th>Số tài khoản nhận</th>
              <th>Số tiền</th>
              <th></th>
            </tr>
          </thead>
          <tbody>
            {rows.map((row, index) => (
              <tr key={index}>
                <td>
                  <input
                    value={row.toAccountNo}
                    onChange={(e) => updateRow(index, "toAccountNo", e.target.value)}
                  />
                </td>
                <td>
                  <input value={row.amount} onChange={(e) => updateRow(index, "amount", e.target.value)} />
                </td>
                <td>
                  <button type="button" onClick={() => removeRow(index)}>
                    Xóa
                  </button>
                </td>
              </tr>
            ))}
          </tbody>
        </table>

        <button type="button" onClick={addRow}>
          + Thêm dòng
        </button>
        <button type="submit">Gửi tất cả</button>
      </form>

      {error && <p>{error}</p>}

      {result && (
        <div>
          <h3>
            Kết quả: {result.success}/{result.total} thành công
          </h3>
          <table border="1">
            <thead>
              <tr>
                <th>Số TK nhận</th>
                <th>Số tiền</th>
                <th>Trạng thái</th>
                <th>Chi tiết</th>
              </tr>
            </thead>
            <tbody>
              {result.results.map((r, index) => (
                <tr key={index}>
                  <td>{r.to_account_no}</td>
                  <td>{r.amount.toLocaleString()}đ</td>
                  <td>{r.status === "success" ? "✅ thành công" : "❌ thất bại"}</td>
                  <td>{r.status === "success" ? r.tx_code : r.error}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}
    </div>
  )
}

export default BulkTransfer
