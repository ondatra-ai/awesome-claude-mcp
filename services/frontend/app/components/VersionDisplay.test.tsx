import { render, screen, waitFor } from '@testing-library/react';
import '@testing-library/jest-dom';
import VersionDisplay from './VersionDisplay';

// Mock fetch
global.fetch = jest.fn();

describe('VersionDisplay', () => {
  beforeEach(() => {
    (fetch as jest.Mock).mockClear();
  });

  it('displays loading state initially', () => {
    (fetch as jest.Mock).mockResolvedValue({
      ok: true,
      json: async () => ({ version: '1.0.0' }),
    });

    render(<VersionDisplay />);
    expect(screen.getByText('Loading version...')).toBeInTheDocument();
  });

  it('displays backend version when fetch succeeds', async () => {
    (fetch as jest.Mock).mockResolvedValue({
      ok: true,
      json: async () => ({ version: '1.0.0' }),
    });

    render(<VersionDisplay />);
    
    await waitFor(() => {
      expect(screen.getByText('Backend version: 1.0.0')).toBeInTheDocument();
    });
  });

  it('displays error message when fetch fails', async () => {
    (fetch as jest.Mock).mockRejectedValue(new Error('Network error'));

    render(<VersionDisplay />);
    
    await waitFor(() => {
      expect(screen.getByText('Error loading version')).toBeInTheDocument();
    });
  });

  it('uses correct API URL from environment', async () => {
    process.env.NEXT_PUBLIC_API_URL = 'http://test-api:8080';
    
    (fetch as jest.Mock).mockResolvedValue({
      ok: true,
      json: async () => ({ version: '1.0.0' }),
    });

    render(<VersionDisplay />);
    
    await waitFor(() => {
      expect(fetch).toHaveBeenCalledWith('http://test-api:8080/version');
    });
  });
});