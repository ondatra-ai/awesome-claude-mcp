import { render, screen } from '@testing-library/react';
import { VersionDisplay } from '@/components/version-display';
import type { IVersionDisplayProps } from '@/interfaces/version-display';

describe('VersionDisplay', () => {
  it('displays loading state', () => {
    const props: IVersionDisplayProps = {
      version: null,
      loadingState: {
        isLoading: true,
        error: null,
      },
    };

    render(<VersionDisplay {...props} />);

    expect(screen.getByTestId('backend-version')).toHaveTextContent(
      'Loading...'
    );
    expect(screen.getByText('Backend Version:')).toBeInTheDocument();
  });

  it('displays version when loaded successfully', () => {
    const props: IVersionDisplayProps = {
      version: { version: '1.0.0' },
      loadingState: {
        isLoading: false,
        error: null,
      },
    };

    render(<VersionDisplay {...props} />);

    expect(screen.getByTestId('backend-version')).toHaveTextContent('1.0.0');
    expect(screen.getByText('Backend Version:')).toBeInTheDocument();
  });

  it('displays error state', () => {
    const props: IVersionDisplayProps = {
      version: null,
      loadingState: {
        isLoading: false,
        error: 'Failed to fetch version',
      },
    };

    render(<VersionDisplay {...props} />);

    expect(screen.getByTestId('backend-version')).toHaveTextContent(
      'Error: Failed to fetch version'
    );
    expect(screen.getByText('Backend Version:')).toBeInTheDocument();
  });

  it('handles null version gracefully', () => {
    const props: IVersionDisplayProps = {
      version: null,
      loadingState: {
        isLoading: false,
        error: null,
      },
    };

    render(<VersionDisplay {...props} />);

    // Should not crash and should handle null version
    expect(screen.getByText('Backend Version:')).toBeInTheDocument();
    expect(screen.getByTestId('backend-version')).toBeInTheDocument();
  });
});
