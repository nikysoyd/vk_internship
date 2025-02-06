import React, { useEffect, useState } from 'react';
import axios from 'axios';

interface PingResult {
  ip: string;
  ping_time: string;
  is_success: boolean;
  last_success: string;
}

const App: React.FC = () => {
  const [results, setResults] = useState<PingResult[]>([]);

  useEffect(() => {
    const fetchData = async () => {
      try {
        const response = await axios.get('http://backend:8080/ping-results');
        setResults(response.data);
      } catch (error) {
        console.error('Error fetching data:', error);
      }
    };

    fetchData();
    const interval = setInterval(fetchData, 5000); // Обновление каждые 5 секунд
    return () => clearInterval(interval);
  }, []);

  return (
    <div>
      <h1>Docker Container Ping Results</h1>
      <table>
        <thead>
          <tr>
            <th>IP Address</th>
            <th>Ping Time</th>
            <th>Last Success</th>
          </tr>
        </thead>
        <tbody>
          {results.map((result, index) => (
            <tr key={index}>
              <td>{result.ip}</td>
              <td>{result.ping_time}</td>
              <td>{result.last_success}</td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
};

export default App;
