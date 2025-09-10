import { render, screen } from '@testing-library/react';
import { WelcomeCard } from '@/components/welcome-card';

describe('WelcomeCard', () => {
  it('renders welcome title', () => {
    render(<WelcomeCard />);

    expect(screen.getByTestId('welcome-title')).toHaveTextContent(
      'Welcome to MCP Google Docs Editor'
    );
  });

  it('renders card description', () => {
    render(<WelcomeCard />);

    expect(
      screen.getByText(
        'This application provides a streamlined interface for document editing operations powered by the Model Context Protocol.'
      )
    ).toBeInTheDocument();
  });

  it('renders document operations feature', () => {
    render(<WelcomeCard />);

    expect(screen.getByTestId('feature-document-ops')).toHaveTextContent(
      'Document Operations'
    );
    expect(screen.getByTestId('feature-document-ops-desc')).toHaveTextContent(
      'Replace, append, prepend, and insert content in Google Docs'
    );
  });

  it('renders AI integration feature', () => {
    render(<WelcomeCard />);

    expect(screen.getByTestId('feature-ai-integration')).toHaveTextContent(
      'AI Integration'
    );
    expect(screen.getByTestId('feature-ai-integration-desc')).toHaveTextContent(
      'Compatible with Claude Code and ChatGPT via MCP protocol'
    );
  });

  it('has proper structure and styling classes', () => {
    const { container } = render(<WelcomeCard />);

    // Check for key structural elements
    expect(container.querySelector('.mx-auto.max-w-2xl')).toBeInTheDocument();
    expect(
      container.querySelector('.grid.grid-cols-1.gap-4.sm\\:grid-cols-2')
    ).toBeInTheDocument();
  });
});
