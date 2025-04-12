const API_BASE = "http://localhost:9000/api/v1/user/wallet";
const WALLET_ID = "3441a239-b1de-4c45-adfb-79fd919f7ab5";
const TOKEN =
  "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6Inpha3lAbWFpbC5jb20iLCJleHBpcmVkX2F0IjoxNzQ0NDU4OTY0LCJpc3N1ZWRfYXQiOjE3NDQ0NDgxNjQsInVzZXJfaWQiOiI1NDczMTdhNy01ODhlLTQxYWQtYmNhNi0yZWU4YWMzZWE0YjIiLCJ1c2VybmFtZSI6Inpha3lkZmxzIn0.wHVASAVEASFxWPQcwtVBH3_ouQovoVZhfl_xJ5eyfxM"; // your token

async function testLargeDeposits() {
  const deposit = async (amount, index) => {
    console.log(`Starting deposit ${index + 1}: ${amount}`);
    const response = await fetch(`${API_BASE}/${WALLET_ID}/deposit`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        Authorization: TOKEN,
      },
      body: JSON.stringify({ amount }),
    });
    const result = await response.json();
    console.log(`Completed deposit ${index + 1}: ${result.success ? "Success" : "Failed"}`);
    return result;
  };

  // Get initial balance
  const initialResp = await fetch(`${API_BASE}/${WALLET_ID}`, {
    headers: { Authorization: TOKEN },
  });
  const {
    data: { balance: initialBalance },
  } = await initialResp.json();
  console.log(`Initial balance: ${initialBalance}`);

  // Different amounts for each deposit
  const amounts = [
    100000, // 100K
    200000, // 200K
    500000, // 500K
    1000000, // 1M
    2000000, // 2M
  ];

  const deposits = amounts.map((amount, index) => deposit(amount, index));
  const results = await Promise.all(deposits);

  // Calculate expected balance
  const successfulTotal = results.filter((r) => r.success).reduce((sum, r) => sum + r.data.amount, 0);
  const expectedBalance = 3900000;

  console.log("\nResults Summary:");
  console.log(`Successful deposits: ${results.filter((r) => r.success).length}`);
  console.log(`Failed deposits: ${results.filter((r) => !r.success).length}`);
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

testLargeDeposits();
