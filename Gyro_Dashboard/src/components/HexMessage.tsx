interface HexMessageProps {
    message?: string[]
}

function HexMessage({ message }: HexMessageProps) {
  return (
    <div className="Message bg-dark text-light" style={{ padding: "10px", marginTop: "10px" }}>
      <h2>Hex Message</h2>
        <div style={{ display: "flex", flexDirection: "column" }}>
            {message?.map((msg, index) => (
            <div key={index} style={{ marginBottom: "10px", color: "grey" }}>
                {msg}
            </div>
            ))}
        </div>
    </div>
  )
}

export default HexMessage