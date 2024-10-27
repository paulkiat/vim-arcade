// adaptive-ui/pages/index.js
import React, { useEffect, useState } from 'react';
import axios from 'axios';

const Home = ({ initialData }) => {
  const [data, setData] = useState(initialData);

  useEffect(() => {
    const interval = setInterval(async () => {
      try {
        const response = await axios.get('/api/heartbeat', {
          headers: {
            'Accept-Encoding': 'gzip, zstd'
          },
          responseType: 'arraybuffer' // Necessary for binary data
        });

        let decompressedData;
        const encoding = response.headers['content-encoding'];

        if (encoding === 'gzip') {
          const decompressed = await decompressGzip(response.data);
          decompressedData = decompressed;
        } else if (encoding === 'zstd') {
          const decompressed = await decompressZstd(response.data);
          decompressedData = decompressed;
        } else {
          decompressedData = JSON.parse(new TextDecoder().decode(response.data));
        }

        setData(decompressedData);
      } catch (error) {
        console.error('Heartbeat failed:', error);
      }
    }, 5000); // Heartbeat every 5 seconds

    return () => clearInterval(interval);
  }, []);

  // Placeholder decompression functions
  async function decompressGzip(data) {
    // Implement gzip decompression on client-side if necessary
    // Typically handled automatically by browsers
    return JSON.parse(new TextDecoder().decode(data));
  }

  async function decompressZstd(data) {
    // Implement zstd decompression on client-side
    // Requires a zstd library for JavaScript (e.g., pako, if available)
    // Alternatively, handle decompression server-side
    // For simplicity, assume server sends JSON if decompression is not implemented
    return JSON.parse(new TextDecoder().decode(data));
  }

  return (
    <div>
      <h1>AdaptiveUI Matchmaking</h1>
      <pre>{JSON.stringify(data, null, 2)}</pre>
    </div>
  );
};

export async function getServerSideProps(context) {
  const res = await fetch('http://localhost:4000/api/initial');
  const initialData = await res.json();

  return { props: { initialData } };
}

export default Home;