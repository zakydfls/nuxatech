const API_BASE = "http://localhost:9000/api/v1/user/wallet";
const WALLET_ID = "3441a239-b1de-4c45-adfb-79fd919f7ab5";
const TOKEN =
  "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6Inpha3lAbWFpbC5jb20iLCJleHBpcmVkX2F0IjoxNzQ0NDU4OTY0LCJpc3N1ZWRfYXQiOjE3NDQ0NDgxNjQsInVzZXJfaWQiOiI1NDczMTdhNy01ODhlLTQxYWQtYmNhNi0yZWU4YWMzZWE0YjIiLCJ1c2VybmFtZSI6Inpha3lkZmxzIn0.wHVASAVEASFxWPQcwtVBH3_ouQovoVZhfl_xJ5eyfxM"; // your token

async function testSmallDeposits() {
  const deposit = async (amount) => {
    const response = await fetch(`${API_BASE}/${WALLET_ID}/deposit`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        Authorization: TOKEN,
      },
      body: JSON.stringify({ amount }),
    });
    return response.json();
  };

  // Get initial balance
  const initialResp = await fetch(`${API_BASE}/${WALLET_ID}`, {
    headers: { Authorization: TOKEN },
  });
  const {
    data: { balance: initialBalance },
  } = await initialResp.json();
  console.log(`Initial balance: ${initialBalance}`);

  // Perform 100 concurrent small deposits
  const amount = 1000; // 1,000 each
  const count = 100;
  const deposits = Array(count)
    .fill()
    .map(() => deposit(amount));

  const results = await Promise.all(deposits);
  const successful = results.filter((r) => r.success).length;

  // Expected final balance
  const expectedBalance = 3900000;
  console.log(`Expected final balance: ${expectedBalance}`);

  console.log("\nWaiting for 5 seconds before checking final balance...");
  await new Promise((resolve) => setTimeout(resolve, 5000));

  // Get actual final balance
  const finalResp = await fetch(`${API_BASE}/${WALLET_ID}`, {
    headers: { Authorization: TOKEN },
  });
  const {
    data: { balance: finalBalance },
  } = await finalResp.json();
  console.log(`Actual final balance: ${finalBalance}`);
  console.log(`Difference: ${finalBalance - expectedBalance}`);
}

testSmallDeposits();
