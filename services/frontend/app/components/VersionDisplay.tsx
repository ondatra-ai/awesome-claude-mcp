'use client';

import React, { useState, useEffect } from 'react';

export default function VersionDisplay() {
  const [version, setVersion] = useState<string>('');
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string>('');

  useEffect(() => {
    const fetchVersion = async () => {
      try {
        const apiUrl = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';
        const response = await fetch(`${apiUrl}/version`);
        
        if (!response.ok) {
          throw new Error('Failed to fetch version');
        }
        
        const data = await response.json();
        setVersion(data.version);
      } catch (err) {
        setError('Error loading version');
        console.error('Failed to fetch version:', err);
      } finally {
        setLoading(false);
      }
    };

    fetchVersion();
  }, []);

  if (loading) {
    return <p className="text-sm text-gray-600">Loading version...</p>;
  }

  if (error) {
    return <p className="text-sm text-red-600">{error}</p>;
  }

  return (
    <p className="text-sm text-gray-600">
      Backend version: {version}
    </p>
  );
}