import { useState, useEffect } from "react"
import { decodeToken, getAccount, transfer, getTransactions, deposit, searchTransaction } from "../api"

function Dashboard({ token, onLogout, onGoToBulkTransfer }) {
  const { user_id: username } = decodeToken(token)

  const [account, setAccount] = useState(null)
  const [transactions, setTransactions] = useState([])
  const [toAccountNo, setToAccountNo] = useState("")
  const [amount, setAmount] = useState("")
  const [message, setMessage] = useState("")
  const [depositAmount, setDepositAmount] = useState("")
  const [depositMessage, setDepositMessage] = useState("")
  const [searchCode, setSearchCode] = useState("")
  const [searchedTx, setSearchedTx] = useState(undefined) // undefined = chưa tìm gì; null = đã tìm nhưng không thấy; object = tìm thấy
  const [searchError, setSearchError] = useState("")

  async function loadData() {
    const acc = await getAccount(token)
    setAccount(acc)
    const txs = await getTransactions(token)
    setTransactions(txs || [])
  }

  useEffect(() => {
    loadData()
  }, [])

  async function handleDeposit(e) {
    e.preventDefault()
    try {
      await deposit(token, depositAmount)
      setDepositMessage("Nạp tiền thành công (demo)")
      setDepositAmount("")
      await loadData()
    } catch (err) {
      setDepositMessage(`Lỗi: ${err.message}`)
    }
  }

  async function handleSearch(e) {
    e.preventDefault()
    setSearchError("")
    try {
      const tx = await searchTransaction(token, searchCode)
      setSearchedTx(tx)
    } catch (err) {
      setSearchedTx(null) // đã tìm nhưng không thấy — khác với "chưa tìm gì" (undefined)
      setSearchError(err.message)
    }
  }

  function clearSearch() {
    setSearchCode("")
    setSearchedTx(undefined)
    setSearchError("")
  }

  async function handleTransfer(e) {
    e.preventDefault()
    try {
      const trx = await transfer(token, toAccountNo, amount)
      setMessage(`Chuyển tiền thành công — mã giao dịch ${trx.tx_code}`)
      setToAccountNo("")
      setAmount("")
      clearSearch()
      await loadData()
    } catch (err) {
      setMessage(`Lỗi: ${err.message}`)
    }
  }

  return (
    <div>
      <h2>Xin chào, {username}</h2>
      <button onClick={onLogout}>Đăng xuất</button>

      {account && (
        <p>
          {account.full_name}
          <br />
          Số tài khoản: <b>{account.account_no}</b>
          <br />
          Số dư: {account.balance.toLocaleString()}đ
        </p>
      )}

      <h3>Nạp tiền</h3>
      <p style={{ fontSize: "13px", opacity: 0.7 }}>
        (DEMO — mô phỏng nộp tiền mặt tại quầy ngân hàng, không kết nối cổng thanh toán thật)
      </p>
      <form onSubmit={handleDeposit}>
        <input placeholder="Số tiền nạp" value={depositAmount} onChange={(e) => setDepositAmount(e.target.value)} />
        <button type="submit">Nạp tiền</button>
      </form>
      {depositMessage && <p>{depositMessage}</p>}

      <h3>Chuyển tiền</h3>
      <form onSubmit={handleTransfer}>
        <input placeholder="Số tài khoản nhận" value={toAccountNo} onChange={(e) => setToAccountNo(e.target.value)} />
        <input placeholder="Số tiền" value={amount} onChange={(e) => setAmount(e.target.value)} />
        <button type="submit">Chuyển</button>
      </form>
      {message && <p>{message}</p>}
      <button onClick={onGoToBulkTransfer}>Chuyển tiền hàng loạt →</button>
      <button
        onClick={() => {
          clearSearch()
          loadData()
        }}
      >
        Làm mới số dư / lịch sử
      </button>

      <h3>Lịch sử giao dịch</h3>
      <form onSubmit={handleSearch}>
        <input placeholder="Tìm theo mã giao dịch" value={searchCode} onChange={(e) => setSearchCode(e.target.value)} />
        <button type="submit">Tìm</button>
      </form>
      {searchError && <p>{searchError}</p>}
      {searchedTx !== undefined && <button onClick={clearSearch}>Xem tất cả lịch sử</button>}

      <table border="1">
        <thead>
          <tr>
            <th>Mã GD</th>
            <th>Từ</th>
            <th>Đến</th>
            <th>Số tiền</th>
            <th>Trạng thái</th>
            <th>Thời gian</th>
          </tr>
        </thead>
        <tbody>
          {(searchedTx === undefined ? transactions : searchedTx ? [searchedTx] : []).map((tx) => (
            <tr key={tx.id}>
              <td>{tx.tx_code}</td>
              <td>{tx.from_acc_id}</td>
              <td>{tx.to_acc_id}</td>
              <td>{tx.amount.toLocaleString()}đ</td>
              <td>{tx.status}</td>
              <td>{new Date(tx.created_at).toLocaleString()}</td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  )
}

export default Dashboard
