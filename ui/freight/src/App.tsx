import { useState, useEffect } from 'react';
import './App.css';

function App() {
  const [code, setCode] = useState("");
  const [loading, setLoading] = useState(false);
  const [jobId, setJobId] = useState("");
  const [output, setOutput] = useState("");
  const [language, setLanguage] = useState("py"); // Default language


  function handleSubmit() {
    setLoading(true);
    setOutput("...")

    fetch("http://localhost:30000/submission", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({
        lang: language,
        code: code,
        job_id: "",
      }),
    })
      .then((response) => response.json())
      .then((data) => {
        console.log("Success:", data);
        setJobId(data["job_id"]);
      })
      .catch((error) => {
        console.error("Error:", error);
      });
  }

  useEffect(() => {
    let interval = 1000;
    if (jobId) {
      interval = setInterval(() => {
        fetch(`http://localhost:30000/check?job_id=${jobId}`)
          .then((response) => response.json())
          .then((data) => {
            if (data['state'] === "SUCCESS" || data.status === "FAILED") {
              setOutput(data.output);
              clearInterval(interval); // Stop polling
              setLoading(false);
            }
          })
          .catch((error) => console.error("Error:", error));
      }, 1000); // Poll every second
    }

    return () => clearInterval(interval); // Cleanup on unmount or jobId change
  }, [jobId]);

  return (
    <>
      <div>
        <h1 style={{color: "gold"}}>🍝 Spaghetti code execution engine</h1>
        <div className="input-container">
        <textarea className="input-box main-input" rows={5}onChange={(event) => setCode(event.target.value)}></textarea>
        <textarea className="input-box side-input" rows={5} value={output} readOnly></textarea>

        </div>
        <br />

        <select value={language} onChange={(event) => setLanguage(event.target.value)}>
          <option value="c">C</option>
          <option value="cpp">C++</option>
          <option value="py">Python</option>
          <option value="go">Go</option>
          <option value="js">JavaScript</option>
        </select>


        {loading ? (
          <button className='rotate-button'>Submit</button>
        ) : (
          <button onClick={handleSubmit}>Submit</button>
        )}
      </div>
    </>
  );
}

export default App;
