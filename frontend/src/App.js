import React, { useState, useEffect } from 'react';
import axios from 'axios';
import './App.css';

function EmployeeNode({ employee, reports }) {
  return (
    <li>
      <div className="employee-card">
        <span className="employee-title">{employee.title}:</span>
        <span className="employee-name">{employee.name}</span>
      </div>
      {reports && reports.length > 0 && (
        <ul>
          {reports.map(({ employee: report, reports: subReports }) => (
            <EmployeeNode 
              key={report.id} 
              employee={report} 
              reports={subReports}
            />
          ))}
        </ul>
      )}
    </li>
  );
}

function App() {
  const [tree, setTree] = useState(null);
  const [error, setError] = useState(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const fetchData = async () => {
      try {
        const response = await axios.get('http://localhost:8080/api/employees');
        // The API returns an array, we'll take the first item as the root
        setTree(response.data[0]);
      } catch (err) {
        console.error('API Error:', err);
        setError('Failed to fetch employee data. Make sure the backend server is running at http://localhost:8080');
      } finally {
        setLoading(false);
      }
    };

    fetchData();
  }, []);

  if (loading) {
    return <div className="loading">Loading employee data...</div>;
  }

  if (error) {
    return <div className="error">{error}</div>;
  }

  if (!tree) {
    return <div className="error">No employee data available</div>;
  }

  return (
    <div className="App">
      <header>
        <h1>Employee Organization Tree</h1>
      </header>
      <main>
        <ul className="tree">
          <EmployeeNode 
            employee={tree.employee} 
            reports={tree.reports} 
          />
        </ul>
      </main>
    </div>
  );
}

export default App; 